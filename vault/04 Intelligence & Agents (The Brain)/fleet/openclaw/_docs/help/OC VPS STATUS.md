---
title: "OC VPS Status"
type: reference
status: active
created: 2026-04-09
updated: 2026-04-12
owner: zeus
tags: [type/reference, domain/so, infra/vps]
ai_summary: "Live snapshot of the OneMind VPS (onemind-main) — hardware, services, containers, ports, paths, and known issues as of April 9, 2026."
---

# OC VPS Status

> Last updated: **2026-04-12 09:25 UTC**

---

## Machine

| Field | Value |
|-------|-------|
| Hostname | `onemind-main` |
| Provider | DigitalOcean |
| OS | Ubuntu 24.04, Linux 6.8.0-71-generic (x86_64) |
| CPU | 4 vCPUs |
| RAM | 16 GB (2.5 GB used, 13 GB available) |
| Disk | 193 GB total, 52 GB used (27%), 142 GB free |
| Uptime | 26 days |
| Node.js | v22.22.1 |
| Public IP | `104.248.4.212` |
| Tailscale IP | `100.77.52.93` |

---

## Running Services (systemd)

| Service | Description | Status |
|---------|-------------|--------|
| **openclaw** | OpenClaw Gateway — Legacy's Runtime | ✅ Active |
| **openclaw-nats-ingester** | OpenClaw NATS Core → SQLite Event Ingester | ✅ Active |
| **legacy-pulse** | NATS Event Bridge → `AMBIENT_STATE.md` | ✅ Active |
| **tailscaled** | Tailscale mesh VPN agent | ✅ Active |
| **docker** | Docker container engine | ✅ Active |
| **datadog-agent** | Monitoring + APM (v7.76.3) | ✅ Active |
| **datadog-agent-trace** | Datadog Trace Agent (APM) | ✅ Active |
| **datadog-agent-installer** | Datadog Installer | ✅ Active |
| **do-agent** | DigitalOcean Monitoring Agent | ✅ Active |
| **droplet-agent** | DigitalOcean Droplet Agent | ✅ Active |

### Degraded Services

| Service | Issue |
|---------|-------|
| **openclaw-test.service** | Failed (test service using `op run` — not production) |

---

## Docker Containers

| Container | Image | Status | Ports |
|-----------|-------|--------|-------|
| **onemind-nats** | `nats:latest` | Up (healthy) | `:4222` (core), `:8222` (monitor), `:1883` (MQTT), `:7422` (leaf), `:9222` (websocket) |
| **nats-dashboard** | `mdawar/nats-dashboard:latest` | Up | `:18300` → 8080 (Tailscale: `http://100.77.52.93:18300`) |

Docker Compose file: `/opt/onemind/docker-compose.yml`

---

## Deployable Services (Stopped / Available)

Services with data volumes preserved on disk — not currently running but ready to restart.

### n8n (Workflow Automation)

| Field | Value |
|-------|-------|
| Image | `n8nio/n8n:latest` |
| Port | `:5678` |
| Data volume | `n8n_data` (preserved on disk) |
| Tailscale URL | `https://onemind-main.tail717ba0.ts.net/n8n` (route still configured) |
| Compose profile | `full` |

**To start:**
```bash
ssh root@104.248.4.212
cd /opt/onemind
docker compose --profile full up -d n8n
```

**To start as a systemd service (re-create):**
```bash
cat > /etc/systemd/system/n8n.service << 'EOF'
[Unit]
Description=OneMind n8n — Workflow Automation Engine
After=network-online.target docker.service

[Service]
Type=simple
User=onemind
WorkingDirectory=/opt/onemind
EnvironmentFile=-/opt/onemind/.env.sa
ExecStart=/usr/bin/docker compose --profile full up n8n
ExecStop=/usr/bin/docker compose --profile full down n8n
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable --now n8n.service
```

**Access:** `http://100.77.52.93:5678` (Tailscale) or via Tailscale serve at `/n8n`

**Why it was stopped:** Not actively used for automations at time of removal (2026-04-12). Data untouched — all workflows still in `n8n_data` volume.

---

## Timers

| Timer | Frequency | Purpose |
|-------|-----------|---------|

| **workspace-push** | Every 5 min | Push OpenClaw workspace changes to git |

---

## OpenClaw

| Field | Value |
|-------|-------|
| Version | `2026.4.11` |
| Config | `/root/.openclaw/openclaw.json` |
| Gateway port | `:18789` (internal), exposed via Tailscale |
| Run as user | `onemind` |

### Agent Workspaces

All 9 agent workspaces at `/home/onemind/onemind-os/openclaw/workspaces/`:

| Workspace | Agent | Domain |
|-----------|-------|--------|
| `legacy` | ⚡ Legacy | All — central orchestrator |
| `spartan` | 💪 Spartan | HP 27-28 — body, fitness |
| `oracle` | 🔮 Oracle | HP 25-26, 29-34 — mind, money, meaning |
| `eden` | 🌱 Eden | LE — sustainability, environment |
| `haven` | 🏡 Haven | LE 50-74 — family, home, heritage |
| `empire` | 💰 Empire | GE 75-99 — wealth, revenue, brand |
| `guardian` | 🛡️ Guardian | SO 11 — security, access control |
| `grid` | ⚙️ Grid | SO 01, 09 — infrastructure, automation |
| `trinity` | 🔺 Trinity | SO 03 — agent coordination |

### Memory Indexing

| Database | Size | Notes |
|----------|------|-------|
| `legacy.sqlite` | 8.7 GB | Primary memory index |
| `empire.sqlite` | 998 MB | Empire agent memory |

Indexing was at ~26% as of April 8 — progressive embedding via OpenAI API ongoing.

### Shared Context Files

Located at `/home/onemind/onemind-os/openclaw/shared/` and symlinked into all workspaces:
- `CODEX.md` — What the codex is, paths, structure, terminology
- `USER.md` — Zeus's details, communication style, priorities
- `MEMORY-STACK.md` — Memory stack configuration reference (added 2026-04-12)

Universal vault search configured via `agents.defaults.memorySearch.extraPaths` → `/home/onemind/onemind-os/codex/`

---

### Memory Stack (Activated 2026-04-12)

Five layers activated on top of the base OpenClaw runtime. Full SOP: `00-24 SO/03 Agents/memory-stack-configuration.md`

```
LAYER 5 — memory-wiki    → obsidian wiki vault from dream reports + MEMORY.md
LAYER 4 — dreaming       → nightly consolidation (03:00 UTC) → writes MEMORY.md
LAYER 3 — QMD v2.1.0    → semantic vector search across full codex vault
LAYER 2 — active-memory  → preference summary injected before each turn
LAYER 1 — lossless-claw  → DAG-based context compaction (never loses messages)
```

#### Installed Plugins

| Plugin | Version | Source | Slot |
|--------|---------|--------|------|
| `lossless-claw` | 0.8.2 | `@martian-engineering/lossless-claw` | `contextEngine` (exclusive) |
| `active-memory` | stock | openclaw built-in | — |
| `memory-core` | stock | openclaw built-in | — |
| `memory-wiki` | stock | openclaw built-in | — |
| `qmd` | 2.1.0 | `@tobilu/qmd` (binary: `/usr/bin/qmd`) | — |

#### Key Config Paths in `~/.openclaw/openclaw.json`

| Key | Value | Purpose |
|-----|-------|---------|
| `memory.backend` | `"qmd"` | Activates QMD as vector search backend |
| `memory.qmd.update.onBoot` | `true` | Re-indexes codex on gateway startup |
| `memory.qmd.sessions.enabled` | `true` | Also indexes session transcripts |
| `agents.defaults.memorySearch.extraPaths` | `["/home/onemind/onemind-os/codex/"]` | Codex vault as QMD collection |
| `agents.defaults.compaction.mode` | `"safeguard"` | lossless-claw sentinel mode |
| `plugins.slots.contextEngine` | `"lossless-claw"` | Exclusive context compaction slot |
| `plugins.entries["active-memory"].config.agents` | all 9 | Every agent gets preference injection |
| `plugins.entries["active-memory"].config.promptStyle` | `"preference-only"` | Minimal, non-intrusive injection |
| `plugins.entries["memory-core"].config.dreaming.enabled` | `true` | Nightly memory consolidation |
| `plugins.entries["memory-wiki"].config.vaultMode` | `"bridge"` | Read memory-core → compile wiki |
| `plugins.entries["memory-wiki"].config.vault.path` | `…/codex/_codex/memory/wiki` | Wiki output in codex |

#### lossless-claw Config

| Field | Value | Meaning |
|-------|-------|---------|
| `freshTailCount` | 64 | Recent messages protected from compaction |
| `leafChunkTokens` | 80000 | Source tokens per leaf before summarization |
| `contextThreshold` | 0.75 | % of context window that triggers compaction |
| `incrementalMaxDepth` | 1 | One condensed pass per turn |
| `summaryModel` | `github-copilot/gpt-4o` | Model used for summarization |

**DB:** `~/.openclaw/lcm.db` | **Agent commands:** `/lcm`, `/lcm doctor`

#### Backups

| File | Created | Notes |
|------|---------|-------|
| `~/.openclaw/openclaw.json.bak.20260412-090153` | 2026-04-12 | Pre-memory-stack backup |

---

## Tailscale Serve (Reverse Proxy)

| URL | Backend |
|-----|---------|
| `https://onemind-main.tail717ba0.ts.net/` | `http://127.0.0.1:18789` (OpenClaw gateway) |
| `https://onemind-main.tail717ba0.ts.net/n8n` | `http://127.0.0.1:5678` (n8n) |
| `https://onemind-main.tail717ba0.ts.net/nats` | `http://127.0.0.1:8222` (NATS monitor) |
| `https://onemind-main.tail717ba0.ts.net:3456/` | `http://127.0.0.1:3456` (onemind-web, when running) |
| `https://onemind-main.tail717ba0.ts.net:8443/` | `http://127.0.0.1:18789` (OpenClaw gateway alt) |

All Tailscale-served endpoints are **tailnet only** — not publicly accessible.

---

## Firewall (UFW)

UFW is **active**. Default policy: deny incoming, allow outgoing.

| Rule | Port | Interface | Notes |
|------|------|-----------|-------|
| ALLOW | 22/tcp | any | SSH |
| ALLOW | 4222 | tailscale0 | NATS Core (Tailscale only) |
| ALLOW | 18789 | tailscale0 | OpenClaw Gateway (Tailscale only) |
| ALLOW | 18300 | tailscale0 | NATS Dashboard (Tailscale only) |

> Docker ports (1883, 4222, 7422, 8222, 9222) are bound to `0.0.0.0` at the Docker level but UFW restricts access to Tailscale interface where configured above. NATS MQTT (1883), leafnode (7422), and WebSocket (9222) remain accessible from the public internet — no UFW rule restricts them.

---

## Key Ports

| Port | Service | Bind | Access |
|------|---------|------|--------|
| 22 | SSH | `0.0.0.0` | Public |
| 443 | Tailscale HTTPS | `100.77.52.93` | Tailscale only |
| 1883 | NATS MQTT | `0.0.0.0` | Public (no UFW rule) |
| 3456 | onemind-web Tailscale serve | `100.77.52.93` | Tailscale only — route configured but onemind-web NOT running on VPS |
| 4222 | NATS Core | `0.0.0.0` | UFW: Tailscale only |
| 5000-5001 | Datadog Agent | `127.0.0.1` | Localhost |
| 5012 | Datadog Trace | `127.0.0.1` | Localhost |
| 7422 | NATS Leafnode | `0.0.0.0` | Public (no UFW rule) |
| 8126 | Datadog APM | `127.0.0.1` | Localhost |
| 8222 | NATS Monitor | `0.0.0.0` | UFW: Tailscale only |
| 8443 | Tailscale HTTPS alt | `100.77.52.93` | Tailscale only |
| 9222 | NATS WebSocket | `0.0.0.0` | Public (no UFW rule) |
| 18300 | NATS Dashboard | `0.0.0.0` | UFW: Tailscale only |
| 18789 | OpenClaw Gateway (main) | `0.0.0.0` | UFW: Tailscale only |
| 18791 | OpenClaw Gateway (WebSocket) | `127.0.0.1` | Localhost |
| 40607 | OpenClaw Gateway (internal) | `127.0.0.1` | Localhost |

---

## Key Paths

| Path | Purpose |
|------|---------|
| `/home/onemind/onemind-os/codex/` | Codex vault (git clone) |
| `/home/onemind/.openclaw/` | OpenClaw config, memory, tasks |
| `/home/onemind/.openclaw/openclaw.json` | Main OpenClaw configuration |
| `/home/onemind/.openclaw/memory/` | Agent memory SQLite databases |
| `/home/onemind/onemind-os/openclaw/workspaces/` | 9 agent workspace directories |
| `/home/onemind/onemind-os/openclaw/shared/` | Shared context files (symlinked into all 9 workspaces) |
| `/root/.openclaw/openclaw.json` | OpenClaw config (also symlinked at `/home/onemind/.openclaw/openclaw.json`) |
| `/root/.openclaw/lcm.db` | lossless-claw DAG database |
| `/root/.openclaw/extensions/lossless-claw/` | lossless-claw extension (v0.8.2) |
| `/usr/bin/qmd` | QMD binary (v2.1.0) |
| `/home/onemind/onemind-os/codex/_codex/memory/wiki/` | memory-wiki output vault |
| `/opt/onemind/` | Docker Compose stack (NATS, OmniDB, n8n) |
| `/opt/onemind/.env.tpl` | 1Password env template |
| `/opt/onemind/.env.sa` | Service account env |

---

## Legacy-Pulse Service

```
[Service]
Type=simple
User=onemind
WorkingDirectory=/home/onemind/onemind-os/_system/openclaw/workspaces/legacy/skills/legacy-pulse
ExecStart=/usr/bin/python3 -u legacy-pulse.py
Restart=always
RestartSec=5
Environment=NATS_URL=nats://localhost:4222
Environment=OPENCLAW_WORKSPACE=/home/onemind/.openclaw/workspace
```

Bridges NATS events into the `AMBIENT_STATE.md` file for agent context awareness.

---

## Known Issues

1. **Discord channel disabled** — bad token; disabled in `openclaw.json` to prevent crash loops
2. **NATS partially exposed** — MQTT (1883), leafnode (7422), and WebSocket (9222) are bound to `0.0.0.0` with no UFW restriction. NATS Core (4222) and dashboard (18300) are restricted to tailscale0 by UFW. Consider adding UFW rules for 1883/7422/9222 if not intentionally public.
3. **legacy-pulse.service degraded** — auto-restart loop with status=200/CHDIR. `legacy-pulse.py` or its working directory (`workspaces/legacy/skills/legacy-pulse`) is missing. Needs investigation.
4. **Memory indexing incomplete** — ~26% indexed, progressive embedding ongoing
5. **openclaw-test.service failed** — test service, non-critical
6. **onemind-web not running on VPS** — Tailscale serve route `:3456` configured but no backend; `http://100.77.52.93:3456` is currently dead

---

## Sync Architecture (Current)

No automated sync is running. VPS codex is a stale git clone.

```
Mac (Obsidian / VS Code)
    ↓ git push (manual)
OneMind-OS/onemind-codex (GitHub)
    ↓ manual git pull on VPS when needed
/home/onemind/onemind-os/codex/
```

**Planned migration:** Replace with Obsidian Headless Sync (`ob sync --continuous`) for real-time bidirectional sync. See CLAUDE.md Phase 1.

---

## Deprecated Services (Removed 2026-04-09)

These services were removed from the VPS. Documented here for reference.

### onemind-pull.service + onemind-pull.timer

**What it did:** Timer-triggered `git pull --rebase --autostash` on the codex vault every 60 seconds. Picked up changes pushed from Mac/Obsidian via GitHub so the VPS had a near-current copy of the vault for agent memory search.

**Why removed:**
- Stuck in a perpetual rebase failure loop (exit code 128)
- Timer kept firing every minute, each attempt failing immediately
- Root cause: diverged history after the `git-filter-repo` cleanup (46.2 GB → 273 MB) on April 5 made the VPS clone incompatible with the force-pushed remote
- Will be replaced by Obsidian Headless Sync which handles ALL files (including binaries) with real-time E2E encrypted sync — no git involvement

**Unit files (deleted):**
- `/etc/systemd/system/onemind-pull.service`
- `/etc/systemd/system/onemind-pull.timer`

**Original service definition:**
```ini
[Service]
Type=oneshot
WorkingDirectory=/home/onemind/onemind-os/codex
ExecStart=/usr/bin/git pull --rebase --autostash -q
User=onemind

[Timer]
OnCalendar=*:*:05
Persistent=true
```

### onemind-watch.service

**What it did:** Filesystem watcher (inotifywait-based) that detected file changes in the codex vault and auto-committed + pushed them to GitHub. Designed so agent edits on VPS would propagate back to Mac/phone via git.

**Why removed:**
- Disabled since 2026-03-17 due to detached HEAD bug
- Would enter detached HEAD state after rebases, then all pushes would fail with non-fast-forward errors
- Even when "running", it was spamming failed push attempts (last seen pushing `AMBIENT_STATE.md` changes that kept getting rejected)
- The auto-commit messages were generic (`agent: update <files>`) with no meaningful context
- Will be replaced by Obsidian Headless Sync (bidirectional, no git needed for real-time sync)

**Unit file (deleted):**
- `/etc/systemd/system/onemind-watch.service`

**Original service definition:**
```ini
[Unit]
Description=OneMind file-watcher — auto-push workspace changes to GitHub

[Service]
Type=simple
WorkingDirectory=/home/onemind/onemind-os/codex
ExecStart=/home/onemind/onemind-os/codex/_system/infrastructure/onemind-watch.sh
User=onemind
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**The watch script** used `inotifywait` to monitor the vault directory, debounced changes for ~3 seconds, then ran `git add . && git commit && git push`. The script may still exist at `_system/infrastructure/onemind-watch.sh` in the vault but is no longer deployed.

### onemind-n8n (container) + n8n.service — Removed 2026-04-12

**What it was:** n8n workflow automation engine. Docker container (`n8nio/n8n:latest`) running on port `:5678`, wrapped in a systemd service. Was running for ~3 weeks.

**Why removed:**
- Not actively used for any current automations
- Volume `n8n_data` preserved on disk — data is safe if n8n needs to come back
- Can be re-deployed via `docker compose --profile full up n8n` if needed

**Service definition (deleted: `/etc/systemd/system/n8n.service`):**
```ini
[Unit]
Description=OneMind n8n — Workflow Automation Engine
After=network-online.target docker.service

[Service]
Type=simple
User=onemind
WorkingDirectory=/home/onemind/onemind-os/infra
EnvironmentFile=-/opt/onemind/.env.sa
ExecStart=/usr/bin/docker compose --profile full up n8n
ExecStop=/usr/bin/docker compose --profile full down n8n
```

**Tailscale serve route** `https://onemind-main.tail717ba0.ts.net/n8n` pointed to `:5678` — no longer valid.

---

### onemind-omnidb (container) — Removed 2026-04-12

**What it was:** PostgreSQL 16 + TimescaleDB + pgvector (`timescale/timescaledb:latest-pg16`), port `:5432`. Was deployed to store NATS events and provide a vector store for agent memory.

**Why removed:**
- The NATS ingester was never wired up — nothing was actually writing to it
- OpenClaw uses its own native memory (SQLite `.sqlite` files per agent) — no need for external Postgres
- ~400 MB RAM freed
- Volume `omnidb_data` preserved on disk

**Replaced by:** `openclaw-nats-ingester.service` — lightweight Python process writing NATS events to `/opt/onemind/events.db` (SQLite, WAL mode, 15.6 MB RAM).

---

### onemind-ingester.service — Renamed 2026-04-12

**What it was:** First deployment name for the NATS → SQLite event ingester, created April 12, 2026 (same session as OmniDB removal).

**Why renamed:** Name didn't clearly indicate which system it belonged to or what it connected. Renamed to `openclaw-nats-ingester.service` for clarity.

**Current name:** `openclaw-nats-ingester.service` — active and running.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
