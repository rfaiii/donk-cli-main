# DONK Tasklist

## v1.1.7
- [x] Boot brand banner: DONK → version → "OH BEAV!" → attribution scramble sequence (fades away on completion)
- [x] Per-theme Accent/Alt color mapping across all 8 themes (`Styles.ThemeColor.Accent`/`.Alt`)
- [x] Cline hosted-gateway free-model support (live catalog fetch + in-app "ADD CLINE API KEY" flow)
- [x] "Other Models" palette entry + Cline catalog dialog
- [x] ALT-colored File Finder: metadata, separators, close button
- [x] Accent CPU/RAM bars; bold underlined project location
- [x] Home-screen Command and File Finder buttons with superfile icon glyphs, themed to the Primary color, stacked with a cushion line
- [x] Stacked NODE/MCP/LSP resource monitors (left) with SKILLS (right) for tighter space use
- [x] Lightened section and status-bar help text for readability; restored `Status.Help` foreground
- [x] Reorganized command-palette menu into documented order
- [x] Restored the "CLI" suffix in the header wordmark (`DONK` → `DONK-CLI`, per `docs/UI_BRANDING.md`)
- [x] Click-away to dismiss the File Finder and the Commands dialog (outside click / `esc` / `[X]`)
- [x] Truncate long File Finder and SKILLS entries with an ellipsis instead of wrapping
- [x] Fixed project-location rendering (no garbled/numeric output) by styling the path once
- [x] Fixed boot banner getting stuck on the 34-rune attribution line so it always fades
- [ ] Docs: refresh `docs/TASKLIST.md` to this version (this file)

## In progress
- [ ] NODE localhost probe test
- [ ] NODE HTTP health check and bearer auth test
- [ ] NODE WebSocket streaming and reconnect test
- [ ] NODE SSH transport and host-key verification test
- [ ] iPhone pairing smoke test against `donk node serve`
- [ ] Offline/online NODE UI state transition test

## Next
- [ ] Implement laptop ↔ iPhone NODE sync/connect flow
- [ ] Update onboarding breadcrumbs for NODE connections
- [ ] Complete mobile companion parity items
- [ ] Deep repo cleanup and docs consolidation
- [ ] Desktop packaging: produce clean `.exe` (Windows), `.dmg` (macOS), and Linux executables; test across Windows, macOS, and mobile
- [ ] Animated ASCII beavers + micro-animations on events (send, error, finish, etc.)
- [ ] Sound-effects library wired into the TUI (completion chime, error tone, etc.)

## Mobile Companion Parity
- [ ] Command palette `/` overlay
- [ ] Command parser `/themes`, `/node`, `/finder`, `/mcp`, `/cd`
- [ ] Theme switching
- [ ] Ollama model management
- [ ] File Finder
- [ ] Desktop settings sync
- [ ] Mobile-specific settings

