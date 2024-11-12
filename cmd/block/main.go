package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/block"
)

func main() {

	ctx := context.Background()
	err := block.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to block account, %v", err)
	}
}
