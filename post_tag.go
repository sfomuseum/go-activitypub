package activitypub

type PostTag struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Href      string `json:"href"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Created   int64  `json:"created"`
}
