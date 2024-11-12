package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub/app/message/process/example"
)

func main() {

	ctx := context.Background()
	err := example.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to process message, %v", err)
	}
}
