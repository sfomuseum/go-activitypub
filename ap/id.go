package ap

import (
	// "log/slog"
	"fmt"

	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
)

func NewId(uris_table *uris.URIs) string {
	uuid := id.NewUUID()
	u := uris.NewURL(uris_table, "/")
	u.Fragment = fmt.Sprintf("as-%s", uuid)

	// slog.Debug("New activitpub ID", "id", u.String())
	return u.String()
}
