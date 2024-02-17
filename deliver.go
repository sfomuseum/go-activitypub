package activitypub

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-activitypub/ap"
)

type DeliverPostToFollowersOptions struct {
	AccountsDatabase  AccountsDatabase
	FollowersDatabase FollowersDatabase
	DeliveryQueue     DeliveryQueue
	Post              *Post
	Hostname          string
	URIs              *URIs
}

func DeliverPostToFollowers(ctx context.Context, opts *DeliverPostToFollowersOptions) error {

	acct, err := opts.AccountsDatabase.GetAccount(ctx, opts.Post.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	followers_cb := func(ctx context.Context, follower_id string) error {

		slog.Info("DELIVER", "to", follower_id)

		post_opts := &DeliverPostOptions{
			From:     acct,
			To:       follower_id,
			Post:     opts.Post,
			Hostname: opts.Hostname,
			URIs:     opts.URIs,
		}

		err := opts.DeliveryQueue.DeliverPost(ctx, post_opts)

		if err != nil {
			return fmt.Errorf("Failed to deliver post to %s, %w", follower_id, err)
		}

		return nil
	}

	err = opts.FollowersDatabase.GetFollowers(ctx, acct.Id, followers_cb)

	if err != nil {
		return fmt.Errorf("Failed to get followers for post author, %w", err)
	}

	return nil
}

func DeliverPostToAccount(ctx context.Context, opts *DeliverPostOptions) (*ap.Activity, error) {

	note, err := opts.Post.AsNote(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive note from post, %w", err)
	}

	from := fmt.Sprintf("%s@%s", opts.From.Id, opts.Hostname)

	to_list := []string{
		opts.To,
	}

	create_activity, err := ap.NewCreateActivity(ctx, from, to_list, note)

	if err != nil {
		return nil, fmt.Errorf("Failed to create activity from post, %w", err)
	}

	post_opts := &PostToAccountOptions{
		From:     opts.From,
		To:       opts.To,
		Message:  create_activity,
		Hostname: opts.Hostname,
		URIs:     opts.URIs,
	}

	activity, err := PostToAccount(ctx, post_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to post to inbox, %w", err)
	}

	return activity, nil
}
