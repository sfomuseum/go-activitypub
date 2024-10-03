package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoDBFollowersTable = &dynamodb.CreateTableInput{
	KeySchema: []types.KeySchemaElement{
		{
			AttributeName: aws.String("Id"), // partition key
			KeyType:       "HASH",
		},
	},
	AttributeDefinitions: []types.AttributeDefinition{
		{
			AttributeName: aws.String("Id"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("FollowerAddress"),
			AttributeType: "S",
		},
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("account_follower"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("FollowerAddress"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},

		{
			IndexName: aws.String("created"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Created"),
					KeyType:       "HASH",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "KEYS_ONLY",
			},
		},
	},
	BillingMode: BILLING_MODE,
	TableName:   &FOLLOWERS_TABLE_NAME,
}
