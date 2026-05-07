# create_google_doc runtime wiring

## MCP

- **Method:** `tools/call`
- **Tool name:** `create_google_doc`
- **Arguments:** JSON object per skill frontmatter (`intent`, `title`, `editors`, optional `body`, `commenters`, `viewers`, …).

## REST (same process)

`POST /api/runtime/create-google-doc` with the same JSON body as MCP `arguments`.

## Environment (skills-mcp-server)

| Variable | Notes |
|----------|--------|
| `GOOGLE_CLIENT_ID` / `GOOGLE_OAUTH_CLIENT_ID` / `JOANNE_GOOGLE_CLIENT_ID` | OAuth client id (first non-empty wins in that order). |
| `GOOGLE_CLIENT_SECRET` / `GOOGLE_OAUTH_CLIENT_SECRET` / `JOANNE_GOOGLE_CLIENT_SECRET` | Client secret. |
| `GOOGLE_REFRESH_TOKEN` / `GOOGLE_OAUTH_REFRESH_TOKEN` / `JOANNE_GOOGLE_REFRESH_TOKEN` | **Required** for server-side API calls. Must be issued with scopes including **Google Docs** and **Drive file** access (`drive.file`). |

Portal-style **`GOOGLE_OAUTH_CLIENT_*`** names are supported so local compose can reuse the same keys as makeacompany-ai OAuth apps.

## Parity

Contract fields mirror **skill-factory** `tools/v1/create-google-doc.tool.json` input shape where applicable. Slack / employee-factory may still use **`joanne-create-doc`** as the historical runtime tool id; this MCP server exposes **`create_google_doc`** as the stable callable name.
