package activitypub

import (
	"context"
	"fmt"
	"log/slog"
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

func GetMessage(ctx context.Context, db MessagesDatabase, account_id int64, note_id int64) (*Message, error) {

	slog.Debug("Get message", "account", account_id, "note", note_id)
	return db.GetMessageWithAccountAndNoteIds(ctx, account_id, note_id)
}

func AddMessage(ctx context.Context, db MessagesDatabase, account_id int64, note_id int64, author_address string) (*Message, error) {

	slog.Debug("Add message", "account", account_id, "note", note_id, "author", author_address)

	m, err := NewMessage(ctx, account_id, note_id, author_address)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new message, %w", err)
	}

	err = db.AddMessage(ctx, m)

	if err != nil {
		return nil, fmt.Errorf("Failed to add message, %w", err)
	}

	return m, nil
}

func UpdateMessage(ctx context.Context, db MessagesDatabase, m *Message) (*Message, error) {

	slog.Debug("Update message", "id", m.Id)

	now := time.Now()
	ts := now.Unix()

	m.LastModified = ts

	err := db.UpdateMessage(ctx, m)
	return m, err
}

func NewMessage(ctx context.Context, account_id int64, note_id int64, author_address string) (*Message, error) {

	slog.Debug("Create new message", "account", account_id, "note", note_id, "author", author_address)

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
