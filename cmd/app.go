package cmd

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/server"
	"dall06/go-cleanapi/utils"
	"time"

	"github.com/patrickmn/go-cache"
)

type app struct {
	stage string
}

func NewApp(s string) *app {
	return &app{
		stage: s,
	}
}

// Load app configuration such as servers, cache and utils
func (a *app) Run() {
	// config
	env := config.NewEnv(a.stage)
	env.LoadStrings()

	mycache := cache.New(5*time.Minute, 10*time.Minute)

	logger := utils.NewLogger(a.stage)
	responses := utils.NewResponsesUtils()
	// server api rest
	s := server.NewServer(*mycache, logger, responses)
	s.Start()
}