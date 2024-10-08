package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub/app/stats/counts/date"
)

func main() {

	ctx := context.Background()
	err := date.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to derive counts, %v", err)
	}
}
