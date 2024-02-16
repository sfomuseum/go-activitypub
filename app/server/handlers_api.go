package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/http/api"
)

func webfingerHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &api.WebfingerHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
		Hostname:         run_opts.Hostname,
	}

	return api.WebfingerHandler(opts)
}

func profileHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &api.ProfileHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
		Hostname:         run_opts.Hostname,
	}

	return api.ProfileHandler(opts)
}

func inboxHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	setupFollowersDatabaseOnce.Do(setupFollowersDatabase)

	if setupFollowersDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupFollowersDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupFollowersDatabaseError)
	}

	opts := &api.InboxHandlerOptions{
		AccountsDatabase:  accounts_db,
		FollowersDatabase: followers_db,
		URIs:              run_opts.URIs,
		Hostname:          run_opts.Hostname,
	}

	return api.InboxHandler(opts)
}
