package www

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
)

type FollowingHandlerOptions struct {
	AccountsDatabase  activitypub.AccountsDatabase
	FollowingDatabase activitypub.FollowingDatabase
	URIs              *activitypub.URIs
}

func FollowingHandler(opts *FollowingHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

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

		if !IsActivityStreamRequest(req, "Accept") {
			logger.Error("Not activitystream request")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		account_name, host, err := activitypub.ParseAddressFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse address from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account name", account_name)

		if host != "" && host != opts.URIs.Hostname {
			logger.Error("Resouce has bunk hostname", "host", host)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

		if err != nil {

			logger.Error("Failed to retrieve account", "error", err)

			if err == activitypub.ErrNotFound {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger = logger.With("account id", acct.Id)

		resource, err := acct.FollowingResource(ctx, opts.URIs, opts.FollowingDatabase)

		if err != nil {
			logger.Error("Failed to create following resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", ap.ACTIVITY_CONTENT_TYPE)

		enc := json.NewEncoder(rsp)
		err = enc.Encode(resource)

		if err != nil {
			logger.Error("Failed to encode following resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
