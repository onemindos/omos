---
title: "OpenClaw — Deployment Management"
type: sop
status: active
created: 2026-04-12
owner: zeus
tags: [type/sop, domain/SO, domain/openclaw, status/active]
ai_summary: "How to push workspace configs, skills, shared context, and full openclaw config to the VPS — rsync, scp patterns, deploy script."
---

# OpenClaw — Deployment Management

How workspace configs, skills, shared context, and the full OpenClaw config get pushed from Mac to the VPS.

---

## VPS Directory Structure

```
/home/onemind/onemind-os/
├── codex/                              ← Vault (git pull from GitHub)
├── openclaw/
│   ├── config.yaml                     ← Main OpenClaw config
│   ├── workspaces/
│   │   ├── legacy/
│   │   │   ├── IDENTITY.md
│   │   │   ├── AGENTS.md
│   │   │   ├── SOUL.md
│   │   │   ├── CODEX.md               ← symlink → shared/CODEX.md
│   │   │   ├── USER.md                ← symlink → shared/USER.md
│   │   │   └── skills/
│   │   ├── spartan/
│   │   ├── oracle/
│   │   ├── eden/
│   │   ├── haven/
│   │   ├── empire/
│   │   ├── guardian/
│   │   ├── grid/
│   │   └── trinity/
│   └── shared/
│       ├── CODEX.md                    ← symlinked into all workspaces
│       └── USER.md                     ← symlinked into all workspaces
```

**Mac source of truth:** `~/Documents/onemind-codex/00-24 SO/01 Intelligence/openclaw/`

---

## Push a Single File

```bash
VPS="root@104.248.4.212"
REMOTE="/home/onemind/onemind-os/openclaw"

# Push an updated workspace identity file
scp /Users/zeuslegacy/Documents/onemind-codex/"00-24 SO (Sovereign Operations)"/01\ Intelligence\ \(AI,\ Memory\ \&\ Tools\)/openclaw/workspaces/legacy/IDENTITY.md \
    $VPS:$REMOTE/workspaces/legacy/IDENTITY.md

# Push a new skill to Legacy
scp -r /Users/zeuslegacy/Documents/onemind-codex/.claude/skills/signals/ \
    $VPS:$REMOTE/workspaces/legacy/skills/signals/

# Push updated main config
scp /path/to/config.yaml $VPS:$REMOTE/config.yaml
```

After pushing config — reload without full restart (if supported), otherwise restart:

```bash
ssh $VPS 'systemctl reload openclaw || systemctl restart openclaw'
```

---

## Bulk Sync (rsync)

Use rsync to push all workspace files at once. **Always dry-run first (`-n`).**

```bash
LOCAL="/Users/zeuslegacy/Documents/onemind-codex/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools)/openclaw"
REMOTE="root@104.248.4.212:/home/onemind/onemind-os/openclaw"

# Dry run — see what would change
rsync -avzn "$LOCAL/" "$REMOTE/"

# Execute
rsync -avz "$LOCAL/" "$REMOTE/"
```

---

## Update Shared Context (CODEX.md / USER.md)

Shared files are symlinked from `shared/` into every workspace. Update the source:

```bash
VPS="root@104.248.4.212"
SHARED="/home/onemind/onemind-os/openclaw/shared"

# Update CODEX.md across all agents
scp openclaw/shared/CODEX.md $VPS:$SHARED/CODEX.md

# Update USER.md across all agents
scp openclaw/shared/USER.md $VPS:$SHARED/USER.md

# Verify symlinks are intact (should show → shared/CODEX.md)
ssh $VPS 'ls -la /home/onemind/onemind-os/openclaw/workspaces/legacy/'
```

No re-linking needed — symlinks already point to `shared/`. The workspace picks up the new file immediately.

---

## Full Deploy Script

At `_system/openclaw/deploy/scripts/openclaw-reload.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

VPS="root@104.248.4.212"
LOCAL="/Users/zeuslegacy/Documents/onemind-codex/00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools)/openclaw"
REMOTE="/home/onemind/onemind-os/openclaw"

echo "Syncing OpenClaw workspace to VPS..."
rsync -avz --delete "$LOCAL/" "$VPS:$REMOTE/"

echo "Restarting OpenClaw..."
ssh "$VPS" "systemctl restart openclaw"

echo "Checking status..."
ssh "$VPS" "systemctl status openclaw --no-pager -l | head -20"

echo "Done."
```

Run via VS Code task: **Deploy: Push to VPS**

---

## Add a New Agent Workspace

```bash
VPS="root@104.248.4.212"
AGENT="new-agent-name"
REMOTE="/home/onemind/onemind-os/openclaw"

# 1. Create workspace dir
ssh $VPS "mkdir -p $REMOTE/workspaces/$AGENT/memory"

# 2. Push workspace files
scp IDENTITY.md SOUL.md AGENTS.md BOOT.md $VPS:$REMOTE/workspaces/$AGENT/

# 3. Symlink shared context
ssh $VPS "
  ln -sf $REMOTE/shared/CODEX.md $REMOTE/workspaces/$AGENT/CODEX.md
  ln -sf $REMOTE/shared/USER.md  $REMOTE/workspaces/$AGENT/USER.md
"

# 4. Add agent to OpenClaw config.yaml
# 5. Reload OpenClaw
ssh $VPS 'systemctl restart openclaw'
```

---

## onemind-web Deployment

For the web UI (Docker on VPS):

```bash
# SSH to VPS and rebuild
ssh root@104.248.4.212 '
  cd /opt/onemind-web
  git pull origin main
  docker compose up -d --build
'

# Check it's running
ssh root@104.248.4.212 'docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"'
```

Data dir must be owned by UID 1000 (`node` user inside container). If fresh install:

```bash
ssh root@104.248.4.212 'chown -R 1000:1000 /opt/onemind-web/data'
```

---

## Related

- [UPDATES.md](UPDATES.md) — updating the OpenClaw binary
- [VPS-SERVICES.md](VPS-SERVICES.md) — all services on VPS
- [VPS-GATEWAY.md](VPS-GATEWAY.md) — gateway config and ports
- [OC VPS STATUS.md](OC%20VPS%20STATUS.md) — live VPS status
