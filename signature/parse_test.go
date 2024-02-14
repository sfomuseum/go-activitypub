package signature

import (
	"testing"
)

func TestParse(t *testing.T) {

	raw := `Signature: keyId="https://my-example.com/actor#main-key",headers="(request-target) host date",signature="Y2FiYW...IxNGRiZDk4ZA=="`

	sig, err := Parse(raw)

	if err != nil {
		t.Fatalf("Failed to parse raw signature, %v", err)
	}

	if sig.String() != raw {
		t.Fatalf("Stringified signature (%s) does not match original (%s)", sig.String(), raw)
	}
}
