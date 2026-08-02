# DONK first demo

Goal: prove the hybrid shell works on your machine in under two minutes — no AI required for the UI tour.

## One-liner (Windows)

```powershell
.\scripts\donkcli.ps1
```

Install so you can type `donkcli` anywhere:

```powershell
.\scripts\install-donkcli.ps1
# new terminal:
donkcli
```

Full guided demo (doctor + walkthrough + launch):

```powershell
.\scripts\demo.ps1
```

Doctor only:

```powershell
.\scripts\demo.ps1 -DoctorOnly
```

Alacritty host (if installed):

```powershell
.\scripts\demo.ps1 -Alacritty
```

## One-liner (Unix)

```bash
./scripts/demo.sh
./scripts/demo.sh --doctor-only
```

## Manual

```bash
cargo build -p donk-terminal --bin donk
cargo run -p donk-terminal --bin donk -- doctor
cargo run -p donk-terminal --bin donk -- demo
cargo run -p donk-terminal --bin donk -- --skip-splash
```

## In-shell tour

| Step | Command | Expect |
|------|---------|--------|
| 1 | `/help` | Command list |
| 2 | `/sys` | CPU/mem bars + processes · esc back |
| 3 | `/files` | Dual panes · tab focus · esc |
| 4 | `/read` | Glow sample markdown |
| 5 | `/read ROADMAP.md` | Repo roadmap rendered |
| 6 | `/animations` | tab effects · space on spring |
| 7 | `/setup` | ←/→ change · enter apply |
| 8 | `/nodes` | Local node id |
| 9 | `/place left` | Snap (Windows) |

Chat AI (optional): `ollama serve` → `/models` → ask a question.

Quit: **esc** leaves a tool · **ctrl+c** exits DONK.

## What “good” looks like

- `donk doctor` prints `[ok]` for version, terminal, config, themes, sync
- `[!!] ai` is fine without Ollama — UI demo still passes
- Status bar shows `node/online` (or offline) and theme name
- Tools open/close without crashing the alt-screen

## Known gaps (ok for v2.5.x demo)

- Mobile relay not live yet (`docs/MOBILE.md`)
- Snap uses primary work area only
- Alacritty fork deferred — host pack + scripts instead
- Layout B (S8/S9) and Crush chrome pills still on the [RESOURCES.md](../RESOURCES.md) / ROADMAP queue

Product version for this demo line: **2.6.1** (`donk-cli-v2.6.1`).
