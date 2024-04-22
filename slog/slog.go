package slog

import (
	"io"
	go_slog "log/slog"
	"os"
)

var logLevel = new(go_slog.LevelVar)

func EnableVerboseLogging() {
	logLevel.Set(go_slog.LevelDebug)
}

func DisableVerboseLogging() {
	logLevel.Set(go_slog.LevelInfo)
}

func ConfigureLogger(logger *go_slog.Logger, verbose bool) {

	go_slog.SetDefault(logger)

	if verbose {
		EnableVerboseLogging()
		logger.Debug("Verbose logging enabled")
	}
}

func Default() *go_slog.Logger {
	return DefaultWithWriter(os.Stderr)
}

func DefaultWithWriter(wr io.Writer) *go_slog.Logger {

	opts := &go_slog.HandlerOptions{
		Level: logLevel,
	}

	// handler := go_slog.NewJSONHandler(wr, opts)
	handler := go_slog.NewTextHandler(wr, opts)

	return go_slog.New(handler)
}
