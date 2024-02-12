package server

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func setupActorDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	actor_db, err = activitypub.NewActorDatabase(ctx, run_opts.ActorDatabaseURI)

	if err != nil {
		setupActorDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}
