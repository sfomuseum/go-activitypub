package api

import (
	_ "encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/sfomuseum/go-activitypub"
)

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
