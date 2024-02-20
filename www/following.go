package www

import (
	"net/http"

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

		rsp.Header().Set("Content-type", ap.ACTIVITY_CONTENT_TYPE)

		return
	}

	return http.HandlerFunc(fn), nil
}
