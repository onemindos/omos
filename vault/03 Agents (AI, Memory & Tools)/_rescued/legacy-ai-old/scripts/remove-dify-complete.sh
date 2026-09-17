#!/bin/bash
###############################################################################
# Legacy AI - Complete Dify Removal Script
# Removes all Dify containers, volumes, configs, and references from VM
###############################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SSH_KEY="/Users/zeuslegacy/Documents/LightsailDefaultKey-us-east-1.pem"
VM_HOST="ubuntu@3.215.29.128"
BACKUP_DIR="$HOME/dify-removal-backup-$(date +%Y%m%d_%H%M%S)"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║       Legacy AI - Complete Dify Removal Script            ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Confirmation
echo -e "${YELLOW}⚠️  WARNING: This will PERMANENTLY remove:${NC}"
echo "  • 7 Dify containers (web, api, worker, nginx, sandbox, postgres, redis)"
echo "  • All Dify Docker volumes and data"
echo "  • Dify configurations from VM"
echo "  • Dify domain from Caddy reverse proxy"
echo ""
echo -e "${RED}This action CANNOT be undone!${NC}"
echo ""
read -p "Type 'DELETE DIFY' to confirm removal: " CONFIRM

if [ "$CONFIRM" != "DELETE DIFY" ]; then
    echo -e "${RED}❌ Removal cancelled${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Confirmed. Starting Dify removal...${NC}"
echo ""

###############################################################################
# Phase 1: Create Backup
###############################################################################

echo -e "${BLUE}═══ Phase 1: Creating Backup ═══${NC}"

mkdir -p "$BACKUP_DIR"

echo "  📦 Backing up Dify data from VM..."
ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    mkdir -p ~/dify-removal-backup
    
    # Backup Dify PostgreSQL database
    if docker ps | grep -q dify-postgres; then
        echo '    • Backing up Dify database...'
        docker exec dify-postgres pg_dump -U dify dify > ~/dify-removal-backup/dify_database.sql 2>/dev/null || true
    fi
    
    # Backup Dify Redis data
    if docker ps | grep -q dify-redis; then
        echo '    • Backing up Dify Redis data...'
        docker exec dify-redis redis-cli --no-auth-warning -a '2MehsKm00ldzVYo6uG06pkn4J8hr5Mnw' SAVE 2>/dev/null || true
        docker cp dify-redis:/data/dump.rdb ~/dify-removal-backup/dify_redis.rdb 2>/dev/null || true
    fi
    
    # Backup Dify configurations
    if [ -d /opt/legacy-ai/stacks/apps/dify ]; then
        echo '    • Backing up Dify configurations...'
        tar -czf ~/dify-removal-backup/dify_configs.tar.gz -C /opt/legacy-ai/stacks/apps dify 2>/dev/null || true
    fi
    
    # Create backup info file
    cat > ~/dify-removal-backup/BACKUP_INFO.txt << 'EOF'
Dify Removal Backup
Date: $(date)
Containers removed: 7 (web, api, worker, nginx, sandbox, postgres, redis)
Volumes removed: 6 (postgres_data, redis_data, app_storage, nginx_cache, nginx_logs, sandbox_dependencies)
EOF
    
    # Compress entire backup
    tar -czf ~/dify-removal-backup.tar.gz -C ~/ dify-removal-backup
    rm -rf ~/dify-removal-backup
    
    echo '    ✅ Backup created: ~/dify-removal-backup.tar.gz'
REMOTE_SCRIPT

# Download backup to local machine
echo "  📥 Downloading backup to local machine..."
scp -i "$SSH_KEY" "$VM_HOST:~/dify-removal-backup.tar.gz" "$BACKUP_DIR/" 2>/dev/null || echo "    ⚠️  Backup download skipped"

echo -e "${GREEN}✅ Phase 1 Complete: Backup created${NC}"
echo "   Local: $BACKUP_DIR/dify-removal-backup.tar.gz (if download succeeded)"
echo "   VM: ~/dify-removal-backup.tar.gz"
echo ""

###############################################################################
# Phase 2: Stop and Remove Dify Containers
###############################################################################

echo -e "${BLUE}═══ Phase 2: Stopping and Removing Dify Containers ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  🛑 Stopping Dify containers..."
    
    # Stop all Dify containers
    for container in dify-web dify-api dify-worker dify-nginx dify-sandbox dify-postgres dify-redis; do
        if docker ps -a | grep -q $container; then
            echo "    • Stopping $container..."
            docker stop $container 2>/dev/null || true
        fi
    done
    
    echo "  🗑️  Removing Dify containers..."
    
    # Remove all Dify containers
    for container in dify-web dify-api dify-worker dify-nginx dify-sandbox dify-postgres dify-redis; do
        if docker ps -a | grep -q $container; then
            echo "    • Removing $container..."
            docker rm -f $container 2>/dev/null || true
        fi
    done
    
    echo "    ✅ All Dify containers removed"
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 2 Complete: Containers removed${NC}"
echo ""

###############################################################################
# Phase 3: Remove Dify Docker Volumes
###############################################################################

echo -e "${BLUE}═══ Phase 3: Removing Dify Docker Volumes ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  💾 Removing Dify volumes..."
    
    # Remove all Dify-related volumes
    for volume in dify_postgres_data dify_redis_data dify_app_storage dify_nginx_cache dify_nginx_logs dify_sandbox_dependencies; do
        if docker volume ls | grep -q $volume; then
            echo "    • Removing volume: $volume..."
            docker volume rm $volume 2>/dev/null || true
        fi
    done
    
    echo "    ✅ All Dify volumes removed"
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 3 Complete: Volumes removed${NC}"
echo ""

###############################################################################
# Phase 4: Remove Dify Configuration Files
###############################################################################

echo -e "${BLUE}═══ Phase 4: Removing Dify Configuration Files ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  📁 Removing Dify configuration directories..."
    
    # Remove Dify stack directory
    if [ -d /opt/legacy-ai/stacks/apps/dify ]; then
        echo "    • Removing /opt/legacy-ai/stacks/apps/dify..."
        sudo rm -rf /opt/legacy-ai/stacks/apps/dify
    fi
    
    # Remove any Dify-related scripts
    if [ -f /opt/legacy-ai/scripts/deploy-dify.sh ]; then
        echo "    • Removing deploy-dify.sh..."
        sudo rm -f /opt/legacy-ai/scripts/deploy-dify.sh
    fi
    
    echo "    ✅ Dify configuration files removed"
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 4 Complete: Configuration files removed${NC}"
echo ""

###############################################################################
# Phase 5: Remove Dify from Caddy Reverse Proxy
###############################################################################

echo -e "${BLUE}═══ Phase 5: Removing Dify from Caddy Configuration ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  🌐 Updating Caddy configuration..."
    
    # Backup current Caddyfile
    if [ -f /opt/legacy-ai/infrastructure/caddy/Caddyfile ]; then
        sudo cp /opt/legacy-ai/infrastructure/caddy/Caddyfile /opt/legacy-ai/infrastructure/caddy/Caddyfile.backup-$(date +%Y%m%d)
        
        # Remove Dify domain block from Caddyfile
        sudo sed -i '/^dify\.zeusops\.dev {/,/^}/d' /opt/legacy-ai/infrastructure/caddy/Caddyfile
        
        echo "    • Removed dify.zeusops.dev from Caddyfile"
        
        # Reload Caddy to apply changes
        if docker ps | grep -q caddy; then
            echo "    • Reloading Caddy..."
            docker exec caddy caddy reload --config /etc/caddy/Caddyfile 2>/dev/null || true
        fi
    fi
    
    echo "    ✅ Caddy configuration updated"
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 5 Complete: Caddy configuration updated${NC}"
echo ""

###############################################################################
# Phase 6: Clean Up Docker Images
###############################################################################

echo -e "${BLUE}═══ Phase 6: Removing Dify Docker Images ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  🖼️  Removing Dify Docker images..."
    
    # Remove Dify images
    docker images | grep langgenius/dify | awk '{print $3}' | xargs -r docker rmi -f 2>/dev/null || echo "    • No Dify images to remove"
    
    # Remove pgvector image if not used by other services
    if ! docker ps -a | grep -v dify | grep -q pgvector; then
        docker images | grep pgvector/pgvector | awk '{print $3}' | xargs -r docker rmi -f 2>/dev/null || echo "    • No pgvector images to remove"
    fi
    
    echo "    ✅ Dify images removed"
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 6 Complete: Docker images removed${NC}"
echo ""

###############################################################################
# Phase 7: Verify Removal on VM
###############################################################################

echo -e "${BLUE}═══ Phase 7: Verifying Removal ═══${NC}"

ssh -i "$SSH_KEY" "$VM_HOST" bash << 'REMOTE_SCRIPT'
    echo "  🔍 Verification Report:"
    echo ""
    
    # Check for Dify containers
    DIFY_CONTAINERS=$(docker ps -a | grep -i dify | wc -l)
    if [ "$DIFY_CONTAINERS" -eq 0 ]; then
        echo "    ✅ Containers: 0 Dify containers found"
    else
        echo "    ⚠️  Containers: $DIFY_CONTAINERS Dify containers still exist"
    fi
    
    # Check for Dify volumes
    DIFY_VOLUMES=$(docker volume ls | grep -i dify | wc -l)
    if [ "$DIFY_VOLUMES" -eq 0 ]; then
        echo "    ✅ Volumes: 0 Dify volumes found"
    else
        echo "    ⚠️  Volumes: $DIFY_VOLUMES Dify volumes still exist"
    fi
    
    # Check for Dify images
    DIFY_IMAGES=$(docker images | grep -i dify | wc -l)
    if [ "$DIFY_IMAGES" -eq 0 ]; then
        echo "    ✅ Images: 0 Dify images found"
    else
        echo "    ⚠️  Images: $DIFY_IMAGES Dify images still exist"
    fi
    
    # Check for Dify in Caddyfile
    if [ -f /opt/legacy-ai/infrastructure/caddy/Caddyfile ]; then
        if grep -q "dify\.zeusops\.dev" /opt/legacy-ai/infrastructure/caddy/Caddyfile; then
            echo "    ⚠️  Caddy: dify.zeusops.dev still in Caddyfile"
        else
            echo "    ✅ Caddy: dify.zeusops.dev removed from configuration"
        fi
    fi
    
    # Check disk space reclaimed
    echo ""
    echo "  💾 Disk Space:"
    df -h / | tail -1 | awk '{print "    Used: "$3" / "$2" ("$5")"}'
    
    echo ""
    echo "  🎯 Remaining Active Services:"
    docker ps --format "table {{.Names}}\t{{.Status}}" | head -15
REMOTE_SCRIPT

echo -e "${GREEN}✅ Phase 7 Complete: Verification finished${NC}"
echo ""

###############################################################################
# Final Report
###############################################################################

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║              DIFY REMOVAL COMPLETE                         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✅ All Dify components have been removed from the VM${NC}"
echo ""
echo "📦 Backup Location:"
echo "   Local: $BACKUP_DIR/dify-removal-backup.tar.gz"
echo "   VM:    ~/dify-removal-backup.tar.gz"
echo ""
echo "🗑️  Removed Components:"
echo "   • 7 Docker containers (web, api, worker, nginx, sandbox, postgres, redis)"
echo "   • 6 Docker volumes (data, storage, cache, logs)"
echo "   • Dify Docker images"
echo "   • Dify configuration files"
echo "   • dify.zeusops.dev from Caddy"
echo ""
echo "🎯 Remaining Services:"
echo "   • Letta (letta.zeusops.dev)"
echo "   • n8n (n8n.zeusops.dev)"
echo "   • Node-RED (nodered.zeusops.dev)"
echo "   • Home Assistant (ha.zeusops.dev)"
echo "   • EMQX (emqx.zeusops.dev)"
echo "   • Portainer (command.zeusops.dev)"
echo ""
echo -e "${YELLOW}📝 Next Steps:${NC}"
echo "   1. Run local repository cleanup: ./scripts/cleanup-dify-references.sh"
echo "   2. Verify all remaining services are working"
echo "   3. Optional: Remove backup after confirming everything works"
echo ""

# Save removal report
cat > "$BACKUP_DIR/REMOVAL_REPORT.txt" << EOF
Dify Removal Report
Date: $(date)
Status: COMPLETE

Removed from VM:
- Containers: 7 (dify-web, dify-api, dify-worker, dify-nginx, dify-sandbox, dify-postgres, dify-redis)
- Volumes: 6 (postgres_data, redis_data, app_storage, nginx_cache, nginx_logs, sandbox_dependencies)
- Images: All langgenius/dify-* images
- Configuration: /opt/legacy-ai/stacks/apps/dify
- Domain: dify.zeusops.dev removed from Caddy

Backup Location:
- Local: $BACKUP_DIR/dify-removal-backup.tar.gz
- VM: ~/dify-removal-backup.tar.gz

To restore (if needed):
1. Extract backup: tar -xzf dify-removal-backup.tar.gz
2. Deploy from backup configurations
3. Restore databases from SQL dumps
EOF

echo "✅ Removal report saved to: $BACKUP_DIR/REMOVAL_REPORT.txt"
echo ""
echo "═══════════════════════════════════════════════════════════"
