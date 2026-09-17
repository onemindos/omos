---
id: so-00-frontmatter
title: "What is Front Matter?"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [domain/so, type/reference, status/active]
ai_summary: "Explains YAML front matter — what it is, how it works, and all required fields for codex files"
---

# What is Front Matter?

Front matter is a block of YAML metadata placed at the **very top** of a Markdown file, enclosed between two lines of triple dashes (`---`).

```markdown
---
title: "My Note"
type: note
status: active
created: 2026-04-12
tags: [type/note, domain/so]
ai_summary: "A one-line summary of what this file is about"
---

# My Note

The actual content of the file starts here.
```

Everything between the `---` markers is metadata *about* the file — not content inside it. Obsidian reads it, agents read it, and the vault build process reads it.

---

## Why It Matters

Front matter is the **primary progressive disclosure mechanism** for the vault.

When an agent needs to understand what's in the vault, it doesn't read every file from top to bottom. That would be slow and expensive. Instead, it reads only the `ai_summary` field from each file's front matter — one line per file — and decides which files are actually relevant to the current task.

Without front matter, an agent must open and read the full content of every file. With front matter, an agent can scan 100 files in the time it takes to read one.

---

## The Full Schema

```yaml
---
id: so-00-20260412-my-note          # Unique identifier: {domain}-{date}-{slug}
title: "Human-Readable Title"        # Required — display name
type: note                           # Required — document type (see types below)
status: active                       # Lifecycle state
created: 2026-04-12                  # Required — ISO date
modified: 2026-04-12                 # Update whenever you edit the file
owner: zeus                          # Who owns this document
tags: [domain/so, type/note]         # Required — at least one tag
related: ["[[Other Note]]"]          # Wiki-links to connected notes
ai_summary: "One-line summary"       # Required — ≤20 words
---
```

---

## Required Fields

| Field | Purpose |
|-------|---------|
| `title` | Human-readable name shown in Obsidian |
| `type` | What kind of document (controls templates and agent behavior) |
| `created` | Creation date in `YYYY-MM-DD` format |
| `tags` | At least one tag for domain routing |
| `ai_summary` | One-line description (≤20 words) — **the most important field** |

---

## Document Types (`type:`)

| Value | Use For |
|-------|---------|
| `note` | General knowledge, observations, ideas |
| `decision` | Decision records with context and rationale |
| `sop` | Standard operating procedures, protocols |
| `project` | Project tracking with goals and milestones |
| `person` | Contact and relationship profiles |
| `review` | Weekly, monthly, quarterly reviews |
| `doctrine` | Core principles and framework definitions |
| `reference` | External references, API docs, specs, how-tos |
| `log` | Time-stamped entries (daily, event-based) |

---

## Tag Namespaces (`tags:`)

Tags use a `namespace/value` pattern. Every file should have at least a `domain/` and `type/` tag.

| Namespace | Values | Example |
|-----------|--------|---------|
| `domain/` | `so`, `hp`, `le`, `ge` | `domain/so` |
| `type/` | `note`, `sop`, `project`, etc. | `type/reference` |
| `status/` | `active`, `draft`, `archived` | `status/active` |
| `priority/` | `high`, `low` | `priority/high` |
| `agent/` | `legacy`, `spartan`, etc. | `agent/legacy` |

---

## The `ai_summary` Rule

The `ai_summary` field deserves special attention. Keep it:

- **≤20 words** — agents read hundreds of these; brevity matters
- **Specific** — say what the file actually contains, not what type of file it is
- **Action-oriented** — "Documents the decision to use NATS over Redis for messaging" beats "A decision record about messaging"

**Bad:** `"A note about infrastructure"`
**Good:** `"VPS setup guide — DigitalOcean droplet, NATS, PostgreSQL, OpenClaw install sequence"`

---

## How Obsidian Uses Front Matter

- The `title` field overrides the filename in the graph view and link previews
- The `tags` field populates Obsidian's tag panel and makes files filterable
- The `status` field can be used in Dataview queries (e.g., show all `status: active` projects)
- The `related` field creates bi-directional links in the knowledge graph

---

## How Agents Use Front Matter

Agents load context in layers (see `INTERFACE.md`):

| Level | What Agents Load |
|-------|-----------------|
| L0 | `_codex/system/` pinned context — always first |
| L1 | `ai_summary` field from every relevant file — fast scan |
| L2 | Full file content — only when L1 says it's relevant |
| L3 | Linked documents via `related:` — only when deeper context needed |

Front matter at L1 means agents can understand the vault's state from reading ~100 one-liners instead of ~100 full documents.

---

## Related

- [[conventions]] — full naming, tagging, and file organization rules
- [[interface]] — agent navigation entry point (progressive disclosure levels)
- `_codex/templates/` — document templates with pre-filled front matter for each type
