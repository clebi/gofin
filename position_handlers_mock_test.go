package main

import (
	"errors"

	"github.com/labstack/echo"
)

type DummyEsPosition struct {
}

func (posStock *DummyEsPosition) AddPosition(position *Position) error {
	return nil
}

type ErrorEsPosition struct {
	Msg string
}

func (posStock *ErrorEsPosition) AddPosition(position *Position) error {
	return errors.New(posStock.Msg)
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
