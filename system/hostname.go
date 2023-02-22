package system

import (
	"os"
	"strings"

	log "toolkit/logger"
)

const (
	toolkitHostname = "TOOLKIT_HOSTNAME"
)

func Hostname() string {
	hostname := os.Getenv(toolkitHostname)
	hostname = strings.TrimSpace(hostname)
	if len(hostname) < 1 {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			log.Warnf("os hostname, err: %s", err.Error())
			hostname = "default"
		}
	}
	return hostname
}
