package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverPostOptions struct {
	From               *Account           `json:"from"`
	To                 string             `json:"to"`
	Post               *Post              `json:"post"`
	PostTags           []*PostTag         `json:"post_tags"`
	URIs               *uris.URIs         `json:"uris"`
	DeliveriesDatabase DeliveriesDatabase `json:"deliveries_database,omitempty"`
	MaxAttempts        int                `json:"max_attempts"`
}

type DeliverPostToFollowersOptions struct {
	AccountsDatabase   AccountsDatabase
	FollowersDatabase  FollowersDatabase
	PostTagsDatabase   PostTagsDatabase
	NotesDatabase      NotesDatabase
	DeliveriesDatabase DeliveriesDatabase
	DeliveryQueue      DeliveryQueue
	Post               *Post
	PostTags           []*PostTag `json:"post_tags"`
	MaxAttempts        int        `json:"max_attempts"`
	URIs               *uris.URIs
}

// DeliverPostToFollowers schedules (en-queues) a `Post` instance to be delivered to all the
// accounts following the creator of the post.
//
// If the body of the post starts with the string "boost:" then the body is treated as a URI
// containing a pointer to the (ActivityPub) object (mostly likely a post/note) being boosted
// and an "?author_address=" query parameter referencing the author of that object and to whom
// the post (being delivered) will also be delivered. The body of the post (the "boost://" string)
// will be trapped and handled differently from "normal" posts in the `DeliverPost` method.
// Specifically it will be delivered as an ActivityPub "Announce" type rather than a "Note".
//
// None of this is ideal. It is a reflection of the intersection of the abstract-factory-pie nature
// of ActivityPub treating everything as a generic activity, the original goal of this package to get
// to basic "social media" style post/follow actions working and the mechanics how databases and
// queue need to be structured to do all that in practice. Eventually all the delivery mechanics
// will be refactored to just working ActvityPub "Activity" blobs but that is not the case today.
//
// This code does know (or care) what creates "boost:" post. It is most likely assumed to be custom
// code configured to read from a `ProcessMessageQueue` instance in the inbox handler but it could
// be defined somewhere else.
func DeliverPostToFollowers(ctx context.Context, opts *DeliverPostToFollowersOptions) error {

	logger := slog.Default()
	logger = logger.With("method", "DeliverPostToFollowers")
	logger = logger.With("post id", opts.Post.Id)
	logger = logger.With("account id", opts.Post.AccountId)

	logger.Info("Deliver post to followers")

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Post.AccountId)

	if err != nil {
		logger.Error("Failed to retrieve account ID for post", "error", err)
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	acct_address := acct.Address(opts.URIs.Hostname)
	logger = logger.With("account address", acct_address)

	followers_cb := func(ctx context.Context, follower_uri string) error {

		already_delivered := false

		deliveries_cb := func(ctx context.Context, d *Delivery) error {

			if d.Success {
				already_delivered = true
			}

			return nil
		}

		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, opts.Post.Id, follower_uri, deliveries_cb)

		if err != nil {
			logger.Info("Failed to retrieve deliveries for post and recipient", "recipient", follower_uri, "error", err)
			return fmt.Errorf("Failed to retrieve deliveries for post (%d) and recipient (%s), %w", opts.Post.Id, follower_uri, err)
		}

		if already_delivered {
			logger.Debug("Post already delivered", "post id", opts.Post.Id, "recipient", follower_uri)
			return nil
		}

		post_opts := &DeliverPostOptions{
			From:               acct,
			To:                 follower_uri,
			Post:               opts.Post,
			PostTags:           opts.PostTags,
			URIs:               opts.URIs,
			DeliveriesDatabase: opts.DeliveriesDatabase,
			MaxAttempts:        opts.MaxAttempts,
		}

		err = opts.DeliveryQueue.DeliverPost(ctx, post_opts)

		if err != nil {
			logger.Error("Failed to schedule post delivery", "recipient", follower_uri, "error", err)
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_uri, err)
		}

		logger.Info("Schedule post delivery", "recipient", follower_uri)
		return nil
	}

	err = opts.FollowersDatabase.GetFollowersForAccount(ctx, acct.Id, followers_cb)

	if err != nil {
		logger.Error("Failed to get followers for post author", "error", err)
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	// tags/mentions

	for _, t := range opts.PostTags {

		err := followers_cb(ctx, t.Name) // name or href?

		if err != nil {
			logger.Error("Failed to deliver message", "to", t.Name, "to id", t.Id, "error", err)
			return fmt.Errorf("Failed to deliver message to %s (%d), %w", t.Name, t.Id, err)
		}
	}

	// If this is a boost/announce post (hint) then extract the post author
	// from the ?author= parameter in order that we can send them a notification
	// of the boost. This is necessary in case they aren't already a follower (of acct)

	if strings.HasPrefix(opts.Post.Body, BOOST_URI_SCHEME) {

		logger.Info("Post is boost")

		parts := strings.SplitN(opts.Post.Body, " ", 2)

		if len(parts) != 2 {
			logger.Error("Invalid boost string")
			return fmt.Errorf("Invalid boost string")
		}

		boost_uri := parts[0]

		u, err := url.Parse(boost_uri)

		if err != nil {
			logger.Error("boost:// post body did not parse", "boost uri", boost_uri, "error", err)
			return fmt.Errorf("Invalid boost:// post body")
		}

		q := u.Query()

		author_addr := q.Get("author_address")

		_, _, err = ParseAddress(author_addr)

		if err != nil {
			logger.Error("Invalid author address", "author_address", author_addr, "error", err)
			return fmt.Errorf("Invalid author address")
		}

		// Invoke the followers_cb (which was set up for GetFollowersForAccount above)

		err = followers_cb(ctx, author_addr)

		if err != nil {
			logger.Error("Failed to deliver message", "error", err)
			return fmt.Errorf("Failed to deliver message to %s , %w", author_addr, err)
		}

	}

	// END OF this is no good to have to replicate this twice...

	return nil
}

// DeliverPost... TBD
// For posts with bodies starting with "boost:" see notes in `DeliverPostToFollowers` above.
func DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	logger := slog.Default()
	logger = logger.With("method", "DeliverPost")
	logger = logger.With("post", opts.Post.Id)
	logger = logger.With("from", opts.From.Id)
	logger = logger.With("to", opts.To)

	logger.Debug("Deliver post", "max attempts", opts.MaxAttempts)

	if opts.MaxAttempts > 0 {

		count_attempts := 0

		deliveries_cb := func(ctx context.Context, d *Delivery) error {
			count_attempts += 1
			return nil
		}

		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, opts.Post.Id, opts.To, deliveries_cb)

		if err != nil {
			logger.Error("Failed to count deliveries for post ID and recipient", "error", err)
			return fmt.Errorf("Failed to count deliveries for post ID and recipient, %w", err)
		}

		logger.Debug("Deliveries attempted", "count", count_attempts)

		if count_attempts >= opts.MaxAttempts {
			logger.Warn("Post has met or exceed max delivery attempts threshold", "max", opts.MaxAttempts, "count", count_attempts)
			return nil
		}
	}

	// Sort out dealing with Snowflake errors sooner...
	delivery_id, _ := id.NewId()

	now := time.Now()
	ts := now.Unix()

	d := &Delivery{
		Id:        delivery_id,
		PostId:    opts.Post.Id,
		AccountId: opts.From.Id, // This is still a bob@bob.com which suggests that we need to store actual inbox addresses...
		Recipient: opts.To,
		Created:   ts,
		Success:   false,
	}

	defer func() {

		now := time.Now()
		ts := now.Unix()

		d.Completed = ts

		logger.Info("Add delivery for post", "delivery id", d.PostId, "recipient", d.Recipient, "success", d.Success)

		err := opts.DeliveriesDatabase.AddDelivery(ctx, d)

		if err != nil {
			logger.Error("Failed to add delivery", "post_id", opts.Post.Id, "recipienct", d.Recipient, "error", err)
		}
	}()

	// START OF check what "kind" of post this is

	// See notes in post.go for why we are doing this. It's not great but it's not awful
	// either. It is a reasonable way to kick the can down the road a little further while
	// we continue to figure things out.

	logger.Info("Deliver post", "body", opts.Post.Body)

	var activity *ap.Activity

	if strings.HasPrefix(opts.Post.Body, BOOST_URI_SCHEME) {

		logger.Info("Post is boost")

		// Boost (announce) activities
		// https://boyter.org/posts/activitypub-announce-post/

		parts := strings.SplitN(opts.Post.Body, " ", 2)

		if len(parts) != 2 {
			logger.Error("Invalid boost string")
			return fmt.Errorf("Invalid boost string")
		}

		boost_uri := parts[0]
		boost_obj := parts[1]

		logger = logger.With("uri", boost_uri)
		logger = logger.With("object", boost_obj)

		u, err := url.Parse(boost_uri)

		if err != nil {
			logger.Error("boost:// post body did not parse", "error", err)
			return fmt.Errorf("Invalid boost:// post body")
		}

		q := u.Query()

		author_addr := q.Get("author_address")

		if author_addr == "" {
			logger.Error("Missing ?author_address= parameter")
			return fmt.Errorf("Missing ?author_address= parameter")
		}

		_, _, err = ParseAddress(author_addr)

		if err != nil {
			logger.Error("Invalid author address", "address", author_addr, "error", err)
			return fmt.Errorf("Invalid author address, %w", err)
		}

		// Apparently this is not necessary? As in the Announce 'cc' property takes
		// an address rather than a URL? Anyway, that's how I can get boosts to show
		// up in the Mastodon web application.

		/*
			author_uri := q.Get("author_uri")

			if author_uri == "" {
				logger.Error("Missing ?author_uri= parameter")
				return fmt.Errorf("Missing ?author_uri= parameter")
			}

			_, err = url.Parse(author_uri)

			if err != nil {
				logger.Error("Invalid author uri", "uri", author_uri, "error", err)
				return fmt.Errorf("Invalid author URI, %w", err)
			}
		*/

		from_uri := opts.From.AccountURL(ctx, opts.URIs).String()
		from_address := opts.From.Address(opts.URIs.Hostname)

		logger = logger.With("from", from_uri)

		boost_activity, err := ap.NewBoostActivity(ctx, opts.URIs, from_uri, author_addr, boost_obj)

		if err != nil {
			logger.Error("Failed to create boost activity", "error", err)
			return fmt.Errorf("Failed to create boost activity")
		}

		activity_id := fmt.Sprintf("%s#boost-from-%s", boost_obj, from_address)
		boost_activity.Id = activity_id

		enc_boost, _ := json.Marshal(boost_activity)
		logger.Info("BOOST", "activity", string(enc_boost))

		activity = boost_activity

	} else {

		// Note (create) activities

		note, err := NoteFromPost(ctx, opts.URIs, opts.From, opts.Post, opts.PostTags)

		if err != nil {
			d.Error = err.Error()
			return fmt.Errorf("Failed to derive note from post, %w", err)
		}

		from_uri := opts.From.AccountURL(ctx, opts.URIs).String()

		to_list := []string{
			opts.To,
		}

		create_activity, err := ap.NewCreateActivity(ctx, opts.URIs, from_uri, to_list, note)

		if err != nil {
			d.Error = err.Error()
			return fmt.Errorf("Failed to create activity from post, %w", err)
		}

		if len(note.Cc) > 0 {
			create_activity.Cc = note.Cc
		}

		// START OF is this really necessary?
		// Also, what if this isn't a post?

		uuid := id.NewUUID()

		post_url := opts.From.PostURL(ctx, opts.URIs, opts.Post)
		post_id := fmt.Sprintf("%s#%s", post_url.String(), uuid)

		create_activity.Id = post_id

		// END OF is this really necessary?

		activity = create_activity
	}

	// END OF check what "kind" of post this is...

	logger = logger.With("activity id", activity.Id)

	d.ActivityId = activity.Id

	post_opts := &PostToAccountOptions{
		From:     opts.From,
		To:       opts.To,
		Activity: activity,
		URIs:     opts.URIs,
	}

	inbox, err := PostToAccount(ctx, post_opts)

	d.Inbox = inbox

	if err != nil {
		logger.Error("Failed to post activity to account", "from", opts.From, "to", opts.To, "error", err)

		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox '%s', %w", opts.To, err)
	}

	d.Success = true

	logger.Info("Posted activity to account", "from", opts.From, "to", opts.To)
	return nil
}
