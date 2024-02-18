package handler

import (
	"net/http"
)

func NullHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		return
	}

	return http.HandlerFunc(fn)
}
