package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/inbox"
	_ "gocloud.dev/docstore/memdocstore"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := inbox.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to retrieve messages", "error", err)
		os.Exit(1)
	}
}
