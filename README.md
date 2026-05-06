# skills-mcp-server

Go-based MCP + REST server for Agent Skills spec content.

## Phase 1 scope

- Agent Skills spec validation (`SKILL.md` frontmatter + directory-name matching)
- Filesystem-backed skill registry
- REST CRUD for skills and skill resources
- MCP tools for list/read/search skills and resources
- Initial runtime target is `read-web` for Joanne and Ross via `agents-mcp-server`

## Local run

```bash
cp .env.dev.example .env.dev
set -a && source .env.dev && set +a
go run ./cmd/skills-mcp-server
```

## Examples (instruction-only skills)

Markdown-first Agent Skills live under [`examples/agent-skills/`](examples/agent-skills/). Copy a folder into your runtime `skills/` directory (see `examples/agent-skills/README.md`) so `list_skills` / `read_skill` can load them—same shape as community Claude-style procedural skills on GitHub.

### Built-in `read-web` skill (seeded)

When `SKILLS_SEED_BUILTIN_READ_WEB` is true (default), startup writes embedded content to `SKILLS_MCP_SERVER_DIR/read-web/` **only if** `read-web/SKILL.md` is not already present—so **`read-web` appears in `list_skills` like a CRUD-created skill**. Tool execution is still the Go-backed `read_web` MCP tool + REST alias; the markdown explains when to use it and how agents should call the tool.

Health:

```bash
curl http://localhost:8081/health
```

## REST

- `GET /api/skills?q=...`
- `POST /api/skills`
- `GET /api/skills/:name`
- `PUT /api/skills/:name`
- `DELETE /api/skills/:name`
- `GET /api/skills/:name/resources`
- `POST /api/skills/:name/resources`
- `GET /api/skills/:name/resources/:path`
- `PUT /api/skills/:name/resources/:path`
- `DELETE /api/skills/:name/resources/:path`
- `POST /api/runtime/read-web`

## MCP

Executable capabilities (`tools/list`, `tools/call`) are registered in Go via [`internal/tools`](internal/tools) (`Registry` + per-tool `Tool` implementations). Agent-facing markdown skills under `SKILLS_MCP_SERVER_DIR` are separate (discovery via `list_skills` / `read_skill`). Callers such as **`agents-mcp-server`** should keep using stable tool names (`read_web`, etc.).

`POST /mcp` with JSON-RPC 2.0:

- `initialize`
- `tools/list`
- `tools/call` for:
  - `list_skills`
  - `read_skill`
  - `search_skills`
  - `list_skill_resources`
  - `read_skill_resource`
  - `get_skill_resource_info`
  - `read_web`
