#!/bin/sh
set -e
REPO="minddrop/kagongjok"

OS=$(uname -s)
ARCH=$(uname -m)

if [ "$OS" = "Darwin" ]; then
    OS="Darwin"
elif [ "$OS" = "Linux" ]; then
    OS="Linux"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

if [ "$ARCH" = "x86_64" ]; then
    ARCH="x86_64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

echo "Detecting latest release for $OS $ARCH from $REPO..."
# Fetch the latest release URL for the matching OS and ARCH
LATEST_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep "browser_download_url" | grep -i "$OS" | grep -i "$ARCH" | grep "\.tar\.gz" | cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "Could not find a release for $OS $ARCH."
    exit 1
fi

echo "Downloading $LATEST_URL..."
TMP_DIR=$(mktemp -d)
curl -sL "$LATEST_URL" -o "$TMP_DIR/release.tar.gz"

echo "Extracting..."
tar -xzf "$TMP_DIR/release.tar.gz" -C "$TMP_DIR"

echo "Installing to /usr/local/bin/kagongjok (may require sudo)..."
sudo mv "$TMP_DIR/kagongjok" /usr/local/bin/kagongjok
sudo chmod +x /usr/local/bin/kagongjok

rm -rf "$TMP_DIR"
echo "Installation complete! Run 'kagongjok' to start."
