#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
GOBIN_DIR=$(go env GOBIN)
if [ -z "$GOBIN_DIR" ]; then
	GOBIN_DIR=$(go env GOPATH)/bin
fi

go build -o "$ROOT/donk-cli" "$ROOT"
mkdir -p "$GOBIN_DIR"
cp "$ROOT/donk-cli" "$GOBIN_DIR/donk-cli"

# Keep the common user-local installation synchronized too. This avoids an
# older copy winning when both directories are present on PATH.
LOCAL_BIN="$HOME/.local/bin"
mkdir -p "$LOCAL_BIN"
cp "$ROOT/donk-cli" "$LOCAL_BIN/donk-cli"

printf 'Installed latest donk-cli to:\n  %s\n  %s\n' "$GOBIN_DIR/donk-cli" "$LOCAL_BIN/donk-cli"