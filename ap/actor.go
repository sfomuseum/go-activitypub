package ap

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/sfomuseum/go-activitypub/crypto"
	"github.com/sfomuseum/go-activitypub/webfinger"
)

// https://www.w3.org/TR/activitystreams-vocabulary/#actor-types

type Actor struct {
	// It has to be an interface because JSON-LD... thanks, JSON-LD...
	Context           []interface{} `json:"@context"`
	Id                string        `json:"id"`
	Type              string        `json:"type"`
	PreferredUsername string        `json:"preferredUsername"`
	Inbox             string        `json:"inbox"`
	Outbox            string        `json:"outbox"`
	PublicKey         PublicKey     `json:"publicKey"`

	Following string `json:"following,omitempty"`
	Followers string `json:"followers,omitempty"`
	Name      string `json:"name,omitempty"`
	Summary   string `json:"summary,omitempty"`
	URL       string `json:"url,omitempty"`
	// Don't omitempty because if you do then false values are omitted
	// ManuallyApprovesFollowers bool   `json:"manuallyApprovesFollowers"`
	Discoverable bool          `json:"discoverable,omitempty"`
	Published    string        `json:"published,omitempty"`
	Icon         Icon          `json:"icon,omitempty"`
	Attachments  []*Attachment `json:"attachment,omitempty"` // Is this just a Mastodon-ism?
}

func (a *Actor) Address() (string, error) {

	u, err := url.Parse(a.Inbox)

	if err != nil {
		return "", fmt.Errorf("Failed to parse inbox URL, %w", err)
	}

	return fmt.Sprintf("%s@%s", a.PreferredUsername, u.Host), nil
}

// Returns the `rsa.PublicKey` instance for 'a'.
func (a *Actor) PublicKeyRSA(ctx context.Context) (*rsa.PublicKey, error) {

	public_key_str := a.PublicKey.PEM

	if public_key_str == "" {
		return nil, fmt.Errorf("Actor missing public key")
	}

	public_key, err := crypto.RSAPublicKeyFromPEM(public_key_str)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse PEM block containing public key, %w", err)

	}

	return public_key, nil
}

func RetrieveActor(ctx context.Context, id string, insecure bool) (*Actor, error) {

	logger := slog.Default()
	logger = logger.With("actor", id)

	logger.Debug("Retrieve actor")

	actor_id, actor_hostname, err := ParseAddress(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse ID, %w", err)
	}

	logger = logger.With("actor id", actor_id)

	webfinger_scheme := "https"

	if insecure {
		webfinger_scheme = "http"
	}

	webfinger_acct := fmt.Sprintf("acct:%s@%s", actor_id, actor_hostname)

	webfinger_q := &url.Values{}
	webfinger_q.Set("resource", webfinger_acct)

	webfinger_u := &url.URL{}
	webfinger_u.Scheme = webfinger_scheme
	webfinger_u.Host = actor_hostname
	webfinger_u.Path = webfinger.Endpoint
	webfinger_u.RawQuery = webfinger_q.Encode()

	webfinger_url := webfinger_u.String()

	logger.Debug("Webfinger URL for resource", "url", webfinger_url)

	logger = logger.With("webfinger", webfinger_url)

	webfinger_rsp, err := http.Get(webfinger_url)

	if err != nil {
		return nil, fmt.Errorf("Failed to perform webfinger (%s) for actor, %w", webfinger_url, err)
	}

	defer webfinger_rsp.Body.Close()

	if webfinger_rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Remote endpoint did not return successfully %d, %s", webfinger_rsp.StatusCode, webfinger_rsp.Status)
	}

	var webfinger_resource *webfinger.Resource

	dec := json.NewDecoder(webfinger_rsp.Body)
	err = dec.Decode(&webfinger_resource)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode webfinger resource, %w", err)
	}

	var profile_url string

	for _, l := range webfinger_resource.Links {

		if l.Rel == "self" && l.Type == "application/activity+json" {
			profile_url = l.HRef
			break
		}
	}

	if profile_url == "" {
		return nil, fmt.Errorf("Failed to derive profile URL from webfinger resource")
	}

	return RetrieveActorWithProfileURL(ctx, profile_url)
}

func RetrieveActorWithProfileURL(ctx context.Context, profile_url string) (*Actor, error) {

	// slog.Info(profile_url)

	logger := slog.Default()
	logger = logger.With("profile url", profile_url)

	logger.Debug("Retrieve actor with profile")

	profile_req, err := http.NewRequestWithContext(ctx, "GET", profile_url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create profile request, %w", err)
	}

	profile_req.Header.Set("Accept", ACTIVITYSTREAMS_ACCEPT_HEADER)

	cl := &http.Client{}

	profile_rsp, err := cl.Do(profile_req)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve profile URL (%s), %w", profile_url, err)
	}

	defer profile_rsp.Body.Close()

	if profile_rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Remote endpoint did not return successfully %d, %s", profile_rsp.StatusCode, profile_rsp.Status)
	}

	var actor *Actor

	dec := json.NewDecoder(profile_rsp.Body)
	err = dec.Decode(&actor)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode profile response, %w", err)
	}

	return actor, nil
}
