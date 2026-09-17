#!/bin/bash
################################################################################
# FIX HEALTH CHECKS FOR NODE-RED AND EMQX
# Purpose: Update docker compose.yml files with proper health checks
################################################################################

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 Fixing Health Checks for Node-RED and EMQX${NC}"
echo ""

################################################################################
# FIX NODE-RED HEALTH CHECK
################################################################################
fix_nodered_health() {
    echo -e "${YELLOW}📝 Fixing Node-RED health check...${NC}"
    
    NODERED_COMPOSE="/opt/legacy-ai/stacks/iot/node-red/docker compose.yml"
    
    if [ ! -f "$NODERED_COMPOSE" ]; then
        echo -e "${YELLOW}⚠️  Node-RED docker compose.yml not found at expected location${NC}"
        # Try alternative locations
        NODERED_COMPOSE=$(find /opt/legacy-ai -name docker compose.yml -path "*/node-red/*" 2>/dev/null | head -1)
        if [ -z "$NODERED_COMPOSE" ]; then
            echo -e "${YELLOW}⚠️  Could not locate Node-RED docker compose.yml${NC}"
            return 1
        fi
    fi
    
    echo "Found: $NODERED_COMPOSE"
    
    # Backup original
    cp "$NODERED_COMPOSE" "${NODERED_COMPOSE}.backup-$(date +%Y%m%d)"
    
    # Check if health check exists
    if grep -q "healthcheck:" "$NODERED_COMPOSE"; then
        echo "  → Updating existing health check"
        # Update the health check section
        sed -i '/healthcheck:/,/start_period:/c\    healthcheck:\n      test: ["CMD", "curl", "-f", "http://localhost:1880/"]\n      interval: 30s\n      timeout: 10s\n      retries: 3\n      start_period: 40s' "$NODERED_COMPOSE"
    else
        echo "  → Adding health check"
        # Add health check before networks or volumes section
        sed -i '/^    networks:/i\    healthcheck:\n      test: ["CMD", "curl", "-f", "http://localhost:1880/"]\n      interval: 30s\n      timeout: 10s\n      retries: 3\n      start_period: 40s' "$NODERED_COMPOSE"
    fi
    
    echo -e "${GREEN}✅ Node-RED health check configured${NC}"
}

################################################################################
# FIX EMQX HEALTH CHECK
################################################################################
fix_emqx_health() {
    echo -e "${YELLOW}📝 Fixing EMQX health check...${NC}"
    
    EMQX_COMPOSE="/opt/legacy-ai/stacks/iot/emqx/docker compose.yml"
    
    if [ ! -f "$EMQX_COMPOSE" ]; then
        echo -e "${YELLOW}⚠️  EMQX docker compose.yml not found at expected location${NC}"
        # Try alternative locations
        EMQX_COMPOSE=$(find /opt/legacy-ai -name docker compose.yml -path "*/emqx/*" 2>/dev/null | head -1)
        if [ -z "$EMQX_COMPOSE" ]; then
            echo -e "${YELLOW}⚠️  Could not locate EMQX docker compose.yml${NC}"
            return 1
        fi
    fi
    
    echo "Found: $EMQX_COMPOSE"
    
    # Backup original
    cp "$EMQX_COMPOSE" "${EMQX_COMPOSE}.backup-$(date +%Y%m%d)"
    
    # Check if health check exists
    if grep -q "healthcheck:" "$EMQX_COMPOSE"; then
        echo "  → Updating existing health check"
        # Update the health check section
        sed -i '/healthcheck:/,/start_period:/c\    healthcheck:\n      test: ["CMD", "/opt/emqx/bin/emqx", "ctl", "status"]\n      interval: 30s\n      timeout: 10s\n      retries: 5\n      start_period: 60s' "$EMQX_COMPOSE"
    else
        echo "  → Adding health check"
        # Add health check before networks or volumes section
        sed -i '/^    networks:/i\    healthcheck:\n      test: ["CMD", "/opt/emqx/bin/emqx", "ctl", "status"]\n      interval: 30s\n      timeout: 10s\n      retries: 5\n      start_period: 60s' "$EMQX_COMPOSE"
    fi
    
    echo -e "${GREEN}✅ EMQX health check configured${NC}"
}

################################################################################
# RESTART CONTAINERS
################################################################################
restart_containers() {
    echo -e "${YELLOW}🔄 Restarting containers to apply changes...${NC}"
    
    # Restart Node-RED
    if docker ps -a | grep -q node-red; then
        echo "  → Restarting Node-RED..."
        NODERED_DIR=$(find /opt/legacy-ai -name docker compose.yml -path "*/node-red/*" -exec dirname {} \; 2>/dev/null | head -1)
        if [ -n "$NODERED_DIR" ]; then
            cd "$NODERED_DIR"
            docker compose up -d --force-recreate
            echo -e "${GREEN}  ✅ Node-RED restarted${NC}"
        fi
    fi
    
    # Restart EMQX
    if docker ps -a | grep -q emqx; then
        echo "  → Restarting EMQX..."
        EMQX_DIR=$(find /opt/legacy-ai -name docker compose.yml -path "*/emqx/*" -exec dirname {} \; 2>/dev/null | head -1)
        if [ -n "$EMQX_DIR" ]; then
            cd "$EMQX_DIR"
            docker compose up -d --force-recreate
            echo -e "${GREEN}  ✅ EMQX restarted${NC}"
        fi
    fi
    
    echo ""
    echo -e "${YELLOW}⏳ Waiting 30 seconds for health checks to stabilize...${NC}"
    sleep 30
}

################################################################################
# VERIFY HEALTH
################################################################################
verify_health() {
    echo -e "${YELLOW}🔍 Verifying health status...${NC}"
    echo ""
    
    # Check Node-RED
    if docker ps | grep -q node-red; then
        NODERED_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' node-red 2>/dev/null || echo "no-healthcheck")
        echo "Node-RED health: $NODERED_HEALTH"
        
        if [ "$NODERED_HEALTH" = "healthy" ]; then
            echo -e "${GREEN}  ✅ Node-RED is healthy${NC}"
        elif [ "$NODERED_HEALTH" = "starting" ]; then
            echo -e "${YELLOW}  ⏳ Node-RED is starting (check again in a minute)${NC}"
        else
            echo -e "${YELLOW}  ⚠️  Node-RED health: $NODERED_HEALTH${NC}"
            echo "  Logs:"
            docker logs --tail 20 node-red
        fi
    fi
    
    echo ""
    
    # Check EMQX
    if docker ps | grep -q emqx; then
        EMQX_HEALTH=$(docker inspect --format='{{.State.Health.Status}}' emqx 2>/dev/null || echo "no-healthcheck")
        echo "EMQX health: $EMQX_HEALTH"
        
        if [ "$EMQX_HEALTH" = "healthy" ]; then
            echo -e "${GREEN}  ✅ EMQX is healthy${NC}"
        elif [ "$EMQX_HEALTH" = "starting" ]; then
            echo -e "${YELLOW}  ⏳ EMQX is starting (check again in a minute)${NC}"
        else
            echo -e "${YELLOW}  ⚠️  EMQX health: $EMQX_HEALTH${NC}"
            echo "  Logs:"
            docker logs --tail 20 emqx
        fi
    fi
    
    echo ""
    echo -e "${BLUE}💡 Tip: Check health again in 1-2 minutes with:${NC}"
    echo "   docker ps --format 'table {{.Names}}\t{{.Status}}'"
}

################################################################################
# MAIN
################################################################################
main() {
    fix_nodered_health
    echo ""
    fix_emqx_health
    echo ""
    restart_containers
    verify_health
    
    echo ""
    echo -e "${GREEN}✅ Health check fixes applied${NC}"
    echo ""
    echo -e "${BLUE}📋 Next steps:${NC}"
    echo "   1. Monitor containers for 2-3 minutes"
    echo "   2. Check: docker ps --filter 'health=unhealthy'"
    echo "   3. View logs if still unhealthy: docker logs <container>"
}

main "$@"
