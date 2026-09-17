# LE Domain Audit Prompt (50-74 Legacy Evolution)

> Paste this into a new Copilot chat session. This agent will audit and clean LE while the main session works on SO.

---

## PROMPT START — COPY BELOW THIS LINE

You are auditing the **LE (Legacy Evolution)** domain group of the OneMind Codex vault. This is a folder-by-folder production cleanup. Zeus (the vault owner) is running parallel sessions — you handle LE while another handles SO.

### Context Files — Read These First

Read these files before touching anything:
- `/Users/zeuslegacy/Documents/onemind-codex/AGENTS.md` — full vault context
- `/Users/zeuslegacy/Documents/onemind-codex/50-74 LE (Legacy Evolution)/README.md` — domain overview

### Your Domain

**LE — Legacy Evolution (50-74): YOUR FOUNDATION**
Path: `/Users/zeuslegacy/Documents/onemind-codex/50-74 LE (Legacy Evolution)/`

Current folders:
```
50 Home (Sanctuary & Property)/            ← 336 MB
51 Marriage (Partnership & Comms)/          ← 24 KB
52 Children (Parenting & Education)/        ← 5.9 GB (biggest)
53 Parents (Elder Care & Support)/          ← 20 KB
54 Heritage (History & Traditions)/         ← 684 MB
55 Health (Family Wellness)/                ← 72 KB
56 Wealth (Estate & Assets)/                ← 2 MB
57 Legal (Documents)/                       ← 224 MB
58 Community (Extended Network)/            ← 4 KB
_legacy/                                    ← 2.8 GB old numbering content
_sorted-duplicates-docs/                    ← 194 MB — likely safe to delete
```

Loose files at domain root:
- `Security Camera with Blue Iris & Deep Learning.md` — route to 50 Home or SO/Infrastructure

### Known Issues to Fix
1. **_legacy/ (2.8 GB)** — old numbering system content, needs triage
2. **_sorted-duplicates-docs/ (194 MB)** — needs verification, likely all duplicates → _to-delete
3. **Children (52) is 5.9 GB** — likely contains PDFs/images, needs audit
4. **Loose .md file** at domain root — route to proper folder
5. **LE is the MOST PERSONAL domain** — family photos, children's records, legal docs. Be respectful and careful.

---

## THE PROCESS (follow this exactly)

### Phase 1: Audit Each Domain Folder (one at a time)

For each numbered domain folder (50-58), do this:

1. **List all files** — `find <folder> -type f | head -100` and `du -sh <folder>`
2. **Classify every file** into categories:
   - **PRODUCTION** — clean, useful, belongs here
   - **STALE** — outdated, broken links, references old structure
   - **JUNK** — empty, placeholder, test files, Notion HTML exports, `.DS_Store`
   - **DUPLICATE** — copy of another file in the vault
   - **ROUTE-OUT** — belongs in a different domain
3. **Present your findings to Zeus** — show counts, ask for approval before executing
4. **Execute approved changes**

### Phase 2: Execute Changes

After Zeus approves each folder:

1. **Move junk/stale/duplicates** to `_to-delete/` subfolder (NEVER delete directly)
2. **Route files** to correct domain folders
3. **Rename files** to lowercase-hyphen convention
4. **Fix stale terminology** (see rules below)
5. **Fix/add frontmatter** where missing

### Phase 3: Handle Special Folders

- **Verify _sorted-duplicates-docs/** — check if these are truly duplicates of files in domain folders. If confirmed, move entire folder to _to-delete
- **Triage _legacy/** — classify contents, route useful files, move rest to _to-delete
- **Route loose root files** — `Security Camera with Blue Iris & Deep Learning.md` to proper location

---

## RULES (critical — follow these exactly)

### Naming Convention
- **Files**: `lowercase-with-hyphens.md` (e.g., `family-budget.md`, `home-maintenance.md`)
- **Domain folders**: `XX Name (Subtitle)/` (e.g., `50 Home (Sanctuary & Property)/`)
- **No spaces in filenames**, no Title-Case, no underscores

### _to-delete Pattern
- **NEVER delete files directly** — always move to `_to-delete/` within the domain group
- Create `_to-delete/` at `/Users/zeuslegacy/Documents/onemind-codex/50-74 LE (Legacy Evolution)/_to-delete/` if it doesn't exist
- Organize by source: `_to-delete/from-50-home/`, `_to-delete/from-52-children/`, etc.

### Terminology Fixes (apply to every file you touch)

| OLD (wrong) | NEW (correct) |
|---|---|
| `UI` (for operating layer) | `SO` |
| `UI Agent` | `SO Agent` or `Grid/Guardian` |
| `domain: UI` (frontmatter) | `domain: SO` |
| `id: ui-` (frontmatter) | `id: so-` |
| `Mother AI` | `Central orchestrator` or `SO Agent` |
| `quadrant` | `domain` or `domain group` |
| `four quadrants` | `3 domains + 1 operating layer` |
| `pillar` | `domain` |
| `#pillar/LE` | `#domain/le` |
| `Your Family` (for LE tagline) | `Your Foundation` |
| `Your Business` (for GE) | `Your Empire` |
| `PI Score` / `PI` | `HPI` (Holistic Performance Index) |
| `LI Score` / `LI` | `LEI` (Legacy Evolution Index) |
| `WI Score` / `WI` | `GEI` (Generational Empire Index) |

### Frontmatter
Every .md file should have YAML frontmatter:
```yaml
---
title: "Human-Readable Title"
type: note | decision | sop | project | person | review | doctrine | reference
status: active | draft | archived
created: YYYY-MM-DD
modified: 2026-04-05
owner: zeus
tags: [domain/le, type/note]
ai_summary: "One-line summary (≤20 words)"
---
```

### Do NOT
- Delete any files without Zeus's explicit approval
- Rename domain group folders (the `50-74 LE (Legacy Evolution)` name is locked)
- Create new top-level folders
- Move files out of _to-delete once placed there
- Change file content beyond terminology fixes and frontmatter
- Be cavalier with family/children/legal documents — when in doubt, ask

---

## WORKFLOW

Start with: **"I'm ready to audit LE. Which domain folder should I start with?"**

Then work through each folder one at a time. Present findings, wait for approval, execute. Repeat.

Suggested order: Start with the smallest folders first (58, 53, 51, 55) to build momentum, then tackle the bigger ones (50, 54, 57, 52, _legacy).

**GO.**


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
