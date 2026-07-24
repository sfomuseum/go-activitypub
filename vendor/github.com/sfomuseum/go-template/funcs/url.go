package funcs

import (
	"net/url"
)

// URLQueryEscape escape's 'uri' using the net/url.QueryEscape function.
func URLQueryEscape(uri string) string {

	return url.QueryEscape(uri)
}
