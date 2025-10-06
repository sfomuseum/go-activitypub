package handlers

import (
	"net/http"
)

func PingPongHandler() (http.Handler, error) {
	return PingHandler("PONG")
}

func PingHandler(response string) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-Type", "text/plain")

		if response != "" {
			rsp.Header().Set("X-ping", response)
		}

		rsp.WriteHeader(http.StatusNoContent)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
