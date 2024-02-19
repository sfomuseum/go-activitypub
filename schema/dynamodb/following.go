package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBFollowingTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("AccountId"), // partition key
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("FollowingAddress"),
			KeyType:       aws.String("RANGE"),
		},
	},
	AttributeDefinitions: []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("FollowingAddress"),
			AttributeType: aws.String("S"),
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &FOLLOWING_TABLE_NAME,
}
