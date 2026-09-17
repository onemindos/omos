# TOOLS.md — Empire Tool Usage Guide

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
memory_search "monthly revenue MRR"          → revenue tracking
memory_search "client pipeline deals"        → deal management
memory_search "content calendar schedule"    → marketing
memory_search "SOP onboarding process"       → operations
memory_search "investment portfolio crypto"  → investments
memory_search "partnership JV terms"         → partnerships
```

When to use: ANY question about business, revenue, clients, marketing, or operations. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/78 Client Pipeline/...
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/80 Marketing/...
read_file /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/83 Investments/...
```

## Zoho CRM (Client & Deal Management)

Use for managing the business pipeline:
- Create and update leads/contacts
- Track deal stages and values
- Log client interactions
- Generate pipeline reports

## Todoist (Task Management)

Use for creating and tracking GE-domain tasks:
- Client follow-up reminders
- Content creation deadlines
- Operational tasks and reviews
- Meeting preparation

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
edit /home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)/78 Client Pipeline/existing-note.md
```

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `75-99 GE/75 Zeus Delacruz Enterprise/` | Personal brand, content |
| `75-99 GE/76 Legacy Acadacom/` | Education platform, courses |
| `75-99 GE/77 Legacy Platform Holdings/` | IP, holding company |
| `75-99 GE/78 Client Pipeline/` | Active clients, deals, revenue |
| `75-99 GE/79 Community (Business)/` | Paid community, member LTV |
| `75-99 GE/80 Marketing/` | Content strategy, social, email |
| `75-99 GE/81 Operations/` | SOPs, team management, tools |
| `75-99 GE/82 Legal (Business)/` | Contracts, compliance |
| `75-99 GE/83 Investments/` | Portfolio, real estate, crypto |
| `75-99 GE/84 Partnerships/` | Strategic alliances, JVs |
