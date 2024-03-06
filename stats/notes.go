package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

func CountNotesForDay(ctx context.Context, notes_db activitypub.NotesDatabase, day string) (int64, error) {

	t, err := time.Parse("2006-01-02", day)

	if err != nil {
		return 0, fmt.Errorf("Failed to parse day, %w", err)
	}

	start := t.Unix()
	end := start + ONEDAY

	return CountNotesForDateRange(ctx, notes_db, start, end)
}

func CountNotesForDateRange(ctx context.Context, notes_db activitypub.NotesDatabase, start int64, end int64) (int64, error) {

	count := int64(0)

	cb := func(ctx context.Context, id int64) error {
		count += 1
		return nil
	}

	err := notes_db.GetNoteIdsForDateRange(ctx, start, end, cb)

	if err != nil {
		return 0, fmt.Errorf("Failed to count notes, %w", err)
	}

	return count, nil
}
