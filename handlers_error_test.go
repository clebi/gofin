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
)

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
