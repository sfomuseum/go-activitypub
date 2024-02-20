package activitypub

import (
	"net/http"
	"testing"
)

func TestParseAddressFromRequest(t *testing.T) {

	tests := map[string][2]string{
		"bob":            [2]string{"bob", ""},
		"@bob":           [2]string{"bob", ""},
		"@bob@localhost": [2]string{"bob", "localhost"},
		"@bob@bob.com":   [2]string{"bob", "bob.com"},
	}

	for addr, expected := range tests {

		path := "/{resource}/inbox"

		req, err := http.NewRequest("GET", path, nil)

		if err != nil {
			t.Fatalf("Failed to create new request for %s, %v", path, err)
		}

		req.SetPathValue("resource", addr)

		name, host, err := ParseAddressFromRequest(req)

		if err != nil {
			t.Fatalf("Failed to parse address '%s', %v", path, err)
		}

		if name != expected[0] {
			t.Fatalf("Unexpected host for '%s'. Expected '%s' but got '%s'", path, expected[0], name)
		}

		if host != expected[1] {
			t.Fatalf("Unexpected host for '%s'. Expected '%s' but got '%s'", path, expected[1], name)
		}

	}

}

func TestParseAddress(t *testing.T) {

	tests := map[string][2]string{
		"bob":            [2]string{"bob", ""},
		"@bob":           [2]string{"bob", ""},
		"@bob@localhost": [2]string{"bob", "localhost"},
		"@bob@bob.com":   [2]string{"bob", "bob.com"},
	}

	for addr, expected := range tests {

		name, host, err := ParseAddress(addr)

		if err != nil {
			t.Fatalf("Failed to parse address '%s', %v", addr, err)
		}

		if name != expected[0] {
			t.Fatalf("Unexpected host for '%s'. Expected '%s' but got '%s'", addr, expected[0], name)
		}

		if host != expected[1] {
			t.Fatalf("Unexpected host for '%s'. Expected '%s' but got '%s'", addr, expected[1], name)
		}

	}
}
