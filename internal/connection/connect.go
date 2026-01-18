package connection

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"kagongjok/internal/config"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// AttemptConnection tries to log in to the Starbucks Wi2 WiFi.
func AttemptConnection(ctx context.Context) error {
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

	// Construct URL
	mac := GetMacAddress()
	randNum := rand.Intn(999999999)
	targetURL := fmt.Sprintf(config.StarbucksLoginURLTemplate, mac, randNum)

	// Navigate with timeout
	if err := page.Navigate(targetURL); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait load: %w", err)
	}

	info, err := page.Info()
	if err != nil {
		return fmt.Errorf("failed to get page info: %w", err)
	}

	if info.Title == "logged in" || info.Title == "at_STARBUCKS_Wi2" {
		slog.Info("Already logged in.")
		// Already logged in
		return nil
	}

	slog.Info("Page loaded, clicking on CONNECT...")

	// Click submit and wait for navigation
	// Check if element exists first
	_, err = page.Element(`input[type="submit"]`)
	if err != nil {
		return fmt.Errorf("failed to find submit button: %w", err)
	}

	// Attempt submit with navigation wait
	err = rod.Try(func() {
		page.MustElement(`input[type="submit"]`).MustClick()
		page.MustWaitNavigation()
	})
	if err != nil {
		return fmt.Errorf("failed to submit login form: %w", err)
	}

	slog.Info("Reading terms of use...")
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	slog.Info("And clicking on agree...")

	// Race condition handling using Rod's Race
	_, err = page.Race().
		Element("#button_accept").MustHandle(func(e *rod.Element) {
		e.MustClick()
		slog.Info("Clicked accept button")
	}).
		Element("#alertArea a").MustHandle(func(e *rod.Element) {
		e.MustClick()
		slog.Info("Clicked retry link")
	}).
		Do()

	if err != nil {
		return fmt.Errorf("failed to click accept or retry: %w", err)
	}

	// Wait for final navigation
	err = page.WaitLoad()
	if err != nil {
		return fmt.Errorf("failed to wait for final load: %w", err)
	}

	info, err = page.Info()
	if err != nil {
		return fmt.Errorf("failed to get info: %w", err)
	}
	finishTitle := info.Title

	slog.Info(fmt.Sprintf("Navigated to %s", finishTitle))

	if finishTitle == "logged in" || finishTitle == "at_STARBUCKS_Wi2" {
		slog.Info("Automatic login successful.")
	} else {
		slog.Info("Automatic login failed.")
		return fmt.Errorf("autologin failed, title: %s", finishTitle)
	}

	return nil
}
