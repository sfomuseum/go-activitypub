package www

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	_ "log/slog"
	"net/http"
	"path/filepath"
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
	Hostname          string
	AllowFollow       bool
	AllowCreate       bool
}

func InboxPostHandler(opts *InboxPostHandlerOptions) (http.Handler, error) {

	http_cl := &http.Client{}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := LoggerWithRequest(req, nil)

		if req.Method != http.MethodPost {
			logger.Error("Method not allowed")
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !IsActivityStreamRequest(req) {
			logger.Error("Not activitystream request")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// START OF TBD...

		// sudo make me a regexp or req.PathId(...)

		account_name := filepath.Base(req.URL.Path)

		logger = logger.With("account", account_name)

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

		sender_address := activity.Actor
		logger = logger.With("sender_address", sender_address)

		sender_name, sender_host, err := activitypub.ParseAccountURI(sender_address)

		if err != nil {
			logger.Error("Failed to parse send ID", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		is_blocked, err := opts.BlocksDatabase.IsBlockedByAccount(ctx, acct.Id, sender_host, sender_name)

		if err != nil {
			logger.Error("Failed to determine if sender is blocked", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if is_blocked {
			logger.Error("Sender is blocked")
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

			if sender_name == acct.Name {
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

		verifier, err := httpsig.NewVerifier(req)

		if err != nil {
			logger.Error("Failed to create signature verifier", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		key_id := verifier.KeyId()
		logger = logger.With("key id", key_id)

		logger.Info("Fetch key for sender", "key_id", key_id)

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

		public_key_str := sender_actor.PublicKey.PEM

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

		// Actually do something

		switch activity.Type {
		case "Follow":

			is_following, err := opts.FollowersDatabase.IsFollowing(ctx, sender_address, acct.Id)

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

			err = opts.FollowersDatabase.AddFollower(ctx, acct.Id, sender_address)

			if err != nil {
				logger.Error("Failed to add follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Undo":

			is_following, err := opts.FollowersDatabase.IsFollowing(ctx, sender_address, acct.Id)

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

			err = opts.FollowersDatabase.RemoveFollower(ctx, acct.Id, sender_address)

			if err != nil {
				logger.Error("Failed to remove follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

		case "Create":

			is_following, err := opts.FollowingDatabase.IsFollowing(ctx, acct.Id, sender_address)

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

			note_id := note.Id
			logger = logger.With("note id", note_id)

			db_note, err := opts.NotesDatabase.GetNoteWithNoteIdAndAuthorAddress(ctx, note_id, sender_address)

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

				db_id, err := activitypub.NewId()

				if err != nil {
					logger.Error("Failed to create new ID for note", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_note = &activitypub.Note{
					Id:            db_id,
					NoteId:        note_id,
					AuthorAddress: sender_address,
					Body:          enc_note,
					Created:       ts,
					LastModified:  ts,
				}

				err = opts.NotesDatabase.AddNote(ctx, db_note)

				if err != nil {
					logger.Error("Failed to add note", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger = logger.With("note id", db_note.Id)
			}

			db_message, err := opts.MessagesDatabase.GetMessageWithAccountAndNoteIds(ctx, acct.Id, db_note.Id)

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

				db_message.LastModified = ts

				err = opts.MessagesDatabase.UpdateMessage(ctx, db_message)

				if err != nil {
					logger.Error("Failed to update message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

			} else {

				db_id, err := activitypub.NewId()

				if err != nil {
					logger.Error("Failed to create new ID for message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_message = &activitypub.Message{
					Id:            db_id,
					NoteId:        db_note.Id,
					AuthorAddress: sender_address,
					AccountId:     acct.Id,
					Created:       ts,
					LastModified:  ts,
				}

				err = opts.MessagesDatabase.AddMessage(ctx, db_message)

				if err != nil {
					logger.Error("Failed to add message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger = logger.With("message id", db_message.Id)
			}

			logger.Info("Note has been added to messages")

		default:
			// pass
		}

		// return acceptance

		acct_address := acct.Address(opts.Hostname)

		accept, err := ap.NewAcceptActivity(ctx, acct_address, sender_address)

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
