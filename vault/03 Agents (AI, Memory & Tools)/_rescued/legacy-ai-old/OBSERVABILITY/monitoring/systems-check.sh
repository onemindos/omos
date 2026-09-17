#!/bin/bash
################################################################################
# LEGACY AI SYSTEMS HEALTH CHECK
# Purpose: Comprehensive health monitoring for all infrastructure components
# Usage: ./scripts/monitoring/systems-check.sh
# Output: Terminal report + JSON file in /tmp/
################################################################################

set -e

# Color codes for terminal output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
REPORT_FILE="/tmp/systems-check-${TIMESTAMP}.json"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔍 LEGACY AI SYSTEMS CHECK - $(date)"
echo "========================================"

# Initialize alerts array
declare -a ALERTS

################################################################################
# 1. AWS VM HEALTH
################################################################################
echo -e "\n${BLUE}📊 AWS VM STATUS:${NC}"

# Disk space
echo -e "\n${YELLOW}Disk Usage:${NC}"
df -h | grep -E '(Filesystem|/$|/var|/home)'

# Memory
echo -e "\n${YELLOW}Memory Usage:${NC}"
free -h

# System uptime
echo -e "\n${YELLOW}System Uptime:${NC}"
uptime

# Failed systemd units
echo -e "\n${YELLOW}Failed Services:${NC}"
FAILED_SERVICES=$(systemctl list-units --state=failed --no-pager --no-legend | wc -l)
if [ "$FAILED_SERVICES" -gt 0 ]; then
    systemctl list-units --state=failed --no-pager
    ALERTS+=("$FAILED_SERVICES systemd service(s) failed")
else
    echo "✅ No failed services"
fi

# Check disk space warnings (>80%)
DISK_USAGE=$(df -h / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 80 ]; then
    ALERTS+=("Disk usage critical: ${DISK_USAGE}%")
fi

# Check memory usage warnings (>85%)
MEMORY_USAGE=$(free | grep Mem | awk '{printf "%.0f", ($3/$2) * 100}')
if [ "$MEMORY_USAGE" -gt 85 ]; then
    ALERTS+=("Memory usage high: ${MEMORY_USAGE}%")
fi

################################################################################
# 2. DOCKER CONTAINERS
################################################################################
echo -e "\n${BLUE}🐳 DOCKER CONTAINERS:${NC}"

# List all containers with status
echo -e "\n${YELLOW}Container Status:${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" --all | head -25

# Count containers
CONTAINERS_TOTAL=$(docker ps -a -q | wc -l)
CONTAINERS_RUNNING=$(docker ps -q | wc -l)
CONTAINERS_STOPPED=$(docker ps -a -q -f status=exited | wc -l)
CONTAINERS_UNHEALTHY=$(docker ps -q -f health=unhealthy | wc -l)

echo -e "\n${YELLOW}Container Summary:${NC}"
echo "  Total: $CONTAINERS_TOTAL"
echo "  Running: $CONTAINERS_RUNNING"
echo "  Stopped: $CONTAINERS_STOPPED"
echo "  Unhealthy: $CONTAINERS_UNHEALTHY"

# Alert on stopped containers that should be running
if [ "$CONTAINERS_STOPPED" -gt 0 ]; then
    echo -e "\n${RED}⚠️  Stopped containers detected:${NC}"
    docker ps -a -f status=exited --format "{{.Names}}\t{{.Status}}"
    ALERTS+=("$CONTAINERS_STOPPED container(s) stopped")
fi

# Alert on unhealthy containers
if [ "$CONTAINERS_UNHEALTHY" -gt 0 ]; then
    echo -e "\n${RED}⚠️  Unhealthy containers detected:${NC}"
    docker ps -f health=unhealthy --format "{{.Names}}\t{{.Status}}"
    ALERTS+=("$CONTAINERS_UNHEALTHY container(s) unhealthy")
fi

# Container resource usage
echo -e "\n${YELLOW}Container Resource Usage:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" | head -25

################################################################################
# 3. APPLICATION HEALTH CHECKS
################################################################################
echo -e "\n${BLUE}🌐 APPLICATION HEALTH CHECKS:${NC}"

# Define all application URLs from DOMAINS-ACCESS-CURRENT.md
declare -A APPS=(
    ["EMQX"]="https://emqx.zeusops.dev"
    ["Home Assistant"]="https://ha.zeusops.dev"
    ["Node-RED"]="https://nodered.zeusops.dev"
    ["n8n"]="https://n8n.zeusops.dev"
    ["Letta"]="https://letta.zeusops.dev"
    ["Portainer"]="https://command.zeusops.dev"
    ["Dashboard"]="https://dashboard.zeusops.dev"
)

APP_HEALTHY=0
APP_UNHEALTHY=0

for APP_NAME in "${!APPS[@]}"; do
    URL="${APPS[$APP_NAME]}"
    
    # Try to connect (allow self-signed certs, follow redirects, 10s timeout)
    if curl -sSf -k -L --max-time 10 "$URL" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✅ ${APP_NAME}${NC}: $URL"
        ((APP_HEALTHY++))
    else
        echo -e "  ${RED}❌ ${APP_NAME}${NC}: $URL (unreachable)"
        ALERTS+=("Application unreachable: $APP_NAME ($URL)")
        ((APP_UNHEALTHY++))
    fi
done

echo -e "\n${YELLOW}Application Summary:${NC}"
echo "  Healthy: $APP_HEALTHY"
echo "  Unhealthy: $APP_UNHEALTHY"

################################################################################
# 4. DOMAIN STATUS
################################################################################
echo -e "\n${BLUE}🌍 DOMAIN STATUS:${NC}"

DOMAINS=("emqx.zeusops.dev" "ha.zeusops.dev" "nodered.zeusops.dev" "n8n.zeusops.dev" "letta.zeusops.dev" "command.zeusops.dev" "dashboard.zeusops.dev")

for DOMAIN in "${DOMAINS[@]}"; do
    # Check DNS resolution
    if dig +short "$DOMAIN" > /dev/null 2>&1; then
        IP=$(dig +short "$DOMAIN" | head -1)
        echo -e "  ${GREEN}✅ ${DOMAIN}${NC} → $IP"
    else
        echo -e "  ${RED}❌ ${DOMAIN}${NC} (DNS resolution failed)"
        ALERTS+=("DNS resolution failed: $DOMAIN")
    fi
done

################################################################################
# 5. SSL CERTIFICATE CHECKS
################################################################################
echo -e "\n${BLUE}🔒 SSL CERTIFICATES:${NC}"

for DOMAIN in "${DOMAINS[@]}"; do
    # Check SSL expiration
    EXPIRY=$(echo | openssl s_client -servername "$DOMAIN" -connect "$DOMAIN:443" 2>/dev/null | openssl x509 -noout -enddate 2>/dev/null | cut -d= -f2)
    
    if [ -n "$EXPIRY" ]; then
        EXPIRY_EPOCH=$(date -j -f "%b %d %T %Y %Z" "$EXPIRY" +%s 2>/dev/null || date -d "$EXPIRY" +%s 2>/dev/null)
        NOW_EPOCH=$(date +%s)
        DAYS_LEFT=$(( ($EXPIRY_EPOCH - $NOW_EPOCH) / 86400 ))
        
        if [ "$DAYS_LEFT" -lt 30 ]; then
            echo -e "  ${RED}⚠️  ${DOMAIN}${NC}: Expires in ${DAYS_LEFT} days ($EXPIRY)"
            ALERTS+=("SSL expiring soon: $DOMAIN ($DAYS_LEFT days)")
        else
            echo -e "  ${GREEN}✅ ${DOMAIN}${NC}: Valid for ${DAYS_LEFT} days"
        fi
    else
        echo -e "  ${RED}❌ ${DOMAIN}${NC}: Could not retrieve SSL certificate"
    fi
done

################################################################################
# 6. DATABASE STATUS
################################################################################
echo -e "\n${BLUE}💾 DATABASE STATUS:${NC}"

# PostgreSQL databases (from DATABASES-CURRENT.md)
POSTGRES_DBS=("n8n-postgres:5432:n8n:n8n")

for DB_INFO in "${POSTGRES_DBS[@]}"; do
    IFS=':' read -r CONTAINER PORT DB_NAME DB_USER <<< "$DB_INFO"
    
    # Try to connect to database
    if docker exec "$CONTAINER" pg_isready -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1; then
        # Get database size
        DB_SIZE=$(docker exec "$CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME'));" 2>/dev/null | xargs)
        echo -e "  ${GREEN}✅ ${CONTAINER}${NC}: $DB_NAME ($DB_SIZE)"
    else
        echo -e "  ${RED}❌ ${CONTAINER}${NC}: $DB_NAME (unreachable)"
        ALERTS+=("Database unreachable: $CONTAINER/$DB_NAME")
    fi
done

# Redis databases
REDIS_DBS=("n8n-redis")

for CONTAINER in "${REDIS_DBS[@]}"; do
    if docker exec "$CONTAINER" redis-cli ping > /dev/null 2>&1; then
        # Get Redis memory usage
        REDIS_MEM=$(docker exec "$CONTAINER" redis-cli info memory 2>/dev/null | grep "used_memory_human" | cut -d: -f2 | tr -d '\r')
        echo -e "  ${GREEN}✅ ${CONTAINER}${NC}: Memory usage: $REDIS_MEM"
    else
        echo -e "  ${RED}❌ ${CONTAINER}${NC}: (unreachable)"
        ALERTS+=("Redis unreachable: $CONTAINER")
    fi
done

################################################################################
# 7. DISK SPACE WARNINGS
################################################################################
echo -e "\n${BLUE}⚠️  DISK SPACE ALERTS:${NC}"

HIGH_DISK=$(df -h | awk 'NR>1 && int($5) > 80 {print $0}')
if [ -n "$HIGH_DISK" ]; then
    echo "$HIGH_DISK"
else
    echo "✅ All partitions below 80% usage"
fi

################################################################################
# 8. RECENT ERRORS
################################################################################
echo -e "\n${BLUE}🚨 RECENT ERRORS (last hour):${NC}"

if command -v journalctl > /dev/null 2>&1; then
    RECENT_ERRORS=$(journalctl --since "1 hour ago" -p err --no-pager -n 20 2>/dev/null)
    if [ -n "$RECENT_ERRORS" ]; then
        echo "$RECENT_ERRORS"
    else
        echo "✅ No errors in the last hour"
    fi
else
    echo "⚠️  journalctl not available"
fi

################################################################################
# 9. DOCKER VOLUMES
################################################################################
echo -e "\n${BLUE}📦 DOCKER VOLUMES:${NC}"

VOLUME_COUNT=$(docker volume ls -q | wc -l)
echo "Total volumes: $VOLUME_COUNT"

# Check for large volumes (>5GB)
echo -e "\n${YELLOW}Large volumes (>5GB):${NC}"
docker system df -v | grep -A 100 "Local Volumes" | awk 'NR>1 && $3 ~ /GB/ && int($3) > 5 {print $1 "\t" $3}'

################################################################################
# 10. BACKUP STATUS
################################################################################
echo -e "\n${BLUE}💾 BACKUP STATUS:${NC}"

# Check if backup directory exists
if [ -d "/var/lib/docker/volumes/legacy-ai-backups" ] || [ -d "/backups" ]; then
    BACKUP_DIR=$([ -d "/var/lib/docker/volumes/legacy-ai-backups" ] && echo "/var/lib/docker/volumes/legacy-ai-backups" || echo "/backups")
    LATEST_BACKUP=$(find "$BACKUP_DIR" -type f -name "*.sql" -o -name "*.tar.gz" 2>/dev/null | sort -r | head -1)
    
    if [ -n "$LATEST_BACKUP" ]; then
        BACKUP_AGE=$(( ($(date +%s) - $(stat -f %m "$LATEST_BACKUP" 2>/dev/null || stat -c %Y "$LATEST_BACKUP")) / 86400 ))
        echo "Latest backup: $(basename "$LATEST_BACKUP")"
        echo "Age: $BACKUP_AGE days"
        
        if [ "$BACKUP_AGE" -gt 2 ]; then
            ALERTS+=("Backup outdated: $BACKUP_AGE days old")
        fi
    else
        echo "⚠️  No backups found"
        ALERTS+=("No backups found")
    fi
else
    echo "⚠️  Backup directory not found"
fi

################################################################################
# 11. GENERATE JSON REPORT
################################################################################

# Convert alerts array to JSON
ALERTS_JSON="["
for i in "${!ALERTS[@]}"; do
    ALERTS_JSON+="\"${ALERTS[$i]}\""
    if [ $i -lt $((${#ALERTS[@]} - 1)) ]; then
        ALERTS_JSON+=","
    fi
done
ALERTS_JSON+="]"

# Determine overall health
if [ ${#ALERTS[@]} -eq 0 ]; then
    OVERALL_HEALTH="healthy"
elif [ ${#ALERTS[@]} -le 2 ]; then
    OVERALL_HEALTH="warning"
else
    OVERALL_HEALTH="critical"
fi

# Create JSON report
cat > "$REPORT_FILE" << EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "overall_health": "$OVERALL_HEALTH",
  "vm_health": {
    "disk_usage_percent": $DISK_USAGE,
    "memory_usage_percent": $MEMORY_USAGE,
    "failed_services": $FAILED_SERVICES
  },
  "containers": {
    "total": $CONTAINERS_TOTAL,
    "running": $CONTAINERS_RUNNING,
    "stopped": $CONTAINERS_STOPPED,
    "unhealthy": $CONTAINERS_UNHEALTHY
  },
  "applications": {
    "healthy": $APP_HEALTHY,
    "unhealthy": $APP_UNHEALTHY
  },
  "alerts": $ALERTS_JSON
}
EOF

################################################################################
# 12. SUMMARY
################################################################################
echo -e "\n========================================"
echo -e "${BLUE}📊 HEALTH SUMMARY:${NC}"
echo -e "  Overall Status: ${OVERALL_HEALTH^^}"
echo -e "  Containers: $CONTAINERS_RUNNING/$CONTAINERS_TOTAL running"
echo -e "  Applications: $APP_HEALTHY/${#APPS[@]} healthy"
echo -e "  Alerts: ${#ALERTS[@]}"

if [ ${#ALERTS[@]} -gt 0 ]; then
    echo -e "\n${RED}⚠️  ACTIVE ALERTS:${NC}"
    printf '  - %s\n' "${ALERTS[@]}"
fi

echo -e "\n✅ Systems check complete"
echo -e "📄 Report saved: $REPORT_FILE"

# Exit with error code if critical
if [ "$OVERALL_HEALTH" == "critical" ]; then
    exit 1
fi

exit 0
