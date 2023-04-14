// Package cmd define the main entry point of an application, as well as any command-line utilities or tools that are part of the application.
package cmd

import "flag"

type flagValues struct {
	port string
}

// Tools is an interface that extend tools
type Tools interface {
	Flags() flagValues
}

type tools struct{}

var _ Tools = (*tools)(nil)

// NewTools is a constructor for tools
func NewTools() Tools {
	return &tools{}
}

func (*tools) Flags() flagValues {
	// run app
	var port string
	flag.StringVar(&port, "p", "8080", "port for http server")
	flag.Parse()

	fv := &flagValues{
		port: port,
	}

	return *fv
}
