package www

import (
	"net/http"
	"time"

	"github.com/aaronland/go-http/v3/slog"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/uris"
)

type OutboxPostHandlerOptions struct {
	AccountsDatabase database.AccountsDatabase
	PostsDatabase    database.PostsDatabase
	URIs             *uris.URIs
}

func OutboxPostHandler(opts *OutboxPostHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		t1 := time.Now()

		defer func() {
			logger.Info("Time to serve request", "ms", time.Since(t1).Milliseconds())
		}()

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !IsActivityStreamRequest(req, "Content-Type") {
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
