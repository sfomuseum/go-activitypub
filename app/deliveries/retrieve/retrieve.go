package retrieve

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-activitypub"
	ap_slog "github.com/sfomuseum/go-activitypub/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	ap_slog.ConfigureLogger(logger, opts.Verbose)

	db, err := activitypub.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate deliveries database, %w", err)
	}

	defer db.Close(ctx)

	// To do: Support retrieving deliveries for post, account, etc.

	d, err := db.GetDeliveryWithId(ctx, opts.DeliveryId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve delivery, %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(d)

	if err != nil {
		return fmt.Errorf("Failed to encode delivery, %w", err)
	}

	return nil
}
