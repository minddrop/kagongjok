package connection

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"kagongjok/internal/provider"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// AttemptConnection tries to log in to the Wi-Fi using the specified provider.
func AttemptConnection(ctx context.Context, p provider.Provider) error {
	slog.Info("Attempting connection...")

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Launch browser headless
	if ctx.Err() != nil {
		return ctx.Err()
	}

	l := launcher.New().Headless(true)
	if path, has := launcher.LookPath(); has {
		l.Bin(path)
	}

	u, err := l.Launch()
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

// EnsureBrowser checks if a browser binary is installed locally.
// If not, it downloads the required browser binary so that browser automation
// is ready even when internet connection is lost later.
func EnsureBrowser(ctx context.Context) error {
	l := launcher.NewBrowser()
	l.Context = ctx

	if l.Validate() == nil {
		slog.Debug("Browser binary requirement already installed", "path", l.BinPath())
		return nil
	}

	slog.Info("Downloading browser binary requirement (Chromium)...")

	path, err := l.Get()
	if err != nil {
		return fmt.Errorf("failed to download browser binary: %w", err)
	}

	slog.Info("Browser binary requirement installed successfully.", "path", path)
	return nil
}


