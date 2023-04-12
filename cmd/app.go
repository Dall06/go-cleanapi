//go:build !coverage
// +build !coverage

// it is just an implementation, the test of the app ocurrs on integration testing in server layer

package cmd

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/server"
	"dall06/go-cleanapi/utils"
	"errors"

	"github.com/go-playground/validator/v10"
)

type app struct {
	port string
}

func NewApp(p string) app {
	return app{
		port: p,
	}
}

// Load app configuration such as servers, cache and utils
func (a app) Main() error {
	conf := config.NewConfig(a.port)
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

	u := utils.NewUUIDGenerator()
	if u == nil {
		return errors.New("empty uid generator repo")
	}

	v := validator.New()
	if v == nil {
		return errors.New("empty validator repo")
	}

	s := server.NewServer(l, jwt, u, *v)
	s.Start()
	return nil
}