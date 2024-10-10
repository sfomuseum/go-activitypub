package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/inbox"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverActivityOptions struct {
	// This is what we used to do. Now we derive it from Activity.Actor
	// From               *activitypub.Account        `json:"from"`
	To       string       `json:"to"`
	Activity *ap.Activity `json:"activity"`
	// PostId is a misnomer. It is what unique 64-bit ID this package derives
	// for things stored in one of the "database" tables. The name is a reflection
	// of early attempts just to figure things out.
	PostId             int64                       `json:"post_id"`
	URIs               *uris.URIs                  `json:"uris"`
	AccountsDatabase   database.AccountsDatabase   `json:"accounts_database,omitempty"`
	DeliveriesDatabase database.DeliveriesDatabase `json:"deliveries_database,omitempty"`
	MaxAttempts        int                         `json:"max_attempts"`
}

type DeliverActivityToFollowersOptions struct {
	AccountsDatabase   database.AccountsDatabase
	FollowersDatabase  database.FollowersDatabase
	PostTagsDatabase   database.PostTagsDatabase
	NotesDatabase      database.NotesDatabase
	DeliveriesDatabase database.DeliveriesDatabase
	DeliveryQueue      DeliveryQueue
	Activity           *ap.Activity
	// PostId is a misnomer. It is what unique 64-bit ID this package derives
	// for things stored in one of the "database" tables. The name is a reflection
	// of early attempts just to figure things out.
	PostId      int64
	Mentions    []*activitypub.PostTag `json:"mentions"`
	MaxAttempts int                    `json:"max_attempts"`
	URIs        *uris.URIs
}

func DeliverActivityToFollowers(ctx context.Context, opts *DeliverActivityToFollowersOptions) error {

	logger := slog.Default()
	logger = logger.With("method", "DeliverActivityToFollowers")
	logger = logger.With("actor", opts.Activity.Actor)

	post_id := opts.PostId
	logger = logger.With("post id", post_id)

	logger.Info("Deliver post to followers")

	acct_name, _, err := activitypub.ParseAddress(opts.Activity.Actor)

	if err != nil {
		logger.Error("Failed to parse (actor) address", "error", err)
		return fmt.Errorf("Failed to parse (actor) address, %w", err)
	}

	acct, err := opts.AccountsDatabase.GetAccountWithName(ctx, acct_name)

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

		deliveries_cb := func(ctx context.Context, d *activitypub.Delivery) error {

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

		post_opts := &DeliverActivityOptions{
			To:                 follower_uri,
			Activity:           opts.Activity,
			PostId:             post_id,
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

	// tags/mentions...

	for _, t := range opts.Mentions {

		err := followers_cb(ctx, t.Name) // name or href?

		if err != nil {
			logger.Error("Failed to deliver message", "to", t.Name, "to id", t.Id, "error", err)
			return fmt.Errorf("Failed to deliver message to %s (%d), %w", t.Name, t.Id, err)
		}
	}

	return nil
}

// DeliverActivity... TBD
// For posts with bodies starting with "boost:" see notes in `DeliverActivityToFollowers` above.
func DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {

	activity := opts.Activity
	from := activity.Actor
	to := opts.To
	post_id := opts.PostId

	logger := slog.Default()
	logger = logger.With("method", "DeliverActivity")
	logger = logger.With("from", from)
	logger = logger.With("to", to)
	logger = logger.With("post id", post_id)

	logger.Info("Deliver activity to recipient")

	acct_name, _, err := activitypub.ParseAddress(from)

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

		err := opts.DeliveriesDatabase.GetDeliveriesWithPostIdAndRecipient(ctx, post_id, to, deliveries_cb)

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

	activity_id := fmt.Sprintf("%s#%s", activity.Type, activity.Id)

	logger = logger.With("delivery id", delivery_id)
	logger = logger.With("activity_id", activity_id)

	now := time.Now()
	ts := now.Unix()

	d := &activitypub.Delivery{
		Id:         delivery_id,
		ActivityId: activity_id,
		PostId:     post_id,
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

	recipient, err := activitypub.RetrieveActor(ctx, to, opts.URIs.Insecure)

	if err != nil {
		return fmt.Errorf("Failed to derive actor for to address, %w", err)
	}

	inbox_uri := recipient.Inbox
	d.Inbox = inbox_uri

	post_opts := &inbox.PostToInboxOptions{
		From:     acct,
		Inbox:    inbox_uri,
		Activity: activity,
		URIs:     opts.URIs,
	}

	err = inbox.PostToInbox(ctx, post_opts)

	if err != nil {
		logger.Error("Failed to post activity to inbox", "error", err)

		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox '%s', %w", recipient, err)
	}

	d.Success = true

	logger.Info("Posted activity to inbox")
	return nil
}
