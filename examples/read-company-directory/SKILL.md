---
name: read-company-directory
description: Read the squad / company directory—who each agent is and what they own. Use when someone asks who works here, names the company directory skill, wants introductions to the team, or asks which specialist handles a lane (sales, build, ops, research, EA).
license: MIT
metadata:
  author: bimross-examples
  version: "1.0"
  kind: instruction-only
compatibility: Works from conversation and any roster the host injects into instructions; does not require external APIs.
---

# Read company directory (squad roster)

## Role

This skill is **markdown-only guidance**: it teaches models how to stay coherent about **who else is in the squad** and **who owns which lane**, so agents can answer “who works here?” and recommend handoffs without contradicting each other.

It is **not** an HR database. The **authoritative** roster for a live deployment is whatever the runtime already injected (system prompts, roster blocks, canonical `<@U…>` mention maps). This file is the **fallback pattern** when the model needs explicit instructions.

## Truth order (highest wins)

1. **Runtime-injected roster** — prepended context, mention tokens, “you are X / your siblings are …” blocks.
2. **What the human said in this thread** — names, roles, “ask Joanne for …”.
3. **This skill** — generic squad model and the reference snapshot below.

If (1) or (2) conflicts with the snapshot, **follow (1) or (2)**.

## Squad model (shared vocabulary)

- Several **specialist agents** cover different outcomes (sales judgment, simplification, automation/build, research intern-style support, executive ops, etc.).
- **Same workspace, different lanes:** prefer handing off to the right lane instead of guessing outside it.
- **Slack / chat UX:** when the product expects it, **peer references use real mention tokens** (`<@U…>`) supplied by the host—not guessed display names—so notifications and routing stay correct.

## Reference snapshot (example deployment)

Use this table only when no richer roster is injected. Treat it as **illustrative**; swap names/roles if your environment defines them differently.

| Agent (example) | Typical lane |
|-------------------|--------------|
| Alex | Revenue, offers, GTM framing |
| Tim | Experiments, systems, networking / decision quality |
| Ross | Automation, implementation, shipping and infra risk |
| Garth | Research, synthesis, follow-through |
| Joanne | Executive operations, docs/email flows, channel onboarding patterns |

When listing “who works here,” keep answers **short and role-shaped**: one line per specialist, no fake tenure or private metrics.

## When to activate

- **Discovery:** “Who works here?”, “What do you each do?”, “Who should I ask about X?”
- **Explicit skill-shaped asks:** “read company directory”, “show the squad”, “who is on the team?”
- **Routing:** recommending another specialist for a task outside your lane.
- **Coordination:** suggesting a handoff without duplicating another agent’s responsibility.

## When not to refuse

- Do **not** answer “I can’t discuss personnel at your organization” when the user means **this product’s configured agent roster** (the squad). That is **not** a request for private employee records; answer from injected context + the reference snapshot, within policy.

## When not to lean on this skill alone

- **Inventing** extra teammates, customers, or titles not in context.

## Output expectations

- **Direct answers** first; optional one-line handoff suggestion if useful.
- If roster is unknown after checking injected context, say you **don’t have the live roster in this session** and suggest the human name the right teammate or check `/admin` / product roster surfaces—**do not** fabricate names.

## Deeper checklist

For repeat asks or onboarding flows, use [references/CHECKLIST.md](references/CHECKLIST.md).
