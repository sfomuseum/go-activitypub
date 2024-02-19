package activitypub

type Following struct {
	AccountId        int64  `json:"account_id"`
	FollowingAddress string `json:"following_address"`
	Created          int64  `json:"created"`
}
