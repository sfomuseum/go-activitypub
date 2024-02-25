package activitypub

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverPostToFollowersOptions struct {
	AccountsDatabase  AccountsDatabase
	FollowersDatabase FollowersDatabase
	DeliveryQueue     DeliveryQueue
	Post              *Post
	URIs              *uris.URIs
}

func DeliverPostToFollowers(ctx context.Context, opts *DeliverPostToFollowersOptions) error {

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Post.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	followers_cb := func(ctx context.Context, follower_uri string) error {

		post_opts := &DeliverPostOptions{
			From: acct,
			To:   follower_uri,
			Post: opts.Post,
			URIs: opts.URIs,
		}

		err := opts.DeliveryQueue.DeliverPost(ctx, post_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_uri, err)
		}

		return nil
	}

	err = opts.FollowersDatabase.GetFollowersForAccount(ctx, acct.Id, followers_cb)

	if err != nil {
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	return nil
}

func DeliverPostToAccount(ctx context.Context, opts *DeliverPostOptions) error {

	slog.Debug("Deliver post", "post", opts.Post.Id, "from", opts.From.Id, "to", opts.To)

	note, err := opts.Post.AsNote(ctx)

	if err != nil {
		return fmt.Errorf("Failed to derive note from post, %w", err)
	}

	from_uri := opts.From.Address(opts.URIs.Hostname)

	to_list := []string{
		opts.To,
	}

	create_activity, err := ap.NewCreateActivity(ctx, opts.URIs, from_uri, to_list, note)

	if err != nil {
		return fmt.Errorf("Failed to create activity from post, %w", err)
	}

	post_opts := &PostToAccountOptions{
		From:    opts.From,
		To:      opts.To,
		Message: create_activity,
		URIs:    opts.URIs,
	}

	err = PostToAccount(ctx, post_opts)

	if err != nil {
		return fmt.Errorf("Failed to post to inbox, %w", err)
	}

	return nil
}
