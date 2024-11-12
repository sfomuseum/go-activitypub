package activitypub

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub/id"
)

type Message struct {
	Id            int64  `json:"id"`
	NoteId        int64  `json:"note_id"`
	AuthorAddress string `json:"author_uri"`
	AccountId     int64  `json:"account_id"`
	Created       int64  `json:"created"`
	LastModified  int64  `json:"created"`
}

func NewMessage(ctx context.Context, account_id int64, note_id int64, author_address string) (*Message, error) {

	db_id, err := id.NewId()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new message ID, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	m := &Message{
		Id:            db_id,
		NoteId:        note_id,
		AccountId:     account_id,
		AuthorAddress: author_address,
		Created:       ts,
		LastModified:  ts,
	}

	return m, nil
}
