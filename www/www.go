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
	logger = logger.With("remote_addr", req.RemoteAddr)

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
