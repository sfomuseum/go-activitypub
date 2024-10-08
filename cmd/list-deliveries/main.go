package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/post/deliveries"
)

func main() {

	ctx := context.Background()
	err := deliveries.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to list deliveries, %v", err)
	}
}
