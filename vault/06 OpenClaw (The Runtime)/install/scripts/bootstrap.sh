#!/bin/bash
# bootstrap.sh — Deploy an OpenClaw node from scratch
#
# Usage:
#   curl -sL https://raw.githubusercontent.com/OneMind-OS/openclaw-core-workspaces/main/bootstrap.sh | bash
# OR after cloning:
#   bash bootstrap.sh
#
# Prerequisites:
#   - Debian 12 / Ubuntu 22+
#   - Root or sudo access
#   - GitHub SSH key configured
#   - OP_SERVICE_ACCOUNT_TOKEN exported (from 1Password SO vault)

set -euo pipefail

ONEMIND_USER=ubuntu
OPENCLAW_DIR=/home/$ONEMIND_USER/.openclaw
REPO_DIR=/home/$ONEMIND_USER/openclaw-workspaces
SCRIPTS_DIR=/home/$ONEMIND_USER
SYSTEMD_DIR=/etc/systemd/system
NODE_VERSION=20

# ── 1. Detect VM (core vs sentinel) ───────────────────────────────────────
HOSTNAME_SHORT=$(hostname -s)
case "$HOSTNAME_SHORT" in
  *core*)     WORKSPACE_REPO="git@github.com:OneMind-OS/openclaw-core-workspaces.git" ;;
  *sentinel*) WORKSPACE_REPO="git@github.com:OneMind-OS/openclaw-sentinel-workspaces.git" ;;
  *)
    echo "Unknown hostname: $HOSTNAME_SHORT"
    read -rp "Workspace repo SSH URL: " WORKSPACE_REPO
    ;;
esac
echo "[bootstrap] Repo: $WORKSPACE_REPO"

# ── 2. Create onemind user if missing ─────────────────────────────────────
if ! id "$ONEMIND_USER" &>/dev/null; then
  useradd -m -s /bin/bash "$ONEMIND_USER"
  echo "[bootstrap] Created user $ONEMIND_USER"
fi

# ── 3. Install dependencies ───────────────────────────────────────────────
apt-get update -qq
apt-get install -y -qq curl git rsync inotify-tools jq unzip

if ! command -v node &>/dev/null; then
  curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash -
  apt-get install -y nodejs
fi

if ! command -v op &>/dev/null; then
  curl -sS https://downloads.1password.com/linux/keys/1password.asc | \
    gpg --dearmor --output /usr/share/keyrings/1password-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/$(dpkg --print-architecture) stable main" | \
    tee /etc/apt/sources.list.d/1password.list
  apt-get update -qq && apt-get install -y 1password-cli
fi

if ! command -v gh &>/dev/null; then
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | \
    dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | \
    tee /etc/apt/sources.list.d/github-cli.list > /dev/null
  apt-get update && apt-get install gh -y
fi

if ! command -v openclaw &>/dev/null; then
  npm install -g @openclaw/cli@latest
fi

echo "[bootstrap] Dependencies installed"

# ── 4. Clone workspace repo ───────────────────────────────────────────────
if [ ! -d "$REPO_DIR/.git" ]; then
  sudo -u "$ONEMIND_USER" git clone "$WORKSPACE_REPO" "$REPO_DIR"
  echo "[bootstrap] Workspace repo cloned to $REPO_DIR"
else
  sudo -u "$ONEMIND_USER" git -C "$REPO_DIR" pull --rebase
  echo "[bootstrap] Workspace repo updated"
fi

# ── 5. Deploy config and workspaces ───────────────────────────────────────
mkdir -p "$OPENCLAW_DIR"
cp "$REPO_DIR/openclaw.json" "$OPENCLAW_DIR/openclaw.json"

for ws_dir in "$REPO_DIR"/workspace*/; do
  ws_name=$(basename "$ws_dir")
  rsync -a "$ws_dir" "$OPENCLAW_DIR/$ws_name/"
done

chown -R "$ONEMIND_USER:$ONEMIND_USER" "$OPENCLAW_DIR"
echo "[bootstrap] Config and workspaces deployed"

# ── 6. Deploy scripts ─────────────────────────────────────────────────────
cp "$REPO_DIR/scripts/openclaw-reload.sh" "$SCRIPTS_DIR/openclaw-reload.sh"
cp "$REPO_DIR/scripts/openclaw-config-sync.sh" "$SCRIPTS_DIR/openclaw-config-sync.sh"
chmod +x "$SCRIPTS_DIR/openclaw-reload.sh" "$SCRIPTS_DIR/openclaw-config-sync.sh"
chown "$ONEMIND_USER:$ONEMIND_USER" \
  "$SCRIPTS_DIR/openclaw-reload.sh" \
  "$SCRIPTS_DIR/openclaw-config-sync.sh"

# ── 7. Install and enable systemd services ────────────────────────────────
for svc in openclaw openclaw-reload openclaw-config-sync; do
  cp "$REPO_DIR/systemd/$svc.service" "$SYSTEMD_DIR/$svc.service"
done
systemctl daemon-reload
systemctl enable openclaw openclaw-reload openclaw-config-sync

# ── 8. Start services ─────────────────────────────────────────────────────
systemctl start openclaw
sleep 3
systemctl start openclaw-reload openclaw-config-sync

echo ""
echo "═══════════════════════════════════════════════════════"
echo "  OpenClaw node deployed successfully"
echo "  Gateway: http://localhost:18789"
echo "  Status: systemctl status openclaw"
echo "  Logs:   journalctl -u openclaw -f"
echo "═══════════════════════════════════════════════════════"
