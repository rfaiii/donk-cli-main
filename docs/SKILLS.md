# DONK Skills

DONK discovers Agent Skills from the default `~/.agents/skills` directory, as
well as project skill directories such as `.agents/skills` and `.crush/skills`.
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
```