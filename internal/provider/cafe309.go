package provider

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"kagongjok/internal/util"

	"github.com/go-rod/rod"
)

const (
	// Cafe309LoginURLTemplate is the format string for the login redirect URL.
	Cafe309LoginURLTemplate = "https://service.wi2.ne.jp/wi2auth/redirect?cmd=login&mac=%s&apname=tunnel%%201&url=http%%3A%%2F%%2Fexample%%2Ecom%%2F%%3F%d"
)

type Cafe309Provider struct{}

func (c *Cafe309Provider) Name() string {
	return "309cafe"
}

func (c *Cafe309Provider) Login(ctx context.Context, page *rod.Page) error {
	slog.Info("Attempting 309 cafe automatic login...")

	mac := util.GetMacAddress()
	randNum := rand.Intn(999999999)
	targetURL := fmt.Sprintf(Cafe309LoginURLTemplate, mac, randNum)

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

	if info.Title == "logged in" || info.Title == "SMaRK" || info.Title == "success" {
		slog.Info("Already logged in.")
		return nil
	}

	slog.Info("Page loaded, looking for preliminary submit or agree button...")

	err = rod.Try(func() {
		if el, err := page.Element(`input[type="submit"]`); err == nil && el != nil {
			el.MustClick()
			page.MustWaitNavigation()
		}
	})

	slog.Info("Reading terms of use (TOP.html)...")
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	slog.Info("And clicking on accept...")

	_, err = page.Race().
		Element("#button_accept").MustHandle(func(e *rod.Element) {
		e.MustClick()
		slog.Info("Clicked accept button")
	}).
		Element("#alertArea a").MustHandle(func(e *rod.Element) { // just in case
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

	if finishTitle == "logged in" || finishTitle == "SMaRK" || finishTitle == "Example Domain" || strings.Contains(info.URL, "example.com") {
		slog.Info("Automatic login successful.")
		return nil
	}

	slog.Info("Automatic login failed.")
	return fmt.Errorf("autologin failed, title: %s", finishTitle)
}
