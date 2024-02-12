package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var actor_database_uri string
var actor_id string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&actor_database_uri, "actor-database-uri", "mem://actors/Id", "...")
	fs.StringVar(&actor_id, "actor-id", "", "...")

	return fs
}
