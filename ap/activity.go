package ap

type Activity struct {
	Context string      `json:"@context"`
	Id      string      `json:"id"`
	Type    string      `json:"type"`
	Actor   string      `json:"actor"`
	To      []string    `json:"to,omitempty"`
	Object  interface{} `json:"object"`
}
