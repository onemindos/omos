#!/bin/bash
################################################################################
# LEGACY AI CRON JOB SETUP
# Purpose: Configure all automated maintenance tasks
# Usage: sudo ./scripts/maintenance/setup-cron.sh
################################################################################

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}⏰ SETTING UP AUTOMATED TASKS${NC}"
echo "========================================"

# Determine script directory (absolute path)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MONITORING_DIR="$SCRIPT_DIR/scripts/monitoring"
MAINTENANCE_DIR="$SCRIPT_DIR/scripts/maintenance"

echo "Repository root: $SCRIPT_DIR"

# Check if scripts exist
if [ ! -f "$MONITORING_DIR/systems-check.sh" ]; then
    echo -e "${RED}❌ systems-check.sh not found${NC}"
    exit 1
fi

if [ ! -f "$MONITORING_DIR/quick-fix.sh" ]; then
    echo -e "${RED}❌ quick-fix.sh not found${NC}"
    exit 1
fi

if [ ! -f "$MAINTENANCE_DIR/backup-all.sh" ]; then
    echo -e "${RED}❌ backup-all.sh not found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All scripts found${NC}"

################################################################################
# CREATE LOG DIRECTORIES
################################################################################
echo -e "\n${YELLOW}📁 Creating log directories...${NC}"

LOG_DIRS=(
    "/var/log/legacy-ai-systems"
    "/var/log/legacy-ai-backups"
    "/var/log/legacy-ai-alerts"
)

for DIR in "${LOG_DIRS[@]}"; do
    if sudo mkdir -p "$DIR" 2>/dev/null; then
        sudo chmod 755 "$DIR"
        echo "  ✅ Created: $DIR"
    else
        # If sudo fails, try without sudo (for local testing)
        mkdir -p "$DIR" 2>/dev/null || echo "  ⚠️  Could not create: $DIR"
    fi
done

################################################################################
# BACKUP EXISTING CRONTAB
################################################################################
echo -e "\n${YELLOW}💾 Backing up existing crontab...${NC}"

if crontab -l > /dev/null 2>&1; then
    BACKUP_FILE="/tmp/crontab-backup-$(date +%Y%m%d-%H%M%S).txt"
    crontab -l > "$BACKUP_FILE" 2>/dev/null || true
    echo "Existing crontab backed up to: $BACKUP_FILE"
else
    echo "No existing crontab found"
fi

################################################################################
# CREATE NEW CRONTAB
################################################################################
echo -e "\n${YELLOW}⏰ Configuring cron jobs...${NC}"

# Create temp file with new cron jobs
TEMP_CRON=$(mktemp)

cat > "$TEMP_CRON" << EOF
# ============================================================================
# LEGACY AI AUTOMATED TASKS
# Generated: $(date)
# ============================================================================

SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin

# ============================================================================
# HEALTH MONITORING
# ============================================================================

# Systems health check - Every 6 hours
0 */6 * * * $MONITORING_DIR/systems-check.sh >> /var/log/legacy-ai-systems/systems-check.log 2>&1

# Quick fixes - Every hour
0 * * * * $MONITORING_DIR/quick-fix.sh >> /var/log/legacy-ai-systems/quick-fix.log 2>&1

# ============================================================================
# BACKUPS
# ============================================================================

# Comprehensive backup - Daily at 2 AM
0 2 * * * $MAINTENANCE_DIR/backup-all.sh /backups >> /var/log/legacy-ai-backups/backup-all.log 2>&1

# Database-only backup - Every 12 hours
0 */12 * * * $MAINTENANCE_DIR/backup-databases.sh >> /var/log/legacy-ai-backups/backup-databases.log 2>&1

# ============================================================================
# CLEANUP & MAINTENANCE
# ============================================================================

# Docker cleanup - Weekly (Sundays at 3 AM)
0 3 * * 0 docker system prune -af --volumes=false >> /var/log/legacy-ai-systems/docker-cleanup.log 2>&1

# Docker image cleanup - Weekly (Sundays at 3:30 AM)
30 3 * * 0 docker image prune -af >> /var/log/legacy-ai-systems/docker-image-cleanup.log 2>&1

# Log rotation - Weekly (Sundays at 4 AM)
0 4 * * 0 find /var/log/legacy-ai-* -type f -name "*.log" -mtime +30 -delete 2>&1

# Temp file cleanup - Daily at 4 AM
0 4 * * * find /tmp -type f -name "legacy-ai-*" -mtime +7 -delete 2>&1
0 4 * * * find /tmp -type f -name "*-backup-*" -mtime +7 -delete 2>&1

# ============================================================================
# SSL CERTIFICATES
# ============================================================================

# SSL certificate renewal check - Daily at 6 AM
0 6 * * * docker exec caddy caddy reload --config /etc/caddy/Caddyfile >> /var/log/legacy-ai-systems/ssl-renewal.log 2>&1 || true

# ============================================================================
# MONITORING & ALERTS
# ============================================================================

# Disk space alert - Every 6 hours
0 */6 * * * DISK_USAGE=\$(df -h / | tail -1 | awk '{print \$5}' | sed 's/%//'); if [ "\$DISK_USAGE" -gt 85 ]; then $MONITORING_DIR/send-alert.sh warning "Disk usage at \${DISK_USAGE}%" "Root partition is running low on space"; fi

# Memory alert - Every 6 hours
30 */6 * * * MEM_USAGE=\$(free | grep Mem | awk '{printf "%.0f", (\$3/\$2) * 100}'); if [ "\$MEM_USAGE" -gt 85 ]; then $MONITORING_DIR/send-alert.sh warning "Memory usage at \${MEM_USAGE}%" "System memory is running high"; fi

# Container health alert - Every hour
45 * * * * UNHEALTHY=\$(docker ps --filter health=unhealthy -q | wc -l); if [ "\$UNHEALTHY" -gt 0 ]; then $MONITORING_DIR/send-alert.sh warning "\$UNHEALTHY unhealthy container(s)" "Some containers are reporting unhealthy status"; fi

# ============================================================================
# END OF CRONTAB
# ============================================================================
EOF

# Display the crontab
echo -e "\n${BLUE}Cron jobs to be installed:${NC}"
cat "$TEMP_CRON"

################################################################################
# INSTALL CRONTAB
################################################################################
echo -e "\n${YELLOW}Installing crontab...${NC}"

if crontab "$TEMP_CRON"; then
    echo -e "${GREEN}✅ Crontab installed successfully${NC}"
else
    echo -e "${RED}❌ Failed to install crontab${NC}"
    rm "$TEMP_CRON"
    exit 1
fi

rm "$TEMP_CRON"

################################################################################
# VERIFY INSTALLATION
################################################################################
echo -e "\n${YELLOW}🔍 Verifying installation...${NC}"

CRON_COUNT=$(crontab -l | grep -v "^#" | grep -v "^$" | wc -l)
echo "Active cron jobs: $CRON_COUNT"

################################################################################
# CREATE WRAPPER SCRIPTS IN /usr/local/bin (OPTIONAL)
################################################################################
echo -e "\n${YELLOW}🔗 Creating system-wide wrapper scripts...${NC}"

if [ "$EUID" -eq 0 ] || sudo -n true 2>/dev/null; then
    # Create wrapper scripts for easy execution
    WRAPPERS=(
        "legacy-check:$MONITORING_DIR/systems-check.sh"
        "legacy-fix:$MONITORING_DIR/quick-fix.sh"
        "legacy-backup:$MAINTENANCE_DIR/backup-all.sh"
        "legacy-alert:$MONITORING_DIR/send-alert.sh"
    )
    
    for WRAPPER in "${WRAPPERS[@]}"; do
        IFS=':' read -r CMD SCRIPT <<< "$WRAPPER"
        
        WRAPPER_PATH="/usr/local/bin/$CMD"
        
        if [ -n "$SCRIPT" ] && [ -f "$SCRIPT" ]; then
            sudo tee "$WRAPPER_PATH" > /dev/null << WRAPPER_EOF
#!/bin/bash
# Wrapper for Legacy AI script: $SCRIPT
exec "$SCRIPT" "\$@"
WRAPPER_EOF
            
            sudo chmod +x "$WRAPPER_PATH"
            echo "  ✅ Created: $CMD"
        fi
    done
    
    echo ""
    echo -e "${GREEN}System-wide commands available:${NC}"
    echo "  - legacy-check    # Run systems health check"
    echo "  - legacy-fix      # Run quick fixes"
    echo "  - legacy-backup   # Run comprehensive backup"
    echo "  - legacy-alert    # Send alert"
else
    echo "⚠️  Skipping system-wide wrappers (requires sudo)"
fi

################################################################################
# TEST CRON JOBS (DRY RUN)
################################################################################
echo -e "\n${YELLOW}🧪 Testing scripts (dry run)...${NC}"

echo "Testing systems-check.sh..."
if bash -n "$MONITORING_DIR/systems-check.sh"; then
    echo -e "  ${GREEN}✅ Syntax OK${NC}"
else
    echo -e "  ${RED}❌ Syntax error${NC}"
fi

echo "Testing quick-fix.sh..."
if bash -n "$MONITORING_DIR/quick-fix.sh"; then
    echo -e "  ${GREEN}✅ Syntax OK${NC}"
else
    echo -e "  ${RED}❌ Syntax error${NC}"
fi

echo "Testing backup-all.sh..."
if bash -n "$MAINTENANCE_DIR/backup-all.sh"; then
    echo -e "  ${GREEN}✅ Syntax OK${NC}"
else
    echo -e "  ${RED}❌ Syntax error${NC}"
fi

################################################################################
# SUMMARY
################################################################################
echo -e "\n========================================"
echo -e "${GREEN}✅ CRON SETUP COMPLETE${NC}"
echo "========================================"
echo ""
echo "📋 Scheduled Tasks:"
echo "  - Systems check: Every 6 hours"
echo "  - Quick fixes: Every hour"
echo "  - Backups: Daily at 2 AM (full), every 12h (databases)"
echo "  - Docker cleanup: Weekly (Sundays 3 AM)"
echo "  - SSL renewal: Daily at 6 AM"
echo "  - Alerts: Every 6 hours (disk/memory) + hourly (containers)"
echo ""
echo "📁 Log Locations:"
echo "  - Systems: /var/log/legacy-ai-systems/"
echo "  - Backups: /var/log/legacy-ai-backups/"
echo "  - Alerts: /var/log/legacy-ai-alerts/"
echo ""
echo "🔧 Management Commands:"
echo "  - View crontab: crontab -l"
echo "  - Edit crontab: crontab -e"
echo "  - View logs: tail -f /var/log/legacy-ai-systems/*.log"
echo ""
echo -e "${YELLOW}⚠️  Important:${NC}"
echo "  1. Verify backup directory exists: /backups"
echo "  2. Configure alert webhooks (optional):"
echo "     export SLACK_WEBHOOK_URL='https://hooks.slack.com/...'"
echo "  3. Ensure Docker is running on system boot"
echo "  4. Monitor logs for first 48 hours"
echo ""
echo -e "${GREEN}✅ Automation is now active!${NC}"
echo ""

exit 0
