package activitypub

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type GetBlockIdsCallbackFunc func(context.Context, int64) error
type GetBlocksCallbackFunc func(context.Context, *Block) error

type BlocksDatabase interface {
	GetBlockIdsForDateRange(context.Context, int64, int64, GetBlockIdsCallbackFunc) error
	GetBlockWithAccountIdAndAddress(context.Context, int64, string, string) (*Block, error)
	GetBlockWithId(context.Context, int64) (*Block, error)
	AddBlock(context.Context, *Block) error
	RemoveBlock(context.Context, *Block) error
	Close(context.Context) error
}

var block_database_roster roster.Roster

// BlocksDatabaseInitializationFunc is a function defined by individual block_database package and used to create
// an instance of that block_database
type BlocksDatabaseInitializationFunc func(ctx context.Context, uri string) (BlocksDatabase, error)

// RegisterBlocksDatabase registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `BlocksDatabase` instances by the `NewBlocksDatabase` method.
func RegisterBlocksDatabase(ctx context.Context, scheme string, init_func BlocksDatabaseInitializationFunc) error {

	err := ensureBlocksDatabaseRoster()

	if err != nil {
		return err
	}

	return block_database_roster.Register(ctx, scheme, init_func)
}

func ensureBlocksDatabaseRoster() error {

	if block_database_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		block_database_roster = r
	}

	return nil
}

// NewBlocksDatabase returns a new `BlocksDatabase` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `BlocksDatabaseInitializationFunc`
// function used to instantiate the new `BlocksDatabase`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterBlocksDatabase` method.
func NewBlocksDatabase(ctx context.Context, uri string) (BlocksDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := block_database_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(BlocksDatabaseInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func BlocksDatabaseSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureBlocksDatabaseRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range block_database_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
