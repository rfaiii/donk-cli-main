# Features

## Current
- **Multi-Model:** choose from a wide range of LLMs or add your own via OpenAI- or Anthropic-compatible APIs
- **Flexible:** switch LLMs mid-session while preserving context
- **Session-Based:** maintain multiple work sessions and contexts per project
- **LSP-Enhanced:** DONK CLI uses LSPs for additional context, just like you do
- **Extensible:** add capabilities via MCPs (`http`, `stdio`, and `sse`)
- **Node Integration:** built-in NODE device discovery, status tracking, and npm script execution from the command palette
- **File Finder:** improved two-pane file browser with fixed pane widths, scrollbar track, and control-safe rendering
- **Command Palette:** quick access to system commands, npm scripts, and custom actions
- **Works Everywhere:** first-class support in every terminal on macOS, Linux, Windows (PowerShell and WSL), Android, FreeBSD, OpenBSD, and NetBSD

## Planned
- **Tool Registry:** centralized slash-command and tool-panel registry in `internal/tools`
- **Health Monitor:** `/health` service checks and resource monitor in Go
- **Ollama / Local Models:** first-class local-model backend support in agent/tools
- **System Monitor:** terminal-native system resource panel
- **Remote Node Discovery:** expand NODE discovery beyond localhost ports to remote IPs
- **Palette-Native Node Flow:** direct palette invocation of node/npm tools with device selection
- **NPM Script Runner:** device-scoped npm script execution with output routing to chat
