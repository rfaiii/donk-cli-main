#!/usr/bin/env bash
# Launch DONK inside Alacritty using the bundled host config.
# Usage:
#   ./scripts/run-alacritty-host.sh
#   ./scripts/run-alacritty-host.sh --skip-splash
#   DONK_BIN=/path/to/donk ./scripts/run-alacritty-host.sh
#
# This script intentionally targets Alacritty only,
# because `/host write` emits an Alacritty-specific host config.
# For other terminals, see: /host launch

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_TEMPLATE="$ROOT/resources/host/alacritty/donk.toml"
RUNTIME_DIR="${TMPDIR:-/tmp}/donk-alacritty-host"
RUNTIME_CONFIG="$RUNTIME_DIR/donk.toml"

SKIP_SPLASH=0
THEME=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --skip-splash) SKIP_SPLASH=1; shift ;;
    --theme) THEME="${2:-}"; shift 2 ;;
    *) echo "Unknown arg: $1" >&2; exit 2 ;;
  esac
done

if [[ ! -f "$CONFIG_TEMPLATE" ]]; then
  echo "Missing host config: $CONFIG_TEMPLATE" >&2
  exit 1
fi

if [[ ! -d "$ROOT" ]]; then
  echo "Repo root missing: $ROOT" >&2
  exit 1
fi

# mkdir -p is safe even if the dir already exists.
mkdir -p "$RUNTIME_DIR"

DONK_BIN="${DONK_BIN:-}"
if [[ -z "$DONK_BIN" ]]; then
  if [[ -x "$ROOT/target/release/donk" ]]; then
    DONK_BIN="$ROOT/target/release/donk"
  elif [[ -x "$ROOT/target/debug/donk" ]]; then
    DONK_BIN="$ROOT/target/debug/donk"
  elif command -v donk >/dev/null 2>&1; then
    DONK_BIN="$(command -v donk)"
  else
    echo "Building donk (debug)..."
    (cd "$ROOT" && cargo build -p donk-terminal --bin donk)
    DONK_BIN="$ROOT/target/debug/donk"
  fi
fi

if [[ ! -x "$DONK_BIN" ]]; then
  echo "donk binary not executable: $DONK_BIN" >&2
  exit 1
fi

ALACRITTY_BIN="${ALACRITTY_BIN:-}"
if [[ -z "$ALACRITTY_BIN" ]]; then
  if command -v alacritty >/dev/null 2>&1; then
    ALACRITTY_BIN="$(command -v alacritty)"
  else
    echo "Alacritty not found on PATH. Install it or set ALACRITTY_BIN." >&2
    exit 1
  fi
fi

ARGS_TOML="[]"
if [[ "$SKIP_SPLASH" -eq 1 && -n "$THEME" ]]; then
  ARGS_TOML='["--skip-splash","'"$THEME"'"]'
elif [[ "$SKIP_SPLASH" -eq 1 ]]; then
  ARGS_TOML='["--skip-splash"]'
elif [[ -n "$THEME" ]]; then
  ARGS_TOML='["'"$THEME"'"]'
fi

# Deliberately quote only what needs quoting; avoid quoting the whole side.
sed -e 's|program = "donk"|program = "'"$DONK_BIN"'"|' \
    -e 's|args = \[\]|args = '"$ARGS_TOML"'|' \
    "$CONFIG_TEMPLATE" > "$RUNTIME_CONFIG"

echo "DONK Alacritty host"
echo "  alacritty: $ALACRITTY_BIN"
echo "  donk:      $DONK_BIN"
echo "  config:    $RUNTIME_CONFIG"

exec "$ALACRITTY_BIN" --config-file "$RUNTIME_CONFIG"
