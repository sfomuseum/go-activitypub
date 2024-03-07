package main

import (
	"context"
	"os"

	"github.com/sfomuseum/go-activitypub/app/stats/counts/date"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := date.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to generate counts for date, %v", err)
		os.Exit(1)
	}
}
