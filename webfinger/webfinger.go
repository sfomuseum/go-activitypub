package webfinger

// https://www.rfc-editor.org/rfc/rfc7033 (webfinger)
// https://www.rfc-editor.org/rfc/rfc7565 (acct:)

const ContentType string = "application/jrd+json"

type Resource struct {
	Subject    string            `json:"subject,omitempty"`
	Aliases    []string          `json:"aliases,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Links      []Link            `json:"links,omitempty"`
}

type Link struct {
	HRef       string             `json:"href"`
	Type       string             `json:"type,omitempty"`
	Rel        string             `json:"rel"`
	Properties map[string]*string `json:"properties,omitempty"`
	Titles     map[string]string  `json:"titles,omitempty"`
}
