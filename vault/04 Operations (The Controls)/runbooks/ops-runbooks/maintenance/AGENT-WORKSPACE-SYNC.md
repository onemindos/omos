---
type: runbook
domain: 03
name: Agent Workspace Sync
status: active
updated: 2026-05-11
---
# Agent Workspace Sync Runbook

## Overview

Agent workspace files (`~/.openclaw/workspace-{id}/`) are **not** auto-synced from the vault. The vault git sync (`onemind-pull.timer`) keeps `/home/ubuntu/onemind/` current, but agent workspaces are a separate directory and must be updated manually or via this procedure.

**When this is needed:** After editing any agent's openclaw workspace files in the vault (`_roster/_COMMAND/{agent}/openclaw/`).

**Files affected:** `AGENTS.md`, `IDENTITY.md`, `SOUL.md`, `TOOLS.md`, `HEARTBEAT.md`, `MEMORY.md`, `USER.md`

---

## Workspace Path Map

| Agent | OpenClaw Workspace |
|-------|--------------------|
| Legacy (`main`) | `~/.openclaw/workspace/` |
| Zenith | `~/.openclaw/workspace-zenith/` |
| Unity | `~/.openclaw/workspace-unity/` |
| Eternal | `~/.openclaw/workspace-eternal/` |
| Sovereignty | `~/.openclaw/workspace-sovereignty/` |
| Sentinel | `~/.openclaw/agents/sentinel/workspace/` *(on omos-sentinel)* |

**Vault source:** `/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/{agent}/openclaw/`

---

## Sync Procedure

### 1. Verify vault is current on VPS

```bash
ssh onemindos-core "cd /home/ubuntu/onemind && git log --oneline -1"
```

Should match latest commit on GitHub. If behind, wait for `onemind-pull.timer` (runs every 10s) or force pull:

```bash
ssh onemindos-core "cd /home/ubuntu/onemind && git pull --rebase"
```

### 2. Copy workspace files for a single agent

```bash
ssh onemindos-core '
VAULT="/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND"
AGENT="legacy"   # change to: zenith | unity | eternal | sovereignty
WS="/home/ubuntu/.openclaw/workspace"  # change per agent — see map above

cp "$VAULT/$AGENT/openclaw/AGENTS.md"    "$WS/AGENTS.md"
cp "$VAULT/$AGENT/openclaw/IDENTITY.md"  "$WS/IDENTITY.md"
cp "$VAULT/$AGENT/openclaw/SOUL.md"      "$WS/SOUL.md"
cp "$VAULT/$AGENT/openclaw/TOOLS.md"     "$WS/TOOLS.md"
cp "$VAULT/$AGENT/openclaw/HEARTBEAT.md" "$WS/HEARTBEAT.md"
cp "$VAULT/$AGENT/openclaw/MEMORY.md"    "$WS/MEMORY.md"
cp "$VAULT/$AGENT/openclaw/USER.md"      "$WS/USER.md"
echo "Done: $AGENT"
'
```

### 3. Sync all 5 omos-core agents at once

```bash
ssh onemindos-core '
VAULT="/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND"
declare -A WS=(
  [legacy]="/home/ubuntu/.openclaw/workspace"
  [zenith]="/home/ubuntu/.openclaw/workspace-zenith"
  [unity]="/home/ubuntu/.openclaw/workspace-unity"
  [eternal]="/home/ubuntu/.openclaw/workspace-eternal"
  [sovereignty]="/home/ubuntu/.openclaw/workspace-sovereignty"
)
for agent in "${!WS[@]}"; do
  src="$VAULT/$agent/openclaw"
  dst="${WS[$agent]}"
  for f in AGENTS.md IDENTITY.md SOUL.md TOOLS.md HEARTBEAT.md MEMORY.md USER.md; do
    [ -f "$src/$f" ] && cp "$src/$f" "$dst/$f"
  done
  echo "Synced: $agent"
done
'
```

### 4. Verify sync

```bash
ssh onemindos-core 'python3 -c "
import os
mappings = {
    \"legacy\":      (\"/home/ubuntu/.openclaw/workspace/AGENTS.md\", \"/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/legacy/openclaw/AGENTS.md\"),
    \"zenith\":      (\"/home/ubuntu/.openclaw/workspace-zenith/AGENTS.md\", \"/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/zenith/openclaw/AGENTS.md\"),
    \"unity\":       (\"/home/ubuntu/.openclaw/workspace-unity/AGENTS.md\", \"/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/unity/openclaw/AGENTS.md\"),
    \"eternal\":     (\"/home/ubuntu/.openclaw/workspace-eternal/AGENTS.md\", \"/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/eternal/openclaw/AGENTS.md\"),
    \"sovereignty\": (\"/home/ubuntu/.openclaw/workspace-sovereignty/AGENTS.md\", \"/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/_roster/_COMMAND/sovereignty/openclaw/AGENTS.md\"),
}
for agent, (ws, vault) in mappings.items():
    match = open(ws).read() == open(vault).read() if os.path.exists(ws) and os.path.exists(vault) else False
    print(f\"{agent:15}| {\"MATCH\" if match else \"DRIFT DETECTED\"}\")
"'
```

### 5. Restart gateway (only needed if IDENTITY/SOUL/TOOLS changed)

AGENTS.md, HEARTBEAT.md, MEMORY.md, USER.md changes take effect on next agent session — no restart needed.
If IDENTITY.md, SOUL.md, or TOOLS.md changed:

```bash
ssh onemindos-core "openclaw gateway restart"
```

---

## Automation Note

A future enhancement: add a post-push hook or `onemind-watch.service` trigger that auto-copies vault agent files to workspaces on change. Until then, run Step 3 above after any vault agent workspace edit.
