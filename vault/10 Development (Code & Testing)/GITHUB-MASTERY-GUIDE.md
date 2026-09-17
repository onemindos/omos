# GitHub Mastery Guide

**Comprehensive guide to GitHub platform, workflows, and personal repository management for One Mind development.**

---

## Table of Contents

1. [Introduction](#introduction)
2. [GitHub CLI (gh)](#github-cli-gh)
3. [Repository Management](#repository-management)
4. [Repository Settings & Best Practices](#repository-settings--best-practices)
5. [README Excellence](#readme-excellence)
6. [LICENSE Selection](#license-selection)
7. [.github Directory](#github-directory)
8. [Issue Templates](#issue-templates)
9. [Pull Request Templates](#pull-request-templates)
10. [Issues & Project Management](#issues--project-management)
11. [Pull Requests & Code Review](#pull-requests--code-review)
12. [Branch Protection Rules](#branch-protection-rules)
13. [GitHub Actions (CI/CD)](#github-actions-cicd)
14. [Workflows for One Mind](#workflows-for-one-mind)
15. [GitHub Secrets & Environments](#github-secrets--environments)
16. [GitHub Projects](#github-projects)
17. [Releases & Tags](#releases--tags)
18. [GitHub Packages](#github-packages)
19. [GitHub Pages](#github-pages)
20. [Security Features](#security-features)
21. [Collaboration & Permissions](#collaboration--permissions)
22. [CODEOWNERS](#codeowners)
23. [GitHub Discussions](#github-discussions)
24. [Repository Templates](#repository-templates)
25. [Personal Dashboard & Notifications](#personal-dashboard--notifications)
26. [GitHub Mobile](#github-mobile)
27. [One Mind GitHub Workflow](#one-mind-github-workflow)
28. [Troubleshooting](#troubleshooting)
29. [Quick Reference Card](#quick-reference-card)

---

## Introduction

**GitHub** is more than just Git hosting — it's a complete development platform with:

- **Code hosting** — Git repositories with web interface
- **Collaboration** — Pull requests, code review, issues
- **Automation** — GitHub Actions for CI/CD
- **Security** — Dependabot, code scanning, secrets management
- **Documentation** — README, wikis, GitHub Pages
- **Project management** — Projects, milestones, labels
- **Package hosting** — npm, Docker, Maven, NuGet

This guide covers **everything you need** to manage your personal repositories and collaborate effectively.

---

## GitHub CLI (gh)

### Installation

```bash
# macOS
brew install gh

# Linux
sudo apt install gh  # Debian/Ubuntu
sudo dnf install gh  # Fedora

# Windows
winget install GitHub.cli
```

### Authentication

```bash
# Login interactively
gh auth login

# Choose:
# 1. GitHub.com
# 2. HTTPS
# 3. Login with a web browser (recommended)
# OR paste a personal access token

# Check status
gh auth status

# Logout
gh auth logout
```

### Essential Commands

```bash
# === Repositories ===
gh repo create                    # Interactive create
gh repo create onemind-os --public
gh repo clone user/repo
gh repo view                      # View current repo
gh repo view user/repo            # View specific repo
gh repo list                      # List your repos
gh repo archive user/repo         # Archive repo
gh repo delete user/repo          # Delete repo (careful!)

# === Issues ===
gh issue list
gh issue list --state open
gh issue list --assignee @me
gh issue list --label bug
gh issue create                   # Interactive
gh issue create --title "Bug: Memory leak" --body "Details..."
gh issue view 123
gh issue close 123
gh issue reopen 123
gh issue comment 123 --body "Fixed in #456"

# === Pull Requests ===
gh pr list
gh pr list --state open
gh pr list --author @me
gh pr create                      # Interactive
gh pr create --title "feat: Add voice interface" --body "..."
gh pr view 456
gh pr checkout 456               # Checkout PR locally
gh pr diff 456
gh pr review 456 --approve
gh pr review 456 --request-changes --body "Needs tests"
gh pr merge 456
gh pr merge 456 --squash
gh pr merge 456 --rebase
gh pr close 456

# === GitHub Actions ===
gh workflow list
gh workflow view
gh workflow run                  # Trigger workflow
gh run list
gh run view 123456789
gh run watch                     # Watch current run
gh run download 123456789        # Download artifacts

# === Releases ===
gh release list
gh release create v1.0.0
gh release create v1.0.0 --title "Release 1.0.0" --notes "..."
gh release upload v1.0.0 ./dist/*
gh release view v1.0.0
gh release delete v1.0.0

# === Gists ===
gh gist create file.py
gh gist create --public file.py
gh gist list
gh gist view abc123
gh gist clone abc123

# === Secrets ===
gh secret list
gh secret set SECRET_NAME
gh secret set SECRET_NAME < secret.txt
gh secret remove SECRET_NAME

# === SSH Keys ===
gh ssh-key list
gh ssh-key add ~/.ssh/id_ed25519.pub --title "MacBook Pro"

# === Extensions ===
gh extension list
gh extension install owner/repo
```

### gh Aliases

```bash
# Add to ~/.config/gh/config.yml
aliases:
  pv: pr view
  pc: pr checkout
  ic: issue create
  co: repo clone
  rv: repo view
```

---

## Repository Management

### Creating Repositories

**Via GitHub CLI:**

```bash
# Interactive
gh repo create

# Quick create
gh repo create onemind-voice --public --description "Voice interface for One Mind"

# Create with README and .gitignore
gh repo create onemind-api --public --gitignore Python --license mit

# Create from template
gh repo create my-new-project --template owner/template-repo

# Create and clone
gh repo create onemind-mobile --public --clone
```

**Via Web Interface:**

1. Go to [github.com/new](https://github.com/new)
2. Fill in:
   - **Repository name** — Short, descriptive, kebab-case
   - **Description** — One-line summary
   - **Visibility** — Public or Private
   - **Initialize** — README, .gitignore, LICENSE
3. Click "Create repository"

### Repository Naming Conventions

| Type             | Pattern             | Example                       |
| ---------------- | ------------------- | ----------------------------- |
| **Main project** | `project-name`      | `onemind-os`                  |
| **Sub-project**  | `project-component` | `onemind-agent`, `onemind-ui` |
| **Library**      | `library-name`      | `memory-toolkit`              |
| **Tool**         | `tool-name`         | `deployment-helper`           |
| **Template**     | `template-name`     | `python-service-template`     |
| **Config**       | `config-type`       | `dotfiles`, `vim-config`      |

**Naming Rules:**

- Use lowercase
- Use hyphens (kebab-case), not underscores
- Be descriptive but concise
- Avoid abbreviations unless widely known
- Don't include your username in the name

### Cloning Repositories

```bash
# Clone via HTTPS
git clone https://github.com/user/repo.git

# Clone via SSH (preferred)
git clone git@github.com:user/repo.git

# Clone with gh CLI
gh repo clone user/repo

# Clone into specific directory
git clone git@github.com:user/repo.git my-folder

# Clone specific branch
git clone -b develop git@github.com:user/repo.git

# Shallow clone (faster, less history)
git clone --depth 1 git@github.com:user/repo.git
```

### Forking Repositories

```bash
# Fork via gh CLI
gh repo fork user/repo

# Fork and clone
gh repo fork user/repo --clone

# Fork with custom name
gh repo fork user/repo --fork-name my-fork
```

---

## Repository Settings & Best Practices

### Essential Settings

**Go to: Settings → General**

| Setting            | Recommendation                                        |
| ------------------ | ----------------------------------------------------- |
| **Description**    | Clear, one-line summary                               |
| **Website**        | Docs URL or live demo                                 |
| **Topics**         | Add 3-8 relevant tags                                 |
| **Features**       | ✅ Wikis (if needed), ✅ Issues, ✅ Projects          |
| **Pull Requests**  | ✅ Allow squash merging, ✅ Auto-delete head branches |
| **Merge button**   | Uncheck ❌ Allow merge commits (unless needed)        |
| **Default branch** | `main`                                                |

### Repository Topics

Add topics for discoverability:

```text
# Example for onemind-os
ai, agent, automation, python, nextjs, solar-punk, personal-os
```

**Topic Guidelines:**

- Use lowercase
- Use hyphens for multi-word topics
- Add language/framework tags (python, typescript, react)
- Add domain tags (ai, automation, iot)
- Max 20 topics

### Social Preview Image

**Settings → General → Social preview**

Create a 1280×640px image with:

- Project name/logo
- Tagline
- Key visual

Tools: [Figma](https://figma.com), [Canva](https://canva.com)

### Archiving Repositories

When a project is no longer maintained:

```bash
# Via CLI
gh repo archive user/repo

# Via web: Settings → General → Danger Zone → Archive
```

**Archived repos:**

- ✅ Read-only
- ❌ Can't accept new issues/PRs
- ⚠️ Can be unarchived later

---

## README Excellence

The README is your project's **front page**. Make it count.

### Anatomy of a Great README

```markdown
# Project Name

> One-line tagline or description

[![Build Status](badge-url)](link)
[![License](badge-url)](LICENSE)

![Demo GIF or Screenshot](media/demo.gif)

## Features

- ✨ Feature 1
- 🚀 Feature 2
- 🔒 Feature 3

## Quick Start

\`\`\`bash

# Installation

npm install project-name

# Usage

project-name init
\`\`\`

## Installation

### Prerequisites

- Node.js 18+
- PostgreSQL 15+

### Steps

1. Clone the repository
   \`\`\`bash
   git clone https://github.com/user/repo.git
   cd repo
   \`\`\`

2. Install dependencies
   \`\`\`bash
   npm install
   \`\`\`

3. Configure environment
   \`\`\`bash
   cp .env.example .env

   # Edit .env with your values

   \`\`\`

4. Run the app
   \`\`\`bash
   npm start
   \`\`\`

## Usage

### Basic Example

\`\`\`javascript
import { Thing } from 'project-name';

const thing = new Thing();
thing.doStuff();
\`\`\`

### Advanced Example

See [examples/](examples/) directory.

## Configuration

| Variable  | Description  | Default  |
| --------- | ------------ | -------- |
| `API_KEY` | Your API key | Required |
| `PORT`    | Server port  | `3000`   |

## Development

\`\`\`bash

# Run tests

npm test

# Run linter

npm run lint

# Build

npm run build
\`\`\`

## Architecture

Brief overview or link to [ARCHITECTURE.md](ARCHITECTURE.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Roadmap

- [x] Feature 1
- [ ] Feature 2 (in progress)
- [ ] Feature 3 (planned)

## License

[MIT](LICENSE) © [Your Name](https://github.com/yourusername)

## Acknowledgments

- Thanks to [person/project]
- Inspired by [project]

## Contact

- **Author:** Your Name
- **Email:** you@example.com
- **Twitter:** [@yourusername](https://twitter.com/yourusername)
  \`\`\`
```

### README Best Practices

1. **Keep it scannable** — Use headers, lists, code blocks
2. **Show, don't tell** — Include screenshots, GIFs, or videos
3. **Get people started fast** — Quick Start section at the top
4. **Link to docs** — Don't put everything in README
5. **Update regularly** — Keep installation steps current
6. **Add badges** — Build status, coverage, version, license
7. **Be welcoming** — Encourage contributions

### Badges

```markdown
![Build](https://github.com/user/repo/workflows/CI/badge.svg)
![Coverage](https://codecov.io/gh/user/repo/branch/main/graph/badge.svg)
![Version](https://img.shields.io/github/v/release/user/repo)
![License](https://img.shields.io/github/license/user/repo)
![Contributors](https://img.shields.io/github/contributors/user/repo)
```

Generate at: [shields.io](https://shields.io)

### Visual Elements

| Element                  | Purpose           | Tool                                                                    |
| ------------------------ | ----------------- | ----------------------------------------------------------------------- |
| **Screenshot**           | Show the UI       | Native capture                                                          |
| **Demo GIF**             | Show interaction  | [Kap](https://getkap.co), [LICEcap](https://www.cockos.com/licecap/)    |
| **Architecture diagram** | Explain structure | [Excalidraw](https://excalidraw.com), [Mermaid](https://mermaid.js.org) |
| **Logo**                 | Brand identity    | Figma, Canva                                                            |

---

## LICENSE Selection

### Common Open Source Licenses

| License          | Permissions                                                              | Conditions                                                | Limitations                                     |
| ---------------- | ------------------------------------------------------------------------ | --------------------------------------------------------- | ----------------------------------------------- |
| **MIT**          | ✅ Commercial use<br>✅ Modification<br>✅ Distribution                  | ℹ️ Include license                                        | ❌ Liability<br>❌ Warranty                     |
| **Apache 2.0**   | ✅ Commercial use<br>✅ Modification<br>✅ Distribution<br>✅ Patent use | ℹ️ Include license<br>ℹ️ State changes                    | ❌ Liability<br>❌ Trademark use<br>❌ Warranty |
| **GPL-3.0**      | ✅ Commercial use<br>✅ Modification<br>✅ Distribution                  | ℹ️ Disclose source<br>ℹ️ Same license<br>ℹ️ State changes | ❌ Liability<br>❌ Warranty                     |
| **BSD 3-Clause** | ✅ Commercial use<br>✅ Modification<br>✅ Distribution                  | ℹ️ Include license                                        | ❌ Liability<br>❌ Use trademark<br>❌ Warranty |

### Which License?

```text
┌─ Want people to share improvements?
│  ├─ Yes → GPL-3.0
│  └─ No ↓
│
└─ Want patent protection?
   ├─ Yes → Apache 2.0
   └─ No → MIT
```

**For One Mind projects:** Use **MIT** for libraries/tools, **Apache 2.0** for main applications.

### Adding a License

```bash
# Via GitHub CLI
gh repo create my-project --license mit

# Or add manually
# 1. Create LICENSE file
# 2. Copy license text from choosealicense.com
# 3. Update copyright year and name
```

**Update README:**

```markdown
## License

[MIT](LICENSE) © 2026 Zeus Delacruz
```

---

## .github Directory

The `.github/` directory contains GitHub-specific configuration.

```text
.github/
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   ├── feature_request.md
│   └── config.yml
├── PULL_REQUEST_TEMPLATE.md
├── workflows/
│   ├── ci.yml
│   ├── deploy.yml
│   └── release.yml
├── FUNDING.yml
├── CODEOWNERS
├── dependabot.yml
└── SECURITY.md
```

### FUNDING.yml

Support your project:

```yaml
# .github/FUNDING.yml
github: Zeus-Delacruz
patreon: onemind
custom: ["https://buymeacoffee.com/zeus", "https://paypal.me/zeus"]
```

### SECURITY.md

Security policy:

```markdown
# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**Do not open a public issue.**

Email security@onemindcodex.com with:

1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)

We will respond within 48 hours.
```

### dependabot.yml

Automated dependency updates:

```yaml
# .github/dependabot.yml
version: 2
updates:
  # Frontend dependencies
  - package-ecosystem: "npm"
    directory: "/ui"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "Zeus-Delacruz"
    labels:
      - "dependencies"
      - "javascript"

  # Python dependencies
  - package-ecosystem: "pip"
    directory: "/agents"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "Zeus-Delacruz"
    labels:
      - "dependencies"
      - "python"

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "monthly"
```

---

## Issue Templates

### Creating Issue Templates

**Create:** `.github/ISSUE_TEMPLATE/bug_report.md`

```markdown
---
name: Bug Report
about: Report a bug to help us improve
title: "[BUG] "
labels: bug
assignees: ""
---

## Bug Description

A clear description of the bug.

## Steps to Reproduce

1. Go to '...'
2. Click on '...'
3. Scroll down to '...'
4. See error

## Expected Behavior

What should happen.

## Actual Behavior

What actually happens.

## Screenshots

If applicable, add screenshots.

## Environment

- **OS:** [e.g., macOS 14.1]
- **Browser:** [e.g., Chrome 120]
- **Version:** [e.g., 1.2.0]

## Additional Context

Any other information.
```

**Create:** `.github/ISSUE_TEMPLATE/feature_request.md`

```markdown
---
name: Feature Request
about: Suggest a new feature
title: "[FEATURE] "
labels: enhancement
assignees: ""
---

## Problem Statement

Describe the problem this feature would solve.

## Proposed Solution

Describe your proposed solution.

## Alternatives Considered

Alternative solutions you've considered.

## Additional Context

Mockups, examples, or related issues.

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
```

### Issue Template Chooser

**Create:** `.github/ISSUE_TEMPLATE/config.yml`

```yaml
blank_issues_enabled: false
contact_links:
  - name: 💬 Discussions
    url: https://github.com/user/repo/discussions
    about: Ask questions and discuss ideas
  - name: 📚 Documentation
    url: https://docs.onemindcodex.com
    about: Read the documentation
```

---

## Pull Request Templates

**Create:** `.github/PULL_REQUEST_TEMPLATE.md`

```markdown
## Description

Brief description of changes.

## Type of Change

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📝 Documentation update
- [ ] 🔨 Refactoring
- [ ] ⚡ Performance improvement

## Related Issues

Fixes #123
Relates to #456

## Changes Made

- Change 1
- Change 2
- Change 3

## Screenshots

If applicable, add screenshots.

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass locally
- [ ] Dependent changes merged

## Additional Notes

Any additional information for reviewers.
```

---

## Issues & Project Management

### Issue Labels

**Standard labels:**

| Label              | Color     | Use                 |
| ------------------ | --------- | ------------------- |
| `bug`              | 🔴 Red    | Something broken    |
| `enhancement`      | 🟢 Green  | New feature         |
| `documentation`    | 🔵 Blue   | Docs improvement    |
| `good first issue` | 🟣 Purple | Easy for newcomers  |
| `help wanted`      | 🟡 Yellow | Need assistance     |
| `question`         | 🟠 Orange | Question/discussion |
| `wontfix`          | ⚫ Black  | Won't be addressed  |
| `duplicate`        | ⚪ Gray   | Duplicate issue     |

**One Mind custom labels:**

| Label            | Color     | Use             |
| ---------------- | --------- | --------------- |
| `agent`          | `#0E8A16` | Agent-related   |
| `ui`             | `#1D76DB` | UI-related      |
| `infrastructure` | `#FBCA04` | Infrastructure  |
| `security`       | `#D73A4A` | Security issue  |
| `performance`    | `#A2EEEF` | Performance     |
| `breaking`       | `#B60205` | Breaking change |

### Managing Issues

```bash
# Create issue
gh issue create --title "Bug: Memory leak" --label bug --assignee @me

# List issues
gh issue list --label bug --state open

# View issue
gh issue view 123

# Add labels
gh issue edit 123 --add-label "performance,urgent"

# Assign
gh issue edit 123 --add-assignee @me

# Close issue
gh issue close 123

# Reopen issue
gh issue reopen 123

# Link to PR
# In PR description:
# Closes #123
# Fixes #456
# Resolves #789
```

### Issue References

```markdown
# Same repo

#123

# Different repo

user/repo#123

# Cross-link

Fixes #123
Closes #456
Resolves #789
Relates to #012
Duplicate of #345
```

### Milestones

**Create milestone:**

1. Go to **Issues → Milestones → New milestone**
2. Fill in:
   - **Title:** v1.0.0 Release
   - **Due date:** 2026-02-15
   - **Description:** Features for version 1.0
3. Create

**Assign issues to milestone:**

```bash
gh issue edit 123 --milestone "v1.0.0"
```

---

## Pull Requests & Code Review

### Creating Pull Requests

```bash
# Via gh CLI (recommended)
gh pr create

# With options
gh pr create \
  --title "feat: Add voice interface" \
  --body "Implements voice commands using LiveKit" \
  --label enhancement \
  --assignee @me \
  --reviewer teammate

# Draft PR
gh pr create --draft

# Convert draft to ready
gh pr ready
```

### PR Best Practices

1. **One purpose per PR** — Don't mix features with fixes
2. **Keep it small** — Aim for <400 lines changed
3. **Write clear title** — Use conventional commits format
4. **Fill out template** — Don't leave sections blank
5. **Link issues** — Use "Fixes #123" to auto-close
6. **Self-review first** — Check your own diff
7. **Request specific reviewers** — Tag people with relevant expertise
8. **Respond to feedback** — Address all comments
9. **Keep updated** — Merge main regularly to avoid conflicts

### PR Title Format

Use [Conventional Commits](https://www.conventionalcommits.org):

```text
type(scope): description

Examples:
feat(agent): Add voice recognition
fix(ui): Fix memory leak in chat
docs(readme): Update installation steps
refactor(api): Simplify error handling
perf(db): Optimize query performance
test(agent): Add unit tests for voice
chore(deps): Update dependencies
```

### Code Review

**As a reviewer:**

```bash
# Checkout PR locally
gh pr checkout 456

# Review changes
gh pr diff 456

# Approve
gh pr review 456 --approve

# Request changes
gh pr review 456 --request-changes --body "Please add tests"

# Comment
gh pr review 456 --comment --body "Looks good overall, minor suggestions"
```

**Review comments:**

````markdown
# Suggest changes (GitHub web UI)

suggestion

```python
# Better approach
def optimized_function():
    return result
```
````

````

**Review guidelines:**

- ✅ **Be constructive** — Suggest improvements, don't just criticize
- ✅ **Be specific** — "This function is complex" → "Consider extracting lines 45-60 to a helper function"
- ✅ **Ask questions** — "Why did you choose this approach?"
- ✅ **Praise good work** — "Nice refactoring!"
- ❌ **Don't be pedantic** — Focus on important issues
- ❌ **Don't nitpick style** — Use automated linters instead

### Merging Strategies

| Strategy | Effect | Use When |
|----------|--------|----------|
| **Squash** | Combines all commits into one | Feature branches (default) |
| **Rebase** | Replays commits on top of base | Linear history preferred |
| **Merge commit** | Creates explicit merge commit | Want full history preserved |

```bash
# Merge via CLI
gh pr merge 456 --squash
gh pr merge 456 --rebase
gh pr merge 456 --merge
````

**One Mind default:** Squash merge

---

## Branch Protection Rules

Protect important branches from force pushes and require checks.

### Setting Up Protection

**Go to:** Settings → Branches → Add rule

**Rule for `main`:**

- ✅ **Require pull request before merging**
  - ✅ Require approvals: 1
  - ✅ Dismiss stale reviews when new commits pushed
  - ✅ Require review from Code Owners
- ✅ **Require status checks to pass**
  - ✅ Require branches to be up to date
  - Status checks: `CI`, `Lint`, `Tests`
- ✅ **Require conversation resolution before merging**
- ✅ **Require signed commits**
- ✅ **Require linear history**
- ✅ **Do not allow bypassing the above settings** (unless admin)

### Rulesets (New)

**Settings → Rules → New ruleset**

Rulesets are more powerful than branch protection:

```yaml
# Example ruleset
name: Main Branch Protection
target: main
enforcement: active

rules:
  - require_pull_request:
      required_approvals: 1
      dismiss_stale_reviews: true
  - require_status_checks:
      strict: true
      checks:
        - ci
        - lint
        - tests
  - require_signatures: true
  - block_force_pushes: true
```

---

## GitHub Actions (CI/CD)

### Workflow Basics

Workflows live in `.github/workflows/` and use YAML.

**Basic structure:**

```yaml
name: Workflow Name
on: [push, pull_request]

jobs:
  job-name:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run command
        run: echo "Hello"
```

### Triggers (on)

```yaml
# Push to specific branches
on:
  push:
    branches: [main, develop]

# Pull requests
on:
  pull_request:
    branches: [main]

# Manual trigger
on:
  workflow_dispatch:
    inputs:
      environment:
        description: 'Deployment environment'
        required: true
        default: 'staging'

# Schedule (cron)
on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight

# Multiple triggers
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  workflow_dispatch:
```

### Jobs & Steps

```yaml
jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: "npm"

      - name: Install dependencies
        run: npm ci

      - name: Run tests
        run: npm test

      - name: Build
        run: npm run build

  deploy:
    needs: build # Wait for build to complete
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main' # Only on main

    steps:
      - name: Deploy
        run: ./deploy.sh
```

### Matrix Builds

Test across multiple versions:

```yaml
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        node-version: [18, 20, 21]

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}
      - run: npm ci
      - run: npm test
```

### Caching Dependencies

```yaml
- name: Setup Node.js
  uses: actions/setup-node@v4
  with:
    node-version: "20"
    cache: "npm" # Automatically caches npm

# Or manual cache
- name: Cache dependencies
  uses: actions/cache@v3
  with:
    path: ~/.npm
    key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-node-
```

### Artifacts

```yaml
- name: Upload coverage
  uses: actions/upload-artifact@v4
  with:
    name: coverage-report
    path: coverage/

- name: Download coverage
  uses: actions/download-artifact@v4
  with:
    name: coverage-report
```

---

## Workflows for One Mind

### CI Workflow

`.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  lint-frontend:
    name: Lint Frontend
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: "npm"
          cache-dependency-path: ui/package-lock.json

      - name: Install dependencies
        working-directory: ui
        run: npm ci

      - name: Run ESLint
        working-directory: ui
        run: npm run lint

      - name: Type check
        working-directory: ui
        run: npm run type-check

  test-frontend:
    name: Test Frontend
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"
          cache: "npm"
          cache-dependency-path: ui/package-lock.json

      - name: Install dependencies
        working-directory: ui
        run: npm ci

      - name: Run tests
        working-directory: ui
        run: npm test -- --coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          directory: ui/coverage

  lint-backend:
    name: Lint Backend
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.12"
          cache: "pip"

      - name: Install dependencies
        run: |
          pip install ruff mypy
          pip install -r agents/legacy/requirements.txt

      - name: Run Ruff
        run: ruff check agents/

      - name: Run MyPy
        run: mypy agents/

  test-backend:
    name: Test Backend
    runs-on: ubuntu-latest

    services:
      postgres:
        image: pgvector/pgvector:pg16
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.12"
          cache: "pip"

      - name: Install dependencies
        run: pip install -r agents/legacy/requirements.txt

      - name: Run tests
        run: pytest tests/ --cov --cov-report=xml
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost/test

      - name: Upload coverage
        uses: codecov/codecov-action@v3

  security:
    name: Security Scan
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: "fs"
          scan-ref: "."
          format: "sarif"
          output: "trivy-results.sarif"

      - name: Upload results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: "trivy-results.sarif"
```

### Deploy Workflow

`.github/workflows/deploy.yml`

```yaml
name: Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    name: Deploy to Production
    runs-on: ubuntu-latest
    environment: production

    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Build and push Docker image
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/onemind-agent:$IMAGE_TAG .
          docker push $ECR_REGISTRY/onemind-agent:$IMAGE_TAG

      - name: Deploy to ECS
        run: |
          aws ecs update-service \
            --cluster onemind-cluster \
            --service legacy-agent \
            --force-new-deployment
```

### Release Workflow

`.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

jobs:
  release:
    name: Create Release
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Full history for changelog

      - name: Generate changelog
        id: changelog
        uses: metcalfc/changelog-generator@v4.1.0
        with:
          myToken: ${{ secrets.GITHUB_TOKEN }}

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          body: ${{ steps.changelog.outputs.changelog }}
          draft: false
          prerelease: ${{ contains(github.ref, 'alpha') || contains(github.ref, 'beta') }}
```

---

## GitHub Secrets & Environments

### Repository Secrets

**Settings → Secrets and variables → Actions → New repository secret**

```bash
# Via CLI
echo "secret_value" | gh secret set SECRET_NAME

# From file
gh secret set SECRET_NAME < secret.txt

# List secrets
gh secret list
```

**Common secrets:**

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `DATABASE_URL`
- `API_KEY`
- `SLACK_WEBHOOK_URL`

**Use in workflow:**

```yaml
steps:
  - name: Deploy
    env:
      API_KEY: ${{ secrets.API_KEY }}
    run: ./deploy.sh
```

### Environment Secrets

**Settings → Environments → New environment**

Create environments: `development`, `staging`, `production`

**Environment settings:**

- ✅ **Required reviewers** — Manual approval before deployment
- ✅ **Wait timer** — Delay before deployment
- ✅ **Deployment branches** — Limit to specific branches

**Use in workflow:**

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production # Uses production secrets

    steps:
      - name: Deploy
        env:
          API_KEY: ${{ secrets.API_KEY }} # From production environment
        run: ./deploy.sh
```

---

## GitHub Projects

Modern project management integrated with GitHub.

### Creating a Project

1. Go to your profile → **Projects → New project**
2. Choose template:
   - **Table** — Spreadsheet view
   - **Board** — Kanban board
   - **Roadmap** — Timeline view
3. Name it: "One Mind Development"

### Views

- **Board** — Kanban columns (To Do, In Progress, Done)
- **Table** — All fields visible
- **Roadmap** — Gantt-style timeline

### Custom Fields

Add custom fields:

- **Status** — To Do, In Progress, Done, Blocked
- **Priority** — Low, Medium, High, Critical
- **Size** — XS, S, M, L, XL
- **Sprint** — Sprint 1, Sprint 2, etc.
- **Area** — Agent, UI, Infrastructure, Docs

### Automation

**Settings → Workflows**

Auto-move items:

- ✅ **Item added to project** → Set status to "To Do"
- ✅ **Pull request merged** → Set status to "Done"
- ✅ **Issue closed** → Set status to "Done"
- ✅ **Pull request opened** → Set status to "In Progress"

### Linking Issues & PRs

```bash
# Add to project
gh project item-add <project-id> --owner @me --url https://github.com/user/repo/issues/123
```

Or use web UI: Issue → Projects → Add to project

---

## Releases & Tags

### Creating Releases

```bash
# Via CLI
gh release create v1.0.0 \
  --title "Release 1.0.0" \
  --notes "## What's New

- Voice interface
- Memory improvements
- Bug fixes

## Breaking Changes

- API endpoint changes

See [CHANGELOG.md](CHANGELOG.md) for details."

# Prerelease
gh release create v1.0.0-beta.1 --prerelease

# Draft release
gh release create v1.0.0 --draft

# Upload assets
gh release create v1.0.0 ./dist/onemind-*.zip
```

### Release Assets

Upload binaries, packages, or documentation:

```bash
gh release upload v1.0.0 \
  ./dist/onemind-macos-arm64.tar.gz \
  ./dist/onemind-linux-x64.tar.gz \
  ./dist/onemind-windows-x64.zip
```

### Automated Releases

Use GitHub Actions:

```yaml
- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    files: |
      dist/*.tar.gz
      dist/*.zip
    generate_release_notes: true
```

### CHANGELOG.md

Keep a changelog following [Keep a Changelog](https://keepachangelog.com):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Feature in progress

## [1.0.0] - 2026-01-23

### Added

- Voice interface with LiveKit integration
- Memory system with vector search
- Real-time event streaming via NATS

### Changed

- Improved agent response time by 40%
- Updated UI to new design system

### Fixed

- Memory leak in chat component
- API timeout issues
- Database connection pooling

### Security

- Updated dependencies with security patches

## [0.9.0] - 2026-01-10

...
```

---

## GitHub Packages

Host packages (npm, Docker, Maven, NuGet, etc.) on GitHub.

### Docker Packages

```bash
# Login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Tag image
docker tag myimage ghcr.io/username/myimage:latest

# Push
docker push ghcr.io/username/myimage:latest

# Pull
docker pull ghcr.io/username/myimage:latest
```

**In workflow:**

```yaml
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Build and push
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: ghcr.io/${{ github.repository }}/agent:latest
```

### npm Packages

**package.json:**

```json
{
  "name": "@username/package-name",
  "publishConfig": {
    "registry": "https://npm.pkg.github.com"
  }
}
```

**.npmrc:**

```
@username:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=${GITHUB_TOKEN}
```

**Publish:**

```bash
npm publish
```

---

## GitHub Pages

Host static sites from your repository.

### Setup

**Settings → Pages**

1. **Source:** Deploy from a branch
2. **Branch:** Select `main` or `gh-pages`
3. **Folder:** `/` (root) or `/docs`

**Your site will be at:** `https://username.github.io/repo`

### Custom Domain

1. Add `CNAME` file with your domain:

   ```
   onemindcodex.com
   ```

2. **Settings → Pages → Custom domain**
   Enter: `onemindcodex.com`

3. Configure DNS:
   ```
   CNAME record: www → username.github.io
   A records: @ → 185.199.108.153, 185.199.109.153, 185.199.110.153, 185.199.111.153
   ```

### Deploy with Actions

```yaml
name: Deploy Pages

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "20"

      - name: Install and build
        run: |
          npm ci
          npm run build

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: "./dist"

      - name: Deploy to GitHub Pages
        uses: actions/deploy-pages@v4
```

---

## Security Features

### Dependabot

Automated security updates.

**Enable:** Settings → Security → Dependabot → Enable

**Configure:** `.github/dependabot.yml` (see earlier section)

### Code Scanning

Find security vulnerabilities.

**Enable:** Security → Code scanning → Set up CodeQL

```yaml
name: CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: "0 0 * * 1" # Weekly on Monday

jobs:
  analyze:
    runs-on: ubuntu-latest
    permissions:
      security-events: write

    steps:
      - uses: actions/checkout@v4

      - name: Initialize CodeQL
        uses: github/codeql-action/init@v2
        with:
          languages: javascript, python

      - name: Autobuild
        uses: github/codeql-action/autobuild@v2

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v2
```

### Secret Scanning

GitHub automatically scans for leaked secrets (API keys, tokens).

**Enable:** Settings → Security → Secret scanning → Enable

**Alerts:** Security → Secret scanning alerts

### Security Advisories

Report security vulnerabilities privately.

**Security → Advisories → New draft security advisory**

---

## Collaboration & Permissions

### Repository Roles

| Role         | Permissions                         |
| ------------ | ----------------------------------- |
| **Read**     | View and clone                      |
| **Triage**   | Manage issues/PRs (no write)        |
| **Write**    | Push to repo, manage issues/PRs     |
| **Maintain** | Manage repo (no sensitive settings) |
| **Admin**    | Full access including deletion      |

### Adding Collaborators

```bash
# Via CLI
gh repo collab add USERNAME --permission write

# Via web: Settings → Collaborators → Add people
```

### Teams

For organizations:

1. **Settings → Teams → New team**
2. Add members
3. Assign team to repositories

**Mention teams:**

```markdown
@org/team-name can you review this?
```

---

## CODEOWNERS

Automatically request reviews from code owners.

**Create:** `.github/CODEOWNERS`

```text
# Default owners for everything
*       @Zeus-Delacruz

# Frontend
/ui/**  @Zeus-Delacruz @frontend-team

# Backend
/agents/**  @Zeus-Delacruz @backend-team

# Infrastructure
/infra/**  @Zeus-Delacruz @devops-team
/docker-compose*.yml  @Zeus-Delacruz @devops-team
/.github/workflows/**  @Zeus-Delacruz @devops-team

# Documentation
/docs/**  @Zeus-Delacruz
*.md  @Zeus-Delacruz

# Specific files
/package.json  @Zeus-Delacruz
/requirements.txt  @Zeus-Delacruz
```

**Syntax:**

```text
# Pattern → Owner(s)

# File
README.md  @username

# Directory
/docs/  @docs-team

# Everything in directory (recursively)
/src/**  @dev-team

# File type anywhere
*.js  @js-experts

# Multiple owners
/critical/**  @user1 @user2 @team1
```

---

## GitHub Discussions

Community forum for your repository.

### Enable Discussions

**Settings → General → Features → ✅ Discussions**

### Categories

Default categories:

- 📣 **Announcements** — Updates from maintainers
- 💡 **Ideas** — Feature suggestions
- 🙏 **Q&A** — Questions from users
- 🙌 **Show and tell** — Community showcases
- 💬 **General** — Everything else

**Create custom categories:** Discussions → Categories → New category

### Using Discussions

```bash
# Via CLI
gh discussion create --title "How to setup?" --body "..."
gh discussion list
gh discussion view 123
```

**Convert issue to discussion:**

Issue → "Convert to discussion" (sidebar)

---

## Repository Templates

Create template repositories to bootstrap new projects.

### Creating a Template

1. Create repository with common structure:

   ```
   template-repo/
   ├── .github/
   │   ├── workflows/
   │   │   └── ci.yml
   │   └── ISSUE_TEMPLATE/
   ├── .gitignore
   ├── LICENSE
   ├── README.md
   └── package.json (or equivalent)
   ```

2. **Settings → General → Template repository → ✅**

### Using a Template

```bash
# Via CLI
gh repo create my-new-project --template username/template-repo

# Via web: Use template button on repo page
```

**One Mind templates:**

- `onemind-service-template` — Python FastAPI service
- `onemind-ui-template` — Next.js UI starter
- `onemind-agent-template` — Agent starter

---

## Personal Dashboard & Notifications

### Dashboard

**github.com** (when logged in)

Shows:

- Recent activity
- Repositories you're working on
- Pull requests assigned to you
- Issues assigned to you
- Discussions you're participating in

### Notifications

**github.com/notifications**

Types:

- **Participating** — You commented, were mentioned, or assigned
- **Watching** — Updates from repos you watch

**Manage watching:**

```bash
# Watch repo
gh repo watch user/repo

# Unwatch
gh repo unwatch user/repo

# List watched
gh repo list --watch
```

**Notification settings:**

Settings → Notifications

- Email frequency
- Web notifications
- Mobile notifications

### Custom Dashboard

Use GitHub CLI with scripts:

```bash
# My issues
gh issue list --assignee @me --state open

# My PRs
gh pr list --author @me --state open

# PRs needing review
gh pr list --search "review-requested:@me"

# Combine in a script
cat > ~/bin/gh-dashboard << 'EOF'
#!/bin/bash
echo "=== My Issues ==="
gh issue list --assignee @me --state open

echo -e "\n=== My PRs ==="
gh pr list --author @me --state open

echo -e "\n=== Review Requested ==="
gh pr list --search "review-requested:@me"
EOF
chmod +x ~/bin/gh-dashboard
```

---

## GitHub Mobile

**App:** [iOS](https://apps.apple.com/app/github/id1477376905) | [Android](https://play.google.com/store/apps/details?id=com.github.android)

**Features:**

- 📱 View repositories
- 💬 Comment on issues/PRs
- ✅ Approve/request changes on PRs
- 🔔 Notifications
- 📊 View Actions runs
- 🔍 Code browsing (limited)

**Use cases:**

- Quick PR approvals
- Respond to comments
- Monitor CI/CD runs
- Triage notifications

---

## One Mind GitHub Workflow

### Repository Structure

```text
Zeus-Delacruz/
├── onemind-os           # Main project (private)
├── onemind-agent        # Agent library (public)
├── onemind-ui-kit       # UI component library (public)
├── onemind-docs         # Documentation site (public)
└── onemind-templates    # Project templates (public)
```

### Workflow

1. **Issue created** — User reports bug or requests feature
2. **Label & assign** — Triage and assign to milestone
3. **Create branch** — `git checkout -b feature/voice-interface`
4. **Develop** — Commit regularly with clear messages
5. **Push** — `git push origin feature/voice-interface`
6. **Open PR** — Use PR template, link issue ("Closes #123")
7. **CI runs** — Automated tests, linting, security scans
8. **Code review** — Request review from team
9. **Address feedback** — Make changes, push updates
10. **Approve** — Reviewer approves PR
11. **Merge** — Squash merge to main
12. **Deploy** — Automated deployment triggered
13. **Close issue** — Auto-closed by "Closes #123"

### Branch Naming

| Type         | Pattern                | Example                   |
| ------------ | ---------------------- | ------------------------- |
| **Feature**  | `feature/description`  | `feature/voice-interface` |
| **Bug fix**  | `bugfix/description`   | `bugfix/memory-leak`      |
| **Hotfix**   | `hotfix/description`   | `hotfix/auth-critical`    |
| **Docs**     | `docs/description`     | `docs/api-reference`      |
| **Refactor** | `refactor/description` | `refactor/simplify-api`   |

### Commit Messages

```text
type(scope): subject

body

footer

# Examples:
feat(agent): Add voice recognition with LiveKit
fix(ui): Resolve memory leak in chat component
docs(readme): Update installation instructions
refactor(api): Simplify error handling
perf(db): Optimize vector search query
test(agent): Add unit tests for voice module
chore(deps): Update dependencies
```

### Labels

| Category     | Labels                                                       |
| ------------ | ------------------------------------------------------------ |
| **Type**     | `bug`, `enhancement`, `documentation`                        |
| **Priority** | `p0-critical`, `p1-high`, `p2-medium`, `p3-low`              |
| **Status**   | `needs-triage`, `in-progress`, `blocked`, `ready-for-review` |
| **Area**     | `agent`, `ui`, `infrastructure`, `docs`                      |
| **Meta**     | `good-first-issue`, `help-wanted`, `duplicate`, `wontfix`    |

---

## Troubleshooting

### Can't Push to Repository

**Error:** `Permission denied (publickey)`

```bash
# Check SSH key
ssh -T git@github.com

# Add SSH key if missing
ssh-keygen -t ed25519 -C "your_email@example.com"
pbcopy < ~/.ssh/id_ed25519.pub  # Copy public key
# Add to: Settings → SSH and GPG keys → New SSH key
```

### PR Shows Too Many Changes

**Issue:** Comparing wrong branches

```bash
# Ensure you branched from correct base
git log --graph --oneline main..feature-branch

# Rebase onto correct base
git rebase main
```

### CI Failing on Actions

```bash
# View logs
gh run view --log

# Re-run failed jobs
gh run rerun 123456789

# Cancel run
gh run cancel 123456789
```

### Secrets Not Available

**Issue:** Secrets not passed to fork PRs (security)

**Solution:** Maintainer must trigger workflows manually

### Can't Merge PR

**Issue:** Branch protection rules not met

**Check:**

- ✅ All required checks passed?
- ✅ Required reviews approved?
- ✅ Conversations resolved?
- ✅ Branch up to date with base?

```bash
# Update branch
git checkout feature-branch
git rebase main
git push --force-with-lease
```

---

## Quick Reference Card

### Essential Commands

```bash
# === Authentication ===
gh auth login
gh auth status

# === Repositories ===
gh repo create my-project --public
gh repo clone user/repo
gh repo view
gh repo list

# === Issues ===
gh issue create
gh issue list --assignee @me
gh issue view 123
gh issue close 123

# === Pull Requests ===
gh pr create
gh pr list --state open
gh pr checkout 456
gh pr review 456 --approve
gh pr merge 456 --squash

# === Actions ===
gh workflow list
gh run view --log
gh run watch

# === Releases ===
gh release create v1.0.0
gh release upload v1.0.0 ./dist/*

# === Secrets ===
gh secret set SECRET_NAME
gh secret list
```

### Keyboard Shortcuts (Web)

| Key      | Action                  |
| -------- | ----------------------- |
| `?`      | Show keyboard shortcuts |
| `s`, `/` | Focus search            |
| `g` `n`  | Go to notifications     |
| `g` `p`  | Go to pull requests     |
| `g` `i`  | Go to issues            |
| `t`      | File finder             |
| `l`      | Filter labels           |
| `e`      | Edit file               |
| `c`      | Create new file         |
| `.`      | Open in github.dev      |

### PR Review Shortcuts

| Key     | Action        |
| ------- | ------------- |
| `c`     | Comment       |
| `n`     | Next file     |
| `p`     | Previous file |
| `g` `d` | Go to diff    |
| `g` `f` | Go to files   |

---

## Related Documents

- [Git Mastery Guide](GIT-MASTERY-GUIDE.md)
- [Development Workflow](<../03%20Protocols%20(SOPs)/Development-Workflow.md>)
- [Code Review Guidelines](Code-Review-Guidelines.md)
- [CI/CD Pipeline](<../09%20Automation%20(Workflows)/CI-CD-Pipeline.md>)

---

**Last Updated:** January 23, 2026  
**Author:** Zeus Delacruz  
**Status:** Active Reference
