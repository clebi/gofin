package handlers

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	addPositionData = "{\"username\":\"test_username\",\"broker\":\"test\",\"symbol\":\"test\"," +
		"\"date\":\"2017-04-20T13:00:45Z\",\"number\":1,\"value\":22,\"cost\":24}"
)

func TestAddPosition(t *testing.T) {
	handlers := &PositionHandlers{
		Context: &Context{
			esPosition: &DummyEsPosition{},
		},
		validator: &DummyStructValidator{},
	}
	req, err := http.NewRequest("POST", "http://test.test/position", bytes.NewBufferString(addPositionData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	c, resp := createEcho(req)
	handlers.AddPosition(c)
	assert.Equal(t, addPositionData, resp.Body.String())
}
