# 2. Backup & Recovery Practices

**Never lose data, recover fast**

---

## Overview

Backups are your **insurance policy**. This practice covers:

- ✅ Daily automated backups
- ✅ S3 cloud backup sync
- ✅ 30-day retention
- ✅ Disaster recovery procedures
- ✅ RTO < 30 minutes, RPO < 24 hours

**Time to implement:** 1 hour
**Maintenance:** Zero (automated)

---

## Quick Setup

```bash
# 1. Test backup script
cd /opt/onemind
./production/2-backups/backup.sh

# 2. Schedule daily backups
crontab -e

# Add this line (daily at 2 AM):
0 2 * * * /opt/onemind/production/2-backups/backup.sh >> /var/log/onemind-backup.log 2>&1

# 3. Configure S3 sync (optional)
# Edit backup.sh and set:
# ENABLE_S3_SYNC=true
# S3_BUCKET="s3://your-bucket/onemind"
```

---

## What Gets Backed Up

### Docker Volumes (Critical Data)

| Volume | Contents | Size | Critical |
|--------|----------|------|----------|
| `postgres_data` | Agent memory, sessions | ~500 MB | ✅ Critical |
| `agno_data` | Agent state | ~100 MB | ✅ Critical |
| `nats_data` | Event stream | ~200 MB | ✅ Critical |
| `n8n_data` | Workflows | ~50 MB | ✅ High |
| `kestra_data` | Workflow definitions | ~50 MB | ⚪ Medium |
| `timescale_data` | Historical events | ~1 GB | ⚪ Low |
| `uptime_data` | Monitoring history | ~100 MB | ⚪ Low |
| `caddy_data` | SSL certificates | ~10 MB | ✅ High |

**Total backup size:** ~2-5 GB (compressed)

### Configuration Files

- `docker-compose.yml` - Service definitions
- `.env` - Environment variables (encrypted)
- `Caddyfile` - Reverse proxy config
- `nats.conf` - NATS configuration

---

## Backup Script Features

The `backup.sh` script:

- ✅ Backs up all Docker volumes
- ✅ Backs up configuration files
- ✅ Creates manifest with checksums
- ✅ Compresses data (tar.gz)
- ✅ Optional S3 sync
- ✅ 30-day retention cleanup
- ✅ Discord notifications
- ✅ Error handling & logging

**Location:** `/backup/onemind/backup_YYYYMMDD_HHMMSS/`

---

## Backup Locations

### Local Storage

```
/backup/onemind/
├── backup_20260119_020000/
│   ├── postgres_data.tar.gz
│   ├── agno_data.tar.gz
│   ├── nats_data.tar.gz
│   ├── n8n_data.tar.gz
│   ├── kestra_data.tar.gz
│   ├── timescale_data.tar.gz
│   ├── uptime_data.tar.gz
│   ├── caddy_data.tar.gz
│   ├── configs.tar.gz
│   └── manifest.txt
├── backup_20260120_020000/
└── backup_20260121_020000/
```

**Retention:** 30 days (automatic cleanup)

### S3 Cloud Storage (Optional)

```
s3://onemind-backups/azure/
├── backup_20260119_020000/
├── backup_20260120_020000/
└── backup_20260121_020000/
```

**Storage class:** Glacier Instant Retrieval (~$5/month for 100 GB)
**Retention:** 90 days (S3 lifecycle policy)

---

## Setup S3 Backup

### 1. Install AWS CLI

```bash
sudo apt install -y awscli
aws --version
```

### 2. Configure AWS Credentials

```bash
aws configure

# Enter:
# AWS Access Key ID: [from 1Password]
# AWS Secret Access Key: [from 1Password]
# Default region: us-east-1
# Default output: json
```

### 3. Create S3 Bucket

```bash
# Create bucket
aws s3 mb s3://onemind-backups

# Or use AWS Console
```

### 4. Enable S3 Sync in Script

```bash
nano production/2-backups/backup.sh

# Change:
ENABLE_S3_SYNC=false

# To:
ENABLE_S3_SYNC=true
S3_BUCKET="s3://onemind-backups/azure"
```

### 5. Test S3 Sync

```bash
./production/2-backups/backup.sh

# Verify upload
aws s3 ls s3://onemind-backups/azure/ --recursive
```

---

## Disaster Recovery

### Scenario: Complete Server Failure

**Recovery steps:**

### 1. Provision New Server

```bash
# Same specs: Ubuntu 22.04, 2 vCPU, 4 GB RAM, 80 GB disk
# Install Docker:
curl -fsSL https://get.docker.com | sh
```

### 2. Download Latest Backup

```bash
# From S3
aws s3 sync s3://onemind-backups/azure/backup_latest /tmp/restore

# Or from local backup (if accessible)
scp -r OLD_SERVER:/backup/onemind/backup_latest /tmp/restore
```

### 3. Clone Repository

```bash
sudo mkdir -p /opt/onemind
sudo chown $USER:$USER /opt/onemind
cd /opt/onemind
git clone https://github.com/your-org/onemind-os.git .
```

### 4. Restore Configuration

```bash
cd /opt/onemind
tar xzf /tmp/restore/configs.tar.gz
```

### 5. Restore Docker Volumes

```bash
# Restore all volumes
for volume in postgres_data agno_data nats_data n8n_data kestra_data timescale_data uptime_data caddy_data; do
    echo "Restoring $volume..."
    docker volume create "onemind_${volume}"
    docker run --rm \
        -v "onemind_${volume}:/target" \
        -v /tmp/restore:/backup:ro \
        alpine tar xzf "/backup/${volume}.tar.gz" -C /target
    echo "✅ $volume restored"
done
```

### 6. Start Services

```bash
docker compose -f containers/cloud/docker-compose.yml up -d
```

### 7. Verify Health

```bash
# Check all services healthy
docker compose ps | grep healthy

# Test APIs
curl http://localhost:8080/health
curl http://localhost:8000/health
```

**Recovery Time:** < 30 minutes
**Data Loss:** < 24 hours (last backup)

---

## Restore Individual Service

If only one service needs restoration:

```bash
# 1. Get latest backup
LATEST_BACKUP=$(ls -t /backup/onemind/ | head -1)

# 2. Stop service
docker compose stop agno

# 3. Remove volume
docker volume rm onemind_agno_data

# 4. Recreate volume
docker volume create onemind_agno_data

# 5. Restore data
docker run --rm \
    -v onemind_agno_data:/target \
    -v "/backup/onemind/$LATEST_BACKUP":/backup:ro \
    alpine tar xzf /backup/agno_data.tar.gz -C /target

# 6. Start service
docker compose start agno

# 7. Verify
docker compose ps | grep agno
curl http://localhost:8080/health
```

---

## Testing Disaster Recovery

**Practice disaster recovery quarterly:**

### Test Restore on Temporary Service

```bash
# 1. Stop test service (non-critical)
docker compose stop portainer

# 2. Delete volume
docker volume rm onemind_portainer_data

# 3. Restore from backup
LATEST_BACKUP=$(ls -t /backup/onemind/ | head -1)
docker volume create onemind_portainer_data
docker run --rm \
    -v onemind_portainer_data:/target \
    -v "/backup/onemind/$LATEST_BACKUP":/backup:ro \
    alpine tar xzf /backup/portainer_data.tar.gz -C /target

# 4. Start service
docker compose start portainer

# 5. Verify works
open http://YOUR_IP:9000
```

**✅ If portainer works, your backup/restore process is verified!**

---

## Monitoring Backups

### Check Last Backup

```bash
# View latest backup
ls -lht /backup/onemind/ | head -5

# Check backup size
du -sh /backup/onemind/backup_*/

# View manifest
cat /backup/onemind/backup_*/manifest.txt
```

### Check Backup Logs

```bash
# View today's backup log
tail -100 /var/log/onemind-backup.log

# Check for errors
grep -i error /var/log/onemind-backup.log

# Check for failures
grep -i fail /var/log/onemind-backup.log
```

### Verify S3 Sync (if enabled)

```bash
# List S3 backups
aws s3 ls s3://onemind-backups/azure/ --recursive

# Compare local vs S3
aws s3 sync /backup/onemind s3://onemind-backups/azure --dryrun
```

---

## Troubleshooting

### Backup Script Fails

**Check disk space:**
```bash
df -h /backup
```

**Check Docker access:**
```bash
docker volume ls
```

**Run backup manually (debug mode):**
```bash
bash -x ./production/2-backups/backup.sh
```

### S3 Sync Fails

**Check AWS credentials:**
```bash
aws s3 ls
```

**Test S3 access:**
```bash
aws s3 ls s3://onemind-backups/
```

**Check IAM permissions:**
- Need: `s3:PutObject`, `s3:GetObject`, `s3:ListBucket`

### Restore Fails

**Check backup integrity:**
```bash
tar -tzf /backup/onemind/backup_latest/postgres_data.tar.gz | head
```

**Check volume exists:**
```bash
docker volume ls | grep onemind
```

**Check permissions:**
```bash
ls -l /backup/onemind/
```

---

## Backup Best Practices

### ✅ DO

- ✅ Test backups regularly (quarterly)
- ✅ Store backups off-site (S3)
- ✅ Monitor backup success/failure
- ✅ Document restore procedures
- ✅ Practice disaster recovery
- ✅ Encrypt backups containing secrets
- ✅ Verify backup integrity

### ❌ DON'T

- ❌ Only have one backup location
- ❌ Never test restores
- ❌ Delete old backups immediately
- ❌ Store backups on same server only
- ❌ Ignore backup failures
- ❌ Commit .env to git
- ❌ Trust backups blindly

---

## Maintenance Schedule

### Daily (Automated)
- ✅ 2:00 AM - Backup runs
- ✅ 2:30 AM - S3 sync (if enabled)

### Weekly (5 minutes)
- Review backup logs
- Verify latest backup exists
- Check S3 sync status

### Monthly (15 minutes)
- Test restore procedure
- Review disk usage
- Clean up old backups manually (if needed)

### Quarterly (30 minutes)
- **Full disaster recovery drill**
- Test restore on new VM
- Update documentation
- Review retention policy

---

## Next Steps

After setting up backups:

1. **Configure Monitoring** → [../3-monitoring/](../3-monitoring/README.md)
2. **Setup Zero-Downtime Deployment** → [../4-deployment/](../4-deployment/README.md)
3. **Enable Automation** → [../5-automation/](../5-automation/README.md)

---

## Resources

- AWS S3 Backup Guide
- Docker Volume Backup Strategies
- Disaster Recovery Best Practices

---

**💾 Your data is now protected! Sleep well knowing you can recover from anything.**
