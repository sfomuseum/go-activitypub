package www

import (
	"net/http"

	"github.com/sfomuseum/go-activitypub/ap"
)

func IsActivityStreamRequest(req *http.Request) bool {

	switch req.Header.Get("Accept") {

	case ap.ACTIVITYSTREAMS_ACCEPT_HEADER:
		return true
	default:
		return false
	}

}
