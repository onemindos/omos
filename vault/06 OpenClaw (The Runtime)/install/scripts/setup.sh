#!/bin/bash
# setup.sh — Full OpenClaw VPS deployment from scratch
# Run as root on a fresh Ubuntu 24.04 instance.
#
# Prerequisites:
#   - Fresh Ubuntu 24.04 (AWS Lightsail)
#   - Root SSH access
#   - GitHub repo: github.com/OneMind-OS/One-Mind
#
# Usage:
#   scp -r openclaw/deploy/ root@YOUR_IP:/tmp/deploy/
#   ssh root@YOUR_IP 'bash /tmp/deploy/scripts/setup.sh'
#
# After running, you still need to:
#   1. Copy openclaw.json to ~/.openclaw/openclaw.json (fill in real values)
#   2. Create /opt/onemind/.env with API keys
#   3. Run: sudo -u ubuntu gh auth login
#   4. Install Tailscale: curl -fsSL https://tailscale.com/install.sh | sh && tailscale up
#   5. Start services: systemctl start openclaw onemind-watch onemind-pull.timer

set -euo pipefail

echo "============================================"
echo "  One Mind — OpenClaw VPS Setup"
echo "============================================"
echo ""

# --- 1. System packages ---
echo "[1/7] Installing system packages..."
apt-get update -qq
apt-get install -y -qq git curl docker.io docker-compose-v2 inotify-tools jq
systemctl enable docker --now
usermod -aG docker ubuntu
echo "  Done."

# --- 2. Node.js + OpenClaw + 1Password CLI ---
echo "[2/7] Installing Node.js 22 + OpenClaw + 1Password CLI..."
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

# --- 3. Clone vault repo ---
echo "[3/7] Cloning One Mind vault..."
VAULT=/home/ubuntu/onemind
if [ ! -d "$VAULT" ]; then
  sudo -u ubuntu git clone https://TOKEN@github.com/OneMind-OS/One-Mind.git "$VAULT"
  echo "  NOTE: Replace TOKEN with actual GitHub token in remote URL"
else
  echo "  Vault already exists at $VAULT"
fi
echo "  Done."

# --- 4. Initialize OpenClaw ---
echo "[4/7] Initializing OpenClaw directory..."
sudo -u ubuntu mkdir -p /home/ubuntu/.openclaw
echo "  Done."

# --- 5. /opt/onemind structure ---
echo "[5/7] Setting up /opt/onemind..."
mkdir -p /opt/onemind
if [ ! -f /opt/onemind/.env.sa ]; then
  echo "  Created /opt/onemind/.env.sa — ADD OP_SERVICE_ACCOUNT_TOKEN"
else
  echo "  /opt/onemind/.env.sa already exists"
fi
echo "  Done."

# --- 6. Deploy scripts + systemd ---
echo "[6/7] Installing scripts and systemd units..."
# Watch script
cp scripts/onemind-watch.sh /usr/local/bin/onemind-watch.sh
chmod +x /usr/local/bin/onemind-watch.sh
chown ubuntu:ubuntu /usr/local/bin/onemind-watch.sh

# Systemd units
cp systemd/onemind-watch.service /etc/systemd/system/
cp systemd/onemind-pull.service /etc/systemd/system/
cp systemd/onemind-pull.timer /etc/systemd/system/
cp systemd/openclaw.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable openclaw onemind-watch onemind-pull.timer
echo "  Done."

# --- 7. Firewall ---
echo "[7/7] Configuring firewall..."
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
echo "  2. Create required items in 1Password (see env.tpl for item names)"
echo "  3. Copy openclaw.json.example to ~/.openclaw/openclaw.json and configure"
echo "  4. Run: sudo -u ubuntu gh auth login"
echo "  5. Install Tailscale: curl -fsSL https://tailscale.com/install.sh | sh && tailscale up"
echo "  6. Start Docker services: cd /opt/onemind && docker compose --profile data up -d"
echo "  7. Start OpenClaw: systemctl start openclaw onemind-watch onemind-pull.timer"
echo "  8. Verify: systemctl status openclaw onemind-watch onemind-pull.timer"
echo ""
