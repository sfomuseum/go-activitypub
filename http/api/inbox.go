package api

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/go-fed/httpsig"
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

		verifier, err := httpsig.NewVerifier(req)

		if err != nil {
			slog.Error("Failed to create signature verifier", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		key_id := verifier.KeyId()
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

		public_key_str := other_actor.PublicKey.PEM
		slog.Info("OTHER", "key", public_key_str)

		if public_key_str == "" {
			slog.Error("Other actor missing public key")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF put me in a function

		public_key_block, _ := pem.Decode([]byte(public_key_str))

		if public_key_block == nil || public_key_block.Type != "PUBLIC KEY" {
			slog.Error("failed to decode PEM block containing public key")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		public_key, err := x509.ParsePKIXPublicKey(public_key_block.Bytes)

		if err != nil {
			slog.Error("Failed to parse PEM block containing public key", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// END OF put me in a function

		algo := httpsig.RSA_SHA512

		err = verifier.Verify(public_key, algo)

		if err != nil {
			slog.Error("Failed to verify signature", "error", err)
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
