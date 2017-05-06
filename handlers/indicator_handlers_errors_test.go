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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	indicatorGetStocksErrorMsg = "indicator_get_stocks_error"
)

var indicatorGetStocksErrorTests = []struct {
	context         *Context
	expectedStatus  int
	expectedMessage string
	indexStockFunc  indexStockFunc
}{
	{
		&Context{sh: &ErrorSchemaDecoder{Msg: indicatorGetStocksErrorMsg}},
		http.StatusInternalServerError,
		indicatorGetStocksErrorMsg,
		nil,
	},
	{
		&Context{
			sh:        &IndicatorSchemaDecoder{Symbols: []string{"ERROR"}},
			validator: &DummyStructValidator{},
		},
		http.StatusInternalServerError,
		indicatorGetStocksErrorMsg,
		createTestIndexStockError(http.StatusInternalServerError, indicatorGetStocksErrorMsg),
	},
	{
		&Context{
			sh:         &IndicatorSchemaDecoder{Symbols: []string{"ERROR"}},
			validator:  &DummyStructValidator{},
			quotesAPI:  &ErrorQuotesAPI{Msg: indicatorGetStocksErrorMsg},
			historyAPI: &ErrorFinanceAPI{Msg: genericErrorMsg},
		},
		http.StatusInternalServerError,
		indicatorGetStocksErrorMsg,
		testIndexStockNoError,
	},
	{
		&Context{
			sh:         &IndicatorSchemaDecoder{Symbols: []string{"ERROR"}},
			validator:  &DummyStructValidator{},
			quotesAPI:  &DummyQuotesAPI{},
			historyAPI: &ErrorFinanceAPI{Msg: genericErrorMsg},
			esStock:    &IndicatorErrorEsStock{index: 0, errs: []error{errors.New(indicatorGetStocksErrorMsg)}},
		},
		http.StatusInternalServerError,
		indicatorGetStocksErrorMsg,
		testIndexStockNoError,
	},
	{
		&Context{
			sh:         &IndicatorSchemaDecoder{Symbols: []string{"ERROR"}},
			validator:  &DummyStructValidator{},
			quotesAPI:  &DummyQuotesAPI{},
			historyAPI: &ErrorFinanceAPI{Msg: genericErrorMsg},
			esStock:    &IndicatorErrorEsStock{index: 0, errs: []error{nil, errors.New(indicatorGetStocksErrorMsg)}},
		},
		http.StatusInternalServerError,
		indicatorGetStocksErrorMsg,
		testIndexStockNoError,
	},
}

func TestGetStocksErrors(t *testing.T) {
	for _, tt := range indicatorGetStocksErrorTests {
		handlers := IndicatorHandlers{
			Context:      tt.context,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
			getDate:      getTestDate,
			indexStock:   tt.indexStockFunc,
		}
		req, err := http.NewRequest("GET", testGetStocksURL, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		c, _ := createEcho(req)
		res := handlers.GetStocks(c)
		assert.NotNil(t, res)
	}
}
