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
	Context string      `json:"@context"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`
	Actor   string      `json:"actor"`
	To      []string    `json:"to,omitempty"`
	Object  interface{} `json:"object"`
}

type Note struct {
	Type         string      `json:"type"`
	Id           string      `json:"id"`
	AttributedTo string      `json:"attributedTo"`
	To           string      `json:"to"`
	Content      interface{} `json:"content"`
	URL          string      `json:"url"`
	Published    string      `json:"published"`
}

type PublicKey struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
	PEM   string `json:"publicKeyPem"`
}
