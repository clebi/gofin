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

package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-playground/validator"
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

// GetDateFunc is responsible for get end date
type GetDateFunc func() time.Time

func getYesterDayDate() time.Time {
	return time.Now().AddDate(0, 0, -1)
}

// StructValidator validates structures
type StructValidator interface {
	Struct(s interface{}) (err error)
}

type errorHandlerFunc func(c echo.Context, status int, err error) error

// StockHandlers is an object containing all the handlers concerning stocks
type StockHandlers struct {
	*Context
	validator    StructValidator
	getDate      GetDateFunc
	errorHandler errorHandlerFunc
}

// NewStockHandlers creates a new stock handlers object
func NewStockHandlers(context *Context) *StockHandlers {
	return &StockHandlers{
		Context:      context,
		validator:    validator.New(),
		getDate:      getYesterDayDate,
		errorHandler: handleError,
	}
}

// handleError writes error to the http channel and logs internal errors
func handleError(c echo.Context, status int, err error) error {
	var msg string
	if status == http.StatusInternalServerError {
		log.Error(err)
		msg = "unknown_error"
	} else {
		msg = err.Error()
	}
	c.JSON(status, errorDesc{Status: "error", Description: msg})
	return nil
}

// History retrieve stocks history
//
//  handlers.History(w, r)
//
// This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) History(c echo.Context) error {
	var params HistoryParams
	if err := handlers.sh.Decode(&params, c.Request().URL.Query()); err != nil {
		return handlers.errorHandler(c, http.StatusInternalServerError, err)
	}
	if err := handlers.validator.Struct(params); err != nil {
		return handlers.errorHandler(c, http.StatusBadRequest, err)
	}
	end := handlers.getDate()
	end = end.Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, params.Days*-1)
	stocks, err := handlers.Context.historyAPI.GetHistory(c.Param("symbol"), start.AddDate(0, 0, params.Window*-1), end)
	if err != nil {
		return handlers.errorHandler(c, http.StatusBadRequest, err)
	}
	for _, stock := range stocks {
		err = handlers.Context.esStock.Index(stock)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
	}
	stocksAgg, err := handlers.Context.esStock.GetStocksAgg(c.Param("symbol"), params.Window, params.Step, start, end)
	if err != nil {
		return handlers.errorHandler(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, stocksAgg)
}
