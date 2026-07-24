package slog

import (
	"log/slog"
	"net/http"
	"sync"
)

var default_logger_attrs = new(sync.Map)

// LoggerWithRequest returns a `slog.Logger` instance with the following keys set: (HTTP) method, user agent, path, remote addr and user ip.
func LoggerWithRequest(req *http.Request, logger *slog.Logger) *slog.Logger {

	attrs_map := map[string]any{
		"method":      req.Method,
		"user agent":  req.Header.Get("User-Agent"),
		"accept":      req.Header.Get("Accept"),
		"path":        req.URL.Path,
		"remote addr": req.RemoteAddr,
		"user ip":     ReadUserIP(req),
	}

	attrs := make([]any, 0)

	if logger == nil {

		logger = slog.Default()

		for k, v := range attrs_map {

			_, exists := default_logger_attrs.LoadOrStore(k, v)

			if !exists {
				attrs = append(attrs, k, v)
			}
		}

	} else {

		for k, v := range attrs_map {
			attrs = append(attrs, k, v)
		}
	}

	return logger.With(attrs...)
}

// ReadUserIP returns the value of the `X-Real-Ip` or `X-Forwarded-For` headers (in that order) or the default remote address reported by 'req'.
func ReadUserIP(req *http.Request) string {

	addr := req.Header.Get("X-Real-Ip")

	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}

	if addr == "" {
		addr = req.RemoteAddr
	}

	return addr
}
