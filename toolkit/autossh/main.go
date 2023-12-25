package main

import (
	"os"

	"toolkit/common"
	"toolkit/logs"
	"toolkit/system"
)

func init() {
	common.InitPreset()
}

const (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

func main() {
	exitCode := EXIT_SUCCESS
	defer func() { os.Exit(exitCode) }()
	logs.Infof("%s", system.Version())

	if err := autossh(); err != nil {
		exitCode = EXIT_FAILURE
		logs.Errorf("%s", err.Error())
	}
}
