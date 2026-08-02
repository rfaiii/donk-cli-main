# Reference Index

This catalog lists key paths in `donk-cli-main` for reuse across migration work. Entries include exact paths and short reuse notes.

## Core Entrypoints
- `main.go` — Cobra entrypoint; calls `cmd.Execute()`
- `internal/app/app.go` — Wires services, agents, pubsub brokers, herdr client
- `internal/cmd/root.go` — Root Cobra command setup
- `internal/cmd/run.go` — Run command with client/server streaming

## Services
- `internal/session` — Session service
- `internal/message` — Message service
- `internal/history` — File/history service
- `internal/permission` — Permission service
- `internal/question` — Question service
- `internal/filetracker` — File tracking service
- `internal/db` — SQLite queries + Goose migrations
- `internal/config` — Configuration store and overrides

## UI
- `internal/ui/model/ui.go` — Main Bubble Tea model
- `internal/ui/dialog` — Dialogs: commands, file picker, sessions, models, notifications
- `internal/ui/styles` — Lipgloss/ultraviolet styling
- `internal/ui/chat` — Chat rendering, tools panel, shell, agent output
- `internal/ui/anim` — UI animation helpers and re-exports from `internal/anim`
- `internal/anim` — Go-native animation library: gallery, spinner, spring, progress
- `old/internal/ui/anim/library/rust-reference` — Original Rust `donk-anim` sources for future ports
- `internal/ui/anim/scenes` — Planned Go scene implementations

## Agent & Tools
- `internal/agent` — Coordinator, permissions, questions, prompts, notifications
- `internal/agent/tools` — Tool implementations: bash, edit, fetch, grep, LSP, MCP, web
- `internal/backend` — Backend provider selection
- `internal/commands` — Command palette and custom/MCP prompt loading
- `internal/tools` — Tool registry and slash commands (pending)

## Shell / LSP / Skills
- `internal/shell` — Shell command execution
- `internal/shellconfig` — Shell configuration
- `internal/lsp` — LSP manager
- `internal/skills` — Skill discovery and builtin skills

## OAuth / Server
- `internal/oauth` — OAuth flows: Copilot, Hyper, MCP
- `internal/server` — HTTP/SSE server

## Docs
- `SYSARCH.md` — System architecture overview with diagram
- `docs/MIGRATION-PLAN.md` — Feature migration order and checklist
- `old/docs/INTEGRATION-PLAN.md` — Integration notes for animation, health, FileBrowser/ChromeDB
- `README.md` — Project readme
