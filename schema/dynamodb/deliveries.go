package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBDeliveriesTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       aws.String("HASH"),
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("PostId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Recipient"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: aws.String("N"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("by_account"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
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
			IndexName: aws.String("by_post"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("PostId"),
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
			IndexName: aws.String("by_recipient"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Recipient"),
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
	},
	BillingMode: BILLING_MODE,
	TableName:   &DELIVERIES_TABLE_NAME,
}
