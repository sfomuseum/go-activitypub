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

	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "A valid sfomuseum/go-activitypub/database.NotesDatabase URI.")
	fs.Int64Var(&note_id, "note-id", 0, "The unique 64-bit note ID to retrieve.")
	fs.BoolVar(&body, "body", false, "Display the (ActivityPub) body of the note.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
