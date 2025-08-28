#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")" && pwd)"
DIST_DIR="${SCRIPT_DIR}/dist"
VERSION="$(cat "${SCRIPT_DIR}/VERSION")"

GOOS=linux  GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-linux-amd64" "$SCRIPT_DIR"
GOOS=linux  GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-linux-arm64" "$SCRIPT_DIR"
GOOS=darwin GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-darwin-amd64" "$SCRIPT_DIR"
GOOS=darwin GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-darwin-arm64" "$SCRIPT_DIR"
