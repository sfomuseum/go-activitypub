package activitypub

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreBlocksDatabase struct {
	BlocksDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterBlocksDatabase(ctx, "awsdynamodb", NewDocstoreBlocksDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterBlocksDatabase(ctx, scheme, NewDocstoreBlocksDatabase)
	}
}

func NewDocstoreBlocksDatabase(ctx context.Context, uri string) (BlocksDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreBlocksDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreBlocksDatabase) GetBlockIdsForDateRange(ctx context.Context, start int64, end int64, cb GetBlockIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var b Block
		err := iter.Next(ctx, &b)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, b.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for block %d, %w", b.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreBlocksDatabase) GetBlockWithId(ctx context.Context, id int64) (*Block, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getBlock(ctx, q)
}

func (db *DocstoreBlocksDatabase) GetBlockWithAccountIdAndAddress(ctx context.Context, account_id int64, host string, name string) (*Block, error) {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)
	q = q.Where("Host", "=", host)
	q = q.Where("Name", "=", name)

	return db.getBlock(ctx, q)

}

func (db *DocstoreBlocksDatabase) getBlock(ctx context.Context, q *gc_docstore.Query) (*Block, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var b Block
		err := iter.Next(ctx, &b)

		if err == io.EOF {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &b, nil
		}
	}

	return nil, ErrNotFound

}

func (db *DocstoreBlocksDatabase) AddBlock(ctx context.Context, block *Block) error {

	return db.collection.Put(ctx, block)
}

func (db *DocstoreBlocksDatabase) UpdateBlock(ctx context.Context, block *Block) error {

	return db.collection.Replace(ctx, block)
}

func (db *DocstoreBlocksDatabase) RemoveBlock(ctx context.Context, block *Block) error {

	return db.collection.Delete(ctx, block)
}

func (db *DocstoreBlocksDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
