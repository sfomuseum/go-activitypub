package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreMessagesDatabase struct {
	MessagesDatabase
	collection *gc_docstore.Collection
}

func init() {
	// ctx := context.Background()
}

func NewDocstoreMessagesDatabase(ctx context.Context, uri string) (MessagesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreMessagesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreMessagesDatabase) GetMessageWithId(ctx context.Context, message_id int64) (*Message, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", message_id)

	return db.getMessage(ctx, q)
}

func (db *DocstoreMessagesDatabase) GetMessageWithAccountAndNoteIds(ctx context.Context, account_id int64, note_id int64) (*Message, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("NoteId", "=", note_id)

	return db.getMessage(ctx, q)
}

func (db *DocstoreMessagesDatabase) getMessage(ctx context.Context, q *gc_docstore.Query) (*Message, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var m Message
		err := iter.Next(ctx, &m)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &m, nil
		}
	}

	return nil, ErrNotFound
}

func (db *DocstoreMessagesDatabase) AddMessage(ctx context.Context, message *Message) error {

	return db.collection.Put(ctx, message)
}

func (db *DocstoreMessagesDatabase) UpdateMessage(ctx context.Context, message *Message) error {

	return db.collection.Replace(ctx, message)
}

func (db *DocstoreMessagesDatabase) RemoveMessage(ctx context.Context, message *Message) error {

	return db.collection.Delete(ctx, message)
}

func (db *DocstoreMessagesDatabase) GetMessagesForAccount(ctx context.Context, account_id int64, following_callback GetMessagesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var m Message
		err := iter.Next(ctx, &m)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := following_callback(ctx, &m)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for message '%s', %w", m.Id, err)
			}
		}
	}

	return nil
}
