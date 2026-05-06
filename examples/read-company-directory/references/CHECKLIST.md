# Company directory — operator checklist

Use this when tightening how agents talk about each other in a new deployment.

## Before answering “who works here?”

- Did the runtime inject a **roster block** or mention map? Prefer that over static examples.
- Are you about to **contradict** another bot’s prior message in-thread? Align or defer, don’t fight.

## Handoffs

- Is the ask **in another lane** (sales vs build vs ops)? Name the right specialist and **why**, one sentence.
- Does the product require **`<@U…>`** tokens for peer pings? Only use tokens **shown in host context**, never invented IDs.

## Safety

- No **private metrics**, customer names, or org charts not provided in context.
- If uncertain: short honest uncertainty beats a confident wrong roster.

## Installing this skill

Same as other examples: enabled by default via **`SKILLS_SEED_EXAMPLES`** + **`SKILLS_EXAMPLES_DIR`**, or copy `read-company-directory` into **`SKILLS_MCP_SERVER_DIR`** yourself.
