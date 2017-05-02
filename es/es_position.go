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
	"errors"
	"fmt"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	symbolsAggName = "symbols"
	numberAggName  = "number"
	costAggName    = "cost"
)

// Position contains all values representing a stock position
type Position struct {
	Username string    `json:"username" validate:"required"`
	Broker   string    `json:"broker" validate:"required"`
	Symbol   string    `json:"symbol" validate:"required"`
	Date     time.Time `json:"date,string" validate:"required"`
	Number   int       `json:"number,int" validate:"required"`
	Value    float64   `json:"value,float" validate:"gt=0"`
	Cost     float64   `json:"cost,float" validate:"required"`
}

func (position Position) String() string {
	return fmt.Sprintf("Username: %s Broker: %s Symbol: %s Date: %s Number: %d Value: %f Cost: %f",
		position.Username, position.Broker, position.Symbol, position.Date, position.Number, position.Value, position.Cost)
}

// PositionAgg contains the list of positions aggregation by symbol
type PositionAgg struct {
	Symbol string
	Number int
	Cost   float64
}

func (position PositionAgg) String() string {
	return fmt.Sprintf("Symbol: %s Number: %d Cost: %f", position.Symbol, position.Number, position.Cost)
}

// IPositionStock contains all es position stock actions
type IPositionStock interface {
	AddPosition(position *Position) error
	GetPositions(username string) ([]PositionAgg, error)
}

// PositionStock manage positons in elasticsearch
type PositionStock struct {
	es *elastic.Client
}

// NewPosition create a new elasticsearch poisitons manager
func NewPosition(es *elastic.Client) IPositionStock {
	return &PositionStock{
		es: es,
	}
}

// AddPosition adds a position into elasticsearch storage
//
//  AddPosition(position)
//
//  return the inserted position
func (posStock *PositionStock) AddPosition(position *Position) error {
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	positionMap := map[string]interface{}{
		"username": position.Username,
		"broker":   position.Broker,
		"date":     position.Date.Format(time.RFC3339),
		"symbol":   position.Symbol,
		"number":   position.Number,
		"value":    position.Value,
		"cost":     position.Cost,
	}
	_, err := posStock.es.Index().
		Index("stock-positions").
		Type("stock_position").
		Id(fmt.Sprintf("%s_%s_%s", position.Broker, position.Date.Format(time.RFC3339), position.Symbol)).
		BodyJson(positionMap).
		Do(esContext)
	if err != nil {
		return err
	}
	return nil
}

// GetPositions gets all the positions concerning a user
//
// GetPositions(username)
//
// return the list of positions
func (posStock *PositionStock) GetPositions(username string) ([]PositionAgg, error) {
	esContext, esCancel := context.WithTimeout(context.Background(), indexTimeout)
	defer esCancel()
	numberAgg := elastic.NewSumAggregation().Field("number")
	costAgg := elastic.NewSumAggregation().Field("cost")
	symbolsAgg := elastic.NewTermsAggregation().
		Field("symbol.keyword").
		SubAggregation(numberAggName, numberAgg).
		SubAggregation(costAggName, costAgg)
	query := elastic.NewQueryStringQuery(fmt.Sprintf("username = %s", username))
	results, err := posStock.es.Search("stock-positions").
		Type("stock_position").
		Query(query).
		Aggregation(symbolsAggName, symbolsAgg).
		Do(esContext)
	if err != nil {
		return nil, err
	}
	resAgg, _ := results.Aggregations.Terms(symbolsAggName)
	positions := make([]PositionAgg, len(resAgg.Buckets))
	for i, bucket := range resAgg.Buckets {
		number, _ := bucket.Sum(numberAggName)
		cost, _ := bucket.Sum(costAggName)
		key, ok := bucket.Key.(string)
		if !ok {
			return nil, errors.New("GetPositions: Bad aggregation key")
		}
		positions[i] = PositionAgg{
			Symbol: key,
			Number: int(*number.Value),
			Cost:   *cost.Value,
		}
	}
	return positions, nil
}
