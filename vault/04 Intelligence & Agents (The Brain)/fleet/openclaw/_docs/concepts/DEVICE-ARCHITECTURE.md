---
title: "Device Architecture — OneMind Node Network"
type: reference
status: active
created: 2026-04-16
owner: zeus
tags: [type/reference, domain/SO-01, topic/infrastructure, topic/nodes]
ai_summary: "Physical device architecture for Zeus's OneMind node network — Mac, VPS, Dell AIO, Beelink, iPhone. Roles, specs, and recommended usage for each device."
---

# Device Architecture — OneMind Node Network

Zeus's physical compute footprint across 6 devices. Each device serves a specific role in the OneMind ecosystem, connected via Tailscale mesh network to the VPS gateway.

---

## Network Topology

```
                    ┌──────────────────────────────────────┐
                    │   VPS: onemindos-cloud-1             │
                    │   DigitalOcean · 64.227.8.23         │
                    │   Tailscale: 100.98.25.10            │
                    │                                      │
                    │   ROLE: GATEWAY + CENTRAL BRAIN      │
                    │   ─────────────────────────────       │
                    │   OpenClaw Gateway (9 agents)         │
                    │   Ollama (LLM inference)              │
                    │   PostgreSQL + pgvector               │
                    │   Tailscale Serve (HTTPS proxy)       │
                    └──────────┬───────────────────────────┘
                               │
            Tailscale Mesh (WireGuard, encrypted)
     ┌─────────┬───────────┼───────────┬──────────┬──────────┐
     │         │           │           │          │          │
┌────▼───┐ ┌───▼────┐ ┌───▼────┐ ┌───▼────┐ ┌───▼─────┐   │
│  Mac   │ │ Dell   │ │Beelink │ │ iPhone │ │ Android │   │
│ M2 Max │ │ AIO    │ │ Mini   │ │        │ │ TV      │   │
│        │ │ Ubuntu │ │        │ │        │ │         │   │
│ SOURCE │ │ SCREEN │ │ EDGE   │ │ MOBILE │ │ AMBIENT │   │
│ OF     │ │ + EDGE │ │ HOME   │ │ ACCESS │ │ DISPLAY │   │
│ TRUTH  │ │ NODE   │ │ NODE   │ │        │ │         │   │
└────────┘ └────────┘ └────────┘ └────────┘ └─────────┘   │
```

---

## Device Inventory

### 1. MacBook Pro (M2 Max) — SOURCE OF TRUTH

| Field | Value |
|-------|-------|
| Model | MacBook Pro (Mac14,6) |
| Chip | Apple M2 Max |
| RAM | 32 GB |
| Role | Primary workstation, vault source of truth |
| Tailscale IP | `100.114.142.55` |
| OpenClaw | App (Direct WSS) + Headless Node (LaunchAgent) |
| Plan | Leave open 24/7 (lid-open or clamshell) |

**What runs here:**
- Obsidian vault (One Mind — the brain, ~78 GB)
- VS Code (development, agent work)
- OpenClaw macOS app (camera, canvas, screen, voice)
- OpenClaw headless node service (LaunchAgent — always-on system commands)
- Git (source of truth → GitHub → VPS)
- All creative work, development, writing

**Node capabilities:** browser, camera, canvas, screen, system, notifications, voice wake, accessibility, AppleScript

**Why Mac stays as source of truth:**
- Only device with full Obsidian + Git + VS Code
- Apple Silicon M2 Max can run local Ollama inference when needed
- Full TCC permissions for screen recording, camera, automation
- Most capable node in the network

---

### 2. VPS: onemindos-cloud-1 — GATEWAY + CENTRAL BRAIN

| Field | Value |
|-------|-------|
| Provider | DigitalOcean |
| Public IP | `64.227.8.23` |
| Tailscale | `100.98.25.10` / `omos-core.tail717ba0.ts.net` |
| CPU | Intel Xeon Gold 6248 @ 2.50GHz |
| RAM | 32 GB |
| Disk | 290 GB (17 GB used) |
| OS | Ubuntu 24.04 |
| OpenClaw | v2026.4.14 — Gateway (9 agents) |

**What runs here:**
- OpenClaw gateway (port 18789, loopback via Tailscale Serve)
- 9 AI agents: Legacy⚡, Spartan💪, Oracle🔮, Eden🌱, Haven🏡, Empire💰, Guardian🛡️, Grid⚙️, Trinity🔺
- Ollama (qwen2.5:7b, qwen3:4b, gemma4, etc.)
- PostgreSQL + pgvector (embeddings, RAG pipeline)
- Tailscale Serve (HTTPS termination)
- Git mirror of One Mind vault (via onemind-pull.timer)

**Not a node** — it's the gateway. All nodes connect TO it.

---

### 3. Dell All-in-One (Ubuntu) — SCREEN + EDGE NODE

| Field | Value |
|-------|-------|
| Model | Dell Inspiron All-in-One (exact model TBD) |
| OS | Ubuntu (installed, minimal use) |
| Current use | External display for MacBook Pro |
| Status | Powered on (used as screen), Ubuntu idle |
| Tailscale | Not yet joined |

**Recommended role: SECOND NODE + LOCAL SERVICES**

The Dell is always on because Zeus uses it as a display. Since it's running Ubuntu, it's a free compute node sitting idle. Recommended uses:

1. **OpenClaw headless node** — Join the Tailscale mesh, install OpenClaw CLI, run `openclaw node install`. The gateway gets a second always-on node with independent system access (Linux).

2. **Monitoring dashboard** — Run a persistent browser window showing the OpenClaw dashboard, Datadog metrics, or a custom One Mind status board. The Dell screen becomes a command center display.

3. **Local Ollama inference** — If the Dell has decent specs (8+ GB RAM), run a small Ollama model locally. Faster than VPS CPU inference for quick queries.

4. **n8n / automation runner** — Move n8n workflows from VPS to local. Keeps VPS focused on the gateway.

**Setup steps:**
```bash
# 1. Join Tailscale mesh
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --ssh

# 2. Install Node.js + OpenClaw
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt-get install -y nodejs
sudo npm install -g openclaw

# 3. Connect as node
openclaw node install --host omos-core.tail717ba0.ts.net --port 443 --tls

# 4. Approve on VPS
ssh root@64.227.8.23 'openclaw nodes approve --all'
```

---

### 4. Beelink Mini PC — EDGE HOME NODE

| Field | Value |
|-------|-------|
| Model | Beelink Mini PC (exact model TBD) |
| OS | Linux (likely) |
| Current use | Was Home Assistant, now idle |
| Status | Available, not in use |
| Tailscale | Not yet joined |

**Recommended role: HOME AUTOMATION + EDGE AI**

The Beelink is purpose-built for always-on edge computing. Low power, compact, silent. Best uses:

1. **Home Assistant (restore it)** — This was its original purpose and it's the right device for it. Home Assistant gives Legacy control over lights, sensors, climate, presence detection. This feeds the HP and LE domains directly.

2. **OpenClaw headless node** — After Home Assistant is running, also install OpenClaw node. Legacy can then invoke Home Assistant commands on the Beelink directly.

3. **Local RAG / embeddings** — Run a lightweight embedding server (nomic-embed-text via Ollama) for local document search without hitting the VPS.

4. **Backup sync endpoint** — Run Syncthing or rclone to keep a local backup of the One Mind vault.

**Setup steps:**
```bash
# 1. Restore Home Assistant
# Option A: Home Assistant OS (dedicated — recommended)
# Flash HAOS to USB and boot from it
# Option B: Docker (if keeping Linux)
docker run -d --name homeassistant \
  --privileged --restart=unless-stopped \
  -v /opt/homeassistant:/config \
  -e TZ=America/New_York \
  --network=host \
  ghcr.io/home-assistant/home-assistant:stable

# 2. Join Tailscale mesh
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --ssh

# 3. Install OpenClaw node (if running Linux, not HAOS)
sudo npm install -g openclaw
openclaw node install --host omos-core.tail717ba0.ts.net --port 443 --tls
```

---

### 5. iPhone — MOBILE ACCESS

| Field | Value |
|-------|-------|
| Model | iPhone (model TBD) |
| Current use | Daily carry |
| OpenClaw access | Via Chowder app (sideloaded) or Telegram bot |
| Tailscale | iOS Tailscale app |

**Role: READ + COMMAND**

iPhone is the mobile command surface. It can't run background services but it's the fastest way to reach Legacy on the go.

**What runs here:**
- Obsidian Mobile (vault read + light edits, via Obsidian Sync when enabled)
- Chowder iOS app (OpenClaw chat, voice, canvas)
- Telegram → Legacy (via Telegram channel on VPS)
- Tailscale (mesh network access)

**iPhone is NOT a node.** iOS doesn't support persistent background services. Communication goes through Telegram or Chowder when the app is open.

---

### 6. Android TV — AMBIENT DISPLAY

| Field | Value |
|-------|-------|
| Model | Android TV (exact model TBD) |
| OS | Android TV |
| Current use | TV |
| Status | Available |
| Tailscale | Not yet joined |

**Role: AMBIENT DASHBOARD + CAST TARGET**

The Android TV sits in the same room as the Mac and Dell. It's a passive display surface for persistent, glanceable information.

**What it does:**
1. **OpenClaw Control UI** — Install Tailscale from Google Play Store → open `https://omos-core.tail717ba0.ts.net` in the TV browser → persistent agent dashboard on the big screen
2. **Home Assistant dashboard** — When Beelink runs HA, display the HA dashboard (weather, cameras, sensors, calendar) on the TV
3. **Chromecast target** — Cast agent-generated reports, slides, media from Mac or phone
4. **Kiosk mode** — Android TV browsers support fullscreen/kiosk — pin a dashboard permanently

**Setup steps:**
```
1. Install Tailscale from Google Play Store on the TV
2. Sign into the tailnet (same account as all other devices)
3. Open the TV's web browser
4. Navigate to https://omos-core.tail717ba0.ts.net
5. Bookmark for quick access
```

**Android TV is NOT a node.** It's a display-only device. No OpenClaw agent runs on it — it just shows dashboards via the browser.

---

## Role Summary

| Device | Role | Always On? | OpenClaw Role | Key Capability |
|--------|------|-----------|---------------|----------------|
| **Mac** | Source of truth | Yes (planned) | App node + headless node | Full vault, dev, camera, canvas |
| **VPS** | Central brain | Yes | Gateway (not a node) | 9 agents, LLM inference, DB |
| **Dell** | Display + edge | Yes (always on) | Headless node (planned) | Monitoring dashboard, local compute |
| **Beelink** | Home automation | Yes (planned) | Headless node (planned) | Home Assistant, edge AI |
| **iPhone** | Mobile access | Intermittent | Not a node (client only) | Chat, voice, quick captures |
| **Android TV** | Ambient display | Yes | Not a node (display only) | Dashboard, HA display, cast target |

---

## Compute Distribution Strategy

```
                    HEAVY COMPUTE         LIGHT / EDGE
                    ─────────────         ────────────
  LLM Inference     VPS (Ollama)          Mac (Ollama, when needed)
  Agent Execution   VPS (Gateway)         —
  Embeddings/RAG    VPS (pgvector)        Beelink (nomic-embed, future)
  Home Automation   —                     Beelink (Home Assistant)
  Vault/Git         Mac (source)          VPS (mirror)
  Development       Mac (VS Code)         —
  Monitoring        Dell (dashboard)      VPS (Datadog)
  Ambient Display   Android TV (dashboards, HA, casts)
  Mobile Chat       —                     iPhone (Chowder/Telegram)
```

**Design principle:** VPS handles anything that needs to be fast and always available. Mac handles anything that needs full vault access or Apple hardware. Edge devices (Dell, Beelink) handle always-on local services that don't need the cloud.

---

## Home Assistant Integration

Two integration directions — both solved with existing tools.

### Direction 1: OpenClaw → HA (agent controls devices)

Install `homeassistant-skill` on the VPS. Agents call HA REST API via Tailscale.

```bash
openclaw skills install homeassistant-skill
```

**Flow:** Legacy (VPS) → `homeassistant-skill` → HA REST API → Beelink HA → physical devices

### Direction 2: Voice → OpenClaw (speak to Legacy via HA)

Install the **OpenClaw HACS Integration** on Home Assistant. This is a community custom integration (v0.1.62, actively maintained) that connects HA directly to your existing VPS gateway using the OpenAI-compatible API (`/v1/chat/completions`).

**GitHub:** `https://github.com/techartdev/OpenClawHomeAssistantIntegration`

What it provides:
- Native **Assist conversation agent** — OpenClaw appears as a selectable agent in Voice Assistants
- **Lovelace chat card** (`custom:openclaw-chat-card`) with history, typing indicator, voice button
- **Service calls** for automations: `openclaw.send_message`, `openclaw.clear_history`, `openclaw.invoke_tool`
- **Event hooks**: `openclaw_message_received`, `openclaw_tool_invoked`
- **Tool calls**: HA enables `execute_service` tools so Legacy can control HA devices directly
- **Voice input**: Browser Web Speech or HA STT provider, wake word support ("hey openclaw")
- **Model selector** entity on dashboard

**Prerequisites (already done on VPS):**
- `gateway.http.endpoints.chatCompletions.enabled: true` ✅ (set 2026-04-16)
- Endpoint verified: `POST /v1/chat/completions` → Legacy responds via `openclaw/legacy` model ✅

**Setup (when Beelink HA is running):**
1. HACS → Custom repositories → add `https://github.com/techartdev/OpenClawHomeAssistantIntegration` → Category: Integration
2. Install → Restart HA
3. Settings → Devices & Services → Add Integration → OpenClaw
4. Manual config (remote instance — no local addon):
   - Gateway Host: `100.98.25.10` (VPS Tailscale IP)
   - Gateway Port: `18789`
   - Gateway Token: from 1Password `SO / OpenClaw Gateway / gateway-token`
   - Use SSL: ❌ (Tailscale handles encryption)
5. Configure → Enable tool calls (lets Legacy control HA devices)
6. Settings → Voice Assistants → set conversation agent to "OpenClaw"

**Voice pipeline:** Wake word → HA STT → HACS conversation agent → `/v1/chat/completions` → Legacy → TTS → speaker

> **Note:** With tool calls enabled in the HACS integration, `homeassistant-skill` on the VPS may not be needed — the integration provides bidirectional control (voice→Legacy and Legacy→HA devices) natively.

---

## Future Expansion

- **Beelink → Home Assistant + OpenClaw node** — Priority. Gives Legacy eyes and ears in the physical space.
- **Dell → Monitoring node** — Low effort. Join Tailscale, install OpenClaw, show dashboard.
- **Android TV → Ambient dashboard** — Lowest effort. Install Tailscale, open browser, bookmark Control UI.
- **Mac → 24/7 node** — Leave lid open. Most powerful node in the mesh.
- **GPU VPS (future)** — If Ollama inference on CPU becomes a bottleneck, add a GPU droplet or switch to a cloud LLM API.
- **Raspberry Pi (future)** — Cheap, silent, perfect for additional edge nodes (car, office, gym bag).

---

## Tailscale Mesh Status

All devices should be on the same Tailscale tailnet (`tail717ba0.ts.net`).

| Device | Tailscale Status | Hostname | Next Step |
|--------|-----------------|----------|-----------|
| Mac | ✅ Connected | `zeuss-macbook-pro` | — |
| VPS | ✅ Connected | `onemindos-cloud-1` | — |
| Dell | ❌ Not joined | — | `curl -fsSL https://tailscale.com/install.sh \| sh && sudo tailscale up --ssh` |
| Beelink | ❌ Not joined | — | Install HAOS first, then join Tailscale from HAOS settings |
| iPhone | ❌ Not joined | — | Install Tailscale from App Store → sign in → same tailnet |
| Android TV | ❌ Not joined | — | Install Tailscale from Google Play → sign in → same tailnet |

---

## Setup Checklist

Everything still pending, ordered by priority.

### ✅ Already Done (2026-04-16)

- [x] Mac OpenClaw app connected — Direct WSS transport, paired
- [x] Mac headless node running — LaunchAgent, always-on
- [x] VPS gateway running — 9 agents, OpenClaw 2026.4.14
- [x] VPS chatCompletions API enabled — `/v1/chat/completions` verified, Legacy responds
- [x] Mac + VPS on Tailscale mesh
- [x] 1Password: gateway entry updated, old entry archived

### 🔴 Priority 1 — Beelink + Home Assistant

- [ ] Flash HAOS to USB/SD and boot Beelink from it (recommended over Docker)
- [ ] Join Beelink to Tailscale from HA Settings → Add-ons → Tailscale
- [ ] Install HACS on HA (Settings → Add-ons or manual install)
- [ ] Add HACS custom repo: `https://github.com/techartdev/OpenClawHomeAssistantIntegration`
- [ ] Install OpenClaw HACS Integration → Restart HA
- [ ] Configure integration — Manual:
  - Host: `100.98.25.10` (VPS Tailscale IP)
  - Port: `18789`
  - Token: from 1Password `SO / OpenClaw Gateway / gateway-token`
  - Use SSL: ❌ (Tailscale handles it)
- [ ] Configure → **Enable tool calls** (Legacy controls HA devices bidirectionally)
- [ ] Settings → Voice Assistants → set conversation agent to **OpenClaw**
- [ ] Add chat card to HA dashboard: `type: custom:openclaw-chat-card`
- [ ] Test: type a message in the chat card → Legacy responds
- [ ] (Optional) `openclaw skills install homeassistant-skill` on VPS — only if agent needs to call HA outside of voice/chat sessions

### 🔴 Priority 2 — Voice Hardware

- [ ] Choose voice satellite hardware:
  - **HA Voice Preview Edition** ($59) — plug-and-play, wake word built in, recommended
  - **ESP32 ATOM Echo** ($13) — DIY, ESPHome, more setup
  - **Any microphone-speaker combo** running ESPHome Wyoming protocol
- [ ] Plug in and pair with HA via Settings → Devices
- [ ] Configure HA voice pipeline: Settings → Voice Assistants → pipeline with OpenClaw as conversation agent
- [ ] Set wake word (default: `hey openclaw`)
- [ ] Test end-to-end: speak → Legacy responds via speaker

### 🟡 Priority 3 — Dell (Monitoring Node)

- [ ] Install Tailscale on Dell: `curl -fsSL https://tailscale.com/install.sh | sh && sudo tailscale up --ssh`
- [ ] Install Node.js 22: `curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && sudo apt-get install -y nodejs`
- [ ] Install OpenClaw: `sudo npm install -g openclaw`
- [ ] Connect as node: `openclaw node install --host omos-core.tail717ba0.ts.net --port 443 --tls`
- [ ] Approve on VPS: `ssh root@64.227.8.23 'openclaw nodes approve --all'`
- [ ] Set up persistent monitoring dashboard in browser (pin to screen)

### 🟡 Priority 4 — Android TV (Ambient Dashboard)

- [ ] Install Tailscale from Google Play Store on the TV
- [ ] Sign in to same tailnet
- [ ] Open browser → `https://omos-core.tail717ba0.ts.net`
- [ ] Set to fullscreen/kiosk mode, bookmark
- [ ] (Optional) Add HA dashboard URL once Beelink HA is running

### 🟡 Priority 5 — iPhone

- [ ] Install Tailscale from App Store → sign in to tailnet
- [ ] Verify Chowder app still works (re-deploy from Xcode if expired)

### ⚪ Priority 6 — VPS Agent Speed

- [ ] Try smaller Ollama model for agent tool-schema: `qwen2.5:1.5b` or `qwen2.5:3b` (already on disk, much faster on CPU)
  ```bash
  ssh root@64.227.8.23 'openclaw config set agents.defaults.model.primary ollama/qwen2.5:3b'
  ```
- [ ] Re-enable `onemind-watch.service` (add detached HEAD guard, see CLAUDE.md)
- [ ] Memory indexing completing progressively — no action, just monitor

### ⚪ Priority 7 — Mac 24/7

- [ ] Enable clamshell mode (close lid, use external display via dock) or simply leave lid open
- [ ] Verify headless node stays connected overnight


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
