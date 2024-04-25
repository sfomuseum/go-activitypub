package www

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/uris"
)

type AccountHandlerOptions struct {
	AccountsDatabase   activitypub.AccountsDatabase
	AliasesDatabase    activitypub.AliasesDatabase
	PropertiesDatabase activitypub.PropertiesDatabase
	URIs               *uris.URIs
	Templates          *template.Template
	RedirectOnAlias    bool
}

type AccountTemplateVars struct {
	Account    *activitypub.Account
	AccountURL string
	// To do: URLs (properties)
	IconURL string
}

func AccountHandler(opts *AccountHandlerOptions) (http.Handler, error) {

	account_t := opts.Templates.Lookup("account")

	if account_t == nil {
		return nil, fmt.Errorf("Failed to retrieve 'account' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		t1 := time.Now()

		defer func() {
			logger.Info("Time to serve request", "ms", time.Since(t1).Milliseconds())
		}()

		switch req.Method {
		case http.MethodGet, http.MethodPost:
			// pass
		default:
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
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

		if err != nil && err != activitypub.ErrNotFound {

			logger.Error("Failed to retrieve account", "error", err)

			if err == activitypub.ErrNotFound {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// START OF lookup account by alias

		if err != nil {

			logger.Debug("Lookup resource by alias")

			alias, err := opts.AliasesDatabase.GetAliasWithName(ctx, account_name)

			if err != nil {

				logger.Error("Failed to retrieve account for resource alias", "alias", account_name, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			acct, err = opts.AccountsDatabase.GetAccountWithId(ctx, alias.AccountId)

			if err != nil {

				logger.Error("Failed to retrieve account with ID for resource alias", "alias", account_name, "account id", alias.AccountId, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			// False by default - currently lacking hooks in the rest of the code
			// for setting this manually (cli flag, etc)

			if opts.RedirectOnAlias {

				acct_u := acct.AccountURL(ctx, opts.URIs)
				logger.Debug("Redirect to account page", "page", acct_u.String())

				http.Redirect(rsp, req, acct_u.String(), http.StatusSeeOther)
				return
			}
		}

		logger = logger.With("account id", acct.Id)

		// To do...
		// activitypub.PropertiesMapForAccount(ctx, acct, opts.PropertiesMap)

		// Check content-type here and HTML or JSON it up...

		if IsActivityStreamRequest(req, "Accept") {

			profile, err := acct.ProfileResource(ctx, opts.URIs)

			if err != nil {
				logger.Error("Failed to derive profile response for resource", "error", err)
				http.Error(rsp, "Not acceptable", http.StatusNotAcceptable)
				return
			}

			// To do: append properties map (above)

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

		icon_path := uris.AssignResource(opts.URIs.Icon, acct.Name)
		icon_url := uris.NewURL(opts.URIs, icon_path)

		// To do: URLs (properties map (above))

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
