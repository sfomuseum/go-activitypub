package activitypub

type Delivery struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	AccountId int64  `json:"account_id"`
	Recipient string `json:"recipient"`
	Created   int64  `json:"created"`
	Completed int64  `json:"completed"`
	Status    int    `json:"status"`
	Error     string `json:"error,omitempty"`
}
