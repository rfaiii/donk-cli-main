# ◇ DONK-CLI 1.1.5 by richard izn avery

**The keyboard-first AI workspace for your projects.**

DONK brings project files, conversations, local models, tools, permissions,
skills, LSPs, and automation into one fast terminal cockpit. Inspect the
workspace, choose the model, and ship.

> **Release 1.1.5:** this release adds the onboarding wizard with ASCII previews,
> terminal-aware macOS DMG launcher, beta testing docs, updated packaging, and
> optional Ghostty shader installation. The mobile companion scaffold is present
> for testing, but the mobile version is currently in **test mode** and not yet
> ready for general use.

## ⚡ Quick start

See [`docs/installation.md`](docs/installation.md) for all platforms.

```sh
cd /path/to/donk-cli-go
go build -o ./donk-cli .
./donk-cli
```

Install the current source build globally:

```sh
./scripts/install-donk-cli.sh
hash -r 2>/dev/null || true
donk-cli --version
donk-cli --help
```

When testing changes, prefer `./donk-cli`; a bare `donk-cli` can launch an old
binary from another directory.

## 📦 Install

### macOS

Use the DMG for a normal Mac app experience.

1. Open `dist/donk-cli_dev_darwin_arm64.dmg`
2. Drag `DONK.app` to `Applications`
3. Double-click `DONK.app`

If macOS blocks it:
- Right-click `DONK.app` → Open
- Or run: `xattr -cr /Applications/DONK.app`

The DMG launcher auto-detects your terminal:
- Ghostty first, then Alacritty, Kitty, WezTerm, iTerm2
- Falls back to Terminal.app

### Windows

Use the EXE/MSI installer or ZIP archive.

### Linux

Use deb/rpm/apk/arch packages, AUR, Nix, or release archives.

See [`docs/ONBOARDING.md`](docs/ONBOARDING.md) for the full onboarding and
installer guide.

## 🧭 Move between projects

```sh
donk-cli --cwd /path/to/project
```

Inside DONK, use `/cd` or `/project`, select a directory, and press `s`. You can
also type `cd ~/Projects/my-project` into the prompt; DONK switches projects
instead of sending that text to the agent. `donk-cli projects` lists recent
projects.

## 🧠 Local Ollama models

DONK discovers all models available from Ollama’s local API. Open `/models`,
`/ollama`, or press `Ctrl+L`.

- `r` refreshes local models
- `p` pulls the highlighted model
- `s` starts Ollama when it is offline
- `Enter` selects and warms the model

The model diamond is gray when unknown, purple while loading, green when ready,
and red if loading fails. Read [`docs/OLLAMA_HOW_TO.md`](docs/OLLAMA_HOW_TO.md).

![DONK models](resources/screenshots/menu-models.jpg)

For smaller Ollama coding models such as `qwen2.5-coder:3b-instruct`, use the
fantasy-free native coder. It talks directly to Ollama with a small tool set,
recovers Qwen-style text tool calls, and buffers raw tool-call JSON from the
terminal display:

```sh
GOWORK=off go build -o "$HOME/bin/codetool" ./cmd/codetool
GOWORK=off go build -o "$HOME/bin/donk-cli" .

DONK_CODETOOL="$HOME/bin/codetool" \
  "$HOME/bin/donk-cli" code --native \
  --model qwen2.5-coder:3b-instruct --stream \
  "Read the relevant files and fix the failing test"
```

Use `--cwd /path/to/project` when the coding task is in another directory. See
[`docs/OLLAMA_HOW_TO.md`](docs/OLLAMA_HOW_TO.md) for diagnostics and
troubleshooting.

## ✨ What DONK includes

- Branded Bubble Tea TUI with bounded, responsive layouts
- Bottom resource bar showing live CPU and RAM usage
- Project File Finder with previews, metadata, hidden files, paging, clipboard,
  and project switching
- `[+]` attachment button wired to FILEFINDER for fast file, photo, video, and
  link attachments
- Sessions, model/provider selection, permissions, MCP, LSP, and skills
- Ollama discovery, model pull, automatic startup, and model warm-up
- Agent Skills from `~/.agents/skills` with master-catalog syncing
- NODE connection support for HTTP/JSON, WebSocket, and SSH transports
- Optional embedded Ghostty shader installation and in-TUI cursor animations
- Persistent project registration and recently accessed project tracking

![DONK home](resources/screenshots/home-menu.jpg)
![DONK file finder](resources/screenshots/file-finder.jpg)

## 🛠 Developer checks

```sh
gofmt -w .
go test ./...
go vet ./...
go build ./...
```

Cross-compile the beta targets:

```sh
GOOS=darwin GOARCH=arm64 go build -o /tmp/donk-darwin-arm64 .
GOOS=darwin GOARCH=amd64 go build -o /tmp/donk-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o /tmp/donk-windows-amd64.exe .
GOOS=linux GOARCH=amd64 go build -o /tmp/donk-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o /tmp/donk-linux-arm64 .
```

## 📦 Install

- **macOS** — use the DMG, Homebrew, or NPM wrapper.
- **Windows** — use the EXE/MSI installer or ZIP archive.
- **Linux** — use deb/rpm/apk/arch packages, AUR, Nix, or release archives.

See [`docs/ONBOARDING.md`](docs/ONBOARDING.md) for the full onboarding and
installer guide.

## 📚 Documentation map

| Document | Use it for |
| --- | --- |
| [`docs/installation.md`](docs/installation.md) | Installation, testing, paths, troubleshooting |
| [`docs/ONBOARDING.md`](docs/ONBOARDING.md) | First-run setup, dependencies, and packaging |
| [`docs/DEPENDENCIES.md`](docs/DEPENDENCIES.md) | Concrete macOS/Windows dependency list |
| [`docs/IMG-RESOURCES.md`](docs/IMG-RESOURCES.md) | Screenshots, videos, and app icon assets |
| [`docs/PACKAGING.md`](docs/PACKAGING.md) | Release automation and packaging configs |
| [`docs/ABOUT.md`](docs/ABOUT.md) | Company and maintainer contact info |
| [`docs/BETA.md`](docs/BETA.md) | Beta testing program, invites, distribution, and feedback |
| [`docs/OLLAMA_HOW_TO.md`](docs/OLLAMA_HOW_TO.md) | Ollama setup and local models |
| [`docs/SKILLS.md`](docs/SKILLS.md) | Agent Skills discovery and syncing |
| [`docs/FILE_FINDER.md`](docs/FILE_FINDER.md) | File Finder behavior |
| [`docs/SYSTEM_ARCHITECTURE.md`](docs/SYSTEM_ARCHITECTURE.md) | Architecture and roadmap |
| [`docs/UI_BRANDING.md`](docs/UI_BRANDING.md) | DONK visual identity |
| [`docs/MOBILE-CLI.md`](docs/MOBILE-CLI.md) | Mobile companion bridge and iOS/Android plans |
| [`docs/DONK-SERVER.md`](docs/DONK-SERVER.md) | Host-side companion server design |
| [`docs/NODE_TRANSPORTS.md`](docs/NODE_TRANSPORTS.md) | NODE connection setup and transports |
| [`docs/NODE_HTTP_PROTOCOL.md`](docs/NODE_HTTP_PROTOCOL.md) | NODE HTTP/JSON protocol |
| [`docs/NODE_WEBSOCKET_PROTOCOL.md`](docs/NODE_WEBSOCKET_PROTOCOL.md) | NODE WebSocket streaming protocol |
| [`docs/NODE_SSH_TRANSPORT.md`](docs/NODE_SSH_TRANSPORT.md) | NODE SSH transport |

![DONK notifications](resources/screenshots/notification-select.jpg)
![DONK pink theme](resources/screenshots/theme-pink.jpg)
![DONK purple theme](resources/screenshots/theme-purple.jpg)
![DONK default green theme](resources/screenshots/home-green.jpg)

## 🧪 Beta report checklist

Include:

1. OS and CPU architecture
2. Terminal application
3. `donk-cli --version` output
4. Exact launch command
5. Project path and whether Ollama was enabled
6. Reproduction steps and terminal output

**See the project. Choose your intelligence. Get to work.**
