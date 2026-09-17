---
id: ui-12-20260111053926
title: "GitHub Integration Configuration"
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
ai_summary: "GitHub Integration Configuration"
---

# GitHub Integration Configuration

## Overview

GitHub serves as the **Version Control & Collaboration Layer** for One Mind. It provides:

1. Vault backup and history via Git
2. Infrastructure code management
3. Documentation versioning
4. Automation via GitHub Actions

## Repository Structure

```
One-Mind/
├── .github/
│   └── workflows/        # GitHub Actions
├── 00-Infrastructure/    # Tracked fully
├── README.md             # Public-facing
└── .gitignore            # Sensitive exclusions
```

## Git Configuration

### .gitignore Essentials

```gitignore
# Secrets
.env
.env.local
*.pem
*.key
credentials.json

# Obsidian
.obsidian/workspace.json
.obsidian/workspace-mobile.json
.obsidian/plugins/*/data.json

# System
.DS_Store
Thumbs.db

# Large media (tracked in 99-Review)
*.mp4
*.mov
*.avi
*.mp3
*.wav

# Sensitive folders
**/Private/
**/Secrets/
```

### Commit Message Convention

```
<type>(<scope>): <description>

Types:
- feat: New feature or capability
- fix: Bug fix
- docs: Documentation changes
- refactor: Code restructuring
- config: Configuration changes
- sync: Vault sync updates

Examples:
- feat(agents): Add IT agent skill for monitoring
- docs(readme): Update architecture diagram
- sync(vault): Daily backup 2024-12-23
```

## GitHub Actions

### Daily Vault Backup

```yaml
# .github/workflows/daily-backup.yml
name: Daily Vault Backup

on:
  schedule:
    - cron: '0 6 * * *'  # 6 AM UTC daily
  workflow_dispatch:      # Manual trigger

jobs:
  backup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Create backup tag
        run: |
          DATE=$(date +%Y-%m-%d)
          git tag -a "backup-$DATE" -m "Daily backup $DATE"
          git push origin "backup-$DATE"

      - name: Notify on failure
        if: failure()
        run: |
          # Send notification to IT agent
          echo "Backup failed" >> $GITHUB_STEP_SUMMARY
```

### Infrastructure Validation

```yaml
# .github/workflows/validate.yml
name: Validate Infrastructure

on:
  push:
    paths:
      - '00-Infrastructure/**'
      - 'scripts/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'

      - name: Install dependencies
        run: pip install -r requirements.txt

      - name: Validate agent configs
        run: python scripts/validate_agents.py

      - name: Check skill syntax
        run: python scripts/validate_skills.py
```

## Branch Strategy

| Branch | Purpose | Protection |
|--------|---------|------------|
| `main` | Stable vault state | Required reviews |
| `sync/*` | Obsidian sync branches | Auto-merge |
| `feature/*` | New capabilities | PR required |
| `hotfix/*` | Emergency fixes | Fast-track review |

## MCP Server Configuration

```json
{
  "name": "github",
  "type": "mcp",
  "config": {
    "command": "npx",
    "args": ["-y", "@anthropic/mcp-github"],
    "env": {
      "GITHUB_TOKEN": "${GITHUB_TOKEN}"
    }
  },
  "tools": [
    "github_create_issue",
    "github_list_issues",
    "github_create_pr",
    "github_get_file",
    "github_search_code"
  ]
}
```

## Agent GitHub Permissions

| Agent | Permissions | Scope |
|-------|-------------|-------|
| LEGACY | Read/Write | All repositories |
| HP | Read | Documentation only |
| LE | None | No GitHub access |
| GE | Read | Business docs |
| IT | Read/Write | Infrastructure repos |

## Issue Labels

| Label | Color | Purpose |
|-------|-------|---------|
| `agent:legacy` | Purple | LEGACY-created issues |
| `agent:it` | Blue | IT Agent issues |
| `priority:critical` | Red | Urgent attention |
| `type:bug` | Orange | Bug reports |
| `type:enhancement` | Green | Feature requests |
| `status:review` | Yellow | Awaiting review |

## Webhook Configuration

```yaml
# GitHub → N8N webhook
webhook_url: https://n8n.your-domain.com/webhook/github
secret: ${GITHUB_WEBHOOK_SECRET}
events:
  - push
  - issues
  - pull_request
  - workflow_run
```

## Environment Variables

```bash
# .env configuration
GITHUB_TOKEN=ghp_your_github_token
GITHUB_OWNER=your-username
GITHUB_REPO=one-mind
GITHUB_WEBHOOK_SECRET=your_webhook_secret
```

## Security Best Practices

1. **Token Scope**: Use fine-grained tokens with minimal permissions
2. **Secret Management**: Never commit tokens; use GitHub Secrets
3. **Branch Protection**: Require reviews for main branch
4. **Audit Logging**: Enable security alerts and audit log
5. **Dependency Scanning**: Enable Dependabot for security updates
