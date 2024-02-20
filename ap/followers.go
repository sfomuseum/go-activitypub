package ap

type Followers struct {
	Context    string `json:"@context"`
	Id         string `json:"id"`
	Type       string `json:"type"`
	TotalItems uint32 `json:"totalItems"`
	First      string `json:"first"`
}
