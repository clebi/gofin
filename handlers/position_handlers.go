package handlers

import (
	"net/http"

	"github.com/clebi/gofin/es"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

// PositionHandlers handles all request to position management
type PositionHandlers struct {
	*Context
	validator    StructValidator
	errorHandler errorHandlerFunc
}

// NewPositionHandlers creates a new position handlers object
func NewPositionHandlers(context *Context) *PositionHandlers {
	return &PositionHandlers{
		Context:      context,
		validator:    validator.New(),
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
