package activitypub

import (
	"context"
)

type NullBlocksDatabase struct {
	BlocksDatabase
}

func init() {
	ctx := context.Background()
	RegisterBlocksDatabase(ctx, "null", NewNullBlocksDatabase)
}

func NewNullBlocksDatabase(ctx context.Context, uri string) (BlocksDatabase, error) {
	db := &NullBlocksDatabase{}
	return db, nil
}

func (db *NullBlocksDatabase) IsBlockedByAccount(ctx context.Context, account_id int64, host string, name string) (bool, error) {
	return false, nil
}

func (db *NullBlocksDatabase) GetBlockWithId(ctx context.Context, block_id int64) (*Block, error) {
	return nil, ErrNotFound
}

func (db *NullBlocksDatabase) GetBlockWithAccountIdAndAddress(ctx context.Context, account_id int64, host string, name string) (*Block, error) {
	return nil, ErrNotFound
}

func (db *NullBlocksDatabase) AddBlock(ctx context.Context, block *Block) error {
	return nil
}

func (db *NullBlocksDatabase) UpdateBlock(ctx context.Context, block *Block) error {
	return nil
}

func (db *NullBlocksDatabase) RemoveBlock(ctx context.Context, block *Block) error {
	return nil
}
