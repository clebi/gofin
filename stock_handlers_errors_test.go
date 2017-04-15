package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	unknownErrorMsg  = "unknown_error"
	unknownErrorResp = "{\"status\":\"error\",\"description\":\"unknown_error\"}"
	badRequestMsg    = "bad_request"
	badRequestResp   = "{\"status\":\"error\",\"description\":\"bad_request\"}"
	decoderErrorMsg  = "decoder_error"
	genericErrorMsg  = "generic_error"
)

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

func TestHistoryErrors(t *testing.T) {
	for _, tt := range errorTests {
		handlers := StockHandlers{
			Context:      tt.context,
			getDate:      tt.getDate,
			validator:    tt.validator,
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
			validator:    tt.validator,
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
