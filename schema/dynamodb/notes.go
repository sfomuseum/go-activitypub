package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBNotesTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Id"),
			KeyType:       aws.String("HASH"), // partition key
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("NoteId"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("AuthorAddress"),
			AttributeType: aws.String("S"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("note_address"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("NoteId"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("AuthorAddress"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &NOTES_TABLE_NAME,
}
