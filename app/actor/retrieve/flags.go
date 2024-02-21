package retrieve

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var address string
var verbose bool
var insecure bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&address, "address", "", "...")
	fs.BoolVar(&verbose, "verbose", false, "...")
	fs.BoolVar(&insecure, "insecure", false, "...")
	return fs
}
