package retrieve

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
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

	notes_db, err := activitypub.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	defer notes_db.Close(ctx)

	n, err := notes_db.GetNoteWithId(ctx, opts.NoteId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve note, %w", err)
	}

	var target interface{}

	if opts.Body {

		var activity *ap.Note

		err := json.Unmarshal([]byte(n.Body), &activity)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal body, %w", err)
		}

		target = activity
	} else {
		target = n
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(target)

	if err != nil {
		return fmt.Errorf("Failed to encode note, %w", err)
	}

	return nil
}
