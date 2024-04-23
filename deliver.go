package activitypub

import (
	"context"
	"fmt"
	"log/slog"
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
	DeliveriesDatabase DeliveriesDatabase
	DeliveryQueue      DeliveryQueue
	Post               *Post
	PostTags           []*PostTag `json:"post_tags"`
	MaxAttempts        int        `json:"max_attempts"`
	URIs               *uris.URIs
}

func DeliverPostToFollowers(ctx context.Context, opts *DeliverPostToFollowersOptions) error {

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Post.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

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
			return fmt.Errorf("Failed to retrieve deliveries for post (%d) and recipient (%s), %w", opts.Post.Id, follower_uri, err)
		}

		if already_delivered {
			slog.Debug("Post already delivered", "post id", opts.Post.Id, "recipient", follower_uri)
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
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_uri, err)
		}

		return nil
	}

	err = opts.FollowersDatabase.GetFollowersForAccount(ctx, acct.Id, followers_cb)

	if err != nil {
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	// tags/mentions

	for _, t := range opts.PostTags {

		err := followers_cb(ctx, t.Name) // name or href?

		if err != nil {
			return fmt.Errorf("Failed to deliver message to %s (%d), %w", t.Name, t.Id, err)
		}
	}

	return nil
}

func DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	logger := slog.Default()
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

	d.ActivityId = create_activity.Id

	post_opts := &PostToAccountOptions{
		From:     opts.From,
		To:       opts.To,
		Activity: create_activity,
		URIs:     opts.URIs,
	}

	inbox, err := PostToAccount(ctx, post_opts)

	d.Inbox = inbox

	if err != nil {
		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox '%s', %w", opts.To, err)
	}

	d.Success = true
	return nil
}
