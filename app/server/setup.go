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
		setupAccountsDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}

func setupFollowersDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	followers_db, err = activitypub.NewFollowersDatabase(ctx, run_opts.FollowersDatabaseURI)

	if err != nil {
		setupFollowersDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}

func setupFollowingDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	following_db, err = activitypub.NewFollowingDatabase(ctx, run_opts.FollowingDatabaseURI)

	if err != nil {
		setupFollowingDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}

func setupNotesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	notes_db, err = activitypub.NewNotesDatabase(ctx, run_opts.NotesDatabaseURI)

	if err != nil {
		setupNotesDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}

func setupMessagesDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	messages_db, err = activitypub.NewMessagesDatabase(ctx, run_opts.MessagesDatabaseURI)

	if err != nil {
		setupMessagesDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}
