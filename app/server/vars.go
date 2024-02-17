package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub"
)

var run_opts *RunOptions

var accounts_db activitypub.AccountsDatabase
var setupAccountsDatabaseOnce sync.Once
var setupAccountsDatabaseError error

var followers_db activitypub.FollowersDatabase
var setupFollowersDatabaseOnce sync.Once
var setupFollowersDatabaseError error

var following_db activitypub.FollowingDatabase
var setupFollowingDatabaseOnce sync.Once
var setupFollowingDatabaseError error
