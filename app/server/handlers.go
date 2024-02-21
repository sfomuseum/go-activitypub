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
	}

	return www.WebfingerHandler(opts)
}

func accountHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &www.AccountHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
		Templates:        run_opts.Templates,
	}

	return www.AccountHandler(opts)
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

	setupFollowingDatabaseOnce.Do(setupFollowingDatabase)

	if setupFollowingDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupFollowingDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupFollowingDatabaseError)
	}

	setupNotesDatabaseOnce.Do(setupNotesDatabase)

	if setupNotesDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupNotesDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupNotesDatabaseError)
	}

	setupMessagesDatabaseOnce.Do(setupMessagesDatabase)

	if setupMessagesDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupMessagesDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupMessagesDatabaseError)
	}

	setupBlocksDatabaseOnce.Do(setupBlocksDatabase)

	if setupBlocksDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupBlocksDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupBlocksDatabaseError)
	}

	opts := &www.InboxPostHandlerOptions{
		AccountsDatabase:  accounts_db,
		FollowersDatabase: followers_db,
		FollowingDatabase: following_db,
		NotesDatabase:     notes_db,
		MessagesDatabase:  messages_db,
		BlocksDatabase:    blocks_db,
		URIs:              run_opts.URIs,
		AllowFollow:       run_opts.AllowFollow,
		AllowCreate:       run_opts.AllowCreate,
	}

	return www.InboxPostHandler(opts)
}

func outboxGetHandlerFunc(ctx context.Context) (http.Handler, error) {

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

	opts := &www.OutboxGetHandlerOptions{
		AccountsDatabase: accounts_db,
		PostsDatabase:    posts_db,
		URIs:             run_opts.URIs,
	}

	return www.OutboxGetHandler(opts)
}

func iconHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	opts := &www.IconHandlerOptions{
		AccountsDatabase: accounts_db,
		URIs:             run_opts.URIs,
	}

	return www.IconHandler(opts)
}

func followingHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupAccountsDatabaseOnce.Do(setupAccountsDatabase)

	if setupAccountsDatabaseError != nil {
		slog.Error("Failed to set up account database configuration", "error", setupAccountsDatabaseError)
		return nil, fmt.Errorf("Failed to set up account database configuration, %w", setupAccountsDatabaseError)
	}

	setupFollowingDatabaseOnce.Do(setupFollowingDatabase)

	if setupFollowingDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupFollowingDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupFollowingDatabaseError)
	}

	opts := &www.FollowingHandlerOptions{
		AccountsDatabase:  accounts_db,
		FollowingDatabase: following_db,
		URIs:              run_opts.URIs,
	}

	return www.FollowingHandler(opts)
}

func followersHandlerFunc(ctx context.Context) (http.Handler, error) {

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

	opts := &www.FollowersHandlerOptions{
		AccountsDatabase:  accounts_db,
		FollowersDatabase: followers_db,
		URIs:              run_opts.URIs,
	}

	return www.FollowersHandler(opts)
}
