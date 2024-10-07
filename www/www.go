package www

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/sfomuseum/go-activitypub/ap"
)

func LoggerWithRequest(req *http.Request, logger *slog.Logger) *slog.Logger {

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With("method", req.Method)
	logger = logger.With("user agent", req.Header.Get("User-Agent"))	
	logger = logger.With("accept", req.Header.Get("Accept"))
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote addr", req.RemoteAddr)
	logger = logger.With("user ip", ReadUserIP(req))

	return logger
}

func IsActivityStreamRequest(req *http.Request, header string) bool {

	raw := req.Header.Get(header)
	accept := strings.Split(raw, ",")

	is_activitystream := false

	for _, h := range accept {

		h = strings.TrimSpace(h)

		switch h {

		case ap.ACTIVITYSTREAMS_ACCEPT_HEADER:
			is_activitystream = true
			break
		case ap.ACTIVITY_CONTENT_TYPE:
			is_activitystream = true
			break
		default:
			continue
		}
	}

	return is_activitystream
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
