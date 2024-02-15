package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/99designs/httpsignatures-go"
	"github.com/google/uuid"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
)

// https://paul.kinlan.me/adding-activity-pub-to-your-static-site/
// https://socialhub.activitypub.rocks/t/understanding-the-activity-pub-follow-request-flow/3033

// https://github.com/go-fed/httpsig

type InboxHandlerOptions struct {
	ActorDatabase activitypub.ActorDatabase
	URIs          *activitypub.URIs
	Hostname      string
}

func InboxHandler(opts *InboxHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("path", req.URL.Path)
		logger = logger.With("remote_addr", req.RemoteAddr)

		// START OF TBD...

		inbox_name := filepath.Base(req.URL.Path)
		resource := fmt.Sprintf("%s@%s", inbox_name, opts.Hostname)

		// END OF TBD...

		logger = logger.With("resource", resource)

		a, err := opts.ActorDatabase.GetActor(ctx, resource)

		if err != nil {
			slog.Error("Failed to retrieve inbox for resource", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		logger.Info("ACTOR", "a", a)

		// END OF verify request

		var activity *ap.Activity

		dec := json.NewDecoder(req.Body)
		err = dec.Decode(&activity)

		if err != nil {
			slog.Error("Failed to decode message body", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		//

		switch activity.Type {
		case "Follow", "Undo":
		default:
			slog.Error("Unsupported activity type", "type", activity.Type)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF verify request

		sig, err := httpsignatures.FromRequest(req)

		if err != nil {
			slog.Error("Failed to derive signature from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		key_id := sig.KeyID
		logger = logger.With("key id", key_id)

		slog.Info("Fetch other", "key_id", key_id)

		other_rsp, err := http.Get(key_id)

		if err != nil {
			slog.Error("Failed to retrieve key ID", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer other_rsp.Body.Close()

		var other_actor *ap.Actor

		dec = json.NewDecoder(other_rsp.Body)
		err = dec.Decode(&other_actor)

		if err != nil {
			slog.Error("Failed to other actor", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		key := other_actor.PublicKey.PEM
		slog.Info("OTHER", "key", key)

		if key == "" {
			slog.Error("Other actor missing public key")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if !sig.IsValid(key, req) {
			slog.Error("Request failed verification", "error", err)
			http.Error(rsp, "Forbidden", http.StatusForbidden)
			return
		}

		// Actually do something

		guid := uuid.New()
		logger.Info(guid.String())

		slog.Info("OKAY ACCEPT", "uuid", guid)
	}

	return http.HandlerFunc(fn), nil
}
