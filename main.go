package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kagongjok/internal/healthcheck"
	"kagongjok/internal/provider"
)

func main() {
	providerName := flag.String("provider", "starbucks", "Wi-Fi provider to connect to (starbucks, 309)")
	flag.Parse()

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

	p, err := provider.GetProvider(*providerName)
	if err != nil {
		slog.Error("Failed to initialize provider", "error", err)
		os.Exit(1)
	}

	// Create context that cancels on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go healthcheck.Start(ctx, p)

	<-ctx.Done()

	slog.Info("Shutting down...")
}
