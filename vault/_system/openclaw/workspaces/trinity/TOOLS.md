# TOOLS.md — Trinity Tool Usage Guide

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
memory_search "Legacy Collective community"     → community info
memory_search "Circle.so spaces and setup"      → platform config
memory_search "content strategy guidelines"     → brand voice
memory_search "email campaign templates"        → campaign patterns
memory_search "Discord server structure"        → channel setup
memory_search "community engagement metrics"    → analytics
```

When to use: ANY question about community, content, or public-facing work. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/75 Community/...
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/76 Marketing/...
read_file /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/04 Resources/...
```

## Circle API (Community Management)

Use for Circle.so operations:
- Create and manage community posts
- Welcome new members
- Monitor space activity
- Manage moderation queue

## Discord (Community Server)

Use for Discord operations:
- Post announcements
- Monitor channel activity
- Respond to community questions
- Manage server engagement

## Content Publisher

Use for publishing content:
- Schedule blog posts
- Format newsletters
- Manage content queue
- Track publishing calendar

## Email Campaigns

Use for email operations:
- Draft newsletters
- Schedule drip sequences
- Track open/click rates
- Manage subscriber lists

## Todoist (Task Management)

Use for tracking community tasks:
- Content deadlines
- Campaign milestones
- Event preparation
- Member follow-ups

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
edit /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/75 Community/existing-note.md
```

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `75-99 GE/75 Community/` | Community strategy, partnerships |
| `75-99 GE/76 Marketing/` | Content strategy, brand guidelines |
| `00-24 SO/04 Resources/` | Templates, brand assets |
