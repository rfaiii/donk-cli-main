# Animation Library

This directory holds animation assets and reference ports for `donk-cli-main`.

## Structure
- `rust-reference/` — Original Rust `donk-anim` sources copied for later porting/reference.
- Future Go-native implementations should live under `internal/ui/anim/` or a new `internal/ui/anim/scenes/` package.

## Porting Order
1. `gallery.rs` metadata → Go anim registry
2. `spinner.rs`, `progress.rs` → trivial wins
3. `matrix.rs`, `space.rs` → procedural effects
4. `mission.rs` scenes → `/scenes` command
5. `boot.rs`, `splash.rs` → boot reel
6. `doomfire.rs`, `cycling.rs` → `/animations` tab
