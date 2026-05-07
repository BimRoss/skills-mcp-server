# create_google_doc runtime wiring

## MCP

- **Method:** `tools/call`
- **Tool name:** `create_google_doc`
- **Arguments:** JSON object per skill frontmatter (`intent`, `title`, `editors`, optional `body`, `commenters`, `viewers`, …).

## REST (same process)

`POST /api/runtime/create-google-doc` with the same JSON body as MCP `arguments`.

Drive ACL grants use **invite notifications** (`sendNotificationEmail=true`) so addresses without an existing Google account can accept editor/commenter/viewer access (Drive rejects silent shares for those recipients).

## Environment (skills-mcp-server)

| Variable | Notes |
|----------|--------|
| **`JOANNE_GOOGLE_CLIENT_ID`** / `GOOGLE_CLIENT_ID` / `GOOGLE_OAUTH_CLIENT_ID` | OAuth client id (**Joanne first** so portal `GOOGLE_OAUTH_*` in the same file does not steal the client paired with Joanne’s refresh token). |
| **`JOANNE_GOOGLE_CLIENT_SECRET`** / … | Client secret (same precedence). |
| **`JOANNE_GOOGLE_REFRESH_TOKEN`** / … | **Required** for server-side API calls; must be issued for the **same** OAuth client as id/secret. Scopes: **Docs** + **`drive.file`**. |

Portal-style **`GOOGLE_OAUTH_CLIENT_*`** remains supported as **fallback** when Joanne keys are unset.

## Parity

Contract fields mirror **skill-factory** `tools/v1/create-google-doc.tool.json` input shape where applicable. Slack / employee-factory may still use **`joanne-create-doc`** as the historical runtime tool id; this MCP server exposes **`create_google_doc`** as the stable callable name.
