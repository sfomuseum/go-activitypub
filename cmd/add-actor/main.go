package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-activitypub/app/actor/add"
	_ "gocloud.dev/docstore/memdocstore"
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
