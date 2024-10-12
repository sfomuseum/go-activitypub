package note

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

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
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer accounts_db.Close(ctx)

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

	from := acct.Address(opts.URIs.Hostname)

	boost_activity, err := ap.NewBoostActivityForNote(ctx, opts.URIs, from, opts.NoteURI)

	if err != nil {
		return fmt.Errorf("Failed to create new boost activity, %v", err)
	}

	// FIX ME. See notes in sfomuseum/go-activitypub/boost.go
	boost_id, _ := id.NewId()

	logger = logger.With("activity id", boost_activity.Id)

	deliver_opts := &queue.DeliverActivityToFollowersOptions{
		AccountsDatabase:   accounts_db,
		FollowersDatabase:  followers_db,
		DeliveriesDatabase: deliveries_db,
		DeliveryQueue:      delivery_q,
		Activity:           boost_activity,
		PostId:             boost_id,
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
