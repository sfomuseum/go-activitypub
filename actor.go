package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/sfomuseum/go-activitypub/ap"
	"github.com/sfomuseum/go-activitypub/webfinger"
)

func RetrieveActor(ctx context.Context, id string, insecure bool) (*ap.Actor, error) {

	actor_id, actor_hostname, err := ParseAddress(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse ID, %w", err)
	}

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

	slog.Debug("Webfinger URL for resource", "resource", actor_id, "url", webfinger_url)

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

	slog.Debug("Profile page for actor", "actor", actor_id, "url", profile_url)

	profile_req, err := http.NewRequestWithContext(ctx, "GET", profile_url, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create profile request, %w", err)
	}

	profile_req.Header.Set("Accept", ap.ACTIVITYSTREAMS_ACCEPT_HEADER)

	cl := &http.Client{}

	profile_rsp, err := cl.Do(profile_req)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve profile URL (%s), %w", profile_url, err)
	}

	defer profile_rsp.Body.Close()

	if profile_rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Remote endpoint did not return successfully %d, %s", profile_rsp.StatusCode, profile_rsp.Status)
	}

	var actor *ap.Actor

	dec = json.NewDecoder(profile_rsp.Body)
	err = dec.Decode(&actor)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode profile response, %w", err)
	}

	// START OF this is what Mastodon does but is it really necessary?
	/*

		u, _ := url.Parse(profile_url)

		webfinger_acct2 := fmt.Sprintf("acct:%s@%s", actor.PreferredUsername, u.Host)

		webfinger_q2 := &url.Values{}
		webfinger_q2.Set("resource", webfinger_acct2)

		webfinger_u2 := &url.URL{}
		webfinger_u2.Scheme = webfinger_scheme
		webfinger_u2.Host = u.Host
		webfinger_u2.Path = webfinger.Endpoint
		webfinger_u2.RawQuery = webfinger_q.Encode()

		webfinger_url2 := webfinger_u2.String()

		slog.Debug("Webfinger URL for resource", "resource", webfinger_acct2, "url", webfinger_url2)

		webfinger_rsp2, err := http.Get(webfinger_url2)

		if err != nil {
			return nil, fmt.Errorf("Failed to perform webfinger (%s) for actor, %w", webfinger_url2, err)
		}

		defer webfinger_rsp2.Body.Close()

		if webfinger_rsp2.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Remote endpoint did not return successfully %d, %s", webfinger_rsp2.StatusCode, webfinger_rsp.Status)
		}

		var webfinger_resource2 *webfinger.Resource

		dec = json.NewDecoder(webfinger_rsp2.Body)
		err = dec.Decode(&webfinger_resource2)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode webfinger resource, %w", err)
		}

		if webfinger_acct2 != webfinger_resource2.Subject {
			return nil, fmt.Errorf("Second webfinger request yields a different subject. Expected '%s' but got '%s'", webfinger_acct2, webfinger_resource2.Subject)
		}

	*/
	// END OF this is what Mastodon does but is it really necessary?

	return actor, nil
}
