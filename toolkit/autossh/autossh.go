package main

import (
	"os"
	"os/signal"
	"syscall"

	"toolkit/logs"
)

func autossh() (err error) {
	if err = readConfInFile(); err != nil {
		return
	}
	if err = Startup(); err != nil {
		return
	} else {
		defer Shutdown()
	}

	// block and listen for signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	s := <-sig
	logs.Infof("recv signal: %s", s.String())
	return
}
