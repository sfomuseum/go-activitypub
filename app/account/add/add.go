package add

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/id"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	slog.SetDefault(logger)

	db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	if opts.PublicKeyURI == "" && opts.PrivateKeyURI != "" {
		return fmt.Errorf("Missing public key URI")
	}

	if opts.PublicKeyURI != "" && opts.PrivateKeyURI == "" {
		return fmt.Errorf("Missing private key URI")
	}

	var private_key_uri string
	var public_key_uri string

	if opts.PublicKeyURI != "" && opts.PrivateKeyURI != "" {
		public_key_uri = opts.PublicKeyURI
		private_key_uri = opts.PrivateKeyURI
	} else {
		private_pem, public_pem, err := crypto.GenerateKeyPair(4096)

		if err != nil {
			return fmt.Errorf("Failed to generate private key, %w", err)
		}

		private_key_uri = fmt.Sprintf("constant://?val=%s", url.QueryEscape(string(private_pem)))
		public_key_uri = fmt.Sprintf("constant://?val=%s", url.QueryEscape(string(public_pem)))
	}

	account_id := opts.AccountId

	if account_id == 0 {

		id, err := id.NewId()

		if err != nil {
			return fmt.Errorf("Failed to create new account ID, %w", err)
		}

		account_id = id
	}

	a := &activitypub.Account{
		Id:            account_id,
		Name:          opts.AccountName,
		DisplayName:   opts.DisplayName,
		Blurb:         opts.Blurb,
		URL:           opts.URL,
		PrivateKeyURI: private_key_uri,
		PublicKeyURI:  public_key_uri,
	}

	a, err = activitypub.AddAccount(ctx, db, a)

	if err != nil {
		return fmt.Errorf("Failed to add new account, %w", err)
	}

	return nil
}
