package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub/app/activity/list"
)

func main() {

	ctx := context.Background()
	err := list.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to list activities, %v", err)
	}
}
