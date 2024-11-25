package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/queue"
)

var run_opts *RunOptions

var accounts_db database.AccountsDatabase
var setupAccountsDatabaseOnce sync.Once
var setupAccountsDatabaseError error

var aliases_db database.AliasesDatabase
var setupAliasesDatabaseOnce sync.Once
var setupAliasesDatabaseError error

var followers_db database.FollowersDatabase
var setupFollowersDatabaseOnce sync.Once
var setupFollowersDatabaseError error

var following_db database.FollowingDatabase
var setupFollowingDatabaseOnce sync.Once
var setupFollowingDatabaseError error

var notes_db database.NotesDatabase
var setupNotesDatabaseOnce sync.Once
var setupNotesDatabaseError error

var messages_db database.MessagesDatabase
var setupMessagesDatabaseOnce sync.Once
var setupMessagesDatabaseError error

var blocks_db database.BlocksDatabase
var setupBlocksDatabaseOnce sync.Once
var setupBlocksDatabaseError error

var posts_db database.PostsDatabase
var setupPostsDatabaseOnce sync.Once
var setupPostsDatabaseError error

var post_tags_db database.PostTagsDatabase
var setupPostTagsDatabaseOnce sync.Once
var setupPostTagsDatabaseError error

var likes_db database.LikesDatabase
var setupLikesDatabaseOnce sync.Once
var setupLikesDatabaseError error

var boosts_db database.BoostsDatabase
var setupBoostsDatabaseOnce sync.Once
var setupBoostsDatabaseError error

var properties_db database.PropertiesDatabase
var setupPropertiesDatabaseOnce sync.Once
var setupPropertiesDatabaseError error

var process_message_queue queue.ProcessMessageQueue
var setupProcessMessageQueueOnce sync.Once
var setupProcessMessageQueueError error

var process_follower_queue queue.ProcessFollowerQueue
var setupProcessFollowerQueueOnce sync.Once
var setupProcessFollowerQueueError error
