package dynamodb

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/SecondaryIndexes.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GSI.html

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ACCOUNTS_TABLE_NAME = "accounts"
var FOLLOWERS_TABLE_NAME = "followers"
var FOLLOWING_TABLE_NAME = "following"
var POSTS_TABLE_NAME = "posts"
var NOTES_TABLE_NAME = "notes"
var MESSAGES_TABLE_NAME = "messages"
var BLOCKS_TABLE_NAME = "blocks"
var DELIVERIES_TABLE_NAME = "deliveries"

var BILLING_MODE = aws.String("PAY_PER_REQUEST")

var DynamoDBTables = map[string]*dynamodb.CreateTableInput{
	ACCOUNTS_TABLE_NAME:   DynamoDBAccountsTable,
	FOLLOWERS_TABLE_NAME:  DynamoDBFollowersTable,
	FOLLOWING_TABLE_NAME:  DynamoDBFollowingTable,
	POSTS_TABLE_NAME:      DynamoDBPostsTable,
	NOTES_TABLE_NAME:      DynamoDBNotesTable,
	MESSAGES_TABLE_NAME:   DynamoDBMessagesTable,
	BLOCKS_TABLE_NAME:     DynamoDBBlocksTable,
	DELIVERIES_TABLE_NAME: DynamoDBDeliveriesTable,
}
