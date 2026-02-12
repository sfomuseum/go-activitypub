// Package database defines interfaces and default implementation for the underlying databases (or database tables) used to store ActivityPub related operations.
package database

import (
	"context"
	"fmt"
	"iter"
	"sync/atomic"
)

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type Where struct {
	Conditions []*Condition `json:"conditions"`
	Relation   string       `json:"relation"`
}

type OrderBy struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type Query struct {
	Where   *Where   `json:"where"`
	Offset  *int     `json:"offset"`
	Limit   *int     `json:"limit"`
	OrderBy *OrderBy `json:"order_by"`
}

type Database[T any] interface {
	AddRecord(context.Context, T) error
	RemoveRecord(context.Context, T) error
	UpdateRecord(context.Context, T) error
	GetRecord(context.Context, int64) (T, error)
	QueryRecords(context.Context, *Query) iter.Seq2[T, error]
	Close() error
}

func ListRecords[T any](ctx context.Context, db Database[T]) error {

	for rec, err := range db.QueryRecords(ctx, nil) {

		if err != nil {
			return err
		}

		fmt.Printf("%s\n", rec)
	}

	return nil
}

func Migrate[T any](ctx context.Context, src Database[T], dst Database[T]) (int64, int64, int64, error) {

	count := int64(0)
	success := int64(0)
	errors := int64(0)

	for v, err := range src.QueryRecords(ctx, nil) {

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
