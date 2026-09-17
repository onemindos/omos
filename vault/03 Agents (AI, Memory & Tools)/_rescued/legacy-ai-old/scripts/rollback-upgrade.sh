#!/bin/bash
################################################################################
# LEGACY AI - VM UPGRADE ROLLBACK
# Purpose: Rollback system to pre-upgrade state if issues occur
# Usage: sudo bash rollback-upgrade.sh <backup-directory>
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check arguments
if [ -z "$1" ]; then
    echo -e "${RED}❌ Error: Backup directory required${NC}"
    echo ""
    echo "Usage: sudo bash $0 <backup-directory>"
    echo ""
    echo "Available backups:"
    ls -d ~/backups/pre-upgrade-* 2>/dev/null || echo "  No backups found"
    exit 1
fi

BACKUP_DIR="$1"

if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}❌ Error: Backup directory does not exist: $BACKUP_DIR${NC}"
    exit 1
fi

echo -e "${RED}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${RED}║            LEGACY AI - UPGRADE ROLLBACK                    ║${NC}"
echo -e "${RED}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${YELLOW}⚠️  WARNING: This will rollback your system to:${NC}"
echo "   Backup: $BACKUP_DIR"
echo ""
cat "$BACKUP_DIR/BACKUP_INFO.txt" 2>/dev/null || echo "   (Backup info not available)"
echo ""

read -p "Continue with rollback? (yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${YELLOW}Rollback cancelled${NC}"
    exit 0
fi

LOG_FILE="/var/log/vm-rollback-$(date +%Y%m%d_%H%M%S).log"

log() {
    echo -e "$1" | tee -a "$LOG_FILE"
}

################################################################################
# PHASE 1: STOP ALL CONTAINERS
################################################################################
log "${BLUE}═══ Phase 1: Stopping all containers ═══${NC}"

log "${YELLOW}🛑 Stopping all running containers...${NC}"
docker ps -q | xargs -r docker stop
log "${GREEN}✅ All containers stopped${NC}"
log ""

################################################################################
# PHASE 2: RESTORE DOCKER CONFIGURATIONS
################################################################################
log "${BLUE}═══ Phase 2: Restoring configurations ═══${NC}"

log "${YELLOW}📦 Restoring docker compose.yml files...${NC}"

if [ -d "$BACKUP_DIR/docker-configs" ]; then
    # Backup current configs first
    log "  → Backing up current configs..."
    mkdir -p /tmp/current-configs-backup
    find /opt/legacy-ai -name "docker compose.yml" -exec cp --parents {} /tmp/current-configs-backup \; 2>/dev/null || true
    
    # Restore from backup
    log "  → Restoring from backup..."
    cp -r "$BACKUP_DIR/docker-configs/"* /opt/legacy-ai/ 2>/dev/null || true
    
    log "${GREEN}✅ Configurations restored${NC}"
else
    log "${YELLOW}⚠️  No docker configs in backup${NC}"
fi
log ""

################################################################################
# PHASE 3: RESTORE CADDY CONFIGURATION
################################################################################
log "${BLUE}═══ Phase 3: Restoring Caddy configuration ═══${NC}"

if [ -f "$BACKUP_DIR/caddy/Caddyfile" ]; then
    log "${YELLOW}🔄 Restoring Caddyfile...${NC}"
    cp "$BACKUP_DIR/caddy/Caddyfile" /opt/legacy-ai/infrastructure/caddy/Caddyfile
    log "${GREEN}✅ Caddyfile restored${NC}"
else
    log "${YELLOW}⚠️  No Caddyfile in backup${NC}"
fi
log ""

################################################################################
# PHASE 4: RESTORE DATABASES
################################################################################
log "${BLUE}═══ Phase 4: Restoring databases ═══${NC}"

log "${YELLOW}💾 Starting database containers for restore...${NC}"

# Start PostgreSQL containers
for db_backup in "$BACKUP_DIR/databases/"*_full.sql; do
    if [ -f "$db_backup" ]; then
        container_name=$(basename "$db_backup" | sed 's/_full.sql//')
        
        log "  → Restoring $container_name..."
        
        # Find and start the container
        COMPOSE_DIR=$(find /opt/legacy-ai -name docker compose.yml -exec grep -l "$container_name" {} \; -exec dirname {} \; 2>/dev/null | head -1)
        
        if [ -n "$COMPOSE_DIR" ]; then
            cd "$COMPOSE_DIR"
            docker compose up -d "$container_name" 2>/dev/null || true
            sleep 5
            
            # Restore database
            docker exec -i "$container_name" psql -U postgres < "$db_backup" 2>&1 | tee -a "$LOG_FILE"
            
            log "${GREEN}  ✅ $container_name restored${NC}"
        else
            log "${YELLOW}  ⚠️  Could not find docker compose for $container_name${NC}"
        fi
    fi
done

# Restore Redis databases
for redis_backup in "$BACKUP_DIR/databases/"*_dump.rdb; do
    if [ -f "$redis_backup" ]; then
        container_name=$(basename "$redis_backup" | sed 's/_dump.rdb//')
        
        log "  → Restoring $container_name..."
        
        # Copy dump file to container
        docker cp "$redis_backup" "$container_name:/data/dump.rdb" 2>/dev/null && \
            docker restart "$container_name" && \
            log "${GREEN}  ✅ $container_name restored${NC}" || \
            log "${YELLOW}  ⚠️  Could not restore $container_name${NC}"
    fi
done

log "${GREEN}✅ Database restore attempted${NC}"
log ""

################################################################################
# PHASE 5: RESTART ALL CONTAINERS
################################################################################
log "${BLUE}═══ Phase 5: Restarting all containers ═══${NC}"

log "${YELLOW}🚀 Starting all containers...${NC}"

# Find all docker compose directories and start them
find /opt/legacy-ai -type f -name "docker compose.yml" -exec dirname {} \; 2>/dev/null | sort -u | while read -r dir; do
    log "  → Starting stack in $dir..."
    cd "$dir"
    docker compose up -d 2>&1 | tee -a "$LOG_FILE"
done

log "${GREEN}✅ All containers started${NC}"
log ""

# Wait for containers to stabilize
log "${YELLOW}⏳ Waiting 30 seconds for containers to stabilize...${NC}"
sleep 30

################################################################################
# PHASE 6: VERIFICATION
################################################################################
log "${BLUE}═══ Phase 6: Verification ═══${NC}"

log "${YELLOW}📊 Container status:${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}" | tee -a "$LOG_FILE"

RUNNING=$(docker ps -q | wc -l)
log ""
log "${CYAN}Running containers: $RUNNING${NC}"

log ""
log "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
log "${GREEN}║              ROLLBACK COMPLETE                             ║${NC}"
log "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
log ""

log "${CYAN}📋 Summary:${NC}"
log "   • Backup used: $BACKUP_DIR"
log "   • Containers running: $RUNNING"
log "   • Current configs backed up: /tmp/current-configs-backup"
log "   • Rollback log: $LOG_FILE"
log ""

log "${YELLOW}⚠️  Next Steps:${NC}"
log "   1. Verify all services are working"
log "   2. Check application functionality"
log "   3. Review logs: tail -f $LOG_FILE"
log "   4. If issues persist, check individual container logs"
log ""

log "${CYAN}📝 To check container logs:${NC}"
log "   docker logs <container-name>"
log ""
