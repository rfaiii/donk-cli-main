# Changelog

## 1.1.5

### Mobile companion bridge

- Added mobile companion bridge planning and server design.
- Added `docs/MOBILE-CLI.md` for the iOS/Android companion roadmap.
- Added `docs/DONK-SERVER.md` for the host-side companion server.
- Added iOS companion app scaffold under `mobile-cli/ios/`.
- **Note:** the mobile companion is present for testing, but the mobile version
  is currently in **test mode** and not yet ready for general use.

### Validation

- Documented companion server architecture and protocol requirements.
- Added mobile companion app scaffolding with SwiftUI.

### Packaging & onboarding

- Added `docs/ONBOARDING.md` for first-run setup and dependency guidance.
- Added `docs/DEPENDENCIES.md` for concrete macOS/Windows install requirements.
- Added `docs/IMG-RESOURCES.md` for screenshot, video, and icon asset guidance.
- Added `docs/ABOUT.md` for company and maintainer contact info.
- Added packaging wrappers under `packaging/` for Homebrew, NPM, macOS, and Windows.
- Added `.goreleaser.dist.windows.yaml` for Windows archive packaging.
- Updated `.goreleaser.yml` tap/scoop ownership and metadata for distribution.
- Added macOS DMG with compiled Swift launcher.
- Launcher auto-detects terminal preference: Ghostty first, then Alacritty, Kitty, WezTerm, iTerm2, Terminal.app.
- Added app icon assets: `.icns` for macOS, `.ico` for Windows, iOS icon exports.

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
