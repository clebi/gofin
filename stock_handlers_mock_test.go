package main

import (
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
