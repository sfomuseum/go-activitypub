package www

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/crypto"
)

type InboxPostHandlerOptions struct {
	AccountsDatabase  activitypub.AccountsDatabase
	FollowersDatabase activitypub.FollowersDatabase
	URIs              *activitypub.URIs
	Hostname          string
}

func InboxPostHandler(opts *InboxPostHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		if req.Method != http.MethodPost {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// START OF TBD...

		// sudo make me a regexp or req.PathId(...)

		account_id := filepath.Base(req.URL.Path)

		logger = logger.With("account", account_id)

		acct, err := opts.AccountsDatabase.GetAccount(ctx, account_id)

		if err != nil {
			logger.Error("Failed to retrieve inbox for account", "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		//

		var activity *ap.Activity

		dec := json.NewDecoder(req.Body)
		err = dec.Decode(&activity)

		if err != nil {
			logger.Error("Failed to decode message body", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		sender_id := activity.Actor
		logger = logger.With("sender_id", sender_id)

		slog.Info("INBOX", "sender", sender_id)

		follower_name, _, err := activitypub.ParseAccountURI(sender_id)

		if err != nil {
			logger.Error("Failed to parse send ID", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if follower_name == acct.Id {
			logger.Error("Can not follow yourself")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		switch activity.Type {
		case "Follow":

			is_following, err := opts.FollowersDatabase.IsFollowing(ctx, sender_id, acct.Id)

			if err != nil {
				logger.Error("Failed to determine if following", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			if is_following {
				logger.Error("Already following")
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

		case "Undo":

			is_following, err := opts.FollowersDatabase.IsFollowing(ctx, sender_id, acct.Id)

			if err != nil {
				logger.Error("Failed to determine if following", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			if !is_following {
				logger.Error("Not following")
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

		case "Create":

			is_following, err := opts.FollowersDatabase.IsFollowing(ctx, sender_id, acct.Id)

			if err != nil {
				logger.Error("Failed to determine if following", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			if !is_following {
				logger.Error("Not following")
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

		default:
			logger.Error("Unsupported activity type", "type", activity.Type)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("activity-type", activity.Type)

		// START OF verify request

		verifier, err := httpsig.NewVerifier(req)

		if err != nil {
			logger.Error("Failed to create signature verifier", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		key_id := verifier.KeyId()
		logger = logger.With("key id", key_id)

		logger.Info("Fetch other", "key_id", key_id)

		other_rsp, err := http.Get(key_id)

		if err != nil {
			logger.Error("Failed to retrieve key ID", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer other_rsp.Body.Close()

		var other_actor *ap.Actor

		dec = json.NewDecoder(other_rsp.Body)
		err = dec.Decode(&other_actor)

		if err != nil {
			logger.Error("Failed to other actor", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		public_key_str := other_actor.PublicKey.PEM

		if public_key_str == "" {
			logger.Error("Other actor missing public key")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		public_key, err := crypto.RSAPublicKeyFromPEM(public_key_str)

		if err != nil {
			logger.Error("Failed to parse PEM block containing public key", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		algo := httpsig.RSA_SHA512

		err = verifier.Verify(public_key, algo)

		if err != nil {
			logger.Error("Failed to verify signature", "error", err)
			http.Error(rsp, "Forbidden", http.StatusForbidden)
			return
		}

		// END OF put me in a function

		// Actually do something

		switch activity.Type {
		case "Follow":

			err = opts.FollowersDatabase.AddFollower(ctx, acct.Id, sender_id)

			if err != nil {
				logger.Error("Failed to add follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Undo":

			err = opts.FollowersDatabase.RemoveFollower(ctx, acct.Id, sender_id)

			if err != nil {
				logger.Error("Failed to remove follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Create":

			slog.Info("CREATE", "object", fmt.Sprintf("%T", activity.Object))

		default:
			// pass
		}

		accept, err := ap.NewAcceptActivity(ctx, acct.Id, sender_id)

		if err != nil {
			logger.Error("Failed to create new accept activity", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger = logger.With("accept", accept.Id)

		rsp.Header().Set("Content-type", "application/activity+json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(accept)

		if err != nil {
			logger.Error("Failed to encode accept activity", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
