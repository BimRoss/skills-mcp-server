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

## Installing this skill

With **`SKILLS_SEED_EXAMPLES=true`** (default), the server copies each skill from **`SKILLS_EXAMPLES_DIR`** (defaults to `./examples`, Docker `/app/examples`) into **`SKILLS_MCP_SERVER_DIR`** when that skill is not already present.

You can also copy manually: `cp -R examples/read-web ./skills/` (adjust paths).
