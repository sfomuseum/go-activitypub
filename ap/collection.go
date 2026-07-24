package ap

type OrderedCollection struct {
	// It has to be an interface because JSON-LD... thanks, JSON-LD...
	Context     []any  `json:"@context"`
	Id          string `json:"id"`
	Summary     string `json:"summary,omitempty"`
	Type        string `json:"type"`
	TotalItems  int    `json:"totalItems"`
	OrderedItem []*any `json:"orderedItems,omitempty"`
}
