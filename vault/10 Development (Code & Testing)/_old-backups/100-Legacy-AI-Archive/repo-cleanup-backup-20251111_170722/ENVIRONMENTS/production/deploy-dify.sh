#!/bin/bash
################################################################################
# DIFY DEPLOYMENT SCRIPT
# Purpose: Deploy/update Dify AI platform
# Usage: ./scripts/deployment/deploy-dify.sh
################################################################################

set -e

APP_NAME="dify"
APP_DIR="/home/ubuntu/Legacy-ai-Hybrid/legacy-ai-stacks/automation/dify"
HEALTH_URL="https://dify.zeusops.dev"
COMPOSE_FILE="docker-compose.yml"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 DEPLOYING DIFY - $(date)${NC}"
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

# Check if containers are running
RUNNING_CONTAINERS=$(docker ps --filter name=dify --format '{{.Names}}' | wc -l)
echo "Currently running Dify containers: $RUNNING_CONTAINERS"

# Backup current state
echo -e "\n${YELLOW}💾 Creating backup before deployment...${NC}"
BACKUP_DIR="/tmp/dify-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Backup database
if docker ps --format '{{.Names}}' | grep -q "^dify-postgres$"; then
    echo "  → Backing up PostgreSQL database..."
    docker exec dify-postgres pg_dump -U postgres -F c -f /tmp/pre-deploy-backup.dump
    docker cp dify-postgres:/tmp/pre-deploy-backup.dump "$BACKUP_DIR/dify-postgres.dump"
    docker exec dify-postgres rm /tmp/pre-deploy-backup.dump
    echo -e "  ${GREEN}✅ Database backed up${NC}"
fi

# Save container configs
echo "  → Saving container configurations..."
docker inspect dify-api > "$BACKUP_DIR/dify-api-config.json" 2>/dev/null || true
echo -e "  ${GREEN}✅ Configurations saved${NC}"

################################################################################
# 2. PULL LATEST DOCKER IMAGES
################################################################################
echo -e "\n${YELLOW}🐳 Pulling latest Docker images...${NC}"

docker-compose pull --quiet

echo -e "${GREEN}✅ Images updated${NC}"

################################################################################
# 3. DEPLOY NEW CONTAINERS
################################################################################
echo -e "\n${YELLOW}🚀 Deploying containers...${NC}"

# Recreate containers with new images
docker-compose up -d --remove-orphans --force-recreate

echo -e "${GREEN}✅ Containers deployed${NC}"

################################################################################
# 4. WAIT FOR HEALTH CHECK
################################################################################
echo -e "\n${YELLOW}🏥 Waiting for health check...${NC}"

MAX_WAIT=90
WAIT_COUNT=0

while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
    sleep 3
    ((WAIT_COUNT+=3))
    
    # Check if health endpoint responds
    if curl -sSf -k "$HEALTH_URL" > /dev/null 2>&1; then
        echo -e "\n${GREEN}✅ Health check passed after ${WAIT_COUNT}s${NC}"
        break
    fi
    
    echo -n "."
done

if [ $WAIT_COUNT -ge $MAX_WAIT ]; then
    echo -e "\n${RED}❌ Health check failed - ROLLING BACK${NC}"
    
    # Rollback
    echo "Rolling back to previous containers..."
    docker-compose down
    
    # Restore from backup
    echo "Restoring database backup..."
    docker-compose up -d dify-postgres
    sleep 5
    docker cp "$BACKUP_DIR/dify-postgres.dump" dify-postgres:/tmp/restore.dump
    docker exec dify-postgres pg_restore -U postgres -c /tmp/restore.dump
    docker exec dify-postgres rm /tmp/restore.dump
    
    # Start old containers
    docker-compose up -d
    
    echo -e "${RED}❌ Deployment failed and rolled back${NC}"
    exit 1
fi

################################################################################
# 5. VERIFY DEPLOYMENT
################################################################################
echo -e "\n${YELLOW}🔍 Verifying deployment...${NC}"

# Check all Dify containers are running
DIFY_CONTAINERS=("dify-api" "dify-worker" "dify-web" "dify-nginx" "dify-sandbox" "dify-postgres" "dify-redis")

ALL_RUNNING=true
for CONTAINER in "${DIFY_CONTAINERS[@]}"; do
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
        STATUS=$(docker inspect --format='{{.State.Status}}' "$CONTAINER")
        echo -e "  ${GREEN}✅ ${CONTAINER}${NC}: $STATUS"
    else
        echo -e "  ${RED}❌ ${CONTAINER}${NC}: NOT RUNNING"
        ALL_RUNNING=false
    fi
done

if [ "$ALL_RUNNING" = false ]; then
    echo -e "\n${RED}❌ Some containers are not running - check logs${NC}"
    exit 1
fi

# Check database connection
echo -e "\n${YELLOW}📊 Checking database...${NC}"
if docker exec dify-postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Database connection OK${NC}"
else
    echo -e "${RED}❌ Database connection failed${NC}"
    exit 1
fi

# Check Redis connection
echo -e "\n${YELLOW}🔴 Checking Redis...${NC}"
if docker exec dify-redis redis-cli ping > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Redis connection OK${NC}"
else
    echo -e "${RED}❌ Redis connection failed${NC}"
    exit 1
fi

################################################################################
# 6. CLEANUP
################################################################################
echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"

# Remove dangling images
docker image prune -f > /dev/null 2>&1

# Keep backup for 24 hours
echo "Backup location: $BACKUP_DIR"
echo "(Will be auto-cleaned after 24 hours)"

################################################################################
# 7. SUMMARY
################################################################################
echo -e "\n========================================"
echo -e "${GREEN}✅ DIFY DEPLOYMENT SUCCESSFUL${NC}"
echo "========================================"
echo "  Deployed at: $(date)"
echo "  Application: https://dify.zeusops.dev"
echo "  Containers: ${#DIFY_CONTAINERS[@]}"
echo "  Backup: $BACKUP_DIR"
echo -e "========================================\n"

exit 0
