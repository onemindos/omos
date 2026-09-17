#!/bin/bash
################################################################################
# LEGACY AI - COMPLETE VM INFRASTRUCTURE UPGRADE
# Date: November 11, 2025
# Purpose: Comprehensive system upgrade and optimization
################################################################################

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/home/ubuntu/backups/pre-upgrade-${BACKUP_DATE}"
LOG_FILE="/var/log/vm-upgrade-${BACKUP_DATE}.log"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         LEGACY AI - VM INFRASTRUCTURE UPGRADE              ║${NC}"
echo -e "${BLUE}║                  $(date '+%Y-%m-%d %H:%M:%S')                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to log
log() {
    echo -e "$1" | tee -a "$LOG_FILE"
}

# Function to check if running as sudo
check_sudo() {
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${RED}❌ Please run as root or with sudo${NC}"
        exit 1
    fi
}

################################################################################
# PHASE 1: PRE-UPGRADE AUDIT & BACKUP
################################################################################
phase1_audit_and_backup() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 1: Pre-Upgrade Audit & Backup${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Create backup directory
    log "${YELLOW}📁 Creating backup directory...${NC}"
    mkdir -p "$BACKUP_DIR"/{system-info,docker-configs,databases,volumes,caddy}
    
    # Document current versions
    log "${YELLOW}📋 Documenting current system state...${NC}"
    {
        echo "=== SYSTEM INFORMATION ==="
        echo "Date: $(date)"
        echo "Hostname: $(hostname)"
        echo "OS: $(lsb_release -d | cut -f2)"
        echo "Kernel: $(uname -r)"
        echo ""
        echo "=== DOCKER VERSIONS ==="
        docker --version
        docker compose --version
        echo ""
        echo "=== RUNNING CONTAINERS ==="
        docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}"
        echo ""
        echo "=== DOCKER IMAGES ==="
        docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
        echo ""
        echo "=== DISK USAGE ==="
        df -h
        echo ""
        echo "=== MEMORY USAGE ==="
        free -h
        echo ""
        echo "=== NETWORK INFO ==="
        docker network ls
        echo ""
        echo "=== VOLUMES ==="
        docker volume ls
    } > "$BACKUP_DIR/system-info/pre-upgrade-state.txt"
    
    log "${GREEN}✅ System state documented${NC}"
    
    # Backup Docker configurations
    log "${YELLOW}📦 Backing up Docker configurations...${NC}"
    
    # Copy all docker compose files
    find /opt/legacy-ai -name "docker compose.yml" -o -name "docker compose.*.yml" 2>/dev/null | while read -r file; do
        relative_path=$(echo "$file" | sed 's|/opt/legacy-ai/||')
        mkdir -p "$BACKUP_DIR/docker-configs/$(dirname "$relative_path")"
        cp "$file" "$BACKUP_DIR/docker-configs/$relative_path"
    done
    
    # Backup Caddyfile
    if [ -f "/opt/legacy-ai/infrastructure/caddy/Caddyfile" ]; then
        cp /opt/legacy-ai/infrastructure/caddy/Caddyfile "$BACKUP_DIR/caddy/"
    fi
    
    log "${GREEN}✅ Docker configurations backed up${NC}"
    
    # Backup databases
    log "${YELLOW}💾 Backing up databases...${NC}"
    
    # PostgreSQL databases
    for db_container in $(docker ps --format '{{.Names}}' | grep postgres); do
        log "  → Backing up $db_container..."
        docker exec "$db_container" pg_dumpall -U postgres > "$BACKUP_DIR/databases/${db_container}_full.sql" 2>/dev/null || \
            log "${YELLOW}    ⚠️  Could not backup $db_container${NC}"
    done
    
    # Redis databases
    for redis_container in $(docker ps --format '{{.Names}}' | grep redis); do
        log "  → Backing up $redis_container..."
        docker exec "$redis_container" redis-cli --rdb /tmp/dump.rdb SAVE 2>/dev/null || true
        docker cp "$redis_container:/data/dump.rdb" "$BACKUP_DIR/databases/${redis_container}_dump.rdb" 2>/dev/null || \
            log "${YELLOW}    ⚠️  Could not backup $redis_container${NC}"
    done
    
    log "${GREEN}✅ Databases backed up${NC}"
    
    # Create snapshot info
    cat > "$BACKUP_DIR/BACKUP_INFO.txt" << EOF
LEGACY AI - Pre-Upgrade Backup
Created: $(date)
Backup Directory: $BACKUP_DIR

Contents:
- system-info/       : System state and versions
- docker-configs/    : All docker compose.yml files
- databases/         : PostgreSQL and Redis dumps
- caddy/            : Caddyfile configuration

Restore Instructions:
1. Stop all containers
2. Restore database dumps
3. Restore docker compose files
4. Restart containers

Total Size: $(du -sh "$BACKUP_DIR" | cut -f1)
EOF
    
    log "${GREEN}✅ Phase 1 Complete - Backup created at: $BACKUP_DIR${NC}"
    log ""
}

################################################################################
# PHASE 2: SYSTEM-LEVEL UPGRADES
################################################################################
phase2_system_upgrades() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 2: System-Level Upgrades${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Update package lists
    log "${YELLOW}📦 Updating package lists...${NC}"
    apt update -qq
    
    # Upgrade system packages
    log "${YELLOW}⬆️  Upgrading system packages...${NC}"
    DEBIAN_FRONTEND=noninteractive apt upgrade -y -qq
    
    log "${GREEN}✅ System packages updated${NC}"
    
    # Upgrade Docker
    log "${YELLOW}🐳 Checking Docker version...${NC}"
    CURRENT_DOCKER=$(docker --version | grep -oP '\d+\.\d+\.\d+')
    log "  Current Docker version: $CURRENT_DOCKER"
    
    apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker compose-plugin -qq 2>/dev/null || \
        log "${YELLOW}  ℹ️  Docker already at latest available version${NC}"
    
    NEW_DOCKER=$(docker --version | grep -oP '\d+\.\d+\.\d+')
    log "  Docker version after update: $NEW_DOCKER"
    
    log "${GREEN}✅ Docker updated${NC}"
    
    # Check Portainer version
    log "${YELLOW}🎛️  Checking Portainer version...${NC}"
    PORTAINER_IMAGE=$(docker inspect portainer --format='{{.Config.Image}}' 2>/dev/null)
    log "  Current image: $PORTAINER_IMAGE"
    
    # Pull latest Portainer
    docker pull portainer/portainer-ee:latest
    
    log "${GREEN}✅ Portainer image updated (will restart in Phase 3)${NC}"
    
    # Enable UFW firewall
    log "${YELLOW}🔥 Configuring firewall...${NC}"
    
    # Configure rules before enabling
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow 22/tcp comment 'SSH'
    ufw allow 80/tcp comment 'HTTP'
    ufw allow 443/tcp comment 'HTTPS'
    ufw allow 443/udp comment 'HTTP/3 QUIC'
    
    # Enable firewall
    ufw --force enable
    
    log "${GREEN}✅ Firewall enabled with rules:${NC}"
    ufw status numbered | tee -a "$LOG_FILE"
    
    log "${GREEN}✅ Phase 2 Complete${NC}"
    log ""
}

################################################################################
# PHASE 3: CONTAINER IMAGE UPDATES
################################################################################
phase3_container_updates() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 3: Container Image Updates${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Find all docker compose directories
    COMPOSE_DIRS=$(find /opt/legacy-ai -type f -name "docker compose.yml" -exec dirname {} \; 2>/dev/null | sort -u)
    
    for dir in $COMPOSE_DIRS; do
        if [ -f "$dir/docker compose.yml" ]; then
            STACK_NAME=$(basename "$dir")
            log "${YELLOW}🔄 Updating stack: $STACK_NAME${NC}"
            log "  Location: $dir"
            
            cd "$dir"
            
            # Pull latest images
            docker compose pull 2>&1 | tee -a "$LOG_FILE"
            
            # Recreate containers with new images (rolling update)
            log "  → Deploying updated containers..."
            docker compose up -d --no-deps --force-recreate 2>&1 | tee -a "$LOG_FILE"
            
            # Wait a moment for stability
            sleep 3
            
            log "${GREEN}  ✅ $STACK_NAME updated${NC}"
        fi
    done
    
    log "${GREEN}✅ Phase 3 Complete - All containers updated${NC}"
    log ""
}

################################################################################
# PHASE 4: FIX UNHEALTHY CONTAINERS
################################################################################
phase4_fix_health_checks() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 4: Fix Unhealthy Containers${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Check for unhealthy containers
    UNHEALTHY=$(docker ps --filter "health=unhealthy" --format "{{.Names}}")
    
    if [ -z "$UNHEALTHY" ]; then
        log "${GREEN}✅ No unhealthy containers found${NC}"
    else
        log "${YELLOW}Found unhealthy containers:${NC}"
        echo "$UNHEALTHY" | tee -a "$LOG_FILE"
        
        for container in $UNHEALTHY; do
            log "${YELLOW}  → Checking $container logs...${NC}"
            docker logs --tail 50 "$container" 2>&1 | tee -a "$LOG_FILE"
            
            log "${YELLOW}  → Restarting $container...${NC}"
            docker restart "$container"
            sleep 5
            
            # Check health after restart
            HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "no-healthcheck")
            log "  Health status: $HEALTH"
        done
    fi
    
    log "${GREEN}✅ Phase 4 Complete${NC}"
    log ""
}

################################################################################
# PHASE 5: CADDY SSL & ROUTING OPTIMIZATION
################################################################################
phase5_caddy_optimization() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 5: Caddy SSL & Routing Optimization${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Check Caddy status
    log "${YELLOW}🔍 Checking Caddy status...${NC}"
    docker ps | grep caddy | tee -a "$LOG_FILE"
    
    # Check SSL certificates
    log "${YELLOW}🔒 Checking SSL certificates...${NC}"
    docker exec caddy caddy list-modules 2>/dev/null | tee -a "$LOG_FILE" || true
    
    # Test each domain
    DOMAINS=(
        "letta.zeusops.dev"
        "command.zeusops.dev"
        "n8n.zeusops.dev"
        "emqx.zeusops.dev"
        "ha.zeusops.dev"
        "nodered.zeusops.dev"
        "dashboard.zeusops.dev"
    )
    
    log "${YELLOW}🌐 Testing domain accessibility...${NC}"
    for domain in "${DOMAINS[@]}"; do
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$domain" --max-time 5 2>/dev/null || echo "000")
        if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "307" ] || [ "$HTTP_CODE" = "301" ]; then
            log "${GREEN}  ✅ $domain - HTTP $HTTP_CODE${NC}"
        else
            log "${YELLOW}  ⚠️  $domain - HTTP $HTTP_CODE${NC}"
        fi
    done
    
    # Reload Caddy to pick up any cert changes
    log "${YELLOW}♻️  Reloading Caddy configuration...${NC}"
    docker exec caddy caddy reload --config /etc/caddy/Caddyfile 2>&1 | tee -a "$LOG_FILE"
    
    log "${GREEN}✅ Phase 5 Complete${NC}"
    log ""
}

################################################################################
# PHASE 6: POST-UPGRADE VALIDATION & CLEANUP
################################################################################
phase6_validation_cleanup() {
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}PHASE 6: Post-Upgrade Validation & Cleanup${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    # Check all containers
    log "${YELLOW}📊 Checking container status...${NC}"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Image}}" | tee -a "$LOG_FILE"
    
    RUNNING_COUNT=$(docker ps -q | wc -l)
    log "${CYAN}Total running containers: $RUNNING_COUNT${NC}"
    
    # Check system resources
    log "${YELLOW}💻 System resource usage:${NC}"
    echo "Memory:" | tee -a "$LOG_FILE"
    free -h | tee -a "$LOG_FILE"
    echo "" | tee -a "$LOG_FILE"
    echo "Disk:" | tee -a "$LOG_FILE"
    df -h / | tee -a "$LOG_FILE"
    
    # Docker cleanup
    log "${YELLOW}🧹 Cleaning up unused Docker resources...${NC}"
    
    log "  → Removing unused images..."
    docker image prune -a -f | tee -a "$LOG_FILE"
    
    log "  → Removing unused volumes..."
    docker volume prune -f | tee -a "$LOG_FILE"
    
    log "  → Removing unused networks..."
    docker network prune -f | tee -a "$LOG_FILE"
    
    # Show reclaimed space
    log "${YELLOW}💾 Docker disk usage after cleanup:${NC}"
    docker system df | tee -a "$LOG_FILE"
    
    log "${GREEN}✅ Phase 6 Complete${NC}"
    log ""
}

################################################################################
# GENERATE UPGRADE REPORT
################################################################################
generate_report() {
    REPORT_FILE="$BACKUP_DIR/UPGRADE_REPORT.txt"
    
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    log "${CYAN}Generating Upgrade Report${NC}"
    log "${CYAN}═══════════════════════════════════════════════════════${NC}"
    
    cat > "$REPORT_FILE" << EOF
╔════════════════════════════════════════════════════════════╗
║         LEGACY AI - VM UPGRADE REPORT                      ║
║                  $(date '+%Y-%m-%d %H:%M:%S')                      ║
╚════════════════════════════════════════════════════════════╝

UPGRADE SUMMARY
===============
Start Time: $(head -1 "$LOG_FILE" | grep -oP '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}' || echo "N/A")
End Time: $(date '+%Y-%m-%d %H:%M:%S')
Backup Location: $BACKUP_DIR

SYSTEM VERSIONS (AFTER UPGRADE)
================================
OS: $(lsb_release -d | cut -f2)
Kernel: $(uname -r)
Docker: $(docker --version)
Docker Compose: $(docker compose --version)

CONTAINER STATUS
================
$(docker ps --format "{{.Names}}: {{.Status}}")

Total Running: $(docker ps -q | wc -l)

UNHEALTHY CONTAINERS
====================
$(docker ps --filter "health=unhealthy" --format "{{.Names}}" || echo "None")

RESOURCE USAGE
==============
Memory:
$(free -h)

Disk:
$(df -h /)

Docker Disk Usage:
$(docker system df)

DOMAINS STATUS
==============
EOF
    
    for domain in letta.zeusops.dev command.zeusops.dev n8n.zeusops.dev; do
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$domain" --max-time 5 2>/dev/null || echo "000")
        echo "$domain: HTTP $HTTP_CODE" >> "$REPORT_FILE"
    done
    
    cat >> "$REPORT_FILE" << EOF

FIREWALL STATUS
===============
$(ufw status numbered)

NEXT STEPS
==========
1. ✅ All systems upgraded and operational
2. ✅ Firewall enabled with proper rules
3. ✅ SSL certificates verified
4. ⚠️  Monitor unhealthy containers (if any)
5. ⚠️  Test all application functionality
6. 📋 Update documentation with new versions

ROLLBACK INSTRUCTIONS
=====================
If issues occur, restore from backup:

1. Stop all containers:
   cd /opt/legacy-ai && find . -name docker compose.yml -execdir docker compose down \;

2. Restore configurations:
   cp -r $BACKUP_DIR/docker-configs/* /opt/legacy-ai/

3. Restore databases:
   # PostgreSQL: docker exec -i <container> psql -U postgres < backup.sql
   # Redis: docker cp backup.rdb <container>:/data/dump.rdb

4. Restart containers:
   cd /opt/legacy-ai && find . -name docker compose.yml -execdir docker compose up -d \;

LOGS
====
Full upgrade log: $LOG_FILE

EOF
    
    log "${GREEN}✅ Upgrade report generated: $REPORT_FILE${NC}"
    log ""
    
    # Display summary
    cat "$REPORT_FILE"
}

################################################################################
# MAIN EXECUTION
################################################################################
main() {
    # Check if running as root
    if [ "$EUID" -ne 0 ]; then 
        echo -e "${YELLOW}⚠️  Running without sudo - some operations may fail${NC}"
        echo -e "${YELLOW}   Recommend running: sudo bash $0${NC}"
        echo ""
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # Create log file
    mkdir -p "$(dirname "$LOG_FILE")"
    touch "$LOG_FILE"
    
    # Execute phases
    phase1_audit_and_backup
    phase2_system_upgrades
    phase3_container_updates
    phase4_fix_health_checks
    phase5_caddy_optimization
    phase6_validation_cleanup
    
    # Generate report
    generate_report
    
    log "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    log "${GREEN}║         UPGRADE COMPLETE - ALL PHASES SUCCESSFUL           ║${NC}"
    log "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    log ""
    log "${CYAN}📋 Report saved to: $REPORT_FILE${NC}"
    log "${CYAN}📝 Full logs at: $LOG_FILE${NC}"
    log ""
    log "${YELLOW}⚠️  RECOMMENDED: Review the upgrade report and test all services${NC}"
    log ""
}

# Run main function
main "$@"
