# CODING-TOOL.md

Design and iteration log for a **simple, `fantasy`-free coding path** for local/Ollama
models and the Cline/Hermes skills. Goal: make small/local models (8 GB MacBook-class
machines, 3B–14B qwen-coders) actually *code* without DONK's `charm.land/fantasy` agent
framework intercepting and reshaping every tool call.

## Why we need this

- DONK's built-in coder routes every conversation through `fantasy` and attaches the **full**
  tool surface — including the interactive `question` tool whose schema is nested
  (`questions[] → choices[] → confirm_title/description`).
- Small qwen models can't follow that 15+ tool / nested-schema contract. They emit the
  schema as chat text or invent tool names (`search_code`, `info`) instead of real
  `tool_calls`. This is invisible when the model is run bare in a terminal (no tool schema
  attached) — which is why outputs differ inside donk-cli.
- `fantasy`'s `openaicompat` provider *does* build correct OpenAI `tools`/`tool_calls`
  (`charm.land/fantasy/providers/openaicompat/language_model_hooks.go:410-538`), but for
  local models the contract is simply too rich to be reliable.

## Design principles

1. **Native Ollama, not OpenAI-compat.** Talk Ollama's own `/api/chat` endpoint (with
   `tools`) directly. No translation layer, no `fantasy`.
2. **Minimal, flat tool set.** Only what a coder needs: `bash`, `view`, `edit`,
   `glob`, `grep`, `ls`. Flat schemas (string/number args only) parse reliably on small
   models.
3. **Lenient tool-call recovery.** If a small model dumps a single valid tool call as text
   (no `tool_calls` array), parse it and execute it instead of echoing it to the user.
4. **Skills = direct exec.** Cline/Hermes are already `exec`'d via `bash`; keep that. No
   model needs to "invoke" a skill tool — just run the command.
5. **Native coder entry point.** Use `donk code --native` to bypass `fantasy` while
   keeping the normal DONK coder session available as the default.
6. **Sandbox-first.** Prototype in `~/Projects/donk-cli-test` (`this` repo), then
   copy the validated path into the main checkout.

## Planned tool set (flat schemas)

| tool   | args                                        | exec                       |
|--------|---------------------------------------------|----------------------------|
| bash   | `command` (string, required), `timeout` (int) | `sh -c <command>`          |
| view   | `path` (string, required)                   | read file bytes            |
| edit   | `path`, `old`, `new`, `count` (int, opt)    | exact string replace       |
| glob   | `pattern` (string, required)                | `filepath.Glob`            |
| grep   | `pattern`, `path`                          | recursive regexp           |
| ls     | `path`                                      | `ReadDir`                  |

## Request/response shape (Ollama native)

Request `POST /api/chat`:
```json
{ "model": "qwen2.5-coder:7b", "stream": false, "messages": [...],
  "tools": [ { "type": "function", "function": { "name": "bash",
    "description": "...", "parameters": { "type": "object",
    "properties": { "command": { "type": "string" }, "timeout": { "type": "integer" } },
    "required": ["command"] } } } ] }
```

Response: `message.content` (text) + optional `message.tool_calls[]` where each is
`{ name, arguments }`. Ollama may return `arguments` as either a JSON object or a
JSON string; `codetool` accepts both forms.

## Lenient recovery rule

If `tool_calls` is empty/absent but `content` parses (via `jsonrepair`-style lenient
parse) into a single object `{name, arguments}` whose `name` matches a known tool and
whose `arguments` parses to that tool's params → execute it; otherwise treat `content`
as final text.

## Skill delegation

`codetool --skill cline "fix failing tests"` → `exec.Command("cline", "--yolo", task)`,
streaming stdout straight to the terminal. Same for `hermes`. The local model is NOT
asked to call a skill tool.

## Why this helps 8GB / small-model users

- One small model, one tiny tool schema → far more likely to tool-call correctly.
- No `fantasy` streaming/combinatorics overhead.
- If a model still can't, the user drops to a bigger local model or delegates to Cline/Hermes
  with one flag — not a different app.

## Iteration log

- [x] Sandbox created at `~/Projects/donk-cli-test` (copy of donk-cli, fresh git repo).
- [x] `cmd/codetool/main.go` prototype: native Ollama `/api/chat`, flat 6-tool set,
      lenient recovery, skill exec. Compiles; `go vet` clean; unit tests green.
- [x] **Root cause found + fixed: bare-name tool calls.** qwen2.5-coder (3b + 7b) emit
      `{"name": glob, "arguments": {...}}` — the tool `name` is an **unquoted bare
      identifier** (and `content` is markdown-fenced), which breaks `json.Unmarshal`,
      so naive recovery drops the call and the model's text is echoed verbatim (the
      original qwen-dump symptom). Fix: `decodeToolCall` fast-paths strict JSON, then
      falls back to regex extraction of `"name"` + `"arguments"` (object or string form)
      + `normalizeToolName` aliasing.
- [x] **Root cause found + fixed: `recGlob` skipped the root directory.** `filepath.Walk(".")`
      hit `isExcludedDir(".")` → `strings.HasPrefix(".", ".") == true` → returned
      `SkipDir` on the root, aborting the walk (visited=0). The temp-dir unit test hid
      this because the temp root basename doesn't start with `.`. Fix: never skip the
      root itself (`path != root` guard) while still pruning `.git`/`.venv` subtrees.
- [x] **Live validated on `qwen2.5-coder:3b-instruct`** (Ollama `127.0.0.1:11434`):
  - TEST C "how many .go under cmd/": recovered bare `glob` → `recGlob` matched
    `cmd/codetool/{main.go,main_test.go,recovery_test.go}`; model answered.
  - TEST D "create todo.txt then read it back": recovered `write` → `[write] written`,
    then `view` → `[view] shopping: milk, eggs, bread` (file confirmed). End-to-end
    file I/O on a 3B model.
- [x] Committed in sandbox (`f9ed7b3`).
- [x] **Bare-regex-value recovery + lenient arg parsing.** Live 3b run dumped
      `{"name":"grep","arguments":{"path":"","pattern:"^.*$"}}` — the regex value
      is unquoted, so `matchBrace` swallowed the closing braces and the call was
      echoed to the UI (the original "scrambled feed" symptom). Fixed by (a)
      `extractToolCallsFromText` attempting `decodeToolCall` on the whole object
      when braces never balance, and (b) `parseArgsTolerant` (plus `bareColonRe`
      repair `"key:"^."` -> `"key":"^."`) so `execTool` can still run malformed
      args. `grep` now also defaults empty `path` to the cwd. Live re-validated:
      recovered `grep` executed and returned `main.go: func main() {`.
- [x] **Streaming (`-stream`).** `chatStream` reads Ollama's NDJSON stream, prints
      content deltas as they arrive and accumulates content + tool_calls for the
      agent loop. Streaming is enabled by default and slow/incomplete streams now
      fail explicitly instead of appearing to load forever.
- [x] **Tool-call UX buffering.** JSON or fenced-JSON tool-call dumps are buffered
      while streaming and suppressed once recovered, so users see `[view]`/`[bash]`
      status rather than raw protocol JSON. Ordinary prose still streams immediately.
- [x] **Native DONK integration.** `donk code --native` launches `codetool`, passes
      through `--cwd`, `--model`, and `--stream`, and finds the sidecar through
      `DONK_CODETOOL` or `PATH`.
- [x] **Live chat validation.** `qwen2.5-coder:3b-instruct` reads a file and answers
      correctly in roughly 2–10 seconds on a MacBook-class local Ollama setup.
- [ ] Benchmark multi-step coding tasks across small and large local models.
- [ ] Decide whether to merge the sidecar into the primary DONK binary later.

## Using the native coder

Build the sidecar and DONK binary:

```sh
cd ~/Projects/donk-cli-main
GOWORK=off go build -o "$HOME/bin/codetool" ./cmd/codetool
GOWORK=off go build -o "$HOME/bin/donk-cli" .
```

Run a local coding task:

```sh
DONK_CODETOOL="$HOME/bin/codetool" \
  "$HOME/bin/donk-cli" code --native \
  --model qwen2.5-coder:3b-instruct --stream \
  "Read the relevant files and fix the failing test"
```

`--native` uses the current working directory unless `--cwd` is provided. The
sidecar talks directly to `OLLAMA_HOST` (default `http://127.0.0.1:11434`). Set
`CODETOOL_DEBUG=1` for turn and tool diagnostics. Each model turn has a bounded
deadline and the overall budget is bounded by `-budget`, so a stalled daemon
returns an error instead of leaving the terminal in a loading state.
