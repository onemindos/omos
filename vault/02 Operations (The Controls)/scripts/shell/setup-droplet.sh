#!/usr/bin/env bash
# =============================================================================
# OneMind OS — DigitalOcean Droplet Setup (Phase 1: Base Infrastructure)
# =============================================================================
#
# Sets up a fresh Ubuntu 24.04 Droplet with:
#   - System hardening (UFW, Fail2ban, unattended-upgrades)
#   - Docker + Docker Compose
#   - Tailscale mesh VPN
#   - Node.js 22 LTS (for OpenClaw)
#   - onemind user (non-root, passwordless sudo)
#   - Repo cloned to /opt/onemind
#   - .env file with auto-generated secrets
#
# DOES NOT install OpenClaw — run install-openclaw.sh after this.
#
# Run on a fresh Ubuntu 24.04 Droplet:
#   curl -sSL https://raw.githubusercontent.com/Zeus-Delacruz/onemindos-v0.4/main/.do/setup-droplet.sh | bash
#
# Or SSH in and run:
#   sudo bash .do/setup-droplet.sh
#
# ─────────────────────────────────────────────────────────────────────────────
#
# DO Architecture (Single Droplet):
#
#   ┌─────────────────────────────────────────────────────────────────┐
#   │  DigitalOcean Droplet (s-2vcpu-4gb-amd / $28/mo)              │
#   │                                                                 │
#   │  ┌─────────────────────────┐  ┌──────────────────────────────┐ │
#   │  │  BARE METAL (host)      │  │  DOCKER COMPOSE              │ │
#   │  │                         │  │                              │ │
#   │  │  • OpenClaw Gateway     │  │  • NATS JetStream    :4222  │ │
#   │  │    (Node.js)   :18789   │  │  • OmniDB (PG16)    :5432  │ │
#   │  │  • Tailscale VPN        │  │  • NATS Ingester            │ │
#   │  │  • UFW Firewall         │  │  • Legacy Bridge            │ │
#   │  │  • Fail2ban             │  │  • n8n (optional)   :5678  │ │
#   │  │  • systemd services     │  │  • LiveKit (opt)    :7880  │ │
#   │  └─────────────────────────┘  └──────────────────────────────┘ │
#   │                                                                 │
#   │  Tailscale IP: 100.x.y.z (only way to reach services)         │
#   └─────────────────────────────────────────────────────────────────┘
#
# Minimum Droplet:   s-2vcpu-4gb      ($24/mo)
# Recommended:       s-2vcpu-4gb-amd  ($28/mo) — better single-thread
# Full profile:      s-4vcpu-8gb-amd  ($56/mo) — n8n + LiveKit
#
# =============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log()  { echo -e "${GREEN}→${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
err()  { echo -e "${RED}✗${NC} $1"; }
info() { echo -e "${BLUE}ℹ${NC} $1"; }

REPO_DIR="/opt/onemind"

echo ""
echo "============================================"
echo "  OneMind OS — Droplet Setup (Phase 1)"
echo "  Base Infrastructure"
echo "============================================"
echo ""

# ─── Preflight ───────────────────────────────────────────────────────────────

if [ "$(id -u)" -ne 0 ]; then
    err "Must run as root. Use: sudo bash .do/setup-droplet.sh"
    exit 1
fi

if ! grep -qi "ubuntu" /etc/os-release 2>/dev/null; then
    warn "This script targets Ubuntu 24.04. Other distros may work but are untested."
fi

# ─── 1. System Updates & Essentials ─────────────────────────────────────────

log "Updating system packages..."
export DEBIAN_FRONTEND=noninteractive
apt update && apt upgrade -y
apt install -y \
    curl git jq wget unzip htop \
    ufw fail2ban \
    unattended-upgrades apt-listchanges \
    ca-certificates gnupg lsb-release

# Auto security updates
log "Enabling unattended security upgrades..."
cat > /etc/apt/apt.conf.d/20auto-upgrades << 'EOF'
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
APT::Periodic::AutocleanInterval "7";
EOF

# ─── 2. Create onemind User ────────────────────────────────────────────────

log "Creating onemind user..."
if ! id "onemind" &>/dev/null; then
    useradd -m -s /bin/bash -G sudo onemind
    echo "onemind ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/onemind
    chmod 440 /etc/sudoers.d/onemind
    log "User 'onemind' created with passwordless sudo."
else
    log "User 'onemind' already exists."
fi

# ─── 3. Docker + Docker Compose ────────────────────────────────────────────

log "Installing Docker..."
if ! command -v docker &>/dev/null; then
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
    log "Docker installed."
else
    log "Docker already installed: $(docker --version)"
fi

usermod -aG docker onemind

if ! docker compose version &>/dev/null; then
    log "Installing Docker Compose plugin..."
    apt install -y docker-compose-plugin
fi
log "Docker Compose: $(docker compose version 2>/dev/null)"

# ─── 4. Node.js 22 LTS (OpenClaw requires Node.js) ─────────────────────────

log "Installing Node.js 22 LTS..."
if ! command -v node &>/dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash -
    apt install -y nodejs
    log "Node.js installed: $(node --version)"
else
    NODE_VER=$(node --version | cut -d'.' -f1 | tr -d 'v')
    if [ "$NODE_VER" -lt 20 ]; then
        warn "Node.js $(node --version) is too old for OpenClaw. Installing v22..."
        curl -fsSL https://deb.nodesource.com/setup_22.x | bash -
        apt install -y nodejs
    else
        log "Node.js already installed: $(node --version)"
    fi
fi

# ─── 5. Tailscale VPN ──────────────────────────────────────────────────────

log "Installing Tailscale..."
if ! command -v tailscale &>/dev/null; then
    curl -fsSL https://tailscale.com/install.sh | sh
    log "Tailscale installed."
else
    log "Tailscale already installed."
fi

# ─── 6. Firewall (UFW) ─────────────────────────────────────────────────────

log "Configuring base firewall..."
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp comment "SSH"
ufw --force enable
log "Firewall enabled (SSH only). Services locked to Tailscale after setup-firewall.sh."

# ─── 7. Fail2ban ───────────────────────────────────────────────────────────

log "Configuring Fail2ban..."
cat > /etc/fail2ban/jail.local << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
maxretry = 3
bantime = 7200
EOF

systemctl enable fail2ban
systemctl restart fail2ban
log "Fail2ban active. SSH: 3 fails → 2hr ban."

# ─── 8. Clone Repository ───────────────────────────────────────────────────

log "Setting up repository at $REPO_DIR..."
if [ ! -d "$REPO_DIR/.git" ]; then
    if ssh -o StrictHostKeyChecking=no -o BatchMode=yes -T git@github.com 2>&1 | grep -q "success"; then
        git clone git@github.com:Zeus-Delacruz/onemindos-v0.4.git "$REPO_DIR"
        log "Cloned via SSH (deploy key)."
    else
        warn "No GitHub deploy key. You'll need to clone manually:"
        warn "  git clone git@github.com:Zeus-Delacruz/onemindos-v0.4.git $REPO_DIR"
        mkdir -p "$REPO_DIR"
    fi
else
    cd "$REPO_DIR" && git pull --ff-only 2>/dev/null || warn "Git pull failed (not fatal)."
    log "Repository already exists."
fi

chown -R onemind:onemind "$REPO_DIR"

# ─── 9. Environment File ───────────────────────────────────────────────────

ENV_FILE="$REPO_DIR/.env"
log "Setting up environment file..."

if [ ! -f "$ENV_FILE" ]; then
    # Auto-generate secrets
    OC_TOKEN=$(openssl rand -hex 32)
    HOOKS_TOKEN=$(openssl rand -hex 32)
    N8N_KEY=$(openssl rand -hex 32)
    DB_PASS="onemind-$(openssl rand -hex 12)"
    EDGE_HOME=$(openssl rand -hex 16)
    EDGE_OFFICE=$(openssl rand -hex 16)
    EDGE_FARM=$(openssl rand -hex 16)

    cat > "$ENV_FILE" << ENVEOF
# =============================================================================
# OneMind OS — Environment Variables
# Generated: $(date -u +"%Y-%m-%dT%H:%M:%SZ") on $(hostname)
# =============================================================================
# Secrets were auto-generated. Fill in API keys marked CHANGE_ME.

# ── LLM Providers ────────────────────────────────────────────────────────────
ANTHROPIC_API_KEY=CHANGE_ME
OPENAI_API_KEY=CHANGE_ME
# GOOGLE_API_KEY=CHANGE_ME

# ── OpenClaw Gateway ─────────────────────────────────────────────────────────
OPENCLAW_GATEWAY_TOKEN=${OC_TOKEN}
OPENCLAW_HOOKS_TOKEN=${HOOKS_TOKEN}
OPENCLAW_GATEWAY_URL=http://host.docker.internal:18789
OPENCLAW_WORKSPACE=/home/onemind/.openclaw/workspace
LEGACY_MODEL=anthropic/claude-opus-4-6

# ── Contacts ─────────────────────────────────────────────────────────────────
ZEUS_WHATSAPP_NUMBER=CHANGE_ME
ZEUS_PHONE_NUMBER=CHANGE_ME
ZEUS_TELEGRAM_ID=CHANGE_ME

# ── Voice (ElevenLabs) ───────────────────────────────────────────────────────
ELEVENLABS_API_KEY=CHANGE_ME
ELEVENLABS_VOICE_ID=CHANGE_ME

# ── Voice Calls (Twilio) ─────────────────────────────────────────────────────
TWILIO_ACCOUNT_SID=CHANGE_ME
TWILIO_AUTH_TOKEN=CHANGE_ME
TWILIO_FROM_NUMBER=CHANGE_ME

# ── Messaging ────────────────────────────────────────────────────────────────
TELEGRAM_BOT_TOKEN=CHANGE_ME
DISCORD_BOT_TOKEN=CHANGE_ME
DISCORD_COMMUNITY_CHANNEL_ID=CHANGE_ME

# ── Web Tools ────────────────────────────────────────────────────────────────
BRAVE_API_KEY=CHANGE_ME

# ── NATS JetStream ───────────────────────────────────────────────────────────
NATS_URL=nats://nats:4222
EDGE_HOME_PASSWORD=${EDGE_HOME}
EDGE_OFFICE_PASSWORD=${EDGE_OFFICE}
EDGE_FARM_PASSWORD=${EDGE_FARM}

# ── OmniDB (PostgreSQL 16 + TimescaleDB + pgvector) ──────────────────────────
OMNIDB_USER=onemind
OMNIDB_PASSWORD=${DB_PASS}
OMNIDB_DB=onemind
OMNIDB_URL=postgresql://onemind:${DB_PASS}@omnidb:5432/onemind

# ── Legacy Bridge ────────────────────────────────────────────────────────────
BRIDGE_AGENT_ID=legacy

# ── n8n (profile: full) ──────────────────────────────────────────────────────
N8N_ENCRYPTION_KEY=${N8N_KEY}
N8N_WEBHOOK_URL=http://host.docker.internal:5678

# ── LiveKit (profile: full) ──────────────────────────────────────────────────
LIVEKIT_API_KEY=devkey
LIVEKIT_API_SECRET=devsecret
LIVEKIT_URL=ws://localhost:7880

# ── Logging ──────────────────────────────────────────────────────────────────
LOG_LEVEL=INFO
ENVEOF

    chown onemind:onemind "$ENV_FILE"
    chmod 600 "$ENV_FILE"
    warn "Environment file created: $ENV_FILE"
    warn "Fill in API keys: nano $ENV_FILE"
else
    log "Environment file already exists."
fi

# ─── 10. Create Directories ────────────────────────────────────────────────

log "Creating directory structure..."
mkdir -p /home/onemind/.openclaw/workspace/{memory,canvas,lobster}
chown -R onemind:onemind /home/onemind/.openclaw

# ─── Summary ────────────────────────────────────────────────────────────────

echo ""
echo "============================================"
echo "  Phase 1 Complete — Base Infrastructure"
echo "============================================"
echo ""
echo "  Installed:"
echo "    Docker:     $(docker --version 2>/dev/null | cut -d' ' -f3 || echo 'error')"
echo "    Compose:    $(docker compose version 2>/dev/null | cut -d' ' -f4 || echo 'error')"
echo "    Node.js:    $(node --version 2>/dev/null || echo 'error')"
echo "    Tailscale:  installed"
echo "    UFW:        active (SSH only)"
echo "    Fail2ban:   active"
echo ""
echo "  ┌──────────────────────────────────────────────────┐"
echo "  │  NEXT STEPS (run in this order):                 │"
echo "  │                                                  │"
echo "  │  1. sudo tailscale up                            │"
echo "  │     (connect to your Tailnet)                    │"
echo "  │                                                  │"
echo "  │  2. sudo bash $REPO_DIR/.do/setup-firewall.sh   │"
echo "  │     (open service ports on Tailscale only)       │"
echo "  │                                                  │"
echo "  │  3. nano $REPO_DIR/.env                          │"
echo "  │     (fill in CHANGE_ME values)                   │"
echo "  │                                                  │"
echo "  │  4. sudo -u onemind bash                         │"
echo "  │        $REPO_DIR/.do/install-openclaw.sh         │"
echo "  │     (install OpenClaw + copy workspace + systemd)│"
echo "  │                                                  │"
echo "  │  5. cd $REPO_DIR && docker compose up -d         │"
echo "  │     (start NATS + Legacy Bridge)                 │"
echo "  │                                                  │"
echo "  │  6. bash $REPO_DIR/.do/health-check.sh           │"
echo "  │     (verify all services running)                │"
echo "  └──────────────────────────────────────────────────┘"
echo ""
