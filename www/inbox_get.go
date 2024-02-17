package www

import (
	"net/http"
)

type InboxGetHandlerOptions struct {
}

func InboxGetHandler(opts *InboxGetHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		// ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		http.Error(rsp, "Not implemented", http.StatusNotImplemented)
		return
	}

	return http.HandlerFunc(fn), nil
}
