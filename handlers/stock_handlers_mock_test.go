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
	"github.com/stretchr/testify/mock"
)

type mockHistoryAPI struct {
	mock.Mock
}

func (mock *mockHistoryAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	args := mock.Called(symbol, start, end)
	stocks := args.Get(0).([]finance.Stock)
	return stocks, args.Error(1)
}

type mockEsStock struct {
	es.Stock
	stockAggs map[string][]es.StocksAgg
}

func (mock *mockEsStock) Index(stock finance.Stock) error {
	return nil
}

func (mock *mockEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]es.StocksAgg, error) {
	return mock.stockAggs[symbol], nil
}

type DummySchemaDecoder struct {
}

func (decoder *DummySchemaDecoder) Decode(dst interface{}, src map[string][]string) error {
	if params, ok := dst.(*HistoryListParams); ok {
		params.Symbols = append(params.Symbols, "TEST")
	}
	return nil
}

type ErrorSchemaDecoder struct {
	Msg string
}

func (decoder *ErrorSchemaDecoder) Decode(dst interface{}, src map[string][]string) error {
	return errors.New(decoder.Msg)
}

type ErrorFinanceAPI struct {
	Msg string
}

func (api *ErrorFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return nil, errors.New(api.Msg)
}

type DummyFinanceAPI struct {
}

func (api *DummyFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return []finance.Stock{}, nil
}

type OneItemFinanceAPI struct {
}

func (api *OneItemFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return []finance.Stock{{Symbol: "TEST"}}, nil
}

type esStockIndexError struct {
	es.Stock
	Msg string
}

func (mock *esStockIndexError) Index(stock finance.Stock) error {
	return errors.New(mock.Msg)
}

type esStockGetStockAggError struct {
	es.Stock
	Msg string
}

func (mock *esStockGetStockAggError) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]es.StocksAgg, error) {
	return nil, errors.New(mock.Msg)
}
