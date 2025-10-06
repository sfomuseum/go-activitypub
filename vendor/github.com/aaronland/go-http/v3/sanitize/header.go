package sanitize

import (
	wof_sanitize "github.com/whosonfirst/go-sanitize"
	go_http "net/http"
	"strconv"
)

func HeaderString(req *go_http.Request, param string) (string, error) {

	raw_value := req.Header.Get(param)
	return wof_sanitize.SanitizeString(raw_value, sn_opts)
}

func HeaderInt64(req *go_http.Request, param string) (int64, error) {

	str_value, err := HeaderString(req, param)

	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(str_value, 10, 64)
}

func HeaderBool(req *go_http.Request, param string) (bool, error) {

	str_value, err := HeaderString(req, param)

	if err != nil {
		return false, err
	}

	if str_value == "" {
		return false, nil
	}

	return strconv.ParseBool(str_value)
}
