package activitypub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sfomuseum/go-activitypub/ap"
)

func RetrieveActor(ctx context.Context, id string) (*ap.Actor, error) {

	actor_id, actor_hostname, err := ParseAccountURI(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse ID, %w", err)
	}

	webfinger_q := &url.Values{}
	webfinger_q.Set("resource", actor_id)

	webfinger_u := &url.URL{}
	webfinger_u.Scheme = "http" // https
	webfinger_u.Host = actor_hostname
	webfinger_u.Path = "/.webfinger" // well-known?
	webfinger_u.RawQuery = webfinger_q.Encode()

	webfinger_uri := webfinger_u.String()

	rsp, err := http.Get(webfinger_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to perform webfinger (%s) for actor, %w", webfinger_uri, err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Remote endpoint did not return successfully %d, %s", rsp.StatusCode, rsp.Status)
	}

	var actor *ap.Actor

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&actor)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode actor, %w", err)
	}

	return actor, nil
}
