package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/notes/retrieve"
)

func main() {

	ctx := context.Background()
	err := retrieve.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to note actor, %v", err)
	}
}
