package profile

// Move this in to /ap

type Resource struct {
	Context           []string  `json:"@content"`
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"publicKey"`

	Following                 string `json:"following"`
	Followers                 string `json:"followers"`
	Name                      string `json:"name"`
	Summary                   string `json:"summary"`
	URL                       string `json:"url"`
	ManuallyApprovesFollowers bool   `json:"manuallyApprovesFollowers"`
	Discoverable              bool   `json:"discoverable"`
	Published                 string `json:"published"`
}

type Icon struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	URL       string `json:"url"`
}

type PublicKey struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}
