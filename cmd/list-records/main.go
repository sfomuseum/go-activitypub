package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/sfomuseum/go-activitypub/database"
)

func main() {

	var database_model string
	var database_uri string
	var verbose bool

	flag.StringVar(&database_model, "database-model", "", "...")
	flag.StringVar(&database_uri, "database-uri", "", "...")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	flag.Parse()

	ctx := context.Background()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	db, err := database.NewAccountsDatabase(ctx, database_uri)

	if err != nil {
		log.Fatalf("Failed to create database for model, %v", err)
	}

	for acct, err := range db.QueryRecords(ctx, nil) {

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d %s\n", acct.Id, acct.Name)
	}

}
