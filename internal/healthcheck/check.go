package healthcheck

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"kagongjok/internal/config"
	"kagongjok/internal/connection"
)

func Start(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		slog.Info("Checking connectivity...")

		var nextWait time.Duration

		if err := check(ctx); err == nil {
			slog.Info("No connectivity problems detected.")
			nextWait = config.CheckInterval
		} else {
			slog.Error("Health check failed", "error", err)

			if err := connection.AttemptConnection(ctx); err != nil {
				if ctx.Err() != nil {
					return
				}
				slog.Error("Connection attempt failed", "error", err)
			}
			// Wait before retrying or checking again after a login attempt
			nextWait = config.LoginRetryDelay
		}

		select {
		case <-time.After(nextWait):
		case <-ctx.Done():
			return
		}
	}
}

func check(ctx context.Context) error {
	// Check internet connectivity by pinging a reliable host
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", config.HealthCheckTarget, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection refused: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
}
