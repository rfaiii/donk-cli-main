# Local Models — Phase 1

Phase 1 adds built-in Ollama runtime discovery without requiring a manually
maintained model JSON file.

The `internal/localmodel` package detects the Ollama executable and daemon,
lists installed models through `/api/tags`, and enriches model metadata through
`/api/show`. It also converts discovered models into the existing Catwalk model
shape so later model-picker integration can reuse the provider architecture.

The default endpoint is `http://127.0.0.1:11434`; `OLLAMA_HOST` can override it.
Phase 3 adds `r` refresh, `p` pull, and `s` start-Ollama actions to the model
picker. Pull uses Ollama's streaming `/api/pull` endpoint and renders status,
percentage, and digest in the picker. Press `c` to cancel an active pull. A
completed pull refreshes the local model list.

Phase 4 centralizes the managed `ollama-local` provider configuration. Selecting
a discovered model writes the local endpoint and model metadata through the
existing workspace configuration API, then refreshes the agent model. Manual
provider JSON remains supported for custom endpoints and advanced options.

Phase 5 adds discovery for LM Studio (`127.0.0.1:1234/v1`) and llama.cpp
(`127.0.0.1:8080/v1`) through their OpenAI-compatible `/v1/models` endpoints.
These runtimes appear as additional local model groups when their servers are
available; they do not replace or overwrite Ollama configuration.

## Native coding path

The managed `ollama-local` provider remains the default integration for normal
DONK conversations. Small coding models can opt into the native sidecar with
`donk code --native`. This path uses Ollama `/api/chat` directly, avoids the
large `fantasy` tool schema, and supports lenient recovery for Qwen-style text
tool calls. See [`OLLAMA_HOW_TO.md`](OLLAMA_HOW_TO.md) and
[`../CODING-TOOL.md`](../CODING-TOOL.md) for installation and diagnostics.