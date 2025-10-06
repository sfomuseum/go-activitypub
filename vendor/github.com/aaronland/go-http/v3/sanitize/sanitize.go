package sanitize

import (
	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

var sn_opts *wof_sanitize.Options

func init() {
	sn_opts = wof_sanitize.DefaultOptions()
}
