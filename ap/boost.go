package ap

import (
	"fmt"
)

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

func NewBoostActivity(from string, to string, object interface{}) (*Activity, error) {
	return NewAnnounceActivity(from, to, object)
}

func NewAnnounceActivity(from string, to string, object interface{}) (*Activity, error) {

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
