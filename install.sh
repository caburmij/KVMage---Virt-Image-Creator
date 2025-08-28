#!/bin/bash

set -e

BINARY_NAME="kvmage"
DEFAULT_INSTALL_DIR="/usr/local/bin"
VERSION_FILE="VERSION"

# Allow override of install directory via env var INSTALL_DIR
INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

# Detect platform
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

BINARY_PATH="dist/${BINARY_NAME}-${OS}-${ARCH}"

if [[ ! -f "$BINARY_PATH" ]]; then
  echo "Binary not found at $BINARY_PATH"
  exit 1
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."

# Check if INSTALL_DIR exists
if [[ ! -d "$INSTALL_DIR" ]]; then
  echo "Install directory $INSTALL_DIR does not exist, creating it with sudo..."
  sudo mkdir -p "$INSTALL_DIR"
fi

# Check if we have write permission to INSTALL_DIR
if [[ -w "$INSTALL_DIR" ]]; then
  cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
  chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
  echo "No write permission to $INSTALL_DIR, installing with sudo..."
  sudo cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
  sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Show version if available
if [[ -f "$VERSION_FILE" ]]; then
  VERSION=$(cat "$VERSION_FILE")
  echo "$BINARY_NAME v$VERSION installed successfully!"
else
  echo "$BINARY_NAME installed successfully!"
fi
