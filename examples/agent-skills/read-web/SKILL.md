---
name: read-web
description: Research the public web and return a concise summary with optional grounding citations. Use when the user asks for current events, look-ups, "what happened with…", competitors, docs, or any fact that should be verified against live sources—not for private repo or Slack-only context.
license: MIT
metadata:
  author: bimross-skills-mcp-server
  version: "1.0"
  execution: mcp-tool-read_web
compatibility: Requires the skills-mcp-server runtime tool read_web (Gemini + optional Google Search grounding). Set GEMINI_API_KEY on the server.
---

# Read web (research)

## Role

This skill is the **Agent Skills** packaging for the built-in **`read_web`** tool. Same discovery rules as any CRUD skill: `list_skills` shows `name` + `description`; `read_skill` loads full instructions.

## When to activate

- Questions about **recent or external** facts where training data may be stale.
- **Explicit research** asks ("look up", "find out", "latest on").
- Short factual checks that benefit from search-backed grounding.

Avoid activating for purely internal coordination ("what did we decide in this thread?") unless the user also asks for external verification.

## Execution (not a local script)

Running this skill means invoking the MCP tool:

- **Tool id:** `read_web`
- **Arguments:** `{ "query": "<user task distilled to a single search/research query>" }`

The server performs Gemini `generateContent` with optional `google_search` grounding; citations may be returned alongside the summary.

## Output expectations

Prefer a **short summary** in the user-visible reply; keep raw citation URLs in structured/tool results when the client supports it.

## More detail

See [references/TOOL.md](references/TOOL.md) for HTTP/MCP surfaces and env vars on this server.
