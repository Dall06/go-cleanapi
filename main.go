// Package main runs main thread of the app
package main

import (
	"dall06/go-cleanapi/cmd"
)

// @title go-cleanapi
// @description Golang REST Api based on Uncle's Bob Clean Arch
// @version 1.0.0
// @host localhost:8080
// @BasePath /go-cleanapi/api/v1
func main() {
	app := cmd.NewApp()

	err := app.Main()
	if err != nil {
		panic(err)
	}
}
