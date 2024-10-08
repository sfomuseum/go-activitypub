package notes

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/database"
)

func AddNote(ctx context.Context, db database.NotesDatabase, uuid string, author string, body string) (*Note, error) {

	n, err := NewNote(ctx, uuid, author, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new note, %w", err)
	}

	err = db.AddNote(ctx, n)

	if err != nil {
		return nil, fmt.Errorf("Failed to add note, %w", err)
	}

	slog.Debug("Return new note", "id", n.Id)
	return n, nil
}
