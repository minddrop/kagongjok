package connection

import (
	"context"
	"fmt"
	"log/slog"

	"kagongjok/internal/provider"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// AttemptConnection tries to log in to the Wi-Fi using the specified provider.
func AttemptConnection(ctx context.Context, p provider.Provider) error {
	slog.Info("Attempting connection...")

	// Launch browser headless
	if ctx.Err() != nil {
		return ctx.Err()
	}

	u, err := launcher.New().Headless(true).Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	// Bind context to page
	page = page.Context(ctx)

	if err := p.Login(ctx, page); err != nil {
		return err
	}

	return nil
}
