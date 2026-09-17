#!/bin/bash
################################################################################
# LEGACY AI QUICK FIX - AUTO-REMEDIATION
# Purpose: Automatically fix common infrastructure issues
# Usage: ./scripts/monitoring/quick-fix.sh
# Run: Manually or via cron (hourly recommended)
################################################################################

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
LOG_FILE="/var/log/legacy-ai-quick-fix-${TIMESTAMP}.log"

echo "🔧 RUNNING AUTOMATED FIXES - $(date)" | tee -a "$LOG_FILE"
echo "========================================" | tee -a "$LOG_FILE"

FIXES_APPLIED=0

################################################################################
# 1. RESTART UNHEALTHY CONTAINERS
################################################################################
echo -e "\n${BLUE}🏥 Checking container health...${NC}" | tee -a "$LOG_FILE"

# Get unhealthy containers (if Docker supports health checks)
UNHEALTHY_CONTAINERS=$(docker ps --filter health=unhealthy -q 2>/dev/null || echo "")

if [ -n "$UNHEALTHY_CONTAINERS" ]; then
    for CONTAINER_ID in $UNHEALTHY_CONTAINERS; do
        CONTAINER_NAME=$(docker inspect --format='{{.Name}}' "$CONTAINER_ID" | sed 's/\///')
        echo -e "${YELLOW}⚠️  Restarting unhealthy container: ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
        
        if docker restart "$CONTAINER_ID" >> "$LOG_FILE" 2>&1; then
            echo -e "${GREEN}✅ Successfully restarted ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
            ((FIXES_APPLIED++))
        else
            echo -e "${RED}❌ Failed to restart ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
        fi
    done
else
    echo "✅ No unhealthy containers found" | tee -a "$LOG_FILE"
fi

################################################################################
# 2. RESTART EXITED CONTAINERS (with restart policy)
################################################################################
echo -e "\n${BLUE}🔄 Checking exited containers...${NC}" | tee -a "$LOG_FILE"

# Get containers that should be running (have restart policy)
EXITED_CONTAINERS=$(docker ps -a --filter status=exited --format '{{.Names}}' 2>/dev/null || echo "")

if [ -n "$EXITED_CONTAINERS" ]; then
    for CONTAINER_NAME in $EXITED_CONTAINERS; do
        # Check restart policy
        RESTART_POLICY=$(docker inspect --format='{{.HostConfig.RestartPolicy.Name}}' "$CONTAINER_NAME" 2>/dev/null || echo "no")
        
        if [[ "$RESTART_POLICY" != "no" ]]; then
            echo -e "${YELLOW}⚠️  Restarting exited container: ${CONTAINER_NAME} (policy: ${RESTART_POLICY})${NC}" | tee -a "$LOG_FILE"
            
            if docker start "$CONTAINER_NAME" >> "$LOG_FILE" 2>&1; then
                echo -e "${GREEN}✅ Successfully started ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
                ((FIXES_APPLIED++))
            else
                echo -e "${RED}❌ Failed to start ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
            fi
        fi
    done
else
    echo "✅ No exited containers found" | tee -a "$LOG_FILE"
fi

################################################################################
# 3. CHECK AND RESTART CONTAINERS WITH HIGH RESTART COUNTS
################################################################################
echo -e "\n${BLUE}🔃 Checking for crashlooping containers...${NC}" | tee -a "$LOG_FILE"

# Find containers that have restarted more than 10 times recently
CRASHLOOP_CONTAINERS=$(docker ps --format '{{.Names}}\t{{.Status}}' | awk -F'\t' '$2 ~ /Restarting/ || $2 ~ /\(/ && int($2) > 10 {print $1}' || echo "")

if [ -n "$CRASHLOOP_CONTAINERS" ]; then
    for CONTAINER_NAME in $CRASHLOOP_CONTAINERS; do
        echo -e "${RED}⚠️  Crashlooping container detected: ${CONTAINER_NAME}${NC}" | tee -a "$LOG_FILE"
        echo "   Saving logs before intervention..." | tee -a "$LOG_FILE"
        
        # Save last 100 lines of logs
        docker logs --tail 100 "$CONTAINER_NAME" > "/tmp/${CONTAINER_NAME}-crash-${TIMESTAMP}.log" 2>&1
        
        echo "   Logs saved to: /tmp/${CONTAINER_NAME}-crash-${TIMESTAMP}.log" | tee -a "$LOG_FILE"
        echo "   Manual intervention recommended - NOT auto-restarting crashloop" | tee -a "$LOG_FILE"
    fi
else
    echo "✅ No crashlooping containers found" | tee -a "$LOG_FILE"
fi

################################################################################
# 4. CLEAN LARGE LOG FILES (>1GB)
################################################################################
echo -e "\n${BLUE}🗑️  Checking for large log files...${NC}" | tee -a "$LOG_FILE"

# Find large log files
LARGE_LOGS=$(find /var/log -type f -size +1G 2>/dev/null || echo "")

if [ -n "$LARGE_LOGS" ]; then
    echo "$LARGE_LOGS" | while read -r LOGFILE; do
        LOGSIZE=$(du -h "$LOGFILE" | cut -f1)
        echo -e "${YELLOW}⚠️  Truncating large log: ${LOGFILE} (${LOGSIZE})${NC}" | tee -a "$LOG_FILE"
        
        # Backup last 1000 lines before truncating
        tail -1000 "$LOGFILE" > "${LOGFILE}.backup-${TIMESTAMP}" 2>/dev/null || true
        
        # Truncate
        if truncate -s 0 "$LOGFILE" 2>/dev/null; then
            echo -e "${GREEN}✅ Truncated ${LOGFILE}${NC}" | tee -a "$LOG_FILE"
            ((FIXES_APPLIED++))
        else
            echo -e "${RED}❌ Failed to truncate ${LOGFILE} (permission denied?)${NC}" | tee -a "$LOG_FILE"
        fi
    done
else
    echo "✅ No large log files (>1GB) found" | tee -a "$LOG_FILE"
fi

################################################################################
# 5. CLEAN DOCKER RESOURCES
################################################################################
echo -e "\n${BLUE}🐳 Cleaning unused Docker resources...${NC}" | tee -a "$LOG_FILE"

# Clean dangling images, stopped containers, unused networks (keep volumes)
if docker system prune -f --volumes=false >> "$LOG_FILE" 2>&1; then
    echo -e "${GREEN}✅ Docker cleanup complete${NC}" | tee -a "$LOG_FILE"
    ((FIXES_APPLIED++))
else
    echo -e "${RED}❌ Docker cleanup failed${NC}" | tee -a "$LOG_FILE"
fi

################################################################################
# 6. CHECK DISK SPACE AND ALERT IF CRITICAL
################################################################################
echo -e "\n${BLUE}💾 Checking disk space...${NC}" | tee -a "$LOG_FILE"

DISK_USAGE=$(df -h / | tail -1 | awk '{print $5}' | sed 's/%//')

if [ "$DISK_USAGE" -gt 90 ]; then
    echo -e "${RED}🚨 CRITICAL: Disk usage at ${DISK_USAGE}%${NC}" | tee -a "$LOG_FILE"
    
    # Find largest directories
    echo "Largest directories in /var:" | tee -a "$LOG_FILE"
    du -sh /var/* 2>/dev/null | sort -rh | head -10 | tee -a "$LOG_FILE"
    
    # Clean Docker build cache if critical
    echo -e "${YELLOW}⚠️  Cleaning Docker build cache...${NC}" | tee -a "$LOG_FILE"
    docker builder prune -f >> "$LOG_FILE" 2>&1 || true
    ((FIXES_APPLIED++))
    
elif [ "$DISK_USAGE" -gt 80 ]; then
    echo -e "${YELLOW}⚠️  Warning: Disk usage at ${DISK_USAGE}%${NC}" | tee -a "$LOG_FILE"
else
    echo -e "${GREEN}✅ Disk usage healthy: ${DISK_USAGE}%${NC}" | tee -a "$LOG_FILE"
fi

################################################################################
# 7. CHECK MEMORY AND RESTART MEMORY-HUNGRY CONTAINERS
################################################################################
echo -e "\n${BLUE}💾 Checking memory usage...${NC}" | tee -a "$LOG_FILE"

MEMORY_USAGE=$(free | grep Mem | awk '{printf "%.0f", ($3/$2) * 100}')

if [ "$MEMORY_USAGE" -gt 90 ]; then
    echo -e "${RED}🚨 CRITICAL: Memory usage at ${MEMORY_USAGE}%${NC}" | tee -a "$LOG_FILE"
    
    # Find memory-hungry containers
    echo "Top memory-consuming containers:" | tee -a "$LOG_FILE"
    docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}\t{{.MemPerc}}" | sort -k3 -rh | head -5 | tee -a "$LOG_FILE"
    
    # Get the most memory-hungry container (excluding databases)
    HUNGRY_CONTAINER=$(docker stats --no-stream --format "{{.Name}}\t{{.MemPerc}}" | grep -v -E "(postgres|redis|mongo)" | sort -k2 -rh | head -1 | cut -f1)
    
    if [ -n "$HUNGRY_CONTAINER" ]; then
        echo -e "${YELLOW}⚠️  Considering restart of ${HUNGRY_CONTAINER} to free memory${NC}" | tee -a "$LOG_FILE"
        echo "   Manual intervention recommended for memory issues" | tee -a "$LOG_FILE"
    fi
    
elif [ "$MEMORY_USAGE" -gt 85 ]; then
    echo -e "${YELLOW}⚠️  Warning: Memory usage at ${MEMORY_USAGE}%${NC}" | tee -a "$LOG_FILE"
else
    echo -e "${GREEN}✅ Memory usage healthy: ${MEMORY_USAGE}%${NC}" | tee -a "$LOG_FILE"
fi

################################################################################
# 8. CHECK AND FIX COMMON PERMISSION ISSUES
################################################################################
echo -e "\n${BLUE}🔐 Checking common permission issues...${NC}" | tee -a "$LOG_FILE"

# Check Docker socket permissions
if [ -S /var/run/docker.sock ]; then
    DOCKER_SOCK_PERMS=$(stat -c %a /var/run/docker.sock 2>/dev/null || stat -f %A /var/run/docker.sock 2>/dev/null)
    
    if [ "$DOCKER_SOCK_PERMS" != "666" ] && [ "$DOCKER_SOCK_PERMS" != "660" ]; then
        echo -e "${YELLOW}⚠️  Docker socket has unusual permissions: ${DOCKER_SOCK_PERMS}${NC}" | tee -a "$LOG_FILE"
    else
        echo "✅ Docker socket permissions OK" | tee -a "$LOG_FILE"
    fi
fi

# Check log directory permissions
if [ -d "/var/log/legacy-ai-backups" ]; then
    if [ ! -w "/var/log/legacy-ai-backups" ]; then
        echo -e "${YELLOW}⚠️  Backup log directory not writable${NC}" | tee -a "$LOG_FILE"
    else
        echo "✅ Backup log directory writable" | tee -a "$LOG_FILE"
    fi
fi

################################################################################
# 9. VERIFY CRITICAL SERVICES ARE RUNNING
################################################################################
echo -e "\n${BLUE}✅ Verifying critical services...${NC}" | tee -a "$LOG_FILE"

# Critical containers that should always be running
CRITICAL_CONTAINERS=("portainer" "caddy" "n8n-main" "n8n-postgres")

for CONTAINER in "${CRITICAL_CONTAINERS[@]}"; do
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
        echo "  ✅ ${CONTAINER} is running" | tee -a "$LOG_FILE"
    else
        echo -e "  ${RED}❌ ${CONTAINER} is NOT running${NC}" | tee -a "$LOG_FILE"
        
        # Attempt to start
        if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
            echo "     Attempting to start ${CONTAINER}..." | tee -a "$LOG_FILE"
            if docker start "$CONTAINER" >> "$LOG_FILE" 2>&1; then
                echo -e "     ${GREEN}✅ Started ${CONTAINER}${NC}" | tee -a "$LOG_FILE"
                ((FIXES_APPLIED++))
            else
                echo -e "     ${RED}❌ Failed to start ${CONTAINER}${NC}" | tee -a "$LOG_FILE"
            fi
        fi
    fi
done

################################################################################
# 10. CHECK DOCKER DAEMON
################################################################################
echo -e "\n${BLUE}🐋 Checking Docker daemon...${NC}" | tee -a "$LOG_FILE"

if docker info > /dev/null 2>&1; then
    echo "✅ Docker daemon is healthy" | tee -a "$LOG_FILE"
else
    echo -e "${RED}❌ Docker daemon is not responding${NC}" | tee -a "$LOG_FILE"
    echo "   Manual intervention required: sudo systemctl restart docker" | tee -a "$LOG_FILE"
fi

################################################################################
# 11. CLEAN OLD LOG FILES
################################################################################
echo -e "\n${BLUE}🗑️  Cleaning old log files...${NC}" | tee -a "$LOG_FILE"

# Remove quick-fix logs older than 7 days
OLD_LOGS=$(find /var/log -name "legacy-ai-quick-fix-*.log" -mtime +7 2>/dev/null || echo "")

if [ -n "$OLD_LOGS" ]; then
    echo "$OLD_LOGS" | while read -r OLD_LOG; do
        rm -f "$OLD_LOG" 2>/dev/null && echo "  Removed: $(basename "$OLD_LOG")" | tee -a "$LOG_FILE"
    done
    ((FIXES_APPLIED++))
else
    echo "✅ No old log files to clean" | tee -a "$LOG_FILE"
fi

################################################################################
# 12. SUMMARY
################################################################################
echo -e "\n========================================" | tee -a "$LOG_FILE"
echo -e "${BLUE}🔧 QUICK FIX SUMMARY:${NC}" | tee -a "$LOG_FILE"
echo "  Fixes applied: $FIXES_APPLIED" | tee -a "$LOG_FILE"
echo "  Timestamp: $(date)" | tee -a "$LOG_FILE"
echo "  Log file: $LOG_FILE" | tee -a "$LOG_FILE"

if [ "$FIXES_APPLIED" -gt 0 ]; then
    echo -e "\n${GREEN}✅ Quick fixes complete - $FIXES_APPLIED action(s) taken${NC}" | tee -a "$LOG_FILE"
else
    echo -e "\n${GREEN}✅ System healthy - no fixes needed${NC}" | tee -a "$LOG_FILE"
fi

# Run a quick systems check after fixes
echo -e "\n${BLUE}Running post-fix health check...${NC}" | tee -a "$LOG_FILE"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "$SCRIPT_DIR/systems-check.sh" ]; then
    # Run systems check but don't fail if it returns error
    bash "$SCRIPT_DIR/systems-check.sh" >> "$LOG_FILE" 2>&1 || true
fi

exit 0
