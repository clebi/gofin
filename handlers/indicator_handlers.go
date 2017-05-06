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

	"github.com/labstack/echo"
)

// Indicator contains all values of a stock indicator
type Indicator struct {
	Symbol   string
	Name     string
	Value    float32
	MM200    float32
	MM50     float32
	MM50D200 float32
	V50      float64
	V200     float64
}

type getStocksParams struct {
	Symbols []string `schema:"symbols"`
}

// IndicatorHandlers handles all request to avergaes requrests
type IndicatorHandlers struct {
	*Context
	getDate      GetDateFunc
	errorHandler errorHandlerFunc
	indexStock   indexStockFunc
}

// NewIndicatorHandlers creates a new averages handlers object
func NewIndicatorHandlers(context *Context) IndicatorHandlers {
	return IndicatorHandlers{
		Context:      context,
		getDate:      getYesterDayDate,
		errorHandler: handleError,
		indexStock:   indexStock,
	}
}

// GetStocks retrieves the indicators for a list of stocks
//
// This function is a handler for http server, it should not be called directly
func (handlers *IndicatorHandlers) GetStocks(c echo.Context) error {
	endDate := handlers.getDate()
	var params getStocksParams
	if handlerErr := getQuery(c, handlers.Context, &params); handlerErr != nil {
		return handlers.errorHandler(c, handlerErr.Status, handlerErr.error)
	}
	indicators := make([]Indicator, len(params.Symbols))
	for i, symbol := range params.Symbols {
		httpErr := handlers.indexStock(handlers.Context, symbol, endDate.AddDate(0, 0, -365), endDate)
		if httpErr != nil {
			return handlers.errorHandler(c, httpErr.Status, httpErr.error)
		}
		quote, err := handlers.quotesAPI.GetQuote(symbol)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
		stockStats200, err := handlers.esStock.GetStockStats(symbol, endDate.AddDate(0, 0, -200), endDate)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
		stockStats50, err := handlers.esStock.GetStockStats(symbol, endDate.AddDate(0, 0, -50), endDate)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
		indicators[i] = Indicator{
			Symbol:   quote.Symbol,
			Name:     quote.Name,
			Value:    quote.LastTradePriceOnly,
			MM200:    quote.TwoHundreddayMovingAverage,
			MM50:     quote.FiftydayMovingAverage,
			MM50D200: quote.FiftydayMovingAverage / quote.TwoHundreddayMovingAverage,
			V200:     stockStats200.StandardDeviation / stockStats200.Avg,
			V50:      stockStats50.StandardDeviation / stockStats50.Avg,
		}
	}
	return c.JSON(http.StatusOK, indicators)
}
