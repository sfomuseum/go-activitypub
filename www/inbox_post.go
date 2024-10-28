package www

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/blocks"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/followers"
	"github.com/sfomuseum/go-activitypub/following"
	"github.com/sfomuseum/go-activitypub/messages"
	"github.com/sfomuseum/go-activitypub/notes"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/queue"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/tidwall/gjson"
)

type InboxPostHandlerOptions struct {
	AccountsDatabase    database.AccountsDatabase
	FollowersDatabase   database.FollowersDatabase
	FollowingDatabase   database.FollowingDatabase
	MessagesDatabase    database.MessagesDatabase
	NotesDatabase       database.NotesDatabase
	PostsDatabase       database.PostsDatabase
	BlocksDatabase      database.BlocksDatabase
	LikesDatabase       database.LikesDatabase
	BoostsDatabase      database.BoostsDatabase
	ProcessMessageQueue queue.ProcessMessageQueue
	URIs                *uris.URIs
	AllowFollow         bool
	AllowCreate         bool
	AllowLikes          bool
	AllowBoosts         bool
	// Allows posts to accounts not followed by author but where account is mentioned in post
	AllowMentions bool

	// TBD but the idea is that after the signature verification
	// and block checks are dealt with the best thing would be to
	// hand off to an activity-specific handler using the http.Next()
	// trick rather than smushing all the code in to this handler.
	// It is still unclear which variables need to be pass down to
	// the final activity handler or how (context.Context probably?)
	// Activities map[string]http.Handler
}

func InboxPostHandler(opts *InboxPostHandlerOptions) (http.Handler, error) {

	http_cl := &http.Client{}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		wg := new(sync.WaitGroup)

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
			logger.Error("Not activitystream request", "content-type", req.Header.Get("Content-Type"))
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// Start by doing some initial sanity checking on the message body

		var activity *ap.Activity
		var activity_r io.Reader

		// START OF necessary to track down message parsing errors...

		limited_r := activitypub.DefaultLimitedReader(req.Body)

		// Make me a flag...
		// Or maybe not...
		log_body := false

		if log_body {

			body, err := io.ReadAll(limited_r)

			if err != nil {
				logger.Error("Failed to read message body", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			logger.Debug("DEBUG", "body", string(body))
			activity_r = bytes.NewReader(body)

		} else {
			activity_r = limited_r
		}

		// END OF necessary to track down message parsing errors...

		dec := json.NewDecoder(activity_r)
		err := dec.Decode(&activity)

		if err != nil {
			logger.Error("Failed to decode message body", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		requestor_address := activity.Actor

		logger = logger.With("requestor_address", requestor_address)
		logger = logger.With("activity_type", activity.Type)

		// There is a not insignificant number of people crawling ActivityPub
		// endpoints issuing "Delete" activities just to see if they will stick...

		switch activity.Type {
		case "Accept":

			// To do: Actually implement this with checks/validation...

			logger.Debug("Received 'Accept' activity", "response code", http.StatusAccepted)
			rsp.WriteHeader(http.StatusAccepted)
			return

		case "Like":

			if !opts.AllowLikes {
				logger.Error("Unsupported activity type, likes are disabled")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		case "Announce":

			if !opts.AllowBoosts {
				logger.Error("Unsupported activity type, boosts are disabled")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		case "Follow":

			if !opts.AllowFollow {
				logger.Error("Unsupported activity type")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		case "Undo":

			if !opts.AllowFollow && !opts.AllowLikes && !opts.AllowBoosts {
				logger.Error("Unsupported activity type")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

			enc_obj, err := json.Marshal(activity.Object)

			if err != nil {
				logger.Error("Failed to marshal activity object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			// Block activities are not supported (yet)

			type_rsp := gjson.GetBytes(enc_obj, "type")

			switch type_rsp.String() {
			case "Follow":

				if !opts.AllowFollow {
					logger.Error("Unsupported activity type, follows are disabled")
					http.Error(rsp, "Not implemented", http.StatusNotImplemented)
					return
				}

			case "Like":

				if !opts.AllowLikes {
					logger.Error("Unsupported activity type, likes are disabled")
					http.Error(rsp, "Not implemented", http.StatusNotImplemented)
					return
				}

			case "Announce":

				if !opts.AllowBoosts {
					logger.Error("Unsupported activity type, boosts are disabled")
					http.Error(rsp, "Not implemented", http.StatusNotImplemented)
					return
				}

			default:
				logger.Error("Unsupported undo activity type", "type", type_rsp.String())
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		case "Create":

			if !opts.AllowCreate {
				logger.Error("Unsupported activity type")
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

			// Ensure we are creating a "Note"

			enc_obj, err := json.Marshal(activity.Object)

			if err != nil {
				logger.Error("Failed to marshal activity object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			type_rsp := gjson.GetBytes(enc_obj, "type")

			switch type_rsp.String() {
			case "Note":
				// Okay
			default:
				logger.Error("Unsupported undo activity type", "type", type_rsp.String())
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		default:
			logger.Debug("Unsupported activity type")
			http.Error(rsp, "Not implemented", http.StatusNotImplemented)
			return
		}

		// Ensure the account being poked exists

		account_name, host, err := ap.ParseAddressFromRequest(req)

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

		// Figure out who is doing the poking

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

			profile_r := activitypub.DefaultLimitedReader(profile_rsp.Body)

			dec = json.NewDecoder(profile_r)
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

			requestor_name, requestor_host, err = ap.ParseAddress(requestor_address)

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

		// Check if the requestor is being blocked

		is_blocked, err := blocks.IsBlockedByAccount(ctx, opts.BlocksDatabase, acct.Id, requestor_host, requestor_name)

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

		// One final sanity check

		switch activity.Type {
		case "Follow", "Undo", "Like", "Announce":

			// Note: We have prevented Block Undo activities above

			if requestor_name == acct.Name {
				logger.Error("Account name mismatch for activity type", "type", activity.Type, "requester", requestor_name, "account", acct.Name)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

		default:
			// pass
		}

		// START OF verify request

		// This is important if the server is running behind some kind of proxy (for example Lambda)
		// or the signature verification will fail

		logger.Debug("Manually set hostname on request", "hostname", opts.URIs.Hostname)
		req.Header.Set("Host", opts.URIs.Hostname)
		req.Host = opts.URIs.Hostname

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

			sender_r := activitypub.DefaultLimitedReader(sender_rsp.Body)

			var sender_actor *ap.Actor

			dec = json.NewDecoder(sender_r)
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

		// END OF verify request

		// Actually do something

		// So really, at this point we should simple have per actitivy type handlers that
		// get passed in by opts and we just hand off to them using the http.Next trick.
		// As written we're going to add another 300+ lines of code which really just makes
		// everything more confusing...
		//
		// But really it's more like another 800+ lines of code now that boosts, likes and
		// replies are handled...

		accept_obj := activity

		logger.Info("PROCESS", "type", activity.Type)
		
		switch activity.Type {
		case "Announce":

			var object_uri string

			switch activity.Object.(type) {
			case string:
				object_uri = activity.Object.(string)
			case map[string]interface{}:

				// This is here because the code in deliver.go expects to _dispatch_
				// "Announce" messages including the note of the post being boosted
				// as the body (object) of the activity. This may or may not be incorrect.
				// I am not sure.

				//logger.Info("DEBUG", "map", activity.Object)
				
				obj_map := activity.Object.(map[string]interface{})
				v, exists := obj_map["url"]

				if !exists {
					logger.Error("Object map for announce missing 'url' key")
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				switch v.(type) {
				case string:
					object_uri = v.(string)
				default:
					logger.Error("Invalid or unsupported type for announce url value", "value", v, "type", fmt.Sprintf("%T", v))
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				// logger.Info("DEBUG", "object_uri", object_uri)

			default:
				logger.Error("Invalid or unsupport activity object type for announce activity", "type", fmt.Sprintf("%T", activity.Object))
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			post, err := posts.GetPostFromObjectURI(ctx, opts.URIs, opts.PostsDatabase, object_uri)

			if err != nil {

				logger.Error("Failed to derive post from object URI", "object uri", object_uri, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			logger = logger.With("post id", post.Id)

			if post.AccountId != acct.Id {
				logger.Error("Trying to act on post for different account", "post account", post.AccountId)
				http.Error(rsp, "Forbidden", http.StatusForbidden)
				return
			}

			boost, err := opts.BoostsDatabase.GetBoostWithPostIdAndActor(ctx, post.Id, activity.Actor)

			if err != nil && err != activitypub.ErrNotFound {
				logger.Error("Failed to derive boost from post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			if boost == nil {

				boost, err = activitypub.NewBoost(ctx, post, activity.Actor)

				if err != nil {
					logger.Error("Failed to create new boost for post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				err = opts.BoostsDatabase.AddBoost(ctx, boost)

				if err != nil {
					logger.Error("Failed to add new boost for post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger.Info("Create new boost", "post id", post.Id, "actor", activity.Actor, "boost", boost.Id)
			}

			// Do we need to defer accept here? Apparently not...

		case "Like":

			var object_uri string

			switch activity.Object.(type) {
			case string:
				object_uri = activity.Object.(string)
			default:
				logger.Error("Invalid or unsupport activity object type", "type", fmt.Sprintf("%T", activity.Object))
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			post, err := posts.GetPostFromObjectURI(ctx, opts.URIs, opts.PostsDatabase, object_uri)

			if err != nil {

				logger.Error("Failed to derive post from object URI", "object uri", object_uri, "error", err)

				if err == activitypub.ErrNotFound {
					http.Error(rsp, "Not found", http.StatusNotFound)
					return
				}

				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			logger = logger.With("post id", post.Id)

			if post.AccountId != acct.Id {
				logger.Error("Trying to act on post for different account", "post account", post.AccountId)
				http.Error(rsp, "Forbidden", http.StatusForbidden)
				return
			}

			like, err := opts.LikesDatabase.GetLikeWithPostIdAndActor(ctx, post.Id, activity.Actor)

			if err != nil && err != activitypub.ErrNotFound {
				logger.Error("Failed to derive like from post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			if like == nil {

				like, err = activitypub.NewLike(ctx, post, activity.Actor)

				if err != nil {
					logger.Error("Failed to create new like for post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				err = opts.LikesDatabase.AddLike(ctx, like)

				if err != nil {
					logger.Error("Failed to add new like for post and actor", "post id", post.Id, "actor", activity.Actor, "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger.Info("Create new like", "post id", post.Id, "actor", activity.Actor, "like", like.Id)
			}

			// Do we need to defer accept here? Apparently not...

		case "Follow":

			is_following, _, err := followers.IsFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

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

			err = followers.AddFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

			if err != nil {
				logger.Error("Failed to create new follower", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			// It is unclear whether it is really necessary to send this request in a deferred
			// function (or whether it can be sent inline before the HTTP 202 response is sent
			// below. On the other there are accept activities which are specifically meant to
			// happen "out-of-band", like follower requests that are manually approved, so the
			// easiest way to think about things is that they will (maybe?) get moved in to its
			// own delivery queue (distinct from posts) to happen after the inbox handler has
			// completed. Basically: Treat every message sent to the ActivityPub inbox as an
			// offline task. I am still trying to determine if that's an accurate assumption.

			defer func() {

				accept_actor := acct.AccountURL(ctx, opts.URIs).String()
				logger.Debug("Send Accept activity in defer", "actor", accept_actor)

				accept, err := ap.NewAcceptActivity(ctx, opts.URIs, accept_actor, accept_obj)

				if err != nil {
					logger.Error("Failed to create new accept activity", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger = logger.With("accept", accept.Id)

				err = acct.SendActivity(ctx, opts.URIs, requestor_actor.Inbox, accept)

				if err != nil {

					logger.Error("Failed to post accept activity to requestor, remove follower", "to", requestor_actor.Inbox, "error", err)

					f, err := followers.GetFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

					if err != nil {
						logger.Error("Failed to retrieve newly created follower to remove", "error", err)
					} else {

						err = opts.FollowersDatabase.RemoveFollower(ctx, f)

						if err != nil {
							logger.Error("Failed to remove follower", "id", f.Id, "error", err)
						}
					}

					// Note: We are not returning a HTTP response since we have already returned HTTP 202 below
					return
				}

			}()

		case "Undo":

			enc_obj, err := json.Marshal(activity.Object)

			if err != nil {
				logger.Error("Failed to marshal activity object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			var object_activity *ap.Activity

			err = json.Unmarshal(enc_obj, &object_activity)

			if err != nil {
				logger.Error("Failed to derive activity from object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			logger = logger.With("object type", object_activity.Type)

			switch object_activity.Type {
			case "Follow":

				// Note: We have prevented Block Undo activities above so we're going to assume it's an undo follow request

				is_following, f, err := followers.IsFollower(ctx, opts.FollowersDatabase, acct.Id, requestor_address)

				if err != nil {
					logger.Error("Failed to determine if following", "error", err)
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				if is_following {

					err = opts.FollowersDatabase.RemoveFollower(ctx, f)

					if err != nil {
						logger.Error("Failed to remove follower", "error", err)
						http.Error(rsp, "Internal server error", http.StatusInternalServerError)
						return
					}
				}

			case "Like":

				logger = logger.With("actor", activity.Actor)

				var object_uri string

				switch object_activity.Object.(type) {
				case string:
					object_uri = object_activity.Object.(string)
				default:
					logger.Error("Invalid or unsupport activity object type", "type", fmt.Sprintf("%T", object_activity.Object))
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				logger = logger.With("object uri", object_uri)

				post, err := posts.GetPostFromObjectURI(ctx, opts.URIs, opts.PostsDatabase, object_uri)

				if err != nil {

					logger.Error("Failed to derive post from object URI", "error", err)

					if err == activitypub.ErrNotFound {
						http.Error(rsp, "Not found", http.StatusNotFound)
						return
					}

					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger = logger.With("post id", post.Id)

				if post.AccountId != acct.Id {
					logger.Error("Trying to act on post for different account", "post account", post.AccountId)
					http.Error(rsp, "Forbidden", http.StatusForbidden)
					return
				}

				like, err := opts.LikesDatabase.GetLikeWithPostIdAndActor(ctx, post.Id, activity.Actor)

				if err != nil && err != activitypub.ErrNotFound {
					logger.Error("Failed to derive like from post and actor", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				if like != nil {

					logger = logger.With("like", like.Id)

					err := opts.LikesDatabase.RemoveLike(ctx, like)

					if err != nil {
						logger.Error("Failed to remove like", "error", err)
						http.Error(rsp, "Internal server error", http.StatusInternalServerError)
						return
					}

					logger.Info("Removed like")
				}

			case "Announce":

				logger = logger.With("actor", activity.Actor)

				var object_uri string

				switch object_activity.Object.(type) {
				case string:
					object_uri = object_activity.Object.(string)
				default:
					logger.Error("Invalid or unsupport activity object type", "type", fmt.Sprintf("%T", object_activity.Object))
					http.Error(rsp, "Bad request", http.StatusBadRequest)
					return
				}

				logger = logger.With("object uri", object_uri)

				post, err := posts.GetPostFromObjectURI(ctx, opts.URIs, opts.PostsDatabase, object_uri)

				if err != nil {

					logger.Error("Failed to derive post from object URI", "error", err)

					if err == activitypub.ErrNotFound {
						http.Error(rsp, "Not found", http.StatusNotFound)
						return
					}

					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				logger = logger.With("post id", post.Id)

				if post.AccountId != acct.Id {
					logger.Error("Trying to act on post for different account", "post account", post.AccountId)
					http.Error(rsp, "Forbidden", http.StatusForbidden)
					return
				}

				boost, err := opts.BoostsDatabase.GetBoostWithPostIdAndActor(ctx, post.Id, activity.Actor)

				if err != nil && err != activitypub.ErrNotFound {
					logger.Error("Failed to derive boost from post and actor", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				if boost != nil {

					logger = logger.With("boost", boost.Id)

					err := opts.BoostsDatabase.RemoveBoost(ctx, boost)

					if err != nil {
						logger.Error("Failed to remove boost", "error", err)
						http.Error(rsp, "Internal server error", http.StatusInternalServerError)
						return
					}

					logger.Info("Removed boost")
				}

			default:
				logger.Error("Unsupported object type for undo", "type", object_activity.Type)
				http.Error(rsp, "Not implemented", http.StatusNotImplemented)
				return
			}

		case "Create":

			enc_obj, err := json.Marshal(activity.Object)

			if err != nil {
				logger.Error("Failed to marshal activity object", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			var note *ap.Note

			err = json.Unmarshal(enc_obj, &note)

			if err != nil {
				logger.Error("Failed to unmarshal activity note", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			logger = logger.With("note id", note.Id)

			is_following, _, err := following.IsFollowing(ctx, opts.FollowingDatabase, acct.Id, requestor_address)

			if err != nil {
				logger.Error("Failed to determine if following", "error", err)
				http.Error(rsp, "Bad request", http.StatusBadRequest)
				return
			}

			logger = logger.With("is following", is_following)

			// If we are following this account then it's all good

			is_allowed := is_following

			// If not following then check to see whether account (being posted to)
			// is mentioned in post (being received)

			if !is_allowed && opts.AllowMentions && len(note.Tags) > 0 {

				// https://github.com/sfomuseum/go-activitypub/issues/3
				// account_url := acct.URL

				// And yet it appears to actually be {ACTOR}.id however this
				// does not work (where "work" means open profile tab) in Ivory
				// yet because... I have no idea
				account_url := acct.AccountURL(ctx, opts.URIs).String()

				for _, t := range note.Tags {

					if t.Href == account_url {
						logger.Info("Post author is not followed but account is mentioned")
						is_allowed = true
						break
					}
				}

				if !is_allowed {
					logger.Warn("Post author is not followed and not found in mentions", "account_url", account_url)
				}
			}

			// If still not allowed (not following, not mentioned) then check to see if the post
			// (being received) is in reply to something that the account (being posted to) wrote

			if !is_allowed {

				if note.InReplyTo == "" {

					if !is_following {
						logger.Error("Not following")
						http.Error(rsp, "Forbidden", http.StatusForbidden)
						return
					}

				} else {

					logger = logger.With("in reply to", note.InReplyTo)

					// If we are not already following the author then

					// Fetch the (in-reply-to) post in question and check to see if it
					// was authored by acct

					is_own_post := false

					post, err := posts.GetPostFromObjectURI(ctx, opts.URIs, opts.PostsDatabase, note.InReplyTo)

					if err != nil && err != activitypub.ErrNotFound {
						logger.Error("Failed to determine if object URI references post", "error", err)
						http.Error(rsp, "Internal server error", http.StatusInternalServerError)
						return
					}

					if post != nil {

						logger = logger.With("post id", post.Id)

						if post.AccountId == acct.Id {
							is_own_post = true
						}
					}

					if !is_own_post {
						logger.Error("Reply-to is not own post")
						http.Error(rsp, "Forbidden", http.StatusForbidden)
						return
					}
				}
			}

			// First store the activity pub note as a local "note" - that is store
			// the message from person (x) exactly once regardless of how many
			// different accounts (on this service) that the note is being delivered
			// to

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

			if db_note != nil {

				logger = logger.With("note id", db_note.Id)

			} else {

				new_note, err := notes.AddNote(ctx, opts.NotesDatabase, note_uuid, requestor_address, string(enc_obj))

				if err != nil {
					logger.Error("Failed to create new note", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_note = new_note
				logger = logger.With("note id", db_note.Id)
			}

			// Now store a "message" which is a pointer to the note associated with the account the
			// note is being delivered to

			db_message, err := messages.GetMessage(ctx, opts.MessagesDatabase, acct.Id, db_note.Id)

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

				db_message, err = messages.UpdateMessage(ctx, opts.MessagesDatabase, db_message)

				if err != nil {
					logger.Error("Failed to update message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

			} else {

				new_message, err := messages.AddMessage(ctx, opts.MessagesDatabase, acct.Id, db_note.Id, requestor_address)

				if err != nil {
					logger.Error("Failed to add message", "error", err)
					http.Error(rsp, "Internal server error", http.StatusInternalServerError)
					return
				}

				db_message = new_message
				logger = logger.With("message id", db_message.Id)
			}

			// Dispatch to handle message queue here... maybe?

			logger.Info("Note has been added to messages")

			wg.Add(1)

			go func() {

				defer wg.Done()

				err = opts.ProcessMessageQueue.ProcessMessage(ctx, db_message.Id)

				if err != nil {
					logger.Error("Failed to process message with process queue", "error", err)
					return
				}

				logger.Info("Delivered message to processing queue")
			}()

		default:
			// pass
		}

		logger.Debug("Inbox post complete", "status", http.StatusAccepted)
		rsp.WriteHeader(http.StatusAccepted)

		wg.Wait()
		return
	}

	return http.HandlerFunc(fn), nil
}
