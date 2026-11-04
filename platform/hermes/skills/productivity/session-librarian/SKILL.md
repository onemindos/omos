---
name: session-librarian
description: "Organize sessions by prompt: find, rename, archive, and propose pruning. Deleting needs execute/override."
version: 1.1.0
author: Hermes Agent + Teknium
license: MIT
platforms: [linux, macos, windows]
metadata:
  hermes:
    tags: [Sessions, Organization, Cleanup, Library, Productivity]
    category: productivity
    related_skills: [weekly-review-planning]
---

# Session Librarian

Manage the session library conversationally: find past sessions about a topic,
summarize what they decided, rename them, and propose stale ones for archive.
Always show the plan before touching anything.

## Authority

Deleting or pruning sessions is a deletion under SOUL rule 4. It needs Zeus to
write "execute" or "override" in chat for that specific deletion. A standing
"clean things up" authorizes proposing, and archiving (reversible) is the most
this skill does on its own after the plan is approved. Renaming a single session
Zeus dictated is fine directly.

## When to use

- "What sessions do I have about X?" / "What did we decide about X?"
- "Rename these sessions to something meaningful."
- "Clean up my session library" / "archive the stale ones."
- "Fork that session into a follow-up focused on Y." / "One session per ticket."

## Surfaces

| Task | Surface |
|---|---|
| Find sessions by topic, read content, summarize decisions | `session_search` tool |
| List/filter by metadata (age, source, cost, tokens) | `hermes sessions list` / `stats` |
| Rename | `hermes sessions rename <session_id> <title...>` |
| Reversible soft-hide | `hermes sessions archive <filters>` |
| Destructive | `hermes sessions delete` / `prune` (needs execute/override) |
| Export before touching anything valuable | `hermes sessions export --session-id <id> --format md` |
| Continue elsewhere | `/branch`, or start a new session citing the summary |

## Procedure

1. **Discover.** `session_search(query=..., limit=5-10)` with varied topic words.
   For metadata sweeps, `hermes sessions list --source telegram --limit 50`.
2. **Summarize.** Use the bookends and match window; dump a whole session only when
   asked. Report each as: link, one-line goal, one-line outcome.
3. **Plan before acting (mandatory for anything that mutates).** A table of renames,
   archives, and proposed deletions with the reason (duplicate of which keeper,
   stale, empty). Wait for the go-ahead.
4. **Act with the safest primitive.** Prefer `archive`. For a destructive command,
   run `--dry-run` first, show it, and run it only after Zeus has written
   "execute" or "override" for it. Offer `export --format md` before deleting
   anything with meaningful content.
5. **Report.** Renames applied, archived (count and how to undo: archived sessions
   stay in the DB and show with `--include-archived`), exported, skipped and why.

## Parallel workstreams

To investigate several tickets, use `delegate_task` with one task per workstream
(each subagent gets its own session) and synthesize the summaries. Do not try to
drive other live sessions.

## Pitfalls

- `session_search` finds content, not metadata; combine it with the CLI when a
  request mixes both.
- Titles are identity for `/resume <title>`. Keep them short, unique and
  prefix-friendly, and warn about collisions.
- Archived is not deleted. Say which one you did.
- Cross-profile session links are read-only from another profile; management
  commands act on the current profile's DB.

## Verification

After a pass, re-run the discovery query and `hermes sessions list` to confirm the
library matches the plan.
