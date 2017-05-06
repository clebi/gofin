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
	"errors"
	"time"

	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
)

type IndicatorSchemaDecoder struct {
	Symbols []string
}

func (decoder *IndicatorSchemaDecoder) Decode(dst interface{}, src map[string][]string) error {
	if indicator, ok := dst.(*getStocksParams); ok {
		indicator.Symbols = decoder.Symbols
	} else {
		return errors.New("bad type for IndicatorSchemaDecoder")
	}
	return nil
}

type IndicatorQuotesAPI struct {
	quotes map[string]*finance.Quote
}

func (api *IndicatorQuotesAPI) GetQuote(symbol string) (*finance.Quote, error) {
	return api.quotes[symbol], nil
}

type IndicatorTestEsStock struct {
	index int
	stats []es.StocksStats
}

func (mock *IndicatorTestEsStock) Index(stock finance.Stock) error {
	return nil
}

func (mock *IndicatorTestEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]es.StocksAgg, error) {
	return nil, nil
}

func (mock *IndicatorTestEsStock) GetStockStats(symbol string, startDate time.Time, endDate time.Time) (*es.StocksStats, error) {
	stats := &mock.stats[mock.index]
	mock.index++
	return stats, nil
}

type IndicatorErrorEsStock struct {
	index int
	errs  []error
}

func (mock *IndicatorErrorEsStock) Index(stock finance.Stock) error {
	return nil
}

func (mock *IndicatorErrorEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]es.StocksAgg, error) {
	return nil, nil
}

func (mock *IndicatorErrorEsStock) GetStockStats(symbol string, startDate time.Time, endDate time.Time) (*es.StocksStats, error) {
	err := mock.errs[mock.index]
	mock.index++
	return nil, err
}
