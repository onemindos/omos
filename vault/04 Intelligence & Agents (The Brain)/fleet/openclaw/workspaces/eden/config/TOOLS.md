# TOOLS.md — Eden Tool Usage Guide

## Codex Access

The Codex is at `/home/onemind/onemind-os/codex/` — your native knowledge search is built-in.

### ⛔ Codex Write Rules (CRITICAL)

The codex directory structure is **immutable** — managed exclusively by Zeus from the local IDE.

- **NEVER create new folders** in the codex. The numbered folder hierarchy (00-12, 25-34, 50-58, 75-84) is fixed.
- **NEVER rename or renumber** existing codex folders.
- **New notes go to `00-24 SO/06 Inbox (Queue)/`** — Zeus will file them to the correct location.
- **If a folder doesn't exist for a topic**, put the note in `06 Inbox (Queue)/` or `99-organize/` — do NOT create a new folder.
- The codex syncs via git. Any rogue folders you create get auto-committed and pollute the repo.

### memory_search (primary — use this first)

Semantically searches your memory + the Codex vault. Uses BM25+vector hybrid search.

```
memory_search "kids school schedule"        → education plans
memory_search "family medical appointments" → health records
memory_search "estate planning trust"       → wealth/legal
memory_search "date night ideas"            → marriage
memory_search "elder care checklist"        → parent support
memory_search "family traditions holidays"  → heritage
```

When to use: ANY question about family, children, marriage, health, estate, or family logistics. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/51 Marriage/...
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/52 Children/...
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/55 Health/...
```

## Todoist (Task Management)

Use for creating and tracking LE-domain tasks:
- School assignments and deadlines
- Medical appointment reminders
- Family event preparation
- Estate document reviews

## File Tools (read / write / edit)

You have direct file access to the codex:

- `read` — read any codex file by path
- `write` — create new notes (ONLY in `06 Inbox (Queue)/` or `99-organize/`, or update existing files in your domain)
- `edit` — make precise edits to existing codex documents

To create a new note:
```
write /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/06 Inbox (Queue)/My-New-Note.md "# Title\n\nContent here"
```

To update an existing file in your domain:
```
edit /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/52 Children/existing-note.md
```

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `50-74 LE/51 Marriage/` | Partnership, communication, planning |
| `50-74 LE/52 Children/` | Parenting, homeschool, development |
| `50-74 LE/53 Parents/` | Elder care, parent support |
| `50-74 LE/54 Heritage/` | Family history, traditions |
| `50-74 LE/55 Health/` | Family wellness, medical records |
| `50-74 LE/56 Wealth/` | Estate, family office, generational wealth |
| `50-74 LE/57 Legal/` | Family legal, trusts, insurance |
| `50-74 LE/58 Community/` | Extended family, neighborhood |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
