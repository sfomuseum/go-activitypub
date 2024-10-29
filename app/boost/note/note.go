package note

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/queue"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	activities_db, err := database.NewActivitiesDatabase(ctx, opts.ActivitiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new activities database, %w", err)
	}

	defer activities_db.Close(ctx)

	followers_db, err := database.NewFollowersDatabase(ctx, opts.FollowersDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate followers database, %w", err)
	}

	defer followers_db.Close(ctx)

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	delivery_q, err := queue.NewDeliveryQueue(ctx, opts.DeliveryQueueURI)

	if err != nil {
		return fmt.Errorf("Failed to create new delivery queue, %w", err)
	}

	logger = logger.With("account", opts.AccountName)

	acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	from_u := acct.AccountURL(opts.URIs.Hostname)
	from := from_u.String()
	
	ap_activity, err := ap.NewBoostActivityForNote(ctx, opts.URIs, from, opts.NoteURI)

	if err != nil {
		return fmt.Errorf("Failed to create new boost AP activity, %w", err)
	}

	activity, err := activitypub.NewActivity(ctx, ap_activity)

	if err != nil {
		return fmt.Errorf("Failed to create new book AP wrapper, %w", err)
	}

	boost_id, err := id.NewId()

	if err != nil {
		return fmt.Errorf("Failed to create new boost ID, %w", err)
	}

	activity.ActivityType = activitypub.BoostActivityType
	activity.ActivityTypeId = boost_id
	activity.AccountId = acct.Id

	logger = logger.With("activity id", activity.Id)
	logger = logger.With("boost id", boost_id)

	err = activities_db.AddActivity(ctx, activity)

	if err != nil {
		logger.Error("Failed to add new activity", "error", err)
		return fmt.Errorf("Failed to add new activity, %w", err)
	}

	deliver_opts := &queue.DeliverActivityToFollowersOptions{
		AccountsDatabase:   accounts_db,
		FollowersDatabase:  followers_db,
		DeliveriesDatabase: deliveries_db,
		DeliveryQueue:      delivery_q,
		Activity:           activity,
		URIs:               opts.URIs,
	}

	logger.Debug("Deliver activity")

	err = queue.DeliverActivityToFollowers(ctx, deliver_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver post, %w", err)
	}

	logger.Info("Delivered boost")
	return nil
}
