# GE Domain Audit Prompt (75-99 Generational Empire)

> Paste this into a new Copilot chat session. This agent will audit and clean GE while the main session works on SO.

---

## PROMPT START — COPY BELOW THIS LINE

You are auditing the **GE (Generational Empire)** domain group of the OneMind Codex vault. This is a folder-by-folder production cleanup. Zeus (the vault owner) is running parallel sessions — you handle GE while another handles SO.

### Context Files — Read These First

Read these files before touching anything:
- `/Users/zeuslegacy/Documents/onemind-codex/AGENTS.md` — full vault context
- `/Users/zeuslegacy/Documents/onemind-codex/75-99 GE (Generational Empire)/README.md` — domain overview

### Your Domain

**GE — Generational Empire (75-99): YOUR EMPIRE**
Path: `/Users/zeuslegacy/Documents/onemind-codex/75-99 GE (Generational Empire)/`

Current folders:
```
75 Brand (Personal & Corporate)/                ← 75 MB
75 Zeus Delacruz Enterprise (Personal Brand)/   ← 13 GB — DUPLICATE with above, needs merge
76 Brands (Portfolio & Entities)/               ← 956 KB
76 Brands (Ventures & Entities)/                ← 1.6 GB — DUPLICATE with above, needs merge
77 Products (Offerings & IP)/                   ← 3.4 GB
78 Marketing (Content & Growth)/                ← 3.4 MB
78 Marketing (Growth & Outreach)/               ← 265 MB — DUPLICATE with above, needs merge
79 Technology (Tech Stack)/                     ← 768 KB
80 Treasury (Finance)/                          ← 1.1 GB
81 Strategy (Analytics)/                        ← 3.2 MB
82 Sales (Revenue & Clients)/                   ← 60 KB
83 Partnerships (Network)/                      ← 19 MB
83 Ventures (Startups & Investments)/           ← 5.4 MB — DUPLICATE with above, needs merge
84 Capital (Fundraising & Finance)/             ← 21 MB
84 Ventures (New Ops & Innovation)/             ← 15 MB — DUPLICATE with above
99 Archive/                                     ← empty
_from-research/                                 ← 3 MB files routed from SO/Research
_legacy/                                        ← 1.6 GB old numbering content
_sorted-duplicates-docs/                        ← 46 MB — likely safe to delete
```

Loose files at domain root:
- `Revenue-Analysis-SecondBrain-vs-OneMind.md` — route to 81 Strategy or 77 Products

### Known Issues to Fix — THIS IS THE MESSIEST DOMAIN
1. **5 DUPLICATE FOLDER PAIRS** — this is the biggest job:
   - `75 Brand (Personal & Corporate)` vs `75 Zeus Delacruz Enterprise (Personal Brand)` (13 GB!)
   - `76 Brands (Portfolio & Entities)` vs `76 Brands (Ventures & Entities)`
   - `78 Marketing (Content & Growth)` vs `78 Marketing (Growth & Outreach)`
   - `83 Partnerships (Network)` vs `83 Ventures (Startups & Investments)`
   - `84 Capital (Fundraising & Finance)` vs `84 Ventures (New Ops & Innovation)`
2. **_legacy/ (1.6 GB)** — old numbering system content, needs triage
3. **_sorted-duplicates-docs/ (46 MB)** — needs verification, likely all duplicates → _to-delete
4. **99 Archive/** — empty, can go to _to-delete
5. **Loose .md file** at domain root — route to proper folder

### How to Handle Duplicate Folder Merges
For each duplicate pair:
1. **List contents of BOTH folders** — understand what's in each
2. **Identify the KEEPER** — usually the one with more content or better name
3. **Present the merge plan to Zeus** — which folder survives, what moves where
4. **After approval**: Move unique files from smaller → bigger, move the emptied folder to _to-delete
5. **Rename the keeper** if needed to have the canonical name

**For domain 75 specifically**: The `75 Zeus Delacruz Enterprise (Personal Brand)/` is 13 GB — that's likely the real content. `75 Brand (Personal & Corporate)/` at 75 MB is likely the shell. Present both to Zeus before deciding.

---

## THE PROCESS (follow this exactly)

### Phase 1: Audit Each Domain Folder (one at a time)

For each numbered domain folder (75-84, 99), do this:

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

### Phase 3: Merge Duplicate Folders

Work through the 5 duplicate pairs. For each:
1. Audit both folders
2. Propose keeper + merge plan
3. Wait for Zeus approval
4. Execute merge
5. Move emptied folder to _to-delete

---

## RULES (critical — follow these exactly)

### Naming Convention
- **Files**: `lowercase-with-hyphens.md` (e.g., `brand-strategy.md`, `revenue-forecast.md`)
- **Domain folders**: `XX Name (Subtitle)/` (e.g., `75 Brand (Personal & Corporate)/`)
- **No spaces in filenames**, no Title-Case, no underscores

### _to-delete Pattern
- **NEVER delete files directly** — always move to `_to-delete/` within the domain group
- Create `_to-delete/` at `/Users/zeuslegacy/Documents/onemind-codex/75-99 GE (Generational Empire)/_to-delete/` if it doesn't exist
- Organize by source: `_to-delete/from-75-brand/`, `_to-delete/from-76-brands/`, etc.

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
| `#pillar/GE` | `#domain/ge` |
| `Your Family` (for LE) | `Your Foundation` |
| `Your Business` (for GE tagline) | `Your Empire` |
| `enterprise` (for GE concept) | `endeavor` |
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
tags: [domain/ge, type/note]
ai_summary: "One-line summary (≤20 words)"
---
```

### Do NOT
- Delete any files without Zeus's explicit approval
- Rename domain group folders (the `75-99 GE (Generational Empire)` name is locked)
- Create new top-level folders
- Move files out of _to-delete once placed there
- Change file content beyond terminology fixes and frontmatter

---

## WORKFLOW

Start with: **"I'm ready to audit GE. Which domain folder should I start with?"**

Then work through each folder one at a time. Present findings, wait for approval, execute. Repeat.

Suggested order:
1. **Quick wins first**: 79, 81, 82, 99 (smallest, cleanest)
2. **Single folders**: 77, 80 (big but no merge needed)
3. **Duplicate merges**: 76, 78, 83, 84 (merge pairs)
4. **The big one last**: 75 (13 GB merge)
5. **Special folders**: _legacy, _sorted-duplicates-docs, _from-research

**GO.**


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
