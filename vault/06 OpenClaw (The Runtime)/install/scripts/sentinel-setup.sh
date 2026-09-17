#!/bin/bash
# =============================================================================
# sentinel-setup.sh — Full Sentinel VPS Setup
# Run as root on the Sentinel droplet after cloud-init completes
# =============================================================================

set -e

MAIN_VPS="104.248.4.212"
SENTINEL_IP="159.203.122.175"
DD_API_KEY="${DD_API_KEY:-REPLACE_WITH_DATADOG_API_KEY}"
DISCORD_TOKEN="${SENTINEL_DISCORD_TOKEN:-REPLACE_WITH_SENTINEL_BOT_TOKEN}"
OPENCLAW_TOKEN="${SENTINEL_OPENCLAW_TOKEN:-REPLACE_WITH_SENTINEL_GATEWAY_TOKEN}"

echo "=== Sentinel VPS Setup ==="
echo "Sentinel IP: $SENTINEL_IP"

# ── 1. Tailscale ──────────────────────────────────────────────────────────────
echo "[1/7] Joining Tailscale..."
# Replace TS_AUTHKEY with a reusable Tailscale auth key from tailscale.com/admin/settings/keys
tailscale up --authkey="${TS_AUTHKEY:-REPLACE_WITH_TAILSCALE_AUTHKEY}" --hostname=onemind-sentinel --accept-routes 2>/dev/null || echo "Tailscale already up"

# ── 2. 1Password CLI ──────────────────────────────────────────────────────────
echo "[2/7] Installing 1Password CLI..."
curl -sS https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /usr/share/keyrings/1password-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/$(dpkg --print-architecture) stable main" | tee /etc/apt/sources.list.d/1password.list
mkdir -p /etc/debsig/policies/AC2D62742012EA22/ && curl -sS https://downloads.1password.com/linux/debian/debsig/1password.pol | tee /etc/debsig/policies/AC2D62742012EA22/1password.pol
mkdir -p /usr/share/debsig/keyrings/AC2D62742012EA22 && curl -sS https://downloads.1password.com/linux/keys/1password.asc | gpg --dearmor -o /usr/share/debsig/keyrings/AC2D62742012EA22/debsig.gpg
apt-get update -q && apt-get install -y -q 1password-cli

# ── 3. SSH key to reach main VPS ─────────────────────────────────────────────
echo "[3/7] Setting up SSH key for main VPS access..."
mkdir -p /home/onemind/.ssh
# Copy the key from main VPS (run this manually or copy key content here)
# ssh-copy-id -i /root/.ssh/id_ed25519.pub root@104.248.4.212 (from main)
echo "→ Manual step: copy SSH key to reach $MAIN_VPS"
echo "  On main VPS: cat /root/.ssh/id_ed25519.pub"
echo "  Then add to /root/.ssh/authorized_keys on sentinel or vice versa"

# ── 4. OpenClaw sentinel profile ──────────────────────────────────────────────
echo "[4/7] Setting up OpenClaw sentinel profile..."
mkdir -p /home/onemind/.openclaw-sentinel

# Create env file
cat > /opt/sentinel/.env.sa << ENVEOF
OPENCLAW_PROFILE=sentinel
OPENCLAW_STATE_DIR=/home/onemind/.openclaw-sentinel
OPENCLAW_CONFIG_PATH=/home/onemind/.openclaw-sentinel/openclaw.json
LOG_LEVEL=info
NATS_URL=nats://104.248.4.212:4222
ENVEOF

# ── 5. Datadog agent ──────────────────────────────────────────────────────────
echo "[5/7] Installing Datadog agent..."
DD_API_KEY="$DD_API_KEY" DD_SITE="datadoghq.com" bash -c "$(curl -L https://s3.amazonaws.com/dd-agent/scripts/install_script_agent7.sh)"

# Tag this host
cat >> /etc/datadog-agent/datadog.yaml << DDEOF
tags:
  - env:production
  - role:sentinel
  - project:onemindos
DDEOF

systemctl restart datadog-agent

# ── 6. systemd service ────────────────────────────────────────────────────────
echo "[6/7] Installing openclaw-sentinel systemd service..."
cat > /etc/systemd/system/openclaw-sentinel.service << SVCEOF
[Unit]
Description=OneMind OS — Sentinel (Rescue Gateway)
After=network-online.target
Wants=network-online.target

[Service]
User=onemind
WorkingDirectory=/home/onemind
EnvironmentFile=/opt/sentinel/.env.sa
ExecStart=/usr/bin/openclaw --profile sentinel gateway --port 19789
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SVCEOF

systemctl daemon-reload
systemctl enable openclaw-sentinel

# ── 7. Daily git pull for workspace ──────────────────────────────────────────
echo "[7/7] Setting up workspace sync..."
cat > /etc/cron.d/sentinel-git-pull << CRONEOF
*/5 * * * * onemind cd /home/onemind/onemind-os && git pull --ff-only origin main 2>/dev/null
CRONEOF

echo ""
echo "=== Setup complete ==="
echo "Next steps:"
echo "  1. op signin and configure service account token in /opt/sentinel/.env.sa"
echo "  2. openclaw --profile sentinel onboard (configure gateway)"
echo "  3. systemctl start openclaw-sentinel"
echo "  4. Connect Sentinel Discord bot token (SENTINEL_DISCORD_TOKEN)"
