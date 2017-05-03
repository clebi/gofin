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

	"github.com/clebi/gofin/es"
	"github.com/labstack/echo"
)

// PositionDisplay contains all fields to display to the client
type PositionDisplay struct {
	es.PositionAgg
	Name  string
	Value float32
}

// PositionHandlers handles all request to position management
type PositionHandlers struct {
	*Context
	errorHandler errorHandlerFunc
}

// NewPositionHandlers creates a new position handlers object
func NewPositionHandlers(context *Context) *PositionHandlers {
	return &PositionHandlers{
		Context:      context,
		errorHandler: handleError,
	}
}

// AddPosition handles http request to save a position
//
// This function is a handler for http server, it should not be called directly
func (handlers *PositionHandlers) AddPosition(c echo.Context) error {
	position := new(es.Position)
	if err := c.Bind(position); err != nil {
		return handlers.errorHandler(c, http.StatusBadRequest, err)
	}
	if err := handlers.validator.Struct(position); err != nil {
		return handlers.errorHandler(c, http.StatusBadRequest, err)
	}
	if err := handlers.esPosition.AddPosition(position); err != nil {
		return handlers.errorHandler(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, position)
}

// GetPositions handles http request to retrieve the user's positions
//
// This function is a handler for http server, it should not be called directly
func (handlers *PositionHandlers) GetPositions(c echo.Context) error {
	positions, err := handlers.esPosition.GetPositions("tester")
	if err != nil {
		return handlers.errorHandler(c, http.StatusInternalServerError, err)
	}
	displayPosition := make([]PositionDisplay, len(positions))
	for i, position := range positions {
		quote, err := handlers.quotesAPI.GetQuote(position.Symbol)
		if err != nil {
			return handlers.errorHandler(c, http.StatusInternalServerError, err)
		}
		displayPosition[i] = PositionDisplay{PositionAgg: position, Name: quote.Name, Value: quote.LastTradePriceOnly}
	}
	return c.JSON(http.StatusOK, displayPosition)
}
