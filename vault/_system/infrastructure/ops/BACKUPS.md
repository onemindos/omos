# OneMind OS — Backup & Recovery

> Last Updated: 2026-03-14
> Covers all backup layers for both VPSes

---

## Backup Overview

| Layer | Main VPS | Sentinel VPS | Frequency | Retention |
|-------|----------|--------------|-----------|-----------|
| **DO Auto Backups** | Enabled (weekly) | Enabled (weekly) | Weekly | 4 snapshots |
| **PostgreSQL Dumps** | `/home/onemind/backups/postgres/` | n/a (no DB) | Daily 2:00 UTC | 7 days |
| **Docker Volume Backups** | `/home/onemind/backups/docker-volumes/` | n/a (no Docker) | Daily 2:30 UTC | 7 days |
| **Off-site (DO Spaces)** | `s3://onemind-backups/onemind-main/` | n/a | Daily 3:00 UTC | Mirrors local |
| **Git (code/config)** | `codex-watch` + `codex-pull.timer` | Git pull cron (1 min) | Continuous | Git history |
| **Workspace backup** | `workspace-backup.sh` (every 15 min) | n/a | Every 15 min | Git history |

**Recovery targets**: RTO < 30 min, RPO < 24 hours (database), RPO < 1 min (code/config)

---

## DigitalOcean Auto Backups

Weekly full-disk snapshots managed by DigitalOcean. ~20% of droplet cost.

| Droplet | ID | Backups | Cost |
|---------|----|---------|------|
| `onemind-main` | 558147996 | Enabled | ~$16.80/mo |
| `onemind-sentinel` | 558205817 | Enabled | ~$2.40/mo |

```bash
# List backups for main
doctl compute droplet backups 558147996

# List backups for sentinel
doctl compute droplet backups 558205817

# Restore from backup (creates action)
doctl compute droplet-action restore <droplet-id> --image-id <backup-id>
```

---

## PostgreSQL Dumps (Main VPS only)

Dumps the `onemind` database from the OmniDB Docker container.

**Script**: `/home/onemind/backups/backup-postgres.sh`
**Schedule**: Daily at 2:00 AM UTC
**Output**: `/home/onemind/backups/postgres/onemind_YYYY-MM-DD_HHMM.sql.gz`
**Retention**: 7 days

```bash
# Manual run
sudo -u onemind /home/onemind/backups/backup-postgres.sh

# Check latest backup
ls -lh /home/onemind/backups/postgres/

# Restore
gunzip -c /home/onemind/backups/postgres/onemind_2026-03-14_0200.sql.gz \
  | docker exec -i onemind-omnidb psql -U onemind -d onemind
```

### Database Details

| Setting | Value |
|---------|-------|
| Container | `onemind-omnidb` |
| Engine | PostgreSQL 16 + TimescaleDB + pgvector |
| User | `onemind` |
| Database | `onemind` |
| Volume | `onemind_omnidb_data` |

---

## Docker Volume Backups (Main VPS only)

Backs up NATS JetStream data and OmniDB data as compressed tarballs.

**Script**: `/home/onemind/backups/backup-docker-volumes.sh`
**Schedule**: Daily at 2:30 AM UTC
**Output**: `/home/onemind/backups/docker-volumes/<volume>_YYYY-MM-DD_HHMM.tar.gz`
**Retention**: 7 days

| Volume | Contents | Typical Size |
|--------|----------|-------------|
| `onemind_nats_data` | NATS JetStream persistence | ~2 KB |
| `onemind_omnidb_data` | PostgreSQL data files | ~9 MB |

```bash
# Manual run (runs as root for Docker access)
/home/onemind/backups/backup-docker-volumes.sh

# Restore a volume (stop container first)
docker compose stop omnidb
docker run --rm -v onemind_omnidb_data:/target -v /home/onemind/backups/docker-volumes:/backup \
  alpine sh -c "rm -rf /target/* && tar xzf /backup/onemind_omnidb_data_2026-03-14_0230.tar.gz -C /target"
docker compose start omnidb
```

---

## Off-site Backup (DO Spaces)

All local backups are synced to DigitalOcean Spaces (S3-compatible object storage).

**Script**: `/home/onemind/backups/backup-offsite.sh`
**Schedule**: Daily at 3:00 AM UTC (runs after local backups complete)
**Bucket**: `onemind-backups` (nyc3 region)
**Path**: `s3://onemind-backups/onemind-main/`

| Setting | Value |
|---------|-------|
| Spaces endpoint | `nyc3.digitaloceanspaces.com` |
| Bucket | `onemind-backups` |
| Access key | `DO00URA2K3WJMZDM9FMA` |
| Config file | `/home/onemind/.s3cfg` |
| Tool | `s3cmd` |

```bash
# Manual sync
sudo -u onemind /home/onemind/backups/backup-offsite.sh

# List bucket contents
sudo -u onemind s3cmd ls s3://onemind-backups/onemind-main/ --recursive

# Download a backup from Spaces
sudo -u onemind s3cmd get s3://onemind-backups/onemind-main/postgres/onemind_2026-03-14_0200.sql.gz /tmp/
```

---

## Git-based Backups (Both VPSes)

Code and configuration are protected by continuous git sync.

### Main VPS

| Service | Purpose | Direction |
|---------|---------|-----------|
| `codex-watch.service` | Auto-commits changes in `codex/` and `workspaces/` | VPS -> GitHub |
| `codex-pull.timer` | Pulls from GitHub every 1 minute | GitHub -> VPS |
| `workspace-backup.sh` | Backs up workspace state every 15 min (onemind crontab) | VPS -> GitHub |

### Sentinel VPS

| Mechanism | Purpose |
|-----------|---------|
| `/etc/cron.d/sentinel-git-pull` | Pulls monorepo from GitHub every 1 minute |

---

## Cron Schedule (Main VPS)

**File**: `/etc/cron.d/onemind-backups`

```cron
# PostgreSQL dump - daily at 2:00 AM UTC
0 2 * * * onemind /home/onemind/backups/backup-postgres.sh >> /home/onemind/backups/postgres.log 2>&1

# Docker volume backup - daily at 2:30 AM UTC
30 2 * * * root /home/onemind/backups/backup-docker-volumes.sh >> /home/onemind/backups/docker-volumes.log 2>&1

# Off-site sync to DO Spaces - daily at 3:00 AM UTC
0 3 * * * onemind /home/onemind/backups/backup-offsite.sh >> /home/onemind/backups/offsite.log 2>&1
```

---

## Disaster Recovery

### Scenario 1: Main VPS dies completely

1. **Restore from DO backup**: `doctl compute droplet-action restore 558147996 --image-id <backup-id>`
2. **Or rebuild from scratch**: Run `openclaw/deploy/scripts/setup.sh` on a new droplet
3. Git pull restores all code/config in under 1 minute
4. Restore DB from Spaces: download latest `.sql.gz` and pipe to psql
5. Sentinel keeps Legacy reachable during recovery

### Scenario 2: Database corruption

1. Stop OmniDB: `cd /opt/onemind && docker compose stop omnidb`
2. Restore latest dump: `gunzip -c backups/postgres/onemind_*.sql.gz | docker exec -i onemind-omnidb psql -U onemind -d onemind`
3. Or restore full volume from tarball backup
4. Start OmniDB: `docker compose start omnidb`

### Scenario 3: Sentinel VPS dies

1. Non-critical - it is a rescue node
2. Rebuild: Create new droplet, run `openclaw/deploy/scripts/sentinel-setup.sh`
3. Install OpenClaw, restore config from git

---

## Backup File Locations (Main VPS)

```
/home/onemind/backups/
    backup-postgres.sh          <- PostgreSQL dump script
    backup-docker-volumes.sh    <- Docker volume backup script
    backup-offsite.sh           <- DO Spaces sync script
    postgres.log                <- Postgres backup log
    docker-volumes.log          <- Volume backup log
    offsite.log                 <- Spaces sync log
    postgres/
        onemind_YYYY-MM-DD_HHMM.sql.gz
    docker-volumes/
        onemind_nats_data_YYYY-MM-DD_HHMM.tar.gz
        onemind_omnidb_data_YYYY-MM-DD_HHMM.tar.gz
```

---

## Off-site Bucket Structure

```
s3://onemind-backups/
    onemind-main/
        backup-postgres.sh
        backup-docker-volumes.sh
        backup-offsite.sh
        postgres/
            onemind_YYYY-MM-DD_HHMM.sql.gz
        docker-volumes/
            onemind_nats_data_YYYY-MM-DD_HHMM.tar.gz
            onemind_omnidb_data_YYYY-MM-DD_HHMM.tar.gz
```
