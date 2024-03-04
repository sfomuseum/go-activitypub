package activitypub

type Boost struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	PostId    int64  `json:"post_id"`
	Creator   string `json:"creator"`
	Created   int64  `json:"created"`
}
