package ap

import (
	"net/http"
	"testing"
)

func TestParseAddressFromRequest(t *testing.T) {

	tests := map[string][2]string{
		"bob":                  [2]string{"bob", ""},
		"@bob":                 [2]string{"bob", ""},
		"@bob@localhost":       [2]string{"bob", "localhost"},
		"@bob@bob.com":         [2]string{"bob", "bob.com"},
		"acct:@bob@bob.com":    [2]string{"bob", "bob.com"},
		"acct:alice@bob.com":   [2]string{"alice", "bob.com"},
		"@bob@localhost:8080":  [2]string{"bob", "localhost:8080"},
		"@doug@localhost:8080": [2]string{"doug", "localhost:8080"},
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
		"bob":                [2]string{"bob", ""},
		"@bob":               [2]string{"bob", ""},
		"@bob@localhost":     [2]string{"bob", "localhost"},
		"@bob@bob.com":       [2]string{"bob", "bob.com"},
		"acct:@bob@bob.com":  [2]string{"bob", "bob.com"},
		"acct:alice@bob.com": [2]string{"alice", "bob.com"},
		// "https://mastodon.social/users/aaronofsfo": [2]string{ "aaronofsfo", "mastodon.social"},
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

func TestParseAddressesFromString(t *testing.T) {

	tests := map[string]int{
		"hello @bob@example.com":                  1,
		"hello @bob@example.com pass the mustard": 1,
		"hello @bob@example.com pass the mustard to @alice@mustard.com and doug@localhost":                                                   2,
		"hello @bob@example.com pass the mustard to @alice@mustard.com and doug@localhost before sending it @doug@bob.com":                   3,
		"hello @bob@example.com pass the mustard to @alice@mustard.com and doug@localhost before sending it @doug@bob.com and max@gmail.com": 3,
		"test mentioning @doug@localhost yo":                                                    1,
		"test mentioning @doug@localhost:8080 and @bob@localhost:8080 yo":                       2,
		"test mentioning @doug@localhost:8080 and @bob@localhost:8080 yo @alice@alice.com":      3,
		`<a href="http://localhost:8080/users/doug/">@doug@localhost</a> nearby`:                1,
		`<a href="http://localhost:8080/users/doug/">@doug@localhost</a> nearby @alice@bob.com`: 2,
	}

	for str, expected_count := range tests {

		addrs, err := ParseAddressesFromString(str)

		if err != nil {
			t.Fatalf("Failed to parse addresses from string (%s), %v", str, err)
		}

		if len(addrs) != expected_count {
			t.Fatalf("Expected %d address for string (%s) but got %d", expected_count, str, len(addrs))
		}
	}
}
