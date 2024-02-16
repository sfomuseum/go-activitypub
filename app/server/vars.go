package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub"
)

var run_opts *RunOptions

var accounts_db activitypub.AccountsDatabase

var setupAccountsDatabaseOnce sync.Once
var setupAccountsDatabaseError error
