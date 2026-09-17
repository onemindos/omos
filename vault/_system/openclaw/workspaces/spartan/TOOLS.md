# TOOLS.md — Spartan Tool Usage Guide

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
memory_search "workout routine chest"     → training programs
memory_search "macros protein target"     → nutrition plans
memory_search "sleep protocol magnesium"  → sleep optimization
memory_search "budget savings rate"       → personal finance
memory_search "reading list 2026"         → learning goals
memory_search "morning routine"           → daily systems
```

When to use: ANY question about Zeus's health, fitness, nutrition, finance, or personal ops. Always search before answering personalized questions.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/27 Body (Health & Wellness)/...
read_file /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/30 Finance/...
```

## Todoist (Task Management)

Use for creating and tracking HP-domain tasks:
- Workout reminders
- Meal prep tasks
- Bill payment reminders
- Learning goals and book reading

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
edit /home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/27 Body (Health & Wellness)/existing-note.md
```

## Weather

Check for outdoor training conditions:
- Temperature, precipitation, wind for outdoor workouts
- Severe weather alerts that affect training plans

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `25-49 HP/25 Identity/` | Self-awareness, values, personal purpose |
| `25-49 HP/26 Intelligence/` | Learning, books, courses, research |
| `25-49 HP/27 Body/` | Fitness, sleep, nutrition, biometrics |
| `25-49 HP/28 Mastery/` | Skills, certifications, career |
| `25-49 HP/29 Life Systems/` | Routines, habits, productivity |
| `25-49 HP/30 Finance/` | Budget, banking, credit, investments |
| `25-49 HP/32 Relationships/` | Social network, friendships |
| `25-49 HP/33 Joy/` | Recreation, travel, hobbies |
| `25-49 HP/34 Spirituality/` | Mindset, philosophy, rituals |
