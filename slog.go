package activitypub

import (
	"io"
	"log/slog"
	"os"
)

func init() {
	// SetupLogger()
}

func SetupLogger() *slog.Logger {
	return SetupLoggerWithWriter(os.Stdout)
}

func SetupLoggerWithWriter(wr io.Writer) *slog.Logger {

	opts := &slog.HandlerOptions{}

	handler := slog.NewTextHandler(wr, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
	return logger
}
