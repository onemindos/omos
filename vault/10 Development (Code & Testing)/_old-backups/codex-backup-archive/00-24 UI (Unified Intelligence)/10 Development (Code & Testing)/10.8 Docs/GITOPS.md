---
id: ui-12-20260111053926
title: "GitOps Workflow Guide"
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
ai_summary: "GitOps Workflow Guide"
---

# GitOps Workflow Guide

## Overview

This repository uses GitHub Actions for GitOps automation:

```
VS Code + Copilot
      │
      ▼
   git push
      │
      ▼
GitHub Actions
      │
      ├──► deploy.yml (on push)
      ├──► backup.yml (daily 2 AM)
      └──► update.yml (weekly Sunday 3 AM)
      │
      ▼
   Your VM
```

## How It Works

### 1. Deploy on Push
When you push to `main`:
1. GitHub Actions triggers `deploy.yml`
2. SSHs into your VM
3. Runs `git pull` + `./scripts/deploy.sh`
4. Services are updated

### 2. Daily Backup
Every day at 2 AM UTC:
1. GitHub Actions triggers `backup.yml`
2. SSHs into VM, runs backup script
3. Downloads n8n workflow exports
4. Commits to `/workflows/` folder
5. Pushes to GitHub

### 3. Weekly Update
Every Sunday at 3 AM UTC:
1. GitHub Actions triggers `update.yml`
2. Creates backup first
3. Pulls latest container images
4. Performs rolling update
5. Runs health checks

## Setup GitHub Secrets

Add these secrets in GitHub → Settings → Secrets:

| Secret | Description |
|--------|-------------|
| `VM_HOST` | Your VM IP (e.g., `3.215.29.128`) |
| `VM_USER` | SSH user (e.g., `ubuntu`) |
| `VM_SSH_KEY` | Private SSH key (contents of .pem file) |
| `S3_BUCKET` | (Optional) S3 bucket for database backups |

### How to Add SSH Key

1. Go to GitHub repo → Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `VM_SSH_KEY`
4. Value: Paste entire contents of your `.pem` file

## Your Workflow

### Making Changes

```bash
# 1. Edit in VS Code (Copilot helps)
# 2. Test locally if possible
# 3. Commit and push

git add .
git commit -m "Update n8n config"
git push

# 4. GitHub Actions deploys automatically
# 5. Check Actions tab for status
```

### Manual Deploy

If needed, trigger manually:
1. Go to GitHub → Actions
2. Select "Deploy to VM"
3. Click "Run workflow"

### Rollback

If something breaks:

```bash
# Option 1: Revert commit
git revert HEAD
git push

# Option 2: Restore from backup
ssh your-vm
cd ~/one-mind-infrastructure
./scripts/restore.sh
```

## Monitoring Deployments

### Check GitHub Actions
- Go to repo → Actions tab
- See all workflow runs
- Click for detailed logs

### Check VM Status
```bash
ssh your-vm
cd ~/one-mind-infrastructure
./scripts/health-check.sh
```

## Workflow Files

| File | Trigger | Action |
|------|---------|--------|
| `.github/workflows/deploy.yml` | Push to main | Deploy changes |
| `.github/workflows/backup.yml` | Daily 2 AM | Backup + commit workflows |
| `.github/workflows/update.yml` | Weekly Sunday | Update containers |

## Customizing

### Change Backup Schedule
Edit `.github/workflows/backup.yml`:
```yaml
schedule:
  - cron: '0 2 * * *'  # Change this cron expression
```

### Change Update Schedule
Edit `.github/workflows/update.yml`:
```yaml
schedule:
  - cron: '0 3 * * 0'  # Currently: Sunday 3 AM
```

### Skip Deploy for Certain Files
Already configured - docs and workflow exports don't trigger deploy:
```yaml
paths-ignore:
  - 'workflows/**'
  - 'docs/**'
  - '*.md'
```

## Troubleshooting

### Deploy Failed
1. Check Actions log for error
2. SSH into VM manually
3. Check `docker logs <container>`
4. Fix issue, push fix

### Backup Failed
1. Check if n8n is running
2. Verify N8N_API_KEY is set
3. Check SSH connectivity

### SSH Connection Failed
1. Verify VM_HOST is correct
2. Check VM_SSH_KEY is complete (including headers)
3. Verify VM firewall allows SSH (port 22)
