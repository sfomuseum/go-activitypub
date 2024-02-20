package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Note struct {
	Id            int64  `json:"id"`
	UUID          string `json:"uuid"`
	AuthorAddress string `json:"author_address"`
	Body          []byte `json:"body"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}

func AddNote(ctx context.Context, db NotesDatabase, uuid string, author string, body []byte) (*Note, error) {

	slog.Debug("Add note", "uuid", uuid, "author", author)

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

func NewNote(ctx context.Context, uuid string, author string, body []byte) (*Note, error) {

	slog.Debug("Create new note", "uuid", uuid, "author", author)

	db_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new note ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	n := &Note{
		Id:            db_id,
		UUID:          uuid,
		AuthorAddress: author,
		Body:          body,
		Created:       ts,
		LastModified:  ts,
	}

	return n, nil
}
