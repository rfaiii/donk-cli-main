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

## Settings & Configuration
- [ ] **Desktop Settings Sync**: Sync generic `options` over the network from the active desktop node.
- [ ] **Mobile-Specific Settings**: Options specific to iOS (Haptics, Push Notifications, Local Storage constraints).
