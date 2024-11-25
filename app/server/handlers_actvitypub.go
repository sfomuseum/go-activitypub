package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"github.com/sfomuseum/go-activitypub/www"
)

func webfingerHandlerFunc(ctx context.Context) (http.Handler, error) {

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

	opts := &www.WebfingerHandlerOptions{
		AccountsDatabase: accounts_db,
		AliasesDatabase:  aliases_db,
		URIs:             run_opts.URIs,
	}

	wf_handler, err := www.WebfingerHandler(opts)

	if err != nil {
		return nil, err
	}

	wf_handler = cors.Default().Handler(wf_handler)
	return wf_handler, nil
}

func inboxPostHandlerFunc(ctx context.Context) (http.Handler, error) {

	// START OF do this concurrently?
	// Probably only marginally faster at the expense of hard-to-follow code...

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

	setupLikesDatabaseOnce.Do(setupLikesDatabase)

	if setupLikesDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupLikesDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupLikesDatabaseError)
	}

	setupBoostsDatabaseOnce.Do(setupBoostsDatabase)

	if setupBoostsDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupBoostsDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupBoostsDatabaseError)
	}

	setupPostsDatabaseOnce.Do(setupPostsDatabase)

	if setupPostsDatabaseError != nil {
		slog.Error("Failed to set up follower database configuration", "error", setupPostsDatabaseError)
		return nil, fmt.Errorf("Failed to set up follower database configuration, %w", setupPostsDatabaseError)
	}

	setupProcessMessageQueueOnce.Do(setupProcessMessageQueue)

	if setupProcessMessageQueueError != nil {
		slog.Error("Failed to set up process message queue", "error", setupProcessMessageQueueError)
		return nil, fmt.Errorf("Failed to set up process message queue, %w", setupProcessMessageQueueError)
	}

	setupProcessFollowerQueueOnce.Do(setupProcessFollowerQueue)

	if setupProcessFollowerQueueError != nil {
		slog.Error("Failed to set up process follow queue", "error", setupProcessFollowerQueueError)
		return nil, fmt.Errorf("Failed to set up process follow queue, %w", setupProcessFollowerQueueError)
	}

	// END OF do this concurrently?

	opts := &www.InboxPostHandlerOptions{
		AccountsDatabase:     accounts_db,
		FollowersDatabase:    followers_db,
		FollowingDatabase:    following_db,
		NotesDatabase:        notes_db,
		MessagesDatabase:     messages_db,
		BlocksDatabase:       blocks_db,
		PostsDatabase:        posts_db,
		LikesDatabase:        likes_db,
		BoostsDatabase:       boosts_db,
		URIs:                 run_opts.URIs,
		AllowFollow:          run_opts.AllowFollow,
		AllowCreate:          run_opts.AllowCreate,
		AllowLikes:           run_opts.AllowLikes,
		AllowBoosts:          run_opts.AllowBoosts,
		AllowMentions:        run_opts.AllowMentions,
		ProcessMessageQueue:  process_message_queue,
		ProcessFollowerQueue: process_follower_queue,
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

func outboxPostHandlerFunc(ctx context.Context) (http.Handler, error) {

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

	opts := &www.OutboxPostHandlerOptions{
		AccountsDatabase: accounts_db,
		PostsDatabase:    posts_db,
		URIs:             run_opts.URIs,
	}

	return www.OutboxPostHandler(opts)
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
		AllowRemote:      run_opts.AllowRemoteIconURI,
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
