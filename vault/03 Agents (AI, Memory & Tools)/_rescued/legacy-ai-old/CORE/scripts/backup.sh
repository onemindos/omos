#!/bin/bash
# Backup Letta Database
# Backs up Letta's PostgreSQL database with Letta-specific memory data

set -e

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/letta_backup_${TIMESTAMP}.sql"

echo "💾 Backing up Letta database..."

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup PostgreSQL database
docker-compose exec -T postgres pg_dump -U letta letta > "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"

echo "✅ Backup saved to: ${BACKUP_FILE}.gz"

# Keep only last 7 backups
ls -t "${BACKUP_DIR}"/letta_backup_*.sql.gz | tail -n +8 | xargs -r rm

echo "🧹 Cleaned old backups (keeping last 7)"
