package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBFollowingTable = &dynamodb.CreateTableInput{
	KeySchema: []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("AccountId"),
			KeyType:       aws.String("HASH"), // partition key
		},
		{
			AttributeName: aws.String("FollowingAddress"),
			KeyType:       aws.String("RANGE"), // partition key
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
	/*
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
	*/
	BillingMode: BILLING_MODE,
	TableName:   &FOLLOWING_TABLE_NAME,
}
