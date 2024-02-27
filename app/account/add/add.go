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

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/id"
)

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

	account_type, err := activitypub.AccountTypeFromString(opts.AccountType)

	if err != nil {
		return fmt.Errorf("Failed to derive account type from string, %w", err)
	}

	icon_uri := ""

	if opts.AccountIconURI != "" {

		if re_http_url.MatchString(opts.AccountIconURI) {

			if !opts.AllowRemoteIconURI {
				return fmt.Errorf("Remote account icon URIs are not allowed")
			}

			icon_u, err := url.Parse(opts.AccountIconURI)

			if err != nil {
				return fmt.Errorf("Failed to parse remote icon URI, %w", err)
			}

			icon_uri = icon_u.String()

		} else {

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
		}
	}

	a := &activitypub.Account{
		Id:            account_id,
		Name:          opts.AccountName,
		AccountType:   account_type,
		DisplayName:   opts.DisplayName,
		Blurb:         opts.Blurb,
		URL:           opts.URL,
		PrivateKeyURI: private_key_uri,
		PublicKeyURI:  public_key_uri,
		IconURI:       icon_uri,
	}

	a, err = activitypub.AddAccount(ctx, db, a)

	if err != nil {
		return fmt.Errorf("Failed to add new account, %w", err)
	}

	return nil
}
