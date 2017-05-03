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

	"github.com/stretchr/testify/assert"
)

const (
	genericErrorMsg = "generic_error"
)

var errorTests = []struct {
	context         *Context
	getDate         GetDateFunc
	expectedStatus  int
	expectedMessage string
}{
	{
		&Context{sh: &ErrorSchemaDecoder{Msg: genericErrorMsg}},
		getTestDate,
		http.StatusInternalServerError,
		genericErrorMsg,
	},
	{
		&Context{sh: &DummySchemaDecoder{}, validator: &ErrorStructValidator{Msg: genericErrorMsg}},
		getTestDate,
		http.StatusBadRequest,
		genericErrorMsg,
	},
	{
		&Context{
			sh:         &DummySchemaDecoder{},
			historyAPI: &ErrorFinanceAPI{Msg: genericErrorMsg},
			validator:  &DummyStructValidator{},
		},
		getTestDate,
		http.StatusBadRequest,
		genericErrorMsg,
	},
	{
		&Context{
			sh:         &DummySchemaDecoder{},
			historyAPI: &OneItemFinanceAPI{},
			esStock:    &ErrorEsStock{Msg: genericErrorMsg},
			validator:  &DummyStructValidator{},
		},
		getTestDate,
		http.StatusInternalServerError,
		genericErrorMsg,
	},
	{
		&Context{
			sh:         &DummySchemaDecoder{},
			historyAPI: &DummyFinanceAPI{},
			esStock:    &ErrorEsStock{Msg: genericErrorMsg},
			validator:  &DummyStructValidator{},
		},
		getTestDate,
		http.StatusInternalServerError,
		genericErrorMsg,
	},
}

func TestHistoryErrors(t *testing.T) {
	for _, tt := range errorTests {
		handlers := StockHandlers{
			Context:      tt.context,
			getDate:      tt.getDate,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
		}
		req, err := http.NewRequest(testHistoryListMethod, testHistoryListRequest, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		c, _ := createEcho(req)
		res := handlers.History(c)
		assert.NotNil(t, res)
	}
}

func TestHistoryListErrors(t *testing.T) {
	for _, tt := range errorTests {
		handlers := StockHandlers{
			Context:      tt.context,
			getDate:      tt.getDate,
			errorHandler: createErrorHandler(t, tt.expectedStatus, tt.expectedMessage),
		}
		req, err := http.NewRequest(testHistoryMethod, testHistoryRequest, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		c, _ := createEcho(req)
		res := handlers.HistoryList(c)
		assert.NotNil(t, res)
	}
}
