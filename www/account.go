package www

import (
	"encoding/json"
	"net/http"
	_ "path/filepath"

	"github.com/sfomuseum/go-activitypub"
)

type AccountHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
}

func AccountHandler(opts *AccountHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		account_name, _, err := activitypub.ParseAddressFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse address from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account name", account_name)

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

		// Check content-type here and HTML or JSON it up...

		if IsActivityStreamRequest(req) {

			profile, err := acct.ProfileResource(ctx, opts.URIs)

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

		rsp.Header().Set("Content-type", "text/html")

		rsp.Write([]byte(account_name))
		return
	}

	return http.HandlerFunc(fn), nil
}
