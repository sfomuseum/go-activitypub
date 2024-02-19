package create

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	ap_dynamodb "github.com/sfomuseum/go-activitypub/schema/dynamodb"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	slog.SetDefault(logger)

	cl, err := aa_dynamodb.NewClientWithURI(ctx, opts.DynamodbClientURI)

	if err != nil {
		return fmt.Errorf("Failed to create dynamodb client, %w", err)
	}

	create_opts := &aa_dynamodb.CreateTablesOptions{
		Tables:  ap_dynamodb.DynamoDBTables,
		Refresh: opts.Refresh,
	}

	return aa_dynamodb.CreateTables(cl, create_opts)
}
