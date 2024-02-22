package www

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/crypto"
)

type InboxPostHandlerOptions struct {
	AccountsDatabase  activitypub.AccountsDatabase
	FollowersDatabase activitypub.FollowersDatabase
	FollowingDatabase activitypub.FollowingDatabase
	MessagesDatabase  activitypub.MessagesDatabase
	NotesDatabase     activitypub.NotesDatabase
	BlocksDatabase    activitypub.BlocksDatabase
	URIs              *activitypub.URIs
	AllowFollow       bool
	AllowCreate       bool
}

func InboxPostHandler(opts *InboxPostHandlerOptions) (http.Handler, error) {

	http_cl := &http.Client{}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		t1 := time.Now()

		defer func() {
			logger.Info("Time to serve request", "ms", time.Since(t1).Milliseconds())
		}()

		if req.Method != http.MethodPost {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !IsActivityStreamRequest(req, "Content-Type") {
			logger.Error("Not activitystream request")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		account_name, host, err := activitypub.ParseAddressFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse address from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		logger = logger.With("account", account_name)

		if host != "" && host != opts.URIs.Hostname {
			logger.Error("Resouce has bunk hostname", "host", host)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

		if err != nil {

			logger.Error("Failed to retrieve inbox for account", "error", err)

			if err == activitypub.ErrNotFound {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger = logger.With("account id", acct.Id)

		//

		var activity *ap.Activity

		dec := json.NewDecoder(req.Body)
		err = dec.Decode(&activity)

		if err != nil {
			logger.Error("Failed to decode message body", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// enc_activity, _ := json.Marshal(activity)
		// logger = logger.With("activity", string(enc_activity))

		requestor_address := activity.Actor
		logger = logger.With("requestor_address", requestor_address)

		var requestor_name string
		var requestor_host string
		var requestor_actor *ap.Actor

		if strings.HasPrefix(requestor_address, "http") {

			logger.Debug("Derive requestor address from URL")

			requestor_u, err := url.Parse(requestor_address)

			if err != nil {
				logger.Error("Failed to parse address URL for requestor", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			profile_req, err := http.NewRequestWithContext(ctx, "GET", requestor_address, nil)

			if err != nil {
				logger.Error("Failed to create actor/profile request for requestor", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			profile_req.Header.Set("Accept", ap.ACTIVITYSTREAMS_ACCEPT_HEADER)

			profile_rsp, err := http_cl.Do(profile_req)

			if err != nil {
				logger.Error("Failed to retrieve actor/profile request for requestor", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			defer profile_rsp.Body.Close()

			if profile_rsp.StatusCode != http.StatusOK {
				logger.Error("Remote endpoint did not return successfully for actor/profile request for requestor", "code", profile_rsp.StatusCode, "status", profile_rsp.Status)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			var actor *ap.Actor

			dec = json.NewDecoder(profile_rsp.Body)
			err = dec.Decode(&actor)

			if err != nil {
				logger.Error("Failed to decode actor/profile response for requestor", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			requestor_actor = actor
			requestor_name = actor.PreferredUsername
			requestor_host = requestor_u.Host

			requestor_address = fmt.Sprintf("%s@%s", requestor_name, requestor_host)
			logger = logger.With("requestor_address", requestor_address)
			logger.Debug("Re-assign requestor_address variable")

		} else {

			requestor_name, requestor_host, err = activitypub.ParseAddress(requestor_address)

			if err != nil {
				logger.Error("Failed to parse requestor address", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			if requestor_name == "" || requestor_host == "" {
				logger.Error("Requestor address missing name or host", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}
		}

		logger = logger.With("requestor_address", requestor_address, "requestor_name", requestor_name, "requestor_host", requestor_host)

		is_blocked, err := activitypub.IsBlockedByAccount(ctx, opts.BlocksDatabase, acct.Id, requestor_host, requestor_name)

		if err != nil {
			logger.Error("Failed to determine if requestor is blocked", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if is_blocked {
			logger.Error("Requestor is blocked")
			http.Error(rsp, "Forbidden", http.StatusForbidden)
			return
		}

		logger = logger.With("activity-type", activity.Type)

		switch activity.Type {
		case "Follow", "Undo":

			if !opts.AllowFollow {
				logger.Error("Unsupported activity type")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

			if requestor_name == acct.Name {
				logger.Error("Can not follow yourself")
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

		case "Create":

			if !opts.AllowCreate {
				logger.Error("Unsupported activity type")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		default:
			logger.Error("Unsupported activity type")
			http.Error(rsp, "Not implemented", http.StatusNotImplemented)
			return
		}

		// START OF verify request

		// This is important if the server is running behind some kind of proxy (for example Lambda)
		// or the signature verification will fail

		if opts.URIs.Hostname != "" {
			logger.Debug("Manually set hostname on request", "hostname", opts.URIs.Hostname)
			req.Header.Set("Host", opts.URIs.Hostname)
			req.Host = opts.URIs.Hostname
		}

		verifier, err := httpsig.NewVerifier(req)

		if err != nil {
			logger.Error("Failed to create signature verifier", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		key_id := verifier.KeyId()
		logger = logger.With("key id", key_id)

		var requestor_public_key ap.PublicKey

		if requestor_actor != nil && requestor_actor.PublicKey.Id == key_id {
			logger.Debug("request public key ID is the same as signature key ID")
			requestor_public_key = requestor_actor.PublicKey
		} else {

			logger.Info("Fetch key for requestor", "key_id", key_id)

			// Please stop calling this sender_
			// Is it safe-not-confusing to call these varaibles requestor_ or should they have a diffrerent prefix...

			sender_req, err := http.NewRequestWithContext(ctx, "GET", key_id, nil)

			if err != nil {
				logger.Error("Failed to create request for key id", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			sender_req.Header.Set("Accept", ap.ACTIVITYSTREAMS_ACCEPT_HEADER)

			sender_rsp, err := http_cl.Do(sender_req)

			if err != nil {
				logger.Error("Failed to retrieve key ID", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			defer sender_rsp.Body.Close()

			var sender_actor *ap.Actor

			dec = json.NewDecoder(sender_rsp.Body)
			err = dec.Decode(&sender_actor)

			if err != nil {
				logger.Error("Failed to other actor", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			requestor_actor = sender_actor
			requestor_public_key = sender_actor.PublicKey
		}

		public_key_str := requestor_public_key.PEM

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

		// Read from signature...
		algo := httpsig.RSA_SHA256

		err = verifier.Verify(public_key, algo)

		if err != nil {
			logger.Error("Failed to verify signature", "error", err, "signature", req.Header.Get("Signature"), "digest", req.Header.Get("Digest"))
			http.Error(rsp, "Forbidden", http.StatusForbidden)
			return
		}

		// Actually do something

		switch activity.Type {
		case "Follow":

			is_following, _, err := activitypub.IsFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

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

			err = activitypub.AddFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

			if err != nil {
				logger.Error("Failed to create new follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Undo":

			is_following, f, err := activitypub.IsFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

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

			err = opts.FollowersDatabase.RemoveFollower(ctx, f)

			if err != nil {
				logger.Error("Failed to remove follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Create":

			is_following, _, err := activitypub.IsFollowing(ctx, opts.FollowingDatabase, acct.Id, requestor_address)

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

			enc_note, err := json.Marshal(activity.Object)

			if err != nil {
				logger.Error("Failed to marshal activity object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			var note *ap.Note

			err = json.Unmarshal(enc_note, &note)

			if err != nil {
				logger.Error("Failed to unmarshal activity note", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			note_uuid := note.Id
			logger = logger.With("note uuid", note_uuid)

			db_note, err := opts.NotesDatabase.GetNoteWithUUIDAndAuthorAddress(ctx, note_uuid, requestor_address)

			switch {
			case err == activitypub.ErrNotFound:
				// pass
			case err != nil:
				logger.Error("Failed to retrive note", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			default:
				// pass
			}

			now := time.Now()
			ts := now.Unix()

			if db_note != nil {

				logger = logger.With("note id", db_note.Id)

				if bytes.Equal(enc_note, db_note.Body) {
					logger.Error("Note already registered")
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				db_note.Body = enc_note
				db_note.LastModified = ts

				err = opts.NotesDatabase.UpdateNote(ctx, db_note)

				if err != nil {
					logger.Error("Failed to update note", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

			} else {

				new_note, err := activitypub.AddNote(ctx, opts.NotesDatabase, note_uuid, requestor_address, enc_note)

				if err != nil {
					logger.Error("Failed to create new note", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_note = new_note
				logger = logger.With("note id", db_note.Id)
			}

			db_message, err := activitypub.GetMessage(ctx, opts.MessagesDatabase, acct.Id, db_note.Id)

			switch {
			case err == activitypub.ErrNotFound:
				// pass
			case err != nil:
				logger.Error("Failed to retrive message", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			default:
				// pass
			}

			if db_message != nil {

				logger = logger.With("message id", db_message.Id)

				db_message, err = activitypub.UpdateMessage(ctx, opts.MessagesDatabase, db_message)

				if err != nil {
					logger.Error("Failed to update message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

			} else {

				new_message, err := activitypub.AddMessage(ctx, opts.MessagesDatabase, acct.Id, db_note.Id, requestor_address)

				if err != nil {
					logger.Error("Failed to add message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_message = new_message
				logger = logger.With("message id", db_message.Id)
			}

			logger.Info("Note has been added to messages")

		default:
			// pass
		}

		// return acceptance

		acct_address := acct.Address(opts.URIs.Hostname)

		accept, err := ap.NewAcceptActivity(ctx, acct_address, requestor_address)

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
