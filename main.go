package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kagongjok/internal/healthcheck"
	"kagongjok/internal/provider"
)

type PrettyHandler struct {
	provider string
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= slog.LevelInfo
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	attrStr := ""
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			attrStr += fmt.Sprintf(" %s=%v", a.Key, a.Value)
			return true
		})
	}

	fmt.Printf("%s [%s] [%s] %s%s\n", timeStr, r.Level.String(), h.provider, r.Message, attrStr)
	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *PrettyHandler) WithGroup(name string) slog.Handler       { return h }

func main() {
	providerName := flag.String("provider", "starbucks", "Wi-Fi provider to connect to (starbucks, 309)")
	flag.Parse()

	// Configure default structured logger
	logger := slog.New(&PrettyHandler{provider: *providerName})
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
