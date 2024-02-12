package sanitize

import (
	"net/http"
)

func RequestString(req *http.Request, param string) (string, error) {

	switch req.Method {

	case "POST":
		return PostString(req, param)
	default:
		return GetString(req, param)
	}

}

func RequestInt(req *http.Request, param string) (int, error) {

	switch req.Method {

	case "POST":
		return PostInt(req, param)
	default:
		return GetInt(req, param)
	}

}

func RequestInt64(req *http.Request, param string) (int64, error) {

	switch req.Method {

	case "POST":
		return PostInt64(req, param)
	default:
		return GetInt64(req, param)
	}

}

func RequestFloat64(req *http.Request, param string) (float64, error) {

	switch req.Method {

	case "POST":
		return PostFloat64(req, param)
	default:
		return GetFloat64(req, param)
	}

}

func RequestBool(req *http.Request, param string) (bool, error) {

	switch req.Method {

	case "POST":
		return PostBool(req, param)
	default:
		return GetBool(req, param)
	}

}
