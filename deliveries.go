package activitypub

type Delivery struct {
	// As in the activity ID
	Id        string `json:"id"`
	PostId    int64  `json:"post_id"`
	AccountId int64  `json:"account_id"`
	Recipient string `json:"recipient"`
	Created   int64  `json:"created"`
	Completed int64  `json:"completed"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}
