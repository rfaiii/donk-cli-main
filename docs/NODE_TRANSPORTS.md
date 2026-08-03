# NODE Transports

DONK supports one device registry with three transport adapters:

| Scheme | Adapter | Best for |
| --- | --- | --- |
| `http://`, `https://` | HTTP/JSON | Request/response execution and health checks |
| `ws://`, `wss://` | WebSocket/JSON | Persistent connections and streaming output |
| SSH configuration | SSH | Existing computers reachable through SSH |

Start a NODE agent with both HTTP and WebSocket routes enabled:

```sh
donk node serve --host 0.0.0.0:7777 --token "$DONK_NODE_TOKEN"
```

The listener exposes `/v1/node/health`, `/v1/node/capabilities`,
`/v1/node/execute`, and the `/v1/node/ws` WebSocket endpoint. Register HTTP or
WebSocket endpoints with `RegisterNodeTransport`. SSH uses `NewSSHTransport` or
`RegisterSSHTransport`.

Credentials remain in transport objects and are not stored in `Device` or shown
in the UI. The home screen and NODE settings dialog show device name and state:
grey offline, red error, green online.

Open **NODE Connections** from the command palette (`/node`) or press
`ctrl+shift+n`. Use `r` in the dialog to refresh discovery. Configured
transport endpoints are shown for each device so you can verify which node is
selected before running Node/NPM commands.

See `NODE_HTTP_PROTOCOL.md`, `NODE_WEBSOCKET_PROTOCOL.md`, and
`NODE_SSH_TRANSPORT.md` for protocol details.