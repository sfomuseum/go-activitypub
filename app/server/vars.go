package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub"
)

var run_opts *RunOptions

var accounts_db activitypub.AccountsDatabase
var setupAccountsDatabaseOnce sync.Once
var setupAccountsDatabaseError error

var aliases_db activitypub.AliasesDatabase
var setupAliasesDatabaseOnce sync.Once
var setupAliasesDatabaseError error

var followers_db activitypub.FollowersDatabase
var setupFollowersDatabaseOnce sync.Once
var setupFollowersDatabaseError error

var following_db activitypub.FollowingDatabase
var setupFollowingDatabaseOnce sync.Once
var setupFollowingDatabaseError error

var notes_db activitypub.NotesDatabase
var setupNotesDatabaseOnce sync.Once
var setupNotesDatabaseError error

var messages_db activitypub.MessagesDatabase
var setupMessagesDatabaseOnce sync.Once
var setupMessagesDatabaseError error

var blocks_db activitypub.BlocksDatabase
var setupBlocksDatabaseOnce sync.Once
var setupBlocksDatabaseError error

var posts_db activitypub.PostsDatabase
var setupPostsDatabaseOnce sync.Once
var setupPostsDatabaseError error

var likes_db activitypub.LikesDatabase
var setupLikesDatabaseOnce sync.Once
var setupLikesDatabaseError error

var boosts_db activitypub.BoostsDatabase
var setupBoostsDatabaseOnce sync.Once
var setupBoostsDatabaseError error

var properties_db activitypub.PropertiesDatabase
var setupPropertiesDatabaseOnce sync.Once
var setupPropertiesDatabaseError error

var process_message_queue activitypub.ProcessMessageQueue
var setupProcessMessageQueueOnce sync.Once
var setupProcessMessageQueueError error
