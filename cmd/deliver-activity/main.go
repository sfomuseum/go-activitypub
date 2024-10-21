package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/activity/deliver"
)

func main() {

	ctx := context.Background()
	err := deliver.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to deliver activity, %v", err)
	}
}
