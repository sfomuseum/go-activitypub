package activitypub

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/uris"
)

type DeliverPostOptions struct {
	From               *Account
	To                 string
	Post               *Post
	Hostname           string
	URIs               *uris.URIs
	DeliveriesDatabase DeliveriesDatabase
}

type DeliverPostToFollowersOptions struct {
	AccountsDatabase   AccountsDatabase
	FollowersDatabase  FollowersDatabase
	DeliveriesDatabase DeliveriesDatabase
	DeliveryQueue      DeliveryQueue
	Post               *Post
	URIs               *uris.URIs
}

func DeliverPostToFollowers(ctx context.Context, opts *DeliverPostToFollowersOptions) error {

	acct, err := opts.AccountsDatabase.GetAccountWithId(ctx, opts.Post.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account ID for post, %w", err)
	}

	followers_cb := func(ctx context.Context, follower_uri string) error {

		post_opts := &DeliverPostOptions{
			From:               acct,
			To:                 follower_uri,
			Post:               opts.Post,
			URIs:               opts.URIs,
			DeliveriesDatabase: opts.DeliveriesDatabase,
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

func DeliverPost(ctx context.Context, opts *DeliverPostOptions) error {

	slog.Debug("Deliver post", "post", opts.Post.Id, "from", opts.From.Id, "to", opts.To)

	note, err := NoteFromPost(ctx, opts.URIs, opts.From, opts.Post)

	if err != nil {
		return fmt.Errorf("Failed to derive note from post, %w", err)
	}

	from_uri := opts.From.AccountURL(ctx, opts.URIs).String()

	to_list := []string{
		opts.To,
	}

	create_activity, err := ap.NewCreateActivity(ctx, opts.URIs, from_uri, to_list, note)

	if err != nil {
		return fmt.Errorf("Failed to create activity from post, %w", err)
	}

	now := time.Now()
	ts := now.Unix()

	d := &Delivery{
		Id:        create_activity.Id,
		PostId:    opts.Post.Id,
		AccountId: opts.From.Id,
		Recipient: opts.To,
		Created:   ts,
	}

	defer func() {

		err := opts.DeliveriesDatabase.AddDelivery(ctx, d)

		if err != nil {
			slog.Error("Failed to add delivery", "post_id", opts.Post.Id, "recipienct", d.Recipient, "error", err)
		}
	}()

	post_opts := &PostToAccountOptions{
		From:     opts.From,
		To:       opts.To,
		Activity: create_activity,
		URIs:     opts.URIs,
	}

	err = PostToAccount(ctx, post_opts)

	now = time.Now()
	ts = now.Unix()

	d.Completed = ts
	d.Success = true

	if err != nil {
		d.Success = false
		d.Error = err.Error()
		return fmt.Errorf("Failed to post to inbox, %w", err)
	}

	return nil
}
