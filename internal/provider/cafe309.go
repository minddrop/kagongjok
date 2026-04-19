package provider

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"kagongjok/internal/config"
	"kagongjok/internal/util"

	"github.com/go-rod/rod"
)

const (
	// Cafe309LoginURLTemplate is the format string for the login redirect URL.
	Cafe309LoginURLTemplate = "https://service.wi2.ne.jp/wi2auth/redirect?cmd=login&mac=%s&apname=tunnel%%201&url=http%%3A%%2F%%2Fexample%%2Ecom%%2F%%3F%d"
)

type Cafe309Provider struct{}

func (c *Cafe309Provider) Name() string {
	return "309"
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
		if el, err := page.Timeout(3 * time.Second).Element(`input[type="submit"]`); err == nil && el != nil {
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
	if err == nil {
		slog.Info(fmt.Sprintf("Navigated to %s", info.Title))
	}

	slog.Info("Verifying internet connectivity...")
	client := &http.Client{Timeout: 3 * time.Second}

	for i := 0; i < 5; i++ {
		req, err := http.NewRequestWithContext(ctx, "GET", config.HealthCheckTarget, nil)
		if err == nil {
			resp, err := client.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == 200 || resp.StatusCode == 204 {
					slog.Info("Automatic login successful.")
					return nil
				}
			}
		}

		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	slog.Info("Automatic login failed (no internet connectivity).")
	title := "unknown"
	if info != nil {
		title = info.Title
	}
	return fmt.Errorf("autologin failed, no internet connectivity established; final title: %s", title)
}
