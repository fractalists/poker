package main

import (
	"net/http"
	"os"
	"poker/internal/api"
	"poker/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		panic(err)
	}

	level, err := logrus.ParseLevel(opts.logLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)

	manager := service.NewManager()
	if opts.dataFile != "" {
		persistentManager, err := service.NewPersistentManager(opts.dataFile)
		if err != nil {
			panic(err)
		}
		manager = persistentManager
	}
	handler := newRootHandler(api.NewServer(manager), opts.webDist)
	if err := http.ListenAndServe(opts.addr, handler); err != nil {
		panic(err)
	}
}
