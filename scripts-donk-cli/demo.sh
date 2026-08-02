#!/usr/bin/env bash
# First DONK demo — build, doctor, then launch the TUI.
# Usage:
#   ./scripts/demo.sh
#   ./scripts/demo.sh --doctor-only
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

DOCTOR_ONLY=0
THEME="donk-dark"
for arg in "$@"; do
  case "$arg" in
    --doctor-only) DOCTOR_ONLY=1 ;;
    --theme=*) THEME="${arg#*=}" ;;
  esac
done

export CARGO_TARGET_DIR="$ROOT/target"
echo "==> Building donk…"
cargo build -p donk-terminal --bin donk
DONK="$ROOT/target/debug/donk"

echo ""
echo "==> doctor"
"$DONK" doctor || true

if [[ "$DOCTOR_ONLY" -eq 1 ]]; then
  exit 0
fi

echo ""
echo "==> demo guide"
"$DONK" demo
echo ""
echo "==> launching TUI (ctrl+c to quit)"
echo "    try: /help  /sys  /files  /read  /animations"
exec "$DONK" --skip-splash --theme "$THEME"
