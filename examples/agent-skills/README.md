# Agent Skills examples (markdown-first)

These folders follow the [Agent Skills specification](https://agentskills.io/specification): each skill is a directory whose name matches the `name` in `SKILL.md` frontmatter.

Most illustrate **instruction-only** skills (no server-side tool code). **`read-web`** is different: it is normal Agent Skills markdown for discovery, but execution goes through the built-in MCP tool `read_web` (see `references/TOOL.md`).

## Using locally

The runtime server reads skills from `SKILLS_MCP_SERVER_DIR` (default `./skills`, gitignored). Copy or symlink an example into that directory:

```bash
mkdir -p skills
cp -R examples/agent-skills/meeting-prep-brief skills/
```

Restart `skills-mcp-server`, then `list_skills` / `read_skill` via REST or MCP.

## Included examples

| Directory | Purpose |
|-----------|---------|
| `meeting-prep-brief` | Short procedural skill with an optional reference file under `references/` |
| `read-web` | Same markdown as seeded into `SKILLS_MCP_SERVER_DIR/read-web` on startup (`SKILLS_SEED_BUILTIN_READ_WEB=true`): shows up in `list_skills` like CRUD skills; execution remains tool `read_web` |

Canonical seed copy for **`read-web`** also lives under `internal/skills/bundled/read-web` (embedded in the binary).

Community repos often publish similar trees (single `SKILL.md` or plus `references/`). This repo does not vendor third-party skills; use these as templates.
