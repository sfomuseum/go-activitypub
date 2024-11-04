package retrieve

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var address string
var verbose bool
var insecure bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("activitypub")

	fs.StringVar(&address, "address", "", "The @user@host address of the actor to retrieve.")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.BoolVar(&insecure, "insecure", false, "A boolean flag indicating whether the host that the -address flag resolves to is running without TLS enabled.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Retrieve an ActivityPub actor by its @user@host address and emit it as a JSON-encoded string..\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
