package create

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var refresh bool

var dynamodb_client_uri string
var table_prefix string

var table_names multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.BoolVar(&refresh, "refresh", false, "...")
	fs.StringVar(&dynamodb_client_uri, "dynamodb-client-uri", "", "...")
	fs.StringVar(&table_prefix, "table-prefix", "", "...")

	fs.Var(&table_names, "table", "...")
	return fs
}
