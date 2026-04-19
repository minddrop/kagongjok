# Kagongjok

Kagongjok is a background script that automatically logs you back into captive portal Wi-Fi networks when your session times out. It pings an external server to check connectivity and uses a headless browser to complete the login portal flow when the internet goes down.

## Usage

You need Go installed on your system. 

```bash
# Clone the repository and install dependencies
go mod download

# Run with the default provider (Starbucks Japan)
go run .

# Run with an alternative provider
go run . -provider=309
```
*Note: You don't need to manually install Chrome or Chromium. The `go-rod` library automatically downloads a headless Chromium binary at runtime if one isn't found.*

## Supported Networks

- `starbucks` (at_STARBUCKS_Wi2, used at Starbucks in Japan)
- `309` (Cafe 309)

## Adding a new provider

The core logic is split between `healthcheck` and `provider`. To add a new Wi-Fi network:

1. Create a new provider file in `internal/provider/` (e.g., `mynetwork.go`).
2. Implement the `Provider` interface:
   ```go
   type Provider interface {
       Name() string
       Login(ctx context.Context, page *rod.Page) error
   }
   ```
3. Use the `page` object to automate the captive portal's login flow using `go-rod`.
4. Register your provider in `internal/provider/provider.go` inside `GetProvider()`.
