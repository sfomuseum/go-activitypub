package activitypub

type Alias struct {
	Name      string `json:"name"` // Unique primary key
	AccountId int64  `json:"id"`
	Created   int64  `json:"created"`
}

func (a *Alias) String() string {
	return a.Name
}
