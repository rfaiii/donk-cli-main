# Mobile Companion Parity Tasklist

The goal is to achieve 1:1 feature parity with the Desktop CLI application by migrating existing features, settings, and workflows to the iOS Companion app.

## Core Interaction
- [ ] **Command Palette (`/`)**: Intercept `/` in the ChatInputBar to open an overlay displaying available commands.
- [ ] **Command Parser**: Handle command execution inside the text field (e.g. `/themes`, `/node`).

## Desktop-to-Mobile Features
- [ ] **Themes (`/themes`)**: Allow switching colors (Material, Charmtone, etc). Mobile-specific logic for applying SwiftUI color palettes.
- [ ] **Ollama Models (`/ollama`)**: Browse, switch, and manage local models.
- [ ] **NODE Connections (`/node`)**: Interface to manage desktop host connections over WebSockets/Bonjour.
- [ ] **File Finder (`/finder`)**: Navigate the filesystem (either locally on mobile, or via the active desktop node).
- [ ] **Change Project (`/cd`)**: Set the active working directory context.
- [ ] **MCP Manager (`/mcp`)**: Add, list, and remove MCP server configurations.

## NODE Connection Testing
- [ ] **Localhost probe test**: verify device discovery finds `127.0.0.1` and link-local listeners.
- [ ] **HTTP health check test**: validate `/v1/node/health` response and bearer auth behavior.
- [ ] **WebSocket streaming test**: validate JSON-RPC execute frames and reconnect logic.
- [ ] **SSH transport test**: validate host-key verification and command execution quoting.
- [ ] **iPhone pairing smoke test**: validate companion app connects to `donk node serve` and appears online in NODE Connections.
- [ ] **Offline/online state test**: validate UI state transitions when the host stops/starts.

## Settings & Configuration
- [ ] **Desktop Settings Sync**: Sync generic `options` over the network from the active desktop node.
- [ ] **Mobile-Specific Settings**: Options specific to iOS (Haptics, Push Notifications, Local Storage constraints).

## Advanced UI & Animations
- [ ] **Web/HTML Parity**: Continuously sync the CSS/JS web prototype so it looks 1:1 with the iOS app.
- [ ] **Live Text Glitch Effect**: Add a subtle text glitch or typing animation under the banner (matching Desktop CLI aesthetics).
- [ ] **Dynamic Backgrounds**: Investigate and implement subtle animated backgrounds (e.g., slow gradient shifts or particle effects) natively on iOS/Android.
- [ ] **Interactive Command Palette**: Evolve the command palette from simple text insertion to interactive sub-menus (e.g., `/themes` opens a picker, context-sensitive onboarding options).
