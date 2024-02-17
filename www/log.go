package www

import (
	"log/slog"
	"net/http"
)

func LoggerWithRequest(req *http.Request, logger *slog.Logger) *slog.Logger {

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With("method", req.Method)
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote_addr", req.RemoteAddr)

	return logger
}
