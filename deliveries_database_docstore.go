package activitypub

import (
	"context"
	"fmt"
	"io"

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

func (db *DocstoreDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*Delivery, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getDelivery(ctx, q)
}

func (db *DocstoreDeliveriesDatabase) GetDeliveriesWithPostIdAndRecipient(ctx context.Context, post_id int64, recipient string, deliveries_callback GetDeliveriesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("PostId", "=", post_id)
	q = q.Where("Recipient", "=", recipient)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var d Delivery
		err := iter.Next(ctx, &d)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := deliveries_callback(ctx, &d)

			if err != nil {
				return fmt.Errorf("Failed to execute deliveries callback for '%d', %w", d.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreDeliveriesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreDeliveriesDatabase) getDelivery(ctx context.Context, q *gc_docstore.Query) (*Delivery, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var d Delivery
		err := iter.Next(ctx, &d)

		if err == io.EOF {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &d, nil
		}
	}

	return nil, ErrNotFound

}
