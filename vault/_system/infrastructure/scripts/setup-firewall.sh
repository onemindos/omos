#!/usr/bin/env bash
# =============================================================================
# OneMind OS — Tailscale Firewall Rules
# =============================================================================
# Opens service ports ONLY on the Tailscale (tailscale0) interface.
# Run once after Tailscale is connected:  sudo bash .do/setup-firewall.sh
#
# IMPORTANT: setup-droplet.sh already configures:
#   - UFW default deny incoming / allow outgoing
#   - SSH allowed from anywhere (for initial access)
#   - Fail2ban active on SSH
#
# This script adds rules so services are reachable ONLY from within
# your Tailscale mesh — never from the public internet.
# =============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${GREEN}[firewall]${NC} $1"; }
warn() { echo -e "${YELLOW}[firewall]${NC} $1"; }
err()  { echo -e "${RED}[firewall]${NC} $1"; }

# ─── Preflight ────────────────────────────────────────────────────────────────

if [ "$(id -u)" -ne 0 ]; then
    err "Must run as root: sudo bash .do/setup-firewall.sh"
    exit 1
fi

# Verify Tailscale is up
if ! ip link show tailscale0 &>/dev/null; then
    err "tailscale0 interface not found."
    err "Run 'sudo tailscale up' first, then re-run this script."
    exit 1
fi

TS_IP=$(tailscale ip -4 2>/dev/null || echo "")
if [ -z "$TS_IP" ]; then
    err "Cannot determine Tailscale IP. Is Tailscale connected?"
    exit 1
fi

log "Tailscale interface detected: $TS_IP"
echo ""

# ─── Service Ports ────────────────────────────────────────────────────────────

echo -e "${CYAN}Opening ports on tailscale0 (Tailscale mesh only):${NC}"
echo ""

# Core infrastructure (always needed)
declare -A CORE_PORTS=(
    ["4222"]="NATS client"
    ["7422"]="NATS leaf nodes"
    ["1883"]="MQTT broker"
    ["8222"]="NATS HTTP monitor"
    ["9222"]="NATS WebSocket"
    ["18789"]="OpenClaw gateway"
)

# Data layer (--profile data)
declare -A DATA_PORTS=(
    ["5432"]="PostgreSQL (OmniDB)"
)

# Full stack (--profile full)
declare -A FULL_PORTS=(
    ["5678"]="n8n automation"
    ["7880"]="LiveKit signaling"
    ["7881"]="LiveKit RTC TCP"
    ["7882"]="LiveKit TURN/TLS"
)

apply_rules() {
    local -n PORTS=$1
    local label=$2

    echo -e "  ${CYAN}── $label ──${NC}"
    for port in "${!PORTS[@]}"; do
        desc="${PORTS[$port]}"
        # Allow on tailscale0 only
        ufw allow in on tailscale0 to any port "$port" proto tcp comment "OneMind: $desc" >/dev/null 2>&1
        log "  $port/tcp — $desc"
    done
    echo ""
}

apply_rules CORE_PORTS "Core (NATS + OpenClaw)"
apply_rules DATA_PORTS "Data Layer"
apply_rules FULL_PORTS "Full Stack (n8n + LiveKit)"

# Also allow NATS cluster if running multi-node in the future
ufw allow in on tailscale0 to any port 6222 proto tcp comment "OneMind: NATS cluster" >/dev/null 2>&1
log "  6222/tcp — NATS cluster (future)"

# UDP for LiveKit media
ufw allow in on tailscale0 to any port 7882 proto udp comment "OneMind: LiveKit TURN/UDP" >/dev/null 2>&1
log "  7882/udp — LiveKit TURN media"

echo ""

# ─── Reload ──────────────────────────────────────────────────────────────────

log "Reloading firewall..."
ufw reload >/dev/null 2>&1

echo ""
echo -e "${CYAN}Current UFW Status:${NC}"
ufw status numbered

echo ""
echo "============================================"
echo -e "${GREEN}  Firewall configured!${NC}"
echo ""
echo "  All service ports open on tailscale0 ONLY."
echo "  Public internet can only reach SSH (22)."
echo ""
echo "  To verify from another Tailscale device:"
echo "    curl http://$TS_IP:18789/health   # OpenClaw"
echo "    curl http://$TS_IP:8222/healthz   # NATS"
echo "============================================"
