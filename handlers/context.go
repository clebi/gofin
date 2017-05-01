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

// Context is the context of the application
import (
	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
	elastic "gopkg.in/olivere/elastic.v5"
)

// SchemaDecoder decodes URL query to struct
type SchemaDecoder interface {
	Decode(dst interface{}, src map[string][]string) error
}

//Context contains resources that needs to be access in http handlers
type Context struct {
	es         *elastic.Client
	sh         SchemaDecoder
	historyAPI finance.HistoryAPI
	quotesAPI  finance.QuotesAPI
	esStock    es.IStock
	esPosition es.IPositionStock
}

//NewContext creates a new context for handlers
func NewContext(
	es *elastic.Client,
	sh SchemaDecoder,
	historyAPI finance.HistoryAPI,
	quotesAPI finance.QuotesAPI,
	esStock es.IStock,
	esPosition es.IPositionStock) *Context {
	return &Context{
		es:         es,
		sh:         sh,
		historyAPI: historyAPI,
		quotesAPI:  quotesAPI,
		esStock:    esStock,
		esPosition: esPosition,
	}
}
