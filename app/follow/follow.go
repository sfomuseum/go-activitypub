package follow

import (
	"context"
	"flag"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/slog"
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

	logger := slog.Default()

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize accounts database, %w", err)
	}

	following_db, err := activitypub.NewFollowingDatabase(ctx, opts.FollowingDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize following database, %w", err)
	}

	messages_db, err := activitypub.NewMessagesDatabase(ctx, opts.MessagesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to initialize messages database, %w", err)
	}

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

	var activity *ap.Activity

	if opts.Undo {
		activity, err = ap.NewUndoFollowActivity(ctx, opts.URIs, follower_address, following_address)
	} else {
		activity, err = ap.NewFollowActivity(ctx, opts.URIs, follower_address, following_address)
	}

	if err != nil {
		return fmt.Errorf("Failed to create follow activity, %w", err)
	}

	post_opts := &activitypub.PostToAccountOptions{
		From:     follower_acct,
		To:       following_address,
		Activity: activity,
		URIs:     opts.URIs,
	}

	_, err = activitypub.PostToAccount(ctx, post_opts)

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
