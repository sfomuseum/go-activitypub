// Package deliver provides methods for delivering activities to external actors.
package deliver

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

// DeliveryActivityOptions defines configuration options for delivering an `activitypub.Activity` instance to an external actor.
type DeliverActivityOptions struct {
	// To is the ActivityPub address of the external actor to deliver the `activitypub.Activity` instance to.
	To string `json:"to"`
	// Activity is the `activitypub.Activity` being delivered.
	Activity *activitypub.Activity `json:"activity"`
	// URIs is a `uris.URIs` instance containing host and domain specific information of the service delivering the activity.
	URIs *uris.URIs `json:"uris"`
	// AccountsDatabase is a `database.AccountsDatabase` instance used to lookup details about the account sending the activity.
	AccountsDatabase database.AccountsDatabase `json:"accounts_database,omitempty"`
	// DelivereisDatabase is a `database.DeliveriesDatabase` instance used to store and retrieve details about the delivery.
	DeliveriesDatabase database.DeliveriesDatabase `json:"deliveries_database,omitempty"`
	// MaxAttempts is the maximum number of times to attempt to deliver the activity.
	MaxAttempts int `json:"max_attempts"`
}

// DeliveryActivity attempts to deliver an `activitypub.Activity` instance to an external actor.
func DeliverActivity(ctx context.Context, opts *DeliverActivityOptions) error {

	logger := slog.Default()
	logger = logger.With("activity id", opts.Activity.Id)

	ap_activity, err := opts.Activity.UnmarshalActivity()

	if err != nil {
		logger.Error("Failed to unmarshal activity", "error", err)
		return fmt.Errorf("Failed to unmarshal activity, %w", err)
	}

	from_uri := ap_activity.Actor
	to := opts.To

	logger = logger.With("from", from_uri)
	logger = logger.With("to", to)

	logger.Info("Deliver activity to recipient")

	// I guess we could just assume that the tail end is the account name but...
	// Note that in ap.ParseAddressFromRequest we rely on the Go 1.22 net/http
	// {resource} placeholder to derive the name...
	actor, err := ap.RetrieveActorWithProfileURL(ctx, from_uri)

	if err != nil {
		logger.Error("Failed to retrieve actor for profile (from) URI", "error", err)
		return fmt.Errorf("Failed to retrieve actor for profile (from) URI, %w", err)
	}

	/*
		acct_name, _, err := ap.ParseAddress(from_uri)

		if err != nil {
			logger.Error("Failed to parse (actor) address", "error", err)
			return fmt.Errorf("Failed to parse (actor) address, %w", err)
		}
	*/

	acct_name := actor.PreferredUsername

	logger = logger.With("account name", acct_name)
	logger.Debug("Lookup account for actor name")

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

		err := opts.DeliveriesDatabase.GetDeliveriesWithActivityIdAndRecipient(ctx, opts.Activity.Id, to, deliveries_cb)

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
		Id:            delivery_id,
		ActivityId:    opts.Activity.Id,
		ActivityPubId: opts.Activity.ActivityPubId,
		AccountId:     acct.Id,
		Recipient:     to,
		Created:       ts,
		Success:       false,
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
		logger.Error("Failed to retrieve (to) actor", "error", err)

		d.Error = err.Error()
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
