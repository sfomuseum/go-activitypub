package activitypub

import (
	"context"
	"fmt"

	aa_docstore "github.com/aaronland/gocloud-docstore"
	gc_docstore "gocloud.dev/docstore"
)

type DocstoreActorDatabase struct {
	ActorDatabase
	collection *gc_docstore.Collection
}

func init() {

	ctx := context.Background()

	RegisterActorDatabase(ctx, "awsdynamodb", NewDocstoreActorDatabase)

	for _, scheme := range gc_docstore.DefaultURLMux().CollectionSchemes() {
		RegisterActorDatabase(ctx, scheme, NewDocstoreActorDatabase)
	}
}

func NewDocstoreActorDatabase(ctx context.Context, uri string) (ActorDatabase, error) {

	col, err := aa_docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open collection, %w", err)
	}

	db := &DocstoreActorDatabase{
		collection: col,
	}

	return db, nil
}
