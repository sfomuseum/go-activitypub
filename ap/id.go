package ap

import (
	"fmt"
	// "log/slog"

	"github.com/sfomuseum/go-activitypub/id"
	"github.com/sfomuseum/go-activitypub/uris"
)

// NewId return a new identifier in the form of a unique URI.
func NewId(uris_table *uris.URIs) string {

	uuid := id.NewUUID()

	u := uris.NewURL(uris_table, uris_table.Root)
	u.Fragment = fmt.Sprintf("as-%s", uuid)

	// slog.Debug("New ap ID", "id", u.String())
	return u.String()
}
