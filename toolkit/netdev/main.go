package main

import (
	"os"

	log "toolkit/logger"
	"toolkit/system"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

func main() {
	exitCode := EXIT_SUCCESS
	defer os.Exit(exitCode)

	log.Init()
	defer func() {
		_ = log.Sync()
	}()
	log.Infof("%s", system.Version())

	if err := notifyInterfaces(); err != nil {
		exitCode = EXIT_FAILURE
		log.Errorf("%s", err.Error())
	}
}
