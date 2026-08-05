# DONK-CLI Coding Assistant

## Problem

Local Ollama models (qwen-coder, gemma4, etc.) often behave like chatbots instead of
autonomous coding agents. They emit chat instead of tool calls, skip the
"read before edit" discipline, and produce verbose explanations rather than
making changes.

## Goals

1. Make coding with local models as reliable as remote APIs.
2. Provide a dedicated `donk code` entry point that configures the right
   system prompt, tool set, and model-fit warnings.
3. Surface coder-agent state in the UI so the user always knows whether they
   are in "chat" or "code" mode.

## Proposed items

### A. Dedicated coder entry point

- Add a `donk code` subcommand that boots a coder-only agent.
- Add a `/code` toggle inside the TUI to switch the active session into coder mode.
- When coder mode activates:
  - Force the coder agent model (`agents.coder.model`).
  - Restrict tools to coding set: `edit`, `write`, `view`, `bash`, `grep`,
    `glob`, `ls`, `lsp_*`, `node`, `npm`, `donk_info`.
  - Enable beast mode automatically or prompt once.

### B. Model-fit tagging and warnings

- Tag Ollama models as `coding` / `chat` by name heuristic in
  `internal/localmodel`.
- In the model picker and `/models` dialog, show a small tag next to each
  local model.
- When a chat-tagged model is assigned to `agents.coder.model`, show a
  warning banner with a one-click suggestion to switch to a coding-capable
  model.

### C. Coder prompt improvements for small models

- Add a stronger "You are a coding agent" preamble to `coder.md.tpl`.
- Reduce multi-tool parallelism instructions for small models.
- Add a fallback: when the model emits plain text instead of a tool call,
  retry once with an explicit "respond with a tool call" nudge.
- Keep the existing strict rules (read before edit, never commit unless
  asked, concise output).

### D. Coder agent status indicator

- Add a small pill in the header showing:
  - `coder` / `chat` mode
  - active coder model
  - files modified this turn
  - beast mode state
- Use the existing header diagonal-details area so it does not grow the UI.

### E. Coding-assistant skills (Cline & Hermes)

- Embed Cline and Hermes as builtin `agentskills.io`-standard skills so the
  current local coder model can delegate to a stronger autonomous agent when it
  is a small model struggling with multi-step coding.
- Skills are surfaced in the coder system prompt as `<available_skills>` (built
  in `internal/agent/prompt/prompt.go`, rendered by
  `internal/agent/templates/coder.md.tpl`). Each entry carries a `<description>`
  trigger and a `donk://skills/<name>/SKILL.md` `<location>`.
- The "LOAD MATCHING SKILLS" rule in the prompt (`coder.md.tpl`) directs the
  agent to `view` the skill's `<location>` before acting. The View tool resolves
  `donk://` URIs from the embedded filesystem (see `internal/agent/tools/view.go`).
- Cline is `user-invocable` (npm global install) and Hermes is `user-invocable`
  (curl install script); both can be invoked directly by the user from chat.
- Both are `type=builtin` -> active at app load, sorted ahead of other coding
  skills in `<available_skills>`.

## Task list

- [x] Audit local model behavior: log tool-call vs chat emissions when a
      local model is used as the coder agent.
- [x] Build `donk code` subcommand with coder-only tool set.
- [x] Add `/code` toggle in TUI to switch active session into coder mode.
- [x] Add model-fit tagging in `internal/localmodel` for Ollama models.
- [x] Show model-fit tags in `/models` dialog and model picker.
- [x] Add warning when a chat-tagged model is assigned to coder agent.
- [x] Improve `internal/agent/templates/coder.md.tpl` for small local models.
- [ ] Add fallback retry when coder model emits text instead of tool call.
- [x] Add coder agent status pill in header.
- [x] Wire Cline & Hermes as coding-assistant skills surfaced via `<available_skills>`.
- [x] Disable Cline/Hermes in VCR `TestCoderAgent` (keep recorded skill set stable).
- [ ] Add tests for new coder mode, model tags, and header indicator.

## Next steps

1. Start with model-fit tagging and warnings (B) so we can classify local
   models before wiring the coder command.
2. Improve `coder.md.tpl` (C) so new coder sessions have the best prompt.
3. Build `donk code` and `/code` toggle (A).
4. Add header status pill (D).
5. Add audit logging (A) and tests.

## Current status

This session wired **Cline** and **Hermes** as builtin coding-assistant skills
(see Proposal E and the task list above). They appear in the coder prompt as
`<available_skills>` and can be invoked by the agent via their
`donk://skills/<name>/SKILL.md` locations, or directly by the user from chat.
Cline sorts first.

VCR `TestCoderAgent` cassettes record only a stable skill set (`donk-config` +
`donk-hooks` + `jq`), with Cline/Hermes disabled in
`internal/agent/common_test.go` (`DisabledSkills`) so the recorded prompt is
byte-level stable — the generated system prompt is not byte-identical to the
cassettes only because of the pre-existing `coder.md.tpl` first-line change
(`Assistant` -> `coding agent`), which predates this work. Re-recording
cassettes with a Hyper API key is the only way to clear that pre-existing
failure; this change introduces no additional cassette diff.

The remaining open item is the text-vs-tool-call fallback retry (Task list).

## Startup and Cline menu (current)

- **Local Ollama models at startup:** `Models.LocalModelsCmd()`
  (`internal/ui/dialog/models.go`) now starts an installed-but-offline Ollama
  daemon on open/refresh, so local models list on first launch without a manual
  `ollama serve`. Onboarding help exposes `r` (refresh) and `s` (start Ollama).
- **Skip model selection:** pressing `esc` on the startup model picker emits
  `ActionSkipOnboarding` and lands on the homescreen (`uiLanding`) with
  `SetupAgents()` + `InitCoderAgent()` called, so users are never forced to
  choose a provider/sign up to get inside the app.
- **Cline / Hermes are always reachable:** built-in and `user-invocable`, they
  load into the command palette at startup regardless of model state. The
  palette renders them title-cased with `/cline` and `/hermes` aliases
  (`internal/ui/dialog/commands.go`). Open the palette (`ctrl+p`), type `/cline`,
  select it to attach the skill, then send a task to delegate to Cline.
- Deeper Cline UX (per-task plan/act mode selector, direct "run Cline on this
  task" without a local model, and settings-menu conflict avoidance) lives in
  [`CLINE_FEATURE.md`](CLINE_FEATURE.md).
