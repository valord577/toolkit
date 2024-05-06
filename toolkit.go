package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"toolkit/cmds"
	"toolkit/logs"
	"toolkit/system"

	"github.com/valord577/clix"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	log.SetOutput(io.Discard)

	undo, err := maxprocs.Set(maxprocs.Logger(logs.Debugf))
	if err != nil {
		undo()
		logs.Warnf("set maxprocs, err: %s", err.Error())
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
		logs.Errorf("%s", err.Error())
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

	cmds := []*clix.Command{
		cmds.AutoIp, cmds.AutoSSH,
	}
	if err := cmd.AddCmd(cmds...); err != nil {
		return err
	}
	return cmd.Execute()
}
