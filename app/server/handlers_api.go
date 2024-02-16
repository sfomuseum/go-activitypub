package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/http/api"
)

func webfingerHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountDatabaseOnce.Do(setupAccountDatabase)

	if setupAccountDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountDatabaseError)
	}

	opts := &api.WebfingerHandlerOptions{
		AccountDatabase: account_db,
		URIs:            run_opts.URIs,
		Hostname:        run_opts.Hostname,
	}

	return api.WebfingerHandler(opts)
}

func profileHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountDatabaseOnce.Do(setupAccountDatabase)

	if setupAccountDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountDatabaseError)
	}

	opts := &api.ProfileHandlerOptions{
		AccountDatabase: account_db,
		URIs:            run_opts.URIs,
		Hostname:        run_opts.Hostname,
	}

	return api.ProfileHandler(opts)
}

func inboxHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountDatabaseOnce.Do(setupAccountDatabase)

	if setupAccountDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountDatabaseError)
	}

	opts := &api.InboxHandlerOptions{
		AccountDatabase: account_db,
		URIs:            run_opts.URIs,
		Hostname:        run_opts.Hostname,
	}

	return api.InboxHandler(opts)
}
