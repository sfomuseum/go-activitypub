package www

import (
	"encoding/json"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/sfomuseum/go-activitypub"
)

type WebfingerHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
	Hostname         string
}

func WebfingerHandler(opts *WebfingerHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		if req.Method != http.MethodGet {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resource, err := sanitize.GetString(req, "resource")

		if err != nil {
			logger.Error("Failed to derive ?resource= parameter", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if resource == "" {
			logger.Error("Empty ?resource= parameter")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("resource", resource)

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, resource)

		if err != nil {
			logger.Error("Failed to retrieve account for resource", "error", err)

			if err == activitypub.ErrNotFound {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger = logger.With("account id", acct.Id)

		wf, err := acct.WebfingerResource(ctx, opts.Hostname, opts.URIs)

		if err != nil {
			logger.Error("Failed to derive webfinger response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/jrd+json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(wf)

		if err != nil {
			logger.Error("Failed to encode webfinger response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
