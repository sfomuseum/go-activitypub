package stats

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub/database"
)

func CountNotesForDateRange(ctx context.Context, notes_db database.NotesDatabase, start int64, end int64) (int64, error) {

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
