package activitypub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/uris"
	"github.com/sfomuseum/iso8601duration"
)

type PostToAccountOptions struct {
	From     *Account
	To       string
	Activity *ap.Activity
	URIs     *uris.URIs
}

type PostToInboxOptions struct {
	// The `Account` instance of the actor sending the Activity.
	From *Account
	// The URL of the inbox where the Activity should be posted.
	Inbox string
	// The `Activity` instance being posted to the inbox.
	Activity *ap.Activity
	URIs     *uris.URIs
	// Log POST requests before they are sent using the default [log/slog] Logger. Note that this will
	// include the HTTP signature sent with the request so you should apply all the necessary care that
	// these values are logged somewhere you don't want unauthorized eyes to see the.
	LogRequest bool
	// Log the body of the POST response if it contains a status code that is not 200-202 or 204 using
	// the default [log/slog] Logger
	LogResponseOnError bool
}

func PostToAccount(ctx context.Context, opts *PostToAccountOptions) (string, error) {

	actor, err := RetrieveActor(ctx, opts.To, opts.URIs.Insecure)

	if err != nil {
		return "", fmt.Errorf("Failed to retrieve actor for %s, %w", opts.To, err)
	}

	inbox_opts := &PostToInboxOptions{
		From:     opts.From,
		Inbox:    actor.Inbox,
		Activity: opts.Activity,
		URIs:     opts.URIs,
	}

	return actor.Inbox, PostToInbox(ctx, inbox_opts)
}

// PostToInbox delivers an Activity message to a specific inbox.
func PostToInbox(ctx context.Context, opts *PostToInboxOptions) error {

	slog.Debug("Post to inbox", "inbox", opts.Inbox)

	enc_req, err := json.Marshal(opts.Activity)

	if err != nil {
		return fmt.Errorf("Failed to marshal follow activity request, %w", err)
	}

	http_req, err := http.NewRequestWithContext(ctx, "POST", opts.Inbox, bytes.NewBuffer(enc_req))

	if err != nil {
		return fmt.Errorf("Failed to create new request to %s, %w", opts.Inbox, err)
	}

	now := time.Now()

	http_req.Header.Set("Content-Type", ap.ACTIVITY_LD_CONTENT_TYPE)

	// RFC 2612 dates are required
	http_req.Header.Set("Date", now.Format(http.TimeFormat))

	// START OF this is necessary for HTTP signature hoohah...
	inbox_u, err := url.Parse(opts.Inbox)

	if err != nil {
		return fmt.Errorf("Failed to parse inbox URL, %w", err)
	}

	http_req.Header.Set("Host", inbox_u.Host)
	// END OF this is necessary for HTTP signature hoohah...

	// Should this value be configurable?
	str_ttl := "PT5M"

	d, err := duration.FromString(str_ttl)

	if err != nil {
		return fmt.Errorf("Failed to derive duration, %w", err)
	}

	ttl := int64(d.ToDuration().Seconds())

	// Created and Expires headers are important for posting to Mastodon
	// https://github.com/mastodon/mastodon/blob/main/app/controllers/concerns/signature_verification.rb#L183

	created := now.Unix()
	expires := created + ttl

	http_req.Header.Set("Created", strconv.FormatInt(created, 10))
	http_req.Header.Set("Expires", strconv.FormatInt(expires, 10))

	// Note that "key_id" here means a pointer to the actor/profile page where the public key
	// for the follower can be retrieved

	profile_url := opts.From.AccountURL(ctx, opts.URIs)

	key_id := profile_url.String()
	slog.Debug("Post to inbox", "inbox", opts.Inbox, "key_id", key_id)

	private_key, err := opts.From.PrivateKeyRSA(ctx)

	if err != nil {
		return fmt.Errorf("Failed to derive private key for from account, %w", err)
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
		"Digest",
		// See the way this is "(created)" and not "Created". That's a go-fed/httpsig thing... or maybe it's a spec thing?
		// https://github.com/go-fed/httpsig/blob/master/signing.go#L220-L229
		"(created)",
		"(expires)",
	}

	signer, _, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, ttl)

	if err != nil {
		return fmt.Errorf("Failed to create new signer, %w", err)
	}

	err = signer.SignRequest(private_key, key_id, http_req, enc_req)

	if err != nil {
		return fmt.Errorf("Failed to sign request, %w", err)
	}

	// https://pkg.go.dev/net/http/httputil#DumpRequest

	if opts.LogRequest {

		dump, err := httputil.DumpRequest(http_req, true)

		if err != nil {
			return fmt.Errorf("Failed to dump request, %w", err)
		}

		slog.Debug("REQUEST", "body", string(dump))
	}

	http_cl := http.Client{}
	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return fmt.Errorf("Failed to execute post to inbox request, %w", err)
	}

	defer http_rsp.Body.Close()

	slog.Debug("Response", "inbox", opts.Inbox, "code", http_rsp.StatusCode, "content-type", http_rsp.Header.Get("Content-Type"))

	// https://www.w3.org/wiki/ActivityPub/Primer/HTTP_status_codes_for_delivery

	switch http_rsp.StatusCode {
	// HTTP 200, 201, 202, 204
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	default:

		if opts.LogResponseOnError {

			body, read_err := io.ReadAll(http_rsp.Body)

			if read_err != nil {
				return fmt.Errorf("Follow request failed %d, %s; read body also failed, %w", http_rsp.StatusCode, http_rsp.Status, read_err)
			}

			slog.Debug("ERROR", "body", string(body))
		}
	}

	return fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
}
