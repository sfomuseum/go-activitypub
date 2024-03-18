package stats

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-activitypub"
)

type Counts struct {
	Date       string `json:"date"`
	Location   string `json:"location"`
	Accounts   int64  `json:"accounts"`
	Blocks     int64  `json:"blocks"`
	Boosts     int64  `json:"boosts"`
	Deliveries int64  `json:"deliveries"`
	Followers  int64  `json:"followers"`
	Following  int64  `json:"following"`
	Likes      int64  `json:"likes"`
	Messages   int64  `json:"messages"`
	Notes      int64  `json:"notes"`
	Posts      int64  `json:"posts"`
}

type CountsForDateOptions struct {
	Date               string
	Location           string
	AccountsDatabase   activitypub.AccountsDatabase
	BlocksDatabase     activitypub.BlocksDatabase
	BoostsDatabase     activitypub.BoostsDatabase
	DeliveriesDatabase activitypub.DeliveriesDatabase
	FollowersDatabase  activitypub.FollowersDatabase
	FollowingDatabase  activitypub.FollowingDatabase
	LikesDatabase      activitypub.LikesDatabase
	MessagesDatabase   activitypub.MessagesDatabase
	NotesDatabase      activitypub.NotesDatabase
	PostsDatabase      activitypub.PostsDatabase
}

func CountsForDate(ctx context.Context, opts *CountsForDateOptions) (*Counts, error) {

	counts := &Counts{
		Date: opts.Date,
	}

	t, err := time.Parse("2006-01-02", opts.Date)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse day, %w", err)
	}

	if opts.Location != "" {

		loc, err := time.LoadLocation(opts.Location)

		if err != nil {
			return nil, fmt.Errorf("Failed to load location, %w", err)
		}

		t = t.In(loc)

		counts.Location = loc.String()
	}

	start := t.Unix()
	end := start + ONEDAY

	slog.Debug("Get counts for date", "date", opts.Date, "location", opts.Location, "start", start, "end", end)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	mu := new(sync.RWMutex)
	remaining := int32(0)

	go func() {

		defer func() {
			done_ch <- true
		}()

		atomic.AddInt32(&remaining, 1)

		i, err := CountAccountsForDateRange(ctx, opts.AccountsDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for accounts, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Accounts = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountBlocksForDateRange(ctx, opts.BlocksDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for blocks, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Blocks = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountBoostsForDateRange(ctx, opts.BoostsDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for boosts, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Boosts = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountDeliveriesForDateRange(ctx, opts.DeliveriesDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for deliveries, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Deliveries = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountFollowersForDateRange(ctx, opts.FollowersDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for followers, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Followers = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountFollowingForDateRange(ctx, opts.FollowingDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for following, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Following = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountLikesForDateRange(ctx, opts.LikesDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for likes, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Likes = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountMessagesForDateRange(ctx, opts.MessagesDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for messages, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Messages = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountNotesForDateRange(ctx, opts.NotesDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for notes, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Notes = i
	}()

	go func() {

		atomic.AddInt32(&remaining, 1)

		defer func() {
			done_ch <- true
		}()

		i, err := CountPostsForDateRange(ctx, opts.PostsDatabase, start, end)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive counts for posts, %w", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		counts.Posts = i
	}()

	for atomic.LoadInt32(&remaining) > 0 {
		select {
		case <-done_ch:
			atomic.AddInt32(&remaining, -1)
		case err := <-err_ch:
			return nil, err
		}
	}

	return counts, nil
}
