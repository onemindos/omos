#!/usr/bin/env bash
# Mint a short-lived (1 h) OpenChoreo API token for the read-only `ai-operator` account.
#
#   openchoreo-token.sh            prints a headers JSON object (Claude Code `headersHelper`)
#   openchoreo-token.sh --raw      prints the bare access token (for Codex:
#                                  export OPENCHOREO_MCP_TOKEN="$(scripts/ai-operator/openchoreo-token.sh --raw)")
#
# Credentials come from ~/.config/onemind/openchoreo-mcp.cred (written by setup.sh, mode 600:
# line 1 username, line 2 app password). Nothing is cached or logged.
set -euo pipefail

CRED_FILE="${OPENCHOREO_MCP_CRED:-$HOME/.config/onemind/openchoreo-mcp.cred}"
TOKEN_URL="${OPENCHOREO_TOKEN_URL:-https://onemind-auth.tail717ba0.ts.net/application/o/token/}"
CLIENT_ID="onemind-openchoreo"

[[ -r "$CRED_FILE" ]] || { echo "missing $CRED_FILE (run scripts/ai-operator/setup.sh)" >&2; exit 1; }
{ read -r username; read -r password; } <"$CRED_FILE"

response=$(curl -fsS --max-time 8 -X POST "$TOKEN_URL" \
  --data-urlencode "grant_type=client_credentials" \
  --data-urlencode "client_id=$CLIENT_ID" \
  --data-urlencode "username=$username" \
  --data-urlencode "password=$password" \
  --data-urlencode "scope=openid profile email") || { echo "token request failed" >&2; exit 1; }

python3 -c '
import json, sys
token = json.load(sys.stdin)["access_token"]
print(token if sys.argv[1] == "--raw" else json.dumps({"Authorization": "Bearer " + token}))
' "${1:-}" <<<"$response"
