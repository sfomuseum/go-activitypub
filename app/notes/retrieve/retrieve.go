package retrieve

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
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

	notes_db, err := database.NewNotesDatabase(ctx, opts.NotesDatabaseURI)

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
