package system

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	version  = "dev"
	datetime = "-"
	flavor   = "default"
)

func Version() string {
	return filepath.Base(os.Args[0]) + " " + version + "." + flavor + " " + datetime +
		" " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
}
