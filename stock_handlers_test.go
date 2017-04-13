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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	finance "github.com/clebi/yfinance"
	schema "github.com/gorilla/Schema"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	testHistoryMethod  = "GET"
	testHistoryRequest = "http://test.test/graph/TEST?days=3&window=2&step=2"
	symbolTest         = "TEST"
	unknownErrorMsg    = "unknown_error"
	unknownErrorResp   = "{\"status\":\"error\",\"description\":\"unknown_error\"}"
	badRequestMsg      = "bad_request"
	badRequestResp     = "{\"status\":\"error\",\"description\":\"bad_request\"}"
	decoderErrorMsg    = "decoder_error"
	genericErrorMsg    = "generic_error"
)

var (
	testEndDate      = getTestDate().Truncate(24 * time.Hour)
	testStartDate    = testEndDate.AddDate(0, 0, -3)
	testStartMovDate = testStartDate.AddDate(0, 0, -2)
)

func TestHistory(t *testing.T) {
	stock := finance.Stock{Open: 1.1, High: 2.2, Low: 3.3, Close: 4.4, Volume: 999, Symbol: symbolTest, Date: finance.YTime{Time: testStartDate}}
	stocks := []finance.Stock{stock}
	mockedHistoryAPI := mockHistoryAPI{}
	mockedHistoryAPI.On("GetHistory", symbolTest, testStartMovDate, testEndDate).Return(stocks, nil)
	stocksAgg := []EsStocksAgg{{Symbol: symbolTest, MsTime: testStartMovDate.Unix() * 1000, AvgClose: 4.4, MovClose: 4.1}}
	stockAggJSON, err := json.Marshal(stocksAgg)
	if err != nil {
		t.Fatal(err)
	}
	mockedEsStock := mockEsStock{}
	mockedEsStock.On("Index", stock).Return(nil)
	mockedEsStock.On("GetStocksAgg", symbolTest, 2, 2, testStartDate, testEndDate).Return(stocksAgg, nil)
	resp := httptest.NewRecorder()
	req, err := http.NewRequest(testHistoryMethod, testHistoryRequest, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	handlers := StockHandlers{
		Context: &Context{
			sh:         schema.NewDecoder(),
			historyAPI: &mockedHistoryAPI,
			esStock:    &mockedEsStock,
		},
		validator: &DummyStructValidator{},
		getDate:   getTestDate,
	}
	e := echo.New()
	c := e.NewContext(req, resp)
	c.SetParamNames("symbol")
	c.SetParamValues(symbolTest)
	handlers.History(c)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Equal(t, string(stockAggJSON), resp.Body.String())
}

var errorTests = []struct {
	context         *Context
	getDate         GetDateFunc
	validator       StructValidator
	expectedStatus  int
	expectedMessage string
}{
	{
		&Context{sh: &ErrorSchemaDecoder{Msg: genericErrorMsg}},
		getTestDate,
		nil,
		http.StatusInternalServerError,
		genericErrorMsg,
	},
	{
		&Context{sh: &DummySchemaDecoder{}},
		getTestDate,
		&ErrorStructValidator{Msg: genericErrorMsg},
		http.StatusBadRequest,
		genericErrorMsg,
	},
	{
		&Context{sh: &DummySchemaDecoder{}, historyAPI: &ErrorFinanceAPI{Msg: genericErrorMsg}},
		getTestDate,
		&DummyStructValidator{},
		http.StatusBadRequest,
		genericErrorMsg,
	},
	{
		&Context{sh: &DummySchemaDecoder{}, historyAPI: &OneItemFinanceAPI{}, esStock: &ErrorEsStock{Msg: genericErrorMsg}},
		getTestDate,
		&DummyStructValidator{},
		http.StatusInternalServerError,
		genericErrorMsg,
	},
	{
		&Context{sh: &DummySchemaDecoder{}, historyAPI: &DummyFinanceAPI{}, esStock: &ErrorEsStock{Msg: genericErrorMsg}},
		getTestDate,
		&DummyStructValidator{},
		http.StatusInternalServerError,
		genericErrorMsg,
	},
}

func TestSchemaDecodeError(t *testing.T) {
	for _, tt := range errorTests {
		handlers := StockHandlers{
			Context:      tt.context,
			getDate:      tt.getDate,
			validator:    tt.validator,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
		}
		req, err := http.NewRequest(testHistoryMethod, testHistoryRequest, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		c, _ := createEcho(req)
		res := handlers.History(c)
		assert.NotNil(t, res)
	}
}

func TestHandleErrorInternal(t *testing.T) {
	c, resp := createEcho(nil)
	handleError(c, http.StatusInternalServerError, errors.New(unknownErrorMsg))
	assert.Equal(t, unknownErrorResp, resp.Body.String())
}

func TestHandleError(t *testing.T) {
	c, resp := createEcho(nil)
	handleError(c, http.StatusBadRequest, errors.New(badRequestMsg))
	assert.Equal(t, badRequestResp, resp.Body.String())
}
