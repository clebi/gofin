package main

import "errors"

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
