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

// StockHandlers is an object containing all the handlers concerning stocks
type StockHandlers struct {
	*Context
}

// History retrieve stocks history
//
//  handlers.History(w, r)
//
//  This function is a handler for http server, it should not be called directly
func (handlers *StockHandlers) History(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	end := time.Now().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -30)
	history := finance.NewHistory()
	stocks, err := history.GetHistory(vars["symbol"], start, end)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	stockEs := NewEsStock(handlers.es)
	for _, stock := range stocks {
		err = stockEs.Index(stock)
		if err != nil {
			log.Error(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	bresp, err := json.Marshal(stocks)
	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(bresp)
}
