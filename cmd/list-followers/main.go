package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/stats/followers/list"
)

func main() {

	ctx := context.Background()
	err := list.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to list followers, %v", err)
	}
}
