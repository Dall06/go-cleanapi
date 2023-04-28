// Package tools define the main entry point of an application, as well as any command-line utilities or tools that are part of the application.
package tools

import (
	"flag"
	"os"
)

// FlagValues means the values obtained as flag parameters of cli
type FlagValues struct {
	Port    string
	Version string
}

// Flags is an interface that extend tools
type Flags interface {
	GetFlags() (*FlagValues, error)
}

type flags struct {
	flagSet *flag.FlagSet
}

var _ Flags = (*flags)(nil)

// NewFlags is a constructor for tools
func NewFlags() Flags {
	return &flags{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
}

func (f *flags) GetFlags() (*FlagValues, error) {
	// run app
	var (
		port    string
		version string
	)
	f.flagSet.StringVar(&port, "p", "8080", "port for http server")
	f.flagSet.StringVar(&version, "v", "0.0.0", "version for http server")

	if err := f.flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	fv := &FlagValues{
		Port:    port,
		Version: version,
	}

	return fv, nil
}
