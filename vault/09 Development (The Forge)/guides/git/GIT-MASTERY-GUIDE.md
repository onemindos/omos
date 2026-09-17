---
id: ui-10-20260123
title: "Git Mastery Guide"
type: reference
path: UI
domain: 00-24
folder: "10"
status: active
created: 2026-01-23
modified: 2026-01-23
owner: zeus
tags: [path/ui, path/ui/10, type/reference, status/active, git, development]
ai_summary: "Comprehensive Git reference covering branches, merging, rebasing, secrets, LFS, and configuration"
---

# Git Mastery Guide

> **"Version control is the foundation of professional software development."**

This guide covers everything you need to master Git for One Mind development.

---

## Table of Contents

### Fundamentals

1. [Core Concepts](#core-concepts)
2. [Branching](#branching)
3. [Commits](#commits)
4. [Merging](#merging)
5. [Rebasing](#rebasing)
6. [Push & Pull](#push--pull)

### Configuration

7. [.gitignore](#gitignore)
8. [.gitkeep](#gitkeep)
9. [.gitattributes](#gitattributes)
10. [Git Config](#git-config)
11. [Git Aliases](#git-aliases)

### Security & Storage

11. [Git Secrets](#git-secrets)
12. [Git LFS](#git-lfs-large-file-storage)
13. [Signed Commits](#signed-commits-gpg)

### Advanced Features

14. [Git Stash](#git-stash)
15. [Git Tags](#git-tags)
16. [Git Cherry-Pick](#git-cherry-pick)
17. [Git Bisect](#git-bisect)
18. [Git Hooks](#git-hooks)
19. [Git Submodules](#git-submodules)
20. [Git Worktrees](#git-worktrees)
21. [Git Reflog](#git-reflog)

### Workflows & Best Practices

22. [Remote Management](#remote-management)
23. [Forking Workflow](#forking-workflow)
24. [Branching Strategies](#branching-strategies)
25. [One Mind Git Workflow](#one-mind-git-workflow)
26. [Troubleshooting](#troubleshooting)
27. [Quick Reference Card](#quick-reference-card)

---

## Core Concepts

### What is Git?

Git is a **distributed version control system** that tracks changes in your code over time. Every developer has a complete copy of the repository history.

```
┌─────────────────────────────────────────────────────────┐
│                    GIT ARCHITECTURE                      │
├─────────────────────────────────────────────────────────┤
│                                                          │
│   Working Directory → Staging Area → Local Repo → Remote │
│        (edit)          (git add)    (git commit) (push)  │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### The Three States

| State         | Location             | Description                  |
| ------------- | -------------------- | ---------------------------- |
| **Modified**  | Working Directory    | Files changed but not staged |
| **Staged**    | Staging Area (Index) | Files marked for next commit |
| **Committed** | Local Repository     | Changes saved to Git history |

---

## Branching

### What is a Branch?

A **branch** is an independent line of development. It's a pointer to a specific commit that moves forward as you make new commits.

```
main:     A ─── B ─── C ─── D ─── E
                      │
feature:              └─── F ─── G ─── H
```

### Branch Commands

```bash
# List all branches
git branch              # Local branches
git branch -r           # Remote branches
git branch -a           # All branches

# Create a new branch
git branch feature-name

# Create and switch to new branch
git checkout -b feature-name
# OR (modern Git)
git switch -c feature-name

# Switch to existing branch
git checkout main
git switch main

# Rename a branch
git branch -m old-name new-name

# Delete a branch
git branch -d feature-name      # Safe delete (only if merged)
git branch -D feature-name      # Force delete

# Delete remote branch
git push origin --delete feature-name
```

### Branch Naming Conventions (One Mind)

```
feature/   → New features           feature/voice-interface
bugfix/    → Bug fixes              bugfix/memory-leak
hotfix/    → Urgent production fix  hotfix/auth-crash
refactor/  → Code refactoring       refactor/agent-tools
docs/      → Documentation only     docs/api-reference
test/      → Test additions         test/linear-integration
chore/     → Maintenance tasks      chore/update-deps
```

---

## Commits

### What is a Commit?

A **commit** is a snapshot of your repository at a specific point in time. Each commit has:

- **SHA hash** — Unique identifier (e.g., `a1b2c3d`)
- **Author** — Who made the change
- **Timestamp** — When it was made
- **Message** — Description of changes
- **Parent(s)** — Previous commit(s)

### Commit Commands

```bash
# Stage files for commit
git add filename.py           # Single file
git add .                     # All changed files
git add -p                    # Interactive staging (patch mode)

# Commit staged changes
git commit -m "Add voice interface"

# Commit with detailed message
git commit                    # Opens editor for multi-line message

# Amend last commit (before pushing)
git commit --amend -m "Better message"
git commit --amend --no-edit  # Keep message, add staged files

# View commit history
git log                       # Full history
git log --oneline             # Compact view
git log --graph --oneline     # Visual branch graph
git log -5                    # Last 5 commits

# View specific commit
git show a1b2c3d

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1
```

### Commit Message Format (One Mind)

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

| Type       | Description                 |
| ---------- | --------------------------- |
| `feat`     | New feature                 |
| `fix`      | Bug fix                     |
| `docs`     | Documentation               |
| `style`    | Formatting (no code change) |
| `refactor` | Code restructuring          |
| `test`     | Adding tests                |
| `chore`    | Maintenance                 |
| `perf`     | Performance improvement     |

**Examples:**

```bash
# Simple
git commit -m "feat(agent): add Linear integration tool"

# With body
git commit -m "fix(memory): resolve duplicate entries

The memory service was creating duplicates when
the same context was saved within 100ms.

Fixes #123"
```

---

## Merging

### What is Merging?

**Merging** combines changes from one branch into another. It creates a new "merge commit" that has two parents.

```
Before merge:
main:     A ─── B ─── C
                │
feature:        └─── D ─── E

After merge:
main:     A ─── B ─── C ─────── M (merge commit)
                │               │
feature:        └─── D ─── E ───┘
```

### Merge Commands

```bash
# Switch to target branch first
git checkout main

# Merge feature branch into main
git merge feature-branch

# Merge with commit message
git merge feature-branch -m "Merge: Add voice interface"

# Abort a conflicted merge
git merge --abort

# View merge conflicts
git status
git diff
```

### Merge Strategies

| Strategy            | Command               | When to Use                   |
| ------------------- | --------------------- | ----------------------------- |
| **Fast-forward**    | Default when possible | Linear history, no divergence |
| **No fast-forward** | `git merge --no-ff`   | Always create merge commit    |
| **Squash**          | `git merge --squash`  | Combine all commits into one  |

```bash
# Fast-forward (default)
git merge feature-branch

# Force merge commit (recommended for features)
git merge --no-ff feature-branch -m "Merge: Feature description"

# Squash merge (clean history)
git merge --squash feature-branch
git commit -m "feat: Add complete voice interface"
```

### Resolving Merge Conflicts

When Git can't auto-merge, you get conflict markers:

```python
<<<<<<< HEAD
def process_memory(self, data):
    return self.store(data)
=======
def process_memory(self, data, context):
    return self.store(data, context)
>>>>>>> feature-branch
```

**Resolution steps:**

```bash
# 1. See conflicted files
git status

# 2. Open and edit files to resolve conflicts
# Remove markers, keep correct code

# 3. Stage resolved files
git add resolved-file.py

# 4. Complete the merge
git commit -m "Merge: Resolve conflicts in memory module"
```

---

## Rebasing

### What is Rebasing?

**Rebasing** moves your branch to start from a different commit. It "replays" your commits on top of another branch, creating a linear history.

```
Before rebase:
main:     A ─── B ─── C ─── D
                │
feature:        └─── E ─── F

After rebase:
main:     A ─── B ─── C ─── D
                            │
feature:                    └─── E' ─── F'
```

### Rebase Commands

```bash
# Rebase current branch onto main
git checkout feature-branch
git rebase main

# Interactive rebase (edit, squash, reorder commits)
git rebase -i HEAD~3        # Last 3 commits
git rebase -i main          # All commits since branching from main

# Continue after resolving conflicts
git rebase --continue

# Abort rebase
git rebase --abort

# Skip a problematic commit
git rebase --skip
```

### Interactive Rebase Options

When you run `git rebase -i`, you get an editor:

```
pick a1b2c3d Add initial voice module
pick b2c3d4e Add speech recognition
pick c3d4e5f Fix typo in voice module
pick d4e5f6g Add text-to-speech

# Commands:
# p, pick   = use commit
# r, reword = use commit, but edit message
# e, edit   = use commit, but stop for amending
# s, squash = meld into previous commit
# f, fixup  = like squash, but discard message
# d, drop   = remove commit
```

**Example: Squash commits**

```
pick a1b2c3d Add initial voice module
squash b2c3d4e Add speech recognition
fixup c3d4e5f Fix typo in voice module
squash d4e5f6g Add text-to-speech
```

Result: One clean commit instead of four.

### Merge vs Rebase

| Aspect               | Merge                 | Rebase                          |
| -------------------- | --------------------- | ------------------------------- |
| **History**          | Preserves all history | Creates linear history          |
| **Commits**          | Adds merge commit     | Rewrites commit hashes          |
| **Conflicts**        | Resolve once          | May resolve multiple times      |
| **Shared branches**  | ✅ Safe               | ⚠️ Dangerous (rewrites history) |
| **Feature branches** | ✅ Good               | ✅ Great (before merging)       |

**Golden Rule:** Never rebase commits that have been pushed to a shared branch.

---

## Push & Pull

### Push (Upload Changes)

```bash
# Push current branch to remote
git push

# Push to specific remote and branch
git push origin main

# Push and set upstream (first time)
git push -u origin feature-branch

# Force push (dangerous - rewrites remote history)
git push --force
git push --force-with-lease  # Safer: fails if remote has new commits

# Push all branches
git push --all origin

# Push tags
git push --tags
```

### Pull (Download Changes)

```bash
# Pull = Fetch + Merge
git pull

# Pull from specific remote/branch
git pull origin main

# Pull with rebase instead of merge
git pull --rebase

# Fetch without merging (see changes first)
git fetch origin
git log origin/main..HEAD    # See difference
git merge origin/main        # Then merge
```

### Fetch vs Pull

| Command     | Action                                   |
| ----------- | ---------------------------------------- |
| `git fetch` | Downloads changes but doesn't apply them |
| `git pull`  | Downloads AND merges changes             |

```bash
# Safe workflow
git fetch origin              # Get updates
git log origin/main..HEAD     # See what's different
git merge origin/main         # Apply when ready
```

---

## .gitignore

### What is .gitignore?

A `.gitignore` file tells Git which files/folders to ignore. These won't be tracked or committed.

### .gitignore Syntax

```gitignore
# Comments start with #

# Ignore specific file
secrets.json

# Ignore by extension
*.log
*.pyc
*.env

# Ignore directory
node_modules/
__pycache__/
.venv/

# Ignore in any subdirectory
**/*.log

# Ignore in root only
/build/

# Negate (don't ignore)
!important.log

# Ignore pattern
temp-*
*.local.*
```

### One Mind .gitignore Template

```gitignore
# ===== One Mind OS .gitignore =====

# ----- Environment & Secrets -----
.env
.env.*
!.env.example
*.local
secrets/
.secrets/

# ----- Python -----
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
venv/
.venv/
env/
.eggs/
*.egg-info/
.pytest_cache/
.mypy_cache/
.ruff_cache/

# ----- Node.js -----
node_modules/
npm-debug.log*
yarn-error.log*
.pnpm-store/

# ----- IDE & Editors -----
.vscode/settings.json
.idea/
*.swp
*.swo
*~
.DS_Store

# ----- Build Outputs -----
dist/
build/
out/
*.build/

# ----- Logs -----
logs/
*.log

# ----- Database -----
*.db
*.sqlite
*.sqlite3

# ----- Testing -----
coverage/
htmlcov/
.coverage

# ----- Obsidian -----
.obsidian/workspace.json
.obsidian/workspace-mobile.json
.obsidian/plugins/*/data.json
.trash/

# ----- OS Files -----
.DS_Store
Thumbs.db
Desktop.ini
```

### Global .gitignore

Set up a global ignore file for all repositories:

```bash
# Create global gitignore
touch ~/.gitignore_global

# Configure Git to use it
git config --global core.excludesfile ~/.gitignore_global
```

---

## .gitkeep

### What is .gitkeep?

**Git doesn't track empty directories.** The `.gitkeep` file is a convention to preserve empty directory structure in your repository.

```text
# Without .gitkeep - directory not tracked
logs/
  (empty)

# With .gitkeep - directory tracked
logs/
  .gitkeep
```

**Note:** `.gitkeep` is **not a Git feature** — it's just a community convention. You could name it anything (`.keep`, `README.md`, etc.), but `.gitkeep` is the standard.

### Why Use .gitkeep?

Preserve directory structure for:

| Directory   | Purpose               |
| ----------- | --------------------- |
| `logs/`     | Application logs      |
| `uploads/`  | User file uploads     |
| `cache/`    | Temporary cache files |
| `tmp/`      | Temporary files       |
| `data/raw/` | Raw data inputs       |
| `dist/`     | Build outputs         |
| `coverage/` | Test coverage reports |

### Creating .gitkeep Files

```bash
# Single directory
mkdir -p logs
touch logs/.gitkeep
git add logs/.gitkeep

# Multiple directories
mkdir -p logs/{error,access,debug}
touch logs/{error,access,debug}/.gitkeep
git add logs/

# Create entire structure
mkdir -p data/{raw,processed,results}
touch data/{raw,processed,results}/.gitkeep
git add data/
```

### One Mind Directory Structure

```bash
# Create standard directories for One Mind projects

# Backend
mkdir -p logs/{agent,api,worker}
touch logs/{agent,api,worker}/.gitkeep

mkdir -p data/{raw,processed,embeddings}
touch data/{raw,processed,embeddings}/.gitkeep

mkdir -p tmp/{uploads,exports}
touch tmp/{uploads,exports}/.gitkeep

# Frontend
mkdir -p public/{uploads,exports}
touch public/{uploads,exports}/.gitkeep

mkdir -p src/assets/{images,fonts,videos}
touch src/assets/{images,fonts,videos}/.gitkeep

# Add all
git add logs/ data/ tmp/ public/ src/assets/
git commit -m "chore: Add directory structure with .gitkeep"
```

### .gitkeep Content

Most `.gitkeep` files are empty, but you can add helpful comments:

```bash
# logs/.gitkeep
cat > logs/.gitkeep << 'EOF'
# Application logs directory
# Logs are gitignored but directory structure is preserved
# See .gitignore for excluded patterns
EOF
```

### .gitkeep vs README.md

| File            | Purpose                    | Best For                           |
| --------------- | -------------------------- | ---------------------------------- |
| **`.gitkeep`**  | Preserve empty directory   | Temporary/output directories       |
| **`README.md`** | Document directory purpose | Directories with important context |

```bash
# Use .gitkeep for output directories
logs/.gitkeep
tmp/.gitkeep

# Use README.md for important directories
docs/README.md
src/components/README.md
```

### Common Pattern with .gitignore

```bash
# Directory structure
logs/
  .gitkeep
  error.log
  access.log
```

**.gitignore:**

```gitignore
# Ignore all logs
logs/*.log
logs/**/*.log

# But keep the directory structure
!logs/.gitkeep
!logs/*/.gitkeep
```

### One-liner to Add .gitkeep to Empty Directories

```bash
# Find empty directories and add .gitkeep
find . -type d -empty -not -path "*/\.*" -exec touch {}/.gitkeep \;

# Then add to Git
git add .
git commit -m "chore: Add .gitkeep to empty directories"
```

### Alternative: .keep

Some projects use `.keep` instead:

```bash
touch logs/.keep
```

Both work the same — choose one convention and stick with it. **One Mind standard: `.gitkeep`**

### When NOT to Use .gitkeep

❌ **Don't use** for:

- Directories that will always have files
- Generated directories (create via build scripts instead)
- `node_modules/`, `.git/`, etc. (managed by tools)

✅ **Do use** for:

- Required empty directories in fresh clones
- Directory structure that needs to exist before first run
- Organizational structure you want preserved

---

## .gitattributes

### What is .gitattributes?

A `.gitattributes` file controls how Git handles files — line endings, diff behavior, merge strategies, and LFS tracking.

### Common Attributes

```gitattributes
# Set default behavior
* text=auto

# Force LF line endings (Unix)
*.sh text eol=lf
*.py text eol=lf

# Force CRLF line endings (Windows)
*.bat text eol=crlf

# Binary files (don't modify)
*.png binary
*.jpg binary
*.pdf binary
*.zip binary

# Mark as generated (exclude from stats)
*.lock linguist-generated

# Custom diff driver
*.md diff=markdown

# Merge strategy
database.json merge=ours
```

### One Mind .gitattributes Template

```gitattributes
# ===== One Mind OS .gitattributes =====

# ----- Default -----
* text=auto

# ----- Source Code -----
*.py text eol=lf diff=python
*.js text eol=lf
*.ts text eol=lf
*.tsx text eol=lf
*.json text eol=lf
*.yaml text eol=lf
*.yml text eol=lf
*.md text eol=lf diff=markdown
*.sh text eol=lf

# ----- Windows Scripts -----
*.bat text eol=crlf
*.cmd text eol=crlf
*.ps1 text eol=crlf

# ----- Binary Files -----
*.png binary
*.jpg binary
*.jpeg binary
*.gif binary
*.ico binary
*.pdf binary
*.zip binary
*.tar.gz binary
*.woff binary
*.woff2 binary
*.ttf binary
*.eot binary

# ----- Git LFS (Large Files) -----
*.psd filter=lfs diff=lfs merge=lfs -text
*.ai filter=lfs diff=lfs merge=lfs -text
*.mp4 filter=lfs diff=lfs merge=lfs -text
*.mp3 filter=lfs diff=lfs merge=lfs -text
*.wav filter=lfs diff=lfs merge=lfs -text
*.model filter=lfs diff=lfs merge=lfs -text
*.weights filter=lfs diff=lfs merge=lfs -text
*.onnx filter=lfs diff=lfs merge=lfs -text

# ----- Lock Files (Generated) -----
pnpm-lock.yaml linguist-generated
package-lock.json linguist-generated
poetry.lock linguist-generated

# ----- Merge Strategies -----
# Keep our version for conflicts
*.lock merge=ours
```

---

## Git Secrets

### Why Git Secrets?

Prevents you from committing sensitive data like API keys, passwords, and tokens.

### Installation

```bash
# macOS
brew install git-secrets

# Linux
git clone https://github.com/awslabs/git-secrets.git
cd git-secrets
sudo make install
```

### Setup

```bash
# Initialize in repository
cd your-repo
git secrets --install

# Add AWS patterns (built-in)
git secrets --register-aws

# Add custom patterns
git secrets --add 'TODOIST_API_TOKEN[=:]\s*[A-Za-z0-9]+'
git secrets --add 'LINEAR_API_KEY[=:]\s*[A-Za-z0-9]+'
git secrets --add 'password[=:]\s*.+'
git secrets --add 'secret[=:]\s*.+'

# Add allowed patterns (false positives)
git secrets --add --allowed 'EXAMPLE_KEY=your-key-here'

# Global installation (all repos)
git secrets --install ~/.git-templates/git-secrets
git config --global init.templateDir ~/.git-templates/git-secrets
```

### Usage

```bash
# Scan entire history
git secrets --scan-history

# Scan specific files
git secrets --scan path/to/file.py

# List patterns
git secrets --list
```

### One Mind Secret Patterns

```bash
# Add One Mind specific patterns
git secrets --add 'NATS_URL=nats://[^[:space:]]+'
git secrets --add 'DATABASE_URL=postgres://[^[:space:]]+'
git secrets --add 'HOME_ASSISTANT_TOKEN=[^[:space:]]+'
git secrets --add 'AWS_ACCESS_KEY_ID=[A-Z0-9]{20}'
git secrets --add 'AWS_SECRET_ACCESS_KEY=[A-Za-z0-9/+=]{40}'
git secrets --add 'ANTHROPIC_API_KEY=sk-ant-[A-Za-z0-9-]+'
git secrets --add 'OPENAI_API_KEY=sk-[A-Za-z0-9]+'
```

### Pre-commit Hook Alternative

Using `pre-commit` framework:

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.4.0
    hooks:
      - id: detect-secrets
        args: ["--baseline", ".secrets.baseline"]
```

---

## Git LFS (Large File Storage)

### What is Git LFS?

Git LFS replaces large files with text pointers inside Git, while storing the actual file contents on a remote server. Essential for:

- ML models (`.onnx`, `.pt`, `.h5`)
- Media files (videos, audio)
- Design files (`.psd`, `.ai`, `.fig`)
- Large datasets

### Installation

```bash
# macOS
brew install git-lfs

# Linux
sudo apt install git-lfs

# Initialize (once per user)
git lfs install
```

### Setup in Repository

```bash
# Track file types
git lfs track "*.psd"
git lfs track "*.mp4"
git lfs track "*.model"
git lfs track "*.onnx"

# This creates/updates .gitattributes
cat .gitattributes

# Add .gitattributes to repo
git add .gitattributes
git commit -m "chore: Configure Git LFS"
```

### Usage

```bash
# Check tracked patterns
git lfs track

# List LFS files
git lfs ls-files

# Check LFS status
git lfs status

# Fetch LFS files
git lfs fetch
git lfs pull

# Push LFS files
git lfs push --all origin main

# Migrate existing files to LFS
git lfs migrate import --include="*.psd" --everything
```

### One Mind LFS Configuration

```gitattributes
# ----- AI/ML Models -----
*.onnx filter=lfs diff=lfs merge=lfs -text
*.pt filter=lfs diff=lfs merge=lfs -text
*.pth filter=lfs diff=lfs merge=lfs -text
*.h5 filter=lfs diff=lfs merge=lfs -text
*.pkl filter=lfs diff=lfs merge=lfs -text
*.model filter=lfs diff=lfs merge=lfs -text
*.weights filter=lfs diff=lfs merge=lfs -text
*.safetensors filter=lfs diff=lfs merge=lfs -text

# ----- Media -----
*.mp4 filter=lfs diff=lfs merge=lfs -text
*.mov filter=lfs diff=lfs merge=lfs -text
*.mp3 filter=lfs diff=lfs merge=lfs -text
*.wav filter=lfs diff=lfs merge=lfs -text

# ----- Design -----
*.psd filter=lfs diff=lfs merge=lfs -text
*.ai filter=lfs diff=lfs merge=lfs -text
*.sketch filter=lfs diff=lfs merge=lfs -text
*.fig filter=lfs diff=lfs merge=lfs -text

# ----- Archives -----
*.zip filter=lfs diff=lfs merge=lfs -text
*.tar.gz filter=lfs diff=lfs merge=lfs -text

# ----- Databases -----
*.sqlite filter=lfs diff=lfs merge=lfs -text
```

---

## One Mind Git Workflow

### Branch Strategy

```
main (production)
  │
  ├── develop (integration)
  │     │
  │     ├── feature/voice-interface
  │     ├── feature/memory-v2
  │     └── bugfix/api-timeout
  │
  └── hotfix/critical-auth-fix
```

### Daily Workflow

```bash
# 1. Start your day - get latest
git checkout develop
git pull origin develop

# 2. Create feature branch
git checkout -b feature/my-feature

# 3. Work and commit frequently
git add .
git commit -m "feat(module): Add initial implementation"

# 4. Keep branch updated
git fetch origin
git rebase origin/develop

# 5. Push for backup/review
git push -u origin feature/my-feature

# 6. When ready - create PR
# (Use GitHub/Linear for PR creation)

# 7. After merge - cleanup
git checkout develop
git pull origin develop
git branch -d feature/my-feature
```

### Commit Frequency

| Situation      | Commit Strategy                 |
| -------------- | ------------------------------- |
| New feature    | Commit at each logical step     |
| Bug fix        | One commit per fix              |
| Refactoring    | Commit each refactor separately |
| Before leaving | Always commit WIP               |

### PR Checklist

```markdown
## PR Checklist

- [ ] Branch is rebased on latest `develop`
- [ ] All tests pass
- [ ] No secrets committed (`git secrets --scan`)
- [ ] Large files use LFS
- [ ] Commit messages follow convention
- [ ] Documentation updated if needed
- [ ] Self-reviewed the diff
```

---

## Troubleshooting

### Common Issues

**Undo last commit (keep changes):**

```bash
git reset --soft HEAD~1
```

**Discard all local changes:**

```bash
git checkout -- .
# OR
git restore .
```

**Remove file from staging:**

```bash
git reset HEAD filename
# OR
git restore --staged filename
```

**Fix wrong branch (move commits):**

```bash
# On wrong branch with new commits
git stash
git checkout correct-branch
git stash pop
```

**Recover deleted branch:**

```bash
git reflog
git checkout -b recovered-branch abc1234
```

**Clean untracked files:**

```bash
git clean -n    # Preview
git clean -f    # Delete files
git clean -fd   # Delete files and directories
```

**Squash last N commits:**

```bash
git reset --soft HEAD~3
git commit -m "Combined commit message"
```

**View file at specific commit:**

```bash
git show abc1234:path/to/file.py
```

**Blame (who changed what):**

```bash
git blame filename.py
git blame -L 10,20 filename.py  # Lines 10-20
```

### Emergency Commands

```bash
# OH NO I pushed secrets!
# 1. Remove from history
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch secrets.json' \
  --prune-empty --tag-name-filter cat -- --all

# 2. Force push
git push --force --all

# 3. Rotate the exposed secrets immediately!

# I broke everything
git reflog                    # Find good state
git reset --hard abc1234      # Go back to it
```

---

## Quick Reference Card

```bash
# ===== DAILY COMMANDS =====
git status                    # What's changed?
git add .                     # Stage all
git commit -m "msg"           # Commit
git push                      # Upload
git pull                      # Download + merge

# ===== BRANCHING =====
git branch                    # List branches
git checkout -b name          # Create + switch
git checkout main             # Switch branch
git merge branch              # Merge into current
git branch -d name            # Delete branch

# ===== HISTORY =====
git log --oneline             # Compact history
git log --graph               # Visual tree
git diff                      # See changes
git show abc123               # View commit

# ===== UNDO =====
git reset --soft HEAD~1       # Undo commit, keep changes
git reset --hard HEAD~1       # Undo commit, discard changes
git checkout -- file          # Discard file changes
git stash                     # Temporarily save changes
git stash pop                 # Restore stashed changes

# ===== REMOTE =====
git remote -v                 # List remotes
git fetch                     # Download without merge
git pull --rebase             # Download with rebase
git push -u origin branch     # Push new branch
```

---

## Git Config

### What is Git Config?

Git configuration is stored at three levels:

| Level      | Location         | Command               |
| ---------- | ---------------- | --------------------- |
| **System** | `/etc/gitconfig` | `git config --system` |
| **Global** | `~/.gitconfig`   | `git config --global` |
| **Local**  | `.git/config`    | `git config --local`  |

Local overrides Global, which overrides System.

### Essential Configuration

```bash
# Identity (required)
git config --global user.name "Zeus Delacruz"
git config --global user.email "zeus@onemindcodex.com"

# Default branch name
git config --global init.defaultBranch main

# Default editor
git config --global core.editor "code --wait"  # VS Code
git config --global core.editor "vim"          # Vim
git config --global core.editor "nano"         # Nano

# Line endings
git config --global core.autocrlf input  # macOS/Linux
git config --global core.autocrlf true   # Windows

# Color output
git config --global color.ui auto

# Push behavior
git config --global push.default current
git config --global push.autoSetupRemote true

# Pull behavior
git config --global pull.rebase true  # Rebase instead of merge

# Credential helper (cache passwords)
git config --global credential.helper cache
git config --global credential.helper 'cache --timeout=3600'  # 1 hour

# macOS Keychain
git config --global credential.helper osxkeychain

# Diff and merge tools
git config --global diff.tool vscode
git config --global difftool.vscode.cmd 'code --wait --diff $LOCAL $REMOTE'
git config --global merge.tool vscode
git config --global mergetool.vscode.cmd 'code --wait $MERGED'
```

### View Configuration

```bash
# List all settings
git config --list

# List with origin (where each setting comes from)
git config --list --show-origin

# Get specific value
git config user.name
git config user.email

# Edit config file directly
git config --global --edit
```

### One Mind Git Config

```ini
# ~/.gitconfig for One Mind development

[user]
    name = Zeus Delacruz
    email = zeus@onemindcodex.com
    signingkey = YOUR_GPG_KEY_ID

[init]
    defaultBranch = main

[core]
    editor = code --wait
    autocrlf = input
    excludesfile = ~/.gitignore_global
    pager = delta  # Better diff viewer

[push]
    default = current
    autoSetupRemote = true

[pull]
    rebase = true

[fetch]
    prune = true  # Auto-remove deleted remote branches

[merge]
    conflictstyle = diff3  # Show base version in conflicts

[diff]
    colorMoved = default

[commit]
    gpgsign = true  # Sign all commits

[rebase]
    autoStash = true  # Auto-stash before rebase

[rerere]
    enabled = true  # Remember conflict resolutions

[alias]
    # See Git Aliases section
```

---

## Git Aliases

### What are Git Aliases?

Shortcuts for frequently used commands. Save time and reduce typing errors.

### Essential Aliases

```bash
# Add aliases
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.st status
git config --global alias.unstage 'reset HEAD --'
git config --global alias.last 'log -1 HEAD'
```

### One Mind Alias Collection

```bash
# ===== Basic Shortcuts =====
git config --global alias.s "status -sb"
git config --global alias.co "checkout"
git config --global alias.cob "checkout -b"
git config --global alias.br "branch"
git config --global alias.ci "commit"
git config --global alias.cm "commit -m"

# ===== Logging =====
git config --global alias.lg "log --oneline --graph --decorate -20"
git config --global alias.lga "log --oneline --graph --decorate --all"
git config --global alias.hist "log --pretty=format:'%C(yellow)%h%Creset %ad | %C(green)%s%Creset %C(red)%d%Creset [%an]' --date=short"
git config --global alias.today "log --since=midnight --oneline --graph"
git config --global alias.week "log --since='1 week ago' --oneline --graph"

# ===== Diff & Changes =====
git config --global alias.d "diff"
git config --global alias.dc "diff --cached"
git config --global alias.ds "diff --stat"
git config --global alias.changed "diff --name-only"

# ===== Staging =====
git config --global alias.a "add"
git config --global alias.aa "add --all"
git config --global alias.ap "add --patch"
git config --global alias.unstage "reset HEAD --"

# ===== Commits =====
git config --global alias.amend "commit --amend --no-edit"
git config --global alias.undo "reset --soft HEAD~1"
git config --global alias.wip "commit -am 'WIP: Work in progress'"

# ===== Branches =====
git config --global alias.branches "branch -a"
git config --global alias.gone "!git fetch -p && git branch -vv | grep ': gone]' | awk '{print \$1}'"
git config --global alias.cleanup "!git gone | xargs -r git branch -D"

# ===== Remotes =====
git config --global alias.pl "pull --rebase"
git config --global alias.ps "push"
git config --global alias.psf "push --force-with-lease"
git config --global alias.sync "!git fetch origin && git rebase origin/main"

# ===== Stash =====
git config --global alias.sl "stash list"
git config --global alias.ss "stash save"
git config --global alias.sp "stash pop"

# ===== Information =====
git config --global alias.aliases "config --get-regexp alias"
git config --global alias.contributors "shortlog -sn --no-merges"
git config --global alias.whoami "config user.email"
```

### Usage Examples

```bash
# Instead of: git status -sb
git s

# Instead of: git checkout -b feature/new-thing
git cob feature/new-thing

# Instead of: git log --oneline --graph --decorate -20
git lg

# Instead of: git commit --amend --no-edit
git amend

# View all aliases
git aliases
```

---

## Signed Commits (GPG)

### Why Sign Commits?

Signed commits prove that YOU made the commit, not someone impersonating you. Required for:

- Open source contributions
- Enterprise environments
- Security-sensitive projects

### Setup GPG Signing

```bash
# 1. Check for existing keys
gpg --list-secret-keys --keyid-format=long

# 2. Generate new key (if needed)
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, no expiration, your email

# 3. Get your key ID
gpg --list-secret-keys --keyid-format=long
# Output shows: sec   rsa4096/YOUR_KEY_ID 2026-01-23
#                     ^^^^^^^^^^^^^ This is your key ID

# 4. Configure Git to use your key
git config --global user.signingkey YOUR_KEY_ID

# 5. Enable signing for all commits
git config --global commit.gpgsign true

# 6. Enable signing for all tags
git config --global tag.gpgsign true
```

### Add Key to GitHub

```bash
# Export your public key
gpg --armor --export YOUR_KEY_ID

# Copy output and add to:
# GitHub → Settings → SSH and GPG keys → New GPG key
```

### Signing Commands

```bash
# Sign a commit manually
git commit -S -m "Signed commit message"

# Sign a tag
git tag -s v1.0.0 -m "Signed release"

# Verify signatures
git log --show-signature
git verify-commit abc1234
git verify-tag v1.0.0
```

### Troubleshooting GPG

```bash
# "gpg failed to sign the data"
export GPG_TTY=$(tty)  # Add to ~/.zshrc or ~/.bashrc

# macOS: Install pinentry
brew install pinentry-mac
echo "pinentry-program /opt/homebrew/bin/pinentry-mac" >> ~/.gnupg/gpg-agent.conf
gpgconf --kill gpg-agent
```

---

## Git Stash

### What is Git Stash?

Temporarily saves uncommitted changes so you can switch contexts without committing half-done work.

```text
Working Directory (dirty) → git stash → Working Directory (clean)
                                ↓
                          Stash Stack
                          [stash@{0}]  ← Most recent
                          [stash@{1}]
                          [stash@{2}]
```

### Stash Commands

```bash
# Save current changes to stash
git stash
git stash save "WIP: voice feature"  # With message

# Include untracked files
git stash -u
git stash --include-untracked

# Include ignored files too
git stash -a
git stash --all

# List all stashes
git stash list

# View stash contents
git stash show              # Summary
git stash show -p           # Full diff
git stash show stash@{2}    # Specific stash

# Apply stash (keep in stack)
git stash apply
git stash apply stash@{2}   # Specific stash

# Apply and remove from stack
git stash pop
git stash pop stash@{2}

# Remove stash without applying
git stash drop stash@{0}

# Clear all stashes
git stash clear

# Create branch from stash
git stash branch feature-from-stash stash@{0}
```

### Stash Workflow

```bash
# Scenario: Working on feature, need to fix urgent bug

# 1. Stash current work
git stash save "WIP: voice interface integration"

# 2. Fix the bug
git checkout main
git checkout -b hotfix/critical-bug
# ... fix bug ...
git commit -m "fix: Critical bug"
git push

# 3. Return to feature
git checkout feature/voice-interface
git stash pop

# Continue working!
```

### Partial Stash

```bash
# Stash only specific files
git stash push -m "Partial stash" path/to/file.py another/file.js

# Interactive stash (choose hunks)
git stash -p
```

---

## Git Tags

### What are Tags?

Tags mark specific points in history, typically used for releases. Unlike branches, tags don't move.

| Type            | Description                | Use Case               |
| --------------- | -------------------------- | ---------------------- |
| **Lightweight** | Just a pointer to a commit | Quick local bookmarks  |
| **Annotated**   | Full object with metadata  | Releases (recommended) |

### Tag Commands

```bash
# List tags
git tag
git tag -l "v1.*"           # Filter pattern

# Create lightweight tag
git tag v1.0.0

# Create annotated tag (recommended)
git tag -a v1.0.0 -m "Release version 1.0.0"

# Tag a specific commit
git tag -a v0.9.0 abc1234 -m "Retroactive tag"

# View tag details
git show v1.0.0

# Push tags to remote
git push origin v1.0.0      # Single tag
git push origin --tags      # All tags

# Delete tag
git tag -d v1.0.0           # Local
git push origin --delete v1.0.0  # Remote

# Checkout tag (detached HEAD)
git checkout v1.0.0

# Create branch from tag
git checkout -b hotfix/v1.0.1 v1.0.0
```

### Semantic Versioning

```text
v1.2.3
│ │ │
│ │ └── PATCH: Bug fixes (backwards compatible)
│ └──── MINOR: New features (backwards compatible)
└────── MAJOR: Breaking changes
```

**Pre-release tags:**

```text
v1.0.0-alpha.1
v1.0.0-beta.1
v1.0.0-rc.1     # Release candidate
v1.0.0
```

### One Mind Release Workflow

```bash
# 1. Ensure you're on main with latest
git checkout main
git pull origin main

# 2. Create annotated tag
git tag -a v1.2.0 -m "Release v1.2.0

Features:
- Voice interface integration
- Memory improvements

Fixes:
- API timeout issues
- Memory leak in agent"

# 3. Push tag
git push origin v1.2.0

# 4. GitHub will trigger release workflow
```

---

## Git Cherry-Pick

### What is Cherry-Pick?

Applies a specific commit from one branch to another without merging the entire branch.

```text
main:     A ─── B ─── C
                      │
feature:  D ─── E ─── F ─── G

# Cherry-pick F to main:
main:     A ─── B ─── C ─── F'
```

### Cherry-Pick Commands

```bash
# Apply single commit to current branch
git cherry-pick abc1234

# Apply multiple commits
git cherry-pick abc1234 def5678

# Apply range of commits
git cherry-pick abc1234..xyz9876    # Excludes first
git cherry-pick abc1234^..xyz9876   # Includes first

# Don't commit automatically
git cherry-pick --no-commit abc1234

# Continue after resolving conflicts
git cherry-pick --continue

# Abort cherry-pick
git cherry-pick --abort

# Skip current commit and continue
git cherry-pick --skip
```

### Use Cases

```bash
# Scenario 1: Hotfix needed on main but was made on feature branch
git checkout main
git cherry-pick abc1234  # The fix commit
git push

# Scenario 2: Backport fix to older release
git checkout release-1.x
git cherry-pick def5678
git push

# Scenario 3: Grab multiple commits
git log feature-branch --oneline  # Find commits
git checkout main
git cherry-pick commit1 commit2 commit3
```

### Cherry-Pick Best Practices

1. **Use sparingly** — Prefer merging when possible
2. **Document** — Note which commits were cherry-picked
3. **Avoid duplicates** — Track what's been picked
4. **Test** — Cherry-picked code may lack dependencies

---

## Git Bisect

### What is Git Bisect?

Binary search through commits to find which one introduced a bug. Finds the problem in O(log n) steps.

```text
Good ─── ? ─── ? ─── ? ─── ? ─── Bad
                  ↓
        Binary search to find bad commit
```

### Bisect Commands

```bash
# Start bisecting
git bisect start

# Mark current commit as bad
git bisect bad

# Mark a known good commit
git bisect good abc1234

# Git checks out middle commit, you test, then:
git bisect good    # If this commit is OK
git bisect bad     # If this commit has the bug

# Repeat until Git finds the first bad commit

# End bisecting
git bisect reset
```

### Bisect Workflow

```bash
# 1. Start bisect
git bisect start

# 2. Mark endpoints
git bisect bad                    # Current is broken
git bisect good v1.0.0            # v1.0.0 was working

# 3. Git checks out a commit in the middle
# Output: Bisecting: 12 revisions left to test

# 4. Run your test
npm test  # or whatever confirms bug

# 5. Mark result
git bisect good  # Test passed
# OR
git bisect bad   # Test failed

# 6. Repeat step 4-5 until:
# Output: abc1234 is the first bad commit

# 7. Examine the commit
git show abc1234

# 8. Reset to original state
git bisect reset
```

### Automated Bisect

```bash
# Use a script to test automatically
git bisect start
git bisect bad HEAD
git bisect good v1.0.0
git bisect run npm test

# Or with a custom script
git bisect run ./test-for-bug.sh
```

The script should exit:

- `0` = good
- `1-127` (except 125) = bad
- `125` = skip (can't test this commit)

---

## Git Hooks

### What are Git Hooks?

Scripts that run automatically at certain points in the Git workflow. Located in `.git/hooks/`.

### Available Hooks

| Hook                   | Trigger               | Use Case                 |
| ---------------------- | --------------------- | ------------------------ |
| **pre-commit**         | Before commit         | Lint, format, tests      |
| **prepare-commit-msg** | After default message | Auto-fill templates      |
| **commit-msg**         | After message entered | Validate message format  |
| **post-commit**        | After commit          | Notifications            |
| **pre-push**           | Before push           | Run full tests           |
| **pre-rebase**         | Before rebase         | Prevent rebase on shared |
| **post-merge**         | After merge           | Install dependencies     |
| **post-checkout**      | After checkout        | Setup environment        |

### Creating Hooks

```bash
# 1. Create hook file
touch .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# 2. Edit the file
```

**Example: pre-commit hook**

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running pre-commit checks..."

# Run linter
if ! npm run lint; then
    echo "❌ Linting failed. Fix errors before committing."
    exit 1
fi

# Check for secrets
if git secrets --scan-history 2>/dev/null | grep -q "FOUND"; then
    echo "❌ Secrets detected! Remove before committing."
    exit 1
fi

# Run tests
if ! npm test; then
    echo "❌ Tests failed. Fix before committing."
    exit 1
fi

echo "✅ All checks passed!"
exit 0
```

**Example: commit-msg hook**

```bash
#!/bin/bash
# .git/hooks/commit-msg

commit_msg=$(cat "$1")
pattern="^(feat|fix|docs|style|refactor|test|chore|perf)(\(.+\))?: .{1,50}"

if ! [[ $commit_msg =~ $pattern ]]; then
    echo "❌ Invalid commit message format!"
    echo "Expected: type(scope): message"
    echo "Example: feat(agent): Add voice interface"
    exit 1
fi

exit 0
```

### Using Husky (Recommended)

Husky makes hooks easy to share with the team:

```bash
# Install
npm install husky --save-dev

# Initialize
npx husky init

# Create pre-commit hook
echo "npm test" > .husky/pre-commit

# Create commit-msg hook
echo 'npx --no -- commitlint --edit "$1"' > .husky/commit-msg
```

### One Mind Hook Setup

```bash
# .husky/pre-commit
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

# Python checks
cd agents/legacy && python -m ruff check . && cd ../..

# Frontend checks
cd ui && pnpm lint && cd ..

# Secrets scan
git secrets --scan

echo "✅ Pre-commit passed"
```

---

## Git Submodules

### What are Submodules?

A way to include one Git repository inside another. The parent repo stores a pointer to a specific commit in the submodule.

```text
onemind-os/           # Main repo
├── agents/
├── ui/
└── shared-utils/     # Submodule → points to commit abc123
```

### Submodule Commands

```bash
# Add a submodule
git submodule add https://github.com/user/repo.git path/to/submodule

# Clone repo with submodules
git clone --recursive https://github.com/user/repo.git
# OR after cloning:
git submodule update --init --recursive

# Update submodule to latest
cd path/to/submodule
git checkout main
git pull
cd ..
git add path/to/submodule
git commit -m "chore: Update submodule"

# Update all submodules
git submodule update --remote --merge

# Remove submodule
git submodule deinit path/to/submodule
git rm path/to/submodule
rm -rf .git/modules/path/to/submodule
```

### Submodule Workflow

```bash
# 1. After pulling main repo
git pull
git submodule update --init --recursive

# 2. Working in submodule
cd path/to/submodule
git checkout main
git pull
# Make changes...
git add .
git commit -m "feat: Add feature"
git push

# 3. Update parent repo
cd ..
git add path/to/submodule
git commit -m "chore: Update submodule to latest"
git push
```

### Submodule Gotchas

| Issue                        | Solution                                 |
| ---------------------------- | ---------------------------------------- |
| Forgot to init               | `git submodule update --init`            |
| Detached HEAD                | `cd submodule && git checkout main`      |
| Teammates don't have changes | They need to run `git submodule update`  |
| Merge conflicts              | Resolve in submodule, then update parent |

---

## Git Worktrees

### What are Worktrees?

Multiple working directories from one repository. Work on multiple branches simultaneously without stashing or cloning.

```text
~/code/
├── onemind-os/              # Main worktree (main branch)
├── onemind-os-feature/      # Second worktree (feature branch)
└── onemind-os-hotfix/       # Third worktree (hotfix branch)
```

### Worktree Commands

```bash
# List worktrees
git worktree list

# Add worktree for existing branch
git worktree add ../onemind-os-feature feature/voice-interface

# Add worktree and create new branch
git worktree add -b feature/new-thing ../onemind-os-new main

# Remove worktree
git worktree remove ../onemind-os-feature
# OR delete folder and:
git worktree prune

# Move worktree
git worktree move ../old-path ../new-path
```

### Worktree Workflow

```bash
# Scenario: Need to fix bug while working on feature

# 1. Working on feature in main directory
# ~/code/onemind-os (feature/voice-interface)

# 2. Create worktree for hotfix
git worktree add ../onemind-os-hotfix -b hotfix/critical-fix main

# 3. Fix bug in hotfix worktree
cd ../onemind-os-hotfix
# ... make fixes ...
git commit -m "fix: Critical bug"
git push

# 4. Return to feature
cd ../onemind-os
# Continue feature work!

# 5. Cleanup when done
git worktree remove ../onemind-os-hotfix
```

### Why Worktrees > Multiple Clones

| Worktrees               | Multiple Clones      |
| ----------------------- | -------------------- |
| Shared `.git` directory | Separate `.git` each |
| Disk efficient          | Uses more disk space |
| Shared config/hooks     | Separate config each |
| Easy to manage          | Can get messy        |

---

## Git Reflog

### What is Reflog?

The **reference log** records every change to HEAD and branch tips. It's your safety net — you can recover almost anything.

```bash
# Reflog is LOCAL only
# Keeps data for 90 days by default
# Even "deleted" commits are recoverable
```

### Reflog Commands

```bash
# View reflog
git reflog
git reflog show HEAD        # Same as above
git reflog show main        # For specific branch

# Output format:
# abc1234 HEAD@{0}: commit: Add feature
# def5678 HEAD@{1}: checkout: moving from main to feature
# ghi9012 HEAD@{2}: pull: Fast-forward

# View with timestamps
git reflog --date=relative

# Expire old entries
git reflog expire --expire=now --all
git gc --prune=now
```

### Recovery with Reflog

```bash
# Scenario 1: Undo a hard reset
git reset --hard HEAD~3     # Oops, went too far!
git reflog                  # Find the commit you want
git reset --hard abc1234    # Recover!

# Scenario 2: Recover deleted branch
git branch -D feature-branch  # Oops!
git reflog                    # Find last commit of that branch
git checkout -b feature-branch abc1234  # Recovered!

# Scenario 3: Recover after bad rebase
git rebase -i HEAD~5        # Made mistakes
git reflog                  # Find state before rebase
git reset --hard HEAD@{5}   # Go back

# Scenario 4: Find lost commit
git reflog --all | grep "lost commit message"
git cherry-pick abc1234
```

### Reflog Best Practices

1. **Check reflog first** before panicking
2. **Don't run `git gc`** until you're sure you don't need old data
3. **Reflog is local** — can't recover remote-only commits
4. **Expires after 90 days** — recover quickly if needed

---

## Remote Management

### What are Remotes?

Remote repositories are versions of your project hosted elsewhere (GitHub, GitLab, etc.).

### Remote Commands

```bash
# List remotes
git remote -v

# Add remote
git remote add origin https://github.com/user/repo.git
git remote add upstream https://github.com/original/repo.git

# Rename remote
git remote rename origin github

# Remove remote
git remote remove upstream

# Change remote URL
git remote set-url origin https://github.com/user/new-repo.git

# Show remote details
git remote show origin
```

### Multi-Remote Workflow

```bash
# Typical fork workflow
origin    → Your fork
upstream  → Original repository

# Setup
git remote add upstream https://github.com/original/repo.git

# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main
git push origin main

# Or with rebase
git fetch upstream
git rebase upstream/main
git push --force-with-lease origin main
```

### Remote Tracking Branches

```bash
# List all remote tracking branches
git branch -r

# See relationship between local and remote
git branch -vv

# Set upstream for current branch
git branch --set-upstream-to=origin/main

# Push and set upstream
git push -u origin feature-branch
```

---

## Forking Workflow

### What is Forking?

Creating a personal copy of someone else's repository. Standard for open source contributions.

```text
Original Repo (upstream)
        ↓ Fork
Your Fork (origin)
        ↓ Clone
Local Repository
```

### Fork Workflow

```bash
# 1. Fork on GitHub (UI button)

# 2. Clone your fork
git clone https://github.com/YOUR-USER/repo.git
cd repo

# 3. Add upstream remote
git remote add upstream https://github.com/ORIGINAL-OWNER/repo.git

# 4. Keep fork updated
git fetch upstream
git checkout main
git merge upstream/main
git push origin main

# 5. Create feature branch
git checkout -b feature/my-contribution

# 6. Work and commit
# ... make changes ...
git add .
git commit -m "feat: Add awesome feature"

# 7. Push to your fork
git push origin feature/my-contribution

# 8. Create Pull Request (GitHub UI)
# Compare: ORIGINAL-OWNER/repo main ← YOUR-USER/repo feature/my-contribution
```

### Keeping Fork Synced

```bash
# Quick sync script
#!/bin/bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
echo "Fork synced with upstream!"
```

---

## Branching Strategies

### Git Flow

Complex, structured. Good for scheduled releases.

```text
main ─────────────────────────────────────────────
  │                    ↑          ↑
  ├── develop ─────────┼──────────┼───────────────
  │     │              │          │
  │     ├── feature/a ─┘          │
  │     │                         │
  │     ├── feature/b ────────────┘
  │     │
  │     └── release/1.0 ──→ main
  │
  └── hotfix/critical ──→ main & develop
```

**Branches:**

- `main` — Production code
- `develop` — Integration branch
- `feature/*` — New features
- `release/*` — Release preparation
- `hotfix/*` — Production fixes

### GitHub Flow

Simple, continuous deployment. Good for SaaS.

```text
main ─── A ─── B ─── C ─── D ─── E ─── F
              │                   │
              └── feature ────────┘
                  (PR + Deploy)
```

**Rules:**

1. `main` is always deployable
2. Branch off `main` for features
3. Open PR when ready
4. Deploy after merge

### Trunk-Based Development

Minimal branches, fast integration. Good for experienced teams.

```text
main ─── A ─── B ─── C ─── D ─── E ─── F
              │     │
              └─────┘ (short-lived feature branches)
```

**Rules:**

1. Everyone commits to `main` (or very short branches)
2. Feature flags for incomplete features
3. CI/CD runs on every commit
4. Deploy multiple times per day

### One Mind Strategy

We use **GitHub Flow** with enhancements:

```bash
main           # Production (protected)
  │
  ├── feature/voice-interface      # Features
  ├── bugfix/memory-leak           # Bug fixes
  ├── hotfix/auth-critical         # Urgent fixes (fast-track)
  └── docs/api-reference           # Documentation
```

**Rules:**

1. `main` is protected, requires PR
2. Feature branches from `main`
3. Squash merge to keep history clean
4. Delete branch after merge
5. Deploy on merge to `main`

---

## Related Documents

- [Development Workflow](<../03%20Protocols%20(SOPs)/Development-Workflow.md>)
- [Code Review Guidelines](Code-Review-Guidelines.md)
- [CI/CD Pipeline](<../09%20Automation%20(Workflows)/CI-CD-Pipeline.md>)

---

**Last Updated:** January 23, 2026  
**Author:** Zeus Delacruz  
**Status:** Active Reference
