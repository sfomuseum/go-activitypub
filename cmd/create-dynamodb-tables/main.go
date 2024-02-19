package main

import (
	"context"
	"os"

	"github.com/sfomuseum/go-activitypub/app/dynamodb/tables/create"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := create.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to create tables", "error", err)
		os.Exit(1)
	}
}
