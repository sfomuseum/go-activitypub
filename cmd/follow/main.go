package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/follow"
)

func main() {

	ctx := context.Background()
	err := follow.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to follow actor, %v", err)
	}
}
