package list

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

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	activities_db, err := database.NewActivitiesDatabase(ctx, opts.ActivitiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create activities database, %w", err)
	}

	defer activities_db.Close(ctx)

	cb := func(ctx context.Context, a *activitypub.Activity) error {

		slog.Info("LOG", "activity", a)
		return nil
	}

	err = activities_db.GetActivities(ctx, cb)

	if err != nil {
		return fmt.Errorf("Failed to get activities, %w", err)
	}

	return nil
}
