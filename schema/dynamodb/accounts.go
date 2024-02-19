package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBAccountsTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Id"),
			KeyType:       aws.String("HASH"), // partition key
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Name"),
			AttributeType: aws.String("S"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("name"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Label"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Id"),
					KeyType:       aws.String("RANGE"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &ACCOUNTS_TABLE_NAME,
}
