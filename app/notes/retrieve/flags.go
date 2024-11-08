package retrieve

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var notes_database_uri string
var note_id int64
var body bool
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("notes")

	fs.StringVar(&notes_database_uri, "notes-database-uri", "", "A registered sfomuseum/go-activitypub/database.NotesDatabase URI.")
	fs.Int64Var(&note_id, "note-id", 0, "The unique 64-bit note ID to retrieve.")
	fs.BoolVar(&body, "body", false, "Display the (ActivityPub) body of the note.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Retrieve a message (note) by its unique go-activitypub ID.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
