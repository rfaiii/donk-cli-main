# Project: BVR-CLI
# Date: 2026-09-04

## Current Goal
Ship v1.1.9 — a stability and polish release: the homescreen beaver mascot
no longer spazzes on every mouse movement (4s pulse + 400ms hover debounce),
the duplicate NODE section is removed, the `--version` beaver is now in BVR
neon green (#3BF66B), and the project is organized with `packaging/icons/`
and a single `Taskfile.yaml`.

## Current Technical Status
- **Stack:** Go, Bubble Tea v2, Lip Gloss v2, Fantasy, SQLite/sqlc, Charmwalk snapshots
- **CLI entry:** `main.go` + `internal/cmd`
- **UI:** `internal/ui/...` Bubble Tea state machine
- **Packaging:** `packaging/` for npm/homebrew/icons, `dist/release/` for binaries, `winget/` for manifest templates
- **Release:** GitHub Releases v1.1.9
- **Shaders:** `internal/shader` embeds Ghostty cursor GLSL shaders and provides `bvr-cli shaders list/install`
- **NODE:** `internal/node` provides device registry, HTTP/WebSocket/SSH transports, and tests
- **Themes:** 8 themes, each with derived Accent + Alt colors (`Styles.ThemeColor.Accent/.Alt`) per the color-mapping table
- **Cline:** hosted gateway (api.cline.bot) integrated; live catalog fetch surfaces free models; `X-Api-Key` auth

## Recent Iterations
- [x] Boot brand banner: BVR → version → "OH BEAV!" → attribution, accent-bg scramble sequence (internal/ui/model/version_banner.go)
- [x] Per-theme Accent/Alt color mapping across all 8 themes (internal/ui/styles/themes.go themeAccentAlt)
- [x] Cline free-model support: live catalog fetch + in-app "ADD CLINE API KEY" flow, X-Api-Key verification
- [x] "Other Models" palette entry + Cline catalog dialog (internal/ui/dialog/other_models.go)
- [x] ALT-colored File Finder: metadata, separators, close button (internal/ui/dialog/filebrowser.go)
- [x] Accent CPU/RAM bars; bold underlined project location; stacked Command (terminal icon) and File Finder (folder icon) home buttons, themed to the Primary color
- [x] Lightened section help text (LOCAL DEVICE / SKILLS / MCP) for readability
- [x] Bumped dev version 1.1.6:beta_v2 → 1.1.6:beta_v3 → release 1.1.7
- [x] Added docs/bvr-cli-llm-integration-spec.md and docs/CLINE_PROVIDER.md
- [x] Bumped project version to v1.1.5 and updated README/CHANGELOG
- [x] Added DMG Swift launcher with Ghostty-first terminal auto-detection
- [x] Created GitHub release v1.1.5 and attached cross-platform binaries
- [x] Hardened npm installer with release API fallback and binary validation
- [x] Added Windows installer zip bundle and Winget manifest template
- [x] Merged maintainer/agent configuration into `AGENTS.md`
- [x] Added beta testing guide in `docs/BETA.md`
- [x] Embedded shader set and added `bvr-cli shaders` command
- [x] Added optional Ghostty config installer for bundled shaders
- [x] Added in-TUI cursor animation support using bundled effects
- [x] Added `[+]` attachment button wired to FILEFINDER for in-ecosystem attachments
- [x] Added NODE connection testing and iPhone pairing docs

## Active Roadblocks
- Public GitHub release asset download URLs returning 404 during installer end-to-end testing
- `internal/agent` tests show fixture/API mismatch unrelated to installer work
- Homebrew/npm/winget manifest paths need real install validation after release asset serving is confirmed
- iOS companion app is test-mode only; full device sync flow not yet implemented

## Next Steps
- [ ] Validate Windows EXE installer flow on Windows x64
- [ ] Validate Windows MSI/equivalent packaging flow
- [ ] Complete npm install path against published release
- [ ] Complete Homebrew tap install against published release
- [ ] Complete Winget manifest validation
- [ ] Finalize beta invite list and distribution notes
- [ ] Implement full iPhone ↔ laptop NODE sync command flow
- [ ] Add NODE connection breadcrumbs/onboarding in TUI

---
*Note: Update this before switching models or switching from local development to another machine.*
