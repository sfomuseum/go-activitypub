package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/sfomuseum/go-activitypub"
)

type ProfileHandlerOptions struct {
	ActorDatabase activitypub.ActorDatabase
	URIs          *activitypub.URIs
	Hostname      string
}

func ProfileHandler(opts *ProfileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("path", req.URL.Path)
		logger = logger.With("remote_addr", req.RemoteAddr)

		resource, err := sanitize.GetString(req, "resource")

		if err != nil {
			slog.Error("Failed to derive ?resource= parameter", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if resource == "" {
			slog.Error("Empty ?resource= parameter")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("resource", resource)

		a, err := opts.ActorDatabase.GetActor(ctx, resource)

		if err != nil {
			slog.Error("Failed to retrieve actor for resource", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		wf, err := a.ProfileResource(opts.Hostname, opts.URIs)

		if err != nil {
			slog.Error("Failed to derive profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/activity+json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(wf)

		if err != nil {
			slog.Error("Failed to encode profile response for resource", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
