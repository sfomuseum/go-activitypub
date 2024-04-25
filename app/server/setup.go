package server

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func setupAccountsDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	accounts_db, err = activitypub.NewAccountsDatabase(ctx, run_opts.AccountsDatabaseURI)

	if err != nil {
		setupAccountsDatabaseError = fmt.Errorf("Failed to set up accounts database, %w", err)
		return
	}
}

func setupAliasesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	aliases_db, err = activitypub.NewAliasesDatabase(ctx, run_opts.AliasesDatabaseURI)

	if err != nil {
		setupAliasesDatabaseError = fmt.Errorf("Failed to set up aliases database, %w", err)
		return
	}
}

func setupFollowersDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	followers_db, err = activitypub.NewFollowersDatabase(ctx, run_opts.FollowersDatabaseURI)

	if err != nil {
		setupFollowersDatabaseError = fmt.Errorf("Failed to set up followers database, %w", err)
		return
	}
}

func setupFollowingDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	following_db, err = activitypub.NewFollowingDatabase(ctx, run_opts.FollowingDatabaseURI)

	if err != nil {
		setupFollowingDatabaseError = fmt.Errorf("Failed to set up following database, %w", err)
		return
	}
}

func setupNotesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	notes_db, err = activitypub.NewNotesDatabase(ctx, run_opts.NotesDatabaseURI)

	if err != nil {
		setupNotesDatabaseError = fmt.Errorf("Failed to set up notes database, %w", err)
		return
	}
}

func setupMessagesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	messages_db, err = activitypub.NewMessagesDatabase(ctx, run_opts.MessagesDatabaseURI)

	if err != nil {
		setupMessagesDatabaseError = fmt.Errorf("Failed to set up messages database, %w", err)
		return
	}
}

func setupBlocksDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	blocks_db, err = activitypub.NewBlocksDatabase(ctx, run_opts.BlocksDatabaseURI)

	if err != nil {
		setupBlocksDatabaseError = fmt.Errorf("Failed to set up blocks database, %w", err)
		return
	}
}

func setupPostsDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	posts_db, err = activitypub.NewPostsDatabase(ctx, run_opts.PostsDatabaseURI)

	if err != nil {
		setupPostsDatabaseError = fmt.Errorf("Failed to set up posts database, %w", err)
		return
	}
}

func setupLikesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	likes_db, err = activitypub.NewLikesDatabase(ctx, run_opts.LikesDatabaseURI)

	if err != nil {
		setupLikesDatabaseError = fmt.Errorf("Failed to set up likes database, %w", err)
		return
	}
}

func setupBoostsDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	boosts_db, err = activitypub.NewBoostsDatabase(ctx, run_opts.BoostsDatabaseURI)

	if err != nil {
		setupBoostsDatabaseError = fmt.Errorf("Failed to set up boosts database, %w", err)
		return
	}
}

func setupPropertiesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	properties_db, err = activitypub.NewPropertiesDatabase(ctx, run_opts.PropertiesDatabaseURI)

	if err != nil {
		setupPropertiesDatabaseError = fmt.Errorf("Failed to set up properties database, %w", err)
		return
	}
}

func setupProcessMessageQueue() {

	ctx := context.Background()
	var err error

	process_message_queue, err = activitypub.NewProcessMessageQueue(ctx, run_opts.ProcessMessageQueueURI)

	if err != nil {
		setupProcessMessageQueueError = fmt.Errorf("Failed to create process message queue, %w", err)
		return
	}

}
