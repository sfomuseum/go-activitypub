package database

import (
	"context"
	"fmt"
	"io"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	"github.com/sfomuseum/go-activitypub"
	gc_docstore "gocloud.dev/docstore"
)

type DocstorePropertiesDatabase struct {
	PropertiesDatabase
	collection *gc_docstore.Collection
}

func init() {
	ctx := context.Background()

	err := RegisterPropertiesDatabase(ctx, "awsdynamodb", NewDocstorePropertiesDatabase)

	if err != nil {
		panic(err)
	}

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		err := RegisterPropertiesDatabase(ctx, scheme, NewDocstorePropertiesDatabase)

		if err != nil {
			panic(err)
		}
	}
}

func NewDocstorePropertiesDatabase(ctx context.Context, uri string) (PropertiesDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstorePropertiesDatabase{
		collection: col,
	}

	return db, nil
}

func (db *DocstorePropertiesDatabase) GetProperties(ctx context.Context, cb GetPropertiesCallbackFunc) error {

	q := db.collection.Query()

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var p activitypub.Property
		err := iter.Next(ctx, &p)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &p)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for property %v, %w", p, err)
			}
		}
	}

	return nil
}

func (db *DocstorePropertiesDatabase) GetPropertiesForAccount(ctx context.Context, account_id int64, cb GetPropertiesCallbackFunc) error {

	q := db.collection.Query()
	q = q.Where("AccountId", "=", account_id)

	iter := q.Get(ctx)
	defer iter.Stop()

	for {

		var p activitypub.Property
		err := iter.Next(ctx, &p)

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Failed to interate, %w", err)
		} else {

			err := cb(ctx, &p)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for property %v, %w", p, err)
			}
		}
	}

	return nil
}

func (db *DocstorePropertiesDatabase) AddProperty(ctx context.Context, property *activitypub.Property) error {
	return db.collection.Put(ctx, property)
}

func (db *DocstorePropertiesDatabase) UpdateProperty(ctx context.Context, property *activitypub.Property) error {
	return db.collection.Replace(ctx, property)
}

func (db *DocstorePropertiesDatabase) RemoveProperty(ctx context.Context, property *activitypub.Property) error {
	return db.collection.Delete(ctx, property)
}

func (db *DocstorePropertiesDatabase) Close(ctx context.Context) error {
	return db.collection.Close()
}

func (db *DocstorePropertiesDatabase) getProperty(ctx context.Context, q *gc_docstore.Query) (*activitypub.Property, error) {

	iter := q.Get(ctx)
	defer iter.Stop()

	var a activitypub.Property
	err := iter.Next(ctx, &a)

	if err == io.EOF {
		return nil, activitypub.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Failed to interate, %w", err)
	} else {
		return &a, nil
	}
}
