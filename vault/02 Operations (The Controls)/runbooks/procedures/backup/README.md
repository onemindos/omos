# Auto-Backup System

> Timestamped snapshots of critical files before modification. Defense-in-depth beyond git.

## How It Works

1. Before modifying a critical file, agent creates a timestamped backup here
2. Backup includes a `manifest.json` describing what was backed up and why
3. Snapshots are git-ignored (local only, don't bloat the repo)
4. Retention: keep last 10 + all < 7 days old + monthly permanents

## Critical Files (Always Backup Before Modifying)

- `ONEMIND-CODEX.md`
- `AGENTS.md`
- `CLAUDE.md`
- `CODEX-FRAMEWORK.md`
- `_codex/system/*.md`
- `_codex/memory/**`

## Directory Structure

```
_system/backup/
├── README.md          ← This file
├── .gitignore         ← Excludes snapshots from git
└── YYYY-MM-DDThh-mm-ss/
    ├── manifest.json  ← What, why, who
    └── {files}        ← Copied files
```

## Manifest Format

```json
{
  "timestamp": "2026-03-27T15:30:00Z",
  "trigger": "pre-modification",
  "reason": "About to update identity files during upgrade deployment",
  "agent": "copilot",
  "files": ["_codex/system/identity.md", "_codex/system/active-goals.md"]
}
```

## Agent Instructions

Before modifying any file in the critical list:
1. Create directory: `_system/backup/{ISO-timestamp}/`
2. Copy the target file(s) into it
3. Write `manifest.json`
4. Proceed with modification
