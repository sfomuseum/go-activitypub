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

	actor, err := activitypub.RetrieveActor(ctx, opts.Address, opts.Insecure)

	if err != nil {
		return fmt.Errorf("Failed to retrieve actor, %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(actor)

	if err != nil {
		return fmt.Errorf("Failed to encode actor, %w", err)
	}

	return nil
}
