---
title: "GitHub"
type: reference
status: active
tags: [platform, development, git, ci-cd]
---
# GitHub

> All repos, CI/CD, collaboration. Primary code hosting.

## Key Repos
- `OneMind-OS/One-Mind` — This vault
- `OneMind-OS/onemind-app` — Flutter cockpit
- `OneMind-OS/openclaw` — Business gateway

## CLI: gh
```bash
gh repo list
gh pr create --fill
gh pr merge
gh issue create --title "..." --body "..."
gh run list
gh auth status
```

## CLI: git
```bash
git status
git add -A && git commit -m "feat: ..."
git push origin main
git log --oneline -20
git stash && git pull --rebase && git stash pop
```

## Integration
- Hermes has native GitHub skills (PR workflow, code review, issues)
- Copilot for code completion
- Actions for CI/CD
