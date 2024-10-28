package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	_ "log/slog"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Regular expression to match "{label}" style substitutions in URL patterns
var re_braces = regexp.MustCompile(`\{([^\}]+)\}`)

// Regular expression to match Go 1.22 "[METHOD] [HOST]/[PATH]" URL definitions. This is probably
// not as robust as it could be.
var re_route = regexp.MustCompile(`^(?:(?:(GET|POST|PUT|HEAD|OPTION|DELETE)\s)?([^\/]+)?)?(.*)$`)

var re_mutex = new(sync.RWMutex)

// Regular expression to replace "{label}" style strings with "([^\/]+)".
func not_a_slash(s string) string {
	return fmt.Sprintf(`([^\/]+)`)
}

// Local cache of wildcard regular expressions that have been derived from patterns with "{label}" style substitutions
var wildcard_matches = make(map[string]*regexp.Regexp)

// Local struct for return key value pairs that have been derived from URL with "{label}" style substitutions
// and their subsequent wildcard matches
type pathValue struct {
	Key   string
	Value string
}

// String returns the key and value pair in a pathValue as a "key=value" formatted string.
func (pv *pathValue) String() string {
	return fmt.Sprintf("%s=%s", pv.Key, pv.Value)
}

// RouteHandlerFunc returns an `http.Handler` instance.
type RouteHandlerFunc func(context.Context) (http.Handler, error)

// RouteHandlerOptions is a struct that contains configuration settings
// for use the RouteHandlerWithOptions method.
type RouteHandlerOptions struct {
	// Handlers is a map whose keys are `http.ServeMux` style routing patterns and whose keys
	// are functions that when invoked return `http.Handler` instances.
	Handlers map[string]RouteHandlerFunc
	// Logger is a `log.Logger` instance used for feedback and error-reporting.
	Logger *log.Logger
}

// RouteHandler create a new `http.Handler` instance that will serve requests using handlers defined in 'handlers'
// with a `log.Logger` instance that discards all writes. Under the hood this is invoking the `RouteHandlerWithOptions`
// method.
func RouteHandler(handlers map[string]RouteHandlerFunc) (http.Handler, error) {

	logger := log.New(io.Discard, "", 0)

	opts := &RouteHandlerOptions{
		Handlers: handlers,
		Logger:   logger,
	}

	return RouteHandlerWithOptions(opts)
}

// RouteHandlerWithOptions create a new `http.Handler` instance that will serve requests using handlers defined
// in 'opts.Handlers'. This is essentially a "middleware" handler than does all the same routing that the default
// `http.ServeMux` handler does but defers initiating the handlers being routed to until they invoked at runtime.
// Only one handler is initialized (or retrieved from an in-memory cache) and served for any given path being by
// a `RouteHandler` request.
//
// The reason this handler exists is for web applications that:
//
//  1. Are deployed as AWS Lambda functions (with an API Gateway integration) using the "lambda://" `server.Server`
//     implementation that have more handlers than you need or want to initiate, but never use, for every request.
//  2. You don't want to refactor in to (n) atomic Lambda functions. That is you want to be able to re-use the same
//     code in both a plain-vanilla HTTP server configuration as well as Lambda + API Gateway configuration.
//
// URL patterns for handlers (passed in the `RouteHandlerOptions` struct) are parsed using a custom implementation
// of Go 1.22's "URL pattern matching". While this implementation has been tested on common patterns it is likely that
// there are still edge-cases, or simply sufficiently sophisticated cases, that will not match. This is considered a
// "known known" and those cases will need to be addressed as they arise. This custom implementation uses regular
// expressions and has not been benchmarked against the Go 1.22 implementation. I would prefer to use the native
// implementation but it is a) code that is private to [net/http] and sufficiently involved that cloning it in to this
// package, and then tracking the changes, seems like a bad idea.
func RouteHandlerWithOptions(opts *RouteHandlerOptions) (http.Handler, error) {

	matches := new(sync.Map)
	patterns := make([]string, 0)

	for p, _ := range opts.Handlers {
		patterns = append(patterns, p)
	}

	// Sort longest to shortest
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.4:src/net/http/server.go;l=2533

	sort.Slice(patterns, func(i, j int) bool {
		return len(patterns[i]) > len(patterns[j])
	})

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		derive_rsp, err := deriveHandler(req, opts.Handlers, matches, patterns)

		if err != nil {
			opts.Logger.Printf("%v", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if derive_rsp == nil {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		if derive_rsp.Method != "" && derive_rsp.Method != req.Method {
			http.Error(rsp, "Method not allow", http.StatusMethodNotAllowed)
			return
		}

		if derive_rsp.Host != "" && derive_rsp.Host != req.Host {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		if derive_rsp.PathValues != nil {

			for _, pv := range derive_rsp.PathValues {
				req.SetPathValue(pv.Key, pv.Value)
			}
		}

		derive_rsp.Handler.ServeHTTP(rsp, req)
		return
	}

	return http.HandlerFunc(fn), nil
}

type deriveHandlerResults struct {
	MatchingPattern string
	Handler         http.Handler
	Method          string
	Host            string
	PathValues      []*pathValue
}

func deriveHandler(req *http.Request, handlers map[string]RouteHandlerFunc, matches *sync.Map, patterns []string) (*deriveHandlerResults, error) {

	ctx := req.Context()
	path := req.URL.Path

	// Basically do what the default http.ServeMux does but inflate the
	// handler (func) on demand at run-time. Handler is cached above.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.4:src/net/http/server.go;l=2363

	// That was before Go 1.22 's pattern routing which makes everything more complicated. What
	// follows is less complicated (or at least less twisty) than the net/http code which is all
	// private and internal and too much to clone in to this package. What follows may still contain
	// gotchas.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/net/http/server.go;l=2320

	var matching_pattern string
	var host string
	var method string

	var path_values []*pathValue

	for _, p := range patterns {

		// slog.Info("ROUTE", "pattern", p)

		// First just try the simple prefix-based approach

		if strings.HasPrefix(path, p) {
			matching_pattern = p
			break
		}

		// Next try to parse out [METHOD] [HOST]/[PATH]

		route_m := re_route.FindStringSubmatch(p)

		if len(route_m) == 0 {
			continue
		}

		method = route_m[1]
		host = route_m[2]
		route_path := route_m[3]

		// slog.Info("ROUTE", "path", route_path)

		if !re_braces.MatchString(route_path) {

			if strings.HasPrefix(path, route_path) {
				matching_pattern = p
				break
			} else {
				continue
			}
		}

		// If there are replace them with a match-up-to-next-forward-slash capture
		// and then use the result to build a new regular expression

		re_mutex.Lock()

		re_wildcard, exists := wildcard_matches[route_path]

		if !exists {

			str_wildcard := re_braces.ReplaceAllStringFunc(route_path, not_a_slash)
			re, err := regexp.Compile(str_wildcard)

			if err != nil {
				re_mutex.Unlock()
				return nil, fmt.Errorf("Failed to compile wildcard regexp, %w", err)
			}

			re_wildcard = re
			wildcard_matches[route_path] = re_wildcard
		}

		re_mutex.Unlock()

		// Does the current path (like the actual request being processed) match the wildcard?

		path_m := re_wildcard.FindStringSubmatch(path)

		if len(path_m) == 0 {
			continue
		}

		// If it does extract all the curly-substitution braces and then use the two matches
		// to build a set of key,value pairs to populate the request's PathValue lookup table

		key_m := re_braces.FindAllStringSubmatch(route_path, -1)

		count_k := len(key_m)
		path_values = make([]*pathValue, count_k)

		for i := 0; i < count_k; i++ {

			key := key_m[i][1]
			value := path_m[i+1]

			pv := &pathValue{
				Key:   key,
				Value: value,
			}

			path_values[i] = pv
		}

		matching_pattern = p
		break
	}

	if matching_pattern == "" {
		return nil, nil
	}

	// slog.Info("ROUTE", "path", path, "matching_pattern", matching_pattern, "pv", path_values)

	var handler http.Handler

	v, exists := matches.Load(matching_pattern)

	if exists {
		handler = v.(http.Handler)
	} else {

		handler_func, ok := handlers[matching_pattern]

		// Don't fill up the matches cache with 404 handlers

		if !ok {
			return nil, nil
		}

		h, err := handler_func(ctx)

		if err != nil {
			return nil, fmt.Errorf("Failed to instantiate handler func for '%s' matching '%s', %v", path, matching_pattern, err)
		}

		handler = h
		matches.Store(matching_pattern, handler)
	}

	rsp := &deriveHandlerResults{
		MatchingPattern: matching_pattern,
		Handler:         handler,
		PathValues:      path_values,
		Method:          method,
		Host:            host,
	}

	return rsp, nil
}
