#!/bin/bash

set -euo pipefail

REPO_URL="https://gitlab.ctos.io/code/kvmage-virt-image-creator.git"
REPO_DIR="kvmage-virt-image-creator"

echo "[*] Cloning repository..."
git clone "$REPO_URL"

echo "[*] Entering repo directory..."
cd "$REPO_DIR"

echo "[*] Creating dist directory..."
mkdir -p dist

echo "[*] Running build.sh..."
bash build.sh

echo "[*] Running install.sh..."
bash install.sh

echo "[*] Cleaning up..."
cd ..
rm -rf "$REPO_DIR"

echo "[*] Done."
