package block

import (
	"flag"

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

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&blocks_database_uri, "blocks-database-uri", "", "...")

	fs.StringVar(&account_name, "account-name", "", "...")

	fs.StringVar(&block_name, "block-name", "*", "...")
	fs.StringVar(&block_host, "block-host", "", "...")

	fs.BoolVar(&undo, "undo", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	return fs
}
