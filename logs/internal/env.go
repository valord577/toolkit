package internal

import (
	"toolkit/system"

	"go.uber.org/zap/zapcore"
)

// text (default) | json
func newZapEncoderFunc() func(zapcore.EncoderConfig) zapcore.Encoder {
	format := system.GetEnvString("TOOLKIT_LOGS_FORMAT")
	switch format {
	case "json":
		return zapcore.NewJSONEncoder
	default:
		return zapcore.NewConsoleEncoder
	}
}

// development trace
func isDebug() bool {
	return system.GetEnvBool("TOOLKIT_LOGS_DEBUG")
}

// Golang style time format template string.
// Default: "2006-01-02 15:04:05.000 -07:00"
func tmfmtLayout() string {
	layout := system.GetEnvString("TOOLKIT_LOGS_TIME_FORMAT")
	if len(layout) < 1 {
		layout = "2006-01-02 15:04:05.000 -07:00"
	}
	return layout
}
