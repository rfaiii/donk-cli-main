# DONK installation and testing

DONK is a terminal-first AI workspace. This guide covers macOS Apple Silicon,
macOS Intel, Windows x64, and Linux x64/arm64.

## Downloaded release archives

1. Download the archive matching your operating system and CPU.
2. Extract it.
3. Put `donk-cli` somewhere on your `PATH`.
4. Open a new terminal and verify:

```sh
donk-cli --version
donk-cli --help
```

### macOS

Use `darwin_arm64` for Apple Silicon and `darwin_x86_64` for Intel Macs:

```sh
chmod +x donk-cli
sudo install -m 0755 donk-cli /usr/local/bin/donk-cli
```

If macOS blocks an unsigned beta binary, approve DONK in **System Settings →
Privacy & Security** and run it again.

### Linux

Use `linux_x86_64` for Intel/AMD 64-bit or `linux_arm64` for ARM64:

```sh
chmod +x donk-cli
sudo install -m 0755 donk-cli /usr/local/bin/donk-cli
```

Without administrator access:

```sh
mkdir -p ~/.local/bin
install -m 0755 donk-cli ~/.local/bin/donk-cli
```

Ensure `~/.local/bin` is on `PATH`.

### Windows x64

Use the `windows_x86_64` ZIP archive. Extract `donk-cli.exe` to a folder such
as `%LOCALAPPDATA%\DONK\bin`, add that folder to the user `PATH`, and open a
new PowerShell window:

```powershell
donk-cli.exe --version
donk-cli.exe --help
```

## Build from source

Requirements: Go 1.26 or newer.

```sh
cd /path/to/donk-cli-go
go build -o ./donk-cli .
./donk-cli
```

To install the current source build:

```sh
./scripts/install-donk-cli.sh
hash -r 2>/dev/null || true
donk-cli --version
```

Always use `./donk-cli` when testing a repository build. A bare `donk-cli` may
resolve to an older copy elsewhere on `PATH`.

## Choose a project

Start DONK in a project:

```sh
donk-cli --cwd /path/to/project
```

Inside the app, use `/cd` or `/project`, then press `s` on a directory to switch.
You can also type `cd ~/Projects/my-project`; DONK restarts there without
sending the command to the AI agent.

## Ollama setup

Install Ollama from [ollama.com](https://ollama.com). DONK discovers every model
reported by Ollama’s local API in `/models` and `/ollama`. Selecting a model
automatically starts and warms it when possible.

```sh
ollama --version
ollama ls
```

The default endpoint is `http://127.0.0.1:11434`. Set `OLLAMA_HOST` for another
endpoint. See [`OLLAMA_HOW_TO.md`](OLLAMA_HOW_TO.md).

## Agent Skills

DONK discovers skills from `~/.agents/skills`. Sync the master catalog with:

```sh
./scripts/install-master-skills.sh
```

See [`SKILLS.md`](SKILLS.md).

## Configuration and data

Use `donk-cli dirs` to inspect platform-specific data locations. `donk-cli
projects` lists recently tracked projects.

## Troubleshooting

- **The app looks unchanged:** rebuild and run `./donk-cli`, not an older copy
  found on `PATH`.
- **Ollama models are missing:** start Ollama, open `/models`, and press `r`.
- **The model is slow initially:** the first request loads model weights; DONK
  warms the selected model after selection.
- **macOS security warning:** approve the binary in Privacy & Security.
- **Windows command not found:** reopen PowerShell after changing `PATH`.

## Beta verification

```sh
donk-cli --version
donk-cli --help
donk-cli dirs
donk-cli projects
```

Then launch DONK, switch to a test project, open `/models`, refresh local
models, select one, and confirm the model indicator becomes ready.