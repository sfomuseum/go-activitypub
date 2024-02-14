package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/signature"
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

		// START OF verify request
		
		// https://github.com/cbodonnell/go-pub/blob/07b0ca374c28f729bee2aaaae55461d54a09fc4c/handlers.go
		// https://github.com/go-fed/httpsig

		// https://blog.joinmastodon.org/2018/07/how-to-make-friends-and-verify-requests/

		sig, err := signature.ParseFromRequest(req)

		if err != nil {
			slog.Error("Failed to parse signature", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("key id", sig.KeyId)

		other_rsp, err := http.Get(sig.KeyId)

		if err != nil {
			slog.Error("Failed to retrieve key ID", "keyId", sig.KeyId)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer other_rsp.Body.Close()

		var other_actor *ap.Actor

		dec := json.NewDecoder(other_rsp.Body)
		err = dec.Decode(&other_actor)

		if err != nil {
			slog.Error("Failed to other actor", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if other_actor.PublicKey.PEM == "" {
			slog.Error("Other actor missing public key")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// END OF verify request
		
		var activity *ap.Activity

		dec = json.NewDecoder(req.Body)
		err = dec.Decode(&activity)

		if err != nil {
			slog.Error("Failed to decode message body", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		//

		//

		guid := uuid.New()

		logger.Info(guid.String())

		switch activity.Type {
		case "Follow":
		case "Undo":
		default:
			slog.Error("Unsupported activity type", "type", activity.Type)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		/*

			https://paul.kinlan.me/adding-activity-pub-to-your-static-site/

			    Parse the POST body and cast it to an Activity object.
			    Parse the signature of the request to verify the message hasn't been tampered with in transit.
			    From the signature HTTP header get the Actor that wants to follow you and fetch their Public Key (from their Actor file).
			    Verify the message with their Public Key

			Now we believe that we have a valid messages.

			If the message is a Follow request

			    See if the Actor trying to follow is already in the db, if they are return;
			    Add the Actor to the followers collection in FireStore
			    Prepare an Accept message to the Actor indicating that the Follow has been accepted and send it.

			If the message is an Undo for a Follow request.

			    Find the data in the followers collection in FireStore
			    Delete it.

		*/

	}

	return http.HandlerFunc(fn), nil
}
