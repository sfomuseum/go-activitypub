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

		// REVISIT ALL OF THIS...

		/*
			resource_id, resource_hostname, err := activitypub.ParseAccountURI(resource)

			if err != nil {
				logger.Error("Failed to parse resource", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			if resource_hostname != opts.Hostname {
				logger.Error("Resource lookup for bunk hostname", "hostname", resource_hostname)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}
		*/

		a, err := opts.AccountsDatabase.GetAccount(ctx, resource)

		if err != nil {
			logger.Error("Failed to retrieve account for resource", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		wf, err := a.WebfingerResource(ctx, opts.Hostname, opts.URIs)

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
