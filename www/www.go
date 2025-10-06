package www

import (
	"net/http"
	"strings"

	"github.com/sfomuseum/go-activitypub/ap"
)

func IsActivityStreamRequest(req *http.Request, header string) bool {

	raw := req.Header.Get(header)
	accept := strings.Split(raw, ",")

	is_activitystream := false

	for _, h := range accept {

		h = strings.TrimSpace(h)

		switch h {

		case ap.ACTIVITYSTREAMS_ACCEPT_HEADER:
			is_activitystream = true
			break
		case ap.ACTIVITY_CONTENT_TYPE:
			is_activitystream = true
			break
		default:
			continue
		}
	}

	return is_activitystream
}
