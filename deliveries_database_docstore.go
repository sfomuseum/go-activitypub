package activitypub

import (
	"context"
	"fmt"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreDeliveriesDatabase struct {
	DeliveriesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterDeliveriesDatabase(ctx, "awsdynamodb", NewDocstoreDeliveriesDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterDeliveriesDatabase(ctx, scheme, NewDocstoreDeliveriesDatabase)
	}

}

func NewDocstoreDeliveriesDatabase(ctx context.Context, uri string) (DeliveriesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreDeliveriesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstoreDeliveriesDatabase) AddDelivery(ctx context.Context, f *Delivery) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreDeliveriesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}
