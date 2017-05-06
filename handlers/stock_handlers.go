// Copyright 2017 Clément Bizeau
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"net/http"
	"time"

	"github.com/clebi/gofin/es"
	"github.com/labstack/echo"
)

type errorDesc struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// HistoryParams contains all the parameters for the history route
type HistoryParams struct {
	Days   int `schema:"days" validate:"gt=0"`
	Window int `schema:"window" validate:"gt=0"`
	Step   int `schema:"step" validate:"gt=0"`
}

// HistoryListParams contains all the parameters for the history list route
type HistoryListParams struct {
	HistoryParams
	Symbols []string `schema:"symbols"`
}

// HandlerERROR represents an error to send through http
type HandlerERROR struct {
	error
	Status int
}

// StructValidator validates structures
type StructValidator interface {
	Struct(s interface{}) (err error)
}

type errorHandlerFunc func(c echo.Context, status int, err error) error

// StockHandlers is an object containing all the handlers concerning stocks
type StockHandlers struct {
	*Context
	getDate      GetDateFunc
	errorHandler errorHandlerFunc
}

// NewStockHandlers creates a new stock handlers object
func NewStockHandlers(context *Context) *StockHandlers {
	return &StockHandlers{
		Context:      context,
		getDate:      getYesterDayDate,
		errorHandler: handleError,
	}
}

func (handlers *StockHandlers) getDates(from time.Time, days int) (time.Time, time.Time) {
	end := from.Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, days*-1)
	return start, end
}

// History retrieve stocks history
//
//  handlers.History(w, r)
//
// This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) History(c echo.Context) error {
	var params HistoryParams
	if handlerErr := getQuery(c, handlers.Context, &params); handlerErr != nil {
		return handlers.errorHandler(c, handlerErr.Status, handlerErr.error)
	}
	start, end := handlers.getDates(handlers.getDate(), params.Days)
	httpErr := indexStock(handlers.Context, c.Param("symbol"), start.AddDate(0, 0, params.Window*-1), end)
	if httpErr != nil {
		return handlers.errorHandler(c, httpErr.Status, httpErr.error)
	}
	stocksAgg, err := handlers.Context.esStock.GetStocksAgg(c.Param("symbol"), params.Window, params.Step, start, end)
	if err != nil {
		return handlers.errorHandler(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, stocksAgg)
}

// HistoryList retrieve a stock history list
//
//  handlers.History(w, r)
//
// This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) HistoryList(c echo.Context) error {
	var params HistoryListParams
	if handlerErr := getQuery(c, handlers.Context, &params); handlerErr != nil {
		return handlers.errorHandler(c, handlerErr.Status, handlerErr.error)
	}
	start, end := handlers.getDates(handlers.getDate(), params.Days)
	var stocks [][]es.StocksAgg
	for _, symbol := range params.Symbols {
		httpErr := indexStock(handlers.Context, symbol, start.AddDate(0, 0, params.Window*-1), end)
		if httpErr != nil {
			return handlers.errorHandler(c, httpErr.Status, httpErr.error)
		}
		stocksAgg, err := handlers.Context.esStock.GetStocksAgg(symbol, params.Window, params.Step, start, end)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
		stocks = append(stocks, stocksAgg)
	}
	return c.JSON(http.StatusOK, stocks)
}
