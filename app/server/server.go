package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
	"github.com/sfomuseum/go-activitypub/webfinger"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	v, err := opts.clone()

	if err != nil {
		return fmt.Errorf("Failed to create local run options, %w", err)
	}

	run_opts = v

	// Use a "route handler" to defer creating any given route until it is
	// invoked. This is useful in "serverless" environments like AWS Lambda.

	webfinger_get := fmt.Sprintf("GET %s", webfinger.Endpoint)
	account_get := fmt.Sprintf("GET %s", run_opts.URIs.Account)
	inbox_post := fmt.Sprintf("POST %s", run_opts.URIs.Inbox)
	outbox_get := fmt.Sprintf("GET %s", run_opts.URIs.Outbox)
	post_get := fmt.Sprintf("GET %s", run_opts.URIs.Post)

	handlers := map[string]handler.RouteHandlerFunc{

		// HTML (human-facing) pages that may need or want custom chrome
		account_get: accountHandlerFunc,
		post_get:    postHandlerFunc,

		// This needs to be fixed in aaronland/go-http-server/handler
		// outbox_post:              outboxPostHandlerFunc,
		run_opts.URIs.Icon:      iconHandlerFunc,
		run_opts.URIs.Following: followingHandlerFunc,
		run_opts.URIs.Followers: followersHandlerFunc,
		webfinger_get:           webfingerHandlerFunc,
		inbox_post:              inboxPostHandlerFunc,
		outbox_get:              outboxGetHandlerFunc,
	}

	log_logger := slog.NewLogLogger(logger.Handler(), slog.LevelInfo)

	route_handler_opts := &handler.RouteHandlerOptions{
		Handlers: handlers,
		Logger:   log_logger,
	}

	route_handler, err := handler.RouteHandlerWithOptions(route_handler_opts)

	if err != nil {
		return fmt.Errorf("Failed to configure route handler, %w", err)
	}

	route_handler = handler.DisabledHandler(run_opts.Disabled, route_handler)

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	if run_opts.CustomHandlers != nil {

		err = run_opts.CustomHandlers(mux)

		if err != nil {
			return fmt.Errorf("Custom handlers function failed, %w", err)
		}
	}

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	logger.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
