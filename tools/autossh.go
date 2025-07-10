package tools

import (
	"os"
	"os/signal"
	"syscall"

	"toolkit/logs"
	"toolkit/tools/autossh"

	"github.com/valord577/clix"
)

var AutoSSH = &clix.Command{
	Name: "autossh",

	Summary: "Service of SSH Tunnels",
	Run: func(*clix.Command, []string) (err error) {
		if err = autossh.ReadConfInFile(); err != nil {
			return
		}
		if err = autossh.Startup(); err != nil {
			return
		}
		defer autossh.Shutdown()

		// block and listen for signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		s := <-sig
		logs.Infof("recv signal: %s", s.String())
		return
	},
}
