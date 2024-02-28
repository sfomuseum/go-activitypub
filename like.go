package activitypub

type Like struct {
	Id int64 `json:"id"`
	AccountId int64 `json:"account_id"`
	PostId int64 `json:"post_id"`
	Address string `json:"address"`
	Created int64 `json:"created"`
}
