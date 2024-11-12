package handler

import (
	"net/http"
)

// DisabledHandler is a middleware handler that returns an HTTP 503 (Service unavailable) error
// if 'disabled' is true. Otherwise it serves 'next'.
func DisabledHandler(disabled bool, next http.Handler) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if disabled {
			http.Error(rsp, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(rsp, req)
	}

	return http.HandlerFunc(fn)
}
