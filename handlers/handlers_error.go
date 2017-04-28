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
