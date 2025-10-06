package sanitize

import (
	"net/http"
	"strconv"

	wof_sanitize "github.com/whosonfirst/go-sanitize"
)

func GetString(req *http.Request, param string) (string, error) {

	q := req.URL.Query()
	raw_value := q.Get(param)
	return wof_sanitize.SanitizeString(raw_value, sn_opts)
}

func GetInt(req *http.Request, param string) (int, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.Atoi(str_value)
}

func GetInt64(req *http.Request, param string) (int64, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseInt(str_value, 10, 64)
}

func GetFloat64(req *http.Request, param string) (float64, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return 0, err
	}

	if str_value == "" {
		return 0, nil
	}

	return strconv.ParseFloat(str_value, 64)
}

func GetBool(req *http.Request, param string) (bool, error) {

	str_value, err := GetString(req, param)

	if err != nil {
		return false, err
	}

	if str_value == "" {
		return false, nil
	}

	return strconv.ParseBool(str_value)
}
