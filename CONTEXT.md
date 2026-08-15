# Project: DONK-CLI
# Date: 2026-08-14

## Current Goal
Ship and verify the v1.1.5 installer and onboarding experience across:
1) Windows EXE
2) Windows x64 MSI/installer flow
3) NPM package
4) Homebrew formula/tap
5) Winget manifest
6) Optional Ghostty shader installation + standalone shader visualizer

Keep CLI-native Bubble Tea onboarding, Ghostty-first launcher detection, and beta docs updated.

## Current Technical Status
- **Stack:** Go, Bubble Tea v2, Lip Gloss v2, Fantasy, SQLite/sqlc, Charmwalk snapshots
- **CLI entry:** `main.go` + `internal/cmd`
- **UI:** `internal/ui/...` Bubble Tea state machine
- **Packaging:** `packaging/` for npm/homebrew, `dist/release/` for binaries, `winget/` for manifest templates
- **Release:** GitHub Releases v1.1.5 with cross-compiled platform binaries
- **Shaders:** `internal/shader` embeds Ghostty cursor GLSL shaders and provides `donk-cli shaders list/install`

## Recent Iterations
- [x] Rewrote onboarding as screenshot/ASCII-guided tour with Opt Out/Continue
- [x] Bumped project version to v1.1.5 and updated README/CHANGELOG
- [x] Added DMG Swift launcher with Ghostty-first terminal auto-detection
- [x] Created GitHub release v1.1.5 and attached cross-platform binaries
- [x] Hardened npm installer with release API fallback and binary validation
- [x] Added Windows installer zip bundle and Winget manifest template
- [x] Merged maintainer/agent configuration into `AGENTS.md`
- [x] Added beta testing guide in `docs/BETA.md`
- [x] Embedded shader set and added `donk-cli shaders` command
- [x] Added optional Ghostty config installer for bundled shaders

## Active Roadblocks
- Public GitHub release asset download URLs returning 404 during installer end-to-end testing
- `internal/agent` tests show fixture/API mismatch unrelated to installer work
- Homebrew/npm/winget manifest paths need real install validation after release asset serving is confirmed

## Next Steps
- [ ] Validate Windows EXE installer flow on Windows x64
- [ ] Validate Windows MSI/equivalent packaging flow
- [ ] Complete npm install path against published release
- [ ] Complete Homebrew tap install against published release
- [ ] Complete Winget manifest validation
- [ ] Finalize beta invite list and distribution notes
- [ ] Fix `internal/agent` test fixture mismatch

---
*Note: Update this before switching models or switching from local development to another machine.*
