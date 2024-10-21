package queue

// TBD replace all instances of ap.Activity with activitypub.Activity ?
// This would allow to get rid of all the PostId stuff...

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverActivityOptions struct {
	To                 string                      `json:"to"`
	Activity           *activitypub.Activity       `json:"activity"`
	URIs               *uris.URIs                  `json:"uris"`
	AccountsDatabase   database.AccountsDatabase   `json:"accounts_database,omitempty"`
	DeliveriesDatabase database.DeliveriesDatabase `json:"deliveries_database,omitempty"`
	MaxAttempts        int                         `json:"max_attempts"`
}

type DeliverActivityToFollowersOptions struct {
	AccountsDatabase   database.AccountsDatabase
	FollowersDatabase  database.FollowersDatabase
	NotesDatabase      database.NotesDatabase
	DeliveriesDatabase database.DeliveriesDatabase
	DeliveryQueue      DeliveryQueue
	Activity           *activitypub.Activity
	Mentions           []*activitypub.PostTag `json:"mentions"`
	MaxAttempts        int                    `json:"max_attempts"`
	URIs               *uris.URIs
}

func DeliverActivityToFollowers(ctx context.Context, opts *DeliverActivityToFollowersOptions) error {

	logger := slog.Default()
	logger = logger.With("activity id", opts.Activity.Id)
	logger = logger.With("account id", opts.Activity.AccountId)

	logger.Info("Deliver activity to followers")

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Activity.AccountId)

	if err != nil {
		logger.Error("Failed to retrieve account ID for post", "error", err)
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	// TBD - compare acct_address and opts.Activity.Actor?
	acct_address := acct.Address(opts.URIs.Hostname)
	logger = logger.With("account address", acct_address)

	followers_cb := func(ctx context.Context, follower_uri string) error {

		already_delivered := false

		deliveries_cb := func(ctx context.Context, d *activitypub.Delivery) error {

			if d.Success {
				already_delivered = true
			}

			return nil
		}

		// This will probably fail because types...?
		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, opts.Activity.Id, follower_uri, deliveries_cb)

		if err != nil {
			logger.Error("Failed to retrieve deliveries for post and recipient", "recipient", follower_uri, "error", err)
			return fmt.Errorf("Failed to retrieve deliveries for post (%d) and recipient (%s), %w", opts.Activity.Id, follower_uri, err)
		}

		if already_delivered {
			logger.Debug("Post already delivered", "recipient", follower_uri)
			return nil
		}

		post_opts := &DeliverActivityOptions{
			To:                 follower_uri,
			Activity:           opts.Activity,
			URIs:               opts.URIs,
			AccountsDatabase:   opts.AccountsDatabase,
			DeliveriesDatabase: opts.DeliveriesDatabase,
			MaxAttempts:        opts.MaxAttempts,
		}

		logger.Debug("Queue deliver activity", "to", follower_uri)

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

	// tags/mentions...

	for _, t := range opts.Mentions {

		logger.Debug("Deliver activity to mention", "mention", t)

		err := followers_cb(ctx, t.Name) // name or href?

		if err != nil {
			logger.Error("Failed to deliver message", "to", t.Name, "to id", t.Id, "error", err)
			return fmt.Errorf("Failed to deliver message to %s (%d), %w", t.Name, t.Id, err)
		}
	}

	ap_activity, err := opts.Activity.UnmarshalActivity()

	if err != nil {
		logger.Error("Failed to unmarshal activity", "error", err)
		return fmt.Errorf("Failed to unmarshal activity, %w", err)
	}

	for _, a := range ap_activity.Cc {

		logger.Debug("Deliver activity to cc", "address", a)

		err := followers_cb(ctx, a)

		if err != nil {
			logger.Error("Failed to deliver activity", "address", a, "error", err)
			return fmt.Errorf("Failed to deliver message to %s , %w", a, err)
		}

	}

	return nil
}

// DeliverActivity... TBD
// For posts with bodies starting with "boost:" see notes in `DeliverActivityToFollowers` above.
func DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {

	logger := slog.Default()
	logger = logger.With("activity id", opts.Activity.Id)

	ap_activity, err := opts.Activity.UnmarshalActivity()

	if err != nil {
		logger.Error("Failed to unmarshal activity", "error", err)
		return fmt.Errorf("Failed to unmarshal activity, %w", err)
	}

	from := ap_activity.Actor
	to := opts.To

	logger = logger.With("from", from)
	logger = logger.With("to", to)

	logger.Info("Deliver activity to recipient")

	acct_name, _, err := ap.ParseAddress(from)

	if err != nil {
		logger.Error("Failed to parse (actor) address", "error", err)
		return fmt.Errorf("Failed to parse (actor) address, %w", err)
	}

	logger = logger.With("account name", acct_name)

	acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, acct_name)

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

		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, opts.Activity.Id, to, deliveries_cb)

		if err != nil {
			logger.Error("Failed to count deliveries for \"post\" ID and recipient", "error", err)
			return fmt.Errorf("Failed to count deliveries for \"post\" ID and recipient, %w", err)
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

	d := &activitypub.Delivery{
		Id: delivery_id,
		// FIX ME activity ID and activitypub ID...
		ActivityId: opts.Activity.ActivityPubId,
		AccountId:  acct.Id,
		Recipient:  to,
		Created:    ts,
		Success:    false,
	}

	defer func() {

		now := time.Now()
		ts := now.Unix()

		d.Completed = ts

		logger.Info("Add delivery for activity", "success", d.Success)

		err := opts.DeliveriesDatabase.AddDelivery(ctx, d)

		if err != nil {
			logger.Error("Failed to add delivery", "error", err)
		}
	}()

	recipient, err := ap.RetrieveActor(ctx, to, opts.URIs.Insecure)

	if err != nil {
		return fmt.Errorf("Failed to derive actor for to address, %w", err)
	}

	inbox_uri := recipient.Inbox
	d.Inbox = inbox_uri

	err = acct.SendActivity(ctx, opts.URIs, inbox_uri, ap_activity)

	if err != nil {
		logger.Error("Failed to post activity to inbox", "error", err)

		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox '%s', %w", recipient, err)
	}

	d.Success = true

	logger.Info("Posted activity to inbox")
	return nil
}
