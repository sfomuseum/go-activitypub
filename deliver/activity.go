package deliver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/queue"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverActivityToFollowersOptions struct {
	AccountsDatabase   database.AccountsDatabase
	FollowersDatabase  database.FollowersDatabase
	PostTagsDatabase   database.PostTagsDatabase
	NotesDatabase      database.NotesDatabase
	DeliveriesDatabase database.DeliveriesDatabase
	DeliveryQueue      queue.DeliveryQueue
	Activity           *ap.Activity
	// PostTags           []*PostTag `json:"post_tags"`
	MaxAttempts int `json:"max_attempts"`
	URIs        *uris.URIs
}

func DeliverActivityToFollowers(ctx context.Context, opts *DeliverActivityToFollowersOptions) error {

	logger := slog.Default()
	logger = logger.With("method", "DeliverActivityToFollowers")
	logger = logger.With("actor", opts.Activity.Actor)

	post_id := fmt.Sprintf("%s-%s", opts.Activity.Type, opts.Activity.Id)
	logger = logger.With("post id", post_id)

	logger.Info("Deliver post to followers")

	acct_name, _, err := activitypub.ParseAddress(opts.Activity.Actor)

	if err != nil {
		logger.Error("Failed to parse (actor) address", "error", err)
		return fmt.Errorf("Failed to parse (actor) address, %w", err)
	}

	acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

	if err != nil {
		logger.Error("Failed to retrieve account ID for post", "error", err)
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	logger = logger.With("account id", acct.Id)

	// TBD - compare acct_address and opts.Activity.Actor?
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

		// This will probably fail because types...?
		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, post_id, follower_uri, deliveries_cb)

		if err != nil {
			logger.Error("Failed to retrieve deliveries for post and recipient", "recipient", follower_uri, "error", err)
			return fmt.Errorf("Failed to retrieve deliveries for post (%d) and recipient (%s), %w", post_id, follower_uri, err)
		}

		if already_delivered {
			logger.Debug("Post already delivered", "recipient", follower_uri)
			return nil
		}

		post_opts := &queue.DeliverActivityOptions{
			From:     acct,
			To:       follower_uri,
			Activity: opts.Activity,
			// PostTags:           opts.PostTags,
			URIs:               opts.URIs,
			DeliveriesDatabase: opts.DeliveriesDatabase,
			MaxAttempts:        opts.MaxAttempts,
		}

		err = opts.DeliveryQueue.DeliverActivity(ctx, post_opts)

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

	// tags/mentions... TBD...

	/*
		for _, t := range opts.PostTags {

			err := followers_cb(ctx, t.Name) // name or href?

			if err != nil {
				logger.Error("Failed to deliver message", "to", t.Name, "to id", t.Id, "error", err)
				return fmt.Errorf("Failed to deliver message to %s (%d), %w", t.Name, t.Id, err)
			}
		}
	*/

	return nil
}

// DeliverActivity... TBD
// For posts with bodies starting with "boost:" see notes in `DeliverActivityToFollowers` above.
func DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {

	actor := opts.Activity.Actor
	recipient := opts.To // TBD...

	post_id := fmt.Sprintf("%s-%s", opts.Activity.Type, opts.Activity.Id)

	logger := slog.Default()
	logger = logger.With("method", "DeliverActivity")
	logger = logger.With("actor", opts.Activity.Actor)
	logger = logger.With("recipient", recipient)
	logger = logger.With("post id", post_id)

	logger.Info("Deliver activity to recipient")

	acct_name, _, err := activitypub.ParseAddress(opts.Activity.Actor)

	if err != nil {
		logger.Error("Failed to parse (actor) address", "error", err)
		return fmt.Errorf("Failed to parse (actor) address, %w", err)
	}

	acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, account_name)

	if err != nil {
		logger.Error("Failed to retrieve account ID for post", "error", err)
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	logger = logger.With("account id", acct.Id)

	logger.Debug("Deliver activity", "max attempts", opts.MaxAttempts)

	if opts.MaxAttempts > 0 {

		count_attempts := 0

		deliveries_cb := func(ctx context.Context, d *activitypub.Delivery) error {
			count_attempts += 1
			return nil
		}

		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, post_id, recipient, deliveries_cb)

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

	logger = logger.With("delivery id", delivery_id)

	now := time.Now()
	ts := now.Unix()

	d := &activitypubDelivery{
		Id:        delivery_id,
		PostId:    post_id,
		AccountId: actor,
		Recipient: recipient,
		Created:   ts,
		Success:   false,
	}

	defer func() {

		now := time.Now()
		ts := now.Unix()

		d.Completed = ts

		logger.Info("Add delivery for post", "success", d.Success)

		err := opts.DeliveriesDatabase.AddDelivery(ctx, d)

		if err != nil {
			logger.Error("Failed to add delivery", "error", err)
		}
	}()

	logger = logger.With("activity id", activity.Id)

	d.ActivityId = activity.Id

	post_opts := &PostToInbox{
		From:     acct,
		To:       recipient,
		Activity: activity,
		URIs:     opts.URIs,
	}

	inbox, err := PostToInbox(ctx, post_opts)

	d.Inbox = inbox

	if err != nil {
		logger.Error("Failed to post activity to inbox", "error", err)

		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox '%s', %w", recipient, err)
	}

	d.Success = true

	logger.Info("Posted activity to inbox")
	return nil
}
