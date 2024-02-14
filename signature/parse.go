package signature

// https://blog.joinmastodon.org/2018/07/how-to-make-friends-and-verify-requests/

import (
	"fmt"
	_ "log"
	"net/http"
	"regexp"
	"strings"
)

var re_pair = regexp.MustCompile(`^([^=]+)="([^\"]+)"$`)

type Signature struct {
	KeyId     string `json:"keyId"`
	Headers   string `json:"headers"`
	Signature string `json:"signature"`
}

func (s *Signature) String() string {
	return fmt.Sprintf(`Signature: keyId="%s",headers="%s",signature="%s"`, s.KeyId, s.Headers, s.Signature)
}

func ParseFromRequest(req *http.Request) (*Signature, error) {

	raw := req.Header.Get("Signature")

	if raw == "" {
		return nil, fmt.Errorf("Missing Signature header")
	}

	return Parse(raw)
}

func Parse(raw string) (*Signature, error) {

	raw = strings.Replace(raw, "Signature:", "", 1)
	raw = strings.TrimSpace(raw)

	parts := strings.Split(raw, ",")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid signature string")
	}

	sig := new(Signature)

	for idx, str_pair := range parts {

		str_pair = strings.TrimSpace(str_pair)

		if !re_pair.MatchString(str_pair) {
			return nil, fmt.Errorf("Invalid signature pair (re) at offset %d", idx)
		}

		m := re_pair.FindStringSubmatch(str_pair)

		k := m[1]
		v := m[2]

		switch k {
		case "keyId":
			sig.KeyId = v
		case "headers":
			sig.Headers = v
		case "signature":
			sig.Signature = v
		default:
			return nil, fmt.Errorf("Invalid key in pair at offset %d", idx)
		}
	}

	return sig, nil
}
