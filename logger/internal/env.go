package internal

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

// text (default) | json
func newZapEncoderFunc() func(zapcore.EncoderConfig) zapcore.Encoder {
	env := os.Getenv("TOOLKIT_LOGS_FORMAT")
	format := strings.TrimSpace(env)
	switch format {
	case "json":
		return zapcore.NewJSONEncoder
	default:
		return zapcore.NewConsoleEncoder
	}
}

// development trace
func isDebug() bool {
	env := os.Getenv("TOOLKIT_LOGS_DEBUG")
	debug, _ := strconv.ParseBool(env)
	return debug
}

// Golang style time format template string.
// Default: "2006-01-02 15:04:05.000 -07:00"
func tmfmtLayout() string {
	env := os.Getenv("TOOLKIT_LOGS_TIME_FORMAT")
	layout := strings.TrimSpace(env)
	if len(layout) < 1 {
		layout = "2006-01-02 15:04:05.000 -07:00"
	}
	return layout
}
