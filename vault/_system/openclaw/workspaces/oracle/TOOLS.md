# TOOLS.md — Oracle Tool Usage Guide

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
memory_search "budget monthly savings"       → 30 Finance
memory_search "credit score strategy"        → 30 Finance
memory_search "investment portfolio"         → 30 Finance
memory_search "reading list books 2026"      → 26 Intelligence
memory_search "morning routine"              → 29 Life Systems
memory_search "meditation practice"          → 34 Spirituality
memory_search "relationship check-in"        → 32 Relationships
memory_search "travel plans vacation"        → 33 Joy
memory_search "purpose identity goals"       → 25 Identity
memory_search "habits streak tracker"        → 29 Life Systems
```

When to use: ANY question about Zeus's finances, mindset, learning, relationships, routines, or inner life. Always search before answering personalized questions.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/25 Identity (Self, Purpose & Spirit)/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/26 Intelligence (Education & Learning)/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/29 Life Systems (Daily Ops)/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/30 Finance (Money)/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/34 Spirituality (Metaphysics)/...
```

## Todoist (Task Management)

Use for creating and tracking Oracle-domain tasks:
- Bill payment reminders
- Learning goals and reading targets
- Habit check-ins
- Financial review dates
- Relationship follow-ups

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
edit /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/30 Finance (Money)/existing-note.md
```

## Weather

Check weather for travel planning and outdoor recreation recommendations.

## Web Search & Fetch

For researching:
- Investment strategies and market context
- Book reviews and course evaluations
- Productivity frameworks and habit science
- Travel destinations and planning
