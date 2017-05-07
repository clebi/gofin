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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
	schema "github.com/gorilla/Schema"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	testHistoryMethod      = "GET"
	testHistoryRequest     = "http://test.test/graph/TEST?days=3&window=2&step=2"
	testHistoryListMethod  = "GET"
	testHistoryListSymbol1 = "TEST.1"
	testHistoryListSymbol2 = "TEST.2"
	testHistoryListRequest = "http://test.test/history/list/?days=3&window=2&step=2&symbols=" + testHistoryListSymbol1 + "&symbols=" + testHistoryListSymbol2
	symbolTest             = "TEST"
)

var (
	testEndDate      = getTestDate().Truncate(24 * time.Hour)
	testStartDate    = testEndDate.AddDate(0, 0, -3)
	testStartMovDate = testStartDate.AddDate(0, 0, -2)
)

func prepareHisotryCall(
	method string,
	request string,
	respBody interface{},
	mockedHistoryAPI finance.HistoryAPI,
	mockedStock es.IStock) (*http.Request, []byte, *StockHandlers, error) {
	req, err := http.NewRequest(method, request, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	stocksAggsJSON, err := json.Marshal(respBody)
	if err != nil {
		return nil, nil, nil, err
	}

	handlers := StockHandlers{
		Context: &Context{
			sh:         schema.NewDecoder(),
			historyAPI: mockedHistoryAPI,
			esStock:    mockedStock,
			validator:  &DummyStructValidator{},
		},
		getDate: getTestDate,
	}
	return req, stocksAggsJSON, &handlers, nil
}

func TestHistory(t *testing.T) {
	stock := finance.Stock{Open: 1.1, High: 2.2, Low: 3.3, Close: 4.4, Volume: 999, Symbol: symbolTest, Date: finance.YTime{Time: testStartDate}}
	stocks := []finance.Stock{stock}
	mockedHistoryAPI := mockHistoryAPI{}
	mockedHistoryAPI.On("GetHistory", symbolTest, testStartMovDate, testEndDate).Return(stocks, nil)
	mockedEsStock := mockEsStock{
		stockAggs: map[string][]es.StocksAgg{
			symbolTest: {{Symbol: symbolTest, MsTime: testStartMovDate.Unix() * 1000, AvgClose: 4.4, MovClose: 4.1}},
		},
	}
	resp := httptest.NewRecorder()
	req, stockAggJSON, handlers, err := prepareHisotryCall(
		testHistoryMethod,
		testHistoryRequest,
		mockedEsStock.stockAggs[symbolTest],
		&mockedHistoryAPI,
		&mockedEsStock)
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	c := e.NewContext(req, resp)
	c.SetParamNames("symbol")
	c.SetParamValues(symbolTest)
	handlers.History(c)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Equal(t, string(stockAggJSON), resp.Body.String())
}

func TestHistoryList(t *testing.T) {
	stocks := [][]finance.Stock{
		{{Open: 1.1, High: 1.2, Low: 2.3, Close: 2.4, Volume: 111, Symbol: testHistoryListSymbol1, Date: finance.YTime{Time: testStartDate}}},
		{{Open: 2.1, High: 2.2, Low: 2.3, Close: 2.4, Volume: 222, Symbol: testHistoryListSymbol2, Date: finance.YTime{Time: testStartDate}}},
	}
	var stocksAggs [][]es.StocksAgg
	mapStocksAggs := map[string][]es.StocksAgg{}
	mockedHistoryAPI := mockHistoryAPI{}
	for _, stockList := range stocks {
		symbol := stockList[0].Symbol
		mockedHistoryAPI.On("GetHistory", symbol, testStartMovDate, testEndDate).Return(stockList, nil)
		stocksAgg := []es.StocksAgg{
			{
				Symbol:   symbol,
				MsTime:   testStartMovDate.Unix() * 1000,
				AvgClose: float64(stockList[0].Close) + float64(0.3),
				MovClose: float64(stockList[0].Close) + float64(0.1),
			},
		}
		mapStocksAggs[symbol] = stocksAgg
		stocksAggs = append(stocksAggs, stocksAgg)
	}

	mockedEsStock := mockEsStock{stockAggs: mapStocksAggs}
	req, stocksAggsJSON, handlers, err := prepareHisotryCall(
		testHistoryListMethod,
		testHistoryListRequest,
		stocksAggs,
		&mockedHistoryAPI,
		&mockedEsStock)
	if err != nil {
		t.Fatal(err)
	}

	c, resp := createEcho(req)
	handlers.HistoryList(c)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Equal(t, string(stocksAggsJSON), resp.Body.String())
}
