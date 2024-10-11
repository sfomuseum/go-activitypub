package ap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Note struct {
	Type         string   `json:"type"`
	Id           string   `json:"id"`
	AttributedTo string   `json:"attributedTo"`
	InReplyTo    string   `json:"inReplyTo,omitempty"`
	Tags         []*Tag   `json:"tag,omitempty"`
	To           []string `json:"to"`
	Cc           []string `json:"cc"`
	Content      string   `json:"content"`
	URL          string   `json:"url"`
	Published    string   `json:"published"`
}

func RetrieveNote(ctx context.Context, uri string) (*Note, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create request, %w", err)
	}

	req.Header.Set("Accept", ACTIVITY_CONTENT_TYPE)

	cl := &http.Client{}

	rsp, err := cl.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to get note, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Note returned unexpected status, %d (%s)", rsp.StatusCode, rsp.Status)
	}

	var n *Note

	dec := json.NewDecoder(rsp.Body)
	err = dec.Decode(&n)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode note, %w", err)
	}

	if n.Type != "Note" {
		return nil, fmt.Errorf("Unexpected type, %s", n.Type)
	}

	return n, nil
}
