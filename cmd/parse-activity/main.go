package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/go-activitypub/ap"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		var a *ap.Activity

		dec := json.NewDecoder(r)
		err = dec.Decode(&a)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		log.Printf("Decoded %s as %s (%s)\n", path, a.Type, a.Id)
	}

}
