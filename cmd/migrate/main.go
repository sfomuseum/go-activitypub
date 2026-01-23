package main

import (
	"context"
	"flag"
	"log"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/database"
)

func main() {

	var from_database_uri string
	var to_database_uri string

	flag.StringVar(&from_database_uri, "from", "", "...")
	flag.StringVar(&to_database_uri, "to", "", "...")
	flag.Parse()

	ctx := context.Background()

	from_db, err := database.NewAccountsDatabase(ctx, from_database_uri)

	if err != nil {
		log.Fatalf("Failed to create from database, %v", err)
	}

	defer from_db.Close()

	to_db, err := database.NewAccountsDatabase(ctx, to_database_uri)

	if err != nil {
		log.Fatalf("Failed to create to database, %v", err)
	}

	defer to_db.Close()

	cb := func(ctx context.Background, acct *activitypub.Account) {
		return to_db.AddAccount(ctx, acct)
	}

	err = from_db.GetAccounts(ctx, cb)

	if err != nil {
		log.Fatalf("Failed to get accounts, %v", err)
	}

}
