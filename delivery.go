package activitypub

type Delivery struct {
	Id            int64  `json:"id"`
	ActivityId    int64  `json:"activity_id"`
	ActivityPubId string `json:"activitypub_id"`
	AccountId     int64  `json:"account_id"`
	Recipient     string `json:"recipient"`
	Inbox         string `json:"inbox"`
	Created       int64  `json:"created"`
	Completed     int64  `json:"completed"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
}
