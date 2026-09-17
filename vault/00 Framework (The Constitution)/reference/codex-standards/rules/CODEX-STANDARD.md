---
id: so-00-codex-standard-20260430
title: Codex Standard — The Law
type: reference
domain: SO
area: "00"
agent: Legacy
status: active
tags:
  - domain/so
  - area/00
  - type/reference
  - topic/standard
  - topic/taxonomy
todoist_project_id: 6g6P9mmQ3Pf36qjG
updated: 2026-04-30
related:
  - "[[00 Framework README]]"
  - "[[conventions]]"
  - "[[00-24 SO (Sovereign Operations)/00 Framework (Doctrine)/reference/codex-standards/Naming-Conventions]]"
  - "[[frontmatter]]"
  - "[[navigation-protocol]]"
  - "[[routing-rules]]"
  - "[[ONEMIND-CODEX]]"
---
# Codex Standard — The Law

> *Every file in the codex follows this standard. This is what makes the system readable by humans, navigable by AI, and renderable as a graph in Obsidian.*

---

## 1. Frontmatter (Required on Every File)

```yaml
---
id: <domain>-<area>-<slug>-<YYYYMMDD>     # unique, stable — never changes after creation
title: "Human-readable title"
type: <type>                               # see types below
domain: SO | HP | LE | GE
area: "00"                                 # folder number
agent: Legacy | Spartan | Oracle | Eden | Empire
status: active | draft | archive
tags: [domain/so, area/00, type/doctrine, topic/xyz]
todoist_project_id: <id>                   # links to Todoist
updated: YYYY-MM-DD
related:                                   # explicit graph edges
  - "[[Related File Name]]"
  - "[[Another File]]"
---
```

**Required**: `id`, `title`, `type`, `domain`, `area`, `status`, `updated`  
**Recommended**: `agent`, `tags`, `related`, `todoist_project_id`

---

## 2. File Types

| type | Use For |
|------|---------|
| `index` | Folder README hub nodes |
| `doctrine` | Canonical truth — vision, constitution, principles |
| `architecture` | System design, data flow, structural decisions |
| `reference` | Standards, conventions, how-to guides |
| `sop` | Step-by-step procedures, runbooks |
| `skill` | OpenClaw skill specs and documentation |
| `concept` | Ideas, research, exploration |
| `note` | Raw captures, Apple Notes imports, brainstorms |
| `template` | Reusable document templates |
| `stub` | Placeholder — content needed |

---

## 3. Tag Taxonomy

```
domain/so     domain/hp     domain/le     domain/ge

area/00   area/01   area/02   area/03   area/04
area/05   area/06   area/25   area/26   ...

type/doctrine   type/architecture   type/reference
type/sop        type/skill          type/concept
type/note       type/template       type/index   type/stub

topic/agents       topic/activation    topic/cadence
topic/constitution topic/governance    topic/infrastructure
topic/openclaw     topic/nats          topic/memory
topic/skills       topic/vision        topic/philosophy
topic/security     topic/devops        topic/research
topic/health       topic/finance       topic/family
topic/business     topic/brand         topic/marketing

status/active   status/draft   status/archive

agent/legacy   agent/spartan   agent/oracle   agent/eden   agent/empire
```

---

## 4. Folder Structure Rules

Every folder has:
- `README.md` — the hub node (type: index) with full cross-links
- Active files at root (no more than ~15 before subfoldering)
- `drafts/` — raw captures, Apple Notes exports, brainstorms (type: note)
- `archive/` — superseded, outdated, or replaced versions

Naming:
- Files: `kebab-case.md` or `SCREAMING-SNAKE.md` for canonical docs
- Folders: `kebab-case/` — descriptive, not abbreviated
- No spaces in new file names (existing files with spaces are OK)

---

## 5. Hub Node Pattern (README.md)

Every folder README must have:
1. Frontmatter with `type: index` and full `related[]` array
2. Domain breadcrumb: `[[ONEMIND-CODEX]] → Domain → Area → **Here**`
3. Table of all files in the folder
4. Folder map showing subfolders
5. Navigation block showing sibling folders
6. Active Tasks section (syncs to Todoist)

---

## 6. See Also Pattern (Every Non-README File)

Every canonical file ends with:

```markdown
---

## See Also

- [[Folder README]] — parent area
- [[Related File]] — why it's related
- [[Sibling README]] — cross-domain link
- [[ONEMIND-CODEX]] — root doctrine
```

This creates the **graph edges** that power:
- Obsidian graph view (visual clusters by domain/area)
- RAG traversal (AI follows edges to find related context)
- Human navigation (next logical doc to read)

---

## 7. Domain → Agent → Todoist Mapping

| Domain | Folders | Agent | Todoist ID |
|--------|---------|-------|------------|
| SO — Sovereign Operations | 00–24 | Legacy | `6g6P9mcc7JMjfPrF` |
| SO 00 Framework | — | Legacy | `6g6P9mmQ3Pf36qjG` |
| SO 01 Intelligence | — | Legacy | `6g6P9rCGG89q6qPr` |
| SO 02 Infrastructure | — | Legacy | `6g6P9p83V2666MP3` |
| SO 03 Command | — | Legacy | `6gM8gPFJ8MM85hv5` |
| SO 04 Security | — | Legacy | `6g6P9qVm3mVVvV35` |
| SO 05 Development | — | Legacy | `6g6P9pqrf5PCR9P5` |
| SO 06 Research | — | Legacy | `6g6P9qrcrJXWHWqF` |
| HP — Holistic Performance | 25–49 | Spartan/Oracle | `6g6P9w4gg33h49FX` |
| LE — Legacy Evolution | 50–74 | Eden | `6g6PGPpm3XgJrPGR` |
| GE — Generational Entrepreneurship | 75–99 | Empire | `6g6PGWvG4f9C2Xr5` |

---

## See Also

- [[00 Framework README]] — parent area
- [[conventions]] — codex writing conventions
- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine)/reference/codex-standards/Naming-Conventions]] — naming rules detail
- [[frontmatter]] — frontmatter field guide
- [[navigation-protocol]] — how agents navigate
- [[routing-rules]] — where to file new content
- [[ONEMIND-CODEX]] — root doctrine
