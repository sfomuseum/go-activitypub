package create

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	tables := ap_dynamodb.DynamoDBTables

	if table_prefix != "" {

		tables_prefixed := make(map[string]*dynamodb.CreateTableInput)

		for name, details := range tables {

			name_prefixed := fmt.Sprintf("%s%s", table_prefix, name)
			details.TableName = aws.String(name_prefixed)

			tables_prefixed[name_prefixed] = details
		}

		tables = tables_prefixed
	}

	logger.Info("CREATE", "refresh", refresh)
	
	create_opts := &aa_dynamodb.CreateTablesOptions{
		Tables:  tables,
		Refresh: opts.Refresh,
	}

	return aa_dynamodb.CreateTables(cl, create_opts)
}
