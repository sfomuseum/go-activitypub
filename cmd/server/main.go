package main

import (
	"context"
	"os"

	_ "github.com/aaronland/gocloud-blob-s3"
	"github.com/sfomuseum/go-activitypub/app/server"
	"github.com/sfomuseum/go-activitypub/slog"
	_ "gocloud.dev/blob/fileblob"
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
