package activitypub

type Property struct {
	Id        int64  `json:"id"`
	AccountId int64  `json:"account_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Created   int64  `json:"created"`
}
