#!/bin/bash
################################################################################
# LEGACY AI - VM UPGRADE DEPLOYMENT SCRIPT
# Purpose: Execute the complete upgrade on the VM
# Usage: Run this script from your local machine
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
VM_HOST="3.215.29.128"
VM_USER="ubuntu"
SSH_KEY="/Users/zeuslegacy/Documents/LightsailDefaultKey-us-east-1.pem"
LOCAL_REPO="/Users/zeuslegacy/Documents/Legacy-AI"

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║    LEGACY AI - VM UPGRADE DEPLOYMENT (Remote Execution)   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to run command on VM
run_on_vm() {
    ssh -i "$SSH_KEY" "$VM_USER@$VM_HOST" "$@"
}

# Function to copy files to VM
copy_to_vm() {
    scp -i "$SSH_KEY" -r "$1" "$VM_USER@$VM_HOST:$2"
}

################################################################################
# PREFLIGHT CHECKS
################################################################################
echo -e "${CYAN}═══ PREFLIGHT CHECKS ═══${NC}"

echo -e "${YELLOW}🔍 Checking SSH connection...${NC}"
if ! run_on_vm "echo 'Connected successfully'"; then
    echo -e "${RED}❌ Failed to connect to VM${NC}"
    exit 1
fi
echo -e "${GREEN}✅ SSH connection successful${NC}"
echo ""

echo -e "${YELLOW}📊 Current VM status:${NC}"
run_on_vm "echo '=== SYSTEM INFO ===' && uname -a && echo '' && echo '=== UPTIME ===' && uptime && echo '' && echo '=== CONTAINERS ===' && docker ps --format 'table {{.Names}}\t{{.Status}}' | head -10"
echo ""

# Confirm with user
echo -e "${YELLOW}⚠️  This will upgrade your entire VM infrastructure:${NC}"
echo "   • System packages (Ubuntu)"
echo "   • Docker Engine & Compose"
echo "   • All 20 container images"
echo "   • Fix health checks (Node-RED, EMQX)"
echo "   • Optimize Caddy SSL configuration"
echo "   • Enable firewall"
echo ""

# Check if running non-interactively or with AUTO_CONFIRM
if [ -n "$AUTO_CONFIRM" ] || [ ! -t 0 ]; then
    echo -e "${GREEN}✅ Auto-confirming upgrade (non-interactive mode)${NC}"
    echo ""
else
    read -p "Continue with upgrade? (yes/no): " -r
    echo
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        echo -e "${YELLOW}Upgrade cancelled${NC}"
        exit 0
    fi
fi

################################################################################
# UPLOAD SCRIPTS TO VM
################################################################################
echo -e "${CYAN}═══ UPLOADING UPGRADE SCRIPTS ═══${NC}"

echo -e "${YELLOW}📤 Creating remote directory...${NC}"
run_on_vm "mkdir -p ~/upgrade-scripts"

echo -e "${YELLOW}📤 Uploading upgrade scripts...${NC}"
copy_to_vm "$LOCAL_REPO/scripts/vm-upgrade-complete.sh" "~/upgrade-scripts/"
copy_to_vm "$LOCAL_REPO/scripts/fix-health-checks.sh" "~/upgrade-scripts/"
copy_to_vm "$LOCAL_REPO/ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized" "~/upgrade-scripts/"

echo -e "${YELLOW}🔐 Setting execute permissions...${NC}"
run_on_vm "chmod +x ~/upgrade-scripts/*.sh"

echo -e "${GREEN}✅ Scripts uploaded${NC}"
echo ""

################################################################################
# EXECUTE UPGRADE ON VM
################################################################################
echo -e "${CYAN}═══ EXECUTING UPGRADE ON VM ═══${NC}"
echo ""

echo -e "${YELLOW}🚀 Starting upgrade process...${NC}"
echo -e "${YELLOW}   This may take 10-15 minutes${NC}"
echo -e "${YELLOW}   You can monitor progress below${NC}"
echo ""

# Execute the main upgrade script
run_on_vm "sudo bash ~/upgrade-scripts/vm-upgrade-complete.sh" || {
    echo -e "${RED}❌ Upgrade failed! Check logs on VM${NC}"
    echo -e "${YELLOW}To check logs: ssh -i $SSH_KEY $VM_USER@$VM_HOST 'tail -100 /var/log/vm-upgrade-*.log'${NC}"
    exit 1
}

echo ""
echo -e "${GREEN}✅ Main upgrade completed${NC}"
echo ""

################################################################################
# APPLY HEALTH CHECK FIXES
################################################################################
echo -e "${CYAN}═══ APPLYING HEALTH CHECK FIXES ═══${NC}"

run_on_vm "bash ~/upgrade-scripts/fix-health-checks.sh"

echo -e "${GREEN}✅ Health checks updated${NC}"
echo ""

################################################################################
# UPDATE CADDY CONFIGURATION
################################################################################
echo -e "${CYAN}═══ UPDATING CADDY CONFIGURATION ═══${NC}"

echo -e "${YELLOW}📝 Backing up current Caddyfile...${NC}"
run_on_vm "sudo cp /opt/legacy-ai/infrastructure/caddy/Caddyfile /opt/legacy-ai/infrastructure/caddy/Caddyfile.backup-\$(date +%Y%m%d)"

echo -e "${YELLOW}📝 Deploying optimized Caddyfile...${NC}"
run_on_vm "sudo cp ~/upgrade-scripts/Caddyfile.optimized /opt/legacy-ai/infrastructure/caddy/Caddyfile"

echo -e "${YELLOW}🔄 Reloading Caddy...${NC}"
run_on_vm "docker exec caddy caddy reload --config /etc/caddy/Caddyfile"

echo -e "${GREEN}✅ Caddy configuration updated${NC}"
echo ""

################################################################################
# FINAL VALIDATION
################################################################################
echo -e "${CYAN}═══ FINAL VALIDATION ═══${NC}"
echo ""

echo -e "${YELLOW}📊 Container status:${NC}"
run_on_vm "docker ps --format 'table {{.Names}}\t{{.Status}}'"
echo ""

echo -e "${YELLOW}🌐 Testing domains:${NC}"
for domain in letta.zeusops.dev command.zeusops.dev n8n.zeusops.dev; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "https://$domain" --max-time 5 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "307" ] || [ "$HTTP_CODE" = "301" ]; then
        echo -e "${GREEN}  ✅ $domain - HTTP $HTTP_CODE${NC}"
    else
        echo -e "${YELLOW}  ⚠️  $domain - HTTP $HTTP_CODE${NC}"
    fi
done
echo ""

echo -e "${YELLOW}💾 System resources:${NC}"
run_on_vm "echo 'Memory:' && free -h && echo '' && echo 'Disk:' && df -h /"
echo ""

################################################################################
# GENERATE SUMMARY
################################################################################
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              UPGRADE DEPLOYMENT COMPLETE                   ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${CYAN}📋 Summary:${NC}"
run_on_vm "LATEST_BACKUP=\$(ls -td ~/backups/pre-upgrade-* 2>/dev/null | head -1) && cat \${LATEST_BACKUP}/UPGRADE_REPORT.txt 2>/dev/null | tail -50 || echo 'Report not found'"
echo ""

echo -e "${CYAN}📍 Next Steps:${NC}"
echo "   1. ✅ Review the upgrade report above"
echo "   2. ✅ Test all applications in browser"
echo "   3. ✅ Monitor for any unhealthy containers"
echo "   4. 📋 Update repository documentation"
echo ""

echo -e "${CYAN}📂 Files on VM:${NC}"
echo "   • Upgrade scripts: ~/upgrade-scripts/"
echo "   • Backup location: Run on VM: ls -la ~/backups/"
echo "   • Logs: /var/log/vm-upgrade-*.log"
echo ""

echo -e "${YELLOW}🔗 Quick links:${NC}"
echo "   • Portainer: https://command.zeusops.dev"
echo "   • n8n: https://n8n.zeusops.dev"
echo "   • Dify: https://dify.zeusops.dev"
echo "   • Letta: https://letta.zeusops.dev"
echo ""

echo -e "${GREEN}✅ All done!${NC}"
