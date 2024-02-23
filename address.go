package activitypub

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
)

const pat_addr string = `(?:acct\:)?@?([^@]+)(?:@(.*))?`

var re_addr = regexp.MustCompile(fmt.Sprintf(`^%s$`, pat_addr))

func ParseAddress(addr string) (string, string, error) {

	if !re_addr.MatchString(addr) {
		return "", "", fmt.Errorf("Failed to parse address")
	}

	m := re_addr.FindStringSubmatch(addr)
	return m[1], m[2], nil
}

func ParseAddressFromRequest(req *http.Request) (string, string, error) {

	resource := req.PathValue("resource")

	if resource == "" {
		return "", "", fmt.Errorf("request is missing {resource} path value")
	}

	slog.Debug("Parse address from request", "path", req.URL.Path, "resource", resource)
	return ParseAddress(resource)
}
