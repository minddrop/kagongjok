# Kagongjok (Cafe Study Helper) ☕️💻

**Stay Connected at Starbucks Japan.**

Kagongjok acts as your personal WiFi assistant. It runs in the background and automatically logs you back into the "at_STARBUCKS_Wi2" network whenever your session times out (usually every hour).

Perfect for long study or work sessions where you don't want to break your flow to click "Connect" again.

## Features

- 🔄 **Auto-Reconnect**: Detects when internet is lost and instantly logs you back in.
- 💨 **Fast**: Connects in seconds.
- 🤖 **Hands-free**: Just run it and forget it.

## How to Use

### Prerequisites

1. **Go**: [Install Go](https://go.dev/dl/) on your laptop.
2. **Chrome/Chromium**: The tool uses a browser in the background to handle the login page.

### Running it

1. Open your terminal.
2. Navigate to the project folder.
3. Run the tool:

    ```bash
    go run .
    ```

4. Keep the terminal open. You will see logs telling you when it's checking connectivity and when it reconnects you.

## How it Works

Every 10 seconds, it ping-checks Google. If the ping fails (because the WiFi timed out), it launches a hidden browser window, navigates to the Starbucks login portal, clicks "Connect" and "Accept", and then goes back to sleep.

## Important Note

This tool is designed for use at Starbucks locations in Japan that use the `at_STARBUCKS_Wi2` network.
