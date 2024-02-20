package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
	"github.com/sfomuseum/go-activitypub"
	ap_slog "github.com/sfomuseum/go-activitypub/slog"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := OptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	ap_slog.ConfigureLogger(logger, opts.Verbose)

	v, err := opts.clone()

	if err != nil {
		return fmt.Errorf("Failed to create local run options, %w", err)
	}

	run_opts = v

	webfinger_get := fmt.Sprintf("GET %s", activitypub.WEBFINGER_URI)
	profile_get := fmt.Sprintf("GET %s", run_opts.URIs.Profile)
	inbox_post := fmt.Sprintf("POST %s", run_opts.URIs.Inbox)
	// outbox_post := fmt.Sprintf("POST %s", run_opts.URIs.Outbox)

	handlers := map[string]handler.RouteHandlerFunc{
		webfinger_get: webfingerHandlerFunc,
		profile_get:   profileHandlerFunc,
		inbox_post:    inboxPostHandlerFunc,
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

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	slog.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
