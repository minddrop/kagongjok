package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-rod/rod"
)

type Cafe309Provider struct{}

func (c *Cafe309Provider) Name() string {
	return "309cafe"
}

func (c *Cafe309Provider) Login(ctx context.Context, page *rod.Page) error {
	slog.Info("Attempting 309 cafe automatic login...")

	// @TODO: Implement navigation, wait, and clicking based on Feasibility Check findings.
	// Step 1: page.Navigate(targetURL)
	// Step 2: Click the agreement button or submit form
	// Step 3: Wait for navigation

	return fmt.Errorf("309 cafe provider is not yet fully implemented")
}
