package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/inbox"
)

func main() {

	ctx := context.Background()
	err := inbox.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to retrieve inbox, %v", err)
	}
}
