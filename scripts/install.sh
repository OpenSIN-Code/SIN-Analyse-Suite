#!/usr/bin/env bash
set -euo pipefail
REPO="OpenSIN-Code/SIN-Analyse-Suite"
BINARY="sin-analyse"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac
LATEST=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
[ -z "$LATEST" ] && { echo "Failed to find latest release"; exit 1; }
URL="https://github.com/$REPO/releases/download/$LATEST/${BINARY}_${OS}_${ARCH}.tar.gz"
TMP=$(mktemp -d)
curl -sL "$URL" | tar xz -C "$TMP"
mv "$TMP/$BINARY" /usr/local/bin/$BINARY
chmod +x /usr/local/bin/$BINARY
rm -rf "$TMP"
echo "Installed $BINARY $LATEST"
