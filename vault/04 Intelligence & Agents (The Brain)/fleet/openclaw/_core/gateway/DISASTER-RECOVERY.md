---
title: "OpenClaw Disaster Recovery Runbook"
type: runbook
status: active
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [type/runbook, domain/SO, domain/02-infra]
ai_summary: "Step-by-step disaster recovery for all OpenClaw VMs. Three recovery paths: config rollback, full VM restore from Lightsail snapshot, and deploy-from-zero with bootstrap.sh."
---

# OpenClaw Disaster Recovery Runbook

> Platform: AWS Lightsail | Tailscale mesh: `tail717ba0.ts.net`
> Last updated: 2026-05-10

---

## Recovery Paths (choose by scenario)

| Scenario | Path | Time |
|----------|------|------|
| Bad config change | Git rollback | ~2 min |
| VM corrupted / data lost | Lightsail snapshot restore | ~10 min |
| VM deleted / region failure | Bootstrap from zero | ~5 min |
| Gateway crashed, won't restart | Process restart | ~1 min |

---

## Path 1 — Config Rollback (bad openclaw.json)

Every config change auto-commits to GitHub within ~15 seconds. Roll back to any point in history.

```bash
# On your Mac — find the last good commit
gh api "repos/OneMind-OS/openclaw-core-workspaces/commits" \
  --jq '.[].commit | .message + " (" + .committer.date + ")"' | head -10

# Get the openclaw.json from that commit
gh api "repos/OneMind-OS/openclaw-core-workspaces/contents/openclaw.json?ref=<COMMIT_SHA>" \
  --jq '.content' | base64 -d > /tmp/openclaw-good.json

# Push it to the VM
scp /tmp/openclaw-good.json root@omos-core.tail717ba0.ts.net:/home/ubuntu/.openclaw/openclaw.json

# Restart gateway on VM
ssh root@omos-core.tail717ba0.ts.net "systemctl restart openclaw || pkill -f 'openclaw gateway' && su - ubuntu -c 'openclaw gateway &'"
```

Alternatively, revert directly on the VM:
```bash
ssh root@omos-core.tail717ba0.ts.net
su - ubuntu
cd ~/openclaw-workspaces
git log --oneline -10                          # find last good commit
git show <HASH>:openclaw.json > ~/.openclaw/openclaw.json
# Gateway will auto-reload or restart it manually
```

---

## Path 2 — Lightsail Snapshot Restore

Use when the VM is corrupted, OS-level issue, or disk problem.

**⚠️ Requires: Automatic snapshots enabled in Lightsail console (enable now if not done)**

```
Lightsail console → Instance → Snapshots → Enable automatic snapshots
Set to: daily at 02:00 UTC, keep 7 snapshots
Do this for: onemindos-core, omos-sentinel, omos-titan
```

### Restore steps:
1. Lightsail console → Snapshots → select most recent snapshot
2. **Create new instance from snapshot** (don't overwrite — create fresh)
3. New instance gets a new public IP — update Tailscale static IP if needed
4. Verify Tailscale reconnects automatically (it should — tailscaled is a service)
5. SSH via Tailscale: `ssh root@omos-core.tail717ba0.ts.net`
6. Verify services: `systemctl status openclaw-config-sync onemind-watch`
7. Test gateway: `curl http://localhost:18789/v1/models`

**Time: ~10 minutes**

---

## Path 3 — Deploy from Zero (bootstrap.sh)

Use when: VM is gone, snapshot not available, or moving to a new region/provider.

```bash
# 1. Create new Lightsail instance
#    - OS: Ubuntu 24.04 LTS
#    - Plan: 4GB+ RAM for Core, 1GB for Sentinel
#    - Enable Tailscale after boot

# 2. SSH into new instance
ssh ubuntu@<NEW_IP>

# 3. Export your 1Password service account token
export OP_SERVICE_ACCOUNT_TOKEN="<from 1Password SO vault → OP Service Account>"

# 4. Run bootstrap (auto-detects Core vs Sentinel from hostname)
curl -sL https://raw.githubusercontent.com/OneMind-OS/openclaw-core-workspaces/main/bootstrap.sh | bash

# For Sentinel:
curl -sL https://raw.githubusercontent.com/OneMind-OS/openclaw-sentinel-workspaces/main/bootstrap.sh | bash
```

**What bootstrap.sh does:**
1. Creates `ubuntu` user (already exists on Lightsail)
2. Installs: node, op, gh, openclaw, inotify-tools, rsync
3. Clones workspace repo → deploys `openclaw.json` + workspace files
4. Copies scripts: `openclaw-reload.sh`, `openclaw-config-sync.sh`
5. Installs + enables systemd services: `openclaw`, `openclaw-reload`, `openclaw-config-sync`
6. Starts all services

**After bootstrap:**
- Re-add instance to Tailscale network
- Update Tailscale hostname if changed
- Verify: `curl http://localhost:18789/v1/models`

**Time: ~5 minutes**

---

## Path 4 — Gateway Crashed (quick restart)

```bash
ssh root@omos-core.tail717ba0.ts.net

# Check what's running
ps aux | grep openclaw | grep -v grep

# If the node process is dead, restart it
su - ubuntu -c 'openclaw gateway --port 18789 &'

# Or if managed by systemd
systemctl restart openclaw

# Verify
sleep 3 && curl http://localhost:18789/v1/models
```

---

## Verification Checklist (post-recovery)

```bash
# 1. Gateway responding
curl http://localhost:18789/v1/models | head -20

# 2. Config sync running
systemctl status openclaw-config-sync

# 3. Vault sync running (Core only)
systemctl status onemind-watch onemind-pull.timer

# 4. OpenClaw process alive
ps aux | grep 'openclaw gateway' | grep -v grep

# 5. Tailscale connected
tailscale status

# 6. Test a config change syncs to GitHub
touch ~/.openclaw/openclaw.json
sleep 15
cd ~/openclaw-workspaces && git log --oneline -3
# Should show a new "config: auto-sync" commit
```

---

## Lightsail Snapshot Schedule (action required)

**Enable automatic daily snapshots for all 3 VMs in the Lightsail console:**

```
AWS Lightsail Console → Instances → onemindos-core
  → Snapshots tab → Enable automatic snapshots
  → Time: 02:00 UTC | Retention: 7 days

Repeat for: omos-sentinel, omos-titan
```

Cost: ~$0.05/GB/month (Core is ~20GB = ~$1/month per VM)

---

## Key Info (keep this accessible)

| Item | Value |
|------|-------|
| Core Tailscale URL | `omos-core.tail717ba0.ts.net` |
| Sentinel Tailscale URL | `omos-sentinel.tail717ba0.ts.net` |
| Core workspace repo | `OneMind-OS/openclaw-core-workspaces` |
| Sentinel workspace repo | `OneMind-OS/openclaw-sentinel-workspaces` |
| 1Password vault ID | `aaeyris5z44pdqspk277u2oppi` |
| Gateway token location | `op://aaeyris5z44pdqspk277u2oppi/OpenClaw Gateway (onemindos-cloud)/token` |
| OP Service Account | `op://aaeyris5z44pdqspk277u2oppi/OP Service Account/token` |

---

## Related

- [[VPS-SERVICES]] — full service inventory and monitoring commands
- [[OPENCLAW-LIVE-STATE]] — current production state
- [[OPENCLAW-CONFIG-REFERENCE]] — config schema and secret locations
