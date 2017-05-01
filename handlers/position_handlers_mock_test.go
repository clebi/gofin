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

import (
	"errors"

	"github.com/clebi/gofin/es"
	finance "github.com/clebi/yfinance"
	"github.com/labstack/echo"
)

type DummyEsPosition struct {
	PositionAgg []es.PositionAgg
}

func (posStock *DummyEsPosition) AddPosition(position *es.Position) error {
	return nil
}

func (posStock *DummyEsPosition) GetPositions(username string) ([]es.PositionAgg, error) {
	return posStock.PositionAgg, nil
}

type ErrorEsPosition struct {
	Msg string
}

func (posStock *ErrorEsPosition) AddPosition(position *es.Position) error {
	return errors.New(posStock.Msg)
}

func (posStock *ErrorEsPosition) GetPositions(username string) ([]es.PositionAgg, error) {
	return nil, errors.New(posStock.Msg)
}

type ErrorEchoBind struct {
	echo.Context
	Msg string
}

func (echo ErrorEchoBind) Bind(interface{}) error {
	return errors.New(echo.Msg)
}

type DummyEchoBind struct {
	echo.Context
}

func (echo DummyEchoBind) Bind(interface{}) error {
	return nil
}

type DummyQuotesAPI struct {
	quote finance.Quote
}

func (quotes DummyQuotesAPI) GetQuote(symbol string) (*finance.Quote, error) {
	return &quotes.quote, nil
}

type ErrorQuotesAPI struct {
	Msg string
}

func (quotes ErrorQuotesAPI) GetQuote(symbol string) (*finance.Quote, error) {
	return nil, errors.New(quotes.Msg)
}
