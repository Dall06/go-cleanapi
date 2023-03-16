package main

import (
	"dall06/go-cleanapi/cmd"
	"flag"
)

func main() {
	// run app
	var stage string
	flag.StringVar(&stage, "s", "dev", "stage of the app")
	flag.Parse()
	
	app := cmd.NewApp(stage)
	app.Run()
}