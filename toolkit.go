package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"toolkit/system"
	"toolkit/tools"

	"github.com/valord577/clix"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	system.StructuredLogging()

	undo, err := maxprocs.Set()
	if err != nil {
		undo()
		slog.Debug("set maxprocs, errmsg: " + err.Error())
	}
}

const (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

func main() {
	exitCode := EXIT_SUCCESS
	defer func() { os.Exit(exitCode) }()

	if err := exec(); err != nil {
		exitCode = EXIT_FAILURE
		slog.Error(err.Error())
	}
}

func exec() error {
	cmd := &clix.Command{
		Name: filepath.Base(os.Args[0]),

		Run: func(*clix.Command, []string) (err error) {
			_, err = os.Stderr.WriteString(system.Version() + "\n")
			return err
		},
	}

	cmd.AddCmd(
		tools.AutoIp,
	)
	return cmd.Execute()
}
