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

// handleError writes error to the http channel and logs internal errors
import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func handleError(c echo.Context, status int, err error) error {
	var msg string
	if status == http.StatusInternalServerError {
		log.Error(err)
		msg = "unknown_error"
	} else {
		msg = err.Error()
	}
	c.JSON(status, errorDesc{Status: "error", Description: msg})
	return nil
}
