# Jank Fixes

This document tracks UI/UX "jank" issues that were identified and fixed in the
DONK TUI, so future contributors can understand the context and avoid regressions.

## 1. LSP "unstarted" status text

**Problem:** User-configured LSP servers that hadn't been triggered yet showed
`unstarted` in the sidebar. This was confusing — "unstarted" sounds like
something is broken.

**Fix:** Changed the status text from `unstarted` to `idle` in
`internal/ui/model/lsp.go`.

LSPs are intentionally started on-demand (when a matching file is opened), so
`idle` is a more accurate description of the configured-but-not-yet-running
state.

## 2. Skill YAML parsing errors (red-dot skills)

**Problem:** User skills with descriptions containing colons followed by spaces
(e.g. `Triggers on: explore the codebase`) failed YAML parsing and showed as
red-dot error states in the sidebar.

**Root cause:** The YAML parser (`gopkg.in/yaml.v3`) treats `key: value` as a
mapping entry. When a description value contains `Triggers on: explore`, the
parser sees the colon as introducing a nested mapping.

**Fix:** Quote description values that contain colons:

```yaml
# Before (breaks)
description: Use the codebase knowledge graph. Triggers on: explore.

# After (works)
description: "Use the codebase knowledge graph. Triggers on: explore."
```

This is a user-side fix — skill authors should quote descriptions that contain
colons. The parser itself is behaving correctly per the YAML spec.

## 3. Banner version display (pseudo-versions)

**Problem:** When DONK is installed via `go install`, the version is set to a
Go pseudo-version like `v0.87.1-0.20260731174531-4d...`. This long string
overflowed the logo banner and looked janky.

**Fix:** Added `version.ShortVersion()` in `internal/version/version.go` which
strips the pseudo-version suffix for display. The function:

- Returns `v0.87.1` for `v0.87.1-0.20260731174531-4d...`
- Preserves prereleases like `v1.0.0-alpha.1`
- Preserves `devel` for local builds

Updated `internal/ui/model/ui.go` to use `version.ShortVersion()` instead of
`version.Version` in both the wide logo and the compact details view.

## 4. Block character logo rendering (font-dependent glyphs)

**Problem:** The wide logo banner renders "DONK" using block characters (█, ▄, ▀).
On terminals whose font does not include these glyphs, the logo renders as
garbled text (e.g. "E69EH") instead of the intended stylized wordmark.

**Fix:** Added a text fallback mode to the logo renderer. When `logo.Opts.Text`
is true, the title is rendered as plain "DONK" / "HYPERDONK" text with the
same gradient styling, rather than block-character letterforms.

- `internal/ui/logo/logo.go`: added `Opts.Text` field; `Render()` now branches
  between block-character and text rendering.
- `internal/ui/model/ui.go`: `renderLogo()` passes `Text: true` so the header
  always uses the readable text variant.

This avoids font-dependent rendering issues while preserving the gradient
styling and diagonal field lines.

## Reproducing the fixes

After making changes, verify with:

```bash
# Build
go build ./...

# Run targeted tests
go test ./internal/version/... ./internal/skills/... ./internal/ui/model/... ./internal/config/...
```
