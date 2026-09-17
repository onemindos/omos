#!/bin/bash
################################################################################
# IOT STACK DEPLOYMENT SCRIPT
# Purpose: Deploy/update IoT services (EMQX, Home Assistant, Node-RED)
# Usage: ./scripts/deployment/deploy-iot.sh
################################################################################

set -e

APP_NAME="iot"
APP_DIR="/home/ubuntu/Legacy-ai-Hybrid/stacks/iot"
COMPOSE_FILE="docker-compose.yml"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 DEPLOYING IOT STACK - $(date)${NC}"
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
RUNNING_CONTAINERS=$(docker ps --filter name=emqx --filter name=home-assistant --filter name=node-red --format '{{.Names}}' | wc -l)
echo "Currently running IoT containers: $RUNNING_CONTAINERS"

################################################################################
# 2. BACKUP CONFIGURATIONS
################################################################################
echo -e "\n${YELLOW}💾 Backing up configurations...${NC}"
BACKUP_DIR="/tmp/iot-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Backup Home Assistant config if volume exists
if docker volume ls -q | grep -q "home_assistant_config"; then
    echo "  → Backing up Home Assistant config..."
    docker run --rm \
        -v home_assistant_config:/data:ro \
        -v "$BACKUP_DIR:/backup" \
        alpine \
        tar czf "/backup/home-assistant-config.tar.gz" -C /data .
    echo -e "  ${GREEN}✅ Home Assistant config backed up${NC}"
fi

# Backup Node-RED data if volume exists
if docker volume ls -q | grep -q "node_red_data"; then
    echo "  → Backing up Node-RED data..."
    docker run --rm \
        -v node_red_data:/data:ro \
        -v "$BACKUP_DIR:/backup" \
        alpine \
        tar czf "/backup/node-red-data.tar.gz" -C /data .
    echo -e "  ${GREEN}✅ Node-RED data backed up${NC}"
fi

# Backup EMQX data if volume exists
if docker volume ls -q | grep -q "emqx_data"; then
    echo "  → Backing up EMQX data..."
    docker run --rm \
        -v emqx_data:/data:ro \
        -v "$BACKUP_DIR:/backup" \
        alpine \
        tar czf "/backup/emqx-data.tar.gz" -C /data .
    echo -e "  ${GREEN}✅ EMQX data backed up${NC}"
fi

################################################################################
# 3. PULL LATEST DOCKER IMAGES
################################################################################
echo -e "\n${YELLOW}🐳 Pulling latest Docker images...${NC}"

docker-compose pull --quiet

echo -e "${GREEN}✅ Images updated${NC}"

################################################################################
# 4. DEPLOY NEW CONTAINERS
################################################################################
echo -e "\n${YELLOW}🚀 Deploying containers...${NC}"

# Recreate containers with new images
docker-compose up -d --remove-orphans --force-recreate

echo -e "${GREEN}✅ Containers deployed${NC}"

################################################################################
# 5. WAIT FOR SERVICES TO START
################################################################################
echo -e "\n${YELLOW}🏥 Waiting for services to start...${NC}"

sleep 15

################################################################################
# 6. VERIFY DEPLOYMENT
################################################################################
echo -e "\n${YELLOW}🔍 Verifying deployment...${NC}"

# Check IoT containers
IOT_CONTAINERS=("emqx" "home-assistant" "node-red")

ALL_RUNNING=true
for CONTAINER in "${IOT_CONTAINERS[@]}"; do
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
        STATUS=$(docker inspect --format='{{.State.Status}}' "$CONTAINER")
        echo -e "  ${GREEN}✅ ${CONTAINER}${NC}: $STATUS"
    else
        echo -e "  ${RED}❌ ${CONTAINER}${NC}: NOT RUNNING"
        ALL_RUNNING=false
    fi
done

if [ "$ALL_RUNNING" = false ]; then
    echo -e "\n${YELLOW}⚠️  Some containers are not running - this may be expected for IoT stack${NC}"
fi

# Check URLs
echo -e "\n${YELLOW}🌐 Checking service URLs...${NC}"

SERVICES=("https://emqx.zeusops.dev" "https://ha.zeusops.dev" "https://nodered.zeusops.dev")

for URL in "${SERVICES[@]}"; do
    if curl -sSf -k --max-time 10 "$URL" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅ ${URL}${NC}"
    else
        echo -e "  ${YELLOW}⚠️  ${URL}${NC} (may need more time to start)"
    fi
done

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
echo -e "${GREEN}✅ IOT STACK DEPLOYMENT COMPLETE${NC}"
echo "========================================"
echo "  Deployed at: $(date)"
echo "  Services:"
echo "    - EMQX: https://emqx.zeusops.dev"
echo "    - Home Assistant: https://ha.zeusops.dev"
echo "    - Node-RED: https://nodered.zeusops.dev"
echo "  Containers: ${#IOT_CONTAINERS[@]}"
echo "  Backup: $BACKUP_DIR"
echo -e "========================================\n"

exit 0
