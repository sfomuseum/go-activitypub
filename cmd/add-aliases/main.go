package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/aliases/add"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := add.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to add actor, %v", err)
		os.Exit(1)
	}
}
