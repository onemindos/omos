---
name: onemind-maintainer
description: Read, search, write, and maintain the One Mind vault. Reads are free — all writes (create, edit, delete, move) require HITL approval via Telegram/Discord before execution. Includes vault health audit checklist for monthly reviews.
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Codex Skill: One Mind Maintainer

Read, search, write, and maintain the One Mind vault. All mutations require HITL approval before execution — approvals route to Telegram and Discord.

---

## Permission Rules — NON-NEGOTIABLE

| Operation | Permission |
|-----------|-----------|
| Read file | ✅ Free — no approval needed |
| Search / grep | ✅ Free — no approval needed |
| List directory | ✅ Free — no approval needed |
| Create file | 🔒 HITL approval required |
| Edit / overwrite file | 🔒 HITL approval required |
| Delete / move file | 🔒 HITL approval required |

Before ANY write, edit, or delete — submit an approval request via exec-approvals. Do NOT proceed until approved. If rejected, surface the reason and stop.

---

## Approval Request Format

Always include:
1. **Operation** — CREATE / EDIT / DELETE
2. **File path** — full path relative to vault root
3. **Why** — one sentence reason
4. **What changes** — exact content or diff summary

---

## Vault Location

`/home/ubuntu/onemind` (VPS: omos-core)

---

## Vault Architecture

```
onemind/
├── 00-24 SO (Sovereign Operations)/   ← Operating layer (AI, infra, security, dev)
├── 25-49 HP (Holistic Performance)/   ← Personal pillar (body, mind, money)
├── 50-74 LE (Legacy Evolution)/       ← Home pillar (family, community)
├── 75-99 GE (Generational Empire)/    ← Business pillar (ventures, brand, wealth)
├── _codex/                            ← Agent context, skills, templates, memory
└── _system/                           ← Operational code only (scripts, configs)
```

---

## Domain Routing

Route content to the correct pillar before writing:

| Content Type | Destination |
|-------------|------------|
| AI, agents, tech, automation, security | `00-24 SO/` |
| Personal growth, fitness, habits, money | `25-49 HP/` |
| Relationships, community, home | `50-74 LE/` |
| Career, ventures, brand, wealth | `75-99 GE/` |
| Unrouted captures | `00-24 SO/99 Inbox (Queue)/` |

---

## Agent Domain Ownership

Only write to your assigned domain. Legacy may write across all.

| Agent | Domains |
|-------|---------|
| ⚡ Legacy | Orchestrator — all pillars |
| 💪 Spartan | HP 27-28 — body, fitness |
| 🔮 Oracle | HP 25-26, 29-34 — mind, money |
| 🌱 Eden | LE — sustainability |
| 🏡 Haven | LE 50-74 — family, home |
| 💰 Empire | GE 75-99 — ventures, brand, wealth |
| 🛡️ Guardian | SO 04 — security |
| ⚙️ Grid | SO 01-02 — infrastructure |
| 🔬 Trinity | SO 06 — research |

---

## Standard Frontmatter

```yaml
---
domain: XX
name: Domain Name (Short Title)
agent: AgentName
status: active
todoist_project_id: XXXXXXXXXXXXXXX
review: YYYY-QQ
updated: YYYY-MM-DD
---
```

## Standard Sections

```
## 🎯 Purpose
## 🏔 Goals
## ⚡ Active Tasks
## 📦 Completed Archive
## 📎 Context & Notes
## 🔗 Related
```

---

## Priority Tags

- `#p1` = Critical
- `#p2` = High
- `#p3` = Medium
- `#p4` = Low

Active tasks use `tid:TODOIST_TASK_ID` as the sync bridge to Todoist.

---

## Task Sync Rule

- **Vault** = source of truth for goals, context, history
- **Todoist** = execution layer for active tasks and due dates
- Completed tasks move to `## 📦 Completed Archive`

---

## Terminology Rules

| ❌ Never use | ✅ Use instead |
|-------------|---------------|
| quadrant / zone | pillar |
| four quadrants | 4 pillars or 3 life pillars + 1 operating layer |
| Codex Method | CODE× Method |
| OneMind Codex | One Mind |
| enterprise / endeavor | Empire (internal), Business (public) |

---

## Commit Behavior

Do NOT manually commit or push. The `onemind-watch.service` on omos-core detects all file changes, auto-commits with message `agent: <filename>`, and pushes to GitHub. Write the file correctly after approval — the rest is automatic.

---

## Vault Health Audit

Run during monthly reviews or when the vault feels disorganized.

### Checklist

1. **Inbox Hygiene** — Count items in `99 Inbox (Queue)/`. If >20, process immediately. Items >30 days old need triage.
2. **Orphan Detection** — Find notes with no incoming links and no tags. Add links, add tags, move to correct domain, or archive.
3. **Frontmatter Completeness** — Every file needs: `title`, `type`, `created`, `tags`. Fix missing fields.
4. **Broken Links** — Check `[[wiki-links]]` to non-existent notes. Update renamed targets, remove dead links, create stubs for planned notes.
5. **Stale Content** — Notes with `modified:` >6 months old. Review accuracy, update or archive.
6. **Duplicate Detection** — Notes covering same topic in different locations. Merge, keep one authoritative version.
7. **Tag Consistency** — Tags follow `namespace/value` pattern. Merge variant spellings.
8. **README Coverage** — Every domain folder has a `README.md`. Create missing ones, update stale ones.

### Schedule

| Task | Frequency |
| --- | --- |
| Inbox hygiene | Weekly |
| Orphan detection | Monthly |
| Frontmatter completeness | Monthly |
| Broken links | Monthly |
| Stale content review | Quarterly |
| Duplicate detection | Quarterly |
| Tag consistency | Quarterly |
| README coverage | After adding new domains |

### Severity

| Issue | Severity |
| --- | --- |
| Inbox overflow (>20 items) | High — process immediately |
| Missing frontmatter | Medium — fix immediately |
| Broken links | Medium — fix or create stubs |
| Duplicate notes | Medium — consolidate during next review |
| Orphaned notes | Low — add context when convenient |
| Stale content | Low — review during quarterly audit |
