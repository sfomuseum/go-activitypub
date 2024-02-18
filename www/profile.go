package www

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/sfomuseum/go-activitypub"
)

type ProfileHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
	Hostname         string
}

func ProfileHandler(opts *ProfileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !IsActivityStreamRequest(req) {
			logger.Error("Not activitystream request")
			http.Error(rsp, "Not implemented", http.StatusNotImplemented)
			return
		}

		// sudo make me a regexp or req.PathId(...)

		account_name := filepath.Base(req.URL.Path)

		logger = logger.With("account name", account_name)

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

		if err != nil {
			logger.Error("Failed to retrieve account", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		logger = logger.With("account id", acct.Id)

		// Check content-type here and HTML or JSON it up...

		profile, err := acct.ProfileResource(ctx, opts.Hostname, opts.URIs)

		if err != nil {
			logger.Error("Failed to derive profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/activity+json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(profile)

		if err != nil {
			logger.Error("Failed to encode profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
