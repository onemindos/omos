---
id: so-01-20260412000013
title: "Context Improves Decision Quality"
type: doctrine
status: active
created: 2026-04-12
owner: zeus
tags: [path/so, path/so/01, type/doctrine, prompt-engineering, context, decision-quality, ai]
ai_summary: "Core principle: the quality of context given to an AI directly determines the quality of its decisions and outputs — garbage in, garbage out at scale"
---

# Context Improves Decision Quality

## The Principle

> The quality of context you give an AI directly determines the quality of its decisions, outputs, and actions.

This is not a soft guideline — it is a hard constraint. Every AI system, from a single chat session to a multi-agent architecture, operates within the boundaries of its context window. What lives in that window determines everything.

Garbage context → mediocre outputs.
Structured, precise context → high-quality decisions.

---

## Why It Matters More Than Model Choice

Most people focus on *which* model to use. The bigger lever is *what you give it*.

| Input Quality | What You Get |
|---|---|
| Vague, unstructured prompt | Generic response, frequent misalignment |
| Clear intent + background | Accurate, relevant output |
| Structured JSON context profile | Consistent, repeatable, high-fidelity results |
| Full operating context (identity + goals + constraints) | Agent-level decision-making quality |

A GPT-4 with poor context loses to a GPT-3.5 with excellent context on most practical tasks.

---

## The 4 Layers of Context

Better context doesn't mean more words. It means the *right* information in the *right* structure.

### Layer 1 — Identity
Who is this for? Who is the AI acting as?

```
"You are Legacy, Zeus's central AI orchestrator.
Zeus is a founder building One Mind OS..."
```

### Layer 2 — Intent
What is the actual goal — not just the immediate task?

```
"Goal: write a LinkedIn post that positions One Mind
as the GTD for AI-native life management."
```

### Layer 3 — Constraints
What must be avoided? What rules are non-negotiable?

```
"Never use the word 'quadrant'. Never use motivational fluff.
Keep it under 150 words."
```

### Layer 4 — Examples / Profile
What does "right" look like? Show it, don't just describe it.

```json
{ "tone": ["direct", "systems-thinking"], "avoid": ["emoji", "passive voice"] }
```

---

## In Practice: The One Mind Context Stack

The entire One Mind architecture is built on this principle:

| Layer | What It Does |
|---|---|
| `_codex/system/` files | Load identity, mission, goals, beliefs into every agent session |
| `USER.md` in workspaces | Zeus's preferences, communication style, priorities |
| `CODEX.md` in workspaces | Vault structure, terminology rules, domain map |
| JSON Context Profiles | Portable, paste-anywhere context for brand, image, identity |
| YAML frontmatter | Per-file context: what kind of document is this, what domain, what status |

Each layer adds precision. Together they enable agents to make decisions that feel like Zeus made them.

---

## Decision Quality Formula

$$\text{Decision Quality} = f(\text{Model Capability} \times \text{Context Quality})$$

- You can't fix a bad context with a better model
- You *can* get excellent results from a mid-tier model with exceptional context
- The multiplier effect: better context on a better model compounds

---

## Practical Rules

1. **Before any important prompt** — ask: *what does the AI need to know that it doesn't have?*
2. **Reusable context → JSON profiles** — don't re-explain brand, identity, or constraints every session
3. **Constraints are as important as instructions** — tell it what NOT to do
4. **Show, don't describe** — examples outperform adjectives
5. **Freshen context on long sessions** — context drift is real; re-inject your profile mid-session if outputs degrade

---

## Related Files

- [[JSON-Context-Profile]] — the format for packaging context as reusable profiles
- [[onemind-brand-voice]] — One Mind brand context profile (production-ready)
- [[zeus-identity]] — Zeus identity profile for agent briefing
- [[onemind-image-ajcp]] — image generation context profile
