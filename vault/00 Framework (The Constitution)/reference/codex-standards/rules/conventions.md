---
id: so-00-conventions
title: "Codex Conventions"
type: reference
status: active
created: 2026-01-11
modified: 2026-04-05
owner: zeus
tags: [domain/so, type/reference, status/active]
ai_summary: "Naming, tagging, frontmatter, and file organization rules for the vault"
---

# Codex Conventions

Rules for naming, tagging, frontmatter, and file organization. Follow these for consistency across the vault.

---

## Folder Naming

### Domain Folders: `XX Name (Subtitle)/`

- `XX` = Domain number (00-99)
- `Name` = Primary name
- `(Subtitle)` = Optional clarification

**Examples:** `00 Framework (Doctrine)/`, `25 Identity (Self & Spirit)/`, `75 Brand (Personal & Corporate)/`

### Sub-folders

Use lowercase-hyphen names within domains. Subfolders do not require numbering.

---

## File Naming

- **Lowercase with hyphens**: `weekly-review.md`, `home-assistant-setup.md`
- **No spaces in filenames** (Obsidian handles display names separately)
- **Domain prefix optional**: Only if the file travels outside its folder context
- **Date-stamped logs**: `YYYY-MM-DD.md` for daily notes, journals, logs

### Special Files

| File Type | Pattern | Example |
|-----------|---------|---------|
| Index/Navigation | `readme.md` or `index.md` | `readme.md` |
| Master docs | `ALL-CAPS.md` | `ONEMIND-CODEX.md` |
| Daily notes | `YYYY-MM-DD.md` | `2026-01-15.md` |
| Quarterly | `qX-YYYY-roadmap.md` | `q1-2026-roadmap.md` |
| Weekly | `week-YYYY-WW.md` | `week-2026-03.md` |

### What NOT to Do

| Don't Do | Why | Do Instead |
|----------|-----|------------|
| `My Notes 2026.md` | Spaces, unclear | `zeus-notes-2026.md` |
| `stuff.md` | Vague | `project-research-notes.md` |
| `temp123.md` | Meaningless | `draft-marketing-plan.md` |
| `file (2).md` | Duplicate naming | `file-v2.md` or rename original |

## Frontmatter (YAML)

Every file should have frontmatter between `---` markers at the top:

```yaml
---
id: {domain-code}-{YYYYMMDD}-{slug}
title: "Human-Readable Title"
type: note | decision | sop | project | person | review | doctrine | reference
status: active | draft | archived | deprecated
created: YYYY-MM-DD
modified: YYYY-MM-DD
owner: zeus
tags: [domain/so, type/note, status/active]
related: ["[[Other Note]]"]
ai_summary: "One-line summary for agent context loading"
---
```

### Required fields

- `title` — human-readable name
- `type` — what kind of document (see types below)
- `created` — creation date
- `tags` — at least one tag
- `ai_summary` — **one-line summary** (≤20 words). This is the primary progressive disclosure mechanism — agents read these at L1 to decide which files to open. Without it, agents must open every file to know what's inside.

### Recommended fields

- `id` — unique identifier for linking
- `status` — lifecycle state
- `related` — wiki-links to connected notes

## Document Types

| Type | Use For |
| --- | --- |
| `note` | General knowledge, observations, ideas |
| `decision` | Decision records with context and rationale |
| `sop` | Standard operating procedures / protocols |
| `project` | Project tracking with goals and status |
| `person` | Contact/relationship profiles |
| `review` | Review documents (weekly, monthly, quarterly) |
| `doctrine` | Core principles and frameworks |
| `reference` | External references, API docs, specs |
| `log` | Time-stamped entries (daily, event) |

## Tags

Tags use a `namespace/value` pattern:

- `domain/so`, `domain/hp`, `domain/le`, `domain/ge` — which domain group
- `type/note`, `type/decision`, `type/sop` — document type
- `status/active`, `status/draft`, `status/archived` — lifecycle
- `priority/high`, `priority/low` — urgency
- `agent/legacy`, `agent/spartan` — which agent owns/created this

## Wiki-Links

- Use `[[Note Name]]` to create connections between notes
- Use `[[Note Name|Display Text]]` for custom link text
- Obsidian renders these as clickable links and builds the knowledge graph
- Agents should follow `related:` frontmatter links when exploring context

## Folder Rules

1. Every domain folder must have a `README.md` explaining what belongs there
2. New notes go to `06 Inbox (Queue)/` unless the domain is obvious
3. Never create new top-level folders (00-99 numbering is fixed)
4. Subfolders within domains are fine for grouping related notes
5. Asset files (images, PDFs) go to `assets/` at the codex root

## Git Commit Messages

Every change to the vault should have an informative commit message. This creates an audit trail of how knowledge evolved.

**Format**: `{action}: {description}`

| Action | When |
| --- | --- |
| `capture` | New note added |
| `update` | Existing note edited |
| `organize` | File moved from inbox to domain |
| `consolidate` | Notes merged |
| `defrag` | Domain restructured |
| `reflect` | Review or reflection created |
| `init` | Vault bootstrapped |
| `archive` | Content archived or deprecated |
| `system` | `system/` pinned context changed |

**Examples**:
- `capture: added VPS migration decision record`
- `consolidate: merged 4 NATS notes into nats-reference.md`
- `defrag: restructured 08 Infrastructure — 34 files → 18 files, added subfolders`
- `system: promoted active-goals.md to pinned context`

## Agent Self-Organization

Agents are explicitly permitted and encouraged to:

1. **Reorganize files within domains** — move, rename, create subfolders as needed during defrag
2. **Update frontmatter** — fix missing fields, improve `ai_summary`, update `modified:` dates
3. **Manage `system/`** — promote frequently-needed files to pinned context, demote stale ones
4. **Create subfolders** — when a domain naturally clusters into 3+ subtopics with 5+ files each
5. **Split oversized files** — files over 300 lines should be evaluated for splitting
6. **Maintain the filetree** — the folder hierarchy and filenames are navigational signals; keep them descriptive

**Agents must NOT**:
- Create new top-level domain group folders (00-99 numbering is fixed)
- Delete files without owner confirmation (archive instead)
- Move files out of `system/` that have `pinned: true` in frontmatter
- Modify Framework doctrine files without explicit instruction
