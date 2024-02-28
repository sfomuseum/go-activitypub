package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var accounts_database_uri string
var aliases_database_uri string

var account_name string
var display_name string
var blurb string
var account_url string
var account_type string

var account_id int64

var public_key_uri string
var private_key_uri string

var account_icon_uri string
var allow_remote_icon_uri bool
var embed_icon_uri bool

var discoverable bool

var aliases multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "...")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "", "...")

	fs.Int64Var(&account_id, "account-id", 0, "...")

	fs.StringVar(&account_name, "account-name", "", "...")
	fs.Var(&aliases, "alias", "...")

	fs.StringVar(&display_name, "display-name", "", "...")
	fs.StringVar(&blurb, "blurb", "", "...")
	fs.StringVar(&account_url, "url", "", "...")
	fs.StringVar(&account_type, "account-type", "Person", "...")

	fs.BoolVar(&discoverable, "discoverable", true, "...")

	fs.StringVar(&public_key_uri, "public-key-uri", "", "...")
	fs.StringVar(&private_key_uri, "private-key-uri", "", "...")

	fs.StringVar(&account_icon_uri, "account-icon-uri", "", "...")
	fs.BoolVar(&allow_remote_icon_uri, "allow-remote-icon-uri", false, "...")
	fs.BoolVar(&embed_icon_uri, "embed-icon-uri", false, "...")
	return fs
}
