package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/www"
)

func webfingerHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &www.WebfingerHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
		Hostname:         run_opts.Hostname,
	}

	return www.WebfingerHandler(opts)
}

func profileHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &www.ProfileHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
		Hostname:         run_opts.Hostname,
	}

	return www.ProfileHandler(opts)
}

func inboxPostHandlerFunc(ctx context.Context) (http.Handler, error) {

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

	opts := &www.InboxPostHandlerOptions{
		AccountsDatabase:  accounts_db,
		FollowersDatabase: followers_db,
		URIs:              run_opts.URIs,
		Hostname:          run_opts.Hostname,
	}

	return www.InboxPostHandler(opts)
}

func inboxGetHandlerFunc(ctx context.Context) (http.Handler, error) {

	/*
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
	*/

	opts := &www.InboxGetHandlerOptions{}

	return www.InboxGetHandler(opts)
}
