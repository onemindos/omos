# OneMind OS - Backup Strategy

> Last Updated: 2026-01-30
> Complete backup configuration for Azure VM and databases

---

## Backup Overview

| Layer | Solution | Frequency | Retention | Status |
|-------|----------|-----------|-----------|--------|
| **VM (Full)** | Azure Backup | Daily 3pm UTC | 30 days | Active |
| **VM Snapshot** | Manual | Before changes | Manual | Available |
| **PostgreSQL** | pg_dump | Daily 3am UTC | 30 days | Active |
| **TimescaleDB** | pg_dump | Daily 3am UTC | 30 days | Active |

---

## Azure Backup (VM-Level)

### Recovery Services Vault

| Property | Value |
|----------|-------|
| **Vault Name** | onemind-backup-vault |
| **Resource Group** | ONEMIND-PRODUCTION |
| **Location** | East US |
| **Policy** | DefaultPolicy |
| **Schedule** | Daily at 3:00 PM UTC (10:00 AM EST) |
| **Retention** | 30 days |

### Check Backup Status

```bash
# List backup jobs
az backup job list \
  --resource-group ONEMIND-PRODUCTION \
  --vault-name onemind-backup-vault \
  --output table

# Check VM backup status
az backup item show \
  --resource-group ONEMIND-PRODUCTION \
  --vault-name onemind-backup-vault \
  --container-name "IaasVMContainer;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --name "VM;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --backup-management-type AzureIaasVM \
  --query '{lastBackup:properties.lastBackupTime, health:properties.healthStatus}' \
  -o table
```

### Trigger Manual Backup

```bash
az backup protection backup-now \
  --resource-group ONEMIND-PRODUCTION \
  --vault-name onemind-backup-vault \
  --container-name "IaasVMContainer;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --item-name "VM;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --backup-management-type AzureIaasVM \
  --retain-until $(date -d '+30 days' +%d-%m-%Y)
```

### Restore from Backup

```bash
# List available recovery points
az backup recoverypoint list \
  --resource-group ONEMIND-PRODUCTION \
  --vault-name onemind-backup-vault \
  --container-name "IaasVMContainer;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --item-name "VM;iaasvmcontainerv2;ONEMIND-PRODUCTION;onemind-vm" \
  --backup-management-type AzureIaasVM \
  --output table

# Restore VM (creates new VM or replaces existing)
# Use Azure Portal for guided restore: portal.azure.com
```

---

## Disk Snapshots

### Current Snapshot

| Property | Value |
|----------|-------|
| **Name** | onemind-snapshot-20260130 |
| **Source Disk** | onemind-vm_OsDisk_1_... |
| **Size** | 30 GB |
| **Created** | 2026-01-30 |

### Create New Snapshot

```bash
# Before major changes, create a snapshot
az snapshot create \
  --resource-group ONEMIND-PRODUCTION \
  --name onemind-snapshot-$(date +%Y%m%d) \
  --source onemind-vm_OsDisk_1_9673580f9734496a96322a908d72db22

# List snapshots
az snapshot list -g ONEMIND-PRODUCTION -o table
```

### Restore from Snapshot

```bash
# Create new disk from snapshot
az disk create \
  --resource-group ONEMIND-PRODUCTION \
  --name onemind-restored-disk \
  --source onemind-snapshot-20260130

# Then swap the OS disk on the VM (requires VM stop)
```

### Delete Old Snapshots

```bash
# List snapshots older than 7 days
az snapshot list -g ONEMIND-PRODUCTION \
  --query "[?timeCreated<='$(date -d '-7 days' -Iseconds)'].name" -o tsv

# Delete specific snapshot
az snapshot delete -g ONEMIND-PRODUCTION -n <snapshot-name>
```

---

## Database Backups

### Configuration

| Database | Container | Backup Location | Schedule |
|----------|-----------|-----------------|----------|
| PostgreSQL | onemind-postgres | /home/zeus/backups | Daily 3am UTC |
| TimescaleDB | onemind-timescaledb | /home/zeus/backups | Daily 3am UTC |

### Backup Script

Location: `/home/zeus/backup-databases.sh`

```bash
#!/bin/bash
# Runs daily at 3am UTC via cron
# Backs up PostgreSQL and TimescaleDB
# Keeps 30 days of backups

~/backup-databases.sh
```

### Cron Job

```bash
# View current cron
crontab -l

# Output:
0 3 * * * /home/zeus/backup-databases.sh >> /var/log/onemind-backup.log 2>&1
```

### Manual Backup

```bash
# SSH to VM
ssh zeus@100.102.21.44

# Run backup
~/backup-databases.sh

# Check backups
ls -lh ~/backups/
```

### Restore Database

```bash
# SSH to VM
ssh zeus@100.102.21.44

# List available backups
ls -lh ~/backups/

# Restore PostgreSQL
gunzip -c ~/backups/postgres_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i onemind-postgres psql -U onemind -d onemind

# Restore TimescaleDB
gunzip -c ~/backups/timescale_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i onemind-timescaledb psql -U onemind -d events
```

---

## Azure Storage Account

### Details

| Property | Value |
|----------|-------|
| **Account Name** | onemindbackups |
| **Resource Group** | ONEMIND-PRODUCTION |
| **Location** | East US |
| **SKU** | Standard_LRS |
| **Containers** | database-backups, timescale-backups |

### Future Enhancement: Cloud Backup

To upload database backups to Azure Blob Storage:

```bash
# Install azcopy
wget -O azcopy.tar.gz https://aka.ms/downloadazcopy-v10-linux
tar -xzf azcopy.tar.gz --strip-components=1
sudo mv azcopy /usr/local/bin/

# Upload backup
azcopy copy ~/backups/postgres_*.sql.gz \
  "https://onemindbackups.blob.core.windows.net/database-backups/?<SAS-TOKEN>"
```

---

## Backup Schedule Summary

| Time (UTC) | Action |
|------------|--------|
| 03:00 | Database backup (PostgreSQL + TimescaleDB) |
| 15:00 | Azure VM backup |

---

## Recovery Procedures

### Scenario 1: Accidental Data Deletion

1. Restore from database backup:
   ```bash
   gunzip -c ~/backups/postgres_LATEST.sql.gz | \
     docker exec -i onemind-postgres psql -U onemind -d onemind
   ```

### Scenario 2: VM Corruption/Failure

1. Go to Azure Portal → Recovery Services vault → onemind-backup-vault
2. Select "Backup items" → "Azure Virtual Machine"
3. Select onemind-vm → "Restore VM"
4. Choose recovery point and restore options

### Scenario 3: Quick Rollback (Before Major Change)

1. Create snapshot before change:
   ```bash
   az snapshot create -g ONEMIND-PRODUCTION \
     -n onemind-pre-change-$(date +%Y%m%d) \
     --source onemind-vm_OsDisk_1_9673580f9734496a96322a908d72db22
   ```
2. If change fails, restore from snapshot

---

## Monitoring

### Check Backup Logs

```bash
# SSH to VM
ssh zeus@100.102.21.44

# View backup logs
tail -50 /var/log/onemind-backup.log

# Check last backup
ls -lt ~/backups/ | head -5
```

### Azure Backup Alerts

Set up alerts in Azure Portal:
1. Go to Recovery Services vault → onemind-backup-vault
2. Alerts + Metrics → Create alert rule
3. Configure for backup failures

---

## Cost Estimate

| Component | Monthly Cost |
|-----------|--------------|
| Azure Backup (30GB VM, 30 days) | ~$12 |
| Snapshots (30GB × 2) | ~$3 |
| Storage Account (minimal) | ~$1 |
| **Total** | **~$16/month** |

---

## Quick Reference

```bash
# Check all backup status
az backup job list -g ONEMIND-PRODUCTION --vault-name onemind-backup-vault -o table

# List snapshots
az snapshot list -g ONEMIND-PRODUCTION -o table

# SSH and check database backups
ssh zeus@100.102.21.44 "ls -lh ~/backups/"

# Run manual database backup
ssh zeus@100.102.21.44 "~/backup-databases.sh"

# Create manual snapshot
az snapshot create -g ONEMIND-PRODUCTION -n onemind-snap-$(date +%Y%m%d) \
  --source onemind-vm_OsDisk_1_9673580f9734496a96322a908d72db22
```
