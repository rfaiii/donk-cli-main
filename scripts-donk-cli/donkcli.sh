#!/usr/bin/env bash
# Launch DONK as the `donkcli` command.
# Usage:
#   ./scripts/donkcli.sh
#   ./scripts/donkcli.sh --skip-splash
#   ./scripts/donkcli.sh doctor
#   ./scripts/donkcli.sh -- --help
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export CARGO_TARGET_DIR="$ROOT/target"

echo "==> Building donkcli…"
cargo build -p donk-cli --bin donkcli
DONKCLI="$ROOT/target/debug/donkcli"

if [[ $# -eq 0 ]]; then
  set -- --skip-splash
fi

echo "==> donkcli $*"
exec "$DONKCLI" "$@"
