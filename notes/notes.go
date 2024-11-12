package notes

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func AddNote(ctx context.Context, db database.NotesDatabase, uuid string, author string, body string) (*activitypub.Note, error) {

	n, err := activitypub.NewNote(ctx, uuid, author, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new note, %w", err)
	}

	err = db.AddNote(ctx, n)

	if err != nil {
		return nil, fmt.Errorf("Failed to add note, %w", err)
	}

	return n, nil
}
