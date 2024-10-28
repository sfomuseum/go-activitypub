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

	fs.StringVar(&dynamodb_client_uri, "dynamodb-client-uri", "", "A valid aaronland/")
	fs.StringVar(&table_prefix, "table-prefix", "", "A optional prefix to assign to each table name.")
	fs.BoolVar(&refresh, "refresh", false, "Refresh tables if already present.")

	fs.Var(&table_names, "table", "Zero or more table names to create. If zero then all the default tables will be created.")
	return fs
}
