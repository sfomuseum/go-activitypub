package dynamodb

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	_ "github.com/aws/aws-sdk-go-v2/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// CreateTablesOptions defines options for the CreateTables method
type CreateTablesOptions struct {
	// A hash map containing table names and their dynamodb.CreateTableInput defintions
	Tables map[string]*aws_dynamodb.CreateTableInput
	// If true and the table already exists, delete and recreate the table
	Refresh bool
	// An optional string to append to each table name as it is created.
	Prefix string
}

// Create one or more tables associated with the dynamodb.DynamoDB instance.
func CreateTables(ctx context.Context, client *aws_dynamodb.Client, opts *CreateTablesOptions) error {

	for _, def := range opts.Tables {

		table_name := *def.TableName

		logger := slog.Default()
		logger = logger.With("table", table_name)

		if opts.Prefix != "" {

			logger.Debug("Assign prefix to table name", "prefix", opts.Prefix)
			table_name = opts.Prefix + table_name

			logger = slog.Default()
			logger = logger.With("table", table_name)
		}

		// To do: Do this concurrently because of the delay waiting for table deletion to complete

		logger.Debug("Check whether table exists")

		has_table, err := HasTable(ctx, client, table_name)

		if err != nil {
			logger.Error("Failed to determine if table exists", "error", err)
			return fmt.Errorf("Failed to determined whether table exists, %w", err)
		}

		if has_table {

			if !opts.Refresh {
				logger.Debug("Table exists, refresh disabled")
				continue
			}

			logger.Debug("Table exists, refresh")

			req := &aws_dynamodb.DeleteTableInput{
				TableName: def.TableName,
			}

			_, err := client.DeleteTable(ctx, req)

			if err != nil {
				logger.Error("Failed to delete table", "error", err)
				return fmt.Errorf("Failed to delete table '%s', %w", table_name, err)
			}

			// Now wait for the deletion to complete...
			logger.Debug("Table deleted, now waiting for completion")

			ctx := context.Background()

			ticker_ctx, ticker_cancel := context.WithTimeout(ctx, 30*time.Second)
			defer ticker_cancel()

			ticker := time.NewTicker(5 * time.Second)

			done_ch := make(chan bool)
			ready_ch := make(chan bool)
			err_ch := make(chan error)

			go func() {
				for {
					select {
					case <-ticker_ctx.Done():
						return
					case <-done_ch:
						return
					case <-ticker.C:

						has_table, err := HasTable(ctx, client, table_name)

						if err != nil {
							logger.Error("Failed to determine if table exists", "error", err)
						} else {
							logger.Debug("Has table check", "has_table", has_table)
						}

						if !has_table {
							ready_ch <- true
						}

					}
				}
			}()

			table_ready := false

			for {
				select {
				case <-ticker_ctx.Done():
					return fmt.Errorf("Ticker to delete table timed out")
				case err := <-err_ch:
					return fmt.Errorf("Failed to delete table, %w", err)
				case <-ready_ch:
					table_ready = true
				}

				if table_ready {
					break
				}
			}

		}

		logger.Debug("Create table")

		_, err = client.CreateTable(ctx, def)

		if err != nil {
			logger.Error("Failed to create table", "error", err)
			return fmt.Errorf("Failed to create table '%s', %w", table_name, err)
		}
	}

	return nil
}

// Return a boolean value indication whether or not the dynamodb.DynamoDB instances contains a table matching table_name.
func HasTable(ctx context.Context, client *aws_dynamodb.Client, table_name string) (bool, error) {

	tables, err := ListTables(ctx, client)

	if err != nil {
		return false, err
	}

	has_table := slices.Contains(tables, table_name)

	return has_table, nil
}

// Return the list of table names associated with the dynamodb.DynamoDB instance.
func ListTables(ctx context.Context, client *aws_dynamodb.Client) ([]string, error) {

	tables := make([]string, 0)

	input := &aws_dynamodb.ListTablesInput{}

	for {

		rsp, err := client.ListTables(ctx, input)

		if err != nil {
			return nil, err
		}

		for _, n := range rsp.TableNames {
			tables = append(tables, n)
		}

		input.ExclusiveStartTableName = rsp.LastEvaluatedTableName

		if rsp.LastEvaluatedTableName == nil {
			break
		}

	}

	return tables, nil
}
