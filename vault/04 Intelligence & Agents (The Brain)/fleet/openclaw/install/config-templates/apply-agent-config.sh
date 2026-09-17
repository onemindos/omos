#!/usr/bin/env bash
# apply-agent-config.sh
# Applies OpenClaw agent workspace files to VPS and sets config values via CLI.
# Run on VPS: bash apply-agent-config.sh
# Or from Mac (requires Tailscale active): ssh root@omos-core.tail717ba0.ts.net 'bash -s' < apply-agent-config.sh

set -euo pipefail

LOCAL_WORKSPACES="/home/ubuntu/onemind/_system/openclaw/workspaces"
OPENCLAW_WORKSPACES="/home/ubuntu/.openclaw"

echo "=== OpenClaw Agent Config Deploy ==="
echo ""

# ─── CONCURRENCY ────────────────────────────────────────────────────────────
echo "[1/4] Setting concurrency limits..."
openclaw config set agents.defaults.maxConcurrent 16
openclaw config set agents.defaults.subagents.maxConcurrent 12
echo "  ✓ Global: 16 | Sub-agent: 12"
echo ""

# ─── WORKSPACE FILE SYNC ────────────────────────────────────────────────────
echo "[2/4] Syncing workspace files..."

AGENTS=(main spartan oracle haven empire grid guardian trinity eden)

for agent in "${AGENTS[@]}"; do
  SRC="$LOCAL_WORKSPACES/$agent"
  DST="$OPENCLAW_WORKSPACES/workspace"
  [[ "$agent" != "main" ]] && DST="$OPENCLAW_WORKSPACES/workspace-$agent"

  if [[ -d "$SRC" ]]; then
    # Copy all .md files from local workspace to openclaw workspace
    for f in "$SRC"/*.md; do
      [[ -f "$f" ]] || continue
      cp "$f" "$DST/$(basename "$f")"
      echo "  ✓ $agent/$(basename "$f")"
    done
  else
    echo "  ⚠ Source not found: $SRC (skipping)"
  fi
done
echo ""

# ─── CRON JOBS ──────────────────────────────────────────────────────────────
echo "[3/4] Registering cron heartbeats..."

# Grid — every 15 min, silent ops check
openclaw cron add --agent grid --schedule "*/15 * * * *" \
  --message "Run HEARTBEAT checks. Log HEARTBEAT_OK or escalate to Legacy if alert condition met." \
  --label "grid-heartbeat" 2>/dev/null || echo "  ⚠ grid-heartbeat already exists or cron not available"

# Guardian — 06:00 + 22:00 ET (UTC-4 in summer = 10:00 + 02:00 UTC)
openclaw cron add --agent guardian --schedule "0 10 * * *" \
  --message "Run HEARTBEAT security checks. Log result or escalate to Legacy if alert condition met." \
  --label "guardian-heartbeat-am" 2>/dev/null || echo "  ⚠ guardian-heartbeat-am already exists"

openclaw cron add --agent guardian --schedule "0 2 * * *" \
  --message "Run HEARTBEAT security checks. Log result or escalate to Legacy if alert condition met." \
  --label "guardian-heartbeat-pm" 2>/dev/null || echo "  ⚠ guardian-heartbeat-pm already exists"

# Trinity — every 10 min, coordination check
openclaw cron add --agent trinity --schedule "*/10 * * * *" \
  --message "Run HEARTBEAT coordination checks. Report HEARTBEAT_OK or escalate to Legacy." \
  --label "trinity-heartbeat" 2>/dev/null || echo "  ⚠ trinity-heartbeat already exists"

# Spartan — morning + evening
openclaw cron add --agent spartan --schedule "0 7 * * *" \
  --message "Run HEARTBEAT morning check. Flag if training streak broken or alert condition met." \
  --label "spartan-heartbeat-am" 2>/dev/null || echo "  ⚠ spartan-heartbeat-am already exists"

openclaw cron add --agent spartan --schedule "0 20 * * *" \
  --message "Run HEARTBEAT evening check. Log workout status." \
  --label "spartan-heartbeat-pm" 2>/dev/null || echo "  ⚠ spartan-heartbeat-pm already exists"

# Oracle — midday
openclaw cron add --agent oracle --schedule "0 13 * * *" \
  --message "Run HEARTBEAT finance and mind check. Flag any alert conditions to Zeus." \
  --label "oracle-heartbeat" 2>/dev/null || echo "  ⚠ oracle-heartbeat already exists"

# Haven — each morning
openclaw cron add --agent haven --schedule "0 8 * * *" \
  --message "Run HEARTBEAT family and home check. Message Zeus if alert condition met." \
  --label "haven-heartbeat" 2>/dev/null || echo "  ⚠ haven-heartbeat already exists"

# Empire — each morning
openclaw cron add --agent empire --schedule "0 9 * * *" \
  --message "Run HEARTBEAT revenue and business check. Message Zeus if alert condition met." \
  --label "empire-heartbeat" 2>/dev/null || echo "  ⚠ empire-heartbeat already exists"

# Eden — Sunday mornings
openclaw cron add --agent eden --schedule "0 9 * * 0" \
  --message "Run HEARTBEAT home environment check. Escalate to Legacy if alert condition met." \
  --label "eden-heartbeat" 2>/dev/null || echo "  ⚠ eden-heartbeat already exists"

echo ""

# ─── RELOAD ─────────────────────────────────────────────────────────────────
echo "[4/4] Reloading agents..."
openclaw agent reload all 2>/dev/null && echo "  ✓ All agents reloaded" || \
  echo "  ⚠ Reload command not available — restart gateway manually: systemctl --user restart openclaw-gateway"
echo ""

echo "=== Deploy complete ==="
echo ""
echo "Verify with:"
echo "  openclaw agent list"
echo "  openclaw cron list"
