// +build !coverage

package utils

import "fmt"


type mockLoggerRepository struct {
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

func NewMockLoggerRepository() Logger {
	return mockLoggerRepository{}
}

func (l mockLoggerRepository) Initialize() error {
    return nil
}

func (l mockLoggerRepository) Warn(message string, args ...interface{}) {
    l.WarnCalled = true
    l.WarnMsg = message
    l.WarnArgs = args
	fmt.Println(l.WarnCalled, l.WarnMsg, l.WarnArgs)
}

func (l mockLoggerRepository) Info(message string, args ...interface{}) {
    l.InfoCalled = true
    l.InfoMsg = message
    l.InfoArgs = args
	fmt.Println(l.InfoCalled, l.InfoMsg, l.InfoArgs)
}

func (l mockLoggerRepository) Error(message string, args ...interface{}) {
    l.ErrorCalled = true
    l.ErrorMsg = message
    l.ErrorArgs = args
	fmt.Println(l.ErrorCalled, l.ErrorMsg, l.ErrorArgs)
}
