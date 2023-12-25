package logs

import (
	"bytes"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Writer(lvl zapcore.Level) io.Writer {
	log := WithOpts(zap.AddCallerSkip(2))
	if log == nil {
		return io.Discard
	}
	return &writer{log: log, lvl: lvl}
}

type writer struct {
	log *Logger
	lvl zapcore.Level
}

func (w *writer) Write(bs []byte) (n int, err error) {
	bs = bytes.TrimRightFunc(bs, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	w.log.Log(w.lvl, string(bs))
	return
}
