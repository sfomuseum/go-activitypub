package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/sfomuseum/go-activitypub"
)

type ProfileHandlerOptions struct {
	AccountDatabase activitypub.AccountDatabase
	URIs            *activitypub.URIs
	Hostname        string
}

func ProfileHandler(opts *ProfileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("path", req.URL.Path)
		logger = logger.With("remote_addr", req.RemoteAddr)

		// START OF TBD...

		account_name := filepath.Base(req.URL.Path)
		resource := fmt.Sprintf("%s@%s", account_name, opts.Hostname)

		// END OF TBD...

		logger = logger.With("resource", resource)

		a, err := opts.AccountDatabase.GetAccount(ctx, resource)

		if err != nil {
			slog.Error("Failed to retrieve account for resource", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		profile, err := a.ProfileResource(ctx, opts.Hostname, opts.URIs)

		if err != nil {
			slog.Error("Failed to derive profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/activity+json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(profile)

		if err != nil {
			slog.Error("Failed to encode profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
