package main

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	positionErrorMsg = "position_error_msg"
)

var positionErrorTests = []struct {
	echo            echo.Context
	context         *Context
	validator       StructValidator
	expectedStatus  int
	expectedMessage string
}{
	{
		&ErrorEchoBind{Msg: positionErrorMsg},
		nil,
		nil,
		http.StatusBadRequest,
		positionErrorMsg,
	},
	{
		&DummyEchoBind{},
		nil,
		&ErrorStructValidator{Msg: positionErrorMsg},
		http.StatusBadRequest,
		positionErrorMsg,
	},
	{
		&DummyEchoBind{},
		&Context{esPosition: &ErrorEsPosition{Msg: positionErrorMsg}},
		&DummyStructValidator{},
		http.StatusBadRequest,
		positionErrorMsg,
	},
}

func TestAddPositionErrors(t *testing.T) {
	for _, tt := range positionErrorTests {
		handlers := PositionHandlers{
			Context:      tt.context,
			validator:    tt.validator,
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
