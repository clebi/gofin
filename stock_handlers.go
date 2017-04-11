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
	"bytes"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/clebi/yfinance"
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

type errorHandlerFunc func(c echo.Context)

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
		errorHandler: handleErrors,
	}
}

func handleErrors(c echo.Context) {
	if err := recover(); err != nil {
		switch err := err.(type) {
		case finance.YApiError:
			c.JSON(http.StatusBadRequest, errorDesc{Status: "error", Description: err.Error()})
		case validator.ValidationErrors:
			errBuff := bytes.NewBufferString("following parameters are invalid: ")
			for _, err := range err {
				errBuff.WriteString(strings.ToLower(err.Field()))
				errBuff.WriteByte(',')
			}
			errBuff.Truncate(errBuff.Len() - 1)
			errorMsg := errBuff.String()
			c.JSON(http.StatusBadRequest, errorDesc{Status: "error", Description: errorMsg})
		default:
			log.Error(err)
			c.JSON(http.StatusInternalServerError, errorDesc{Status: "error", Description: "unknown_error"})
		}
	}
}

// History retrieve stocks history
//
//  handlers.History(w, r)
//
// This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) History(c echo.Context) error {
	defer handlers.errorHandler(c)
	var params HistoryParams
	if err := handlers.sh.Decode(&params, c.Request().URL.Query()); err != nil {
		panic(err)
	}
	if err := validator.New().Struct(params); err != nil {
		panic(err)
	}
	end := handlers.getDate()
	end = end.Truncate(24 * time.Hour)
	start := end.AddDate(0, 0, params.Days*-1)
	stocks, err := handlers.Context.historyAPI.GetHistory(c.Param("symbol"), start.AddDate(0, 0, params.Window*-1), end)
	if err != nil {
		panic(err)
	}
	for _, stock := range stocks {
		err = handlers.Context.esStock.Index(stock)
		if err != nil {
			panic(err)
		}
	}
	stocksAgg, err := handlers.Context.esStock.GetStocksAgg(c.Param("symbol"), params.Window, params.Step, start, end)
	if err != nil {
		panic(err)
	}
	return c.JSON(http.StatusOK, stocksAgg)
}
