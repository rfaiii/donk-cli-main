# TaskList — macOS support (native Apple Silicon + Intel)

Use this as the implementation/verification checklist for a first-class macOS release.
Keep items concrete, testable, and minimal; avoid adding frameworks without a caller.

## 1. Build/runtime prereqs
- [x] Repo builds on macOS with a standard Rust toolchain
- [x] Non-interactive build does not depend on stale `TMPDIR`
- [x] Duplicate binary target warning removed from `crates/donk-cli/Cargo.toml`
- [ ] Add explicit `x86_64-apple-darwin` test target to CI or local verification matrix
- [ ] Add explicit `aarch64-apple-darwin` test target to CI or local verification matrix
- [ ] Document minimum macOS version support once known (e.g. `macos-version-min`)

## 2. Host detection + launch
- [x] `TerminalInfo::detect()` classifies Terminal.app / iTerm2 / Alacritty / kitty / WezTerm / Ghostty on macOS
- [x] `/host write` rewrites bundled Alacritty host template with actual donk binary path
- [x] `/host launch` on macOS tries Alacritty → kitty → iTerm2 → WezTerm → Ghostty
- [x] `/host launch` in `donk-tui` calls `launch_host_mac()` on macOS instead of Alacritty-only path
- [ ] `/host write` + `/host launch` paths are covered by at least one real launch test per host on macOS
- [ ] `doctor` reports *expected* terminal class in interactive sessions, not just `unknown` artifacts

## 3. Framework additions / integrations needed on macOS
- [x] `core-foundation` / `objc2-core-foundation` usage audit:
  - [x] window restore/snap only works on Windows today—macOS path now falls back cleanly instead of erroring
- [ ] `portable-pty` fork/session behavior is tested on macOS PTY semantics
- [ ] Graphics/shader work is clearly gated on supported terminals; macOS fallback path is non-panicking
- [ ] Confirm Ghostty shader/theme install path wiring works on macOS (`~/Library/Application Support/Ghostty/...`)
- [ ] Confirm Alacritty config directories are handled on macOS:
  - [x] `~/.config/alacritty/` exists on this Mac
  - [ ] macOS-store Alacritty config path if different

## 4. Native UX on macOS
- [ ] Confirm Command key handling mapping is intentional; do not silently override typical macOS shortcuts
- [ ] Verify copy/paste/select works in interactive TUI on macOS terminals
- [ ] High-DPI / font scaling sanity check for bundled/display fonts
- [ ] Accessibility/voice-over is noted as “unsupported” or implemented explicitly

## 5. Distribute / install
- [x] `scripts/run-alacritty-host.sh` is canonical macOS host launcher with safe quoting/validation
- [ ] Add macOS install docs:
  - [ ] Homebrew tap or release artifact install path
  - [ ] `brew install --cask` candidate if packaging is desired
- [ ] Gate/Codesign decision recorded:
  - [ ] Notarized app vs unsigned CLI tool
  - [ ] If GUI host helper is added later, account for sandbox/notarization cost

## 6. Testing matrix
- [ ] Interactive smoke test:
  - [ ] Terminal.app
  - [ ] iTerm2
  - [ ] Alacritty
  - [ ] kitty
  - [ ] WezTerm
  - [ ] Ghostty if installed
- [x] Non-interactive smoke test:
  - [x] `TERM=dumb` / CI-style `donk doctor` exit is stable and informative
- [x] Build verification:
  - [x] `cargo check --workspace`
  - [x] `cargo build --release -p donk-terminal -p donk-cli`
  - [ ] `cargo test --workspace` once tests exist

## 7. Intel compatibility
- [ ] Verify CI or local `x86_64-apple-darwin` builds succeed
- [ ] Verify native `aarch64-apple-darwin` builds succeed
- [ ] Document whether universal binary is a goal:
  - [ ] If yes: add `lipo`/`universal2` packaging steps
  - [ ] If no: document supported architectures clearly

## 8. Crash/no-signal safety
- [ ] Confirm macOS crash-reporter/reopen behavior is acceptable
- [ ] Ensure `donk` can be relaunched from host script without zombie PTY sessions
- [ ] Confirm signal handling on macOS matches intended shutdown behavior

## 9. Docs + release hygiene
- [x] `README.md` shows macOS host docs and install commands
- [x] `ROADMAP.md` references macOS as part of shipped hybrid host surface
- [x] `ARCHITECTURE.md` host model section reflects Alacritty/kitti/iTerm2/WezTerm/Ghostty fallback behavior
- [x] `REFERENCE.md` keeps exact macOS paths for host config and runtime dirs

## 10. Done criteria
- [ ] Full interactive smoke test on at least one shipped macOS host is green
- [x] `donk doctor` is informative on macOS and passes in `TERM=dumb`
- [x] `/host status|write|launch` has working macOS coverage
- [ ] Builds clean on both Intel and Apple Silicon targets
- [x] Docs updated and reviewed
