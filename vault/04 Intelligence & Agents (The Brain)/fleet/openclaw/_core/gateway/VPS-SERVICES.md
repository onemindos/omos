---
title: "VPS Services & Sync Architecture"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [type/reference, domain/SO, domain/02-infra]
ai_summary: "Current systemd services, sync architecture, and runtime state for all OpenClaw VMs (Core + Sentinel). Last verified 2026-05-10."
---

# VPS Services & Sync Architecture

> Last verified: 2026-05-10
> Platform: AWS Lightsail (NOT DigitalOcean)
> OS: Ubuntu 24.04 LTS — default user is `ubuntu` (not `onemind`)

---

## VM Inventory

| VM | Hostname | Public IP | Tailscale URL | Role |
|----|----------|-----------|--------------|------|
| **onemindos-core** | omos-core | 64.227.8.23 | `omos-core.tail717ba0.ts.net` | Primary — 11 agents |
| **omos-sentinel** | ip-172-26-6-151 | (Lightsail private) | `omos-sentinel.tail717ba0.ts.net` | Sentinel — monitoring |
| **omos-titan** | omos-titan | — | — | Workhorse — services TBD |

> SSH access is via Tailscale only. Direct IP SSH is blocked on Core (port 22 refused).
> Always connect via: `ssh root@omos-core.tail717ba0.ts.net`

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                  omos-core (PRIMARY VM)                       │
│  /home/ubuntu/.openclaw/openclaw.json  ← live config         │
│  /home/ubuntu/.openclaw/workspace*/    ← 11 agent workspaces │
│                                                              │
│  [openclaw gateway]        ← port 18789 (loopback + TS)     │
│  [watch.sh]                ← OpenClaw built-in workspace sync│
│  [onemind-watch.service]   ← vault sync → GitHub            │
│  [onemind-pull.timer]      ← GitHub → vault (every 10s)     │
│  [openclaw-config-sync]    ← config changes → GitHub        │
│  [openclaw-rotate-token]   ← token rotation (25min timer)   │
└───────────────┬──────────────────────────────────────────────┘
                │ git push (openclaw-config-sync, 15s debounce)
                ▼
  GitHub: OneMind-OS/openclaw-core-workspaces (PRIVATE)
                │
                ▼ clone → bootstrap.sh → new VM live in 5min

┌──────────────────────────────────────────────────────────────┐
│                  omos-sentinel (SENTINEL VM)                  │
│  /home/ubuntu/.openclaw/openclaw.json                        │
│  /home/ubuntu/.openclaw/workspace/                           │
│                                                              │
│  [openclaw gateway]        ← port 18789                     │
│  [watch.sh]                ← built-in workspace sync        │
│  [openclaw-config-sync]    ← config changes → GitHub        │
└───────────────┬──────────────────────────────────────────────┘
                │ git push
                ▼
  GitHub: OneMind-OS/openclaw-sentinel-workspaces (PRIVATE)
```

---

## Core VM — Active Services

### `openclaw` (gateway)
OpenClaw AI gateway — runs all 11 agents.

| Field | Value |
|-------|-------|
| Process | `/usr/bin/node /usr/lib/node_modules/openclaw/dist/index.js gateway --port 18789` |
| User | `ubuntu` |
| Port | `18789` (loopback only — Tailscale Serve proxies HTTPS) |
| Config | `/home/ubuntu/.openclaw/openclaw.json` |
| Version | v2026.5.2 |

```bash
ps aux | grep openclaw
journalctl -u openclaw -f
```

---

### `onemind-watch.service` — Vault → GitHub sync
Watches vault files and auto-commits changes to `OneMind-OS/One-Mind`.

| Field | Value |
|-------|-------|
| Script | `/usr/local/bin/onemind-watch.sh` |
| User | `ubuntu` |
| Status | `active (running)` |

```bash
systemctl status onemind-watch
journalctl -u onemind-watch -f
```

---

### `onemind-pull.timer` — GitHub → Vault (every 10s)
Pulls vault updates from GitHub into the local vault path.

| Field | Value |
|-------|-------|
| Interval | ~10 seconds |
| Status | `active` |

```bash
systemctl status onemind-pull.timer
journalctl -u onemind-pull -f
```

---

### `openclaw-config-sync.service` ← added 2026-05-10
Watches `openclaw.json`, `agents/`, `identity/` for changes and auto-commits to `openclaw-core-workspaces` GitHub repo.

| Field | Value |
|-------|-------|
| Script | `/home/ubuntu/openclaw-config-sync.sh` |
| User | `ubuntu` |
| Repo dir | `/home/ubuntu/openclaw-workspaces/` |
| Watches | `openclaw.json`, `agents/`, `identity/` |
| Debounce | 3s + 10s cooldown |
| Status | `active (running)` ✅ |

```bash
systemctl status openclaw-config-sync
journalctl -u openclaw-config-sync -f
# Live commits: github.com/OneMind-OS/openclaw-core-workspaces
```

---

### `openclaw-rotate-token.timer`
Rotates the gateway token every 25 minutes.

| Field | Value |
|-------|-------|
| Interval | 25 minutes |
| Status | `active` |

---

### `watch.sh` (OpenClaw built-in)
OpenClaw's own workspace watcher — watches all workspace dirs and hot-reloads agent context.

| Field | Value |
|-------|-------|
| Process | `/bin/bash /home/ubuntu/.openclaw/watch.sh` |
| Watches | All `workspace-*/` dirs under `.openclaw/` |
| Managed by | OpenClaw internals (not systemd) |

---

## Sentinel VM — Active Services

### `openclaw` (gateway)
Single-agent Sentinel gateway.

| Field | Value |
|-------|-------|
| Process | `/usr/bin/node /usr/lib/node_modules/openclaw/dist/index.js gateway --port 18789` |
| User | `ubuntu` |
| Agents | 1 (Sentinel/main) |

---

### `openclaw-config-sync.service` ← added 2026-05-10
Same pattern as Core — watches config and auto-commits to `openclaw-sentinel-workspaces`.

| Field | Value |
|-------|-------|
| Script | `/home/ubuntu/openclaw-config-sync.sh` |
| Repo dir | `/home/ubuntu/openclaw-workspaces/` |
| Status | `active (running)` ✅ |

```bash
systemctl status openclaw-config-sync
journalctl -u openclaw-config-sync -f
```

---

### `watch.sh` (OpenClaw built-in)
Watches `workspace/` for Sentinel agent context changes.

---

## Config Sync — Verified Working (2026-05-10)

```
Core (OneMind-OS/openclaw-core-workspaces):
  1db8ba4 config: auto-sync 2026-05-10T08:38:29Z  ✅
  e4850b5 config: auto-sync 2026-05-10T08:38:13Z
  7c1b00f config: auto-sync 2026-05-10T08:37:58Z

Sentinel (OneMind-OS/openclaw-sentinel-workspaces):
  ee6fb34 config: auto-sync 2026-05-10T08:44:17Z  ✅
```

Any change to `openclaw.json`, `agents/`, or `identity/` commits within ~15 seconds.

---

## Monitoring Commands

```bash
# SSH into VMs (Tailscale only)
ssh root@omos-core.tail717ba0.ts.net
ssh root@omos-sentinel.tail717ba0.ts.net

# Check all services (Core)
systemctl status onemind-watch onemind-pull.timer openclaw-config-sync openclaw-rotate-token.timer

# Live config sync log
journalctl -u openclaw-config-sync -f

# OpenClaw gateway process
ps aux | grep openclaw | grep -v grep

# Check workspace repo commit history
su - ubuntu -c 'cd ~/openclaw-workspaces && git log --oneline -5'

# Test sync (touch config → should commit within 15s)
su - ubuntu -c 'touch ~/.openclaw/openclaw.json && sleep 15 && cd ~/openclaw-workspaces && git log --oneline -3'
```

---

## Troubleshooting

### Full Sync Chain

```
Mac (Obsidian / Claude Code)
    ↓ 1. git commit    ← Obsidian Git auto-commits (every few min)
    ↓ 2. git push      ← MUST HAPPEN or VPS never sees it
GitHub (OneMind-OS/One-Mind) ← SOURCE OF TRUTH
    ↓ 3. onemind-pull.timer (every 10s, automatic)
/home/ubuntu/onemind/ on VPS
    ↓ 4. onemind-watch.service detects changes → commits → pushes
GitHub ← updated
    ↓ 5. Mac: git pull to receive VPS-side changes
```

### Common Issues

| Problem | Check | Fix |
|---------|-------|-----|
| **Mac changes not appearing on VPS** | `git status` on Mac — look for "ahead of origin" | `git push origin main` — Obsidian committed but didn't push |
| **Mac push succeeded but VPS still old** | `ssh omos-core "systemctl status onemind-pull.timer"` | Wait 10s, or force: `ssh omos-core "systemctl start onemind-pull.service"` |
| **VPS changes not appearing on Mac** | `git fetch && git status` on Mac | `git pull` — Mac doesn't auto-pull |
| **Merge conflict on VPS pull** | `ssh omos-core "journalctl -u onemind-pull --since '5 min ago'"` | SSH in, resolve manually: `cd /home/ubuntu/onemind && git pull --rebase` |
| Config not committing | `systemctl status openclaw-config-sync` | `systemctl restart openclaw-config-sync` |
| Workspace repo diverged | `git -C ~/.openclaw status` shows "diverged" | `git stash push -u && git pull --rebase origin main && git stash pop && git push` |
| SSH refused on Core IP | Port 22 blocked | Use Tailscale: `ssh omos-core.tail717ba0.ts.net` |
| known_hosts mismatch | Host keys rotated | `ssh-keygen -R <hostname>` then reconnect |
| OpenClaw not responding | `ps aux \| grep openclaw` | Restart process or see [[DISASTER-RECOVERY]] |

### Quick Diagnostic (run from Mac)

```bash
# 1. Am I ahead? (committed but not pushed)
git -C ~/Documents/onemind status | grep "ahead"

# 2. Force push if needed
git -C ~/Documents/onemind push origin main

# 3. Verify VPS received it (wait 10s first)
ssh omos-core "git -C /home/ubuntu/onemind log --oneline -1"

# 4. Force VPS pull if timer seems stuck
ssh omos-core "systemctl start onemind-pull.service"

# 5. Check pull timer health
ssh omos-core "systemctl status onemind-pull.timer"
```

### Obsidian Git Plugin — Ensure Auto-Push

If you keep hitting "committed but not pushed", check Obsidian Git settings:
- **Auto push after commit:** ON
- **Push interval (minutes):** 1-2 (or enable "Push on backup")
- Without auto-push, commits stay local until manually pushed

---

## Related

- [[OPENCLAW-LIVE-STATE]] — current gateway state, agents, skills
- [[OPENCLAW-CONFIG-REFERENCE]] — config schema and 1Password secret locations
- [[DISASTER-RECOVERY]] — full disaster recovery runbook
- Workspace repos: `OneMind-OS/openclaw-core-workspaces` · `OneMind-OS/openclaw-sentinel-workspaces`
