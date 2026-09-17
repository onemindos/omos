---
id: ui-12-20260111053926
title: "Backup and Restore Guide"
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
ai_summary: "Backup and Restore Guide"
---

# Backup and Restore Guide

## Overview

Your data is your most valuable asset. This guide ensures you never lose:
- n8n workflows and credentials
- Monitoring configurations
- Custom dashboards

## What Gets Backed Up

| Data | Location | Backup Method |
|------|----------|---------------|
| n8n workflows | PostgreSQL | Database dump |
| n8n credentials | PostgreSQL (encrypted) | Database dump |
| n8n execution history | PostgreSQL | Database dump |
| Grafana dashboards | Grafana volume | File copy |
| Uptime Kuma config | SQLite volume | File copy |
| Environment secrets | `.env` file | File copy |

## Automatic Backups

### Enable Daily Backups
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * cd /path/to/one-mind-infrastructure && ./scripts/backup.sh >> /var/log/onemind-backup.log 2>&1
```

### Backup to S3
Configure in `.env`:
```bash
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
S3_BUCKET=your-backup-bucket
S3_REGION=us-east-1
```

## Manual Backup

```bash
./scripts/backup.sh
```

Output:
```
backups/
├── n8n_db_20241219_120000.sql.gz
├── n8n_workflows_20241219_120000.json
├── grafana_20241219_120000.tar.gz
├── uptime_kuma_20241219_120000.tar.gz
└── env_20241219_120000.backup
```

## Restore Procedures

### Full Restore
```bash
./scripts/restore.sh 20241219_120000
```

### Restore n8n Database Only
```bash
# Stop n8n services
docker compose -f apps/automation/n8n/docker-compose.yml stop n8n-main n8n-worker-1 n8n-worker-2

# Restore
gunzip -c backups/n8n_db_YYYYMMDD.sql.gz | docker exec -i n8n-postgres psql -U n8n -d n8n

# Restart
docker compose -f apps/automation/n8n/docker-compose.yml up -d
```

### Restore Specific Workflows
1. Go to n8n UI
2. Settings → Import from File
3. Select `n8n_workflows_YYYYMMDD.json`

## Critical: The Encryption Key

**`N8N_ENCRYPTION_KEY` is the most important secret!**

Without it:
- All stored credentials (API keys, passwords) are **permanently lost**
- Workflows will exist but won't work
- You'll need to re-enter every credential

### Where to Store It
1. Password manager (1Password, Bitwarden)
2. Encrypted note in secure location
3. Printed and stored in safe (disaster recovery)

### Recovery Without Key
If you lose the key, you must:
1. Restore database (workflows work)
2. Manually re-add all credentials
3. Test all workflows

## Multi-Site Sync

For multiple instances with synced configs:

### Git-Based Workflow Sync
```bash
# Export workflows to JSON
curl -X GET "https://n8n.site1.com/api/v1/workflows" \
  -H "X-N8N-API-KEY: $API_KEY" > workflows/all.json

# Commit to git
git add workflows/
git commit -m "Sync workflows $(date +%Y%m%d)"
git push

# On other sites, pull and import
git pull
# Import via n8n API or UI
```

### Shared Credentials
Credentials cannot be exported (security). For multi-site:
1. Use environment variables where possible
2. Document credentials in password manager
3. Apply manually to each instance

## Disaster Recovery

### Complete Server Loss

1. **Spin up new VM**
2. **Install Docker** (see SETUP.md)
3. **Clone repo**
   ```bash
   git clone your-repo
   cd one-mind-infrastructure
   ```
4. **Restore .env** from secure backup
5. **Deploy infrastructure**
   ```bash
   ./scripts/deploy.sh
   ```
6. **Restore database**
   ```bash
   # Download from S3 or copy from backup location
   ./scripts/restore.sh latest
   ```

### Recovery Time Objective
With proper backups: **< 1 hour**


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
