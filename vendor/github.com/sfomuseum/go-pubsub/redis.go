package pubsub

import (
	"fmt"
	"net/url"
	"strconv"
)

const REDIS_DEFAULT_HOST string = "localhost"
const REDIS_DEFAULT_PORT int = 6379

func RedisConfigFromURL(u *url.URL) (string, string, error) {

	channel := u.Host

	q := u.Query()

	host := REDIS_DEFAULT_HOST
	port := REDIS_DEFAULT_PORT

	if q.Has("host") {
		host = q.Get("host")
	}

	if q.Has("port") {
		str_port := q.Get("port")

		v, err := strconv.Atoi(str_port)

		if err != nil {
			return "", "", fmt.Errorf("Failed to parse ?port= parameter, %w", err)
		}

		port = v
	}

	if q.Has("channel") {
		channel = q.Get("channel")
	}

	if channel == "" {
		return "", "", fmt.Errorf("Empty or missing ?channel= parameter")
	}

	endpoint := fmt.Sprintf("%s:%d", host, port)
	return endpoint, channel, nil
}
