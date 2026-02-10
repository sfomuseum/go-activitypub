package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-activitypub/database"
)

func main() {

	var from_database_uri string
	var to_database_uri string
	var database_label string
	var verbose bool

	flag.StringVar(&from_database_uri, "from", "", "...")
	flag.StringVar(&to_database_uri, "to", "", "...")
	flag.StringVar(&database_label, "database", "", "...")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	flag.Parse()

	ctx := context.Background()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

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

	switch database_label {
	case "accounts":

		err := database.MigrateAccountsDatabaseFromURIs(ctx, from_database_uri, to_database_uri, &count, &success, &errors)

		if err != nil {
			log.Fatalf("Failed to migrate accounts database, %v", err)
		}

	case "activities":

		err := database.MigrateActivitiesDatabaseFromURIs(ctx, from_database_uri, to_database_uri, &count, &success, &errors)

		if err != nil {
			log.Fatalf("Failed to migrate activities database, %v", err)
		}

	case "aliases":

		err := database.MigrateAliasesDatabaseFromURIs(ctx, from_database_uri, to_database_uri, &count, &success, &errors)

		if err != nil {
			log.Fatalf("Failed to migrate database, %v", err)
		}

	case "post":

		err := database.MigratePostsDatabaseFromURIs(ctx, from_database_uri, to_database_uri, &count, &success, &errors)

		if err != nil {
			log.Fatalf("Failed to migrate database, %v", err)
		}

	default:
		slog.Error("Unsupported database", "database", database_label)
	}

	done_ch <- true

	slog.Info("Total records processed", "count", atomic.LoadInt64(&count), "success", atomic.LoadInt64(&success), "errors", atomic.LoadInt64(&errors))
}
