---
title: "CLI — git"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/git]
ai_summary: "Git CLI quick reference for the One Mind stack — commit, push, branch, remotes, and repo map."
---

# git — Version Control CLI

> Tracks all markdown + code. Deploy trigger: push to GitHub → Cloudflare Pages builds codex-live.

---

## Repos

| Repo | Local Path | Remote |
|------|-----------|--------|
| onemind-codex (private vault) | `~/Documents/onemind-codex` | `github.com/OneMind-OS/onemind-codex` |
| onemind-codex-live (public site) | `~/Documents/onemind-codex-live` | `github.com/OneMind-OS/onemind-codex-live` |
| onemind-web (platform) | `~/Documents/onemind-web` | `github.com/One-Mind-OS/onemind-web` |
| onemind-codex-template | `~/Documents/onemind-codex-template` | `github.com/OneMind-OS/onemind-codex-template` |
| Personal_AI_Infrastructure (PAI) | `~/Documents/Personal_AI_Infrastructure` | `github.com/danielmiessler/fabric` (upstream) |

---

## Daily Commands

```bash
# Status
git status
git log --oneline -10

# Commit and push
git add -A
git commit -m "so: update agent skills"
git push

# Pull latest
git pull origin main

# Check remotes
git remote -v
```

## Branch Commands

```bash
git checkout -b feature/branch-name
git checkout main
git merge feature/branch-name
git branch -d feature/branch-name   # delete after merge
```

## Useful Utilities

```bash
# See what changed
git diff
git diff --staged

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Discard unstaged changes to a file
git checkout -- filename.md

# See commit history with files changed
git log --oneline --stat

# Blame (who changed what line)
git blame filename.md
```

## .gitignore (codex vault)

```
_docs/
_images/
_library/
.obsidian/
*.DS_Store
```

Binary assets (PDFs, images) sync via Obsidian Sync — git tracks markdown + code only.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
