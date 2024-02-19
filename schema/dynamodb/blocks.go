package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBBlocksTable = &dynamodb.CreateTableInput{
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
			AttributeName: aws.String("Label"),
			AttributeType: aws.String("S"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("label"),
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
	TableName:   &BLOCKS_TABLE_NAME,
}
