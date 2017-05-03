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
	"bytes"
	"net/http"
	"testing"

	finance "github.com/clebi/yfinance"
	"github.com/stretchr/testify/assert"
)

const (
	testGetStocksURL       = "http://test.test/indicator"
	testGetStocksResultStr = "[{\"Symbol\":\"TEST1\",\"Name\":\"TEST_NAME_1\",\"Value\":1.1,\"MM200\":1.3,\"MM50\":1.2," +
		"\"MM200D50\":0.923077},{\"Symbol\":\"TEST2\",\"Name\":\"TEST_NAME_2\",\"Value\":2.1,\"MM200\":2.3,\"MM50\":2.2," +
		"\"MM200D50\":0.9565218}]"
)

func TestGetStocks(t *testing.T) {
	handlers := &IndicatorHandlers{
		Context: &Context{
			sh:        &IndicatorSchemaDecoder{Symbols: []string{"TEST1", "TEST2"}},
			validator: &DummyStructValidator{},
			quotesAPI: &IndicatorQuotesAPI{quotes: map[string]*finance.Quote{
				"TEST1": &finance.Quote{
					Symbol:                     "TEST1",
					Name:                       "TEST_NAME_1",
					LastTradePriceOnly:         1.1,
					FiftydayMovingAverage:      1.2,
					TwoHundreddayMovingAverage: 1.3,
					Volume: 14,
				},
				"TEST2": &finance.Quote{
					Symbol:                     "TEST2",
					Name:                       "TEST_NAME_2",
					LastTradePriceOnly:         2.1,
					FiftydayMovingAverage:      2.2,
					TwoHundreddayMovingAverage: 2.3,
					Volume: 24,
				},
			}},
		},
	}
	req, err := http.NewRequest("POST", testGetStocksURL, bytes.NewBufferString(addPositionData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	c, resp := createEcho(req)
	handlers.GetStocks(c)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
	assert.Equal(t, testGetStocksResultStr, resp.Body.String())
}
