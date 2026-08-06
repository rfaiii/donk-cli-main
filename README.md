# ◇ DONK-CLI 1.1.3

**The keyboard-first AI workspace for your projects.**

DONK brings project files, conversations, local models, tools, permissions,
skills, LSPs, and automation into one fast terminal cockpit. Inspect the
workspace, choose the model, and ship.

> **Release 1.1.3:** this release adds a native Ollama coding path for smaller
> local models and a bottom resource bar that shows live CPU and RAM usage.

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
- Sessions, model/provider selection, permissions, MCP, LSP, and skills
- Ollama discovery, model pull, automatic startup, and model warm-up
- Agent Skills from `~/.agents/skills` with master-catalog syncing
- NODE connection foundations for HTTP/JSON, WebSocket, and SSH transports
- Persistent project registration and recently accessed project tracking

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

## 📚 Documentation map

| Document | Use it for |
| --- | --- |
| [`docs/installation.md`](docs/installation.md) | Installation, testing, paths, troubleshooting |
| [`docs/OLLAMA_HOW_TO.md`](docs/OLLAMA_HOW_TO.md) | Ollama setup and local models |
| [`docs/SKILLS.md`](docs/SKILLS.md) | Agent Skills discovery and syncing |
| [`docs/FILE_FINDER.md`](docs/FILE_FINDER.md) | File Finder behavior |
| [`docs/SYSTEM_ARCHITECTURE.md`](docs/SYSTEM_ARCHITECTURE.md) | Architecture and roadmap |
| [`docs/UI_BRANDING.md`](docs/UI_BRANDING.md) | DONK visual identity |

## 🧪 Beta report checklist

Include:

1. OS and CPU architecture
2. Terminal application
3. `donk-cli --version` output
4. Exact launch command
5. Project path and whether Ollama was enabled
6. Reproduction steps and terminal output

**See the project. Choose your intelligence. Get to work.**