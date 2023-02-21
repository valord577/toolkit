package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

const (
	toolkitHostname = "TOOLKIT_HOSTNAME"
)

func Device() string {
	hostname := os.Getenv(toolkitHostname)
	hostname = strings.TrimSpace(hostname)
	if len(hostname) < 1 {
		hostname = "devel"
	}
	return hostname
}

func Period() time.Duration {
	return time.Duration(c.Netdev.Period) * time.Second
}

func Receiver() string {
	return c.Netdev.Receiver
}

func SmtpCnf() smtp {
	return smtp{
		Host: c.Smtp.Host,
		Port: c.Smtp.Port,
		User: c.Smtp.User,
		Pass: c.Smtp.Pass,

		SslOnConnect: c.Smtp.SslOnConnect,
	}
}

var c appJsonc

const (
	toolkitConfigFilepath = "TOOLKIT_CONFIG_PATH"
)

func ReadInFile() error {
	fp := os.Getenv(toolkitConfigFilepath)
	return readInFile(fp)
}

func readInFile(file string) error {
	bs, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	jsonBs := ignoreComments(bs)
	return json.Unmarshal(jsonBs, &c)
}
