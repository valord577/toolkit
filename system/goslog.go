package system

import (
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func StructuredLogging() {
	slog.SetDefault(slog.New(newLogHandler(os.Stderr)))
}

// text (default) | json
func newLogHandler(w io.Writer) slog.Handler {
	layout := timeLayout()
	replace := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format(layout))
			}
		case slog.LevelKey:
			if l, ok := a.Value.Any().(slog.Level); ok {
				if l == 0 {
					a.Value = slog.StringValue("+0")
				} else if l > 0 {
					a.Value = slog.StringValue("+" + strconv.FormatInt(int64(l), 10))
				} else {
					a.Value = slog.StringValue(strconv.FormatInt(int64(l), 10))
				}
			}
		}
		return a
	}
	opt := &slog.HandlerOptions{
		AddSource: false, Level: slog.LevelInfo, ReplaceAttr: replace,
	}
	if isDebug() {
		opt.Level = slog.LevelDebug
	}

	switch GetEnvString("TOOLKIT_LOGS_FORMAT") {
	case "json":
		return slog.NewJSONHandler(w, opt)
	default:
		return slog.NewTextHandler(w, opt)
	}
}

// development trace
func isDebug() bool {
	return GetEnvBool("TOOLKIT_LOGS_DEBUG")
}

// Golang style time format template string.
// Default: "2006-01-02 15:04:05.000 -07:00"
func timeLayout() string {
	layout := GetEnvString("TOOLKIT_LOGS_TIME_FORMAT")
	if layout == "" {
		layout = "2006-01-02 15:04:05.000 -07:00"
	}
	return layout
}
