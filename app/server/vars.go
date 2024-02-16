package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub"
)

var run_opts *RunOptions

var account_db activitypub.AccountDatabase

var setupAccountDatabaseOnce sync.Once
var setupAccountDatabaseError error
