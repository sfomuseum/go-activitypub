package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/post/create"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := create.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to post message", "error", err)
		os.Exit(1)
	}
}
