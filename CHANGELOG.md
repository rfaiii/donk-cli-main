# Changelog

## 1.1.4

### Mobile companion bridge

- Added mobile companion bridge planning and server design.
- Added `docs/MOBILE-CLI.md` for the iOS/Android companion roadmap.
- Added `docs/DONK-SERVER.md` for the host-side companion server.
- Added iOS companion app scaffold under `mobile-cli/ios/`.

### Validation

- Documented companion server architecture and protocol requirements.
- Added mobile companion app scaffolding with SwiftUI.

## 1.1.3

### Native Ollama coding

- Added `donk-cli code --native` for small local Ollama coding models.
- Added the `cmd/codetool` sidecar with a native Ollama `/api/chat` client and
  a minimal flat coding tool set.
- Added recovery for Qwen-style text tool calls, fenced JSON, truncated view
  calls, object-form arguments, and `file_path` aliases.
- Added bounded Ollama turns and explicit incomplete-stream errors so stalled
  local models do not leave the terminal loading indefinitely.
- Buffered likely tool-call JSON during streaming so protocol payloads do not
  clutter the user-facing response.
- Added installation, diagnostics, and local-model documentation for the native
  coder path.

### Bottom resource bar

- Added a live CPU and RAM resource monitor bar at the bottom of the TUI.
- The bar updates every second and shows CPU usage percentage and heap RAM
  usage percentage with token-driven theme colors.

### Validation

The native path was live-tested with `qwen2.5-coder:3b-instruct` against a local
Ollama daemon on macOS. Chat/file-read turns completed successfully in roughly
2–10 seconds during validation.
