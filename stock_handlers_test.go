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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	finance "github.com/clebi/yfinance"
	schema "github.com/gorilla/Schema"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	symbolTest = "TEST"
)

type mockHistoryAPI struct {
	mock.Mock
}

func (mock *mockHistoryAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	args := mock.Called(symbol, start, end)
	stocks := args.Get(0).([]finance.Stock)
	return stocks, args.Error(1)
}

func getTestDate() time.Time {
	time, _ := time.Parse(time.RFC3339, "2016-12-13T01:24:23Z")
	return time
}

type mockEsStock struct {
	mock.Mock
}

func (mock *mockEsStock) Index(stock finance.Stock) error {
	args := mock.Called(stock)
	return args.Error(0)
}

func (mock *mockEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]EsStocksAgg, error) {
	args := mock.Called(symbol, movAvgWindow, step, startDate, endDate)
	stocks := args.Get(0).([]EsStocksAgg)
	return stocks, args.Error(1)
}

func TestHistory(t *testing.T) {
	testEndDate := getTestDate().Truncate(24 * time.Hour)
	testStartDate := testEndDate.AddDate(0, 0, -3)
	testStartMovDate := testStartDate.AddDate(0, 0, -2)
	stock := finance.Stock{Open: 1.1, High: 2.2, Low: 3.3, Close: 4.4, Volume: 999, Symbol: symbolTest, Date: finance.YTime{Time: testStartDate}}
	stocks := []finance.Stock{stock}
	mockedHistoryAPI := mockHistoryAPI{}
	mockedHistoryAPI.On("GetHistory", symbolTest, testStartMovDate, testEndDate).Return(stocks, nil)
	stocksAgg := []EsStocksAgg{{Symbol: symbolTest, MsTime: testStartMovDate.Unix() * 1000, AvgClose: 4.4, MovClose: 4.1}}
	stcoskAggJson, err := json.Marshal(stocksAgg)
	if err != nil {
		t.Fatal(err)
	}
	mockedEsStock := mockEsStock{}
	mockedEsStock.On("Index", stock).Return(nil)
	mockedEsStock.On("GetStocksAgg", symbolTest, 2, 2, testStartDate, testEndDate).Return(stocksAgg, nil)
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://test.test/graph/TEST?days=3&window=2&step=2", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	errorHandlerNone := func(resp http.ResponseWriter, req *http.Request) {
		assert.Nil(t, recover())
	}
	handlers := StockHandlers{
		Context: &Context{
			sh:         schema.NewDecoder(),
			historyAPI: &mockedHistoryAPI,
			esStock:    &mockedEsStock,
		},
		getDate:      getTestDate,
		errorHandler: errorHandlerNone,
	}
	params := httprouter.Params{httprouter.Param{Key: "symbol", Value: symbolTest}}
	handlers.History(resp, req, params)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Equal(t, string(stcoskAggJson), resp.Body.String())
}
