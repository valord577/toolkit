package system

import (
	"os"
	"strings"
	"time"
)

const (
	sysHostname = "TOOLKIT_SYS_HOSTNAME"
	sysTimeZone = "TOOLKIT_SYS_TIME_ZONE"
)

func Hostname() string {
	hostname := os.Getenv(sysHostname)
	hostname = strings.TrimSpace(hostname)
	if len(hostname) < 1 {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "default"
		}
	}
	return hostname
}

func TimeZone() *time.Location {
	env := os.Getenv(sysTimeZone)
	env = strings.TrimSpace(env)
	loc, err := time.LoadLocation(env)
	if err != nil {
		loc = time.Local
	}
	return loc
}
