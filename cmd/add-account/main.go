package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sfomuseum/go-activitypub/app/account/add"
)

func main() {

	ctx := context.Background()
	err := add.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to add account, %v", err)
	}
}
