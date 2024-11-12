package add

import (
	"flag"
	"fmt"
	"os"

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

	fs.StringVar(&accounts_database_uri, "accounts-database-uri", "", "A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.")
	fs.StringVar(&aliases_database_uri, "aliases-database-uri", "", "A registered sfomuseum/go-activitypub/database.AliasesDatabase URI.")
	fs.StringVar(&properties_database_uri, "properties-database-uri", "", "A registered sfomuseum/go-activitypub/database.PropertiesDatabase URI.")

	fs.Int64Var(&account_id, "account-id", 0, "An optional unique identifier to assign to the account being created. If 0 then an ID will be generated automatically.")

	fs.StringVar(&account_name, "account-name", "", "The user (preferred) name for the account being created.")
	fs.Var(&aliases_list, "alias", "Zero or more aliases for the account being created.")

	fs.StringVar(&display_name, "display-name", "", "The display name for the account being created.")
	fs.StringVar(&blurb, "blurb", "", "The descriptive blurb (caption) for the account being created.")
	fs.StringVar(&account_url, "url", "", "The URL for the account being created.")
	fs.StringVar(&account_type, "account-type", "Person", "The type of account being created. Valid options are: Person, Service.")

	fs.BoolVar(&discoverable, "discoverable", true, "Boolean flag indicating whether the account should be discoverable.")

	fs.StringVar(&public_key_uri, "public-key-uri", "", "A valid `gocloud.dev/runtimevar` referencing the PEM-encoded public key for the account.")
	fs.StringVar(&private_key_uri, "private-key-uri", "", "A valid `gocloud.dev/runtimevar` referencing the PEM-encoded private key for the account.")

	fs.StringVar(&account_icon_uri, "account-icon-uri", "", "A valid `gocloud.dev/blob` URI (as in the bucket URI + filename) referencing the icon URI for the account.")
	fs.BoolVar(&allow_remote_icon_uri, "allow-remote-icon-uri", false, "Allow the -account-icon-uri flag to specify a remote URI.")
	fs.BoolVar(&embed_icon_uri, "embed-icon-uri", false, "If true then assume the -account-icon-uri flag references a local file and read its body in to a base64-encoded value to be stored with the account record.")

	fs.Var(&properties_kv, "property", "Zero or more {KEY}={VALUE} properties to be assigned to the new account.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Add a new ActivityPub account.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
