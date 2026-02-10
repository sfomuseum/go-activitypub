// Package database defines interfaces and default implementation for the underlying databases (or database tables) used to store ActivityPub related operations.
package database

import (
	"context"
	_ "fmt"
	"iter"
	_ "net/url"
)

// This is work in progress. Eventually things will be updated
// to use this and {FOO}_{BAR}Database to implement the {FOO}
// interface.

type Database[T any] interface {
	AddRecord(context.Context, T) error
	RemoveRecord(context.Context, T) error
	UpdateRecord(context.Context, T) error
	GetRecord(context.Context, int64) (T, error)
	ListRecords(context.Context) iter.Seq2[T, error]
	Close() error
}

func Migrate[T any](ctx context.Context, src Database[T], dst Database[T]) error {

	for v, err := range src.ListRecords(ctx) {

		if err != nil {
			return err
		}

		err = dst.AddRecord(ctx, v)

		if err != nil {
			return err
		}
	}

	return nil
}
