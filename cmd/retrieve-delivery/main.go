package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/deliveries/retrieve"
)

func main() {

	ctx := context.Background()
	err := retrieve.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to retrieve delivery, %v", err)
	}
}
