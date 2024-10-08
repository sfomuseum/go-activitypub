package retrieve

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/slog"
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

	db, err := database.NewDeliveriesDatabase(ctx, opts.DeliveriesDatabaseURI)

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
