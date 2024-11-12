package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/deliver"
	"github.com/sfomuseum/go-activitypub/uris"
)

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
	logger = logger.With("activity type", opts.Activity.ActivityType)
	logger = logger.With("activity type id", opts.Activity.ActivityTypeId)
	logger = logger.With("account id", opts.Activity.AccountId)

	logger.Info("Deliver activity to followers")

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Activity.AccountId)

	if err != nil {
		logger.Error("Failed to retrieve account ID for post", "error", err)
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	// TBD - compare acct_address and opts.Activity.Actor?
	acct_address := acct.Address(opts.URIs.Hostname)
	logger = logger.With("from address", acct_address)

	followers_cb := func(ctx context.Context, follower_address string) error {

		logger.Info("Process follower (for activity deliver)", "follower", follower_address)

		already_delivered := false

		deliveries_cb := func(ctx context.Context, d *activitypub.Delivery) error {

			if d.Success {
				logger.Info("Delivery (to follower) already happened", "delivery id", d.Id, "activity id", d.ActivityId, "recipient", d.Recipient)
				already_delivered = true
			}

			return nil
		}

		// This will probably fail because types...?
		err := opts.DeliveriesDatabase.GetDeliveriesWithActivityIdAndRecipient(ctx, opts.Activity.Id, follower_address, deliveries_cb)

		if err != nil {
			logger.Error("Failed to retrieve deliveries for post and recipient", "recipient", follower_address, "error", err)
			return fmt.Errorf("Failed to retrieve deliveries for post (%d) and recipient (%s), %w", opts.Activity.Id, follower_address, err)
		}

		if already_delivered {
			logger.Info("Activity already delivered", "recipient", follower_address)
			return nil
		}

		deliver_opts := &deliver.DeliverActivityOptions{
			To:                 follower_address,
			Activity:           opts.Activity,
			URIs:               opts.URIs,
			AccountsDatabase:   opts.AccountsDatabase,
			DeliveriesDatabase: opts.DeliveriesDatabase,
			MaxAttempts:        opts.MaxAttempts,
		}

		logger.Info("Queue deliver activity", "to", follower_address)

		err = opts.DeliveryQueue.DeliverActivity(ctx, deliver_opts)

		if err != nil {
			logger.Error("Failed to schedule post delivery", "recipient", follower_address, "error", err)
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_address, err)
		}

		logger.Info("Activity delivery complete", "recipient", follower_address)
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

		logger.Info("Deliver activity to cc", "address", a)

		err := followers_cb(ctx, a)

		if err != nil {
			logger.Error("Failed to deliver activity", "address", a, "error", err)
			return fmt.Errorf("Failed to deliver message to %s , %w", a, err)
		}

	}

	return nil
}
