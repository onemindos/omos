---
id: so-01-legacy-sync
title: "Legacy Sync Manifest"
type: sync
domain: SO
area: "01"
status: active
created: 2026-01-11
updated: 2026-05-01
owner: zeus
agent: Legacy
tags: [domain/so, area/01, type/sync, agent/legacy]
ai_summary: "Vault sync queue — what Legacy needs to read and absorb into memory"
---

# LEGACY SYNC MANIFEST

> **Purpose:** This file tells Legacy AI what has changed and needs her attention.
> **Location:** `02 Agents/02.1 Mother AI/02.1.1 Index/_LEGACY-SYNC.md`
> **Protocol:** Legacy checks this file periodically and processes pending items.

---

## How This Works

1. **Zeus or Claude** updates this file when something important changes
2. **Legacy** checks this file on schedule (cron) or when notified
3. **Legacy** processes each `pending` item:
   - Reads the linked document
   - Updates her relevant memory block
   - Marks item as `synced`
4. **Archived** items move to bottom after 7 days

---

## Sync Queue

### Pending Review

```yaml
# Items Legacy needs to process
# Format: - [status] | priority | category | path | summary

- [pending] | P1 | system | LEGACY-AI-SYSTEM.md | Core system doc - persona, tools, protocols
- [pending] | P1 | system | 00-24 SO/.../00 Framework/00.1 Index/LIFE-OS-ARCHITECTURE.md | Master architecture - how 4 pillars connect
- [pending] | P1 | skills | 00-24 SO/.../02 Agents/02.1 Mother AI/02.1.2 Skills/_SKILLS-INDEX.md | Skills framework - what Legacy can do
- [pending] | P1 | navigation | 00-24 SO/.../01 Command/.../01.2 Atlas/01.2.1 Index/_Atlas.md | Vault navigation - all 45 areas mapped
- [pending] | P2 | tasks | 00-24 SO/.../01 Command/.../01.1 Core/01.1.1 Index/TODOIST-STRUCTURE-GUIDE.md | Todoist structure - projects, labels, filters
- [pending] | P2 | tasks | 00-24 SO/.../01 Command/.../01.1 Core/01.1.1 Index/todoist-inbox-query.md | Todoist query blocks for Obsidian
- [pending] | P3 | research | LETTA-RESEARCH-LEGACY-AI.md | Letta technical research - memory blocks, limits
```

### Recently Synced

```yaml
# Items Legacy has processed (moves here after sync)
# Format: - [synced] | date | category | path | what was updated

# (none yet)
```

---

## Memory Block Mapping

When Legacy processes a sync item, she updates the relevant memory block:

| Category | Memory Block | What to Update |
|----------|--------------|----------------|
| `system` | `persona` | Core identity, capabilities, protocols |
| `navigation` | `global_state` | Vault structure, paths, areas |
| `tasks` | `active_tasks` | Task patterns, filters, workflows |
| `user` | `user_context` | Zeus preferences, habits, focus |
| `hp` | `hp_context` | Health, finance, personal data |
| `le` | `le_context` | Family, home, legacy data |
| `ge` | `ge_context` | Business, revenue, ventures data |
| `research` | `archival_memory` | Store in archival, not blocks |

---

## Sync Protocol

```yaml
# Legacy's sync procedure
sync_protocol:
  trigger:
    - cron: "0 */4 * * *"  # Every 4 hours
    - manual: "Legacy, sync"
    - webhook: "On vault push to GitHub"

  steps:
    1. read_file: "_LEGACY-SYNC.md"
    2. parse: Extract [pending] items
    3. for_each_item:
       - read_file: item.path
       - extract: Key information relevant to category
       - update_memory: item.memory_block
       - mark_synced: Update status to [synced] with date
    4. report: "Synced X items, updated Y memory blocks"

  on_complete:
    - Update last_sync timestamp in global_state
    - Log sync summary to archival_memory
```

---

## Adding Items to Queue

### For Zeus (Manual)
Add a line under "Pending Review":
```
- [pending] | P1 | category | path/to/file.md | Brief description
```

### For Claude (When editing vault)
After significant edits, add sync entry:
```
- [pending] | P2 | category | path/to/edited/file.md | What changed
```

### Priority Levels
- `P1` = Critical - sync immediately
- `P2` = Important - sync on next cycle
- `P3` = Low - sync when convenient

---

## Vault Quick Reference

```yaml
vault:
  root: "OneMind-Codex"
  github: "github.com/[repo]/OneMind-Codex"

  pillars:
    UI: "00-24 SO (Sovereign Operations)"
    HP: "25-49 HP (Holistic Performance)"
    LE: "50-74 LE (Legacy Evolution)"
    GE: "75-99 GE (Generational Empire)"

  key_paths:
    atlas: "00-24 SO/.../01.2 Atlas/.../_Atlas.md"
    sync: "00-24 SO/.../02.1 Mother AI/.../_LEGACY-SYNC.md"
    system: "LEGACY-AI-SYSTEM.md"
    tasks: "00-24 SO/.../03 Tasks/..."

  index_pattern: "[area]/[area.1] Index/_Index.md"
```

---

## Document Reference System

### Task → Doc Linking
Every Todoist task can reference vault docs:
```
Task Description:
📁 Path: 25-49 HP/.../30 Finance/30.1 Index/_Index.md
📋 UID: HP-30-001
🔗 [[30 Finance (Money)]]
```

### Doc → Task Linking
Every doc can reference its Todoist filter:
```yaml
# In document frontmatter or footer
todoist_filter: "##HP ##30-Finance"
area_code: "30"
domain: "HP"
```

---

## Last Sync Status

```yaml
last_sync:
  timestamp: null  # Legacy updates this
  items_processed: 0
  memory_blocks_updated: []
  next_scheduled: null
```

---

## Archived Syncs

```yaml
# Old sync items (older than 7 days)
# Format: - [archived] | sync_date | category | path

# (none yet)
```

---

*This manifest is Legacy's source of truth for staying current with the vault.*
*Updated: 2026-01-04*
