package follow

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"io"
	"os"

	"github.com/99designs/httpsignatures-go"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
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

	db, err := activitypub.NewActorDatabase(ctx, opts.AccountDatabaseURI)

	if err != nil {
		return fmt.Errorf("Failed to create new database, %w", err)
	}

	acct, err := db.GetActor(ctx, opts.AccountId)

	if err != nil {
		return fmt.Errorf("Failed to retrieve account %s, %w", opts.AccountId, err)
	}

	acct_url, err := acct.ProfileURL(ctx, opts.URIs)

	if err != nil {
		return fmt.Errorf("Failed to derive profile URL for account, %w", err)
	}

	acct_url.Scheme = "http"

	follower := acct_url.String()

	follow_req, err := ap.NewFollowActivity(ctx, follower, opts.Follow)

	if err != nil {
		return fmt.Errorf("Failed to create follow activity, %w", err)
	}

	enc_req, err := json.Marshal(follow_req)

	if err != nil {
		return fmt.Errorf("Failed to marshal follow activity request, %w", err)
	}

	http_req, err := http.NewRequestWithContext(ctx, "POST", opts.Inbox, bytes.NewBuffer(enc_req))

	if err != nil {
		return fmt.Errorf("Failed to create new request to %s, %w", opts.Inbox, err)
	}

	key_id := follower

	public_key, err := acct.PublicKey(ctx)

	if err != nil {
		return fmt.Errorf("Failed to get private key, %w", err)
	}

	err = httpsignatures.DefaultSha256Signer.SignRequest(key_id, public_key, http_req)

	if err != nil {
		return fmt.Errorf("Failed to sign request, %w", err)
	}

	slog.Info("OK", "signature", http_req.Header.Get("Signature"))

	http_cl := http.Client{}

	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return fmt.Errorf("Failed to execute follow request, %w", err)
	}

	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
	}

	io.Copy(os.Stdout, http_rsp.Body)

	return nil
}
