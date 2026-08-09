package connection

import (
	"context"
	"testing"
	"time"

	"github.com/go-rod/rod/lib/launcher"
)

func TestEnsureBrowser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := EnsureBrowser(ctx)
	if err != nil {
		t.Fatalf("EnsureBrowser failed: %v", err)
	}

	found, has := launcher.LookPath()
	if !has || found == "" {
		t.Fatalf("Expected browser binary to be found, got found=%q, has=%v", found, has)
	}
}

func TestEnsureBrowserIdempotent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	err := EnsureBrowser(ctx)
	if err != nil {
		t.Fatalf("EnsureBrowser failed: %v", err)
	}
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("Cached EnsureBrowser took too long: %v", elapsed)
	}
}
