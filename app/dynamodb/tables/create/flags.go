package create

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var refresh bool

var dynamodb_client_uri string
var table_prefix string

var table_names multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.StringVar(&dynamodb_client_uri, "dynamodb-client-uri", "", "A valid aaronland/gocloud-docstore URI (dynamodb:// or awsdynamodb://).")
	fs.StringVar(&table_prefix, "table-prefix", "", "A optional prefix to assign to each table name.")
	fs.BoolVar(&refresh, "refresh", false, "Refresh tables if already present.")

	fs.Var(&table_names, "table", "Zero or more table names to create. If zero then all the default tables will be created.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Create (or refresh) the DynamoDB tables necessary for use with the go-activitypub package.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
