---
name: codex
description: >
  Native interface to the OneMind Codex — the structured knowledge vault. Use when
  Legacy or any agent needs to capture new information, search the vault with tiered
  loading, organize items into 00-99 domains, run reflection cycles, consolidate
  scattered notes, or audit vault health. Activate when user mentions "codex",
  "vault", "knowledge base", "save this", "file this", "what do we know about",
  "weekly review", "organize the inbox", or "vault health".
metadata:
  author: zeus
  version: "1.0.0"
  domain: codex
  emoji: "📚"
---

# Codex — Knowledge Vault Interface

The codex is the source of truth for all knowledge, plans, decisions, and documentation.
This skill provides structured access following the protocols defined in `_codex/`.

## Architecture

```
User request
  → Legacy activates codex skill
    → Tool reads _codex/ for conventions and instructions
      → Executes operation on the vault
        → Returns result
```

The skill reads `_codex/INTERFACE.md` for vault structure, `_codex/CONVENTIONS.md` for
rules, and `_codex/skills/*.md` for operation-specific protocols.

## Available Tools

| Tool                | Description                                                         |
| ------------------- | ------------------------------------------------------------------- |
| `codex_capture`     | Create a new note in the vault with proper frontmatter and template |
| `codex_search`      | Find information using tiered loading (L0 → L1 → L2)                |
| `codex_organize`    | Move an item from inbox to the correct 00-99 domain                 |
| `codex_reflect`     | Run a reflection cycle (daily, weekly, monthly)                     |
| `codex_consolidate` | Merge scattered notes into curated knowledge                        |
| `codex_health`      | Audit vault: orphans, missing frontmatter, inbox overflow           |
| `codex_tree`        | Show folder structure at a given depth                              |

## Quick Reference

### Capture a new note

```
codex_capture({
  title: "NATS Stream Configuration",
  type: "note",
  domain: "08",
  content: "JetStream requires explicit stream creation...",
  tags: ["domain/so", "topic/nats"]
})
```

### Search the vault

```
codex_search({
  query: "VPS setup",
  domain: "08",
  level: "L1"
})
```

Levels: "L0" (folder scan), "L1" (README + file list), "L2" (full file read).
Omit `domain` to search all domain groups. Omit `level` for auto (starts L0, deepens as needed).

### Organize inbox items

```
codex_organize({ process_all: true })
```

Or move a specific file:

```
codex_organize({ file: "06 Inbox (Queue)/nats-notes.md", target_domain: "08" })
```

### Run reflection

```
codex_reflect({ cadence: "daily" })
codex_reflect({ cadence: "weekly" })
codex_reflect({ cadence: "monthly" })
```

### Check vault health

```
codex_health()
```

Returns: inbox count, orphan count, files missing frontmatter, broken links.

### Show folder tree

```
codex_tree({ path: "00-24 SO", depth: 2 })
```

## Environment Variables

| Variable     | Description                                                   |
| ------------ | ------------------------------------------------------------- |
| `CODEX_PATH` | Path to the codex vault (default: auto-detect from workspace) |

## Protocols

This skill follows the operational protocols defined in the vault itself:

- **Search**: `_codex/skills/search.md` — tiered loading (L0/L1/L2)
- **Capture**: `_codex/skills/capture.md` — classification, templates, frontmatter
- **Organize**: `_codex/skills/organize.md` — inbox → domain mapping
- **Reflect**: `_codex/skills/reflect.md` — daily/weekly/monthly cycles
- **Consolidate**: `_codex/skills/consolidate.md` — merge scattered → curated
- **Maintain**: `_codex/skills/maintain.md` — health audit checklist

## Notes

- The codex is markdown on disk — all operations are file reads/writes
- Always follow `_codex/CONVENTIONS.md` for frontmatter format
- Use `_codex/templates/` when creating new notes
- New items go to `06 Inbox (Queue)/` unless the domain is obvious
- Never create new top-level folders in the codex
