# OneMind OS — 1Password Secrets Template
# Used by: op run --env-file=/opt/onemind/.env.tpl -- openclaw gateway
#
# Format: KEY={{ op://vault-name/item-name/field-name
# All items live in vault: "00-24 UI (Unified Intelligence)"
#
# Items created by deploy/scripts/create-1p-items.sh:
#   "Azure OpenAI"    → api-key, endpoint
#   "VPS OpenClaw"    → gateway-token, hooks-token, mcp-client-secret
#   "Discord Bot"     → token
#   "ElevenLabs"      → api-key, voice-id  (REMOVED — Azure OpenAI TTS used instead)
#   "Anthropic"       → api-key
#   "OpenAI"          → api-key, tts-base-url
#   "Deepgram"        → api-key  (STT for voice pipeline)
#   "GitHub"          → token
#   "LiveKit"         → api-key, api-secret, studio-access-token
#   "OmniDB"          → url, db, user, password
#   "BlueBubbles"     → password  (not yet in vault — in .env.sa)
#   "Edge Devices"    → farm-password, home-password, office-password  (not yet in vault — TODO)
#   "N8N"             → encryption-key  (not yet in vault — in .env.sa)
#   "Gradient"        → api-key
#   "VPS OP Vault"    → vault-id
#
# Static/non-secret vars (NATS_URL, LOG_LEVEL, channel IDs, etc.)
# live in /opt/onemind/.env.sa — NOT here.
#
# OP_SERVICE_ACCOUNT_TOKEN is also in .env.sa — NOT here.
#
# Skill-specific (uncomment after creating the 1Password item):
#   "ClickUp"         → api-key
#   "Lob"             → api-key
#   "Twilio"          → account-sid, auth-token, phone-number
#   "Resend"          → api-key
#   "Spotify"         → client-id, client-secret, refresh-token
#   "Home Assistant"  → token
#   "Todoist"         → api-token

# ── LLM ─────────────────────────────────────────────────────────────────────
AZURE_OPENAI_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Azure OpenAI/api-key
AZURE_OPENAI_ENDPOINT=op://aaeyris5z44pdqspk277u2oppi/Azure OpenAI/endpoint
ANTHROPIC_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Anthropic/api-key
OPENAI_API_KEY=op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key
# OPENAI_TTS_BASE_URL — no longer needed; TTS now routed via AZURE_OPENAI_ENDPOINT
# OPENAI_TTS_BASE_URL=op://aaeyris5z44pdqspk277u2oppi/OpenAI/tts-base-url

# ── OpenClaw Gateway ─────────────────────────────────────────────────────────
OPENCLAW_GATEWAY_TOKEN=op://aaeyris5z44pdqspk277u2oppi/VPS OpenClaw/gateway-token
OPENCLAW_HOOKS_TOKEN=op://aaeyris5z44pdqspk277u2oppi/VPS OpenClaw/hooks-token
MCP_CLIENT_SECRET=op://aaeyris5z44pdqspk277u2oppi/VPS OpenClaw/mcp-client-secret

# ── Discord ──────────────────────────────────────────────────────────────────
DISCORD_BOT_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Discord Bot/token

# ── Voice ────────────────────────────────────────────────────────────────────
# ElevenLabs removed — Azure OpenAI TTS used instead (covered by Azure credits)
# ELEVENLABS_API_KEY=op://aaeyris5z44pdqspk277u2oppi/ElevenLabs/api-key
# ELEVENLABS_VOICE_ID=op://aaeyris5z44pdqspk277u2oppi/ElevenLabs/voice-id

# ── Voice Pipeline (STT + TTS) ───────────────────────────────────────────────
# Deepgram — Speech-to-Text (Nova-3, used by LiveKit voice agent)
DEEPGRAM_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Deepgram/api-key

# ── LiveKit ──────────────────────────────────────────────────────────────────
LIVEKIT_API_KEY=op://aaeyris5z44pdqspk277u2oppi/LiveKit/api-key
LIVEKIT_API_SECRET=op://aaeyris5z44pdqspk277u2oppi/LiveKit/api-secret
STUDIO_ACCESS_TOKEN=op://aaeyris5z44pdqspk277u2oppi/LiveKit/studio-access-token

# ── GitHub ───────────────────────────────────────────────────────────────────
GITHUB_TOKEN=op://aaeyris5z44pdqspk277u2oppi/GitHub/token

# ── Database ─────────────────────────────────────────────────────────────────
OMNIDB_URL=op://aaeyris5z44pdqspk277u2oppi/OmniDB/url
OMNIDB_DB=op://aaeyris5z44pdqspk277u2oppi/OmniDB/db
OMNIDB_USER=op://aaeyris5z44pdqspk277u2oppi/OmniDB/user
OMNIDB_PASSWORD=op://aaeyris5z44pdqspk277u2oppi/OmniDB/password

# ── Smart Home ───────────────────────────────────────────────────────────────
BLUEBUBBLES_PASSWORD=op://aaeyris5z44pdqspk277u2oppi/BlueBubbles/password
# EDGE_* — items not yet created in vault; TODO: create 1P item
# EDGE_FARM_PASSWORD=op://aaeyris5z44pdqspk277u2oppi/Edge Devices/farm-password
# EDGE_HOME_PASSWORD=op://aaeyris5z44pdqspk277u2oppi/Edge Devices/home-password
# EDGE_OFFICE_PASSWORD=op://aaeyris5z44pdqspk277u2oppi/Edge Devices/office-password

# ── Automation ───────────────────────────────────────────────────────────────
# N8N_ENCRYPTION_KEY — not yet in 1Password vault; stored in .env.sa
# N8N_ENCRYPTION_KEY=op://aaeyris5z44pdqspk277u2oppi/N8N/encryption-key
GRADIENT_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Gradient/api-key

# ── 1Password Skill ──────────────────────────────────────────────────────────
OP_VAULT_ID=op://aaeyris5z44pdqspk277u2oppi/VPS OP Vault/vault-id

# ── Command Center Frontends ─────────────────────────────────────────────────
# SwarmClaw (agent orchestration dashboard)
# CREDENTIAL_SECRET=op://aaeyris5z44pdqspk277u2oppi/SwarmClaw VPS/credential-secret
# ACCESS_KEY=op://aaeyris5z44pdqspk277u2oppi/SwarmClaw VPS/access-key
# OpenClaw Studio (native agent dashboard)
# STUDIO_ACCESS_TOKEN=op://aaeyris5z44pdqspk277u2oppi/OpenClaw Studio VPS/access-token

# ── Skill-specific (uncomment after creating the 1Password item) ─────────────
# CLICKUP_API_KEY=op://aaeyris5z44pdqspk277u2oppi/ClickUp/api-key
# LOB_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Lob/api-key
# TWILIO_ACCOUNT_SID=op://aaeyris5z44pdqspk277u2oppi/Twilio/account-sid
# TWILIO_AUTH_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Twilio/auth-token
# TWILIO_PHONE_NUMBER=op://aaeyris5z44pdqspk277u2oppi/Twilio/phone-number
# RESEND_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Resend/api-key
# SPOTIFY_CLIENT_ID=op://aaeyris5z44pdqspk277u2oppi/Spotify/client-id
# SPOTIFY_CLIENT_SECRET=op://aaeyris5z44pdqspk277u2oppi/Spotify/client-secret
# SPOTIFY_REFRESH_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Spotify/refresh-token
# HOME_ASSISTANT_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Home Assistant/token
CLOUDFLARE_API_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Cloudflare - OneMind OS/Global API Key
# HOME_ASSISTANT_URL=http://homeassistant.local:8123
# FRIGATE_URL=http://frigate.local:5000
# TODOIST_API_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Todoist/api-token
# ZOHO_MCP_URL=op://aaeyris5z44pdqspk277u2oppi/Zoho MCP/url
# DIGITALOCEAN_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Digital Ocean -  API Personal token/credential
# DD_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Datadog API Key/credential
# SENTINEL_DISCORD_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Sentinel Discord Bot/token

# ── AWS Bedrock ───────────────────────────────────────────────────────────────
AWS_ACCESS_KEY_ID=op://aaeyris5z44pdqspk277u2oppi/AWS Bedrock/username
AWS_SECRET_ACCESS_KEY=op://aaeyris5z44pdqspk277u2oppi/AWS Bedrock/password
AWS_DEFAULT_REGION=us-east-1
