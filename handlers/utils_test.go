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

func testIndexStockNoError(context *Context, symbol string, start time.Time, end time.Time) *HandlerERROR {
	return nil
}

func createTestIndexStockError(status int, msg string) indexStockFunc {
	return func(context *Context, symbol string, start time.Time, end time.Time) *HandlerERROR {
		return &HandlerERROR{
			Status: status,
			error:  errors.New(msg),
		}
	}
}
