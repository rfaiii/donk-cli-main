# DONK Mobile Companions (Phase 4)

Thin clients — **not** full terminal emulators. They attach to a desktop DONK node.

## Roles
| Client | Stack (planned) | Syncs |
|--------|-----------------|-------|
| iOS / iPad | SwiftUI | Conversations, file history, node status |
| Android | Jetpack Compose | Same |
| Server node | Headless Linux | Relay / API for paired clients |

## Local-first contract
1. Desktop node owns `~/.config/donk/sync/nodes.json` + `sessions.json`
2. Pairing: desktop `/pair` → 6-char code → mobile enters code (trust token issued)
3. Offline mode (`/offline on`): AI stays on local Ollama/mods; hub/relay skipped
4. Manual bridge today: `/sync export path.json` ↔ `/sync import path.json`

## Protocol (next)
- Prefer WebSocket / HTTP/2 relay with device identity
- Attempt direct LAN; fall back to relay
- Encrypt in transit; redact secrets from synced history

## Status in v2.2.0
Identity, peers, pairing codes, session meta store, offline flag, export/import — **shipped as local stubs**. Live mobile apps + network relay = follow-on.
