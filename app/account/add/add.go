package add

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/id"
)

// Reconcile with www/icon.go
var re_http_url = regexp.MustCompile(`^https?\:\/\/(.*)`)

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

	accounts_db, err := activitypub.NewAccountsDatabase(ctx, opts.AccountsDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate accounts database, %w", err)
	}

	defer accounts_db.Close(ctx)

	aliases_db, err := activitypub.NewAliasesDatabase(ctx, opts.AliasesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to instantiate aliases database, %w", err)
	}

	defer aliases_db.Close(ctx)

	properties_db, err := activitypub.NewPropertiesDatabase(ctx, opts.PropertiesDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create properties database, %w", err)
	}

	defer properties_db.Close(ctx)

	// START OF check for existing account name and aliases

	acct_taken, err := activitypub.IsAccountNameTaken(ctx, accounts_db, opts.AccountName)

	if err != nil {
		return fmt.Errorf("Failed to determine if account name is taken, %w", err)
	}

	if acct_taken {
		return fmt.Errorf("Account name is not available")
	}

	for _, name := range opts.Aliases {

		alias_taken, err := activitypub.IsAliasNameTaken(ctx, aliases_db, name)

		if err != nil {
			return fmt.Errorf("Failed to determine if alias name '%s' is taken, %w", name, err)
		}

		if alias_taken {
			return fmt.Errorf("Account name '%s' is not available", name)
		}
	}

	// END OF check for existing account name and aliases

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

	account_type, err := activitypub.AccountTypeFromString(opts.AccountType)

	if err != nil {
		return fmt.Errorf("Failed to derive account type from string, %w", err)
	}

	icon_uri := ""

	if opts.AccountIconURI != "" {

		switch {
		case re_http_url.MatchString(opts.AccountIconURI):

			if !opts.AllowRemoteIconURI {
				return fmt.Errorf("Remote account icon URIs are not allowed")
			}

			icon_u, err := url.Parse(opts.AccountIconURI)

			if err != nil {
				return fmt.Errorf("Failed to parse remote icon URI, %w", err)
			}

			icon_uri = icon_u.String()

		case opts.EmbedIconURI:

			r, err := os.Open(opts.AccountIconURI)

			if err != nil {
				return fmt.Errorf("Failed to open icon URI for reading, %w", err)
			}

			defer r.Close()

			data, err := io.ReadAll(r)

			if err != nil {
				return fmt.Errorf("Failed to read icon URI, %w", err)
			}

			br := bytes.NewReader(data)

			_, format, err := image.Decode(br)

			if err != nil {
				return fmt.Errorf("Failed to decode icon URI, %w", err)
			}

			b64 := base64.StdEncoding.EncodeToString(data)

			icon_uri = fmt.Sprintf("data:image/%s;base64,%s", format, b64)

		default:
			icon_uri = opts.AccountIconURI
		}
	}

	a := &activitypub.Account{
		Id:            account_id,
		Name:          opts.AccountName,
		AccountType:   account_type,
		DisplayName:   opts.DisplayName,
		Blurb:         opts.Blurb,
		Discoverable:  opts.Discoverable,
		URL:           opts.URL,
		PrivateKeyURI: private_key_uri,
		PublicKeyURI:  public_key_uri,
		IconURI:       icon_uri,
	}

	a, err = activitypub.AddAccount(ctx, accounts_db, a)

	if err != nil {
		return fmt.Errorf("Failed to add new account, %w", err)
	}

	// Properties

	_, _, err = activitypub.ApplyPropertiesUpdates(ctx, properties_db, a, opts.Properties)

	if err != nil {
		return fmt.Errorf("Account created (%d) but failed to apply properties, %w", a.Id, err)
	}

	// Aliases

	for _, name := range opts.Aliases {

		a, err := aliases_db.GetAliasWithName(ctx, name)

		if err != nil && err != activitypub.ErrNotFound {
			return fmt.Errorf("Failed to retrieve alias for name '%s', %w", name, err)
		}

		if a != nil {

			if a.AccountId != account_id {
				return fmt.Errorf("Alias '%s' is already in use", name)
			}
		}

		now := time.Now()
		ts := now.Unix()

		a = &activitypub.Alias{
			Name:      name,
			AccountId: account_id,
			Created:   ts,
		}

		err = aliases_db.AddAlias(ctx, a)

		if err != nil {
			return fmt.Errorf("Failed to add alias for name '%s', %w", name, err)
		}
	}

	return nil
}
