package follow

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/iso8601duration"
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

	acct, err := db.GetAccount(ctx, opts.AccountId)

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

	if opts.Undo {
		follow_req.Type = "Undo"
	}

	enc_req, err := json.Marshal(follow_req)

	if err != nil {
		return fmt.Errorf("Failed to marshal follow activity request, %w", err)
	}

	// START OF make me common code...

	http_req, err := http.NewRequestWithContext(ctx, "POST", opts.Inbox, bytes.NewBuffer(enc_req))

	if err != nil {
		return fmt.Errorf("Failed to create new request to %s, %w", opts.Inbox, err)
	}

	now := time.Now()
	http_req.Header.Set("Date", now.Format(time.RFC3339))

	key_id := follower

	private_key, err := acct.PrivateKeyRSA(ctx)

	if err != nil {
		return fmt.Errorf("Failed to derive private key for account, %w", err)
	}

	// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-1.1
	// https://pkg.go.dev/github.com/go-fed/httpsig

	prefs := []httpsig.Algorithm{httpsig.RSA_SHA512, httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256

	headersToSign := []string{
		httpsig.RequestTarget,
		"Date",
		"Digest",
	}

	str_ttl := "PT1M"

	d, err := duration.FromString(str_ttl)

	if err != nil {
		return fmt.Errorf("Failed to derive duration, %w", err)
	}

	ttl := int64(d.ToDuration().Seconds())

	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, ttl)

	if err != nil {
		return fmt.Errorf("Failed to create new signer, %w", err)
	}

	err = signer.SignRequest(private_key, key_id, http_req, enc_req)

	if err != nil {
		return fmt.Errorf("Failed to sign request, %w", err)
	}

	http_cl := http.Client{}

	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return fmt.Errorf("Failed to execute follow request, %w", err)
	}

	// logger.Info("Response", "code", http_rsp.StatusCode)

	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
	}

	var activity *ap.Activity

	dec := json.NewDecoder(http_rsp.Body)
	err = dec.Decode(&activity)

	if err != nil {
		return fmt.Errorf("Failed to decode response, %w", err)
	}

	if activity.Type != "Accept" {
		return fmt.Errorf("Unexpected activity type, %s", activity.Type)
	}

	// END OF make make common code

	// Check actor/object pairs here...

	if undo {
		logger.Info("Unfollowing successful")
	} else {
		logger.Info("Following successful")
	}

	return nil
}
