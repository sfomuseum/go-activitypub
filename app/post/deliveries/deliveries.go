package deliveries

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
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

	logger := slog.Default()

	deliveries_db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create instantiate deliveries database, %w", err)
	}

	defer deliveries_db.Close(ctx)

	count := 0

	deliveries_cb := func(ctx context.Context, d *activitypub.Delivery) error {
		logger.Info("D", "d", d)
		count += 1
		return nil
	}

	err = deliveries_db.GetDeliveriesWithPostIdAndRecipient(ctx, opts.PostId, opts.Recipient, deliveries_cb)

	if err != nil {
		return fmt.Errorf("Failed to load deliveries, %w", err)
	}

	logger.Info(fmt.Sprintf("%d deliveries", count))
	return nil
}
