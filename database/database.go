// Package database defines interfaces and default implementation for the underlying databases (or database tables) used to store ActivityPub related operations.
package database

import (
	"context"
	_ "fmt"
	"iter"
	_ "net/url"
	"sync/atomic"
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

func Migrate[T any](ctx context.Context, src Database[T], dst Database[T]) (int64, int64, int64, error) {

	count := int64(0)
	success := int64(0)
	errors := int64(0)

	for v, err := range src.ListRecords(ctx) {

		defer atomic.AddInt64(&count, 1)

		if err != nil {
			atomic.AddInt64(&errors, 1)
			return count, success, errors, err
		}

		err = dst.AddRecord(ctx, v)

		if err != nil {
			atomic.AddInt64(&errors, 1)
			return count, success, errors, err
		}

		atomic.AddInt64(&success, 1)
	}

	return count, success, errors, nil
}
