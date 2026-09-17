#!/usr/bin/env bash
# create-1p-items.sh — Migrate /opt/onemind/.env secrets to 1Password,
# patch .env.sa with static vars, test op run, and restart openclaw.
#
# Run ONCE after adding OP_SERVICE_ACCOUNT_TOKEN to /opt/onemind/.env.sa:
#   bash /opt/onemind/create-1p-items.sh

set -euo pipefail

VAULT="00-24 UI (Unified Intelligence)"
OLD_ENV="/opt/onemind/.env"
ENV_SA="/opt/onemind/.env.sa"
ENV_TPL="/opt/onemind/.env.tpl"

# ── Guards ────────────────────────────────────────────────────────────────────
if [[ ! -f "$OLD_ENV" ]]; then
  echo "ERROR: $OLD_ENV not found — nothing to migrate."
  exit 1
fi
if [[ ! -f "$ENV_SA" ]]; then
  echo "ERROR: $ENV_SA not found — run setup.sh first."
  exit 1
fi
if [[ ! -f "$ENV_TPL" ]]; then
  echo "ERROR: $ENV_TPL not found — deploy env.tpl first."
  exit 1
fi

# ── Load SA token ─────────────────────────────────────────────────────────────
# shellcheck disable=SC1090
source "$ENV_SA"
if [[ -z "${OP_SERVICE_ACCOUNT_TOKEN:-}" ]]; then
  echo ""
  echo "ERROR: OP_SERVICE_ACCOUNT_TOKEN is not set in $ENV_SA"
  echo ""
  echo "Add it now:"
  echo "  read -s -p 'SA token: ' T && echo \"\" && sed -i \"s|^OP_SERVICE_ACCOUNT_TOKEN=.*|OP_SERVICE_ACCOUNT_TOKEN=\$T|\" $ENV_SA"
  echo ""
  exit 1
fi
export OP_SERVICE_ACCOUNT_TOKEN
echo "✓ Service account token loaded"

# ── Helpers ───────────────────────────────────────────────────────────────────
# Read a value from old .env (handles quoted and unquoted values)
get_env() {
  grep -m1 "^${1}=" "$OLD_ENV" 2>/dev/null | cut -d= -f2- | sed 's/^["'"'"']//;s/["'"'"']$//' || true
}

# Create item only if it doesn't already exist
create_item() {
  local title="$1"; shift
  if op item get "$title" --vault "$VAULT" &>/dev/null 2>&1; then
    echo "  ↷  '$title' already exists — skipping"
  else
    op item create \
      --vault "$VAULT" \
      --category "API Credential" \
      --title "$title" \
      "$@" &>/dev/null
    echo "  ✓  Created '$title'"
  fi
}

# Append a KEY=VALUE to .env.sa only if not already present
append_static() {
  local key="$1"
  local val
  val=$(get_env "$key")
  if [[ -z "$val" ]]; then
    return  # not in old .env, skip
  fi
  if grep -q "^${key}=" "$ENV_SA" 2>/dev/null; then
    echo "  ↷  $key already in .env.sa"
  else
    echo "${key}=${val}" >> "$ENV_SA"
    echo "  ✓  $key → .env.sa"
  fi
}

# ── Step 1: Create 1Password items ────────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════"
echo " Step 1 — Creating 1Password items"
echo "══════════════════════════════════════════════════════"

create_item "Azure OpenAI" \
  "api-key[password]=$(get_env AZURE_OPENAI_API_KEY)" \
  "endpoint[text]=$(get_env AZURE_OPENAI_ENDPOINT)"

create_item "VPS OpenClaw" \
  "gateway-token[password]=$(get_env OPENCLAW_GATEWAY_TOKEN)" \
  "hooks-token[password]=$(get_env OPENCLAW_HOOKS_TOKEN)" \
  "mcp-client-secret[password]=$(get_env MCP_CLIENT_SECRET)"

create_item "Discord Bot" \
  "token[password]=$(get_env DISCORD_BOT_TOKEN)"

create_item "ElevenLabs" \
  "api-key[password]=$(get_env ELEVENLABS_API_KEY)" \
  "voice-id[text]=$(get_env ELEVENLABS_VOICE_ID)"

create_item "Deepgram" \
  "api-key[password]=$(get_env DEEPGRAM_API_KEY)"

create_item "Anthropic" \
  "api-key[password]=$(get_env ANTHROPIC_API_KEY)"

create_item "OpenAI" \
  "api-key[password]=$(get_env OPENAI_API_KEY)" \
  "tts-base-url[text]=$(get_env OPENAI_TTS_BASE_URL)"

create_item "GitHub" \
  "token[password]=$(get_env GITHUB_TOKEN)"

create_item "LiveKit" \
  "api-key[password]=$(get_env LIVEKIT_API_KEY)" \
  "api-secret[password]=$(get_env LIVEKIT_API_SECRET)" \
  "studio-access-token[password]=$(get_env STUDIO_ACCESS_TOKEN)"

create_item "OmniDB" \
  "url[text]=$(get_env OMNIDB_URL)" \
  "db[text]=$(get_env OMNIDB_DB)" \
  "user[text]=$(get_env OMNIDB_USER)" \
  "password[password]=$(get_env OMNIDB_PASSWORD)"

create_item "BlueBubbles" \
  "password[password]=$(get_env BLUEBUBBLES_PASSWORD)"

create_item "Edge Devices" \
  "farm-password[password]=$(get_env EDGE_FARM_PASSWORD)" \
  "home-password[password]=$(get_env EDGE_HOME_PASSWORD)" \
  "office-password[password]=$(get_env EDGE_OFFICE_PASSWORD)"

create_item "N8N" \
  "encryption-key[password]=$(get_env N8N_ENCRYPTION_KEY)"

create_item "Gradient" \
  "api-key[password]=$(get_env GRADIENT_API_KEY)"

# VPS OP Vault item — stores the vault ID itself (used by the OP skill)
echo ""
echo "  Resolving vault ID for 'VPS OP Vault' item..."
VAULT_ID=$(op vault get "$VAULT" --format json | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
create_item "VPS OP Vault" \
  "vault-id[text]=$VAULT_ID"

# ── Step 2: Add static vars to .env.sa ───────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════"
echo " Step 2 — Patching .env.sa with static runtime vars"
echo "══════════════════════════════════════════════════════"
echo "" >> "$ENV_SA"
echo "# --- Runtime config (non-secret, added by create-1p-items.sh) ---" >> "$ENV_SA"

append_static "DISCORD_COMMUNITY_CHANNEL_ID"
append_static "HEARTBEAT_CHANNEL_ID"
append_static "ZEUS_DISCORD_USER_ID"
append_static "BRIDGE_AGENT_ID"
append_static "FAMILY_WHATSAPP_NUMBERS"
append_static "LEGACY_MODEL"
append_static "LIVEKIT_URL"
append_static "LOG_LEVEL"
append_static "N8N_WEBHOOK_URL"
append_static "OPENCLAW_MCP_GATEWAY_URL"
append_static "OPENCLAW_WORKSPACE"

# ── Step 3: Test op run resolution ───────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════"
echo " Step 3 — Testing op run secret resolution"
echo "══════════════════════════════════════════════════════"
RESOLVED=$(op run --env-file="$ENV_TPL" -- env 2>&1 | grep -E "^(AZURE|OPENCLAW_GATEWAY|DISCORD_BOT|ELEVENLABS_API)" | sed 's/=.*/=<RESOLVED>/' || true)
if [[ -z "$RESOLVED" ]]; then
  echo "  ERROR: op run returned no resolved vars — check env.tpl references"
  echo "  Run manually: op run --env-file=$ENV_TPL -- env"
  exit 1
fi
echo "$RESOLVED"
echo "  ✓  op run resolution working"

# ── Step 4: Reload and restart openclaw ──────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════"
echo " Step 4 — Reloading openclaw service"
echo "══════════════════════════════════════════════════════"
systemctl daemon-reload
systemctl restart openclaw
sleep 4

STATUS=$(systemctl is-active openclaw || true)
if [[ "$STATUS" == "active" ]]; then
  echo "  ✓  openclaw is active"
  systemctl status openclaw --no-pager | grep -E "Active:|Main PID:" | sed 's/^/  /'
else
  echo "  ✗  openclaw status: $STATUS"
  echo ""
  journalctl -u openclaw -n 30 --no-pager
  exit 1
fi

# ── Done ──────────────────────────────────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════"
echo " ✅  Migration complete!"
echo "══════════════════════════════════════════════════════"
echo ""
echo "  All secrets now live in 1Password vault: '$VAULT'"
echo "  Only the SA token remains on disk: $ENV_SA"
echo ""
echo "  Verify it's working:"
echo "    journalctl -u openclaw -f"
echo ""
echo "  Once confirmed, delete the old plaintext env:"
echo "    rm /opt/onemind/.env"
echo ""
