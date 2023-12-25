package internal

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

// text (default) | json
const toolkitLogsFormat = "TOOLKIT_LOGS_FORMAT"

func newZapEncoderFunc() func(zapcore.EncoderConfig) zapcore.Encoder {
	env := os.Getenv(toolkitLogsFormat)
	env = strings.TrimSpace(env)
	switch env {
	case "json":
		return zapcore.NewJSONEncoder
	default:
		return zapcore.NewConsoleEncoder
	}
}

// development trace
const toolkitLogsDebug = "TOOLKIT_LOGS_DEBUG"

func isDebug() bool {
	env := os.Getenv(toolkitLogsDebug)
	debug, _ := strconv.ParseBool(env)
	return debug
}

// Golang style time format template string.
// Default: "2006-01-02 15:04:05.000 -07:00"
const toolkitLogsTimeFormat = "TOOLKIT_LOGS_TIME_FORMAT"

func timeFormat() string {
	env := os.Getenv(toolkitLogsTimeFormat)
	env = strings.TrimSpace(env)
	if len(env) < 1 {
		env = "2006-01-02 15:04:05.000 -07:00"
	}
	return env
}
