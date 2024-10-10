package follow

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/inbox"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	accounts_db, err := database.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	following_db, err := database.NewFollowingDatabase(ctx, opts.FollowingDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize following database, %w", err)
	}

	defer following_db.Close(ctx)

	messages_db, err := database.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize messages database, %w", err)
	}

	defer messages_db.Close(ctx)

	follower_acct, err := accounts_db.GetAccountWithName(ctx, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountName, err)
	}

	follower_id := follower_acct.Id

	// See this? It is important to pass the fully-qualifier follower URI so the
	// endpoint receiving the follow activity can figure out where (which hostname)
	// to make a webfinger/profile query.

	follower_address := follower_acct.Address(opts.URIs.Hostname)
	following_address := opts.FollowAddress

	logger = logger.With("follower", follower_address)
	logger = logger.With("following", following_address)

	following_actor, err := activitypub.RetrieveActor(ctx, following_address, opts.URIs.Insecure)

	if err != nil {
		return fmt.Errorf("Failed to retrieve actor for %s, %w", following_address, err)
	}

	following_inbox := following_actor.Inbox
	logger = logger.With("inbox", following_inbox)

	var activity *ap.Activity

	if opts.Undo {
		logger.Info("Create unfollow activity")
		activity, err = ap.NewUndoFollowActivity(ctx, opts.URIs, follower_address, following_address)
	} else {
		logger.Info("Create follow activity")
		activity, err = ap.NewFollowActivity(ctx, opts.URIs, follower_address, following_address)
	}

	if err != nil {
		return fmt.Errorf("Failed to create follow activity, %w", err)
	}

	// enc := json.NewEncoder(os.Stdout)
	// enc.Encode(activity)

	post_opts := &inbox.PostToInboxOptions{
		From:     follower_acct,
		Inbox:    following_inbox,
		Activity: activity,
		URIs:     opts.URIs,
	}

	err = inbox.PostToInbox(ctx, post_opts)

	if err != nil {
		return fmt.Errorf("Failed to deliver follow activity, %w", err)
	}

	if undo {

		f, err := following_db.GetFollowing(ctx, follower_id, following_address)

		if err != nil {
			return fmt.Errorf("Failed to retrieve following, %w", err)
		}

		err = following_db.RemoveFollowing(ctx, f)

		if err != nil {
			return fmt.Errorf("Unfollow request was successful but unable to register unfollowing locally, %w", err)
		}

		msg_cb := func(ctx context.Context, m *activitypub.Message) error {
			logger.Info("Remove message", "id", m.Id)
			return messages_db.RemoveMessage(ctx, m)
		}

		err = messages_db.GetMessagesForAccountAndAuthor(ctx, follower_id, following_address, msg_cb)

		logger.Info("Unfollowing successful")
		return nil
	}

	f, err := activitypub.NewFollowing(ctx, follower_id, following_address)

	if err != nil {
		return fmt.Errorf("Failed to create new following, %w", err)
	}

	err = following_db.AddFollowing(ctx, f)

	if err != nil {
		return fmt.Errorf("Follow request was successful but unable to register following locally, %w", err)
	}

	logger.Info("Following successful")
	return nil
}
