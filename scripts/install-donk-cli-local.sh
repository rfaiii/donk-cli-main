#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "${INSTALL_DIR}"

echo "Building donk-cli..."
cd "${ROOT}"
go build -o "${INSTALL_DIR}/donk-cli" "${ROOT}"

echo ""
echo "Installed donk-cli to:"
echo "  ${INSTALL_DIR}/donk-cli"
echo ""
echo "Make sure ${INSTALL_DIR} is in your PATH:"
echo "  export PATH=\"${INSTALL_DIR}:\${PATH}\""
echo ""
echo "Test with: donk-cli --version"
