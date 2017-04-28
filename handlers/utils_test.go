package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func getTestDate() time.Time {
	time, _ := time.Parse(time.RFC3339, "2016-12-13T01:24:23Z")
	return time
}

func createEcho(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	return c, resp
}

func createErrorHandler(t *testing.T, expectedStatus int, expectedErrorMsg string) errorHandlerFunc {
	return func(c echo.Context, status int, err error) error {
		assert.NotNil(t, c)
		assert.Equal(t, expectedStatus, status)
		assert.Equal(t, expectedErrorMsg, err.Error())
		return err
	}
}
