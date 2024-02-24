package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/iso8601duration"
)

type PostToAccountOptions struct {
	From    *Account
	To      string
	Message interface{}
	URIs    *uris.URIs
}

type PostToInboxOptions struct {
	From    *Account
	Inbox   string
	Message interface{}
	URIs    *uris.URIs
}

func PostToAccount(ctx context.Context, opts *PostToAccountOptions) (*ap.Activity, error) {

	actor, err := RetrieveActor(ctx, opts.To, opts.URIs.Insecure)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve actor for %s, %w", opts.To, err)
	}

	inbox_opts := &PostToInboxOptions{
		From:    opts.From,
		Inbox:   actor.Inbox,
		Message: opts.Message,
		URIs:    opts.URIs,
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

	http_req.Header.Set("Content-Type", ap.ACTIVITY_LD_CONTENT_TYPE)

	http_req.Header.Set("Date", now.Format(time.RFC3339))
	http_req.Header.Set("Host", opts.URIs.Hostname)
	http_req.Host = opts.URIs.Hostname

	// Note that "key_id" here means a pointer to the actor/profile page where the public key
	// for the follower can be retrieved

	profile_url := opts.From.AccountURL(ctx, opts.URIs)

	key_id := profile_url.String()
	slog.Debug("Post to inbox", "inbox", opts.Inbox, "key_id", key_id)

	private_key, err := opts.From.PrivateKeyRSA(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive private key for from account, %w", err)
	}

	// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-1.1
	// https://pkg.go.dev/github.com/go-fed/httpsig

	// prefs := []httpsig.Algorithm{httpsig.RSA_SHA512, httpsig.RSA_SHA256}

	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256

	headersToSign := []string{
		httpsig.RequestTarget,
		"Host",
		"Date",
		// "Digest",
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

	// START OF debugging
	// https://pkg.go.dev/net/http/httputil#DumpRequest

	dump, err := httputil.DumpRequest(http_req, true)

	if err != nil {
		return nil, fmt.Errorf("Failed to dump request, %w", err)
	}

	slog.Debug("REQUEST", "body", string(dump))

	// END OF debugging

	http_cl := http.Client{}
	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute post to inbox request, %w", err)
	}

	defer http_rsp.Body.Close()

	slog.Debug("Response", "inbox", opts.Inbox, "code", http_rsp.StatusCode, "content-type", http_rsp.Header.Get("Content-Type"))

	// https://www.w3.org/wiki/ActivityPub/Primer/HTTP_status_codes_for_delivery

	switch http_rsp.StatusCode {
	// HTTP 201, 202, 204
	case http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil, nil
	case http.StatusOK:

		/*
			var activity *ap.Activity

			activity_r := DefaultLimitedReader(http_rsp.Body)

			dec := json.NewDecoder(activity_r)
			err = dec.Decode(&activity)

			if err != nil {
				return nil, fmt.Errorf("Failed to decode response, %w", err)
			}

			if activity.Type != "Accept" {
				return nil, fmt.Errorf("Unexpected activity type, %s", activity.Type)
			}

			return activity, nil
		*/

		return nil, nil
	default:
		//
	}

	return nil, fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
}
