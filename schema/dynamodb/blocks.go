package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoDBBlocksTable = &dynamodb.CreateTableInput{
	KeySchema: []types.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       "HASH",
		},
	},
	AttributeDefinitions: []types.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("Host"),
			AttributeType: "S",
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: "N",
		},
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("account_and_host"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Host"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("host"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Host"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_created"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Created"),
					KeyType:       "HASH",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "KEYS_ONLY",
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &BLOCKS_TABLE_NAME,
}
