package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub/app/post/deliver"
)

func main() {

	ctx := context.Background()
	err := deliver.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to deliver post, %v", err)
	}
}
