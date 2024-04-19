package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-activitypub/www"
)

func accountHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	setupAliasesDatabaseOnce.Do(setupAliasesDatabase)

	if setupAliasesDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAliasesDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAliasesDatabaseError)
	}
	
	opts := &www.AccountHandlerOptions{
		AccountsDatabase: accounts_db,
		AliasesDatabase: aliases_db,		
		URIs:             run_opts.URIs,
		Templates:        run_opts.Templates,
	}

	h, err := www.AccountHandler(opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create account handler, %w", err)
	}

	if run_opts.AccountHandlerMiddleware != nil {
		h = run_opts.AccountHandlerMiddleware(h)
	}

	return h, nil
}

func postHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	setupPostsDatabaseOnce.Do(setupPostsDatabase)

	if setupPostsDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupPostsDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupPostsDatabaseError)
	}

	opts := &www.PostHandlerOptions{
		AccountsDatabase: accounts_db,
		PostsDatabase:    posts_db,
		URIs:             run_opts.URIs,
		Templates:        run_opts.Templates,
	}

	h, err := www.PostHandler(opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to create post handler, %w", err)
	}

	if run_opts.AccountHandlerMiddleware != nil {
		h = run_opts.AccountHandlerMiddleware(h)
	}

	return h, nil
}
