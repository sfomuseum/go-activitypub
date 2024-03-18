package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/stats/followers/list"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := list.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to list followers", "error", err)
		os.Exit(1)
	}
}
