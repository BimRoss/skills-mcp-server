# Example Agent Skills

Each subdirectory is a standalone skill tree (`SKILL.md` + optional `references/`). By default the server **seeds** these into `SKILLS_MCP_SERVER_DIR` on startup when missing (see root `README.md`); you can still add or override skills manually.

| Directory | Notes |
|-----------|--------|
| `read-web` | Markdown for discovery; runtime calls MCP tool `read_web` — see `references/TOOL.md` |
| `read-company-directory` | Instruction-only squad / company directory |
