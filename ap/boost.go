package ap

import (
	"context"
	"fmt"
)

/*

https://boyter.org/posts/activitypub-announce-post/

In the case of follower instances, they receive the announce, then fetch the content of the announce using the object field.

*/

/*

{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://mastodon.social/users/aaronofsfo/statuses/112361739040787763/activity",
  "type": "Announce",
  "actor": "https://mastodon.social/users/aaronofsfo",
  "to": [
    "https://www.w3.org/ns/activitystreams#Public"
  ],
  "cc": [
    "https://collection.sfomuseum.org/ap/onview",
    "https://mastodon.social/users/aaronofsfo/followers"
  ],
  "object": "https://collection.sfomuseum.org/ap/@onview/posts/1785022876366147584"
}

*/

func NewBoostActivity(ctx context.Context, from string, to string, object interface{}) (*Activity, error) {
	return NewAnnounceActivity(ctx, from, to, object)
}

func NewAnnounceActivity(ctx context.Context, from string, to string, object interface{}) (*Activity, error) {

	activity := &Activity{
		Context: ACTIVITYSTREAMS_CONTEXT,
		Type:    "Announce",
		Actor:   from,
		To: []string{
			fmt.Sprintf("%s#Public", ACTIVITYSTREAMS_CONTEXT),
		},
		Cc: []string{
			//followers...
			to,
		},
		Object: object,
	}

	return activity, nil
}
