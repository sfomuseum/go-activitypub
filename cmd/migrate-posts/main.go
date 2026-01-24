package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func main() {

	var from_database_uri string
	var to_database_uri string
	var verbose bool

	flag.StringVar(&from_database_uri, "from", "", "...")
	flag.StringVar(&to_database_uri, "to", "", "...")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	flag.Parse()

	ctx := context.Background()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	slog.Debug("Set up from database")

	from_ctx, from_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer from_cancel()

	from_db, err := database.NewPostsDatabase(from_ctx, from_database_uri)

	if err != nil {
		log.Fatalf("Failed to create from database, %v", err)
	}

	defer from_db.Close(ctx)

	slog.Debug("Set up to database")

	to_ctx, to_cancel := context.WithTimeout(ctx, 5*time.Second)
	defer to_cancel()

	to_db, err := database.NewPostsDatabase(to_ctx, to_database_uri)

	if err != nil {
		log.Fatalf("Failed to create to database, %v", err)
	}

	defer to_db.Close(ctx)

	count := int64(0)
	success := int64(0)
	errors := int64(0)
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	done_ch := make(chan bool)

	go func() {
		for {
			select {
			case <-done_ch:
				return
			case <-ticker.C:
				slog.Info("Records processed", "count", atomic.LoadInt64(&count), "success", atomic.LoadInt64(&success), "errors", atomic.LoadInt64(&errors))
			}
		}
	}()

	cb := func(ctx context.Context, acct *activitypub.Post) error {

		defer atomic.AddInt64(&count, 1)

		slog.Debug("Add", "Post", acct.Id)
		err := to_db.AddPost(ctx, acct)

		if err != nil {
			slog.Error("Failed to add Post", "Post", acct.Id, "error", err)
			atomic.AddInt64(&errors, 1)
		} else {
			atomic.AddInt64(&success, 1)			
		}
		
		return nil
	}

	slog.Debug("Retrieve post")
	err = from_db.GetPosts(ctx, cb)

	if err != nil {
		log.Fatalf("Failed to get posts, %v", err)
	}

	done_ch <- true

	slog.Info("Total records processed", "count", atomic.LoadInt64(&count), "success", atomic.LoadInt64(&success), "errors", atomic.LoadInt64(&errors))
}
