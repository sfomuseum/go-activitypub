package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/boost/note"
)

func main() {

	ctx := context.Background()
	err := note.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to create post, %v", err)
	}
}
