package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBMessagesTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       aws.String("HASH"),
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("AuthorAddress"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("NoteId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: aws.String("N"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("author"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("AuthorAddress"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
		{
			IndexName: aws.String("note"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("NoteId"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &MESSAGES_TABLE_NAME,
}
