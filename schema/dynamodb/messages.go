package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoDBMessagesTable = &dynamodb.CreateTableInput{
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
			AttributeName: aws.String("AuthorAddress"),
			AttributeType: "S",
		},
		{
			AttributeName: aws.String("NoteId"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: "N",
		},
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("author"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AuthorAddress"),
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
			IndexName: aws.String("account_note"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("NoteId"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("account_author"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("AuthorAddress"),
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
	TableName:   &MESSAGES_TABLE_NAME,
}
