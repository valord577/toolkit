package system

import (
	"os"
	"strconv"
	"strings"
)

func Hostname() string {
	hostname := GetEnvString("TOOLKIT_SYS_HOSTNAME")
	if len(hostname) < 1 {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "development"
		}
	}
	return hostname
}

func GetEnvString(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
func GetEnvInt(key string) int {
	env := strings.TrimSpace(os.Getenv(key))
	value, _ := strconv.ParseInt(env, 10, 32)
	return int(value)
}
func GetEnvBool(key string) bool {
	env := strings.TrimSpace(os.Getenv(key))
	value, _ := strconv.ParseBool(env)
	return value
}
