package ap

type Attachment struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	URL       string `json:"url"`
}
