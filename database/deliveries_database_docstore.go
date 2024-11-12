package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreDeliveriesDatabase struct {
	DeliveriesDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	err := RegisterDeliveriesDatabase(ctx, "awsdynamodb", NewDocstoreDeliveriesDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterDeliveriesDatabase(ctx, scheme, NewDocstoreDeliveriesDatabase)

		if err != nil {
			panic(err)
		}
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

func (db *DocstoreDeliveriesDatabase) GetDeliveryIdsForDateRange(ctx context.Context, start int64, end int64, cb GetDeliveryIdsCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("Created", ">=", start)
	q = q.Where("Created", "<=", end)

	iter := q.Get(ctx, "Id")
	defer iter.Stop()

	for {

		var d activitypub.Delivery
		err := iter.Next(ctx, &d)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {
			err := cb(ctx, d.Id)

			if err != nil {
				return fmt.Errorf("Failed to invoke callback for delivery %d, %w", d.Id, err)
			}
		}
	}

	return nil
}

func (db *DocstoreDeliveriesDatabase) AddDelivery(ctx context.Context, f *activitypub.Delivery) error {

	return db.collection.Put(ctx, f)
}

func (db *DocstoreDeliveriesDatabase) GetDeliveryWithId(ctx context.Context, id int64) (*activitypub.Delivery, error) {

	q := db.collection.Query()
	q = q.Where("Id", "=", id)

	return db.getDelivery(ctx, q)
}

func (db *DocstoreDeliveriesDatabase) GetDeliveries(ctx context.Context, deliveries_callback GetDeliveriesCallbackFunc) error {

	q := db.collection.Query()
	return db.getDeliveriesWithQuery(ctx, q, deliveries_callback)
}

func (db *DocstoreDeliveriesDatabase) GetDeliveriesWithActivityIdAndRecipient(ctx context.Context, activity_id int64, recipient string, deliveries_callback GetDeliveriesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("ActivityId", "=", activity_id)
	q = q.Where("Recipient", "=", recipient)

	return db.getDeliveriesWithQuery(ctx, q, deliveries_callback)
}

func (db *DocstoreDeliveriesDatabase) GetDeliveriesWithActivityPubIdAndRecipient(ctx context.Context, activity_pub_id string, recipient string, deliveries_callback GetDeliveriesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("ActivityPubId", "=", activity_pub_id)
	q = q.Where("Recipient", "=", recipient)

	return db.getDeliveriesWithQuery(ctx, q, deliveries_callback)
}

func (db *DocstoreDeliveriesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstoreDeliveriesDatabase) getDelivery(ctx context.Context, q *gc_docstore.Query) (*activitypub.Delivery, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var d activitypub.Delivery
		err := iter.Next(ctx, &d)

		if err == io.EOF {
			return nil, activitypub.ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("Failed to interate, %w", err)
		} else {
			return &d, nil
		}
	}

	return nil, activitypub.ErrNotFound

}

func (db *DocstoreDeliveriesDatabase) getDeliveriesWithQuery(ctx context.Context, q *gc_docstore.Query, deliveries_callback GetDeliveriesCallbackFunc) error {

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var d activitypub.Delivery
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
