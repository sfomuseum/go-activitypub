package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DynamoDBDeliveriesTable = &dynamodb.CreateTableInput{
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
			AttributeName: aws.String("ActivityId"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("ActivityPubId"),
			AttributeType: "S",
		},
		{
			AttributeName: aws.String("AccountId"),
			AttributeType: "N",
		},
		{
			AttributeName: aws.String("Recipient"),
			AttributeType: "S",
		},
		{
			AttributeName: aws.String("Created"),
			AttributeType: "N",
		},
	},
	GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String("by_account"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("AccountId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_activity_id"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ActivityId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_activitypub_id"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ActivityPubId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_recipient"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("Recipient"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Created"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_activity_recipient"),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ActivityId"),
					KeyType:       "HASH",
				},
				{
					AttributeName: aws.String("Recipient"),
					KeyType:       "RANGE",
				},
			},
			Projection: &types.Projection{
				ProjectionType: "ALL",
			},
		},
		{
			IndexName: aws.String("by_created"),
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
	TableName:   &DELIVERIES_TABLE_NAME,
}
