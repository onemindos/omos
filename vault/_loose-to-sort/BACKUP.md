# OM OFFSITE — Backup Purpose & Retention

> **DO NOT DELETE** until a second backup of binaries exists (Obsidian Sync, Dropbox, or Time Machine).

---

## What This Is

This drive contains a **pre-filter-repo snapshot** of the OneMind Codex vault, taken on **2026-04-02** before `.git` history was cleaned.

| Stat | Value |
|------|-------|
| **Total size** | ~126 GB |
| **Total files** | ~152,000 |
| **Domains** | All 4 (SO, HP, LE, GE) verified present |
| **Snapshot date** | 2026-04-02 |
| **Source path** | `/Users/zeuslegacy/Documents/onemind-codex/` |

---

## Why This Backup Exists

On 2026-04-02, the onemind-codex vault underwent a major cleanup:

1. **git-filter-repo** stripped all binary files (PDFs, images, videos, etc.) from git history
2. `.git` went from **48 GB → 273 MB** (46.2 GB was dead LFS cache, rest was binary blobs in history)
3. The cleanup was **irreversible** — `git push --force` overwrote GitHub history (3,099 commits rewritten)
4. This drive is the **only copy** of the pre-cleanup `.git` state

### What's on this drive that exists nowhere else

- The original `.git/` directory with full binary history (pre-filter-repo)
- Old LFS objects cache (44 GB of `.git/lfs/objects/`)
- Any binary files that may have been in git history but not in the working tree

### What's on this drive AND on the Mac

- All markdown files (also on GitHub)
- All PDFs in `_docs/` (~46.7 GB)
- All images in `_images/` (~29.4 GB)
- All domain folder content

---

## Current Backup Situation (April 2026)

| Content | Mac | GitHub | OM OFFSITE | Other Backup |
|---------|-----|--------|------------|--------------|
| Markdown + code | ✅ | ✅ | ✅ | None |
| PDFs (`_docs/`) | ✅ | ❌ (gitignored) | ✅ | **None** |
| Images (`_images/`) | ✅ | ❌ (gitignored) | **None** |
| Old git history | ❌ (rewritten) | ❌ (force-pushed) | ✅ | **None** |

**OM OFFSITE is the only backup of ~78 GB of binary files (PDFs + images).**

If the Mac drive fails before another backup exists, this drive is the recovery path.

---

## When It's Safe to Delete

Delete the `onemind-codex/` folder from this drive when **ANY** of these are true:

1. **Obsidian Sync is active** — vault (including binaries) mirrors to Obsidian Cloud
2. **Dropbox mirror is running** — another cloud copy of all files
3. **Time Machine is backing up** the vault to another drive
4. **A fresh snapshot** has been taken to another external drive
5. **You've confirmed** (after weeks of normal use) nothing was lost in the filter-repo cleanup

**The safest trigger:** One additional backup of the binaries beyond the Mac — whether that's Obsidian Sync, Dropbox, Time Machine, or a new drive snapshot.

---

## What Happened (Timeline)

| Date | Event |
|------|-------|
| 2026-04-02 | Full vault snapshot copied to OM OFFSITE |
| 2026-04-02 | Verified: 126 GB, 152K files, all 4 domains present |
| 2026-04-02 | Deleted `.git/lfs/objects/` (44 GB dead cache) from Mac |
| 2026-04-02 | Ran `git-filter-repo` — stripped ~60 binary extensions from all 3,099 commits |
| 2026-04-02 | `.git` on Mac: 48 GB → 273 MB |
| 2026-04-02 | `git push --force origin main` — 237 MB uploaded, GitHub history rewritten |
| 2026-04-02 | `.gitattributes` cleaned — all LFS rules removed |
| 2026-04-02 | LFS hooks removed from repo |
| 2026-04-05 | Product & distribution strategy added to AGENTS.md / CLAUDE.md across 4 repos |

---

## Related Drives

| Drive | Mount Point | Role |
|-------|-------------|------|
| **OM PRIMARY** | `/Volumes/OM PRIMARY/onemind-media/` | Master organized media library (~3.2 TB) |
| **OM LOCAL** | `/Volumes/OM LOCAL/onemind-media/` | Working drive — staging, rescue |
| **OM OFFSITE** | `/Volumes/OM OFFSITE/` | **This drive** — codex vault backup |

---

*Created: 2026-04-05 by Zeus + Legacy (Copilot agent session)*
