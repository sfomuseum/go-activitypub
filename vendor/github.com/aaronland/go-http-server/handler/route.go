package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

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

		handler, err := deriveHandler(req, opts.Handlers, matches, patterns)

		if err != nil {
			opts.Logger.Printf("%v", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if handler == nil {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		handler.ServeHTTP(rsp, req)
		return
	}

	return http.HandlerFunc(fn), nil
}

func deriveHandler(req *http.Request, handlers map[string]RouteHandlerFunc, matches *sync.Map, patterns []string) (http.Handler, error) {

	ctx := req.Context()
	path := req.URL.Path

	// Basically do what the default http.ServeMux does but inflate the
	// handler (func) on demand at run-time. Handler is cached above.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.4:src/net/http/server.go;l=2363

	var matching_pattern string

	for _, p := range patterns {
		if strings.HasPrefix(path, p) {
			matching_pattern = p
			break
		}
	}

	if matching_pattern == "" {
		return nil, nil
	}

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

	return handler, nil
}
