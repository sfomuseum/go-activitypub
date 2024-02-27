package activitypub

import (
	"testing"
)

func TestAccountType(t *testing.T) {

	string_tests := map[AccountType]string{
		PersonType:  "Person",
		ServiceType: "Service",
	}

	for account_t, expected := range string_tests {

		if account_t.String() != expected {
			t.Fatalf("Failed to derive string for '%v', got '%s' but expected '%s'", account_t, account_t.String(), expected)
		}
	}

	numeric_tests := map[AccountType]uint32{
		PersonType:  1,
		ServiceType: 2,
	}

	for account_t, expected := range numeric_tests {

		if uint32(account_t) != expected {
			t.Fatalf("Failed to derive numeric value for '%v', got '%d' but expected '%d'", account_t, account_t, expected)
		}
	}

}
