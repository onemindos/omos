---
name: onemind-sync
description: Keep the One Mind vault and Todoist in sync. Vault is source of truth for goals and context. Todoist is the execution layer. The tid: field is the bridge. Use when syncing tasks, pulling vault state, or pushing git changes.
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Codex Skill: One Mind Sync

Keep the vault and Todoist in sync. The vault is the source of truth for goals and context. Todoist is the execution layer for active tasks and due dates. The `tid:` field is the bridge between both.

---

## Trigger

Run when:
- A task is completed in Todoist → move to vault archive
- A task is created in vault → ensure it exists in Todoist
- A `tid:` field is missing from an active task
- Zeus asks for a sync check

---

## The Sync Bridge

Every active task in the vault uses this format:
```
- [ ] #p1 Task title | due:YYYY-MM-DD | tid:TODOIST_TASK_ID
```

`tid:` is the Todoist task ID. Without it, the task is untracked.

---

## Sync Rules

| Vault state | Todoist state | Action |
|------------|--------------|--------|
| Task has no `tid:` | — | Create in Todoist, write `tid:` back to vault (HITL) |
| Task marked `[x]` in vault | Still open | Close in Todoist |
| Task closed in Todoist | Still `[ ]` in vault | Move to `## 📦 Completed Archive` in vault (HITL) |
| Task deleted in Todoist | Still in vault | Flag to Zeus — do not auto-delete |
| Task in vault with `due:` | No due date in Todoist | Update Todoist due date |

---

## Priority Mapping

| Vault tag | Todoist priority |
|-----------|----------------|
| `#p1` | P1 (urgent) |
| `#p2` | P2 (high) |
| `#p3` | P3 (medium) |
| `#p4` | P4 (low) |

---

## Domain → Todoist Project Mapping

Find the `todoist_project_id` in each domain's frontmatter:
```yaml
todoist_project_id: XXXXXXXXXXXXXXX
```

Always create tasks in the correct project, not the inbox.

---

## Rules
- Vault is ALWAYS the source of truth for goals, context, history
- Never delete a vault task without Zeus approval
- All vault writes require HITL approval (see onemind-maintainer skill)
- Todoist is execution only — priorities, due dates, completion status
- Completed tasks move to `## 📦 Completed Archive` with date prefix: `YYYY-MM-DD — Task description`
