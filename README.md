# skills-mcp-server

Go-based MCP + REST server for Agent Skills spec content.

## Phase 1 scope

- Agent Skills spec validation (`SKILL.md` frontmatter + directory-name matching)
- Filesystem-backed skill registry
- REST CRUD for skills and skill resources
- MCP tools for list/read/search skills and resources, **plus create/update/delete skills** (`create_skill`, `update_skill`, `delete_skill`)
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

Executable capabilities (`tools/list`, `tools/call`) are registered in Go via [`internal/tools`](internal/tools) (`Registry` + per-tool `Tool` implementations). Agent-facing markdown skills under `SKILLS_MCP_SERVER_DIR` are separate (discovery via `list_skills` / `read_skill`). **Mutating skill CRUD** is also via MCP (`create_skill`, `update_skill`, `delete_skill`)—same backing store as REST. Callers such as **`agents-mcp-server`** should keep using stable tool names.

### MCP security note

The `/mcp` JSON-RPC endpoint is **not authenticated** in the default deployment: anything that can reach it can invoke tools, including **filesystem writes** from `create_skill` / `update_skill` / `delete_skill`. Run only on a trusted network (e.g. cluster-internal to `agents-mcp-server`). For production hardening, prefer network policy / mesh boundaries; optional future shared-secret or auth wrapper if you expose MCP beyond the mesh.

Skill content changes apply to the **runtime skills directory** (`SKILLS_MCP_SERVER_DIR`, often a PVC)—not automatic git commits to this repo.

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
  - `create_skill` — write new `<name>/SKILL.md` (requires name, description, instructions)
  - `update_skill` — replace existing skill markdown fields
  - `delete_skill` — remove skill directory (requires name; returns error if missing)
  - `read_web`
  - `create_google_doc` (OAuth: **`JOANNE_GOOGLE_*`** preferred in shared env; see `.env.dev.example`)

## Production (Kubernetes)

Runtime env for the cluster comes from Secret **`skills-mcp-server-runtime`** in namespace **`skills-mcp-server`** (see `rancher-admin` deployment). `rancher-admin/scripts/sync-app-pull-secrets.sh` only copies **`dockerhub-pull`**; it does **not** set `GEMINI_API_KEY` or Google OAuth.

**Skills on disk:** GitOps mounts **`SKILLS_MCP_SERVER_DIR`** (`/app/skills`) from a **ReadWriteOnce PVC** (`skills-mcp-server-data` in `rancher-admin/admin/apps/skills-mcp-server/deployment.yaml`) so skills created via MCP/REST survive pod restarts. If the PVC is **Pending**, set **`spec.storageClassName`** on that PVC to a provisioner your admin cluster exposes (cluster default SC is used when the field is omitted—tune per environment).

Push or refresh the runtime secret from a trusted machine (`.env.prod` in this repo, or e.g. `ENV_FILE=../agents-mcp-server/.env.prod` if you keep Joanne Google keys there):

```bash
./scripts/update-rancher-secrets.sh
# or
ENV_FILE=/path/to/.env.prod ./scripts/update-rancher-secrets.sh
```

`makeacompany-ai/scripts/update-rancher-secrets.sh` updates **`makeacompany-ai-runtime-secrets`** (portal `GOOGLE_OAUTH_*`); that does **not** populate this service—`create_google_doc` needs the keys on **`skills-mcp-server-runtime`**.
