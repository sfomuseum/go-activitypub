package ap

type Actor struct {
	Context           []string  `json:"@content"`
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferredUsername"`
	Inbox             string    `json:"inbox"`
	PublicKey         PublicKey `json:"publicKey"`
}

type Activity struct {
	Context string `json:"@context"`
	Id      string `json:"id"`
	Type    string `json:"type"`
	Actor   string `json:"actor"`
	Object  string `json:"object"`
}

type PublicKey struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}
