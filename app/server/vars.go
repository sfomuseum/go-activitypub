package server

import (
	"sync"

	"github.com/sfomuseum/go-activitypub"
)

var run_opts *RunOptions

var actor_db activitypub.ActorDatabase

var setupActorDatabaseOnce sync.Once
var setupActorDatabaseError error
