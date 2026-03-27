#!/bin/bash

set -e

REPO="iyushkarki/csys"
OWNER="iyushkarki"
REPO_NAME="csys"
BIN_NAME="csys"

echo "Installing csys..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  arm64) ARCH="arm64" ;;
  amd64|x86_64) ARCH="amd64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    echo "Fallback: Install with: go install github.com/$REPO@latest"
    exit 1
    ;;
esac

case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    echo "Fallback: Install with: go install github.com/$REPO@latest"
    exit 1
    ;;
esac

echo "Detecting system: $OS-$ARCH"

LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p')

if [ -z "$LATEST" ]; then
  echo "No releases found on GitHub"
  echo "Fallback: Install with Go:"
  echo "   go install github.com/$REPO@latest"
  echo ""
  echo "Ensure Go 1.19+ is installed: go version"
  exit 1
fi

echo "Downloading $LATEST for $OS-$ARCH..."

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST/${BIN_NAME}-${OS}-${ARCH}"

if ! curl -fsSL --head "$DOWNLOAD_URL" > /dev/null 2>&1; then
  echo "Binary not found at: $DOWNLOAD_URL"
  echo "Fallback: Install with Go:"
  echo "   go install github.com/$REPO@latest"
  exit 1
fi

TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$BIN_NAME"
chmod +x "$TEMP_DIR/$BIN_NAME"

if command -v sudo &> /dev/null; then
  BIN_PATH="/usr/local/bin/$BIN_NAME"
  ALIAS_PATH="/usr/local/bin/cs"
  echo "Requires sudo to install to $BIN_PATH"
  sudo mv "$TEMP_DIR/$BIN_NAME" "$BIN_PATH"
  sudo ln -sf "$BIN_PATH" "$ALIAS_PATH"
else
  BIN_PATH="$HOME/.local/bin/$BIN_NAME"
  ALIAS_PATH="$HOME/.local/bin/cs"
  mkdir -p "$HOME/.local/bin"
  mv "$TEMP_DIR/$BIN_NAME" "$BIN_PATH"
  ln -sf "$BIN_PATH" "$ALIAS_PATH"

  if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
    echo ""
    echo "Warning: $HOME/.local/bin is not in your PATH"
    echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc):"
    echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
  fi
fi

install_completions() {
  SHELL_NAME=$(basename "$SHELL")

  case "$SHELL_NAME" in
    zsh)
      COMP_DIR="${HOME}/.zsh/completions"
      mkdir -p "$COMP_DIR"
      "$BIN_PATH" completion zsh > "$COMP_DIR/_csys" 2>/dev/null
      sed 's/csys/cs/g; s/#compdef csys/#compdef cs/' "$COMP_DIR/_csys" > "$COMP_DIR/_cs" 2>/dev/null

      ZSHRC="${HOME}/.zshrc"
      if [ -f "$ZSHRC" ]; then
        if ! grep -q 'fpath.*\.zsh/completions' "$ZSHRC" 2>/dev/null; then
          echo "" >> "$ZSHRC"
          echo 'fpath=(~/.zsh/completions $fpath)' >> "$ZSHRC"
          echo 'autoload -Uz compinit && compinit' >> "$ZSHRC"
        fi
      fi
      echo "Zsh completions installed to $COMP_DIR"
      echo "Run: source ~/.zshrc (or open a new terminal)"
      ;;
    bash)
      if [ -d "/etc/bash_completion.d" ] && command -v sudo &> /dev/null; then
        "$BIN_PATH" completion bash | sudo tee /etc/bash_completion.d/csys > /dev/null 2>&1
        "$BIN_PATH" completion bash | sed 's/csys/cs/g' | sudo tee /etc/bash_completion.d/cs > /dev/null 2>&1
        echo "Bash completions installed to /etc/bash_completion.d/"
      else
        COMP_DIR="${HOME}/.local/share/bash-completion/completions"
        mkdir -p "$COMP_DIR"
        "$BIN_PATH" completion bash > "$COMP_DIR/csys" 2>/dev/null
        "$BIN_PATH" completion bash | sed 's/csys/cs/g' > "$COMP_DIR/cs" 2>/dev/null
        echo "Bash completions installed to $COMP_DIR"
      fi
      echo "Open a new terminal to activate completions"
      ;;
    fish)
      COMP_DIR="${HOME}/.config/fish/completions"
      mkdir -p "$COMP_DIR"
      "$BIN_PATH" completion fish > "$COMP_DIR/csys.fish" 2>/dev/null
      "$BIN_PATH" completion fish | sed 's/csys/cs/g' > "$COMP_DIR/cs.fish" 2>/dev/null
      echo "Fish completions installed to $COMP_DIR"
      ;;
    *)
      echo "Shell completions: run 'csys completion --help' to set up manually"
      ;;
  esac
}

install_completions

echo ""
echo "csys installed successfully!"
echo "Location: $BIN_PATH"
echo "Alias: cs -> $BIN_PATH"
echo ""
echo "Quick start:"
echo "   csys              Show system overview"
echo "   cs                Same as csys (shorter!)"
echo "   cs --live         Live monitoring"
echo "   cs gac \"msg\"      Git: add + commit"
echo "   cs gsync          Git: reset to origin/main"
echo "   cs --help         Show all options"
echo ""

$BIN_PATH --version 2>/dev/null || $BIN_PATH --help | head -2

echo ""
echo "Happy monitoring!"
