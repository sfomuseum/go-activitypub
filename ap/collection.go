package ap

type OrderedCollection struct {
	Context     []string       `json:"@context"`
	Id          string         `json:"id"`
	Summary     string         `json:"summary,omitempty"`
	Type        string         `json:"type"`
	TotalItems  int            `json:"totalItems"`
	OrderedItem []*interface{} `json:"orderedItems,omitempty"`
}
