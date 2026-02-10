// Package database defines interfaces and default implementation for the underlying databases (or database tables) used to store ActivityPub related operations.
package database

import (
	"context"
	_ "fmt"
	"iter"
	_ "net/url"
)

// @sql://boosts:@foo...

type Database[T any] interface {
	Iterate(context.Context) iter.Seq2[T, error]
	Add(context.Context, T) error
}

func Migrate[T any](ctx context.Context, src Database[T], dst Database[T]) error {

	for v, err := range src.Iterate(ctx) {

		if err != nil {
			return err
		}

		err = dst.Add(ctx, v)

		if err != nil {
			return err
		}
	}

	return nil
}
