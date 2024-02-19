package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBBlocksTable = &dynamodb.CreateTableInput{
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
			AttributeName: aws.String("AccountId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Host"),
			AttributeType: aws.String("S"),
		},
		/*
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
		*/
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("account_and_address"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Host"),
					KeyType:       aws.String("RANGE"),
				},
				/*
					{
						AttributeName: aws.String("Name"),
						KeyType:       aws.String("HASH"),
					},
				*/
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &BLOCKS_TABLE_NAME,
}
