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
	"context"
	"time"

	finance "github.com/clebi/yfinance"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	indexName    = "stocks-hist"
	indexType    = "stock_day"
	indexTimeout = 3 * time.Second
)

// IEsStock contains elasticsearch manager actions
type IEsStock interface {
	Index(stock finance.Stock) error
}

// EsStock manage stocks in elasticsearch
type EsStock struct {
	es *elastic.Client
}

// NewEsStock create a new elasticsearch manager object
func NewEsStock(es *elastic.Client) IEsStock {
	return &EsStock{
		es: es,
	}
}

// Index is used to index a stock into elasticsearch
func (esStock *EsStock) Index(stock finance.Stock) error {
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	stockMap := map[string]interface{}{
		"date":   stock.Date.Format(time.RFC3339),
		"open":   stock.Open,
		"high":   stock.High,
		"low":    stock.Low,
		"close":  stock.Close,
		"volume": stock.Volume,
		"symbol": stock.Symbol,
	}
	_, err := esStock.es.Index().
		Index(indexName).
		Type(indexType).
		Id(stock.Symbol + "_" + stock.Date.Format(finance.DateFormat)).
		BodyJson(stockMap).
		Do(esContext)
	if err != nil {
		return err
	}
	return nil
}
