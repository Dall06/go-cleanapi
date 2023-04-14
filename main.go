// Package main runs main thread of the app
package main

import (
	"dall06/go-cleanapi/cmd"
)

func main() {
	tools := cmd.NewTools()
	app := cmd.NewApp(tools)

	err := app.Main()
	if err != nil {
		panic(err)
	}
}
