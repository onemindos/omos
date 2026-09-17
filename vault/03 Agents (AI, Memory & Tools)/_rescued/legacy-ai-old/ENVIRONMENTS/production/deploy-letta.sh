#!/bin/bash
################################################################################
# LETTA DEPLOYMENT SCRIPT
# Purpose: Deploy/update Letta AI assistant
# Usage: ./scripts/deployment/deploy-letta.sh
################################################################################

set -e

APP_NAME="letta"
APP_DIR="/home/ubuntu/Legacy-ai-Hybrid/legacy-ai-stacks/ai/letta"
HEALTH_URL="https://letta.zeusops.dev"
COMPOSE_FILE="docker-compose.yml"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 DEPLOYING LETTA - $(date)${NC}"
echo "========================================"

# Check if running on production VM
if [ ! -d "$APP_DIR" ]; then
    echo -e "${RED}❌ App directory not found: $APP_DIR${NC}"
    echo "This script must be run on the production VM"
    exit 1
fi

cd "$APP_DIR"

################################################################################
# 1. PRE-DEPLOYMENT CHECKS
################################################################################
echo -e "\n${YELLOW}📋 Pre-deployment checks...${NC}"

# Check if docker-compose exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo -e "${RED}❌ $COMPOSE_FILE not found in $APP_DIR${NC}"
    exit 1
fi

# Check if container is running
if docker ps --format '{{.Names}}' | grep -q "^letta$"; then
    echo "Letta container is currently running"
else
    echo "Letta container is not running (will be started)"
fi

################################################################################
# 2. BACKUP LETTA DATA
################################################################################
echo -e "\n${YELLOW}💾 Backing up Letta data...${NC}"
BACKUP_DIR="/tmp/letta-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Backup Letta volume if exists
if docker volume ls -q | grep -q "letta_data"; then
    echo "  → Backing up Letta data volume..."
    docker run --rm \
        -v letta_data:/data:ro \
        -v "$BACKUP_DIR:/backup" \
        alpine \
        tar czf "/backup/letta-data.tar.gz" -C /data .
    echo -e "  ${GREEN}✅ Letta data backed up${NC}"
fi

################################################################################
# 3. PULL LATEST DOCKER IMAGE
################################################################################
echo -e "\n${YELLOW}🐳 Pulling latest Docker image...${NC}"

docker-compose pull --quiet

echo -e "${GREEN}✅ Image updated${NC}"

################################################################################
# 4. DEPLOY NEW CONTAINER
################################################################################
echo -e "\n${YELLOW}🚀 Deploying container...${NC}"

# Recreate container with new image
docker-compose up -d --force-recreate

echo -e "${GREEN}✅ Container deployed${NC}"

################################################################################
# 5. WAIT FOR HEALTH CHECK
################################################################################
echo -e "\n${YELLOW}🏥 Waiting for health check...${NC}"

MAX_WAIT=60
WAIT_COUNT=0

while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
    sleep 3
    ((WAIT_COUNT+=3))
    
    # Check if health endpoint responds
    if curl -sSf -k --max-time 10 "$HEALTH_URL" > /dev/null 2>&1; then
        echo -e "\n${GREEN}✅ Health check passed after ${WAIT_COUNT}s${NC}"
        break
    fi
    
    echo -n "."
done

if [ $WAIT_COUNT -ge $MAX_WAIT ]; then
    echo -e "\n${YELLOW}⚠️  Health check timeout - container may need more time${NC}"
fi

################################################################################
# 6. VERIFY DEPLOYMENT
################################################################################
echo -e "\n${YELLOW}🔍 Verifying deployment...${NC}"

if docker ps --format '{{.Names}}' | grep -q "^letta$"; then
    STATUS=$(docker inspect --format='{{.State.Status}}' "letta")
    echo -e "  ${GREEN}✅ letta${NC}: $STATUS"
else
    echo -e "  ${RED}❌ letta${NC}: NOT RUNNING"
    exit 1
fi

################################################################################
# 7. CLEANUP
################################################################################
echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"

# Remove dangling images
docker image prune -f > /dev/null 2>&1

echo "Backup location: $BACKUP_DIR"

################################################################################
# 8. SUMMARY
################################################################################
echo -e "\n========================================"
echo -e "${GREEN}✅ LETTA DEPLOYMENT SUCCESSFUL${NC}"
echo "========================================"
echo "  Deployed at: $(date)"
echo "  Application: https://letta.zeusops.dev"
echo "  Container: letta"
echo "  Backup: $BACKUP_DIR"
echo -e "========================================\n"

exit 0
