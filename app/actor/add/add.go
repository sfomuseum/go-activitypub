package add

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/crypto"
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

	db, err := activitypub.NewActorDatabase(ctx, opts.ActorDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	private_key, err := crypto.GeneratePrivateKey(4096)

	if err != nil {
		return fmt.Errorf("Failed to generate private key, %w", err)
	}

	private_pem := crypto.EncodePrivateKeyToPEM(private_key)

	// Is thi really PEM...?
	public_pem, err := crypto.GeneratePublicKey(&private_key.PublicKey)

	if err != nil {
		return fmt.Errorf("Failed to generate public key, %w", err)
	}

	private_key_uri := fmt.Sprintf("constant://?val=%s", url.QueryEscape(string(private_pem)))
	public_key_uri := fmt.Sprintf("constant://?val=%s", url.QueryEscape(string(public_pem)))

	a := &activitypub.Actor{
		Id:            opts.ActorId,
		PrivateKeyURI: private_key_uri,
		PublicKeyURI:  public_key_uri,
	}

	err = db.AddActor(ctx, a)

	if err != nil {
		return fmt.Errorf("Failed to add new actor, %w", err)
	}

	return nil
}
