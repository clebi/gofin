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
	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
)

const (
	defaultServerURL = ":9000"
)

// Context is the context of the application
// It contains resources that needs to be access in http handlers
type Context struct {
	es *elastic.Client
}

func main() {
	// Initialize logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Initialize elasticsearch client
	es, err := elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize app context
	context := Context{
		es: es,
	}

	stockHandlers := StockHandlers{
		Context: &context,
	}
	router := mux.NewRouter()
	router.HandleFunc("/graph/{symbol}", stockHandlers.History)
	log.WithFields(log.Fields{"url": defaultServerURL}).Info("Start server")
	log.Fatal(http.ListenAndServe(defaultServerURL, router))
}
