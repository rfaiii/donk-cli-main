# DONK CLI

<p align="center">
    <a href="https://github.com/rfaiii/donk-cli-main/releases"><img src="https://img.shields.io/github/release/rfaiii/donk-cli-main" alt="Latest Release"></a>
    <a href="https://github.com/rfaiii/donk-cli-main/actions/workflows/build.yml/badge.svg" alt="Build Status"></a>
</p>

<p align="center">Your new coding bestie, now available in your favourite terminal.<br />Your tools, your code, and your workflows, wired into your LLM of choice.</p>
<p align="center">终端里的编程新搭档，<br />无缝接入你的工具、代码与工作流，全面兼容主流 LLM 模型。</p>

## Features

- **Multi-Model:** choose from a wide range of LLMs or add your own via OpenAI- or Anthropic-compatible APIs
- **Flexible:** switch LLMs mid-session while preserving context
- **Session-Based:** maintain multiple work sessions and contexts per project
- **LSP-Enhanced:** DONK CLI uses LSPs for additional context, just like you do
- **Extensible:** add capabilities via MCPs (`http`, `stdio`, and `sse`)
- **Node Integration:** built-in NODE device discovery, status tracking, and npm script execution from the command palette
- **File Finder:** improved two-pane file browser with fixed pane widths, scrollbar track, and control-safe rendering
- **Command Palette:** quick access to system commands, npm scripts, and custom actions
- **Works Everywhere:** first-class support in every terminal on macOS, Linux, Windows (PowerShell and WSL), Android, FreeBSD, OpenBSD, and NetBSD

## Installation

Use a package manager:

```bash
# Go install
go install github.com/rfaiii/donk-cli-main@latest

# Build from source
git clone https://github.com/rfaiii/donk-cli-main.git
cd donk-cli-main
go build ./...
```

On illumos (OpenIndiana, OmniOS), the command above works as-is. Only native OS notifications are unavailable there; terminal-based notifications (OSC) and the terminal bell still work. On Oracle Solaris, add `-tags sqlite3_dotlk` so the local database uses dot-file locking:

```bash
go install -tags sqlite3_dotlk github.com/rfaiii/donk-cli-main@latest
```

## Getting Started

The quickest way to get started is to choose a model from the model picker. Follow the steps to authenticate and you'll be good to go.

## API Keys

You can use DONK CLI with many providers such as Anthropic, OpenAI, Gemini, OpenRouter, and more. Press `ctrl+l` to open the model picker, choose the provider of your choice, and paste your API key.

That said, you can also set environment variables for preferred providers:

| Environment Variable        | Provider                                           |
| --------------------------- | -------------------------------------------------- |
| `HYPER_API_KEY`             | Charm Hyper                                        |
| `ANTHROPIC_API_KEY`         | Anthropic                                          |
| `OPENAI_API_KEY`            | OpenAI                                             |
| `VERCEL_API_KEY`            | Vercel AI Gateway                                  |
| `GEMINI_API_KEY`            | Google Gemini                                      |
| `ZAI_API_KEY`               | Z.ai                                               |
| `MINIMAX_API_KEY`           | MiniMax                                            |
| `SYNTHETIC_API_KEY`         | Synthetic                                          |
| `HF_TOKEN`                  | Hugging Face Inference                             |
| `CEREBRAS_API_KEY`          | Cerebras                                           |
| `OPENROUTER_API_KEY`        | OpenRouter                                         |
| `IONET_API_KEY`             | io.net                                             |
| `ALIBABA_SINGAPORE_API_KEY` | Alibaba (Singapore)                                |
| `ALIBABA_US_API_KEY`        | Alibaba (United States)                            |
| `GROQ_API_KEY`              | Groq                                               |
| `AVIAN_API_KEY`             | Avian                                              |
| `OPENCODE_API_KEY`          | OpenCode Zen & Go                                  |
| `VERTEXAI_PROJECT`          | Google Cloud VertexAI (Gemini)                     |
| `VERTEXAI_LOCATION`         | Google Cloud VertexAI (Gemini)                     |
| `AWS_ACCESS_KEY_ID`         | Amazon Bedrock (Claude)                            |
| `AWS_SECRET_ACCESS_KEY`     | Amazon Bedrock (Claude)                            |
| `AWS_REGION`                | Amazon Bedrock (Claude)                            |
| `AWS_PROFILE`               | Amazon Bedrock (Custom Profile)                    |
| `AWS_BEARER_TOKEN_BEDROCK`  | Amazon Bedrock                                     |
| `AZURE_OPENAI_API_ENDPOINT` | Azure OpenAI models                                |
| `AZURE_OPENAI_API_KEY`      | Azure OpenAI models (optional when using Entra ID) |
| `AZURE_OPENAI_API_VERSION`  | Azure OpenAI models                                |
| `MOONSHOT_API_KEY`          | Moonshot                                           |

Also note that DONK CLI can support nearly any provider, including [Local Models](#local-models). For more info see [Custom Providers](#custom-providers) below.

### Local Models

DONK CLI works with local model runners such as Ollama, llama.cpp, LM Studio, and MLX. Use the model picker or add them manually in your config.

## Configuration

> [!TIP]
> DONK CLI ships with a builtin skill for configuring itself. Most of the time you can just tell what you want it to configure and it will get the job done.

DONK CLI runs great with no configuration. That said, if you do need or want to customize DONK CLI, you can, with a `donkrc`.

A `donkrc` is just Bash with some DONK CLI-specific builtins. It’s a lot like a `.bashrc`, just for your DONK CLI. Because DONK CLI has a native, built-in Bash interpreter, Bash-based config works identically across all platforms, including Windows.

For example:

```bash
# Add Ollama.
provider add ollama --type ollama --base-url "http://localhost:11434/v1"

# Register a model on Ollama.
model add ollama/llama3.3 --name "Llama 3.3" --context-window 128000

# Auto-approve some tools.
permissions allow view edit

# Add an MCP server.
mcp add github --type http --url "https://api.githubcopilot.com/mcp/" --header Authorization "Bearer $GH_PAT"
```

Configuration can be added either local to the project itself, or globally, with the following priority:

| Priority | Unix-like                 | Windows                               |
| -------- | ------------------------- | ------------------------------------- |
| 1        | `./.donkrc`               | `.\\.donkrc`                          |
| 2        | `./donkrc`                | `.\\donkrc`                           |
| 3        | `~/.config/donk/donkrc`   | `%USERPROFILE%\\.config\\donk\\donkrc` |

(DONK CLI respects the [XDG Base Directory Specification][xdg], so your paths may differ depending on your `XDG_CONFIG_HOME` value. Data directories such as `~/.local/share/donk` and `%LOCALAPPDATA%\\donk` contain JSON state only; DONK CLI does not execute a `donkrc` from them.)

[xdg]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

What about the old JSON format? It’s still supported, but it should be considered deprecated. See: [the config docs](./docs/config/) for details.

> [!TIP]
> You can override the user and data config locations by setting:
>
> - `DONK_GLOBAL_CONFIG`
> - `DONK_GLOBAL_DATA`

As an additional note, DONK CLI also stores ephemeral data, such as application state, in one additional location. This is state and should not be edited by hand, nor should it be considered configuration.

```bash
# Unix
$HOME/.local/share/donk/donk.json

# Windows
%LOCALAPPDATA%\\donk\\donk.json
```

#### A note on security

Both `donkrc` and `donk.json` are trusted code; `donkrc` runs in a full shell, and any `$(...)` in `donk.json` runs at load time. Don't launch DONK CLI in a directory whose config you haven't reviewed, and don't randomly `source` files from the internet into your config.

### LSPs

DONK CLI can use LSPs for additional context to help inform its decisions, just like you would. LSPs can be added manually like so:

```bash
# donkrc

lsp add go --command "gopls" --env "GOTOOLCHAIN go1.24.5"
lsp add typescript --command "typescript-language-server" --args --stdio
lsp add nix --command "nil"
```

### MCPs

DONK CLI also supports Model Context Protocol (MCP) servers through three transport types: `stdio` for command-line servers, `http` for HTTP endpoints, and `sse` for Server-Sent Events.

```bash
# donkrc

# Add a local MCP server that runs a Node.js script.
mcp add filesystem --command node --args /path/to/mcp-server.js \
  --timeout 120 --disabled-tools some-tool-name --env NODE_ENV production

# Add a GitHub MCP server that uses an API token.
mcp add github --type http --url "https://api.githubcopilot.com/mcp/" \
  --timeout 120 --header Authorization "Bearer $GH_PAT" \
  --disabled-tools create_issue --disabled-tools create_pull_request

# Add a streaming MCP server that uses SSE.
mcp add streaming-service --type sse --url "https://example.com/mcp/sse" \
  --timeout 120 --header API-Key "$API_KEY"
```

#### MCP OAuth

HTTP and SSE MCP servers that require OAuth can use DONK CLI's built-in authorization-code flow instead of a static `Authorization` header. Set `"oauth": true` to enable it:

```json
{
  "mcp": {
    "linear": {
      "type": "http",
      "url": "https://mcp.linear.app/mcp",
      "oauth": true
    }
  }
}
```

##### Pre-registered clients

Some servers (GitHub, Slack) don't support dynamic client registration. For those, register an OAuth app with the provider and supply the credentials directly. All values support shell expansion:

```json
{
  "mcp": {
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "oauth": true,
      "oauth_client_id": "Iv1.abc123def456",
      "oauth_client_secret": "$GITHUB_MCP_SECRET",
      "oauth_callback_port": 40704
    }
  }
}
```

When `oauth_client_id` is set, DONK CLI skips dynamic client registration and authenticates as the specified client. When omitted, DONK CLI attempts dynamic registration automatically (works with Linear, Notion, and other servers that support RFC 7591).

### Hooks

DONK CLI has preliminary support for hooks. For details, see [the hook guide](./docs/hooks/).

### Sharing a workspace across clients

When DONK CLI is run against a shared backend, clients are grouped into **workspaces** keyed by their resolved `--cwd`. Two clients with the same `--cwd` join the same underlying workspace, so they share the session list, message history, permission queue, LSP, and MCP state.

Joining is implicit: pointing a second client at the same working directory attaches it to the existing workspace. Each new invocation, however, starts in its own fresh session by default. To pick up the conversation another client already has open, use the session manager (the session picker) and select it. Sessions surface two signals there:

- `IsBusy` is set while an agent turn is in flight for that session.
- `AttachedClients` reports how many clients are currently viewing it.

A non-zero `AttachedClients` (often combined with `IsBusy`) is the cue that a session is "in progress" on another client and joining it will mirror that view live.

The first client to create a workspace fixes its process-wide flags. In particular, `--beastmode` and `--debug` follow a **first-wins** rule: later clients that arrive at the same `--cwd` with different values for those flags do not change the running workspace. A debug log line is emitted recording the mismatch, and the workspace keeps the flags it was created with.

A workspace lives as long as at least one client has an SSE event stream open against it. When the last stream disconnects, the workspace is torn down. There is a short grace window right after `POST /v1/workspaces` so a client that has created the workspace but not yet opened its event stream does not get reaped before it can attach.

### Global context files

DONK CLI automatically includes two files for cross-project instructions. Think of these are personal additions to the system prompt.

- `~/.config/donk/AGENTS.md`: DONK CLI-specific rules that would confuse other agentic coding tools. If you only use DONK CLI, this is the only one you need to edit.
- `~/.config/AGENTS.md`: generic instructions that other coding tools might read. Avoid referring to DONK CLI-specific features or workflows here. You probably only care about this if you use multiple agentic coding tools and want to share instructions between them.

You can customize these paths with `option global-context-path`. Repeat the command to add multiple paths:

```bash
# Load a single markdown file.
option global-context-path "~/path/to/custom/context/file.md"

# Recursively load all Markdown files in the folder.
option global-context-path "/full/path/to/folder/of/files/"
```

### Ignoring Files

DONK CLI respects `.gitignore` files by default, but you can also create a `.donkignore` file to specify additional files and directories that DONK CLI should ignore. This is useful for excluding files that you want in version control but don't want DONK CLI to consider when providing context.

The `.donkignore` file uses the same syntax as `.gitignore` and can be placed in the root of your project or in subdirectories.

### Allowing Tools

By default, DONK CLI will ask you for permission before running tool calls. If you'd like, you can allow tools to be executed without prompting you for permissions. Use this with care.

```bash
permissions allow view ls grep edit mcp_context7_get-library-doc
```

### Disabling Built-In Tools

You can also deny tools, hiding them from the agent entirely:

```bash
permissions deny bash sourcegraph
```

To disable tools from MCP servers, see the [MCP config section](#mcps).

## Contributing

Contributions are welcome. Please open an issue or pull request on GitHub.

____________________________________________________________________

/\./\./\./\./\./\./\./\/\./\./\./\./\./\./\./\/\./\./\./\./\./\./\./\
 — — — DONK-CLI — — —  — — — DONK-CLI — — —  — — — DONK-CLI — — — 
/\./\./\./\./\./\./\./\/\./\./\./\./\./\./\./\/\./\./\./\./\./\./\./\
____________________________________________________________________
 
 AUTHOR: RICHARD AIZEN AVERY III 
 CONTRIBUTOR: CRAZY JEFF(REY) LIVINGSTON
 FUNDED BY: STEVE THE BEAV JEPPSON 
 DEDICATED TO: MY DOOPZ JENNIFER TRAISTER & 
 FELIX THE TORNADO AVERY + THE BEAV!

“DONK” is a production created by Richard Aizen Avery III. 
© 2026 by Averydevz Fullstack Development and All Rights Reserved. 
Detroit, Michigan. https://donkai.tech/contact for questions.
