---
name: sentinel
description: "OneMind OS Sentinel — unified runbook for omos-sentinel and omos-core OpenClaw nodes. Covers infrastructure topology, agent profiles, model configs, channels, plugins, workspace identities, recovery procedures, cron, and operational commands. WHEN: sentinel status, omos status, omos restart, omos down, omos update, omos logs, omos diagnose, omos up, gateway config, openclaw config, sentinel models, sentinel channels, recovery webhook, watchdog, agent config, model routing, bedrock config, discord bot config, telegram bot config, node topology, openclaw troubleshooting."
license: MIT
metadata:
  author: Zeus / OneMind OS
  version: "1.0.0"
  updated: "2026-05-03"
---

# SENTINEL — OneMind OS Unified Runbook

> Authoritative operational reference for the entire OneMind OS agent infrastructure.
> Two OpenClaw gateway nodes, multi-provider model routing, Discord + Telegram channels, recovery automation.

---

## 1. NODE TOPOLOGY

| Node | SSH | Tailscale IP | Tailscale FQDN | Role | Service Type | Hardware |
|------|-----|-------------|----------------|------|-------------|----------|
| **omos-core** | `root@100.123.126.34` | 100.123.126.34 | omos-core.tail717ba0.ts.net | Primary gateway, multi-agent, 80+ models | systemd SYSTEM (`/etc/systemd/system/`) | 8 vCPU, 30GB RAM, 619GB disk |
| **omos-sentinel** | `ubuntu@100.66.235.38` | 100.66.235.38 | omos-sentinel.tail717ba0.ts.net | Secondary gateway, single agent, 24 models | systemd USER (`~/.config/systemd/user/`) | 2 vCPU, 3.7GB RAM, 77GB disk (Lightsail) |
| **omos-titan** | `ubuntu@100.103.197.23` | 100.103.197.23 | omos-titan.tail717ba0.ts.net | Postiz social media, Temporal workflows | — | — |

**Software:** OpenClaw v2026.5.2 on both nodes, Node.js 22.22.2, npm global install at `/usr/lib/node_modules/openclaw/`

---

## 2. OMOS-SENTINEL — DETAILED PROFILE

### 2.1 Agent

| Field | Value |
|-------|-------|
| Agent ID | `main` |
| Agent Name | **Sentinel** |
| Default | Yes (only agent) |
| Identity | OneMind Actual — military-grade C2 operator |
| Primary Model | `github-copilot/claude-sonnet-4.6` [GH] |
| Subagent Model | `github-copilot/claude-haiku-4.5` [GH] |
| Fallback #1 | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` [bedrock] |
| Fallback #2 | `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0` [bedrock] |
| Sandbox | Off |

### 2.2 Models (24 total)

**GitHub Copilot (21 models) — tagged [GH]:**

| Alias | Model ID | Role |
|-------|----------|------|
| Sonnet 4.6 [GH] | github-copilot/claude-sonnet-4.6 | **Primary default** |
| Haiku 4.5 [GH] | github-copilot/claude-haiku-4.5 | **Subagent default** |
| Opus 4.7 [GH] | github-copilot/claude-opus-4.7 | On-demand |
| Opus 4.6 [GH] | github-copilot/claude-opus-4.6 | On-demand |
| Opus 4.5 [GH] | github-copilot/claude-opus-4.5 | On-demand |
| Sonnet 4.5 [GH] | github-copilot/claude-sonnet-4.5 | On-demand |
| Sonnet 4 [GH] | github-copilot/claude-sonnet-4 | On-demand |
| GPT-5.4 [GH] | github-copilot/gpt-5.4 | On-demand |
| GPT-5.4 Mini [GH] | github-copilot/gpt-5.4-mini | On-demand |
| GPT-5.4 Nano [GH] | github-copilot/gpt-5.4-nano | On-demand |
| GPT-5.3 Codex [GH] | github-copilot/gpt-5.3-codex | On-demand |
| GPT-5.2 [GH] | github-copilot/gpt-5.2 | On-demand |
| GPT-5.2 Codex [GH] | github-copilot/gpt-5.2-codex | On-demand |
| GPT-5 Mini [GH] | github-copilot/gpt-5-mini | On-demand |
| GPT-4.1 [GH] | github-copilot/gpt-4.1 | On-demand |
| Gemini 3.1 Pro [GH] | github-copilot/gemini-3.1-pro | On-demand |
| Gemini 3 Flash [GH] | github-copilot/gemini-3-flash | On-demand |
| Gemini 2.5 Pro [GH] | github-copilot/gemini-2.5-pro | On-demand |
| Grok Code Fast [GH] | github-copilot/grok-code-fast-1 | On-demand |
| Goldeneye [GH] | github-copilot/goldeneye | On-demand |
| Raptor Mini [GH] | github-copilot/raptor-mini | On-demand |

**AWS Bedrock (3 models) — tagged [bedrock]:**

| Alias | Model ID | Role |
|-------|----------|------|
| Sonnet 4.6 [bedrock] | amazon-bedrock/us.anthropic.claude-sonnet-4-6 | **Fallback #1** |
| Haiku 4.5 [bedrock] | amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0 | **Fallback #2** |
| Opus 4.6 [bedrock] | amazon-bedrock/us.anthropic.claude-opus-4-6 | On-demand |

### 2.3 Channels

**Telegram:**
- Bot: `@omos_sentinel_bot`
- Mode: polling, partial streaming
- DM Policy: open
- Group: `-1003905828404` → topic `195` → agentId `main`
- Commands: native auto

**Discord:**
- Bot: `@Sentinel`
- Guild: `1453411552181883027` (OnemindOS)
- Channel: `1482355265092190319` (sentinel-chat), requireMention=true
- Group Policy: allowlist
- Streaming: off
- Thread Bindings: enabled, idle 24h

**Bindings:**
- `agentId=main` ← `channel=telegram`
- `agentId=main` ← `channel=discord`

### 2.4 Plugins (12 enabled)

| Plugin | Purpose |
|--------|---------|
| `amazon-bedrock` | AWS Bedrock model provider |
| `github-copilot` | GitHub Copilot model provider |
| `browser` | Web browsing/scraping |
| `device-pair` | Device pairing (Mac ↔ Sentinel) |
| `discord` | Discord channel (npm-installed: `@openclaw/discord`) |
| `file-transfer` | File ops between paired nodes |
| `memory-core` | Core memory system |
| `openshell` | Sandbox shell execution |
| `skill-workshop` | Skill capture/management |
| `telegram` | Telegram channel |
| `thread-ownership` | Prevents multi-agent thread conflicts |
| `webhooks` | Inbound webhook → TaskFlow bridge |

`plugins.allow`: discord, telegram, github-copilot, amazon-bedrock, browser, device-pair, file-transfer, memory-core, openshell, skill-workshop, thread-ownership, webhooks

### 2.5 Gateway Config

| Setting | Value |
|---------|-------|
| Port | 18789 |
| Bind | loopback |
| Mode | local |
| Auth Mode | token |
| Gateway Token | `dtef9yHLbdpjw9R64xlnXMMg7tRGm8Ib` |
| Tailscale | serve mode |
| Allowed Origins | localhost:18789, 127.0.0.1:18789, omos-sentinel.tail717ba0.ts.net |
| Trusted Proxies | 127.0.0.1 |

### 2.6 Workspace Identity Files

Located at `~/.openclaw/workspace/`:

| File | Purpose |
|------|---------|
| `IDENTITY.md` | Name=OneMind, Instance=omos-sentinel, Callsign=OneMind Actual, military operator persona |
| `SOUL.md` | Core directives — military brevity, standing orders, ROE, voice mode protocol |
| `AGENTS.md` | Full operating manual — session startup, memory protocol, communication, multi-agent awareness |
| `USER.md` | Zeus profile — Eastern TZ, prefers direct comms, mv over cp, never delete |
| `MEMORY.md` | Long-term memory — operator notes, infra state, lessons learned |
| `HEARTBEAT.md` | Periodic ops checklist — infra ping, memory maintenance, Todoist, workspace hygiene |
| `TOOLS.md` | Infrastructure reference — servers, models, cloud services, key paths |

**Sub-Workspaces:**

| Workspace | Role |
|-----------|------|
| `workspace-coder/` | Coding specialist — PRs, refactoring, Claude Code tasks |
| `workspace-comms/` | Read-only comms agent — Telegram/Discord message handling |
| `workspace-media/` | Media processing — transcription, PDF, video frames |
| `workspace-ops/` | Infrastructure ops — heartbeats, GitHub workflows, health checks |

### 2.7 Systemd Service

**Unit file:** `~/.config/systemd/user/openclaw-gateway.service`
```
ExecStart=/usr/bin/node /usr/lib/node_modules/openclaw/dist/index.js gateway --port 18789
Restart=always
RestartSec=5
EnvironmentFile=-/home/ubuntu/.openclaw/gateway.systemd.env
```

**Bedrock drop-in:** `~/.config/systemd/user/openclaw-gateway.service.d/bedrock.conf`
```
[Service]
Environment=AWS_ACCESS_KEY_ID=AKIA6CU7M5XTA7IAMIP5
Environment=AWS_SECRET_ACCESS_KEY=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
Environment=AWS_REGION=us-east-1
Environment=AWS_PROFILE=
```

**Key rule:** `AWS_PROFILE=` (empty) MUST be set to override the Lightsail instance role. Without this, Bedrock auth fails with STS AssumeRole errors.

### 2.8 Config Auto-Restore Protection

OpenClaw compares `openclaw.json` to `.bak` on startup. If config is smaller than backup, it auto-restores from `.bak`. **Always write BOTH files with identical content** when editing config.

```bash
# Safe config edit pattern:
python3 -c '
import json
cfg = json.load(open("/home/ubuntu/.openclaw/openclaw.json"))
# ... make changes ...
for path in ["/home/ubuntu/.openclaw/openclaw.json", "/home/ubuntu/.openclaw/openclaw.json.bak"]:
    with open(path, "w") as f:
        json.dump(cfg, f, indent=2)
        f.write("\n")
'
systemctl --user restart openclaw-gateway.service
```

### 2.9 Cron / Scheduled Tasks

**(none configured on sentinel)**

---

## 3. OMOS-CORE — DETAILED PROFILE

### 3.1 Agents (9 total)

| Agent ID | Name | Default | Purpose |
|----------|------|---------|---------|
| `legacy` | — | **Yes** | Primary operator agent (Zeus's main interface) |
| `spartan` | — | No | — |
| `oracle` | — | No | — |
| `eden` | — | No | — |
| `haven` | — | No | — |
| `empire` | — | No | — |
| `guardian` | — | No | — |
| `grid` | — | No | — |
| `trinity` | — | No | — |

**Model Routing:**
- Primary: `github-copilot/claude-sonnet-4.6` [gh]
- Fallback #1: `amazon-bedrock/us.anthropic.claude-sonnet-4-6` [aws]
- Fallback #2: `amazon-bedrock/us.anthropic.claude-opus-4-6-v1`
- Fallback #3: `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0` [aws]
- Subagent: `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0`

### 3.2 Models (80+ total)

**Provider breakdown:**

| Provider | Tag | Count | Examples |
|----------|-----|-------|----------|
| github-copilot | [gh] | ~30 | Sonnet 4.6, Opus 4.7, GPT-5.5, Gemini 3.1 Pro, Grok Code Fast |
| amazon-bedrock | [aws] | ~15 | Sonnet 4.6, Opus 4.6, Haiku 4.5, Nova Pro/Lite/Micro/Premier, Llama 4, DeepSeek R1, Mistral Large 3 |
| microsoft-foundry | [az] | ~10 | GPT-5, GPT-4.1, o4-mini |
| google | [gg] | ~9 | Gemini 3.1 Pro, 3 Flash, 2.5 Pro/Flash |
| litellm | [do] | ~18 | Qwen3.5, DeepSeek V3.2, Llama 4 Maverick, Nemotron, GLM 5, Kimi K2.5 |

### 3.3 Channels

**Discord:**
- Bot: `@OneMind OS`
- Guild: `1453411552181883027` (OnemindOS) — 14 channels enabled, requireMention=false
- Group Policy: allowlist
- Thread Bindings: enabled

**Telegram:**
- Bot: `@onemindos_bot`
- DM Policy: pairing
- Group: `-1003905828404` — 9 topics (1, 94–101), all bound to `legacy` agent
- Streaming: partial with tool progress

### 3.4 Plugins (enabled)

github-copilot, amazon-bedrock, browser, telegram, discord

(device-pair, openai, microsoft-foundry, google, litellm, talk-voice, phone-control, file-transfer, memory-core, bonjour all disabled)

### 3.5 Cron

```json
{
  "enabled": true,
  "maxConcurrentRuns": 2,
  "retry": { "maxAttempts": 3, "retryOn": ["rate_limit","overloaded","network","server_error"] },
  "sessionRetention": "48h"
}
```

### 3.6 Advanced Config

- **Memory Search:** Hybrid (vector 0.7 + text 0.3), MMR lambda 0.7, temporal decay 30d half-life, extra paths: `/root/onemind-codex`. Currently **disabled**.
- **Sandbox:** mode=non-main, scope=session, Docker image=`openclaw-sandbox:bookworm-slim`, workspace rw
- **Compaction:** safeguard mode, memory flush enabled
- **Subagents:** max 4 children/agent, max 6 concurrent, max depth 2

### 3.7 Bedrock Drop-in

`/etc/systemd/system/openclaw-gateway.service.d/bedrock.conf` — same IAM keys as sentinel, plus `AWS_DEFAULT_REGION=us-east-1`

Additional drop-in configs: `azure.conf`, `deepgram.conf`, `digitalocean.conf`, `fal.conf`, `google.conf`, `groq.conf`, `vertex.conf`

---

## 4. CREDENTIALS REFERENCE

| Item | Value |
|------|-------|
| **Sentinel TG bot** | `@omos_sentinel_bot` / token: `8684619971:AAH7q0DhymiQgno9mcP3jIrBAxXYXWG8kiI` |
| **Sentinel Discord bot** | `@Sentinel` / token: `MTQ4MjMxOTg2NzUyNzgyNzU0OA.GErDWc.neOMR7WoIA21mqxCwRt83DPUvkBzZumY2fQ8uA` |
| **Core TG bot** | `@onemindos_bot` / token: `8632096831:AAENM2qBvrp3CS069yrNhIcOEMh096vSbI8` |
| **Discord Guild** | `1453411552181883027` (OnemindOS) |
| **Sentinel Discord channel** | `1482355265092190319` (sentinel-chat) |
| **Sentinel TG group** | `-1003905828404` topic `195` |
| **Gateway token (sentinel)** | `dtef9yHLbdpjw9R64xlnXMMg7tRGm8Ib` |
| **AWS Bedrock IAM** | `AKIA6CU7M5XTA7IAMIP5` / `8f3M6VB8oDwyJdS/eeB/V+9ABYIJh1/W9OPUnm2A` / us-east-1 |

---

## 5. RECOVERY WEBHOOK (Sentinel VPS)

**Base URL:** `http://recovery.onemindos.com:18790`
**Auth:** `?token=om-recovery-2026-watchdog` (ALWAYS query param — never Authorization header)

### Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/status` | Gateway HTTP status, version, disk, memory |
| POST | `/restart-gateway` | Restart openclaw-gateway, wait 15s, return health |
| POST | `/logs` | Last 50 lines gateway logs |
| POST | `/update-gateway` | `npm install -g openclaw@latest`, restart, health check |
| GET | `/health` | Liveness check |

**Do NOT invent endpoint names. Only use the exact paths above.**

### Recovery Code Snippets

```python
# Status
import urllib.request, json
req = urllib.request.Request("http://recovery.onemindos.com:18790/status?token=om-recovery-2026-watchdog", data=b"{}", headers={"Content-Type":"application/json"}, method="POST")
print(json.loads(urllib.request.urlopen(req, timeout=15).read()))

# Restart
import urllib.request, json
req = urllib.request.Request("http://recovery.onemindos.com:18790/restart-gateway?token=om-recovery-2026-watchdog", data=b"{}", headers={"Content-Type":"application/json"}, method="POST")
print(json.loads(urllib.request.urlopen(req, timeout=40).read()))

# Logs
import urllib.request, json
req = urllib.request.Request("http://recovery.onemindos.com:18790/logs?token=om-recovery-2026-watchdog", data=b"{}", headers={"Content-Type":"application/json"}, method="POST")
print(json.loads(urllib.request.urlopen(req, timeout=15).read()))

# Update
import urllib.request, json
req = urllib.request.Request("http://recovery.onemindos.com:18790/update-gateway?token=om-recovery-2026-watchdog", data=b"{}", headers={"Content-Type":"application/json"}, method="POST")
print(json.loads(urllib.request.urlopen(req, timeout=120).read()))
```

---

## 6. DIGITALOCEAN FALLBACK (Nuclear Option)

Use after 3 consecutive webhook failures.

| Item | Value |
|------|-------|
| Droplet | onemindos-cloud |
| Droplet ID | 565014820 |
| IP | 64.227.8.23 |
| DO Token | `dop_v1_cf514923cdc71ecd431052d3d10ee98885fbc2e88b8f4c7badba640705bafbb3` |

```python
# Power cycle
import urllib.request, json
req = urllib.request.Request("https://api.digitalocean.com/v2/droplets/565014820/actions",
    data=b'{"type":"power_cycle"}',
    headers={"Authorization": "Bearer dop_v1_cf514923cdc71ecd431052d3d10ee98885fbc2e88b8f4c7badba640705bafbb3", "Content-Type": "application/json"},
    method="POST")
print(json.loads(urllib.request.urlopen(req, timeout=15).read()))
```

**Also noted:** omos-titan (Droplet ID: 567310143, IP: 104.248.226.252)

---

## 7. OMOS COMMANDS

### `omos status`
Full VPS health snapshot.
1. POST `/status` via recovery webhook
2. Report: gateway_http, openclaw_version, disk usage, memory usage
3. If gateway_http != 200 → immediately execute `omos down` protocol

### `omos down`
Gateway confirmed down. Skip health check. Act immediately.
1. Alert Zeus via Discord + Telegram: "⚠️ Gateway down. Attempting restart now..."
2. POST `/restart-gateway`
3. If recovered → "✅ Gateway restored after attempt X"
4. If still down → pull logs, diagnose, retry every 2 min
5. Retry up to 10 times
6. After 10 failures OR webhook unreachable → call_zeus (phone)
7. Always send final report: time down, attempts made, root cause if found

### `omos restart`
Force restart even if gateway appears healthy.
1. Alert Zeus: "🔄 Forcing gateway restart..."
2. POST `/restart-gateway`
3. Report result
4. If still_down → escalate to `omos down` protocol

### `omos update`
Update OpenClaw to latest version.
1. Alert Zeus: "🔄 Starting OpenClaw update..."
2. POST `/update-gateway` (timeout 120s)
3. Report: version_before, version_after, gateway_health
4. If unhealthy → execute `omos down` protocol

### `omos logs`
Last 50 lines of gateway logs.
1. POST `/logs`
2. Flag any ERROR, WARN, OOM, port conflict, or plugin crash lines

### `omos diagnose`
Full diagnostic — proactive issue identification.
1. POST `/status` → check gateway, disk, memory
2. POST `/logs` → scan for errors
3. Flag: disk > 80%, memory > 85%, any ERROR in logs, gateway_http != 200
4. Report with severity: OK / WARN / ALERT

### `omos up`
Confirm gateway is healthy after recovery.
1. POST `/status`
2. If 200 → "✅ Gateway confirmed up. Version: X. Disk: Y. Memory: Z."
3. If not 200 → execute `omos down` protocol

### `omos commands`
Print this command reference.

### `omos commands test`
Full system check on all endpoints and integrations.
1. Test /health, /status, /logs, /update-gateway, /restart-gateway
2. Test DigitalOcean API (read-only)
3. Test Tailscale path
4. Test Discord + Telegram direct APIs
5. Report pass/fail table

---

## 8. WATCHDOG AUTO-RECOVERY

**Critical:** Always run `/status` FIRST before sending any alerts.

1. POST `/status` → if gateway_http = 200 → **FALSE ALARM**, do nothing, log it
2. Confirmed down → alert Zeus via Discord + Telegram: "⚠️ Gateway down. Attempting restart..."
3. POST `/restart-gateway` → check health_after_15s
4. Recovered → notify "✅ Gateway restored after attempt X"
5. Still down → POST `/logs`, diagnose, retry every 2 min up to 10 attempts
6. Crash loop (fails within 60s) → report full diagnosis each attempt
7. 10 failures OR webhook unreachable → call_zeus (phone)
8. Always send final report

### Recovery Chain
1. `/restart-gateway` via webhook
2. Retry 3x, 2 min apart
3. DigitalOcean API power cycle
4. call_zeus (phone — emergencies only)

---

## 9. NOTIFICATION TOOLS (bypass OpenClaw)

These work even when the gateway is down:

| Tool | Target |
|------|--------|
| `send_discord_alert` | Zeus Discord DM (direct Discord API) |
| `telegram_send_alert` | Zeus Telegram alerts thread (direct Telegram API) |
| `call_zeus` | Phone call (EMERGENCIES ONLY: 10 failures or VPS unreachable) |

---

## 10. OPERATIONAL GOTCHAS

### Config Editing
- **Always write both `.json` and `.json.bak`** — auto-restore protection compares sizes
- Restart gateway after config changes: `systemctl --user restart openclaw-gateway.service` (sentinel) or `systemctl restart openclaw-gateway.service` (core)
- Sentinel gateway takes ~12s to fully boot

### Providers vs Plugins
- `amazon-bedrock` and `github-copilot` are **MODEL PROVIDERS**, not regular plugins
- They appear in `plugins.entries` for enable/disable but are configured via `agents.defaults.models`
- Don't confuse them with npm-installed plugins like `discord`

### AWS Bedrock Auth
- `AWS_PROFILE=assumed` in env **overrides** explicit IAM keys — AWS SDK uses instance role which lacks Bedrock permissions
- Fix: set `AWS_PROFILE=` (empty) in bedrock.conf drop-in to force explicit key usage
- Lightsail instance role ≠ Bedrock permissions

### Sentinel-specific
- User-level systemd: `systemctl --user` (not `systemctl` as root)
- SSH user: `ubuntu` (not root)
- 3.7GB RAM — monitor memory usage, gateway typically uses 389-504MB
- Clobbered config backups accumulate at `~/.openclaw/openclaw.json.clobbered.*` — safe to clean periodically

### Core-specific
- System-level systemd: `systemctl` (root)
- SSH user: `root`
- 30GB RAM — much more headroom
- Docker sandbox available (openclaw-sandbox:bookworm-slim)
- Multiple env drop-ins: azure, deepgram, digitalocean, fal, google, groq, vertex

---

## 11. KEY FILE PATHS

### Sentinel (`ubuntu@100.66.235.38`)
```
~/.openclaw/openclaw.json              # Main config
~/.openclaw/openclaw.json.bak          # Must match config (auto-restore protection)
~/.openclaw/gateway.systemd.env        # Environment vars (should be empty/minimal)
~/.openclaw/devices/paired.json        # Paired device registry
~/.openclaw/workspace/                 # Agent workspace (IDENTITY, SOUL, AGENTS, etc.)
~/.openclaw/workspace-coder/           # Coder sub-workspace
~/.openclaw/workspace-comms/           # Comms sub-workspace
~/.openclaw/workspace-media/           # Media sub-workspace
~/.openclaw/workspace-ops/             # Ops sub-workspace
~/.openclaw/skill-workshop/            # Captured skills
~/.openclaw/tasks/runs.sqlite          # Task run history
~/.openclaw/npm/                       # npm-installed plugins (discord)
~/.openclaw/patches/                   # Source patches (voice)
~/.config/systemd/user/openclaw-gateway.service        # Service unit
~/.config/systemd/user/openclaw-gateway.service.d/     # Drop-in dir
```

### Core (`root@100.123.126.34`)
```
/root/.openclaw/openclaw.json          # Main config
/root/.openclaw/workspace/             # Agent workspace
/etc/systemd/system/openclaw-gateway.service           # Service unit
/etc/systemd/system/openclaw-gateway.service.d/        # Drop-in dir (bedrock, azure, google, etc.)
```

---

## 12. QUICK REFERENCE COMMANDS

```bash
# Sentinel
ssh ubuntu@100.66.235.38 'systemctl --user status openclaw-gateway.service --no-pager | head -8'
ssh ubuntu@100.66.235.38 'systemctl --user restart openclaw-gateway.service'
ssh ubuntu@100.66.235.38 'openclaw models list 2>&1'
ssh ubuntu@100.66.235.38 'openclaw channels status --probe 2>&1'
ssh ubuntu@100.66.235.38 'openclaw plugins list 2>&1 | grep enabled'
ssh ubuntu@100.66.235.38 'journalctl --user -u openclaw-gateway.service -n 50 --no-pager'
ssh ubuntu@100.66.235.38 'openclaw infer model run --model "github-copilot/claude-sonnet-4.6" --prompt "test" --gateway 2>&1 | tail -3'

# Core
ssh root@100.123.126.34 'systemctl status openclaw-gateway.service --no-pager | head -8'
ssh root@100.123.126.34 'systemctl restart openclaw-gateway.service'
ssh root@100.123.126.34 'openclaw models list 2>&1'
ssh root@100.123.126.34 'openclaw channels status --probe 2>&1'
ssh root@100.123.126.34 'journalctl -u openclaw-gateway.service -n 50 --no-pager'
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
