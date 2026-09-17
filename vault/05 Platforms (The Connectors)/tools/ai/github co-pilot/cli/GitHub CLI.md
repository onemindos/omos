---
title: "CLI — gh (GitHub CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/github]
ai_summary: "GitHub CLI (gh) quick reference — issues, PRs, repos, workflows, and release management for the One Mind GitHub org."
---

# gh — GitHub CLI

> Manage GitHub without leaving the terminal. Auth + org: `OneMind-OS` / `One-Mind-OS`.

---

## Auth

```bash
gh auth status
gh auth login
gh auth refresh
```

---

## Issues

```bash
# List open issues
gh issue list --repo OneMind-OS/onemind-codex

# Create issue
gh issue create --title "Fix SO inbox triage" --body "..." --label "so"

# View issue
gh issue view 42

# Close issue
gh issue close 42 --comment "Done."
```

---

## Pull Requests

```bash
# Create PR from current branch
gh pr create --title "feat: add CLI reference" --body "..."

# List PRs
gh pr list --repo OneMind-OS/onemind-codex

# View PR
gh pr view 12

# Merge PR
gh pr merge 12 --squash

# Check out a PR locally
gh pr checkout 12
```

---

## Repos

```bash
# View repo info
gh repo view OneMind-OS/onemind-codex

# Clone a repo
gh repo clone OneMind-OS/onemind-codex

# Create a new repo
gh repo create OneMind-OS/new-repo --private

# List repos in org
gh repo list OneMind-OS
```

---

## Workflows (GitHub Actions)

```bash
# List workflows
gh workflow list --repo OneMind-OS/onemind-codex-live

# Run a workflow manually
gh workflow run deploy.yml

# View run status
gh run list --repo OneMind-OS/onemind-codex-live

# Watch a run live
gh run watch
```

---

## Releases

```bash
# List releases
gh release list --repo OneMind-OS/onemind-web

# Create release
gh release create v1.2.0 --title "v1.2.0" --notes "..."

# Download release asset
gh release download v1.2.0
```

---

## Secrets (GitHub Actions)

```bash
# Set a secret
gh secret set SECRET_NAME --body "value"

# List secrets
gh secret list
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
