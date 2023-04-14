//go:build !coverage
// +build !coverage

// Package cmd define the main entry point of an application, as well as any command-line utilities or tools that are part of the application.
// it is just an implementation, the test of the app ocurrs on integration testing in server layer
package cmd

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/server"
	"dall06/go-cleanapi/utils"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// App is an interface that extends app
type App interface {
	Main() error
}

var _ App = (*app)(nil)

type app struct {
	tools Tools
}

// NewApp is a constructor for app
func NewApp(t Tools) App {
	return &app{
		tools: t,
	}
}

// Main app configuration such as servers, cache and utils
func (a *app) Main() error {
	flagValues := a.tools.Flags()
	prt := flagValues.port
	ver := flagValues.version

	conf := config.NewConfig(prt, ver)
	err := conf.SetConfig()
	if err != nil {
		return err
	}

	jwt := utils.NewJWT()
	if jwt == nil {
		return errors.New("empty jwt repo")
	}

	l := utils.NewLogger()
	if l == nil {
		return errors.New("empty logger repo")
	}
	err = l.Initialize()
	if err != nil {
		return fmt.Errorf("error when init logger %v: ", err)
	}

	u := utils.NewUUIDGenerator()
	if u == nil {
		return errors.New("empty uid generator repo")
	}

	v := validator.New()
	if v == nil {
		return errors.New("empty validator repo")
	}

	s := server.NewServer(l, jwt, u, *v)
	if err := s.Start(); err != nil {
		return fmt.Errorf("error when starting the server %v: ", err)
	}

	return nil
}
