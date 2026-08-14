# DONK-CLI Dependencies

Concrete dependency reference for end users and maintainers.
This covers everything a new user or installer needs on macOS and Windows.

## End-user runtime dependencies

### Required

| Dependency | Purpose | Notes |
| --- | --- | --- |
| DONK binary | Terminal UI runtime | No Go runtime required after build |
| Terminal emulator | Bubble Tea rendering surface | macOS Terminal, Windows Terminal, Alacritty, iTerm2 |
| Network access | Provider auth, model discovery, updates | Required during first run and `/login` |

### Optional

| Dependency | Purpose | Auto-install | Notes |
| --- | --- | --- | --- |
| Ollama | Local models | Guided | Recommended for offline-capable setups |
| Provider API key | Cloud providers | Guided | Anthropic, OpenAI, Gemini, etc. |
| Git | Project detection and workspace state | No | Used when available |
| LSP servers | Code intelligence | No | Optional per language |

## macOS-specific requirements

### User install

- macOS 13.0+ for signed DMG builds.
- Rosetta 2 on Apple Silicon if using Intel-only terminal plugins.
- Internet connection for first-run model/provider setup.

### Developer/build requirements

- Go 1.26.5+
- Xcode Command Line Tools for CGO-neutral builds: `xcode-select --install`
- Optional: `goreleaser` for release builds
- Optional: `wails` or `electron` if building a future desktop wrapper

## Windows-specific requirements

### User install

- Windows 10/11 x64 or ARM64.
- Windows Terminal or another modern terminal host.
- Internet connection for first-run setup.

### Developer/build requirements

- Go 1.26.5+
- MinGW/MSYS2 or Visual Studio Build Tools only if CGO dependencies require it
- Optional: `goreleaser` for release builds
- Optional: `nsis` or `wixtoolset` for installer packaging

## Build dependencies

### Required

- Go 1.26.5+
- Git

### Charm libraries

DONK depends on Charm ecosystem libraries. In the `rfaiii/donk-cli-main`
repo, the expected modules are:

- `charm.land/bubbletea/v2`
- `charm.land/bubbles/v2`
- `charm.land/lipgloss/v2`
- `charm.land/glamour/v2`
- `charm.land/catwalk`
- `charm.land/fantasy`
- `charm.land/fang/v2`
- `charm.land/log/v2`
- `charm.land/x/vcr`
- `github.com/charmbracelet/ultraviolet`
- `github.com/charmbracelet/x/*`

If these are not published under `charm.land`, map them to the local
`CRUSH/` workspace equivalents or vendored copies before building.

### Go workspace note

This repo may expect a shared `go.work` that includes sibling Charm module
paths. If `go build ./...` fails from this repo alone, use one of:

- `go build .` from `donk-cli-main` only
- Add this repo to the shared workspace
- Replace missing `charm.land/*` modules with local paths in `go.mod`

## Ollama setup

### macOS

```sh
brew install ollama
ollama --version
ollama serve &
ollama pull llama3.1:8b
```

### Windows

```powershell
winget install Ollama.Ollama
ollama --version
ollama serve
ollama pull llama3.1:8b
```

## Provider setup

DONK supports multiple providers. Required items per provider:

| Provider | Required | Auto-detect |
| --- | --- | --- |
| Ollama | `ollama` binary + pulled model | Yes |
| Anthropic | API key | No |
| OpenAI | API key | No |
| Gemini | API key or OAuth | No |
| Custom/local | Endpoint + key/model | No |

Use `/login` or the models dialog to complete setup.

## Verification

```sh
donk-cli --version
donk-cli --help
donk-cli dirs
donk-cli projects
```

Then open DONK, run `/models`, refresh, and confirm a model becomes ready.
