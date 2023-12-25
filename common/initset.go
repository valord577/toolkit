package common

import (
	"io"
	"log"
	"os"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap/zapcore"

	"toolkit/logs"
	"toolkit/system"
)

func InitPreset() {
	log.SetOutput(logs.Writer(zapcore.WarnLevel))

	undo, err := maxprocs.Set(maxprocs.Logger(logs.Debugf))
	if err != nil {
		undo()
		logs.Warnf("set maxprocs, err: %s", err.Error())
	}

	if len(os.Args) > 1 {
		io.WriteString(os.Stderr, system.Version()+"\n")
		os.Exit(0)
	}
}
