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
	"github.com/clebi/yfinance"
	"github.com/labstack/echo"
	"github.com/rs/cors"
	elastic "gopkg.in/olivere/elastic.v5"

	schema "github.com/gorilla/Schema"
)

const (
	defaultServerURL = ":9000"
)

// SchemaDecoder decodes URL query to struct
type SchemaDecoder interface {
	Decode(dst interface{}, src map[string][]string) error
}

// Context is the context of the application
// It contains resources that needs to be access in http handlers
type Context struct {
	es         *elastic.Client
	sh         SchemaDecoder
	historyAPI finance.HistoryAPI
	esStock    IEsStock
}

func main() {
	// Initialize logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Initialize elasticsearch client
	es, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	sh := schema.NewDecoder()
	sh.IgnoreUnknownKeys(true)
	// Initialize app context
	context := Context{
		es:         es,
		sh:         sh,
		historyAPI: finance.NewHistory(),
		esStock:    NewEsStock(es),
	}

	stockHandlers := NewStockHandlers(&context)
	router := echo.New()
	router.GET("/graph/:symbol", stockHandlers.History)
	handler := cors.Default().Handler(router)
	log.WithFields(log.Fields{"url": defaultServerURL}).Info("Start server")
	log.Fatal(http.ListenAndServe(defaultServerURL, handler))
}
