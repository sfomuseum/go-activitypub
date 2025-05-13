package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullMessagesDatabase struct {
	MessagesDatabase
}

func init() {

	ctx := context.Background()

	err := RegisterMessagesDatabase(ctx, "null", NewNullMessagesDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullMessagesDatabase(ctx context.Context, uri string) (MessagesDatabase, error) {

	db := &NullMessagesDatabase{}
	return db, nil
}

func (db *NullMessagesDatabase) GetMessageIdsForDateRange(ctx context.Context, start int64, end int64, cb GetMessageIdsCallbackFunc) error {
	return nil
}

func (db *NullMessagesDatabase) GetMessageWithId(ctx context.Context, message_id int64) (*activitypub.Message, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullMessagesDatabase) GetMessageWithAccountAndNoteIds(ctx context.Context, account_id int64, note_id int64) (*activitypub.Message, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullMessagesDatabase) AddMessage(ctx context.Context, message *activitypub.Message) error {
	return nil
}

func (db *NullMessagesDatabase) UpdateMessage(ctx context.Context, message *activitypub.Message) error {
	return nil
}

func (db *NullMessagesDatabase) RemoveMessage(ctx context.Context, message *activitypub.Message) error {
	return nil
}

func (db *NullMessagesDatabase) GetMessagesForAccount(ctx context.Context, account_id int64, callback_func GetMessagesCallbackFunc) error {
	return nil
}

func (db *NullMessagesDatabase) GetMessagesForAccountAndAuthor(ctx context.Context, account_id int64, author_address string, callback_func GetMessagesCallbackFunc) error {
	return nil
}

func (db *NullMessagesDatabase) Close(ctx context.Context) error {
	return nil
}
