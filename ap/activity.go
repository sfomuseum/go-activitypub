package ap

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	// "net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/sfomuseum/iso8601duration"
)

// https://www.w3.org/TR/activitystreams-vocabulary/

// Activity is a struct encapsulating an ActivityPub activity.
type Activity struct {
	// Context needs to be a "whatever" (interface{}) because ActivityPub (JSON-LD)
	// mixes and matches string URIs, arbritrary data structures and arrays of string
	// URIs and arbritrary data structures in @context...
	Context interface{} `json:"@context,omitempty"`
	// Id is the unique identifier for the activity.
	Id string `json:"id"`
	// Type is the name of the activity being performed.
	Type string `json:"type"`
	// Actor is the URI of the person (actor) performing the activity. Note: This is a fully-qualified "profile" URI and not a "@user@host" address.
	Actor string `json:"actor"`
	// To is the list of URIs the activity should be delivered to. Note: The fact that this can be something like [ "https://www.w3.org/ns/activitystreams#Public" ] makes me wonder what the point of this property is (since the relevant issue is the inbox that the encoded activity is delivered to).
	To []string `json:"to,omitempty"`
	// CC is the list of URIs the activity should be copied to. Note: It's not clear what the point of this unless the purpose of this property (and the "To" property) is for an activity to double as a complete record, inclusive of every address it should be delivered to, that can be scheduled for asynchronous delivery.
	Cc []string `json:"cc,omitempty"`
	// Audience limits visibility to just the specified users.
	Audience string `json:"audience,omitempty"`
	// Object is body of the activity itself.
	Object interface{} `json:"object,omitempty"`
	// The RFC3339 date that the activity was published.
	Published string `json:"published,omitempty"`
}

type PostToInboxOptions struct {
	// KeyId is a pointer (URI) to the actor/profile page where the public key of the actor posting the activity can be retrieved.
	KeyId string
	// The private key of the actor posting the activity used to sign the message.
	PrivateKey *rsa.PrivateKey
	// The URL of the inbox where the Activity should be posted.
	Inbox string
	// Log POST requests before they are sent using the default [log/slog] Logger. Note that this will
	// include the HTTP signature sent with the request so you should apply all the necessary care that
	// these values are logged somewhere you don't want unauthorized eyes to see the.
	LogRequest bool
	// Log the body of the POST response if it contains a status code that is not 200-202 or 204 using
	// the default [log/slog] Logger
	LogResponseOnError bool
}

// PostToInbox delivers an Activity message to a specific inbox.
func (activity *Activity) PostToInbox(ctx context.Context, key_id string, private_key *rsa.PrivateKey, inbox_uri string) error {

	logger := slog.Default()

	logger = logger.With("activity id", activity.Id)
	logger = logger.With("to", activity.To)
	logger = logger.With("inbox", inbox_uri)

	logger.Info("Post activity to inbox")

	enc_req, err := json.Marshal(activity)

	if err != nil {
		logger.Error("Failed to marshal activity", "error", err)
		return fmt.Errorf("Failed to marshal follow activity request, %w", err)
	}

	http_req, err := http.NewRequestWithContext(ctx, "POST", inbox_uri, bytes.NewBuffer(enc_req))

	if err != nil {
		logger.Error("Failed to create new request for activity", "error", err)
		return fmt.Errorf("Failed to create new request to %s, %w", inbox_uri, err)
	}

	now := time.Now()

	http_req.Header.Set("Content-Type", ACTIVITY_LD_CONTENT_TYPE)

	// RFC 2612 dates are required
	http_req.Header.Set("Date", now.Format(http.TimeFormat))

	// START OF this is necessary for HTTP signature hoohah...
	inbox_u, err := url.Parse(inbox_uri)

	if err != nil {
		logger.Error("Failed to parse inbox URI", "error", err)
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
		logger.Error("Failed to create new HTTP signer", "error", err)
		return fmt.Errorf("Failed to create new signer, %w", err)
	}

	err = signer.SignRequest(private_key, key_id, http_req, enc_req)

	if err != nil {
		logger.Error("Failed to sign request", "error", err)
		return fmt.Errorf("Failed to sign request, %w", err)
	}

	// https://pkg.go.dev/net/http/httputil#DumpRequest

	/*
		if opts.LogRequest {

			dump, err := httputil.DumpRequest(http_req, true)

			if err != nil {
				return fmt.Errorf("Failed to dump request, %w", err)
			}

			logger.Debug("REQUEST", "body", string(dump))
		}
	*/

	http_cl := http.Client{}
	http_rsp, err := http_cl.Do(http_req)

	if err != nil {
		return fmt.Errorf("Failed to execute post to inbox request, %w", err)
	}

	defer http_rsp.Body.Close()

	logger.Info("Response", "code", http_rsp.StatusCode, "content-type", http_rsp.Header.Get("Content-Type"))

	// https://www.w3.org/wiki/ActivityPub/Primer/HTTP_status_codes_for_delivery

	switch http_rsp.StatusCode {
	// HTTP 200, 201, 202, 204
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	default:

		/*
			if opts.LogResponseOnError {

				body, read_err := io.ReadAll(http_rsp.Body)

				if read_err != nil {
					return fmt.Errorf("Follow request failed %d, %s; read body also failed, %w", http_rsp.StatusCode, http_rsp.Status, read_err)
				}

				logger.Debug("ERROR", "body", string(body))
			}
		*/
	}

	return fmt.Errorf("Follow request failed %d, %s", http_rsp.StatusCode, http_rsp.Status)
}
