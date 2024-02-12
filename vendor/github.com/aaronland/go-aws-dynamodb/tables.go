package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
)

// CreateTablesOptions defines options for the CreateTables method
type CreateTablesOptions struct {
	// A hash map containing table names and their dynamodb.CreateTableInput defintions
	Tables map[string]*aws_dynamodb.CreateTableInput
	// If true and the table already exists, delete and recreate the table
	Refresh bool
}

// Create one or more tables associated with the dynamodb.DynamoDB instance.
func CreateTables(client *aws_dynamodb.DynamoDB, opts *CreateTablesOptions) error {

	for table_name, def := range opts.Tables {

		has_table, err := HasTable(client, table_name)

		if err != nil {
			return fmt.Errorf("Failed to determined whether table exists, %w", err)
		}

		if has_table {

			if !opts.Refresh {
				continue
			}

			req := &aws_dynamodb.DeleteTableInput{
				TableName: aws.String(table_name),
			}

			client.DeleteTable(req)
		}

		def.TableName = aws.String(table_name)

		_, err = client.CreateTable(def)

		if err != nil {
			return fmt.Errorf("Failed to create table '%s', %w", table_name, err)
		}
	}

	return nil
}

// Return a boolean value indication whether or not the dynamodb.DynamoDB instances contains a table matching table_name.
func HasTable(client *aws_dynamodb.DynamoDB, table_name string) (bool, error) {

	tables, err := ListTables(client)

	if err != nil {
		return false, err
	}

	has_table := false

	for _, name := range tables {

		if name == table_name {
			has_table = true
			break
		}
	}

	return has_table, nil
}

// Return the list of table names associated with the dynamodb.DynamoDB instance.
func ListTables(client *aws_dynamodb.DynamoDB) ([]string, error) {

	tables := make([]string, 0)

	input := &aws_dynamodb.ListTablesInput{}

	for {

		rsp, err := client.ListTables(input)

		if err != nil {
			return nil, err
		}

		for _, n := range rsp.TableNames {
			tables = append(tables, *n)
		}

		input.ExclusiveStartTableName = rsp.LastEvaluatedTableName

		if rsp.LastEvaluatedTableName == nil {
			break
		}

	}

	return tables, nil
}
