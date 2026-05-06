---
name: meeting-prep-brief
description: Produce a tight pre-meeting brief from stated goals, attendees, and constraints. Use when someone asks for meeting prep, a 1:1 outline, or "what should we cover" before a call.
license: MIT
metadata:
  author: bimross-examples
  version: "1.0"
  kind: instruction-only
compatibility: Requires only conversation context; optional calendar or doc URLs pasted by the user.
---

# Meeting prep brief

## When to activate

Use this skill when the user wants structure before a meeting—not live facilitation during the meeting.

## Inputs to infer or ask for once

If missing, ask a single short follow-up (one message) for the most blocking gap only.

- **Meeting type** (1:1, standup, customer, internal decision, etc.)
- **Goal** (decision, alignment, status, discovery)
- **Attendees / roles** (who decides, who needs to leave aligned)
- **Time budget** (15 / 30 / 45 / 60 minutes)

## Output shape

Produce **five sections**, concise bullets, no preamble fluff:

1. **Objective** — one sentence outcome.
2. **Agenda** — ordered topics with rough time boxes summing to the budget.
3. **Decisions needed** — explicit "we need to decide X by end."
4. **Risks / unknowns** — what could derail or remains fuzzy.
5. **Follow-ups** — owners + deadlines only where obvious.

## Tone

Direct, operator-ready. No motivational filler.

## Optional depth

If the user references long background, summarize constraints only—then point to [references/CHECKLIST.md](references/CHECKLIST.md) for a longer reusable checklist pattern.
