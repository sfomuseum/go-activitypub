package add

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var accounts_database_uri string
var aliases_database_uri string
var properties_database_uri string

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

var aliases_list multi.MultiString
var properties_kv multi.KeyValueString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A valid sfomuseum/go-activitypub/database.AccountsDatabase URI.")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "", "A valid sfomuseum/go-activitypub/database.AliasesDatabase URI.")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "A valid sfomuseum/go-activitypub/database.PropertiesDatabase URI.")

	fs.Int64Var(&account_id, "account-id", 0, "An optional unique identifier to assign to the account being created. If 0 then an ID will be generated automatically.")

	fs.StringVar(&account_name, "account-name", "", "The user (preferred) name for the account being created.")
	fs.Var(&aliases_list, "alias", "Zero or more aliases for the account being created.")

	fs.StringVar(&display_name, "display-name", "", "The display name for the account being created.")
	fs.StringVar(&blurb, "blurb", "", "The descriptive blurb (caption) for the account being created.")
	fs.StringVar(&account_url, "url", "", "The URL for the account being created.")
	fs.StringVar(&account_type, "account-type", "Person", "The type of account being created. Valid options are: Person, Service.")

	fs.BoolVar(&discoverable, "discoverable", true, "Boolean flag indicating whether the account should be discoverable.")

	fs.StringVar(&public_key_uri, "public-key-uri", "", "...")
	fs.StringVar(&private_key_uri, "private-key-uri", "", "...")

	fs.StringVar(&account_icon_uri, "account-icon-uri", "", "...")
	fs.BoolVar(&allow_remote_icon_uri, "allow-remote-icon-uri", false, "...")
	fs.BoolVar(&embed_icon_uri, "embed-icon-uri", false, "...")

	fs.Var(&properties_kv, "property", "Zero or more {KEY}={VALUE} properties to be assigned to the new account.")
	return fs
}
