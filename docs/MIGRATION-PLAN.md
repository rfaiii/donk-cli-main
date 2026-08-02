# DONK Migration Plan — Feature-by-Feature Go Port

## Goal
Migrate all portable features/assets from the old Rust/Python `donk-cli` into the Go master project `donk-cli-main`, one feature at a time. Keep the Go Bubble Tea app as the complete UI backbone. No Rust chrome, no Alacritty app, no external terminal wrapper.

## Migration Order

### Phase 1: Animation Library
- [x] Create animation reference library (`internal/ui/anim/library/rust-reference/`)
- [x] Port gallery metadata (`gallery.rs` → `internal/ui/anim/gallery.go`)
- [x] Port spinner presets (`spinner.rs` → `internal/ui/anim/spinner.go`)
- [x] Port spring physics (`spring.rs` → `internal/ui/anim/spring.go`)
- [x] Port progress bar (`progress.rs` → `internal/ui/anim/progress.go`)
- [ ] Create `internal/anim/` library folder for Go-native animation code
- [ ] Port `matrix.rs` → `internal/anim/matrix.go`
- [ ] Port `space.rs` → `internal/anim/space.go`
- [ ] Port `mission.rs` scenes → `internal/anim/scenes/`
- [ ] Wire animations into UI

### Phase 2: Tools & Commands
- [ ] Create `internal/tools/` library folder
- [ ] Port tool definitions from Rust `tools.rs`
- [ ] Port slash command parser
- [ ] Add `/files`, `/health`, `/sys`, `/read`, `/setup`, `/scenes`, `/animations` commands
- [ ] Wire commands into command palette

### Phase 3: Health Monitor
- [ ] Create `internal/health/` package
- [ ] Implement service checkers
- [ ] Add `/health` slash command + view
- [ ] Wire into SQLite history

### Phase 4: Ollama Integration
- [ ] Review existing Ollama discovery/config
- [ ] Expand local model capabilities
- [ ] Add Node.js integration if needed

### Phase 5: Node.js Integration
- [ ] Review existing Node dependencies
- [ ] Add Node-based tools/scripts
- [ ] Wire into agent tool system

## Current State
- Go master: `/Users/richavery/Projects/donk-cli-main`
- Rust reference: `/Users/richavery/Projects/donk-cli`
- Animation library: `internal/ui/anim/library/rust-reference/`
- Ported animations: `internal/ui/anim/{gallery,spinner,spring,progress}.go`
