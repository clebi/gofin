package es

import (
	"context"
	"fmt"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Position contains all values representing a stock position
type Position struct {
	Username string    `json:"username" validate:"required"`
	Broker   string    `json:"broker" validate:"required"`
	Symbol   string    `json:"symbol" validate:"required"`
	Date     time.Time `json:"date,string" validate:"required"`
	Number   int       `json:"number,int" validate:"gt=0"`
	Value    float64   `json:"value,float" validate:"gt=0"`
	Cost     float64   `json:"cost,float" validate:"gt=0"`
}

// IEsPositionStock contains all es position stock actions
type IEsPositionStock interface {
	AddPosition(position *Position) error
}

// EsPositionStock manage positons in elasticsearch
type EsPositionStock struct {
	es *elastic.Client
}

// NewEsPosition create a new elasticsearch poisitons manager
func NewEsPosition(es *elastic.Client) IEsPositionStock {
	return &EsPositionStock{
		es: es,
	}
}

// AddPosition adds a position into elasticsearch storage
//
//  AddPosition(position)
//
//  return the inserted position
func (posStock *EsPositionStock) AddPosition(position *Position) error {
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
