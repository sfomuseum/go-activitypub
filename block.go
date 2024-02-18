package activitypub

type Block struct {
	Id           int64  `json:"id"`
	AccountId    int64  `json:"account_id"`
	Name         string `json:"name"`
	Host         string `json:"host"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"lastmodified"`
}
