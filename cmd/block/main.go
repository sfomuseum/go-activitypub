package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/block"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := block.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to process block", "error", err)
		os.Exit(1)
	}
}
