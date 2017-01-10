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

package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/clebi/yfinance"
	"github.com/gorilla/mux"
)

type errorDesc struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// StockHandlers is an object containing all the handlers concerning stocks
type StockHandlers struct {
	*Context
}

func (handlers *StockHandlers) handleErrors(resp http.ResponseWriter, req *http.Request) {
	if err := recover(); err != nil {
		var errorMsg string
		switch err := err.(type) {
		case finance.YApiError:
			errorMsg = err.Error()
			resp.WriteHeader(http.StatusBadRequest)
		default:
			log.Error(err)
			errorMsg = "unknown_error"
			resp.WriteHeader(http.StatusInternalServerError)
		}
		bresp, err := json.Marshal(errorDesc{
			Status:      "error",
			Description: errorMsg,
		})
		if err != nil {
			return
		}
		resp.Write(bresp)
	}
}

// History retrieve stocks history
//
//  handlers.History(w, r)
//
//  This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) History(resp http.ResponseWriter, req *http.Request) {
	defer handlers.handleErrors(resp, req)
	vars := mux.Vars(req)
	end := time.Now().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -30)
	history := finance.NewHistory()
	stocks, err := history.GetHistory(vars["symbol"], start, end)
	if err != nil {
		panic(err)
	}
	stockEs := NewEsStock(handlers.es)
	for _, stock := range stocks {
		err = stockEs.Index(stock)
		if err != nil {
			panic(err)
		}
	}
	bresp, err := json.Marshal(stocks)
	if err != nil {
		panic(err)
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(bresp)
}
