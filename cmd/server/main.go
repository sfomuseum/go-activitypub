package main

import (
	"context"
	"os"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "github.com/sfomuseum/go-pubsub/publisher"
	_ "gocloud.dev/blob/fileblob"
	
	"github.com/sfomuseum/go-activitypub/app/server"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to add actor, %v", err)
		os.Exit(1)
	}
}
