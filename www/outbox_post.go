package www

import (
	"net/http"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

type OutboxPostHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	PostsDatabase    activitypub.PostsDatabase
	URIs             *activitypub.URIs
}

func OutboxPostHandler(opts *OutboxPostHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := LoggerWithRequest(req, nil)

		t1 := time.Now()

		defer func() {
			logger.Info("Time to serve request", "ms", time.Since(t1).Milliseconds())
		}()

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !IsActivityStreamRequest(req) {
			logger.Error("Not activitystream request")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger.Error("Forbidden")
		http.Error(rsp, "Forbidden", http.StatusForbidden)
		return
	}

	return http.HandlerFunc(fn), nil
}