# Codex Domain Schema — One Mind OS

## Standard Domain File Format

Every domain folder contains a `README.md` as the primary index file.
Sub-projects earn their own `.md` file when complexity warrants it.

---

## Template

```markdown
---
domain: XX
name: Domain Name (Short Title)
agent: AgentName
status: active
todoist_project_id: XXXXXXXXXXXXXXX
review: YYYY-QQ
updated: YYYY-MM-DD
related: ["[[RelatedDomain1]]", "[[RelatedDomain2]]"]
depth: 2   # 1=pillar root  2=domain  3=subfolder  4=deep content
---

## 🎯 Purpose
One sentence — what this domain exists to accomplish.

## 🏔 Goals
- [ ] Goal 1 (Q2 2026)
- [ ] Goal 2 (ongoing)

## ⚡ Active Tasks
- [ ] #p1 Task title | due:YYYY-MM-DD | tid:TODOIST_TASK_ID
- [ ] #p2 Task title | tid:TODOIST_TASK_ID

## 📦 Completed Archive
- [x] YYYY-MM-DD — Task description

## 📎 Context & Notes
Reference docs, decisions, links, SOPs go here.

## 🔗 Related
- [[Domain Name]] — relationship note
```

---

## Frontmatter Field Reference

| Field | Required | Purpose |
| ----- | -------- | ------- |
| `domain` | ✅ | Pillar code (HP / SO / LE / GE) |
| `name` | ✅ | Domain display name |
| `agent` | ✅ | Owning agent (Legacy, Zenith, Unity, Eternal, Sovereignty) |
| `status` | ✅ | active / archived / stub |
| `updated` | ✅ | Last modified date YYYY-MM-DD |
| `review` | ✅ | Next review quarter YYYY-QQ |
| `ai_summary` | ✅ | One sentence — what this file contains. Agent scans this without reading the body. |
| `tags` | ✅ | `[domain/hp, type/index, status/active]` |
| `todoist_project_id` | ✗ | Todoist sync bridge (domains only) |
| `related` | ✗ | Explicit graph edges — `["[[FileA]]", "[[FileB]]"]` — agent reads connections without wikilink parsing |
| `depth` | ✗ | Hierarchy level: `1`=pillar root · `2`=domain · `3`=subfolder · `4`=deep content. Agents prioritize depth 1-2 for orientation, 3-4 for specific answers. |

---

## Priority Tags
- `#p1` = Critical / P1 in Todoist
- `#p2` = High / P2 in Todoist
- `#p3` = Medium / P3 in Todoist
- `#p4` = Low / P4 in Todoist

## Sync Rules
- **Codex = Source of Truth** for projects, goals, context, history
- **Todoist = Execution Layer** for active tasks, due dates, priorities
- `tid:` field = Todoist task ID — the sync bridge between both systems
- Completed tasks move to `## 📦 Completed Archive` section
- Use SYNC-INDEX.md to map domains ↔ Todoist project IDs

## Agent Ownership
| Domain Range | Agent |
|-------------|-------|
| SO 00-10 | ⚡ Legacy / 🌐 Sovereignty |
| HP 25-34 | 🏔 Zenith (Spartan / Oracle skills) |
| LE 50-57 | 🌿 Unity (Haven / Eden skills) |
| GE 75-83 | 👑 Eternal (Empire skills) |
