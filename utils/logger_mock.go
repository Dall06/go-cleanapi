//go:build !coverage
// +build !coverage

// Package utils is a package that provides general method for the api usage
package utils

import "fmt"

type loggerMock struct {
	WarnCalled bool
	WarnMsg    string
	WarnArgs   []interface{}

	InfoCalled bool
	InfoMsg    string
	InfoArgs   []interface{}

	ErrorCalled bool
	ErrorMsg    string
	ErrorArgs   []interface{}
}

// NewLoggerMock is a consturcotr to genrate a logger mock
func NewLoggerMock() Logger {
	return &loggerMock{}
}

func (l loggerMock) Initialize() error {
	return nil
}

func (l loggerMock) Warn(message string, args ...interface{}) {
	l.WarnCalled = true
	l.WarnMsg = message
	l.WarnArgs = args
	fmt.Println(l.WarnCalled, l.WarnMsg, l.WarnArgs)
}

func (l loggerMock) Info(message string, args ...interface{}) {
	l.InfoCalled = true
	l.InfoMsg = message
	l.InfoArgs = args
	fmt.Println(l.InfoCalled, l.InfoMsg, l.InfoArgs)
}

func (l loggerMock) Error(message string, args ...interface{}) {
	l.ErrorCalled = true
	l.ErrorMsg = message
	l.ErrorArgs = args
	fmt.Println(l.ErrorCalled, l.ErrorMsg, l.ErrorArgs)
}
