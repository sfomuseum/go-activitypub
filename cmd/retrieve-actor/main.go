package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/actor/retrieve"
)

func main() {

	ctx := context.Background()
	err := retrieve.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to retrieve actor, %v", err)
	}
}
