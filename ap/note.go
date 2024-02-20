package ap

type Note struct {
	Type         string `json:"type"`
	Id           string `json:"id"`
	AttributedTo string `json:"attributedTo"`
	To           string `json:"to"`
	Content      string `json:"content"`
	URL          string `json:"url"`
	Published    string `json:"published"`
}
