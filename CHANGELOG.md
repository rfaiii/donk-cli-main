# Changelog

## 1.1.3

### Native Ollama coding

- Added `donk-cli code --native` for small local Ollama coding models.
- Added the `cmd/codetool` sidecar with a native Ollama `/api/chat` client and
  a minimal flat coding tool set.
- Added recovery for Qwen-style text tool calls, fenced JSON, truncated view
  calls, object-form arguments, and `file_path` aliases.
- Added bounded Ollama turns and explicit incomplete-stream errors so stalled
  local models do not leave the terminal loading indefinitely.
- Buffered likely tool-call JSON during streaming so protocol payloads do not
  clutter the user-facing response.
- Added installation, diagnostics, and local-model documentation for the native
  coder path.

### Validation

The native path was live-tested with `qwen2.5-coder:3b-instruct` against a local
Ollama daemon on macOS. Chat/file-read turns completed successfully in roughly
2–10 seconds during validation.