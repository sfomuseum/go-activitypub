package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

// Note is a message (or post) delivered to an account.
type Note struct {
	Id            int64  `json:"id"`
	UUID          string `json:"uuid"`
	AuthorAddress string `json:"author_address"`
	Body          string `json:"body"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"lastmodified"`
}

func NewNote(ctx context.Context, uuid string, author string, body string) (*Note, error) {

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
