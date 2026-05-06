# read_web runtime wiring

This skill’s instructions point agents at the **`read_web`** tool. Implementation lives in server code (`internal/readweb`), not under `scripts/`.

## MCP

- Method: `tools/call`
- Tool name: `read_web`
- Arguments: `{ "query": "string" }`

## REST (same backend)

`POST /api/runtime/read-web` with JSON `{ "query": "..." }` returns `fallbackText`, `finalSummary`, and `citations`.

## Environment (server)

| Variable | Purpose |
|----------|---------|
| `GEMINI_API_KEY` | Required for live calls |
| `GEMINI_MODEL` | Model id (default `gemini-2.5-flash` in config) |
| `ENABLE_WEB_RESEARCH` | When true, attaches Google Search grounding |

## Seeding this folder

On startup, if `SKILLS_SEED_BUILTIN_READ_WEB` is true and `./read-web/SKILL.md` is missing under `SKILLS_MCP_SERVER_DIR`, the server writes this bundled tree so **`read-web` appears like any other skill** in `list_skills`.
