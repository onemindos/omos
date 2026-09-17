#!/usr/bin/env bash
# deploy-frontends.sh — Deploy openclaw-studio + swarmclaw to VPS
#
# IMPORTANT: Always build on the VPS, never rsync a local build.
# Native Node.js addons (better-sqlite3) are platform-specific — a macOS build
# will crash with "invalid ELF header" on Linux. Building on the VPS ensures the
# correct Linux binaries and prevents .env.local from leaking local secrets.
#
# Run from repo root:
#   bash openclaw/deploy/scripts/deploy-frontends.sh
#
# To deploy only SwarmClaw:
#   bash openclaw/deploy/scripts/deploy-frontends.sh --swarmclaw-only
#
# To deploy only Studio:
#   bash openclaw/deploy/scripts/deploy-frontends.sh --studio-only

set -euo pipefail

VPS="root@104.248.4.212"
REPO_DIR="/home/onemind/onemind-os"
STUDIO_DIR="$REPO_DIR/command-center/openclaw-studio"
SWARM_DIR="$REPO_DIR/command-center/swarmclaw"
DEPLOY_DIR="$REPO_DIR/openclaw/deploy"

DEPLOY_STUDIO=true
DEPLOY_SWARM=true

for arg in "$@"; do
  case "$arg" in
    --swarmclaw-only) DEPLOY_STUDIO=false ;;
    --studio-only)    DEPLOY_SWARM=false ;;
  esac
done

echo "=== Pulling latest code on VPS ==="
ssh "$VPS" "su - onemind -c 'cd $REPO_DIR && git pull --ff-only'"

if $DEPLOY_SWARM; then
  echo ""
  echo "=== Building SwarmClaw on VPS ==="
  # Build on VPS so native addons (better-sqlite3) compile for Linux.
  # Secrets (ACCESS_KEY, CREDENTIAL_SECRET) come from systemd override.conf —
  # never from .env.local. Any auto-generated .env.local is removed after build.
  ssh "$VPS" "
    cd $SWARM_DIR
    npm install --prefer-offline 2>&1 | tail -3
    npm run build 2>&1 | tail -8
    rm -f .next/standalone/.env.local
    chown -R onemind:onemind .next/
    echo 'SwarmClaw build done'
  "

  echo ""
  echo "=== Restarting SwarmClaw ==="
  ssh "$VPS" "
    cp $DEPLOY_DIR/systemd/swarmclaw.service /etc/systemd/system/swarmclaw.service
    systemctl daemon-reload
    systemctl enable swarmclaw
    systemctl restart swarmclaw
    sleep 4
    journalctl -u swarmclaw -n 6 --no-pager 2>/dev/null | tail -6
  "
fi

if $DEPLOY_STUDIO; then
  echo ""
  echo "=== Building openclaw-studio on VPS ==="
  ssh "$VPS" "
    cd $STUDIO_DIR
    npm install --prefer-offline 2>&1 | tail -3
    npm run build 2>&1 | tail -5
    chown -R onemind:onemind .next/ 2>/dev/null || true
    echo 'Studio build done'
  "

  echo ""
  echo "=== Restarting openclaw-studio ==="
  ssh "$VPS" "
    cp $DEPLOY_DIR/systemd/openclaw-studio.service /etc/systemd/system/openclaw-studio.service 2>/dev/null || true
    systemctl daemon-reload
    systemctl enable openclaw-studio 2>/dev/null || true
    systemctl restart openclaw-studio 2>/dev/null || true
    echo 'Studio restarted'
  "
fi

echo ""
echo "=== Tailscale serve (idempotent) ==="
ssh "$VPS" "
  tailscale serve --bg --yes --https=3443 http://127.0.0.1:3000 2>/dev/null || true
  tailscale serve --bg --yes --https=8443 http://127.0.0.1:3456 2>/dev/null || true
"

echo ""
echo "=== Service status ==="
ssh "$VPS" "systemctl is-active openclaw swarmclaw 2>/dev/null && systemctl status swarmclaw --no-pager -l | head -10"

echo ""
echo "=== Done ==="
echo "  Gateway:   https://onemind-main.tail717ba0.ts.net        (port 443)"
echo "  Studio:    https://onemind-main.tail717ba0.ts.net:3443"
echo "  SwarmClaw: https://onemind-main.tail717ba0.ts.net:8443"
