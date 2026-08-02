# DONK Integration Plan — Animations, Health, FileBrowser

## Clarifications

### Rust Chrome ≠ ChromeDB
- **Rust chrome**: UI visual styling/frame work from `crates/donk-tui/src/chrome.rs` — borders, sidebars, headers, status bars, color palettes, slash chips, keyboard footers. This is **presentation only**, not a database.
- **ChromeDB**: A proposed service-registry/database layer for `/health` and service metadata. This is **data/storage**, not UI chrome.
- These are separate concerns. The Go `FileBrowser` can adopt Rust-inspired chrome styling without touching ChromeDB.

---

## 1) Animation Library — Status

### Completed
- ✅ Created `internal/ui/anim/library/` in `donk-cli-main`
- ✅ Copied all Rust `donk-anim/src/*.rs` files into `internal/ui/anim/library/rust-reference/`
- ✅ Added README with porting order
- ✅ Committed and pushed as `cee40bc`

### What’s Stored For Later Port
| Rust source | Animation | Priority | Notes |
|-------------|-----------|----------|-------|
| `gallery.rs` | Gallery metadata | High | Data-only, easiest port |
| `spinner.rs` | Spinner | High | Trivial Go port |
| `progress.rs` | Progress bar | High | Trivial Go port |
| `matrix.rs` | Matrix rain | Medium | Procedural columns |
| `space.rs` | Space field | Medium | Procedural stars |
| `mission.rs` | 8 scenes | Medium | Algorithms portable |
| `boot.rs` | Boot sequence | Medium | ASCII frames |
| `splash.rs` | Splash intro | Medium | Timing dependent |
| `doomfire.rs` | Doom fire | Low | Procedural palette |
| `cycling.rs` | Char cycling | Low | ASCII effect |
| `eyes.rs` | Eyes | Low | ASCII art dependent |
| `vanish.rs` | Vanish | Low | ASCII dependent |
| `showcase.rs` | Showcase | Low | Depends on gallery |
| `registry.rs` | Registry | Reference | Go version will differ |
| `lib.rs` | Module glue | Reference | Go uses packages |

### Next Steps
1. Port `gallery.rs` metadata → Go anim registry
2. Port `spinner.rs` + `progress.rs` → Go `internal/ui/anim/`
3. Wire `/animations` tab UI to use Go spinner
4. Port `matrix.rs` → Go background effect
5. Port `mission.rs` scenes → `/scenes` command

---

## 2) `/health` Resource Monitor Plan

### Current State
- Rust repo: placeholder renderer only (`render_health_placeholder`)
- Go base: has `/v1/health` readiness endpoint for server self-check
- No service registry or historical health data in either repo

### Proposed Scope
| Check | Implementation | Data store |
|-------|---------------|------------|
| Process alive | `ps` / `pgrep` | SQLite |
| Port listening | `lsof` / TCP dial | SQLite |
| HTTP endpoint | `http.Get` with timeout | SQLite |
| DNS resolve | `net.LookupHost` | SQLite |
| Disk usage | `syscall.Statfs` | SQLite |
| Server self-check | `GET /v1/health` | In-memory |

### Integration Plan
1. Add `health/` package under `internal/` in `donk-cli-main`
2. Define `ServiceConfig` struct for `crush.json`
3. Implement checker interface: `Check(ctx) (Status, error)`
4. Store history in SQLite via existing `internal/db`
5. Add `/health` slash command + Bubble Tea view
6. Add `/health` JSON API endpoint if needed
7. Wire alerts to notification system

### No Blockers Identified
- SQLite persistence: ✅ exists
- HTTP client: ✅ exists
- Shell exec: ✅ exists
- Notification system: ✅ exists
- Bubble Tea view plumbing: ✅ exists

---

## 3) FileBrowser / Superfile Integration Plan

### Current State
- Go `FileBrowser` dialog exists at `internal/ui/dialog/filebrowser.go`
- Superfile-style dual-pane UX partially implemented
- No `/files` slash command yet; has `/finder` and `ctrl+shift+f`
- Rust repo has `/files` dual-pane chrome polish we can reference

### What To Keep
- Go `FileBrowser` struct and Bubble Tea dialog integration
- Preview pane, metadata, clipboard support
- Command palette integration (`/finder`, `ctrl+shift+f`)

### What To Port / Polish
- Rust `/files` dual-pane chrome styling → Go `FileBrowser` Draw()
- Rust slash chip `/files` → add `/files` alias alongside `/finder`
- Rust file preview UX → enhance Go preview pane
- Rust clipboard integration → Go clipboard package

### Chrome vs ChromeDB Clarification
- **Chrome**: UI styling for the file browser panels
- **ChromeDB**: Optional service-registry DB layer for `/health` and `/files` service browser mode
- These can be done independently; ChromeDB is not required for FileBrowser polish

### Proposed Integration Steps
1. Add `/files` slash alias to command palette
2. Polish `FileBrowser.Draw()` with Rust-inspired chrome styling
3. Add dual-pane mode toggle in `FileBrowser`
4. Add service browser sub-mode if ChromeDB is implemented
5. Wire `/health` mini-view into `/files` sidebar

---

## 4) Migration Task List

### Immediate
- [ ] Port `gallery.rs` metadata → Go anim registry
- [ ] Port `spinner.rs` → Go `internal/ui/anim/`
- [ ] Add `/files` slash alias
- [ ] Document ChromeDB schema proposal

### Phase 1 (Week 1)
- [ ] Port `progress.rs`, `matrix.rs`, `space.rs`
- [ ] Polish `FileBrowser.Draw()` with chrome styling
- [ ] Implement `/health` package with 3 check types
- [ ] Add `/health` slash command + view

### Phase 2 (Week 2)
- [ ] Port `mission.rs` scenes → `/scenes` command
- [ ] Add dual-pane mode to `FileBrowser`
- [ ] Implement `/health` SQLite history
- [ ] Wire `/health` mini-view into `/files`

### Phase 3 (Week 3)
- [ ] Port `boot.rs`, `splash.rs` → boot reel
- [ ] Implement `/animations` tab UI
- [ ] Add ChromeDB service registry
- [ ] Add service browser sub-mode to `/files`

### Phase 4 (Week 4)
- [ ] Port `doomfire.rs`, `cycling.rs`
- [ ] Polish all animations with theme colors
- [ ] Update docs to `donk-cli-main`
- [ ] Create release

---

## 5) Open Decisions
- ChromeDB schema: extend existing SQLite or new table?
- `/files` vs `/finder`: keep both or consolidate?
- Animation tab: standalone `/animations` or embed in chat?
- Health history retention: 7 days? 30 days? unlimited?
