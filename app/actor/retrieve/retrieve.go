package retrieve

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub/ap"
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

	// logger := slog.Default()

	actor, err := ap.RetrieveActor(ctx, opts.Address, opts.Insecure)

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
