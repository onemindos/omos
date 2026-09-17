# WATCHDOG PROTOCOL
**Domain:** SO 01 — Infrastructure  
**Owner:** Grid ⚙️  
**Status:** ACTIVE ✅  
**Last Updated:** 2026-04-28  

---

## Purpose

The Watchdog Protocol is the external self-healing system for One Mind OS infrastructure. It solves the dead-man's switch problem: if the OpenClaw gateway goes down, the agents inside it can't fix themselves. This protocol uses an external cloud agent (Letta Legacy-AI) triggered by two independent monitoring layers to detect outages, attempt automated recovery, and escalate to Zeus if recovery fails.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                  MONITORING LAYERS                   │
│                                                      │
│  Layer 1: VPS systemd timer (every 2 min, LOCAL)    │
│  └─ Survives: gateway crash, service failure         │
│                                                      │
│  Layer 2: Datadog Synthetic Monitor (every 2 min)   │
│  └─ Survives: full VPS outage, network failure       │
│  └─ Locations: AWS us-east-1 + us-west-2 (offsite)  │
└──────────────────┬──────────────────────────────────┘
                   │ GATEWAY_DOWN signal
                   ▼
┌─────────────────────────────────────────────────────┐
│           LETTA LEGACY-AI (Cloud Agent)              │
│  Agent ID: agent-da585f6f-189f-472a-974c-29b9578f9b55│
│  Platform: Letta Cloud (app.letta.com)               │
│  Model: claude-opus-4-5                              │
│  Memory: gateway_watchdog runbook (core memory)      │
│                                                      │
│  Tools:                                              │
│  - Bash (SSH into VPS)                               │
│  - send_discord_alert (Zeus #alerts)                 │
│  - telegram_send_alert (Telegram backup)             │
│  - call_zeus (phone escalation)                      │
│  - comms, fetch_webpage, BashOutput                  │
│  - run_code_with_tools, letta, todoist               │
└──────────────────┬──────────────────────────────────┘
                   │ SSH + notify
                   ▼
┌─────────────────────────────────────────────────────┐
│            VPS: onemindos-cloud-1                    │
│            root@64.227.8.23                          │
│            Service: openclaw-gateway                 │
│            Port: 18789                               │
└─────────────────────────────────────────────────────┘
```

---

## Components

### Layer 1 — VPS Local Watchdog

| Item | Value |
|---|---|
| Script | `/usr/local/bin/gateway-watchdog.sh` |
| Service | `/etc/systemd/system/gateway-watchdog.service` |
| Timer | `/etc/systemd/system/gateway-watchdog.timer` |
| Interval | Every 2 minutes |
| Log | `/var/log/gateway-watchdog.log` |
| State file | `/var/tmp/gateway-watchdog-state` |

**What it does:**
- Curls `http://localhost:18789/health` every 2 minutes
- If HTTP 200 → silent, do nothing
- If down → writes first-detection timestamp to state file, POSTs `GATEWAY_DOWN` message to Letta Legacy-AI API
- If still down on next cycle → POSTs `GATEWAY_STILL_DOWN` with elapsed time
- On recovery → clears state file, logs resolved

### Layer 2 — Datadog External Synthetic Monitor

| Item | Value |
|---|---|
| Monitor ID | `zwx-r5a-s9j` |
| Name | OpenClaw Gateway Health — onemindos-cloud |
| URL | `http://64.227.8.23:18789/health` |
| Interval | Every 120 seconds |
| Locations | aws:us-east-1, aws:us-west-2 |
| Assertion | HTTP status 200 |
| Webhook | `letta-gateway-watchdog` → Letta Legacy-AI API |
| Dashboard | https://us5.datadoghq.com |

**Why this layer exists:**
Layer 1 runs ON the VPS. If the VPS itself goes down (power, network, OOM), Layer 1 can't run. Datadog checks from AWS locations completely outside the VPS — covers full VPS outage scenarios.

### Letta Legacy-AI Agent

| Item | Value |
|---|---|
| Agent ID | `agent-da585f6f-189f-472a-974c-29b9578f9b55` |
| Platform | Letta Cloud |
| Model | claude-opus-4-5 |
| Memory block | `Legacy AI Persona` (contains full runbook) |
| API Key | 1Password → SO vault → "onemindos - open claw/letta code - api key" |

**Why Letta:**
- Cloud-hosted — survives VPS outage completely
- Persistent memory — always knows the runbook, VPS credentials, notification channels
- Webhook-activatable — both monitoring layers POST to it directly
- Has SSH execution capability via Bash tool
- Has direct lines to Zeus (Discord, Telegram, phone)

---

## Recovery Runbook (Legacy-AI Protocol)

Legacy-AI executes this on every `GATEWAY_DOWN` trigger:

### Step 1 — Confirm Outage
```bash
ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 root@64.227.8.23 \
  "curl -s http://localhost:18789/health"
```
- HTTP 200 → false alarm, notify Zeus "false alarm, gateway is up"
- SSH fails → VPS may be unreachable → skip to escalation

### Step 2 — Notify Zeus (Immediately)
- Discord: `send_discord_alert` → "⚠️ Gateway down detected. Attempting restart now..."
- Telegram: `telegram_send_alert` → backup notification
- Do NOT wait for confirmation — send and proceed

### Step 3 — Restart Attempt
```bash
ssh root@64.227.8.23 "systemctl --user restart openclaw-gateway"
# Wait 15 seconds
ssh root@64.227.8.23 "curl -s http://localhost:18789/health"
```
- If HTTP 200 → notify Zeus: "✅ Gateway restored after attempt X"
- If still down → proceed to Step 4

### Step 4 — Log Analysis
```bash
ssh root@64.227.8.23 "journalctl --user -u openclaw-gateway -n 100"
```
- **OOM killed** → check `free -h`, kill memory hogs if safe, retry
- **Port conflict** → find and kill process on port 18789, retry
- **Config error** → report exact error to Zeus, do NOT auto-fix config
- **Plugin crash** → identify plugin, attempt disable if possible

### Step 5 — Crash Loop Detection
- Crash loop = gateway restarts but crashes again within 60 seconds
- On crash loop: pull full logs, report diagnosis to Zeus every attempt
- Keep retrying every 2 minutes

### Step 6 — Escalation
- After 10 failed attempts → `call_zeus` (phone call)
- If VPS completely unreachable → `call_zeus` immediately
- Do not stop trying — continue attempts even after phone call

### Step 7 — Resolution Report
Always send final status regardless of outcome:
- Time gateway was first down
- Number of restart attempts
- Root cause (if identified)
- Resolution or escalation status

---

## Notification Chain

| Priority | Tool | When |
|---|---|---|
| 1st | `send_discord_alert` | Every attempt, status update |
| 2nd | `telegram_send_alert` | Same — backup channel |
| Emergency | `call_zeus` | 10 failures OR VPS unreachable |

---

## VPS Gateway Reference

| Item | Value |
|---|---|
| Host | `root@64.227.8.23` |
| Tailscale | `100.98.25.10` |
| Service name | `openclaw-gateway` |
| Service type | USER systemd service |
| Health endpoint | `http://localhost:18789/health` |
| Port | `18789` |
| Logs | `journalctl --user -u openclaw-gateway -n 100` |
| Config | `/root/.openclaw/openclaw.json` |

**Key commands:**
```bash
# Restart
systemctl --user restart openclaw-gateway

# Status
systemctl --user status openclaw-gateway

# Health check
curl -s http://localhost:18789/health

# Logs
journalctl --user -u openclaw-gateway -n 100

# Disk
df -h

# Memory
free -h
```

---

## Credentials

| Credential | Location |
|---|---|
| Letta API key | 1Password SO → "onemindos - open claw/letta code - api key" |
| Datadog API key | 1Password SO → "Datadog - Watch dog Api" |
| Datadog App key | 1Password SO → "Datadog" → app key field |
| VPS root password | 1Password SO → "VPS Root Password (onemindos-cloud)" |

---

## Ops Commands

```bash
# Check watchdog timer status
systemctl status gateway-watchdog.timer

# View live watchdog log
tail -f /var/log/gateway-watchdog.log

# Run watchdog manually (test)
systemctl start gateway-watchdog.service

# Check next trigger time
systemctl list-timers gateway-watchdog.timer

# Disable watchdog (maintenance)
systemctl stop gateway-watchdog.timer

# Re-enable watchdog
systemctl start gateway-watchdog.timer
```

---

## Coverage Matrix

| Failure Scenario | Layer 1 Detects | Layer 2 Detects | Auto-Recovery |
|---|---|---|---|
| Gateway process crash | ✅ <2 min | ✅ <2 min | ✅ via Legacy-AI |
| Gateway OOM killed | ✅ <2 min | ✅ <2 min | ✅ via Legacy-AI |
| Port conflict | ✅ <2 min | ✅ <2 min | ✅ via Legacy-AI |
| Plugin crash loop | ✅ <2 min | ✅ <2 min | ✅ via Legacy-AI |
| Config corruption | ✅ <2 min | ✅ <2 min | ⚠️ Reports to Zeus |
| Full VPS outage | ❌ (VPS is down) | ✅ <2 min | ⚠️ Zeus notified |
| Network outage | ❌ | ✅ <2 min | ⚠️ Zeus notified |

---

## Build History

| Date | Action |
|---|---|
| 2026-04-28 | Watchdog Protocol designed and deployed |
| 2026-04-28 | Legacy-AI agent configured with tools + runbook |
| 2026-04-28 | VPS systemd timer deployed (Layer 1) |
| 2026-04-28 | Datadog synthetic monitor created (Layer 2) |
| 2026-04-28 | Datadog webhook → Letta wired |
| 2026-04-28 | Full pipeline verified end-to-end |

---

---

## Context — Why This Was Built

This protocol was designed and deployed on 2026-04-28 in the same session where OpenClaw 4.25 and 4.26 were reviewed.

**OpenClaw 4.26** (released 2026-04-28) introduced Google Live browser Talk sessions, Cerebras as a bundled provider, memory embedding improvements for Ollama, and the `openclaw migrate` CLI tool. Full release notes: https://github.com/openclaw/openclaw/releases/tag/v2026.4.26

**OpenClaw 4.25** introduced a full TTS overhaul (Azure Speech, ElevenLabs v3, Volcengine, Inworld, Xiaomi), plugin cold registry (faster boot), full OpenTelemetry coverage, browser automation hardening, PWA + Web Push in the Control UI, and install/update hardening across all platforms. Full release notes: https://github.com/openclaw/openclaw/releases/tag/v2026.4.25

The Watchdog Protocol was conceived because OpenClaw 4.26 brought new voice/realtime features that increase the risk of gateway instability during upgrades. The protocol ensures zero manual intervention is required for recovery — the system heals itself.

**Note on Datadog coverage:** The full VPS outage blind spot (Layer 1 can't trigger if the VPS itself is down) is fully closed. Datadog synthetic monitor `zwx-r5a-s9j` checks from AWS us-east-1 and us-west-2 — completely external to the VPS. Both layers are live and operational.

---

*Built by Grid ⚙️ — One Mind OS Infrastructure*

---

## Recovery Webhook Server (v2 — Added 2026-04-28)

After initial deployment, three issues were identified and resolved:

### Issues Found
1. **Bash/SSH tool blocked in Letta cloud** — Letta's sandbox cannot make outbound SSH connections. Legacy-AI returned `GATEWAY RECOVERY BLOCKED — Bash tool unavailable`.
2. **Discord + Telegram tools outdated** — Both were calling `localhost:8100` (OpenClaw comms API). When the gateway is down, this fails too — circular dependency.
3. **Cloudflare WAF** was blocking POST requests from Letta's cloud IPs to the recovery subdomain.

### Fixes Applied

**1. Recovery Webhook Server**
- Script: `/usr/local/bin/recovery-webhook.py`
- Service: `/etc/systemd/system/recovery-webhook.service` (auto-start, always-on)
- Port: `18790` (open in DO firewall)
- URL: `http://recovery.onemindos.com:18790`
- Auth: `?token=om-recovery-2026-watchdog` (query param — CF strips headers)
- DNS: `recovery.onemindos.com` A record → 64.227.8.23 (CF proxy OFF)
- Credential: token stored in 1Password SO vault

Endpoints:
```
POST /restart-gateway  → restarts openclaw-gateway, waits 15s, returns health
POST /status           → gateway HTTP, service state, disk, memory
POST /logs             → last 50 lines of gateway logs
GET  /health           → liveness check
```

**2. Discord tool** — rewrote to call Discord Bot API directly (DM to Zeus)
- Bot token in tool source
- DM Channel: `1453585411925606615`
- Zero dependency on OpenClaw

**3. Telegram tool** — rewrote to call Telegram Bot API directly
- Bot token + chat ID in tool source
- Posts to Alerts thread (thread 38)
- Zero dependency on OpenClaw

### Verified End-to-End
From Letta cloud → `http://recovery.onemindos.com:18790/status?token=...` → returns VPS health ✅
From Letta cloud → `send_discord_alert` → Discord DM to Zeus ✅
Full pipeline test passed 2026-04-28.

---

## OpenClaw Management Layer (Added 2026-04-28)

Legacy-AI is not just a watchdog — it is the external OpenClaw management layer.
Zeus can send commands directly to Legacy-AI in Letta chat to control the gateway.

### Command Reference

| Command | What happens |
|---|---|
| `omos down` | Confirmed down — skip health check, go straight to restart protocol, notify Zeus |
| `omos up` | Check gateway health, report status |
| `omos restart` | Force restart gateway even if healthy |
| `omos update` | npm install -g openclaw@latest → restart → verify → report version before/after |
| `omos status` | Full snapshot — gateway HTTP, OpenClaw version, disk, memory |
| `omos logs` | Last 50 lines of gateway logs |
| `omos diagnose` | Full diagnostic — logs, disk, memory, error identification |

### Agent Optimisation (2026-04-28)

Memory blocks removed (irrelevant to infra role):
- `business_goals` (3D printer — not infra)
- `Family Information` (not infra)
- `morning brief` (3-line stub)
- `onemind_documentation_system` (referenced unavailable tools)
- `gateway_watchdog` (duplicate — superseded by Legacy AI Persona block)

Remaining blocks:
- `Legacy AI Persona` — full runbook + command routing + recovery webhook (single source of truth)
- `pillar_routing` — hive mind orchestration context

### Recovery Webhook v3

Added `/update-gateway` endpoint:
- Runs `npm install -g openclaw@latest`
- Restarts gateway
- Returns version before/after + health status

### Maintenance Schedule

Weekly Todoist recurring task: every Monday 9am ET
Task ID: `6gVV7rm3pv8c9gRW`
Checklist: `omos status` → `omos logs` → `omos update` → `omos diagnose`

---

## Runbook Expansion (2026-04-28)

Full per-command instructions added to Legacy AI Persona memory block.
False alert fix: Datadog monitor updated — now requires BOTH AWS locations to fail + 2 retries + 120s sustained before firing.
Update fix: Legacy-AI corrected — `omos update` POSTs to `/update-gateway` directly, never asks Zeus to deploy manually.

### Per-Command Detail

**omos status** — POST /status → report gateway_http, version, disk, memory. If not 200 → trigger down protocol.
**omos down** — Skip health check. Notify Zeus, POST /restart-gateway, retry up to 10x, escalate via call_zeus if all fail.
**omos restart** — Force restart even if healthy. Notify before + after. If still_down → escalate.
**omos update** — POST /update-gateway (handles npm install + restart internally). Report version before/after. If fails → down protocol.
**omos logs** — POST /logs → return + flag ERROR/WARN/OOM/port conflict lines.
**omos diagnose** — POST /status + /logs → flag disk >80%, memory >85%, errors. Report severity per item.
**omos up** — POST /status → confirm 200. If not → down protocol immediately.
**omos commands** — No tool call. Print command reference from memory.

### False Alert Fix
Datadog monitor `zwx-r5a-s9j` settings updated:
- min_failure_duration: 120s (was 0)
- min_location_failed: 2 (both AWS locations must fail)
- retry count: 2 (with 30s interval between retries)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/operations (Ops Procedures)/security (Security Ops)/security|security]]
