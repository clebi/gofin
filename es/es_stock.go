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

package es

import (
	"context"
	"encoding/json"
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
	statsAggregationName    = "stats"
)

type stockValue struct {
	Date string `json:"date"`
}

// StocksStats contains all stats concerning a stock
type StocksStats struct {
	Symbol            string
	StandardDeviation float64
	Avg               float64
}

// StocksAgg is the a stock aggregation
type StocksAgg struct {
	Symbol   string  `json:"symbol"`
	MsTime   int64   `json:"mstime"`
	AvgClose float64 `json:"close"`
	MovClose float64 `json:"mv_close"`
}

// IStock contains elasticsearch manager actions
type IStock interface {
	Index(stock finance.Stock) error
	GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]StocksAgg, error)
	GetStockStats(symbol string, startDate time.Time, endDate time.Time) (*StocksStats, error)
	GetDateForNumPoint(symbol string, numPoints int, endDate time.Time) (*time.Time, error)
}

// Stock manage stocks in elasticsearch
type Stock struct {
	es *elastic.Client
}

// NewStock create a new elasticsearch manager object
func NewStock(es *elastic.Client) IStock {
	return &Stock{
		es: es,
	}
}

// Index is used to index a stock into elasticsearch
func (esStock *Stock) Index(stock finance.Stock) error {
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
func (esStock *Stock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]StocksAgg, error) {
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
	stocks := make([]StocksAgg, len(resAgg.Buckets))
	for i, bucket := range resAgg.Buckets {
		avg, _ := bucket.Avg(avgCloseAggregationName)
		mov, movOk := bucket.MovAvg(movCloseAggregationName)
		stocks[i] = StocksAgg{
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

// GetStockStats retrives the stats about a stock
//
// 	GetStockStats("CW8.PA", startDate, endDate)
//
// return the stock stats
func (esStock *Stock) GetStockStats(symbol string, startDate time.Time, endDate time.Time) (*StocksStats, error) {
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	query := elastic.NewQueryStringQuery(fmt.Sprintf("symbol = %s AND date: [%s TO %s]",
		symbol, startDate.Format(finance.DateFormat), endDate.Format(finance.DateFormat)))
	statsAgg := elastic.NewExtendedStatsAggregation().Field("close")
	results, err := esStock.es.Search(indexName).
		Type(indexType).
		Query(query).
		Aggregation(statsAggregationName, statsAgg).
		Size(0).
		Do(esContext)
	if err != nil {
		return nil, err
	}
	resAgg, _ := results.Aggregations.ExtendedStats(statsAggregationName)
	return &StocksStats{
		Symbol:            symbol,
		StandardDeviation: *resAgg.StdDeviation,
		Avg:               *resAgg.Avg,
	}, nil
}

// GetDateForNumPoint compute the start date to get a number of data points
//
// 	GetDateForNumPoint("CW8.PA", endDate)
//
// returns the start date
func (esStock *Stock) GetDateForNumPoint(symbol string, numPoints int, endDate time.Time) (*time.Time, error) {
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	startDate := endDate.AddDate(0, 0, int(float64(numPoints)*2)*-1)
	query := elastic.NewQueryStringQuery(fmt.Sprintf("symbol = %s AND date: [%s TO %s]",
		symbol, startDate.Format(finance.DateFormat), endDate.Format(finance.DateFormat)))
	results, err := esStock.es.Search(indexName).
		Type(indexType).
		Query(query).
		Sort("date", false).
		From(numPoints - 1).
		Size(1).
		Do(esContext)
	if err != nil {
		return nil, err
	}
	var value stockValue
	err = json.Unmarshal(*results.Hits.Hits[0].Source, &value)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(time.RFC3339, value.Date)
	if err != nil {
		return nil, err
	}
	return &date, nil
}
