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

Use `donk-cli` as the binary name and command on every platform.

- **macOS** — DMG, Homebrew, NPM, or release archive.
- **Windows** — EXE, MSI, ZIP, NPM, or winget.
- **Linux** — release archive, NPM, or Homebrew on Linux.

After install, verify with:

```sh
donk-cli --version
donk-cli --help
```

See [`docs/installation.md`](docs/installation.md) for platform-specific steps,
troubleshooting, and package-manager commands.

## 📚 Documentation map

| Document | Use it for |
| --- | --- |
| [`docs/installation.md`](docs/installation.md) | Installation, testing, paths, troubleshooting |
| [`docs/ONBOARDING.md`](docs/ONBOARDING.md) | First-run setup, dependencies, and packaging |
| [`docs/IMG-RESOURCES.md`](docs/IMG-RESOURCES.md) | Screenshots, videos, and app icon assets |
| [`docs/PACKAGING.md`](docs/PACKAGING.md) | Release automation and packaging configs |
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
| [`docs/TASKLIST.md`](docs/TASKLIST.md) | Active task tracking |

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

## 📬 Contact

- **Maintainer:** Richard Aizen Avery III
- **Email:** averydevz@outlook.com
- **GitHub:** https://github.com/richavery/donk-cli-main
- **Issues:** https://github.com/richavery/donk-cli-main/issues

**See the project. Choose your intelligence. Get to work.**
