package retrieve

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var notes_database_uri string
var note_id int64
var body bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("notes")

	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "...")
	fs.Int64Var(&note_id, "note-id", 0, "...")
	fs.BoolVar(&body, "body", false, "...")
	fs.BoolVar(&verbose, "verbose", false, "...")

	return fs
}
