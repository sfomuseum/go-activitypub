package database

import (
	"context"

	"github.com/sfomuseum/go-activitypub"
)

type NullBlocksDatabase struct {
	BlocksDatabase
}

func init() {
	ctx := context.Background()
	err := RegisterBlocksDatabase(ctx, "null", NewNullBlocksDatabase)

	if err != nil {
		panic(err)
	}
}

func NewNullBlocksDatabase(ctx context.Context, uri string) (BlocksDatabase, error) {
	db := &NullBlocksDatabase{}
	return db, nil
}

func (db *NullBlocksDatabase) GetBlockIdsForDateRange(ctx context.Context, start int64, end int64, cb GetBlockIdsCallbackFunc) error {
	return nil
}

func (db *NullBlocksDatabase) IsBlockedByAccount(ctx context.Context, account_id int64, host string, name string) (bool, error) {
	return false, nil
}

func (db *NullBlocksDatabase) GetBlockWithId(ctx context.Context, block_id int64) (*activitypub.Block, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullBlocksDatabase) GetBlockWithAccountIdAndAddress(ctx context.Context, account_id int64, host string, name string) (*activitypub.Block, error) {
	return nil, activitypub.ErrNotFound
}

func (db *NullBlocksDatabase) AddBlock(ctx context.Context, block *activitypub.Block) error {
	return nil
}

func (db *NullBlocksDatabase) UpdateBlock(ctx context.Context, block *activitypub.Block) error {
	return nil
}

func (db *NullBlocksDatabase) RemoveBlock(ctx context.Context, block *activitypub.Block) error {
	return nil
}
