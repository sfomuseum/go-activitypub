package ap

import (
	"context"
	"testing"
)

func TestRetrieveNote(t *testing.T) {

	tests := map[string]string{
		"https://collection.sfomuseum.org/ap/@N977JE/posts/1844576284093452288": "https://collection.sfomuseum.org/ap/@N977JE",
	}

	ctx := context.Background()

	for uri, author := range tests {

		n, err := RetrieveNote(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to fetch note for %s, %v", uri, err)
		}

		if n.AttributedTo != author {
			t.Fatalf("Invalid author. Expected '%s' but got '%s'", author, n.AttributedTo)
		}
	}
}
