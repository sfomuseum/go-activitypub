package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sfomuseum/go-activitypub/app/post"
	_ "gocloud.dev/docstore/memdocstore"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := post.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to follow actor", "error", err)
		os.Exit(1)
	}
}
