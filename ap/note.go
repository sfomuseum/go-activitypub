package ap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Note struct {
	// Type is the type of the note (aka "Note").
	Type string `json:"type"`
	// Id is the unique identifier for the note.
	Id string `json:"id"`
	// The URI of the actor that the note is attributed to.
	AttributedTo string `json:"attributedTo"`
	// ...
	InReplyTo string `json:"inReplyTo,omitempty"`
	// Zero or more tags associated with the note.
	Tags []*Tag `json:"tag,omitempty"`
	// To is the list of URIs the activity should be delivered to.
	To []string `json:"to"`
	// CC is the list of URIs the activity should be copied to.
	Cc []string `json:"cc,omitempty"`
	// The body of the note.
	Content string `json:"content"`
	// The permanent URL of the post.
	URL string `json:"url"`
	// The RFC3339 date that the activity was published.
	Published string `json:"published"`
	// Zero or more attachments to include with the note.
	Attachments []*Attachment `json:"attachment,omitempty"`
}

// Retrieve note fetches and unmarshals the "application/activity+json" representation of 'uri'.
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
