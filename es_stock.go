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
	"fmt"
	"math"
	"time"

	finance "github.com/clebi/yfinance"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	indexName               = "stocks-hist"
	indexType               = "stock_day"
	indexTimeout            = 3 * time.Second
	timeAggregationName     = "time_agg"
	avgCloseAggregationName = "avg_close"
	movCloseAggregationName = "mov_close"
)

// EsStocksAgg is the a stock aggregation
type EsStocksAgg struct {
	Symbol   string  `json:"symbol"`
	MsTime   int64   `json:"mstime"`
	AvgClose float64 `json:"close"`
	MovClose float64 `json:"mv_close"`
}

// IEsStock contains elasticsearch manager actions
type IEsStock interface {
	Index(stock finance.Stock) error
	GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]EsStocksAgg, error)
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

// GetStocksAgg retrieves aggregations of stock values by dates
//
//  GetStocksAgg("TEST", startDate, endDate)
//
// returns an array ofg stocks aggregations
func (esStock *EsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]EsStocksAgg, error) {
	movStartDate := startDate.AddDate(0, 0, movAvgWindow*-1)
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	query := elastic.NewQueryStringQuery(fmt.Sprintf("symbol = %s AND date: [%s TO %s]",
		symbol, movStartDate.Format(finance.DateFormat), endDate.Format(finance.DateFormat)))
	avgCloseAgg := elastic.NewAvgAggregation().Field("close")
	movCloseAgg := elastic.NewMovAvgAggregation().BucketsPath(avgCloseAggregationName).
		Window(int(math.Ceil(float64(movAvgWindow) / float64(step))))
	minDateAgg := elastic.NewMinAggregation().Field("date")
	selectAgg := elastic.NewBucketSelectorAggregation().
		AddBucketsPath("avg_close", avgCloseAggregationName).
		AddBucketsPath("date", "min_date").
		Script(elastic.NewScript(fmt.Sprintf("params.avg_close > 0 && params.date >= %dL", startDate.Unix()*1000)))
	timeAgg := elastic.NewDateHistogramAggregation().Field("date").Interval(fmt.Sprintf("%dd", step)).
		SubAggregation(avgCloseAggregationName, avgCloseAgg).
		SubAggregation(movCloseAggregationName, movCloseAgg).
		SubAggregation("min_date", minDateAgg).
		SubAggregation("selector", selectAgg)
	results, err := esStock.es.Search(indexName).
		Type(indexType).
		Query(query).
		Aggregation(timeAggregationName, timeAgg).
		Size(0).
		Do(esContext)
	if err != nil {
		return nil, err
	}
	resAgg, _ := results.Aggregations.DateHistogram(timeAggregationName)
	stocks := make([]EsStocksAgg, len(resAgg.Buckets))
	for i, bucket := range resAgg.Buckets {
		avg, _ := bucket.Avg(avgCloseAggregationName)
		mov, movOk := bucket.MovAvg(movCloseAggregationName)
		stocks[i] = EsStocksAgg{
			Symbol:   symbol,
			MsTime:   int64(bucket.Key),
			AvgClose: *avg.Value,
		}
		if movOk {
			stocks[i].MovClose = *mov.Value
		} else {
			stocks[i].MovClose = *avg.Value
		}
	}
	return stocks, nil
}
