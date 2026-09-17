# HP Domain Audit Prompt (25-49 Holistic Performance)

> Paste this into a new Copilot chat session. This agent will audit and clean HP while the main session works on SO.

---

## PROMPT START — COPY BELOW THIS LINE

You are auditing the **HP (Holistic Performance)** domain group of the OneMind Codex vault. This is a folder-by-folder production cleanup. Zeus (the vault owner) is running parallel sessions — you handle HP while another handles SO.

### Context Files — Read These First

Read these files before touching anything:
- `/Users/zeuslegacy/Documents/onemind-codex/AGENTS.md` — full vault context
- `/Users/zeuslegacy/Documents/onemind-codex/25-49 HP (Holistic Performance)/README.md` — domain overview

### Your Domain

**HP — Holistic Performance (25-49): YOUR SELF**
Path: `/Users/zeuslegacy/Documents/onemind-codex/25-49 HP (Holistic Performance)/`

Current folders:
```
25 Identity (Self, Purpose & Spirit)/
26 Intelligence (Education & Learning)/    ← 36 GB (biggest folder)
27 Body (Health & Wellness)/
28 Mastery (Skills & Career)/
29 Life Systems (Daily Ops)/
30 Finance (Money)/                        ← DUPLICATE with next
30 Finance (Stewardship & Wealth)/         ← DUPLICATE — needs merge
31 Sovereignty (Legal)/
32 Relationships (Social)/
33 Joy (Recreation & Travel)/
34 Spirituality (Metaphysics)/
_from-research/                            ← files routed from SO/Research
_legacy/                                   ← 4.5 GB old numbering content
_unsorted-pdfs/                            ← 412 MB needs triage
_unsorted-images/                          ← empty
images/                                    ← 25 MB — should be in domain folders
pdfs/                                      ← empty
```

Loose files at domain root:
- `Drivers License Management.md` — route to proper domain folder
- `Essential_For_Your_Well-Being__Your_Safety_And_Security-_Part_2.md` — route or _to-delete

### Known Issues to Fix
1. **Duplicate domain 30** — `Finance (Money)` vs `Finance (Stewardship & Wealth)` — merge into one
2. **_legacy/ (4.5 GB)** — old numbering system content, needs triage
3. **_unsorted-pdfs/ (412 MB)** — needs classification and routing
4. **images/ folder** — stray images, should be in domain subfolders or _to-delete
5. **Loose .md files** at domain root — route to proper folders
6. **Empty folders** — `pdfs/` and `_unsorted-images/` can go to _to-delete

---

## THE PROCESS (follow this exactly)

### Phase 1: Audit Each Domain Folder (one at a time)

For each numbered domain folder (25-34), do this:

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

- **Merge duplicate domain 30** — combine both Finance folders into one, move the emptied one to _to-delete
- **Triage _legacy/** — classify contents, route what's useful, move rest to _to-delete
- **Triage _unsorted-pdfs/** — identify each PDF, route to proper domain folder
- **Route _from-research/** — files already routed from SO, verify they're in the right HP subfolder

---

## RULES (critical — follow these exactly)

### Naming Convention
- **Files**: `lowercase-with-hyphens.md` (e.g., `weekly-review.md`, `workout-plan.md`)
- **Domain folders**: `XX Name (Subtitle)/` (e.g., `30 Finance (Stewardship & Wealth)/`)
- **No spaces in filenames**, no Title-Case, no underscores

### _to-delete Pattern
- **NEVER delete files directly** — always move to `_to-delete/` within the domain group
- Create `_to-delete/` at `/Users/zeuslegacy/Documents/onemind-codex/25-49 HP (Holistic Performance)/_to-delete/` if it doesn't exist
- Organize by source: `_to-delete/from-25-identity/`, `_to-delete/from-26-intelligence/`, etc.

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
| `#pillar/HP` | `#domain/hp` |
| `Your Family` (for LE) | `Your Foundation` |
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
tags: [domain/hp, type/note]
ai_summary: "One-line summary (≤20 words)"
---
```

### Do NOT
- Delete any files without Zeus's explicit approval
- Rename domain group folders (the `25-49 HP (Holistic Performance)` name is locked)
- Create new top-level folders
- Move files out of _to-delete once placed there
- Change file content beyond terminology fixes and frontmatter

---

## WORKFLOW

Start with: **"I'm ready to audit HP. Which domain folder should I start with?"**

Then work through each folder one at a time. Present findings, wait for approval, execute. Repeat.

Suggested order: Start with the smallest folders first (29, 32, 33, 25) to build momentum before the big ones (26 at 36 GB, 30 merge, _legacy triage).

**GO.**


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
