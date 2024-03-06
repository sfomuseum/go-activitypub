package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBAliasesTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Name"), // partition key
			KeyType:       aws.String("HASH"),
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Name"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: aws.String("N"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("account"),
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
			IndexName: aws.String("by_created"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Created"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("KEYS_ONLY"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &ALIASES_TABLE_NAME,
}
