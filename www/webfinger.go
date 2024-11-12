package www

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaronland/go-http-sanitize"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/go-activitypub/webfinger"
)

type WebfingerHandlerOptions struct {
	AccountsDatabase database.AccountsDatabase
	AliasesDatabase  database.AliasesDatabase
	URIs             *uris.URIs
}

func WebfingerHandler(opts *WebfingerHandlerOptions) (http.Handler, error) {

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

		name, host, err := ap.ParseAddress(resource)

		if err != nil {
			logger.Error("Failed to parse address (resource)", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account name", name)

		if host != "" && host != opts.URIs.Hostname {
			logger.Error("Resouce has bunk hostname", "host", host)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, name)

		if err != nil && err != activitypub.ErrNotFound {

			logger.Error("Failed to retrieve account for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// START OF lookup account by alias

		if err != nil {

			logger.Debug("Lookup resource by alias")

			alias, err := opts.AliasesDatabase.GetAliasWithName(ctx, name)

			if err != nil {

				logger.Error("Failed to retrieve account for resource alias", "alias", name, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			acct, err = opts.AccountsDatabase.GetAccountWithId(ctx, alias.AccountId)

			if err != nil {

				logger.Error("Failed to retrieve account with ID for resource alias", "alias", name, "account id", alias.AccountId, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			wf_u := acct.WebfingerURL(ctx, opts.URIs)
			logger.Debug("Redirect to webfinger endpoint", "endpoint", wf_u.String())

			http.Redirect(rsp, req, wf_u.String(), http.StatusSeeOther)
			return
		}

		// END OF lookup account by alias

		logger = logger.With("account id", acct.Id)

		wf, err := acct.WebfingerResource(ctx, opts.URIs)

		if err != nil {
			logger.Error("Failed to derive webfinger response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", webfinger.ContentType)

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
