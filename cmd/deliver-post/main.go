package main

import (
	"context"
	"os"

	"github.com/sfomuseum/go-activitypub/app/post/deliver"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := deliver.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to post message", "error", err)
		os.Exit(1)
	}
}
