package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var actor_id string
var database_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&database_uri, "database-uri", "mem://actors/Id", "...")
	fs.StringVar(&actor_id, "actor-id", "", "...")

	return fs
}
