package main

import (
	"context"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/account/get"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := get.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to get account, %v", err)
		os.Exit(1)
	}
}
