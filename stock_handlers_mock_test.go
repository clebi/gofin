package main

import (
	"errors"
	finance "github.com/clebi/yfinance"
	"github.com/stretchr/testify/mock"
	"time"
)

type mockHistoryAPI struct {
	mock.Mock
}

func (mock *mockHistoryAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	args := mock.Called(symbol, start, end)
	stocks := args.Get(0).([]finance.Stock)
	return stocks, args.Error(1)
}

type mockEsStock struct {
	mock.Mock
}

func (mock *mockEsStock) Index(stock finance.Stock) error {
	args := mock.Called(stock)
	return args.Error(0)
}

func (mock *mockEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]EsStocksAgg, error) {
	args := mock.Called(symbol, movAvgWindow, step, startDate, endDate)
	stocks := args.Get(0).([]EsStocksAgg)
	return stocks, args.Error(1)
}

type DummySchemaDecoder struct {
}

func (decoder *DummySchemaDecoder) Decode(dst interface{}, src map[string][]string) error {
	return nil
}

type ErrorSchemaDecoder struct {
	Msg string
}

func (decoder *ErrorSchemaDecoder) Decode(dst interface{}, src map[string][]string) error {
	return errors.New(decoder.Msg)
}

type DummyStructValidator struct {
}

func (validator *DummyStructValidator) Struct(s interface{}) error {
	return nil
}

type ErrorStructValidator struct {
	Msg string
}

func (validator *ErrorStructValidator) Struct(s interface{}) error {
	return errors.New(validator.Msg)
}

type ErrorFinanceAPI struct {
	Msg string
}

func (api *ErrorFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return nil, errors.New(api.Msg)
}

type DummyFinanceAPI struct {
}

func (api *DummyFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return []finance.Stock{}, nil
}

type OneItemFinanceAPI struct {
}

func (api *OneItemFinanceAPI) GetHistory(symbol string, start time.Time, end time.Time) ([]finance.Stock, error) {
	return []finance.Stock{{Symbol: "TEST"}}, nil
}

type ErrorEsStock struct {
	Msg string
}

func (es *ErrorEsStock) Index(stock finance.Stock) error {
	return errors.New(es.Msg)
}

func (es *ErrorEsStock) GetStocksAgg(symbol string, movAvgWindow int, step int, startDate time.Time, endDate time.Time) ([]EsStocksAgg, error) {
	return nil, errors.New(es.Msg)
}
