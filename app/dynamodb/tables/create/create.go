package create

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	aa_dynamodb "github.com/aaronland/go-aws-dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ap_dynamodb "github.com/sfomuseum/go-activitypub/schema/dynamodb"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	cl, err := aa_dynamodb.NewClient(ctx, opts.DynamodbClientURI)

	if err != nil {
		return fmt.Errorf("Failed to create dynamodb client, %w", err)
	}

	tables := ap_dynamodb.DynamoDBTables

	if len(table_names) > 0 {

		tables_named := make(map[string]*dynamodb.CreateTableInput)

		for _, n := range table_names {

			details, exists := tables[n]

			if !exists {
				return fmt.Errorf("Missing table definition for '%s'", n)
			}

			tables_named[n] = details
		}

		tables = tables_named
	}

	if table_prefix != "" {

		tables_prefixed := make(map[string]*dynamodb.CreateTableInput)

		for name, details := range tables {

			name_prefixed := fmt.Sprintf("%s%s", table_prefix, name)
			details.TableName = aws.String(name_prefixed)

			tables_prefixed[name_prefixed] = details
		}

		tables = tables_prefixed
	}

	create_opts := &aa_dynamodb.CreateTablesOptions{
		Tables:  tables,
		Refresh: opts.Refresh,
	}

	err = aa_dynamodb.CreateTables(ctx, cl, create_opts)

	if err != nil {
		logger.Error("Failed to create tables", "error", err)
	}

	return nil
}
