#!/bin/bash

set -e

REPO="iyushkarki/csys"
OWNER="iyushkarki"
REPO_NAME="csys"
BIN_NAME="csys"

echo "🚀 Installing csys..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  arm64) ARCH="arm64" ;;
  amd64|x86_64) ARCH="amd64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    echo "💡 Fallback: Install with: go install github.com/$REPO@latest"
    exit 1
    ;;
esac

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  *)
    echo "❌ Unsupported OS: $OS"
    echo "💡 Fallback: Install with: go install github.com/$REPO@latest"
    exit 1
    ;;
esac

echo "📦 Detecting system: $OS-$ARCH"

LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p')

if [ -z "$LATEST" ]; then
  echo "⚠️  No releases found on GitHub"
  echo "💡 Fallback: Install with Go:"
  echo "   go install github.com/$REPO@latest"
  echo ""
  echo "📝 Ensure Go 1.19+ is installed: go version"
  exit 1
fi

echo "📥 Downloading $LATEST for $OS-$ARCH..."

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST/${BIN_NAME}-${OS}-${ARCH}"

if ! curl -fsSL --head "$DOWNLOAD_URL" > /dev/null 2>&1; then
  echo "❌ Binary not found at: $DOWNLOAD_URL"
  echo "💡 Fallback: Install with Go:"
  echo "   go install github.com/$REPO@latest"
  exit 1
fi

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$BIN_NAME"
chmod +x "$TEMP_DIR/$BIN_NAME"

if command -v sudo &> /dev/null; then
  BIN_PATH="/usr/local/bin/$BIN_NAME"
  echo "📝 Requires sudo to install to $BIN_PATH"
  sudo mv "$TEMP_DIR/$BIN_NAME" "$BIN_PATH"
else
  BIN_PATH="$HOME/.local/bin/$BIN_NAME"
  mkdir -p "$HOME/.local/bin"
  mv "$TEMP_DIR/$BIN_NAME" "$BIN_PATH"

  if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
    echo ""
    echo "⚠️  Warning: $HOME/.local/bin is not in your PATH"
    echo "📝 Add this to your shell profile (~/.bashrc, ~/.zshrc, etc):"
    echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
  fi
fi

echo ""
echo "✅ csys installed successfully!"
echo "📍 Location: $BIN_PATH"
echo ""
echo "🎯 Quick start:"
echo "   csys              # Show system overview"
echo "   csys --live       # Live monitoring (updates every 2s)"
echo "   csys --help       # Show all options"
echo ""

$BIN_PATH --version 2>/dev/null || $BIN_PATH --help | head -2

echo ""
echo "Happy monitoring! 🚀"
