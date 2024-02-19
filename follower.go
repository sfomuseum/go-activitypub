package activitypub

type Follower struct {
	AccountId       int64  `json:"account_id"`
	FollowerAddress string `json:"follower_address"`
	Created         int64  `json:"created"`
}
