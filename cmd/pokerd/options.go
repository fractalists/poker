package main

import (
	"flag"
)

type options struct {
	addr     string
	logLevel string
	webDist  string
}

func parseOptions(args []string) (options, error) {
	fs := flag.NewFlagSet("pokerd", flag.ContinueOnError)
	addr := fs.String("addr", "127.0.0.1:8080", "http listen address")
	logLevel := fs.String("log-level", "info", "debug|info|warn|error")
	webDist := fs.String("web-dist", "web/dist", "compiled web application directory")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}

	return options{
		addr:     *addr,
		logLevel: *logLevel,
		webDist:  *webDist,
	}, nil
}
