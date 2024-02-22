package www

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
)

type OutboxGetHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	PostsDatabase    activitypub.PostsDatabase
	URIs             *activitypub.URIs
}

func OutboxGetHandler(opts *OutboxGetHandlerOptions) (http.Handler, error) {

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

		outbox_url := acct.OutboxURL(ctx, opts.URIs)

		col := &ap.OrderedCollection{
			Context: []string{
				"https://www.w3.org/ns/activitystreams",
			},
			Id:         outbox_url.String(),
			Type:       "OrderedCollection",
			TotalItems: 0,
		}

		rsp.Header().Set("Content-type", ap.ACTIVITY_CONTENT_TYPE)

		enc := json.NewEncoder(rsp)
		err = enc.Encode(col)

		if err != nil {
			logger.Error("Failed to encode collection resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
