package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/deliveries/retrieve"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := retrieve.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to retrieve deliveries, %v", err)
		os.Exit(1)
	}
}
