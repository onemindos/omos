#!/bin/bash
################################################################################
# LEGACY AI COMPREHENSIVE BACKUP SYSTEM
# Purpose: Backup all critical data including databases, volumes, configs
# Usage: ./scripts/maintenance/backup-all.sh [backup_dir]
# Schedule: Run daily via cron at 2 AM
################################################################################

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BACKUP_BASE_DIR="${1:-/backups}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="${BACKUP_BASE_DIR}/${TIMESTAMP}"
LOG_FILE="/var/log/legacy-ai-backup-${TIMESTAMP}.log"
RETENTION_DAYS=7

echo "💾 STARTING COMPREHENSIVE BACKUP - $(date)" | tee -a "$LOG_FILE"
echo "========================================" | tee -a "$LOG_FILE"
echo "Backup location: $BACKUP_DIR" | tee -a "$LOG_FILE"

# Create backup directory
mkdir -p "$BACKUP_DIR"/{volumes,databases,configs,stacks}

BACKUP_ERRORS=0
BACKUP_SUCCESS=0

################################################################################
# 1. BACKUP DOCKER VOLUMES
################################################################################
echo -e "\n${BLUE}📦 Backing up Docker volumes...${NC}" | tee -a "$LOG_FILE"

VOLUMES=$(docker volume ls -q)
VOLUME_COUNT=$(echo "$VOLUMES" | wc -w)

echo "Found $VOLUME_COUNT volumes to backup" | tee -a "$LOG_FILE"

for VOLUME in $VOLUMES; do
    echo "  → Backing up volume: $VOLUME" | tee -a "$LOG_FILE"
    
    # Backup volume using temporary alpine container
    if docker run --rm \
        -v "${VOLUME}:/data:ro" \
        -v "${BACKUP_DIR}/volumes:/backup" \
        alpine \
        tar czf "/backup/${VOLUME}.tar.gz" -C /data . >> "$LOG_FILE" 2>&1; then
        
        BACKUP_SIZE=$(du -h "${BACKUP_DIR}/volumes/${VOLUME}.tar.gz" | cut -f1)
        echo -e "    ${GREEN}✅ ${VOLUME}${NC} (${BACKUP_SIZE})" | tee -a "$LOG_FILE"
        ((BACKUP_SUCCESS++))
    else
        echo -e "    ${RED}❌ Failed to backup ${VOLUME}${NC}" | tee -a "$LOG_FILE"
        ((BACKUP_ERRORS++))
    fi
done

################################################################################
# 2. BACKUP POSTGRESQL DATABASES
################################################################################
echo -e "\n${BLUE}🐘 Backing up PostgreSQL databases...${NC}" | tee -a "$LOG_FILE"

# n8n PostgreSQL database
echo "  → Backing up n8n-postgres..." | tee -a "$LOG_FILE"
if docker exec n8n-postgres pg_dump -U n8n -d n8n -F c -f /tmp/n8n-backup.dump >> "$LOG_FILE" 2>&1; then
    docker cp n8n-postgres:/tmp/n8n-backup.dump "${BACKUP_DIR}/databases/n8n-postgres-${TIMESTAMP}.dump"
    docker exec n8n-postgres rm /tmp/n8n-backup.dump
    
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/databases/n8n-postgres-${TIMESTAMP}.dump" | cut -f1)
    echo -e "    ${GREEN}✅ n8n-postgres${NC} (${BACKUP_SIZE})" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
else
    echo -e "    ${RED}❌ Failed to backup n8n-postgres${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_ERRORS++))
fi

# Dify PostgreSQL database
echo "  → Backing up dify-postgres..." | tee -a "$LOG_FILE"
if docker exec dify-postgres pg_dump -U postgres -F c -f /tmp/dify-backup.dump >> "$LOG_FILE" 2>&1; then
    docker cp dify-postgres:/tmp/dify-backup.dump "${BACKUP_DIR}/databases/dify-postgres-${TIMESTAMP}.dump"
    docker exec dify-postgres rm /tmp/dify-backup.dump
    
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/databases/dify-postgres-${TIMESTAMP}.dump" | cut -f1)
    echo -e "    ${GREEN}✅ dify-postgres${NC} (${BACKUP_SIZE})" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
else
    echo -e "    ${RED}❌ Failed to backup dify-postgres${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_ERRORS++))
fi

################################################################################
# 3. BACKUP REDIS DATABASES
################################################################################
echo -e "\n${BLUE}🔴 Backing up Redis databases...${NC}" | tee -a "$LOG_FILE"

# n8n Redis
echo "  → Backing up n8n-redis..." | tee -a "$LOG_FILE"
if docker exec n8n-redis redis-cli --rdb /tmp/n8n-redis.rdb save >> "$LOG_FILE" 2>&1; then
    sleep 2  # Wait for save to complete
    docker cp n8n-redis:/data/dump.rdb "${BACKUP_DIR}/databases/n8n-redis-${TIMESTAMP}.rdb"
    
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/databases/n8n-redis-${TIMESTAMP}.rdb" | cut -f1)
    echo -e "    ${GREEN}✅ n8n-redis${NC} (${BACKUP_SIZE})" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
else
    echo -e "    ${RED}❌ Failed to backup n8n-redis${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_ERRORS++))
fi

# Dify Redis
echo "  → Backing up dify-redis..." | tee -a "$LOG_FILE"
if docker exec dify-redis redis-cli --rdb /tmp/dify-redis.rdb save >> "$LOG_FILE" 2>&1; then
    sleep 2  # Wait for save to complete
    docker cp dify-redis:/data/dump.rdb "${BACKUP_DIR}/databases/dify-redis-${TIMESTAMP}.rdb"
    
    BACKUP_SIZE=$(du -h "${BACKUP_DIR}/databases/dify-redis-${TIMESTAMP}.rdb" | cut -f1)
    echo -e "    ${GREEN}✅ dify-redis${NC} (${BACKUP_SIZE})" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
else
    echo -e "    ${RED}❌ Failed to backup dify-redis${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_ERRORS++))
fi

################################################################################
# 4. BACKUP ENVIRONMENT FILES (ENCRYPTED)
################################################################################
echo -e "\n${BLUE}🔐 Backing up environment files (encrypted)...${NC}" | tee -a "$LOG_FILE"

# Find all .env files in stacks directory
ENV_FILES=$(find stacks/ -name ".env" -type f 2>/dev/null || echo "")

if [ -n "$ENV_FILES" ]; then
    ENV_COUNT=$(echo "$ENV_FILES" | wc -l)
    echo "Found $ENV_COUNT .env files" | tee -a "$LOG_FILE"
    
    # Check if GPG is available
    if command -v gpg >/dev/null 2>&1; then
        echo "$ENV_FILES" | while read -r ENVFILE; do
            APP_NAME=$(basename "$(dirname "$ENVFILE")")
            echo "  → Encrypting: $ENVFILE" | tee -a "$LOG_FILE"
            
            # Encrypt with passphrase (use environment variable or prompt)
            if [ -n "$BACKUP_ENCRYPTION_KEY" ]; then
                # Use provided key
                echo "$BACKUP_ENCRYPTION_KEY" | gpg --batch --yes --passphrase-fd 0 \
                    --symmetric --cipher-algo AES256 \
                    --output "${BACKUP_DIR}/configs/${APP_NAME}.env.gpg" \
                    "$ENVFILE" >> "$LOG_FILE" 2>&1
            else
                # No encryption key provided - copy unencrypted with warning
                echo -e "    ${YELLOW}⚠️  No encryption key set - storing unencrypted${NC}" | tee -a "$LOG_FILE"
                cp "$ENVFILE" "${BACKUP_DIR}/configs/${APP_NAME}.env"
            fi
            
            if [ $? -eq 0 ]; then
                echo -e "    ${GREEN}✅ ${APP_NAME}.env${NC}" | tee -a "$LOG_FILE"
                ((BACKUP_SUCCESS++))
            else
                echo -e "    ${RED}❌ Failed to backup ${APP_NAME}.env${NC}" | tee -a "$LOG_FILE"
                ((BACKUP_ERRORS++))
            fi
        done
    else
        echo -e "${YELLOW}⚠️  GPG not available - copying env files unencrypted${NC}" | tee -a "$LOG_FILE"
        echo "$ENV_FILES" | while read -r ENVFILE; do
            APP_NAME=$(basename "$(dirname "$ENVFILE")")
            cp "$ENVFILE" "${BACKUP_DIR}/configs/${APP_NAME}.env"
            echo -e "  ${GREEN}✅ ${APP_NAME}.env${NC} (unencrypted)" | tee -a "$LOG_FILE"
            ((BACKUP_SUCCESS++))
        done
    fi
else
    echo "⚠️  No .env files found" | tee -a "$LOG_FILE"
fi

################################################################################
# 5. BACKUP PORTAINER STACKS
################################################################################
echo -e "\n${BLUE}📚 Backing up Portainer stacks...${NC}" | tee -a "$LOG_FILE"

if [ -d "infrastructure/portainer" ]; then
    cp -r infrastructure/portainer "${BACKUP_DIR}/stacks/portainer"
    echo -e "${GREEN}✅ Portainer configurations backed up${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
else
    echo -e "${YELLOW}⚠️  Portainer infrastructure directory not found${NC}" | tee -a "$LOG_FILE"
fi

# Backup all docker-compose files
echo "  → Backing up docker-compose files..." | tee -a "$LOG_FILE"
find stacks/ -name "docker-compose.yml" -o -name "docker-compose.*.yml" 2>/dev/null | while read -r COMPOSE_FILE; do
    RELATIVE_PATH=$(echo "$COMPOSE_FILE" | sed 's|^stacks/||')
    TARGET_DIR="${BACKUP_DIR}/stacks/$(dirname "$RELATIVE_PATH")"
    mkdir -p "$TARGET_DIR"
    cp "$COMPOSE_FILE" "$TARGET_DIR/"
    echo "    ✅ $(basename "$COMPOSE_FILE")" | tee -a "$LOG_FILE"
done
((BACKUP_SUCCESS++))

################################################################################
# 6. BACKUP NGINX CONFIGURATIONS
################################################################################
echo -e "\n${BLUE}🌐 Backing up Nginx configurations...${NC}" | tee -a "$LOG_FILE"

# Check for Nginx configs in various locations
NGINX_LOCATIONS=("legacy-ai-stacks/automation/dify/nginx" "infrastructure/nginx" "legacy-ai-stacks/automation/*/nginx")

for LOCATION_PATTERN in "${NGINX_LOCATIONS[@]}"; do
    for NGINX_DIR in $LOCATION_PATTERN; do
        if [ -d "$NGINX_DIR" ]; then
            RELATIVE_PATH=$(echo "$NGINX_DIR" | sed 's|^legacy-ai-stacks/||;s|^infrastructure/||')
            TARGET_DIR="${BACKUP_DIR}/configs/nginx/${RELATIVE_PATH}"
            mkdir -p "$TARGET_DIR"
            cp -r "$NGINX_DIR"/* "$TARGET_DIR/" 2>/dev/null || true
            echo -e "  ${GREEN}✅ Backed up${NC}: $NGINX_DIR" | tee -a "$LOG_FILE"
            ((BACKUP_SUCCESS++))
        fi
    done
done

################################################################################
# 7. BACKUP CADDY CONFIGURATION
################################################################################
echo -e "\n${BLUE}📦 Backing up Caddy configuration...${NC}" | tee -a "$LOG_FILE"

# Try to extract Caddyfile from running container
if docker ps --format '{{.Names}}' | grep -q "^caddy$"; then
    if docker cp caddy:/etc/caddy/Caddyfile "${BACKUP_DIR}/configs/Caddyfile" 2>/dev/null; then
        echo -e "${GREEN}✅ Caddyfile backed up${NC}" | tee -a "$LOG_FILE"
        ((BACKUP_SUCCESS++))
    else
        echo -e "${YELLOW}⚠️  Could not extract Caddyfile from container${NC}" | tee -a "$LOG_FILE"
    fi
fi

# Backup Caddy data volume (SSL certificates)
if docker volume ls -q | grep -q "caddy_data"; then
    echo "  → Backing up Caddy SSL certificates..." | tee -a "$LOG_FILE"
    docker run --rm \
        -v caddy_data:/data:ro \
        -v "${BACKUP_DIR}/configs:/backup" \
        alpine \
        tar czf "/backup/caddy-ssl-certs.tar.gz" -C /data . >> "$LOG_FILE" 2>&1
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Caddy SSL certificates backed up${NC}" | tee -a "$LOG_FILE"
        ((BACKUP_SUCCESS++))
    fi
fi

################################################################################
# 8. BACKUP DOCUMENTATION
################################################################################
echo -e "\n${BLUE}📚 Backing up documentation...${NC}" | tee -a "$LOG_FILE"

if [ -d "docs" ]; then
    cp -r docs "${BACKUP_DIR}/documentation"
    echo -e "${GREEN}✅ Documentation backed up${NC}" | tee -a "$LOG_FILE"
    ((BACKUP_SUCCESS++))
fi

################################################################################
# 9. CREATE BACKUP MANIFEST
################################################################################
echo -e "\n${BLUE}📝 Creating backup manifest...${NC}" | tee -a "$LOG_FILE"

cat > "${BACKUP_DIR}/BACKUP-MANIFEST.md" << EOF
# Backup Manifest
**Timestamp**: $(date -u +%Y-%m-%dT%H:%M:%SZ)  
**Backup ID**: ${TIMESTAMP}  
**Location**: ${BACKUP_DIR}

## Summary
- **Success**: ${BACKUP_SUCCESS} items
- **Errors**: ${BACKUP_ERRORS} items
- **Total Size**: $(du -sh "$BACKUP_DIR" | cut -f1)

## Contents

### Docker Volumes ($(ls -1 "${BACKUP_DIR}/volumes" 2>/dev/null | wc -l) volumes)
$(ls -1h "${BACKUP_DIR}/volumes" 2>/dev/null || echo "None")

### Databases
- PostgreSQL: n8n, Dify
- Redis: n8n, Dify

### Configuration Files
$(ls -1 "${BACKUP_DIR}/configs" 2>/dev/null || echo "None")

### Docker Compose Stacks
$(find "${BACKUP_DIR}/stacks" -name "docker-compose*.yml" 2>/dev/null | wc -l) compose files

## Restore Instructions
See: /docs/guides/RESTORE.md

## Retention Policy
Backups older than ${RETENTION_DAYS} days will be automatically deleted.

## Encryption
$(if [ -n "$BACKUP_ENCRYPTION_KEY" ]; then echo "Environment files encrypted with GPG/AES256"; else echo "⚠️  Environment files NOT encrypted"; fi)

---
Generated by: scripts/maintenance/backup-all.sh
EOF

echo -e "${GREEN}✅ Manifest created${NC}" | tee -a "$LOG_FILE"

################################################################################
# 10. CLEANUP OLD BACKUPS
################################################################################
echo -e "\n${BLUE}🗑️  Cleaning up old backups...${NC}" | tee -a "$LOG_FILE"

if [ -d "$BACKUP_BASE_DIR" ]; then
    OLD_BACKUPS=$(find "$BACKUP_BASE_DIR" -maxdepth 1 -type d -mtime +${RETENTION_DAYS} ! -path "$BACKUP_BASE_DIR" 2>/dev/null || echo "")
    
    if [ -n "$OLD_BACKUPS" ]; then
        OLD_COUNT=$(echo "$OLD_BACKUPS" | wc -l)
        echo "Found $OLD_COUNT old backup(s) to remove (>${RETENTION_DAYS} days)" | tee -a "$LOG_FILE"
        
        echo "$OLD_BACKUPS" | while read -r OLD_BACKUP; do
            BACKUP_AGE=$(( ($(date +%s) - $(stat -c %Y "$OLD_BACKUP" 2>/dev/null || stat -f %m "$OLD_BACKUP")) / 86400 ))
            echo "  Removing: $(basename "$OLD_BACKUP") (${BACKUP_AGE} days old)" | tee -a "$LOG_FILE"
            rm -rf "$OLD_BACKUP"
        done
    else
        echo "✅ No old backups to remove" | tee -a "$LOG_FILE"
    fi
fi

################################################################################
# 11. OPTIONAL: UPLOAD TO S3
################################################################################
# Uncomment and configure if using AWS S3 for remote backups
# echo -e "\n${BLUE}☁️  Uploading to S3...${NC}" | tee -a "$LOG_FILE"
# if command -v aws >/dev/null 2>&1; then
#     S3_BUCKET="your-backup-bucket"
#     if aws s3 sync "$BACKUP_DIR" "s3://${S3_BUCKET}/legacy-ai-backups/${TIMESTAMP}/" >> "$LOG_FILE" 2>&1; then
#         echo -e "${GREEN}✅ Uploaded to S3${NC}" | tee -a "$LOG_FILE"
#     else
#         echo -e "${RED}❌ S3 upload failed${NC}" | tee -a "$LOG_FILE"
#         ((BACKUP_ERRORS++))
#     fi
# else
#     echo -e "${YELLOW}⚠️  AWS CLI not available - skipping S3 upload${NC}" | tee -a "$LOG_FILE"
# fi

################################################################################
# 12. SUMMARY
################################################################################
echo -e "\n========================================" | tee -a "$LOG_FILE"
echo -e "${BLUE}💾 BACKUP SUMMARY:${NC}" | tee -a "$LOG_FILE"
echo "  Successful: $BACKUP_SUCCESS items" | tee -a "$LOG_FILE"
echo "  Errors: $BACKUP_ERRORS items" | tee -a "$LOG_FILE"
echo "  Location: $BACKUP_DIR" | tee -a "$LOG_FILE"
echo "  Total size: $(du -sh "$BACKUP_DIR" | cut -f1)" | tee -a "$LOG_FILE"
echo "  Manifest: ${BACKUP_DIR}/BACKUP-MANIFEST.md" | tee -a "$LOG_FILE"
echo "  Log file: $LOG_FILE" | tee -a "$LOG_FILE"

if [ "$BACKUP_ERRORS" -eq 0 ]; then
    echo -e "\n${GREEN}✅ Backup completed successfully!${NC}" | tee -a "$LOG_FILE"
    exit 0
else
    echo -e "\n${YELLOW}⚠️  Backup completed with $BACKUP_ERRORS error(s)${NC}" | tee -a "$LOG_FILE"
    exit 1
fi
