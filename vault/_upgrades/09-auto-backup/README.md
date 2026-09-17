---
id: upgrade-09-backup
title: "Upgrade 09: Auto-Backup System"
type: blueprint
status: deployed
created: 2026-03-27
deployed: 2026-03-27
owner: zeus
tags: [type/blueprint, status/deployed, domain/01-infrastructure]
destination: "_system/backup/"
---

# Upgrade 09: Auto-Backup System

> ✅ **DEPLOYED** to `_system/backup/` — commit `28356c78` (2026-03-27)

> **Source**: PAI's TELOS auto-backup (timestamped snapshots before every change)
> **Destination**: `_system/backup/`
> **Why**: Git is great for version control, but it requires explicit commits. PAI creates timestamped backups *before* modifying critical files, allowing instant rollback without git archaeology. Defense-in-depth for the most important files.

---

## What Gets Backed Up

### Critical Files (always backup before modifying)

| File | Why It's Critical |
|------|-------------------|
| `ONEMIND-CODEX.md` | Canonical master document |
| `_codex/system/*.md` | Always-loaded identity context |
| `_codex/memory/**` | Persistent agent memory |
| `AGENTS.md` | Universal agent instructions |
| `CLAUDE.md` | Claude-specific instructions |
| `CODEX-FRAMEWORK.md` | Method definition |

### Backup Triggers

1. **Before any agent modifies a critical file** — Automatic snapshot
2. **Before bulk operations** (renaming, refactoring) — Pre-operation checkpoint
3. **Before deploying to VPS** — Pre-deploy snapshot
4. **Manual trigger** — Zeus can request a backup anytime

## Backup Format

```
_system/backup/
├── README.md
├── 2026-03-27T15-30-00/
│   ├── manifest.json        ← What was backed up and why
│   ├── ONEMIND-CODEX.md     ← Snapshot
│   ├── system/
│   │   ├── identity.md
│   │   ├── mission.md
│   │   └── active-goals.md
│   └── memory/
│       └── learning/
│           └── signals.jsonl
├── 2026-03-28T09-00-00/
│   └── ...
└── .gitignore              ← Exclude from git (backups are local-only)
```

## Manifest Format

```json
{
  "timestamp": "2026-03-27T15:30:00Z",
  "trigger": "pre-modification",
  "reason": "About to update ONEMIND-CODEX.md during brand repositioning",
  "agent": "copilot",
  "files": [
    "ONEMIND-CODEX.md",
    "_codex/system/identity.md",
    "_codex/system/active-goals.md"
  ]
}
```

## Retention Policy

- **Keep last 10 backups** — Auto-prune oldest when exceeding
- **Keep all backups < 7 days old** — Never prune recent
- **Git-ignored** — Backups stay local, don't bloat the repo
- **Exception**: Monthly backups are kept permanently (1st of each month)

## How to Deploy

1. Create `_system/backup/` with README and `.gitignore`
2. Add backup instructions to relevant hooks (execute-pre)
3. Agents create snapshot before modifying critical files
4. Prune automatically based on retention policy

## PAI Equivalent

PAI's TELOS system creates timestamped backups in `~/.claude/PAI/USER/TELOS/backups/` before any Update workflow modifies personal context files. Each backup includes the full file + a changelog entry in `updates.md`.
