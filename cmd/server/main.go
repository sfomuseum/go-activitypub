package main

import (
	"context"
	"log"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/sfomuseum/go-pubsub/publisher"
	_ "gocloud.dev/blob/fileblob"

	"github.com/sfomuseum/go-activitypub/app/server"
)

func main() {

	ctx := context.Background()
	err := server.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run server, %v", err)
	}
}
