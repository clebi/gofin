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
	"github.com/labstack/echo"
)

type DummyEsPosition struct {
}

func (posStock *DummyEsPosition) AddPosition(position *es.Position) error {
	return nil
}

type ErrorEsPosition struct {
	Msg string
}

func (posStock *ErrorEsPosition) AddPosition(position *es.Position) error {
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
