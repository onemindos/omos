# Codex Skill: Defragment

How to restructure and reorganize the vault when it becomes cluttered over long-horizon use. This is distinct from consolidation (which merges related notes) — defragmentation restructures the *entire hierarchy* of a domain or the whole vault.

> Inspired by Letta Context Repositories' memory defragmentation skill: "backs up the agent's memory filesystem, then launches a subagent that reorganizes files, splitting large files, merging duplicates, and restructuring into a clean hierarchy of 15-25 focused files."

---

## When to Defragment

- A domain folder has 50+ files and feels unwieldy
- Multiple domains have overlapping content that should be restructured
- The vault has grown organically and the hierarchy no longer reflects reality
- After a major life/project change that shifts priorities
- Quarterly (or when `_codex/skills/maintain.md` audit scores poorly)

## What Defrag Does That Consolidate Doesn't

| Consolidate | Defragment |
|---|---|
| Merges scattered notes on same topic | Restructures the entire folder hierarchy |
| Operates on individual notes | Operates on whole domains or vault |
| Preserves existing structure | May create/remove subfolders |
| Bottom-up (note-by-note) | Top-down (architecture-first) |
| Safe to run anytime | Requires backup first |

## The Defragmentation Process

### Step 0: Backup

**Always backup before defragmenting.**

```bash
git add -A && git commit -m "pre-defrag: snapshot before restructure"
git tag "pre-defrag-$(date +%Y%m%d)"
```

If something goes wrong, you can restore:

```bash
git reset --hard pre-defrag-YYYYMMDD
```

### Step 1: Audit the Target

Choose what to defragment:
- A single domain folder (e.g., `08 Infrastructure/`)
- A domain group (e.g., all of `00-24 SO/`)
- The entire vault

For the target, catalog:
- Total file count
- Files by type (note, decision, SOP, etc.)
- Files by status (active, draft, archived, deprecated)
- Average file size (flag anything over 1000 lines)
- Files with no incoming links (orphans)
- Files with stale `modified:` dates (>6 months)

### Step 2: Design the New Structure

Based on the audit, decide what the domain should look like after defrag:

**Target: 15-25 focused files per domain** (Letta's recommendation). Achieve this by:

- **Split** oversized files (>300 lines) into focused subtopics
- **Merge** overlapping files into single authority documents
- **Archive** deprecated, stale, or superseded content
- **Create subfolders** if a domain naturally has 3+ distinct subtopics with 5+ files each
- **Remove empty subfolders** that add hierarchy without purpose

### Step 3: Execute the Restructure

For each change:

1. Move/rename files to their new locations
2. Update frontmatter (`modified:`, `tags:`, `related:`)
3. Update all `[[wiki-links]]` that point to moved files
4. Update the domain README.md to reflect the new structure

**Order of operations** (prevents broken links):
1. Merges first (combine files, update links)
2. Splits second (break out subtopics from large files)
3. Moves third (relocate files to new subfolders or domains)
4. Archives last (remove deprecated content)

### Step 4: Update `system/` Directory

After restructuring, re-evaluate what's in `_codex/system/`:
- Are the pinned files still the most essential?
- Did any new files emerge that should be pinned?
- Did any pinned files become less relevant?

### Step 5: Verify and Commit

1. Run `_codex/skills/maintain.md` health audit on the restructured areas
2. Ensure no broken links
3. Ensure all files have valid frontmatter
4. Commit with descriptive message:

```bash
git add -A && git commit -m "defrag: restructured [target] — [N] files → [M] files, [summary of changes]"
```

---

## Defragmentation Patterns

### Pattern: Flat Folder → Organized Subfolders

```
BEFORE:
08 Infrastructure/
├── vps-setup.md
├── nats-config.md
├── nats-streams.md
├── nats-subjects.md
├── docker-compose.md
├── livekit-setup.md
├── tailscale-config.md
├── firewall-rules.md
├── ssl-certs.md
├── port-reference.md
├── backup-strategy.md
├── monitoring.md
└── ... (30+ files)

AFTER:
08 Infrastructure/
├── README.md
├── vps/
│   ├── setup.md
│   ├── firewall.md
│   ├── ssl.md
│   └── monitoring.md
├── nats/
│   ├── config.md (merged from config + streams + subjects)
│   └── bridge.md
├── services/
│   ├── docker-compose.md
│   ├── livekit.md
│   └── tailscale.md
└── ops/
    ├── backup.md
    └── port-reference.md
```

### Pattern: Oversized File → Focused Files

```
BEFORE:
10 Development/project-notes.md (800 lines covering frontend, backend, deployment, testing)

AFTER:
10 Development/
├── frontend-architecture.md (~200 lines)
├── backend-api-design.md (~200 lines)
├── deployment-pipeline.md (~150 lines)
└── testing-strategy.md (~150 lines)
```

### Pattern: Cross-Domain Cleanup

```
BEFORE:
02 Agents/legacy-deployment.md (really about infrastructure)
08 Infrastructure/agent-monitoring.md (really about agents)
12 Research/vps-performance.md (mixed infra + research)

AFTER:
02 Agents/monitoring.md (agent-specific monitoring)
08 Infrastructure/deployment.md (all deployment procedures)
12 Research/infrastructure-performance-analysis.md (research-only)
```

---

## Safety Rules

1. **Always tag the commit before defrag** — you must be able to revert
2. **Don't defrag and add new content simultaneously** — restructure is its own operation
3. **Preserve git history** — use `git mv` when possible so history follows the file
4. **Update all links** — broken links after defrag means the defrag was incomplete
5. **One domain at a time** — don't defrag the whole vault in one pass unless it's small
