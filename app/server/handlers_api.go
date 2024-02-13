package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/http/api"
)

func webfingerHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupActorDatabaseOnce.Do(setupActorDatabase)

	if setupActorDatabaseError != nil {
		slog.Error("Failed to set up actor database configuration", "error", setupActorDatabaseError)
		return nil, fmt.Errorf("Failed to set up actor database configuration, %w", setupActorDatabaseError)
	}

	opts := &api.WebfingerHandlerOptions{
		ActorDatabase: actor_db,
		URIs:          run_opts.URIs,
		Hostname:      run_opts.Hostname,
	}

	return api.WebfingerHandler(opts)
}

func profileHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupActorDatabaseOnce.Do(setupActorDatabase)

	if setupActorDatabaseError != nil {
		slog.Error("Failed to set up actor database configuration", "error", setupActorDatabaseError)
		return nil, fmt.Errorf("Failed to set up actor database configuration, %w", setupActorDatabaseError)
	}

	opts := &api.ProfileHandlerOptions{
		ActorDatabase: actor_db,
		URIs:          run_opts.URIs,
		Hostname:      run_opts.Hostname,
	}

	return api.ProfileHandler(opts)
}

func inboxHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupActorDatabaseOnce.Do(setupActorDatabase)

	if setupActorDatabaseError != nil {
		slog.Error("Failed to set up actor database configuration", "error", setupActorDatabaseError)
		return nil, fmt.Errorf("Failed to set up actor database configuration, %w", setupActorDatabaseError)
	}

	opts := &api.InboxHandlerOptions{
		ActorDatabase: actor_db,
		URIs:          run_opts.URIs,
		Hostname:      run_opts.Hostname,
	}

	return api.InboxHandler(opts)
}
