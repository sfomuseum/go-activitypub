package www

import (
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/ap"
)

func LoggerWithRequest(req *http.Request, logger *slog.Logger) *slog.Logger {

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With("method", req.Method)
	logger = logger.With("accept", req.Header.Get("Accept"))
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote addr", req.RemoteAddr)
	logger = logger.With("user ip", ReadUserIP(req))

	return logger
}

func IsActivityStreamRequest(req *http.Request) bool {

	switch req.Header.Get("Accept") {

	case ap.ACTIVITYSTREAMS_ACCEPT_HEADER:
		return true
	case ap.ACTIVITY_CONTENT_TYPE:
		return true
	default:
		return false
	}

}

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
