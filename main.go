package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kagongjok/internal/healthcheck"
)

func main() {
	// Configure default structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	slog.Info("Started...")

	// Create context that cancels on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go healthcheck.Start(ctx)

	<-ctx.Done()

	slog.Info("Shutting down...")
}
