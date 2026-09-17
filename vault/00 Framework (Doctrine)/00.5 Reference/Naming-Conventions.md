---
id: ui-00-20260111053925
title: "Naming Conventions"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/00]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Naming Conventions"
---

# Naming Conventions

> **File and folder naming standards for the One Mind vault**

---

## Folder Naming

### Format: `XX Name (Subtitle)/`

**Pattern:**
- `XX` = Domain number (00-99)
- `Name` = Primary name
- `(Subtitle)` = Optional clarification

**Examples:**
- `00 Framework (Doctrine)/`
- `25 Identity (Self & Spirit)/`
- `75 Zeus Delacruz Enterprise (Personal Brand)/`

---

## Sub-Folder Naming

### Format: `XX.Y Name/` or `XX.Y.Z Name/`

**Pattern:**
- `XX` = Parent domain
- `Y` = Sub-category
- `Z` = Sub-sub-category (optional)

**Examples:**
- `02.3 Agents/`
- `02.3.14 Memory Block Architecture/`

---

## File Naming

### General Files: `Name-With-Dashes.md`

**Pattern:**
- Use dashes (kebab-case)
- Capitalize first letter of each word
- `.md` extension

**Examples:**
- `Agent-Configuration.md`
- `Daily-Rhythm-Template.md`
- `Cross-Pillar-Routing.md`

---

### Special Files

| File Type | Pattern | Example |
|-----------|---------|---------|
| Index/Navigation | `INDEX.md` or `README.md` | `INDEX.md` |
| Master docs | `ALL-CAPS.md` | `ONEMIND-CODEX.md` |
| Dates | `YYYY-MM-DD.md` | `2026-01-15.md` |
| Quarterly | `QX-YYYY-Roadmap.md` | `Q1-2026-Roadmap.md` |
| Weekly | `Week-YYYY-WW.md` | `Week-2026-03.md` |

---

## Domain-Specific Prefixes

### Prevent Overlap Across Domains

| Domain | Prefix Pattern | Example |
|--------|---------------|---------|
| **GE** | Business/Entity names | `Legacy Acadacom`, `Zeus Delacruz Enterprise` |
| **HP** | `Zeus -` | `Zeus - Personal Finance`, `Zeus - Fitness Plan` |
| **LE** | `Delacruz -` | `Delacruz - Family Wealth`, `Delacruz - Estate` |
| **UI** | `XX-[Name] Ops` | `01-Legacy AI Ops`, `02-Appstack Ops` |

---

## Wikilink Conventions

### Internal Vault Links: `[[Path/To/File]]`

**Patterns:**
- Full path from vault root: `[[00-24 UI/02 Agents/Agent-Config]]`
- Relative path: `[[../02 Agents/Agent-Config]]`
- Display text: `[[ONEMIND-CODEX|The Master Doc]]`

---

## Tag Conventions

### Tags: `#category/subcategory`

**Common patterns:**
- `#pillar/HP` - HP domain content
- `#pillar/LE` - LE domain content
- `#pillar/GE` - GE domain content
- `#pillar/UI` - UI domain content
- `#status/active` - Active project
- `#status/archived` - Archived content
- `#type/project` - Project file
- `#type/area` - Area file
- `#type/workflow` - Workflow/protocol

---

## Frontmatter (YAML) Template

```yaml
---
id: note-{unique-id}
title: "Note Title"
type: concept | procedure | reference | experience | project | area
pillar: GE | HP | LE | UI | cross
status: inbox | draft | active | complete | archived
created: 2026-01-01
modified: 2026-01-01
owner: zeus
tags: [topic/subtopic, pillar/name]
links: []
ai_summary: "One-line summary for agent parsing"
---
```

---

## What NOT to Do

### ❌ Bad Practices

| Don't Do | Why | Do Instead |
|----------|-----|------------|
| `My Notes 2026.md` | Spaces, unclear | `Zeus-Notes-2026.md` |
| `stuff.md` | Vague | `Project-Research-Notes.md` |
| `temp123.md` | Meaningless | `Draft-Marketing-Plan.md` |
| `XXX TODO.md` | Caps, space | `Action-Items.md` |
| `file (2).md` | Duplicate naming | `File-v2.md` or rename original |
| `really_long_file_name_that_goes_on.md` | Too long, underscores | `Concise-Name.md` |

---

## Archive Naming

### When Archiving Content

**Pattern:** Move to Archive folder, preserve original name

**Structure:**
```
XX Archive/
├── YYYY/                    ← Year folder
│   ├── Original-File.md     ← Keep original name
│   └── Project-Name/        ← Keep structure if multi-file
```

**Example:**
```
00-24 UI/09 Archive/
├── 2025/
│   ├── Old-Agent-Config.md
│   └── Deprecated-Workflow.md
```

---

## Project Naming

### Active Projects: Clear, Action-Oriented

**Pattern:** `{Verb}-{Object}.md` or `{Goal-Outcome}.md`

**Examples:**
- `Launch-Credit-Repair-Business.md`
- `Achieve-750-Credit-Score.md`
- `Build-Homestead-Infrastructure.md`

---

## Agent-Specific Naming

### Agent Configuration Files

**Pattern:** `{Agent-Name}-Agent-Config.md`

**Examples:**
- `HP-Agent-Config.md`
- `LE-Agent-Config.md`
- `GE-Agent-Config.md`
- `UI-Agent-Config.md`

### Skills

**Pattern:** `skill-{skill-name}.md`

**Examples:**
- `skill-morning-briefing.md`
- `skill-inbox-processing.md`
- `skill-metric-tracking.md`

---

## Business Entity Naming

### Pattern: Clear Legal Names

**Format:** `Legal Entity Name`

**Examples:**
- `Legacy Platform Holdings`
- `Legacy Acadacom`
- `Legacy Restore`

**Folder structure:**
```
75-99 GE/
├── 76 Legacy Acadacom/
│   ├── Entity-Info.md
│   ├── Financial-Records.md
│   └── Legal-Documents/
```

---

## Quick Reference Table

| Content Type | Naming Pattern | Example |
|--------------|----------------|---------|
| Folders (top-level) | `XX Name (Subtitle)/` | `00 Framework (Doctrine)/` |
| Sub-folders | `XX.Y Name/` | `02.3 Agents/` |
| Files | `Name-With-Dashes.md` | `Agent-Configuration.md` |
| Dates | `YYYY-MM-DD.md` | `2026-01-15.md` |
| Projects | `{Verb}-{Object}.md` | `Launch-Product.md` |
| Archives | Original name in `XX Archive/YYYY/` | `09 Archive/2025/Old-File.md` |
| Tags | `#category/subcategory` | `#pillar/HP` |

---

## Name Changes & History

**When renaming files:**
1. Update all [[wikilinks]] that reference it
2. Add redirect note in old location (optional)
3. Log the change in daily note
4. Update agent memory if significant

**Example redirect note:**
```markdown
# Old File Name

This file has been renamed to:
[[New-File-Name]]

Moved on: 2026-01-15
Reason: Clarity and consistency
```

---

*Follow these conventions for consistency across the vault. AI agents rely on predictable naming.*
