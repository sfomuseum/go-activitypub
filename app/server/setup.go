package server

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func setupAccountsDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	accounts_db, err = activitypub.NewAccountsDatabase(ctx, run_opts.AccountsDatabaseURI)

	if err != nil {
		setupAccountsDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}
