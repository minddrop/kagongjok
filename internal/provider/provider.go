package provider

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
)

// Provider abstracts the captive portal login logic for different Wi-Fi networks.
type Provider interface {
	// Name returns the identifier of the Wi-Fi provider.
	Name() string
	// Login attempts to authenticate with the captive portal.
	Login(ctx context.Context, page *rod.Page) error
}

// GetProvider returns a Provider implementation by its name.
func GetProvider(name string) (Provider, error) {
	switch name {
	case "starbucks":
		return &StarbucksProvider{}, nil
	case "309cafe":
		return &Cafe309Provider{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
