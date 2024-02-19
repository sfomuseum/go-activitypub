package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var refresh bool

var dynamodb_client_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.BoolVar(&refresh, "refresh", false, "...")
	fs.StringVar(&dynamodb_client_uri, "dynamodb-client-uri", "", "...")

	return fs
}
