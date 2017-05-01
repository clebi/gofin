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
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/clebi/gofin/es"
	"github.com/clebi/gofin/handlers"
	"github.com/clebi/yfinance"
	"github.com/labstack/echo"
	"github.com/rs/cors"
	elastic "gopkg.in/olivere/elastic.v5"

	schema "github.com/gorilla/Schema"
)

const (
	defaultServerURL = ":9000"
)

func main() {
	// Initialize logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Initialize elasticsearch client
	esClient, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	sh := schema.NewDecoder()
	sh.IgnoreUnknownKeys(true)
	// Initialize app context
	context := handlers.NewContext(
		esClient,
		sh,
		finance.NewHistory(),
		finance.NewQuotes(),
		es.NewStock(esClient),
		es.NewPosition(esClient),
	)

	stockHandlers := handlers.NewStockHandlers(context)
	positionHandlers := handlers.NewPositionHandlers(context)
	router := echo.New()
	router.GET("/history/:symbol", stockHandlers.History)
	router.GET("/history/list", stockHandlers.HistoryList)
	router.POST("/position", positionHandlers.AddPosition)
	router.GET("/position", positionHandlers.GetPositions)
	handler := cors.Default().Handler(router)
	log.WithFields(log.Fields{"url": defaultServerURL}).Info("Start server")
	log.Fatal(http.ListenAndServe(defaultServerURL, handler))
}
