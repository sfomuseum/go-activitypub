package ap

import (
	"context"
	"testing"
)

func TestRetrieveActor(t *testing.T) {

	ctx := context.Background()

	tests := map[string]string{
		"N977JE@collection.sfomuseum.org":  "N977JE@collection.sfomuseum.org",
		"@N977JE@collection.sfomuseum.org": "N977JE@collection.sfomuseum.org",
	}

	for addr, expected_addr := range tests {

		a, err := RetrieveActor(ctx, addr, false)

		if err != nil {
			t.Fatalf("Failed to retrieve actor '%s', %v", addr, err)
		}

		addr2, err := a.Address()

		if err != nil {
			t.Fatalf("Failed to derive address for actor '%s', %v", addr, err)
		}

		if addr2 != expected_addr {
			t.Fatalf("Expected '%s' to be the same as '%s'", addr2, expected_addr)
		}
	}
}

func TestRetrieveActorWithProfileURL(t *testing.T) {

	ctx := context.Background()

	tests := map[string]string{
		"https://collection.sfomuseum.org/ap/N977JE": "N977JE@collection.sfomuseum.org",
	}

	for profile_url, expected_addr := range tests {

		a, err := RetrieveActorWithProfileURL(ctx, profile_url)

		if err != nil {
			t.Fatalf("Failed to retrieve actor '%s', %v", profile_url, err)
		}

		addr2, err := a.Address()

		if err != nil {
			t.Fatalf("Failed to derive address for actor '%s', %v", profile_url, err)
		}

		if addr2 != expected_addr {
			t.Fatalf("Expected '%s' to be the same as '%s'", addr2, expected_addr)
		}
	}
}
