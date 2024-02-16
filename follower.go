package activitypub

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Follower struct {
	// The account doing the following
	FollowerId string `json:"follower_id"`
	// The account being followed
	FollowingId  string `json:"following_id"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}

func AddFollower(ctx context.Context, db FollowersDatabase, a *Follower) (*Follower, error) {

	now := time.Now()
	ts := now.Unix()

	a.Created = ts
	a.LastModified = ts

	err := db.AddFollower(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to add follower, %w", err)
	}

	return a, nil
}

func UpdateFollower(ctx context.Context, db FollowersDatabase, a *Follower) (*Follower, error) {

	now := time.Now()
	ts := now.Unix()

	a.LastModified = ts

	err := db.UpdateFollower(ctx, a)

	if err != nil {
		return nil, fmt.Errorf("Failed to update follower, %w", err)
	}

	return a, nil
}

func ParseFollowerURI(uri string) (string, string, error) {

	parts := strings.Split(uri, "@")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid address")
	}

	return parts[0], parts[1], nil
}
