# Mobile CLI

Plans and outline for the iPhone and Android companion variants, plus the
host-side companion server that bridges them to `donk-cli`.

## Goals

- Provide a mobile-first control surface for Donk CLI sessions.
- Keep the existing Go CLI intact; add a lightweight companion bridge instead of rewriting the core in another runtime.
- Support secure local/network transport between the host CLI and mobile clients.

## Source of Truth

This plan is based on:
`/Users/richavery/.hermes/attachments/donk_cli_companion_architecture.html`

Companion server design details live in:
[`docs/DONK-SERVER.md`](DONK-SERVER.md)

## Topology

- Host machine runs the normal Donk CLI.
- Optional companion server/bridge runs alongside the CLI to expose secure mobile-facing commands.
- Mobile clients connect via a private/link-local or authenticated channel.

## Essentials

1. Host-side bridge module
2. Mobile client variant(s)
3. Secure auth/transport model
4. Session and telemetry bridge
5. Installation and docs

## Proposed Structure

```
mobile-cli/
  README.md
  proto.md
  ios/
  android/
```

## Roadmap

1. Capture requirements and constraints.
2. Define the bridge protocol and data model.
3. Prototype the host-side entrypoint.
4. Prototype one mobile client variant.
5. Add tests, docs, and packaging.

## Current status

- iOS companion app scaffolded in `mobile-cli/ios/`.
- Companion server design documented in `docs/DONK-SERVER.md`.
- Host-side bridge module has not been implemented yet.

## Next Steps

- Implement the companion server in Go.
- Wire mobile client discovery/connection to the host bridge.
- Lock down the minimum viable command set.
