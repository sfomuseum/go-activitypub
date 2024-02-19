package slog

import (
	"io"
	go_slog "log/slog"
	"os"
)

func Default() *go_slog.Logger {
	return DefaultWithWriter(os.Stderr)
}

func DefaultWithWriter(wr io.Writer) *go_slog.Logger {
	handler := go_slog.NewJSONHandler(wr, nil)
	return go_slog.New(handler)
}
