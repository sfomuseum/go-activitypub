package messages

import (
	"context"
	"fmt"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func GetMessage(ctx context.Context, db database.MessagesDatabase, account_id int64, note_id int64) (*activitypub.Message, error) {

	return db.GetMessageWithAccountAndNoteIds(ctx, account_id, note_id)
}

func AddMessage(ctx context.Context, db database.MessagesDatabase, account_id int64, note_id int64, author_address string) (*activitypub.Message, error) {

	m, err := activitypub.NewMessage(ctx, account_id, note_id, author_address)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new message, %w", err)
	}

	err = db.AddMessage(ctx, m)

	if err != nil {
		return nil, fmt.Errorf("Failed to add message, %w", err)
	}

	return m, nil
}

func UpdateMessage(ctx context.Context, db database.MessagesDatabase, m *activitypub.Message) (*activitypub.Message, error) {

	now := time.Now()
	ts := now.Unix()

	m.LastModified = ts

	err := db.UpdateMessage(ctx, m)
	return m, err
}
