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

## Skills directory

Runtime skills live under **`SKILLS_MCP_SERVER_DIR`** (default `./skills`). You can add them manually (REST CRUD, `cp`, CI) or rely on **startup seeding** from the repo’s [`examples/`](examples/) tree.

### Seeding from `examples/` (default on)

On startup, each **direct subdirectory** of **`SKILLS_EXAMPLES_DIR`** that contains a top-level `SKILL.md` is copied into `SKILLS_MCP_SERVER_DIR/<name>/` **only if** that skill does not already exist there (first write wins; idempotent across restarts).

| Env | Default | Meaning |
|-----|---------|--------|
| `SKILLS_SEED_EXAMPLES` | `true` | Set `false` to disable seeding entirely |
| `SKILLS_EXAMPLES_DIR` | `examples` (relative to process cwd) | Source tree; Docker image sets `/app/examples` |

The Docker image **`COPY`s `examples/`** into `/app/examples` and sets `SKILLS_EXAMPLES_DIR` so Compose volumes get populated on first boot.

For **`read-web`**, markdown is for discovery; execution is still the Go **`read_web`** MCP tool + REST alias. For **`create-google-doc`**, discovery is markdown under `examples/create-google-doc/`; execution is **`create_google_doc`** + `POST /api/runtime/create-google-doc`.

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
- `POST /api/runtime/create-google-doc`

## MCP

Executable capabilities (`tools/list`, `tools/call`) are registered in Go via [`internal/tools`](internal/tools) (`Registry` + per-tool `Tool` implementations). Agent-facing markdown skills under `SKILLS_MCP_SERVER_DIR` are separate (discovery via `list_skills` / `read_skill`). Callers such as **`agents-mcp-server`** should keep using stable tool names (`read_web`, `create_google_doc`, etc.).

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
  - `create_google_doc` (Google OAuth refresh token + Docs/Drive scopes; see `.env.dev.example`)
