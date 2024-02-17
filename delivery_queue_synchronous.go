package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

type SynchronousDeliveryQueue struct {
	DeliveryQueue
}

func init() {
	ctx := context.Background()
	RegisterDeliveryQueue(ctx, "synchronous", NewSynchronousDeliveryQueue)
}

func NewSynchronousDeliveryQueue(ctx context.Context, uri string) (DeliveryQueue, error) {
	q := &SynchronousDeliveryQueue{}
	return q, nil
}

func (q *SynchronousDeliveryQueue) DeliverPost(ctx context.Context, p *Post, follower_id string) error {

	to := []string{
		follower_id,
	}

	create_activity, err := p.AsCreateActivity(ctx, to)

	enc_activity, err := json.Marshal(create_activity)

	if err != nil {
		return fmt.Errorf("Failed to encode activity, %w", err)
	}

	slog.Info("POST", "activity", string(enc_activity))

	// follower_ids should be in the form of @USER@HOSTNAME

	follower_id = "bob@localhost:8080"
	actor, err := RetrieveActor(ctx, follower_id)

	if err != nil {
		return fmt.Errorf("Failed to retrieve actor resource for %s, %w", follower_id, err)
	}

	actor_inbox := actor.Inbox
	slog.Info("ACTOR", "inbox", actor_inbox)

	/*
		// START OF...

		inbox_uri := actor_inbox	// string
		enc_req := enc_activity		// []byte
		private_key := "FIX ME"		// *rsa.PrivateKey
		key_id := "FIX ME"	// string, but...what?

		http_req, err := http.NewRequestWithContext(ctx, "POST", inbox_uri, bytes.NewBuffer(enc_req))

		if err != nil {
			return fmt.Errorf("Failed to create new request to %s, %w", opts.Inbox, err)
		}

		now := time.Now()
		http_req.Header.Set("Date", now.Format(time.RFC3339))

		key_id := follower

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
	*/

	return nil
}
