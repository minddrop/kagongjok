# Kagongjok

Kagongjok is a background script that automatically logs you back into captive portal Wi-Fi networks when your session times out. It pings an external server to check connectivity and uses a headless browser to complete the login portal flow when the internet goes down.

## Installation

### Option 1: Quick Install (macOS / Linux)
You can easily install the latest pre-compiled release by running:
```bash
curl -sSfL https://raw.githubusercontent.com/minddrop/kagongjok/main/install.sh | sh
```

### Option 2: Install with Go (All platforms)
If you already have Go installed, the easiest way is:
```bash
go install github.com/minddrop/kagongjok@latest
```

### Option 3: Build from Source (All platforms)
If you prefer to clone the repository and build it yourself:
```bash
# Clone the repository and build
git clone https://github.com/minddrop/kagongjok
cd kagongjok
go build -o kagongjok .
```

### Option 4: Download Pre-compiled Binary
You can download the pre-compiled binary (`.zip` for Windows, `.tar.gz` for macOS/Linux) directly from the [Releases page](https://github.com/minddrop/kagongjok/releases).

## Usage

If you installed via the quick install script or built the binary:
```bash
# Run with the default provider (Starbucks Japan)
kagongjok

# Run with an alternative provider
kagongjok 309
```

If you prefer to run it using Go directly without installing:
```bash
go run . 309
```

## Uninstallation

Depending on how you installed Kagongjok, use the matching uninstallation command:

If you installed via **Option 1** (Quick Install script):
```bash
sudo rm /usr/local/bin/kagongjok
```

If you installed via **Option 2** (`go install`):
```bash
rm $(go env GOPATH)/bin/kagongjok
```

If you downloaded the binary or built it from source manually, simply delete the `kagongjok` executable file.

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
