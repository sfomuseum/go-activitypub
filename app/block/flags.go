package block

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var accounts_database_uri string
var blocks_database_uri string

var account_name string

var block_name string
var block_host string

var undo bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("block")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A known sfomuseum/go-activitypub/AccountsDatabase URI.")
	fs.StringVar(&blocks_database_uri, "blocks-database-uri", "", "A known sfomuseum/go-activitypub/BlocksDatabase URI.")

	fs.StringVar(&account_name, "account-name", "", "The name of the account doing the blocking.")

	fs.StringVar(&block_name, "block-name", "*", "The name of the account being blocked. If \"*\" then all the accounts associated with the blocked host will be blocked.")
	fs.StringVar(&block_host, "block-host", "", "The name of the host associated with the account being blocked.")

	fs.BoolVar(&undo, "undo", false, "Undo an existing block.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Manage the blocking of third-parties on behalf of a registered sfomuseum/go-activity account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
