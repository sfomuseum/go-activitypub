package www

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/sfomuseum/go-activitypub"
)

type AccountHandlerOptions struct {
	AccountsDatabase activitypub.AccountsDatabase
	URIs             *activitypub.URIs
	Templates        *template.Template
}

type AccountTemplateVars struct {
	Account    *activitypub.Account
	AccountURL string
	IconURL    string
}

func AccountHandler(opts *AccountHandlerOptions) (http.Handler, error) {

	account_t := opts.Templates.Lookup("account")

	if account_t == nil {
		return nil, fmt.Errorf("Failed to retrieve 'account' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

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

			return
		}

		acct.PrivateKeyURI = "constant://?val="

		account_url := acct.AccountURL(ctx, opts.URIs)

		icon_path := activitypub.AssignResource(opts.URIs.Icon, acct.Name)
		icon_url := activitypub.NewURL(opts.URIs, icon_path)

		vars := AccountTemplateVars{
			Account:    acct,
			IconURL:    icon_url.String(),
			AccountURL: account_url.String(),
		}

		rsp.Header().Set("Content-type", "text/html")

		err = account_t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
