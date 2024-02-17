package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/iso8601duration"
)

type PostToAccountOptions struct {
	From     *Account
	To       string
	Message  interface{}
	Hostname string
	URIs     *URIs
}

type PostToInboxOptions struct {
	From     *Account
	Inbox    string
	Message  interface{}
	Hostname string
	URIs     *URIs
}

func PostToAccount(ctx context.Context, opts *PostToAccountOptions) (*ap.Activity, error) {

	actor, err := RetrieveActor(ctx, opts.To)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve actor for %s, %w", opts.To, err)
	}

	inbox_opts := &PostToInboxOptions{
		From:     opts.From,
		Inbox:    actor.Inbox,
		Message:  opts.Message,
		Hostname: opts.Hostname,
		URIs:     opts.URIs,
	}

	return PostToInbox(ctx, inbox_opts)
}

func PostToInbox(ctx context.Context, opts *PostToInboxOptions) (*ap.Activity, error) {

	enc_req, err := json.Marshal(opts.Message)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal follow activity request, %w", err)
	}

	http_req, err := http.NewRequestWithContext(ctx, "POST", opts.Inbox, bytes.NewBuffer(enc_req))

	if err != nil {
		return nil, fmt.Errorf("Failed to create new request to %s, %w", opts.Inbox, err)
	}

	now := time.Now()
	http_req.Header.Set("Date", now.Format(time.RFC3339))

	// So "key_id" here means a pointer to the actor/profile page where the public key for the follower can be retrieved

	profile_url, err := opts.From.ProfileURL(ctx, opts.Hostname, opts.URIs)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive profile URL for follower, %w", err)
	}

	key_id := profile_url.String()
	slog.Info("FOLLOW", "key_id", key_id)

	private_key, err := opts.From.PrivateKeyRSA(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive private key for from account, %w", err)
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
		return nil, fmt.Errorf("Failed to derive duration, %w", err)
	}

	ttl := int64(d.ToDuration().Seconds())

	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, ttl)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new signer, %w", err)
	}

	err = signer.SignRequest(private_key, key_id, http_req, enc_req)

	if err != nil {
		return nil, fmt.Errorf("Failed to sign request, %w", err)
	}

	http_cl := http.Client{}

	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute follow request, %w", err)
	}

	// logger.Info("Response", "code", http_rsp.StatusCode)

	defer http_rsp.Body.Close()

	if http_rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
	}

	var activity *ap.Activity

	dec := json.NewDecoder(http_rsp.Body)
	err = dec.Decode(&activity)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode response, %w", err)
	}

	if activity.Type != "Accept" {
		return nil, fmt.Errorf("Unexpected activity type, %s", activity.Type)
	}

	return activity, nil
}
