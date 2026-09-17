#!/usr/bin/env bash
# =============================================================================
# OneMind OS — Health Check
# =============================================================================
# Checks all services and reports status.
# Run anytime:  bash .do/health-check.sh
# =============================================================================

set -uo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
WARN=0

pass() { echo -e "  ${GREEN}✓${NC} $1"; ((PASS++)); }
fail() { echo -e "  ${RED}✗${NC} $1"; ((FAIL++)); }
warn() { echo -e "  ${YELLOW}⚠${NC} $1"; ((WARN++)); }
skip() { echo -e "  ${YELLOW}–${NC} $1 (skipped)"; }

REPO_DIR="${REPO_DIR:-/opt/onemind}"

echo ""
echo "============================================"
echo "  OneMind OS — Health Check"
echo "  $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
echo "============================================"
echo ""

# ─── System ──────────────────────────────────────────────────────────────────

echo "System:"
# Disk
DISK_PCT=$(df / | awk 'NR==2 {print $5}' | tr -d '%')
if [ "$DISK_PCT" -lt 80 ]; then
    pass "Disk usage: ${DISK_PCT}%"
elif [ "$DISK_PCT" -lt 90 ]; then
    warn "Disk usage: ${DISK_PCT}% (getting full)"
else
    fail "Disk usage: ${DISK_PCT}% (CRITICAL)"
fi

# Memory
MEM_PCT=$(free | awk '/Mem:/ {printf "%.0f", $3/$2*100}')
if [ "$MEM_PCT" -lt 80 ]; then
    pass "Memory usage: ${MEM_PCT}%"
elif [ "$MEM_PCT" -lt 90 ]; then
    warn "Memory usage: ${MEM_PCT}% (high)"
else
    fail "Memory usage: ${MEM_PCT}% (CRITICAL)"
fi

# Load
LOAD=$(uptime | awk -F'load average:' '{print $2}' | awk -F',' '{print $1}' | xargs)
CPUS=$(nproc)
LOAD_INT=$(echo "$LOAD" | cut -d'.' -f1)
if [ "$LOAD_INT" -lt "$CPUS" ]; then
    pass "Load average: $LOAD (${CPUS} CPUs)"
else
    warn "Load average: $LOAD (${CPUS} CPUs) — high"
fi

echo ""

# ─── Tailscale ───────────────────────────────────────────────────────────────

echo "Network:"
if command -v tailscale &>/dev/null; then
    TS_STATUS=$(tailscale status --json 2>/dev/null | jq -r '.Self.Online // false' 2>/dev/null || echo "false")
    TS_IP=$(tailscale ip -4 2>/dev/null || echo "none")
    if [ "$TS_STATUS" = "true" ]; then
        pass "Tailscale: connected ($TS_IP)"
    else
        fail "Tailscale: not connected"
    fi
else
    fail "Tailscale: not installed"
fi

# UFW
if command -v ufw &>/dev/null; then
    UFW_STATUS=$(sudo ufw status 2>/dev/null | head -1)
    if echo "$UFW_STATUS" | grep -q "active"; then
        pass "Firewall: active"
    else
        warn "Firewall: inactive"
    fi
else
    warn "UFW: not installed"
fi

echo ""

# ─── Docker ──────────────────────────────────────────────────────────────────

echo "Docker:"
if command -v docker &>/dev/null && docker info &>/dev/null; then
    pass "Docker daemon: running"
else
    fail "Docker daemon: not running"
    echo ""
    echo "Cannot check Docker services. Fix Docker first."
    echo ""
    exit 1
fi

# Check each container
check_container() {
    local name=$1
    local required=$2
    local port=$3

    STATUS=$(docker inspect --format='{{.State.Status}}' "$name" 2>/dev/null || echo "missing")
    HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "$name" 2>/dev/null || echo "none")

    if [ "$STATUS" = "running" ]; then
        if [ "$HEALTH" = "healthy" ]; then
            pass "$name: running (healthy)"
        elif [ "$HEALTH" = "none" ]; then
            pass "$name: running"
        else
            warn "$name: running ($HEALTH)"
        fi

        # Port check if specified
        if [ -n "$port" ]; then
            if docker exec "$name" sh -c "wget -qO- http://localhost:$port >/dev/null 2>&1 || curl -sf http://localhost:$port >/dev/null 2>&1" 2>/dev/null; then
                pass "  Port $port: responding"
            fi
        fi
    elif [ "$STATUS" = "missing" ]; then
        if [ "$required" = "required" ]; then
            fail "$name: not found"
        else
            skip "$name: not deployed (optional)"
        fi
    else
        if [ "$required" = "required" ]; then
            fail "$name: $STATUS"
        else
            warn "$name: $STATUS"
        fi
    fi
}

check_container "onemind-nats"          "required" "8222"
check_container "onemind-omnidb"        "optional" ""
check_container "onemind-nats-ingester" "optional" ""
check_container "onemind-legacy-bridge" "required" ""
check_container "onemind-n8n"           "optional" "5678"
check_container "onemind-livekit"       "optional" "7880"

echo ""

# ─── NATS Specific ───────────────────────────────────────────────────────────

echo "NATS:"
NATS_HEALTH=$(curl -sf http://localhost:8222/healthz 2>/dev/null)
if [ $? -eq 0 ]; then
    pass "NATS API: responding"

    # Connection count
    NATS_CONNZ=$(curl -sf http://localhost:8222/connz 2>/dev/null)
    if [ $? -eq 0 ]; then
        CONN_COUNT=$(echo "$NATS_CONNZ" | jq '.num_connections // 0' 2>/dev/null || echo "?")
        pass "NATS connections: $CONN_COUNT"
    fi

    # JetStream status
    NATS_JSZ=$(curl -sf http://localhost:8222/jsz 2>/dev/null)
    if [ $? -eq 0 ]; then
        JS_STREAMS=$(echo "$NATS_JSZ" | jq '.streams // 0' 2>/dev/null || echo "?")
        JS_MSGS=$(echo "$NATS_JSZ" | jq '.messages // 0' 2>/dev/null || echo "?")
        pass "JetStream: $JS_STREAMS streams, $JS_MSGS messages"
    fi
else
    fail "NATS API: not responding on :8222"
fi

echo ""

# ─── OmniDB ─────────────────────────────────────────────────────────────────

echo "OmniDB:"
if docker inspect onemind-omnidb &>/dev/null; then
    # Test PG connection
    DB_CHECK=$(docker exec onemind-omnidb pg_isready -U onemind 2>/dev/null)
    if [ $? -eq 0 ]; then
        pass "PostgreSQL: accepting connections"

        # Check extensions
        TSDB=$(docker exec onemind-omnidb psql -U onemind -d onemind -tAc "SELECT extversion FROM pg_extension WHERE extname='timescaledb'" 2>/dev/null || echo "")
        PGVEC=$(docker exec onemind-omnidb psql -U onemind -d onemind -tAc "SELECT extversion FROM pg_extension WHERE extname='vector'" 2>/dev/null || echo "")

        if [ -n "$TSDB" ]; then
            pass "TimescaleDB: v$TSDB"
        else
            warn "TimescaleDB: not loaded"
        fi

        if [ -n "$PGVEC" ]; then
            pass "pgvector: v$PGVEC"
        else
            warn "pgvector: not loaded"
        fi

        # Event count
        EVT_COUNT=$(docker exec onemind-omnidb psql -U onemind -d onemind -tAc "SELECT count(*) FROM events" 2>/dev/null || echo "error")
        if [ "$EVT_COUNT" != "error" ]; then
            pass "Events stored: $EVT_COUNT"
        fi
    else
        fail "PostgreSQL: not accepting connections"
    fi
else
    skip "OmniDB: not deployed (use --profile data)"
fi

echo ""

# ─── OpenClaw Gateway ───────────────────────────────────────────────────────

echo "OpenClaw Gateway:"
if command -v openclaw &>/dev/null; then
    pass "OpenClaw: installed ($(openclaw --version 2>/dev/null || echo 'yes'))"
else
    fail "OpenClaw: not installed"
fi

# Check systemd service
if systemctl is-active --quiet openclaw 2>/dev/null; then
    pass "Service: active (systemd)"
    UPTIME=$(systemctl show openclaw --property=ActiveEnterTimestamp --value 2>/dev/null || echo "unknown")
    pass "  Since: $UPTIME"
elif pgrep -f "openclaw gateway" &>/dev/null; then
    pass "Gateway: running (manual/foreground)"
else
    fail "Gateway: not running"
fi

# Health endpoint
GW_HEALTH=$(curl -sf http://localhost:18789/health 2>/dev/null)
if [ $? -eq 0 ]; then
    pass "Gateway API: responding on :18789"
else
    # Try Tailscale IP
    TS_IP=$(tailscale ip -4 2>/dev/null || echo "")
    if [ -n "$TS_IP" ]; then
        GW_HEALTH=$(curl -sf "http://$TS_IP:18789/health" 2>/dev/null)
        if [ $? -eq 0 ]; then
            pass "Gateway API: responding on $TS_IP:18789"
        else
            fail "Gateway API: not responding"
        fi
    else
        fail "Gateway API: not responding on :18789"
    fi
fi

echo ""

# ─── Workspace ───────────────────────────────────────────────────────────────

echo "Workspace:"
OC_WORKSPACE="${HOME:-/home/onemind}/.openclaw/workspace"
if [ -d "$OC_WORKSPACE" ]; then
    pass "Workspace: $OC_WORKSPACE"

    for file in SOUL.md AGENTS.md HEARTBEAT.md USER.md; do
        if [ -f "$OC_WORKSPACE/$file" ]; then
            pass "  $file: present"
        else
            warn "  $file: MISSING"
        fi
    done

    SKILL_COUNT=$(ls -d "$OC_WORKSPACE/skills/"*/ 2>/dev/null | wc -l)
    pass "  Skills: $SKILL_COUNT installed"

    # Backup check
    if [ -d "$OC_WORKSPACE/.git" ]; then
        LAST_COMMIT=$(cd "$OC_WORKSPACE" && git log --oneline -1 2>/dev/null || echo "no commits")
        pass "  Backup: git active ($LAST_COMMIT)"
    else
        warn "  Backup: no git initialized"
    fi
else
    fail "Workspace: not found at $OC_WORKSPACE"
fi

echo ""

# ─── Summary ────────────────────────────────────────────────────────────────

echo "============================================"
echo "  Results: ${GREEN}${PASS} passed${NC}, ${RED}${FAIL} failed${NC}, ${YELLOW}${WARN} warnings${NC}"
echo "============================================"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo -e "${RED}Some checks failed. Review above and fix issues.${NC}"
    exit 1
elif [ "$WARN" -gt 0 ]; then
    echo -e "${YELLOW}All critical checks passed with some warnings.${NC}"
    exit 0
else
    echo -e "${GREEN}All checks passed. OneMind OS is healthy.${NC}"
    exit 0
fi
