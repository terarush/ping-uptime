#!/bin/bash

TARGET_DIR="/usr/local/bin"

echo " Installing ping-uptime (Auto-detecting version)"

echo "Detecting the latest release version from GitHub..."
VERSION=$(curl -s https://api.github.com/repos/terarush/ping-uptime/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)

if [ -z "$VERSION" ]; then
    echo "Failed to detect version automatically. Using fallback v0.0.3..."
    VERSION="v0.0.3"
else
    echo "Latest version detected: $VERSION"
fi

DOWNLOAD_URL="https://github.com/terarush/ping-uptime/releases/download/${VERSION}/ping-uptime"

echo "Downloading binary from GitHub..."
curl -L "$DOWNLOAD_URL" -o /tmp/ping-uptime

if [ $? -ne 0 ]; then
    echo "Failed to download file. Please check your internet connection."
    exit 1
fi

chmod +x /tmp/ping-uptime

echo "Moving binary to $TARGET_DIR..."
sudo mv /tmp/ping-uptime "$TARGET_DIR/ping-uptime"

if [ $? -ne 0 ]; then
    echo "Failed to move file. Make sure you have sudo privileges."
    exit 1
fi

DETECTED_SHELL=$(basename "$SHELL")
RC_FILE=""

if [ "$DETECTED_SHELL" = "zsh" ]; then
    RC_FILE="$HOME/.zshrc"
elif [ "$DETECTED_SHELL" = "bash" ]; then
    RC_FILE="$HOME/.bashrc"
fi

if [ -n "$RC_FILE" ] && [ -f "$RC_FILE" ]; then
    if ! grep -q "$TARGET_DIR" "$RC_FILE"; then
        echo "Adding $TARGET_DIR to PATH in $RC_FILE..."
        echo "export PATH=\"$TARGET_DIR:\$PATH\"" >> "$RC_FILE"
        echo "Shell configuration updated successfully."
        echo "Please run: source $RC_FILE to refresh your terminal."
    else
        echo "PATH $TARGET_DIR is already registered in $RC_FILE."
    fi
else
    echo "Unknown shell or configuration file not found. Please ensure $TARGET_DIR is in your PATH manually."
fi

echo "Installation Complete! Try running: ping-uptime -h"
