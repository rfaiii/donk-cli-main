# DONK Skills

DONK discovers Agent Skills from the default `~/.agents/skills` directory, as
well as project skill directories such as `.agents/skills` and `.donk/skills`.
Each skill must be a directory containing a valid `SKILL.md` with frontmatter.

The local master catalog at `~/Documents/AI-SKILLS/MASTER-AI-SKILLS.md` is an
inventory of skill directories, not a single skill file. To make the catalog
available to DONK without copying or overwriting existing skills, run:

```sh
./scripts/install-master-skills.sh
```

The installer creates symlinks from `~/.agents/skills/<name>` to each source
skill directory. Existing entries are preserved. Because the links point at the
source directories, edits and updates in `~/Documents/AI-SKILLS` are available
to DONK automatically after the next skill discovery.

Override paths when needed:

```sh
DONK_MASTER_SKILLS_DIR=/path/to/AI-SKILLS \
DONK_SKILLS_DIR=/path/to/skills \
./scripts/install-master-skills.sh

## Builtin skills

DONK also embeds a set of builtin skills directly into the binary via
`internal/skills/builtin/`. These are always available — no installation or
symlinks needed — and are surfaced to the coding agent as `<available_skills>`
in the coder system prompt. Builtin skills use virtual
`donk://skills/<name>/SKILL.md` locations, which the View tool resolves from the
embedded filesystem.

| Skill            | Description                                                  |
| ---------------- | ------------------------------------------------------------ |
| `donk-config`    | DONK configuration help                                      |
| `donk-hooks`     | Authoring, configuring, and debugging hooks                  |
| `jq`             | jq JSON processor usage guide                                |
| `cline`          | Autonomous terminal coding agent (delegates a coding job)    |
| `hermes`         | Self-improving autonomous agent runtime and model router     |

`cline` and `hermes` are coding-assistant skills: when the local coder model is a
small model that struggles with multi-step coding, the coder system prompt
(`internal/agent/templates/coder.md.tpl`) directs the agent to hand off to Cline
or Hermes. Both follow the
[agentskills.io](https://agentskills.io) open standard — the same format DONK
skills use — so project rules and skills can be shared. See
[`DONK_CODER.md`](DONK_CODER.md) for the coder-entry design.

User skills with the same name as a builtin override the builtin (last
occurrence wins in `skills.Deduplicate()`).