package provider

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"kagongjok/internal/util"

	"github.com/go-rod/rod"
)

const (
	// StarbucksLoginURLTemplate is the format string for the login redirect URL.
	StarbucksLoginURLTemplate = "https://service.wi2.ne.jp/wi2auth/redirect?cmd=login&mac=%s&essid=%%20&apname=tunnel%%201&apgroup=&url=http%%3A%%2F%%2Fexample%%2Ecom%%2F%%3F%d"
)

type StarbucksProvider struct{}

func (s *StarbucksProvider) Name() string {
	return "starbucks"
}

func (s *StarbucksProvider) Login(ctx context.Context, page *rod.Page) error {
	slog.Info("Attempting Starbucks automatic login...")

	mac := util.GetMacAddress()
	randNum := rand.Intn(999999999)
	targetURL := fmt.Sprintf(StarbucksLoginURLTemplate, mac, randNum)

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
		return nil
	}

	slog.Info("Page loaded, clicking on CONNECT...")

	_, err = page.Element(`input[type="submit"]`)
	if err != nil {
		return fmt.Errorf("failed to find submit button: %w", err)
	}

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
		return nil
	}

	slog.Info("Automatic login failed.")
	return fmt.Errorf("autologin failed, title: %s", finishTitle)
}
