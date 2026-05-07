---
name: create-google-doc
description: Create a Google Doc with plain text, set a title, and share it with collaborators by email (editors, optional commenters/viewers). Use when the user asks for a shared Google document, a write-up in Docs, or meeting notes as a Doc—not for Markdown-only or Slack-native artifacts alone.
license: MIT
metadata:
  author: bimross-skills-mcp-server
  version: "1.0"
  execution: mcp-tool-create_google_doc
  skillFactoryParity: skill.contract.v1 / create-google-doc
compatibility: Requires skills-mcp-server OAuth env (Google Docs + Drive `drive.file` delegated via refresh token). Client id/secret may use `GOOGLE_CLIENT_ID` / `GOOGLE_OAUTH_CLIENT_ID` style names; a **refresh token** with Docs scopes is mandatory for API calls.
---

# Create Google Doc

## Role

Packaged **Agent Skill** for the **`create_google_doc`** MCP tool. Aligns with **skill-factory** `create-google-doc` / `joanne-create-doc` intent: structured doc creation with sharing.

## When to activate

- User wants a **Google Doc** (not Notion, not a Slack canvas-only summary unless they say Doc).
- They name emails or say “share with …” / editors / viewers.
- Thread content should become **document body** (pass as `body` or as rich `intent` when `body` is omitted).

**Mutating / side effects:** creates cloud content and Drive permissions—treat as **confirm-before-run** in product surfaces that use human confirmation UX.

## Execution

Invoke MCP tool:

- **Tool id:** `create_google_doc`
- **Required JSON arguments:** `intent` (string), `title` (string), `editors` (array of email strings)
- **Optional:** `body` (plain text; if empty, `intent` is used as the doc body), `commenters`, `viewers`, `type`, `length` (hints only; `type`/`length` are ignored by this server build)

`editors` / `commenters` / `viewers` must be **valid email addresses** for Google accounts.

## Output

Tool structured content includes **`fallbackText`**, **`finalSummary`**, **`docId`**, **`docUrl`** (edit link). Prefer surfacing the URL to the user.

## Wiring

See [references/TOOL.md](references/TOOL.md) for env vars and the REST alias on this server.
