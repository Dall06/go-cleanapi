package cmd

import (
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
}