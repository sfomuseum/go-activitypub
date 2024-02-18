package handler

import (
	"net/http"
)

func NotImplementedHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		http.Error(rsp, "Not Implemented", http.StatusNotImplemented)
		return
	}

	return http.HandlerFunc(fn)
}
