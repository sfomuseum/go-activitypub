package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBPostsTable = &dynamodb.CreateTableInput{
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
			AttributeName: aws.String("UUID"),
			AttributeType: aws.String("S"),
		},
	},
	GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
		{
			IndexName: aws.String("uuid"),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("UUID"),
					KeyType:       aws.String("HASH"),
				},
			},
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &POSTS_TABLE_NAME,
}
