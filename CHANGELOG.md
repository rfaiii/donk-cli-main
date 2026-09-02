# Changelog

## 1.1.7

### Homescreen UI & animations

- New boot brand banner plays a sequence at startup (each element held ~2s):
  the wordmark **DONK**, the version (`v1.1.7`), **OH BEAV!**, and the
  attribution `created by RICHARD AIZEN AVERY III`. Letters scramble into place
  with a slow fade; the bar uses the theme **ACCENT** color for its background
  and the app **BACKGROUND** color for its text so it pops
  (`internal/ui/model/version_banner.go`).
- Project location on the homescreen is now **bold**, **underlined**, and
  rendered in the accent color for readability against the dark background.
- Added a prominent **MODEL / PROVIDER** line directly on the homescreen
  (accent values, alt labels).
- Added a `[ "/" OPENS COMMANDS ]` command-palette button on the left side of
  the homescreen, with the `[ OPEN FILE FINDER ]` button immediately adjacent
  on the right; both use the accent color.
- CPU/RAM resource monitors now use the active **ACCENT** color (e.g. pink bars
  on the green theme).
- Lightened the section heading / help text under LOCAL DEVICE, SKILLS, and
  MCP so it reads against the dark surface.

### Global theme & color mapping

- Introduced per-theme **Accent** and **Alt** colors, exposed as
  `Styles.ThemeColor.Accent` / `.Alt`, following the mapping table:
  Green→(Pink,Purple), Pink→(Purple,Blue), Purple→(Blue,Orange),
  Blue→(Orange,White), Orange→(White,Yellow), White→(Yellow,Red),
  Yellow→(Red,Green), Red→(Green,Pink).
- Command buttons use the accent color for text/highlighting and the alt color
  for metadata (shortcut column).

### File Finder

- File Finder metadata, rule separators, and the close button now use the
  active **Alt** color instead of the default red/muted tones.

### Cline provider integration

- Added Cline's hosted API gateway (`api.cline.bot`) as a first-class known
  provider so Cline-hosted models appear in the model picker and onboarding.
- Auth uses the `X-Api-Key` header; the key is read from `$CLINE_API_KEY` or
  the existing `providers.cline.api_key` config entry (in-app API key dialog
  works as with any other provider).
- Default model catalog lives in `internal/config/cline.go` and can be
  extended/overridden per-user via `providers.cline.models` in config.
- Added an "Other Models" entry to the `/` command palette (aliases: other,
  cline, external) that opens a dedicated dialog listing the Cline catalog;
  selecting a model prompts for the Cline API key if none is stored yet.
- The Other Models dialog now includes an "ADD CLINE API KEY" row when no
  Cline key is configured. Selecting it saves the key (verified via the
  `X-Api-Key` header) and immediately loads the full live Cline catalog,
  including Cline's free models, without a restart.
- Cline key verification (`TestConnection`) now authenticates via the
  `X-Api-Key` header so valid Cline keys are accepted.

### Menu reorganization

- Reorganized the main command-palette menu into the documented order: New
  Session, Sessions, Switch Model, Open File Finder, Change Project, Node
  Connections, Ollama Models, Other Models, Ollama How To, Themes,
  Notification Style, Toggle Beast Mode, Toggle Code Mode, Toggle Help,
  Initialize Project, Disable Background Color, Quit.

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
- Added `docs/installation.md` for first-run setup and dependency guidance.
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
