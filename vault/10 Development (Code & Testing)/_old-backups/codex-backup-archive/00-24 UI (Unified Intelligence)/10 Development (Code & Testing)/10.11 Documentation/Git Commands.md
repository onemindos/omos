---
id: ui-12-20260111053926
title: "Git Commands"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/12]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Git Commands"
---

**Core Git Commands You Need:**

**Repository Management:**

- `git clone` - Copy a repository
- `git init` - Start a new repo
- `git remote` - Manage remote connections
- `git fetch` - Download without merging
- `git pull` - Fetch + merge
- `git push` - Send your commits

**Daily Workflow:**

- `git status` - See what's changed
- `git add` - Stage files
- `git commit` - Save your changes
- `git log` - View history
- `git show` - Inspect commits
- `git stash` - Temporarily shelve changes

**Branching:**

- `git branch` - List/create branches
- `git checkout` / `git switch` - Change branches
- `git merge` - Combine branches
- `git rebase` - Rewrite history cleanly
- `git cherry-pick` - Grab specific commits

**Undo/Fix Things:**

- `git reset` - Undo commits (dangerous but powerful)
- `git revert` - Undo safely with new commit
- `git restore` - Discard changes
- `git clean` - Remove untracked files

**Advanced Tools:**

- `git bisect` - Binary search to find bugs
- `git blame` - See who changed what
- `git reflog` - Safety net for "lost" commits
- `git submodule` - Manage nested repos
- `git worktree` - Multiple working directories

**Security/Secrets:**

- `git-crypt` - You know this one
- `git-secret` - Alternative to git-crypt
- `git filter-branch` / `git filter-repo` - Remove sensitive data from history

**Hooks:**

- `.git/hooks/` - Automate actions on commit, push, etc.
- `pre-commit` framework - Better hook management

**For your infrastructure work:** Focus on `git submodule`, `git worktree`, and hooks. You'll need them for managing multiple repos across your Legacy Platform infrastructure. And get `git filter-repo` installed now—when you need to purge secrets from history, you need it _fast_.