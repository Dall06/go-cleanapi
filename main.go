package main

import (
	"dall06/go-cleanapi/cmd"
	"flag"
)

func main() {
	// run app
	var port string
	flag.StringVar(&port, "p", "8080", "port for http server")
	flag.Parse()
	
	app := cmd.NewApp(port)
	
	err := app.Main()
	if err != nil {
		panic(err)
	}
}