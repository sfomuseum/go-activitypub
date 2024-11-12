package activitypub

// This probably just needs to be renamed as "tag" and associated with an
// activitypub.ActivityTypeStatus and activitypub.ActivityTypeId but right
// now it is tightly coupled with "posts" (or "notes")

type PostTag struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Href      string `json:"href"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Created   int64  `json:"created"`
}
