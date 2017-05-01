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

	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
	"github.com/stretchr/testify/assert"
)

const (
	addPositionData = "{\"username\":\"test_username\",\"broker\":\"test\",\"symbol\":\"test\"," +
		"\"date\":\"2017-04-20T13:00:45Z\",\"number\":1,\"value\":22,\"cost\":24}"
	getPositionsData = "[{\"Symbol\":\"test_agg\",\"Number\":5,\"Cost\":14,\"Name\":\"TEST NAME\",\"Value\":15}]"
)

func TestAddPosition(t *testing.T) {
	handlers := &PositionHandlers{
		Context: &Context{
			esPosition: &DummyEsPosition{},
		},
		validator: &DummyStructValidator{},
	}
	req, err := http.NewRequest("POST", "http://test.test/position", bytes.NewBufferString(addPositionData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	c, resp := createEcho(req)
	handlers.AddPosition(c)
	assert.Equal(t, addPositionData, resp.Body.String())
}

func TestGetPositions(t *testing.T) {
	handlers := &PositionHandlers{
		Context: &Context{
			quotesAPI: &DummyQuotesAPI{quote: finance.Quote{Name: "TEST NAME", LastTradePriceOnly: 15}},
			esPosition: &DummyEsPosition{
				PositionAgg: []es.PositionAgg{{Symbol: "test_agg", Number: 5, Cost: 14}},
			},
		},
	}
	req, err := http.NewRequest("POST", "http://test.test/position", bytes.NewBufferString(addPositionData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	c, resp := createEcho(req)
	handlers.GetPositions(c)
	assert.Equal(t, getPositionsData, resp.Body.String())
}
