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
	"testing"

	"github.com/clebi/gofin/es"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	positionErrorMsg         = "position_error_msg"
	positionQuoteAPIErrorMsg = "position_quote_api_error_msg"
)

var addPositionErrorTests = []struct {
	echo            echo.Context
	context         *Context
	expectedStatus  int
	expectedMessage string
}{
	{
		&ErrorEchoBind{Msg: positionErrorMsg},
		nil,
		http.StatusBadRequest,
		positionErrorMsg,
	},
	{
		&DummyEchoBind{},
		&Context{validator: &ErrorStructValidator{Msg: positionErrorMsg}},
		http.StatusBadRequest,
		positionErrorMsg,
	},
	{
		&DummyEchoBind{},
		&Context{
			esPosition: &ErrorEsPosition{Msg: positionErrorMsg},
			validator:  &ErrorStructValidator{Msg: positionErrorMsg},
		},
		http.StatusBadRequest,
		positionErrorMsg,
	},
}

func TestAddPositionErrors(t *testing.T) {
	for _, tt := range addPositionErrorTests {
		handlers := PositionHandlers{
			Context:      tt.context,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
		}
		_, err := http.NewRequest(testHistoryListMethod, testHistoryListRequest, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		res := handlers.AddPosition(tt.echo)
		assert.NotNil(t, res)
	}
}

var getPositionErrorTests = []struct {
	echo            echo.Context
	context         *Context
	expectedStatus  int
	expectedMessage string
}{
	{
		&DummyEchoBind{},
		&Context{esPosition: &ErrorEsPosition{Msg: positionErrorMsg}},
		http.StatusInternalServerError,
		positionErrorMsg,
	},
	{
		&DummyEchoBind{},
		&Context{
			esPosition: &DummyEsPosition{PositionAgg: []es.PositionAgg{{Symbol: "test_agg", Number: 5, Cost: 14}}},
			quotesAPI:  &ErrorQuotesAPI{Msg: positionQuoteAPIErrorMsg},
		},
		http.StatusInternalServerError,
		positionQuoteAPIErrorMsg,
	},
}

func TestGetPositionsErrors(t *testing.T) {
	for _, tt := range getPositionErrorTests {
		handlers := PositionHandlers{
			Context:      tt.context,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
		}
		_, err := http.NewRequest(testHistoryListMethod, testHistoryListRequest, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		res := handlers.GetPositions(tt.echo)
		assert.NotNil(t, res)
	}
}
