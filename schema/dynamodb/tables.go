package dynamodb

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/SecondaryIndexes.html
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GSI.html

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ACCOUNTS_TABLE_NAME = "accounts"
var ACTIVITIES_TABLE_NAME = "activities"
var ALIASES_TABLE_NAME = "aliases"
var PROPERTIES_TABLE_NAME = "properties"
var FOLLOWERS_TABLE_NAME = "followers"
var FOLLOWING_TABLE_NAME = "following"
var POSTS_TABLE_NAME = "posts"
var POST_TAGS_TABLE_NAME = "post_tags"
var NOTES_TABLE_NAME = "notes"
var MESSAGES_TABLE_NAME = "messages"
var BLOCKS_TABLE_NAME = "blocks"
var DELIVERIES_TABLE_NAME = "deliveries"
var LIKES_TABLE_NAME = "likes"
var BOOSTS_TABLE_NAME = "boosts"

var BILLING_MODE = types.BillingModePayPerRequest

var DynamoDBTables = map[string]*dynamodb.CreateTableInput{
	ACCOUNTS_TABLE_NAME:   DynamoDBAccountsTable,
	ACTIVITIES_TABLE_NAME: DynamoDBActivitiesTable,
	ALIASES_TABLE_NAME:    DynamoDBAliasesTable,
	FOLLOWERS_TABLE_NAME:  DynamoDBFollowersTable,
	FOLLOWING_TABLE_NAME:  DynamoDBFollowingTable,
	POSTS_TABLE_NAME:      DynamoDBPostsTable,
	POST_TAGS_TABLE_NAME:  DynamoDBPostTagsTable,
	NOTES_TABLE_NAME:      DynamoDBNotesTable,
	MESSAGES_TABLE_NAME:   DynamoDBMessagesTable,
	BLOCKS_TABLE_NAME:     DynamoDBBlocksTable,
	DELIVERIES_TABLE_NAME: DynamoDBDeliveriesTable,
	LIKES_TABLE_NAME:      DynamoDBLikesTable,
	BOOSTS_TABLE_NAME:     DynamoDBBoostsTable,
	PROPERTIES_TABLE_NAME: DynamoDBPropertiesTable,
}
