#!/bin/bash
# setup.sh — Full OpenClaw VPS deployment from scratch
# Run as root on a fresh Ubuntu 24.04 droplet.
#
# Prerequisites:
#   - Fresh Ubuntu 24.04 (DigitalOcean, AWS, etc.)
#   - Root SSH access
#   - GitHub repo: github.com/One-Mind-OS/onemindos-v0.4
#
# Usage:
#   scp -r openclaw/deploy/ root@YOUR_IP:/tmp/deploy/
#   ssh root@YOUR_IP 'bash /tmp/deploy/scripts/setup.sh'
#
# After running, you still need to:
#   1. Copy openclaw.json to ~/.openclaw/openclaw.json (fill in real values)
#   2. Create /opt/onemind/.env with API keys
#   3. Run: sudo -u onemind gh auth login
#   4. Install Tailscale: curl -fsSL https://tailscale.com/install.sh | sh && tailscale up
#   5. Start services: systemctl start openclaw codex-watch openclaw-reload codex-pull.timer

set -euo pipefail

echo "============================================"
echo "  OneMind OS — OpenClaw VPS Setup"
echo "============================================"
echo ""

# --- 1. System user ---
echo "[1/8] Creating onemind user..."
if ! id onemind &>/dev/null; then
  useradd -m -s /bin/bash onemind
  echo "onemind ALL=(ALL) NOPASSWD: /usr/bin/systemctl restart openclaw" >> /etc/sudoers.d/onemind
  chmod 440 /etc/sudoers.d/onemind
fi
echo "  Done."

# --- 2. System packages ---
echo "[2/8] Installing system packages..."
apt-get update -qq
apt-get install -y -qq git curl docker.io docker-compose-v2 inotify-tools jq
systemctl enable docker --now
usermod -aG docker onemind
echo "  Done."

# --- 3. Node.js + OpenClaw + 1Password CLI ---
echo "[3/8] Installing Node.js 22 + OpenClaw + 1Password CLI..."
if ! command -v node &>/dev/null; then
  curl -fsSL https://deb.nodesource.com/setup_22.x | bash -
  apt-get install -y -qq nodejs
fi
npm install -g openclaw
echo "  OpenClaw version: $(openclaw --version 2>/dev/null || echo 'installed')"
# Install 1Password CLI (op) — used by openclaw.service to inject secrets at runtime
if ! command -v op &>/dev/null; then
  curl -sS https://downloads.1password.com/linux/keys/1password.asc \
    | gpg --dearmor --output /usr/share/keyrings/1password-archive-keyring.gpg
  echo "deb [arch=amd64 signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/amd64 stable main" \
    | tee /etc/apt/sources.list.d/1password.list
  apt-get update -qq && apt-get install -y -qq 1password-cli
fi
echo "  op version: $(op --version 2>/dev/null || echo 'installed')"
echo "  Done."

# --- 4. Clone repo ---
echo "[4/8] Cloning monorepo..."
REPO=/home/onemind/onemind-os
if [ ! -d "$REPO" ]; then
  sudo -u onemind git clone https://github.com/One-Mind-OS/onemindos-v0.4.git "$REPO"
else
  echo "  Repo already exists at $REPO"
fi
echo "  Done."

# --- 5. Initialize OpenClaw ---
echo "[5/8] Initializing OpenClaw directory..."
sudo -u onemind mkdir -p /home/onemind/.openclaw
# Run setup-symlinks
sudo -u onemind bash "$REPO/openclaw/deploy/scripts/setup-symlinks.sh"
echo "  Done."

# --- 6. /opt/onemind structure ---
echo "[6/8] Setting up /opt/onemind..."
mkdir -p /opt/onemind
# Bootstrap file: only OP_SERVICE_ACCOUNT_TOKEN + static non-secret vars.
# All other secrets fetched from 1Password at runtime by op run.
if [ ! -f /opt/onemind/.env.sa ]; then
  cp "$REPO/openclaw/deploy/env.sa.example" /opt/onemind/.env.sa
  chmod 600 /opt/onemind/.env.sa
  echo "  Created /opt/onemind/.env.sa — ADD OP_SERVICE_ACCOUNT_TOKEN"
else
  echo "  /opt/onemind/.env.sa already exists"
fi
# 1Password secrets template (references, not real values — safe to copy as-is)
cp "$REPO/openclaw/deploy/env.tpl" /opt/onemind/.env.tpl
echo "  Copied /opt/onemind/.env.tpl"
# Copy docker-compose for NATS + PostgreSQL
if [ -f "$REPO/docker-compose.yml" ]; then
  cp "$REPO/docker-compose.yml" /opt/onemind/docker-compose.yml
fi
echo "  Done."

# --- 7. Deploy scripts ---
echo "[7/8] Installing scripts and systemd units..."
# Scripts
cp "$REPO/openclaw/deploy/scripts/codex-watch.sh" /home/onemind/codex-watch.sh
cp "$REPO/openclaw/deploy/scripts/openclaw-reload.sh" /home/onemind/openclaw-reload.sh
cp "$REPO/openclaw/deploy/scripts/workspace-backup.sh" /home/onemind/workspace-backup.sh
chmod +x /home/onemind/codex-watch.sh /home/onemind/openclaw-reload.sh /home/onemind/workspace-backup.sh
chown onemind:onemind /home/onemind/codex-watch.sh /home/onemind/openclaw-reload.sh /home/onemind/workspace-backup.sh

# Systemd units
cp "$REPO/openclaw/deploy/systemd/"*.service /etc/systemd/system/
cp "$REPO/openclaw/deploy/systemd/"*.timer /etc/systemd/system/
systemctl daemon-reload
systemctl enable openclaw codex-watch openclaw-reload codex-pull.timer

# Cron job
(sudo -u onemind crontab -l 2>/dev/null | grep -v workspace-backup; echo "*/15 * * * * /home/onemind/workspace-backup.sh") | sudo -u onemind crontab -
echo "  Done."

# --- 8. Firewall ---
echo "[8/8] Configuring firewall..."
if command -v ufw &>/dev/null; then
  ufw allow 22/tcp comment "SSH"
  ufw allow in on tailscale0 to any port 4222 comment "NATS (Tailscale)"
  ufw allow in on tailscale0 to any port 18789 comment "OpenClaw Gateway (Tailscale)"
  ufw --force enable
fi
echo "  Done."

echo ""
echo "============================================"
echo "  Setup complete!"
echo "============================================"
echo ""
echo "Remaining manual steps:"
echo "  1. Add OP_SERVICE_ACCOUNT_TOKEN to /opt/onemind/.env.sa"
echo "  2. Create required items in 1Password (see deploy/env.tpl for item names)"
echo "  3. Copy openclaw.json.example to ~/.openclaw/openclaw.json and configure"
echo "  3. Run: sudo -u onemind gh auth login"
echo "  4. Install Tailscale: curl -fsSL https://tailscale.com/install.sh | sh && tailscale up"
echo "  5. Start Docker services: cd /opt/onemind && docker compose --profile data up -d"
echo "  6. Start OpenClaw: systemctl start openclaw codex-watch openclaw-reload codex-pull.timer"
echo "  7. Verify: systemctl status openclaw"
echo ""
