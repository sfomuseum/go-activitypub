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

	ctx := context.Background()

	RegisterMessagesDatabase(ctx, "awsdynamodb", NewDocstoreMessagesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterMessagesDatabase(ctx, scheme, NewDocstoreMessagesDatabase)
	}
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

func (db *DocstoreMessagesDatabase) GetMessageIdsForDateRange(ctx context.Context, start int64, end int64, cb GetMessageIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var m Message
		err := iter.Next(ctx, &m)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, m.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for message %d, %w", m.Id, err)
			}
		}
	}

	return nil
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

func (db *DocstoreMessagesDatabase) GetMessagesForAccount(ctx context.Context, account_id int64, callback_func GetMessagesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	return db.getMessagesWithCallback(ctx, q, callback_func)
}

func (db *DocstoreMessagesDatabase) GetMessagesForAccountAndAuthor(ctx context.Context, account_id int64, author_address string, callback_func GetMessagesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("AuthorAddress", "=", author_address)

	return db.getMessagesWithCallback(ctx, q, callback_func)
}

func (db *DocstoreMessagesDatabase) getMessagesWithCallback(ctx context.Context, q *gc_docstore.Query, callback_func GetMessagesCallbackFunc) error {

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

			err := callback_func(ctx, &m)

			if err != nil {
				return fmt.Errorf("Failed to execute following callback for message '%d', %w", m.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreMessagesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
