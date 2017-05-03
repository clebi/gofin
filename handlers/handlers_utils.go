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

	"github.com/labstack/echo"
)

func getQuery(c echo.Context, context *Context, params interface{}) *HandlerERROR {
	if err := context.sh.Decode(params, c.Request().URL.Query()); err != nil {
		return &HandlerERROR{error: err, Status: http.StatusInternalServerError}
	}
	if err := context.validator.Struct(params); err != nil {
		return &HandlerERROR{error: err, Status: http.StatusBadRequest}
	}
	return nil
}
