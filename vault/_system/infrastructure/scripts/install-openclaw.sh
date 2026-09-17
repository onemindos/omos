#!/usr/bin/env bash
# =============================================================================
# OneMind OS — Install OpenClaw Gateway (Phase 2)
# =============================================================================
# Installs OpenClaw, copies workspace files, creates systemd service.
#
# Run as the onemind user:
#   sudo -u onemind bash .do/install-openclaw.sh
#
# Prerequisites:
#   - Node.js 22+ installed (from setup-droplet.sh)
#   - Repository cloned to /opt/onemind
#   - .env file configured
#
# =============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}→${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
err()  { echo -e "${RED}✗${NC} $1"; }

REPO_DIR="/opt/onemind"
OC_WORKSPACE="$HOME/.openclaw/workspace"
OC_CONFIG="$HOME/.openclaw"

echo ""
echo "============================================"
echo "  OneMind OS — Install OpenClaw (Phase 2)"
echo "============================================"
echo ""

# ─── Preflight ───────────────────────────────────────────────────────────────

if [ "$(id -u)" -eq 0 ]; then
    err "Do NOT run as root. Run as onemind user:"
    err "  sudo -u onemind bash $REPO_DIR/.do/install-openclaw.sh"
    exit 1
fi

if ! command -v node &>/dev/null; then
    err "Node.js not found. Run setup-droplet.sh first."
    exit 1
fi

if [ ! -d "$REPO_DIR" ]; then
    err "Repository not found at $REPO_DIR. Run setup-droplet.sh first."
    exit 1
fi

# ─── 1. Install OpenClaw ───────────────────────────────────────────────────

log "Installing OpenClaw..."
if ! command -v openclaw &>/dev/null; then
    npm install -g openclaw
    log "OpenClaw installed: $(openclaw --version 2>/dev/null || echo 'installed')"
else
    log "OpenClaw already installed: $(openclaw --version 2>/dev/null || echo 'yes')"
    log "Updating to latest..."
    npm update -g openclaw
fi

# ─── 2. Initialize Workspace ───────────────────────────────────────────────

log "Setting up OpenClaw workspace..."
mkdir -p "$OC_WORKSPACE"/{memory,canvas,lobster,skills}

# ─── 3. Copy Workspace Files ───────────────────────────────────────────────

log "Copying workspace files from repo..."
WORKSPACE_SRC="$REPO_DIR/openclaw/workspace"

# Core identity files
for file in SOUL.md AGENTS.md BOOT.md HEARTBEAT.md IDENTITY.md TOOLS.md USER.md; do
    if [ -f "$WORKSPACE_SRC/$file" ]; then
        cp "$WORKSPACE_SRC/$file" "$OC_WORKSPACE/$file"
        log "  Copied $file"
    else
        warn "  Missing $file in repo workspace"
    fi
done

# Bridge script
if [ -f "$WORKSPACE_SRC/bridge.py" ]; then
    cp "$WORKSPACE_SRC/bridge.py" "$OC_WORKSPACE/bridge.py"
    log "  Copied bridge.py"
fi

# Project manager
if [ -f "$WORKSPACE_SRC/project_manager.py" ]; then
    cp "$WORKSPACE_SRC/project_manager.py" "$OC_WORKSPACE/project_manager.py"
    log "  Copied project_manager.py"
fi

# ─── 4. Copy Skills ────────────────────────────────────────────────────────

log "Copying skills..."
if [ -d "$WORKSPACE_SRC/skills" ]; then
    cp -r "$WORKSPACE_SRC/skills/"* "$OC_WORKSPACE/skills/" 2>/dev/null || true
    SKILL_COUNT=$(ls -d "$OC_WORKSPACE/skills/"*/ 2>/dev/null | wc -l)
    log "  Copied $SKILL_COUNT skills"
fi

# Install skill dependencies
log "Installing skill dependencies..."
for skill_dir in "$OC_WORKSPACE/skills/"*/; do
    if [ -f "$skill_dir/package.json" ]; then
        (cd "$skill_dir" && npm install --production --silent 2>/dev/null) || warn "  Failed: $(basename $skill_dir)"
    fi
done

# ─── 5. Configure openclaw.json ────────────────────────────────────────────

log "Configuring openclaw.json..."
if [ -f "$WORKSPACE_SRC/openclaw.json" ]; then
    cp "$WORKSPACE_SRC/openclaw.json" "$OC_CONFIG/openclaw.json"
    log "  Copied openclaw.json to $OC_CONFIG/"

    # Reminder about env vars
    warn "openclaw.json uses \${ENV_VAR} syntax."
    warn "Make sure .env is sourced or vars are exported before starting gateway."
else
    warn "No openclaw.json found in repo. Running 'openclaw setup'..."
    openclaw setup 2>/dev/null || true
fi

# ─── 6. Create Environment Loader ──────────────────────────────────────────

log "Creating environment loader..."
cat > "$HOME/.openclaw/load-env.sh" << 'LOADER'
#!/usr/bin/env bash
# Load .env file for OpenClaw gateway
ENV_FILE="/opt/onemind/.env"
if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
fi
LOADER
chmod +x "$HOME/.openclaw/load-env.sh"

# ─── 7. Create systemd Service ─────────────────────────────────────────────

log "Creating systemd service..."

# Need sudo for systemd
SYSTEMD_FILE="/etc/systemd/system/openclaw.service"
sudo tee "$SYSTEMD_FILE" > /dev/null << SERVICEEOF
[Unit]
Description=OneMind OS — OpenClaw Gateway (Legacy's Runtime)
Documentation=https://github.com/Zeus-Delacruz/onemindos-v0.4
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
User=onemind
Group=onemind
WorkingDirectory=/home/onemind/.openclaw

# Load environment from .env
EnvironmentFile=-/opt/onemind/.env

# OpenClaw gateway command
ExecStart=$(which openclaw) gateway
ExecReload=/bin/kill -HUP \$MAINPID

# Restart on failure (self-healing)
Restart=always
RestartSec=10
StartLimitBurst=5
StartLimitIntervalSec=300

# Resource limits
LimitNOFILE=65536
MemoryMax=2G

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=openclaw

[Install]
WantedBy=multi-user.target
SERVICEEOF

sudo systemctl daemon-reload
sudo systemctl enable openclaw
log "systemd service created and enabled."

# ─── 8. Create Workspace Backup Cron ───────────────────────────────────────

log "Setting up workspace backup..."
warn "Workspace backup requires a private GitHub repo."
warn "Create: github.com/Zeus-Delacruz/openclaw-workspace-backup"
echo ""

cat > "$HOME/.openclaw/backup.sh" << 'BACKUP'
#!/usr/bin/env bash
# Auto-push OpenClaw workspace to GitHub every 15 minutes
cd ~/.openclaw/workspace || exit 1

# Initialize git if needed
if [ ! -d .git ]; then
    git init
    git remote add origin git@github.com:Zeus-Delacruz/openclaw-workspace-backup.git 2>/dev/null || true
fi

git add -A
git commit -m "auto-backup $(date +%Y-%m-%d_%H:%M)" --allow-empty 2>/dev/null
git push origin main 2>/dev/null || git push -u origin main 2>/dev/null || true
BACKUP
chmod +x "$HOME/.openclaw/backup.sh"

# Add cron (won't duplicate if already exists)
(crontab -l 2>/dev/null | grep -v "openclaw/backup" ; echo "*/15 * * * * /home/onemind/.openclaw/backup.sh >> /home/onemind/.openclaw/backup.log 2>&1") | crontab -

log "Backup cron installed (every 15 minutes)."

# ─── Summary ────────────────────────────────────────────────────────────────

echo ""
echo "============================================"
echo "  Phase 2 Complete — OpenClaw Installed"
echo "============================================"
echo ""
echo "  OpenClaw:    $(openclaw --version 2>/dev/null || echo 'installed')"
echo "  Config:      $OC_CONFIG/openclaw.json"
echo "  Workspace:   $OC_WORKSPACE/"
echo "  Service:     openclaw.service (systemd)"
echo "  Backup:      every 15 min → GitHub"
echo ""
echo "  Workspace contents:"
ls -1 "$OC_WORKSPACE/" 2>/dev/null | sed 's/^/    /'
echo ""
echo "  Skills installed:"
ls -d "$OC_WORKSPACE/skills/"*/ 2>/dev/null | xargs -I{} basename {} | sed 's/^/    /' || echo "    (none)"
echo ""
echo "  ┌──────────────────────────────────────────────────┐"
echo "  │  START OPENCLAW:                                 │"
echo "  │                                                  │"
echo "  │  # Manual (foreground, see output):              │"
echo "  │  source ~/.openclaw/load-env.sh                  │"
echo "  │  openclaw gateway                                │"
echo "  │                                                  │"
echo "  │  # Service (background, auto-restart):           │"
echo "  │  sudo systemctl start openclaw                   │"
echo "  │  sudo journalctl -u openclaw -f                  │"
echo "  │                                                  │"
echo "  │  # Then start Docker services:                   │"
echo "  │  cd /opt/onemind && docker compose up -d         │"
echo "  └──────────────────────────────────────────────────┘"
echo ""
