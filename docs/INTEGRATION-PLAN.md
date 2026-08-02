# DONK Integration Plan — Go-Backbone Migration

## Clarification
- **UI backbone**: `donk-cli-go` Bubble Tea app is the complete UI. No Rust chrome, no Alacritty app, no external terminal wrapper.
- **Migration goal**: Port features/assets/algorithms/data from the old Rust/Python project into the existing Go app. Do not rebuild UI chrome.
- **Retain**: Go `FileBrowser`, theme system, Bubble Tea layout, existing slash commands, agent tools, SQLite persistence.
- **Discard**: Rust `chrome.rs`, Ratatui renderers, crate scaffolding, Python runtime, Alacritty host scripts.

---

## 1) Animation Library — Status

### Completed
- ✅ Created `internal/ui/anim/library/` with all Rust `donk-anim` sources in `rust-reference/`
- ✅ Documented porting order and priority
- ✅ Committed/pushed as `cee40bc`

### Port Strategy
Port algorithms/data into Go `internal/ui/anim/` or `internal/ui/anim/scenes/`. Do not recreate Rust chrome/UI wrappers. Use Bubble Tea `tea.Cmd`/`tea.Msg` only.

| Rust source | Go target | Priority | Effort |
|-------------|-----------|----------|--------|
| `gallery.rs` | `anim/registry.go` | High | Low — data only |
| `spinner.rs` | `anim/spinner.go` | High | Low — trivial |
| `progress.rs` | `anim/progress.go` | High | Low — trivial |
| `matrix.rs` | `anim/matrix.go` | Medium | Medium |
| `space.rs` | `anim/space.go` | Medium | Medium |
| `mission.rs` | `scenes/` package | Medium | Medium |
| `boot.rs` | `anim/boot.go` | Medium | Medium |
| `splash.rs` | `anim/splash.go` | Medium | Medium |
| `doomfire.rs` | `anim/doomfire.go` | Low | Low |
| `cycling.rs` | `anim/cycling.go` | Low | Low |
| `eyes.rs` | `anim/eyes.go` | Low | Low |
| `vanish.rs` | `anim/vanish.go` | Low | Low |
| `showcase.rs` | `anim/showcase.go` | Low | Low |
| `registry.rs` | Reference only | — | — |
| `lib.rs` | Reference only | — | — |

### Next Steps
1. Port `gallery.rs` metadata → Go anim registry
2. Port `spinner.rs` + `progress.rs` → Go `internal/ui/anim/`
3. Wire `/animations` tab UI to use Go spinner
4. Port `matrix.rs`, `space.rs` → background effects
5. Port `mission.rs` scenes → `/scenes` command

---

## 2) `/health` Resource Monitor

### Current State
- Rust repo: placeholder renderer only
- Go base: `/v1/health` readiness endpoint exists for server self-check
- No service registry or historical health data

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
1. Add `internal/health/` package in `donk-cli-main`
2. Define `ServiceConfig` struct for `crush.json`
3. Implement checker interface: `Check(ctx) (Status, error)`
4. Store history in SQLite via existing `internal/db`
5. Add `/health` slash command + Bubble Tea view
6. Add `/health` JSON API endpoint
7. Wire alerts to notification system

### No Blockers
- SQLite persistence: ✅
- HTTP client: ✅
- Shell exec: ✅
- Notification system: ✅
- Bubble Tea view plumbing: ✅

---

## 3) FileBrowser / Superfile

### Current State
- Go `FileBrowser` dialog at `internal/ui/dialog/filebrowser.go`
- Superfile-style dual-pane UX partially implemented
- No `/files` slash command yet; has `/finder` and `ctrl+shift+f`

### What To Keep
- Go `FileBrowser` struct and Bubble Tea dialog integration
- Preview pane, metadata, clipboard support
- Command palette integration (`/finder`, `ctrl+shift+f`)

### What To Polish
- Add `/files` slash alias alongside `/finder`
- Enhance `FileBrowser.Draw()` with theme-aware styling
- Add dual-pane mode toggle if valuable
- Add service browser sub-mode later if ChromeDB is implemented

---

## 4) Migration Task List

### Immediate
- [ ] Port `gallery.rs` metadata → Go anim registry
- [ ] Port `spinner.rs` + `progress.rs` → Go `internal/ui/anim/`
- [ ] Add `/files` slash alias
- [ ] Document ChromeDB schema proposal

### Phase 1
- [ ] Port `matrix.rs`, `space.rs`
- [ ] Polish `FileBrowser.Draw()` with theme-aware styling
- [ ] Implement `/health` package with 3 check types
- [ ] Add `/health` slash command + view

### Phase 2
- [ ] Port `mission.rs` scenes → `/scenes` command
- [ ] Add dual-pane mode to `FileBrowser`
- [ ] Implement `/health` SQLite history
- [ ] Wire `/health` mini-view into `/files`

### Phase 3
- [ ] Port `boot.rs`, `splash.rs` → boot reel
- [ ] Implement `/animations` tab UI
- [ ] Add ChromeDB service registry
- [ ] Add service browser sub-mode to `/files`

### Phase 4
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
