package server

import (
	"context"
	"fmt"

	"github.com/sfomuseum/go-activitypub"
)

func setupAccountDatabase() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	account_db, err = activitypub.NewAccountDatabase(ctx, run_opts.AccountDatabaseURI)

	if err != nil {
		setupAccountDatabaseError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}
