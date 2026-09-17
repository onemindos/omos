---
id: dashboard-onemind-ops
title: OneMind OS — Operations Dashboard
tags: [infrastructure, dashboard, services, homepage, ops]
owner: legacy
status: active
updated: 2026-06-08
---

# OneMind OS — Operations Dashboard

> **Updated 2026-06-08:** Added AgentMemory (:3113), SwarmClaw (:3456), Langfuse (:3000), LiveKit, Google Cloud, DigitalOcean, Proton Mail, NATS Fabric for Agents. Removed ~~Paperclip~~ *(removed 2026-06-08)*. HA token refreshed.

**Dashboard URL:** `http://localhost:3010` / `http://omos-mbp.tail717ba0.ts.net:3010`  
**Engine:** Homepage v1.13.1 (Next.js standalone, macOS LaunchAgent)  
**Config:** `~/homepage/config/` (git-versionable YAML)  
**LaunchAgent:** `com.onemind.homepage` — auto-starts on login  
**Health Monitor:** Uptime Kuma `:3001`

---

## Top Bar Widgets

| Widget | Description |
|--------|-------------|
| System Resources | omos-mbp CPU / RAM / Disk (5s refresh) |
| Date / Time | Short date + 24h clock |
| Weather | Greenwich, CT — OpenMeteo (imperial) |
| Search | DuckDuckGo (new tab) |

---

## Service Groups

### 🤖 AI Command

| Service | URL | Port | Description |
|---------|-----|------|-------------|
| Hermes WebUI | `localhost:8787` | 8787 | Chat & Task Interface |
| Hermes Home | `hermes-agent.nousresearch.com` | — | Official Site & Downloads |
| Hermes Dashboard | `localhost:9119`
| AgentMemory | `localhost:3113` | Persistent agent memory viewer (MCP) | | 9119 | Agent Monitor — Live TUI Sessions |
| Agent Web UI | `localhost:3300` | 3300 | Synadia Harness — Claude / Pi / OpenClaw |
| agenttap — The Stream | `omos-mbp.tail717ba0.ts.net:8088` | 8088 | Wireshark for AI agents — live NATS agent-SDK traffic visualizer |
| OpenClaw | `omos-mbp.tail717ba0.ts.net:3210` | 3210 | Device Gateway Agent — Controller |
| Claude | `claude.ai` | — | Anthropic AI assistant |
| ChatGPT | `chat.openai.com` | — | OpenAI assistant & API platform |
| Perplexity | `perplexity.ai` | — | AI-powered search & research |
| Grok | `grok.com` | — | xAI assistant |

### 🏠 Home & Farm

| Service | URL | Port | Description |
|---------|-----|------|-------------|
| Home Assistant | `omos-ha.tail717ba0.ts.net:8123` | 8123 | Smart Home Hub — Lights, Sensors, Automations |
| FarmBot | `github.com/FarmBot/farmbot_os` | — | Open-source precision farming robot OS |
| Obsidian Vault | `obsidian://open` | local | Knowledge Base — OneMind OS Vault |

### 🔌 NATS / Synadia

| Service | URL | Description |
|---------|-----|-------------|
| Synadia NGS | `app.ngs.global` | Managed NATS Super-Cluster — Global Edge |
| Synadia Cloud | `synadia.com` | Synadia — NATS Platform & Products |
| Synadia Agents | `github.com/synadia-ai/synadia-agents` | NATS-native Agent Framework |
| nats-io |
| NATS Fabric for Agents | `nats.io` | NATS connective technology hub | `github.com/nats-io` | Official NATS GitHub org |

### 🔧 Infrastructure

| Service | URL | Port | Description |
|---------|-----|------|-------------|
| NATS Bus | `localhost:8222` | 8222 | NATS v2.14 — Message Bus Monitor |
| NATS Dashboard | `localhost:4321` | 4321 | Visual Dashboard — Streams & Consumers |
| NetBird | `net.onemindos.dev` | — | WireGuard Mesh VPN — Peer Network Manager |
| Tailscale Admin | `login.tailscale.com/admin` | — | Tailnet Fleet — omos-* nodes + MagicDNS |
| AWS Console | `us-east-1.console.aws.amazon.com` | — | AWS — Cloud Infrastructure (us-east-1) |
| Azure Portal | `portal.azure.com` | — | Microsoft Cloud Portal |
| Google Cloud | `console.cloud.google.com` | GCP Console |
| DigitalOcean | `cloud.digitalocean.com` | Cloud VMs & Droplets |
| Cloudflare | `dash.cloudflare.com` | — | DNS, CDN, Zero Trust & Security |
| ownCloud OCIS | `localhost:9200` | 9200 | Self-Hosted Cloud — Files, Calendar, Contacts |
| Uptime Kuma | `localhost:3001` | 3001 | Service Monitor — Uptime & Alerts |
| Netdata | `localhost:19999` | 19999 | Real-Time System Metrics — CPU, Memory, Disk |
| LiveKit | `cloud.livekit.io` | Real-time voice/video infrastructure |
| Node-RED | `localhost:1880` | 1880 | Flow-Based Automation & IoT Wiring |

### 📡 Protocols

| Service | URL | Description |
|---------|-----|-------------|
| Zenoh | `zenoh.io` | Zero-overhead pub/sub for edge, cloud & embedded |
| Reticulum | `reticulum.network` | Cryptography-based networking for unreliable links |
| OpenMANET | `openmanet.net` | Open mobile ad-hoc mesh networking |

### 🔐 Security

| Service | URL | Description |
|---------|-----|-------------|
| Cloaked | `my.cloaked.com` | Privacy masks for email, phone & identity |
| 1Password | `start.1password.com` | Secrets, credentials & vault manager |

### 💼 Business

| Service | URL | Description |
|---------|-----|-------------|
| zeusdelacruz.com | `zeusdelacruz.com` | Personal Brand — Speaking, Services, Authority |
| weareonemind.com | `weareonemind.com` | OneMind OS — Business Site on GoHighLevel |
| GoHighLevel CRM | `app.gohighlevel.com` | CRM · Funnels · Automations · Pipelines |
| Geelark | `geelark.com` | Cloud phone & social media marketing automation |
| Skool Community | `skool.com` | Wealth Warriors — Community Platform |
| Skool Profile | `skool.com/@anthony-delacruz-9523` | Anthony Delacruz profile |

### 💰 Business Finance

| Service | URL | Description |
|---------|-----|-------------|
| Mercury Bank | `app.mercury.com/panorama` | Business Banking — Accounts & Transactions |
| Secret | `joinsecret.com` | Startup deals, SaaS credits & founder perks |
| FounderPass | `founderpass.com` | Founder discounts & startup tool credits |
| Lightning Network | `lightning.network` | L2 Bitcoin instant payments protocol |

### 🌍 Intelligence

| Service | URL | Port | Description |
|---------|-----|------|-------------|
| Crucix | `localhost:3117` | 3117 | OSINT Intelligence Terminal — 27 live feeds (satellites, flights, radiation, markets, conflict) |
| WorldWideView Marketplace | `marketplace.worldwideview.dev` | — | Intelligence data marketplace |
| WorldWideView (omos-cloud) | `100.111.189.116:3002` | 3002 | World Monitor — Real-Time Global Intelligence Dashboard (omos-cloud, Tailscale) |

### 📚 Resources: AI & Agents

| Service | URL | Description |
|---------|-----|-------------|
| AAIF | `aaif.io` | Agentic AI Foundation — AGENTS.md / MCP / agentgateway / Goose |
| Hermes Docs | `hermes-agent.nousresearch.com/docs` | Hermes Agent — Official Documentation |
| Hermes Atlas | `hermesatlas.com` | Community map — 100+ tools, skills & ecosystem |
| AgentSkills.io | `agentskills.io/home` | Open standard for reusable agent skills |
| Context7 Docs | `context7.com/docs/overview` | Live docs & code context for AI agents |
| Langfuse Docs | `langfuse.com/docs` | LLM observability & tracing documentation |

### 📚 Resources: NATS & Infra

| Service | URL | Description |
|---------|-----|-------------|
| NATS Docs | `docs.nats.io` | NATS.io — Official documentation |
| Composio Docs | `docs.composio.dev` | 2000+ integrations, MCP & tool API reference |
| ClickHouse Docs | `clickhouse.com/docs` | High-performance analytics DB documentation |
| Go Docs | `go.dev` | Official language site, docs & packages |

### 📚 Resources: Robotics & Edge

| Service | URL | Description |
|---------|-----|-------------|
| OpenMind | `openmind.com` | AI robotics platform (OM1) |
| Zenoh Use Cases | `zenoh.io/usecases/` | Use cases for edge, robotics, automotive & IoT |
| ROS Docs | `docs.ros.org` | Official docs, tutorials & API reference |
| ArduPilot Docs | `ardupilot.org/ardupilot/` | Full documentation for drones, planes & rovers |
| PX4 Docs | `docs.px4.io` | Official flight control firmware documentation |
| Gazebo Docs | `gazebosim.org/docs` | Getting started with robot simulation |

### 📚 Resources: Business

| Service | URL | Description |
|---------|-----|-------------|
| GoHighLevel Docs | `help.gohighlevel.com` | Help center, CRM & funnel docs |
| Skool Help | `help.skool.com` | Creator & community platform help docs |
| Meta Ads Docs | `facebook.com/business/learn` | Ads Manager learning & business resources |
| Google Ads Help | `support.google.com/google-ads` | Campaign setup, billing & optimization docs |

### 📚 Resources: Security

| Service | URL | Description |
|---------|-----|-------------|
| 1Password Docs | `developer.1password.com` | Developer docs, CLI, SDKs & Secrets Automation |
| Tailscale Docs | `tailscale.com/kb` | Knowledge base, mesh setup & ACL reference |
| Headscale Docs | `headscale.net/docs/` | Self-hosted Tailscale control plane docs |

### 🏥 Personal Health

| Service | URL | Description |
|---------|-----|-------------|
| HeadsUp Health | `headsuphealth.com` | Personal health data aggregator & dashboard |

### 🛠 Dev & Ops

| Service | URL | Description |
|---------|-----|-------------|
| GitHub — OneMind Org | `github.com/onemindos` | onemindos — Private Repos & Code |
| GitHub — Enterprise | `github.com/enterprises/onemind` | OneMind Enterprise — Billing, SSO, Policies |
| GitHub — Zeus | `github.com/Zeus-Delacruz` | Zeus Delacruz — Personal Profile & Public Repos |
| Todoist | `app.todoist.com` | Task Manager — Daily Ops & Projects |
| Fastmail |
| Proton Mail | `mail.proton.me` | zeus.vaults@pm.me (E2E encrypted) | `app.fastmail.com` | zeuslegacy@fastmail.com — Primary Inbox |
| Webflow | `webflow.com/dashboard` | zeusdelacruz.com Design |
| Composio | `dashboard.composio.dev` | Tool Integration Platform — Connected Apps |

### 📦 Repos

| Service | URL | Description |
|---------|-----|-------------|
| hermes-workspace | `github.com/outsourc-e/hermes-workspace` | Native web workspace for Hermes |
| hermes-webui | `github.com/nesquena/hermes-webui` | Lightweight web UI for Hermes Agent |
| agent-sessions | `github.com/jazzyalex/agent-sessions` | Agent session management |
| opensourceprojects.dev | `opensourceprojects.dev` | Curated open source project discovery |
| coral | `github.com/withcoral/coral` | Agent coordination framework |
| osiris | `github.com/simplifaisoul/osiris` | Agent framework |
| dimos | `github.com/dimensionalOS/dimos` | Humanoid robot OS |
| OM1 | `github.com/OpenMind/OM1` | OpenMind robot firmware & SDK |
| design.md | `github.com/google-labs-code/design.md` | Google's DESIGN.md token spec for coding agents |
| kokoro-tts | `github.com/nazdridoy/kokoro-tts` | Lightweight local TTS engine |
| FarmBot | `github.com/orgs/FarmBot/repositories` | FarmBot GitHub org — All repositories |
| cactus | `github.com/cactus-compute/cactus` | On-device LLM runtime for mobile & edge |
| Understand-Anything | `github.com/Lum1104/Understand-Anything` | Universal multimodal understanding framework |

### 🎙 AI Voice

| Service | URL | Description |
|---------|-----|-------------|
| kokoro-tts | `github.com/nazdridoy/kokoro-tts` | Lightweight local TTS |
| ElevenLabs | `elevenlabs.io` | AI voice generation & cloning |
| Cartesia | `cartesia.ai` | Ultra-low-latency voice AI (real-time TTS) |
| OpenAI TTS | `platform.openai.com/docs/guides/text-to-speech` | GPT-4o-audio voices |
| FFmpeg | `ffmpeg.org` | Universal audio/video processing & transcoding |
| Moshi | `github.com/kyutai-labs/moshi` | Real-time speech-to-speech voice AI |

### 🧰 AI Tools

| Service | URL | Description |
|---------|-----|-------------|
| Herdr | `herdr.dev` | Terminal-native agent runtime — tmux persistence + agent state + API |
| Herdr Docs | `herdr.dev/docs` | Herdr — Official documentation |
| Apify | `apify.com` | Web scraping, browser automation & AI agent actors |
| Context7 | `context7.com` | Live docs & code context for AI coding agents |
| MCP Market | `mcpmarket.com` | Marketplace for MCP servers & integrations |
| Composio | `composio.dev` | 2000+ tool integrations for AI agents |
| SkillsMP | `skillsmp.com` | Skills marketplace for AI agents |
| Firecrawl | `firecrawl.dev` | Web crawling & scraping API for AI agents |
| Crawlee Python | `crawlee.dev/python` | Python web scraping & browser automation |
| Screenpipe | `screenpipe.com` | Screen + audio capture for AI context & memory |
| Cactus | `cactuscompute.com` | On-device / edge LLM runtime |

### 📣 Social

| Service | URL | Description |
|---------|-----|-------------|
| YouTube | `youtube.com/@zeuslegacydelacruz` | @zeuslegacydelacruz |
| LinkedIn | `linkedin.com/in/zeuslegacydelacruz/` | Zeus Legacy Delacruz |
| X Console | `console.x.com` | X Developer Console — API, apps & credentials |
| X / Twitter | `x.com/zeuslegacy` | @zeuslegacy |
| Instagram | `instagram.com/zeuslegacydelacruz/` | @zeuslegacydelacruz |
| Meta Business | `business.facebook.com` | Ads, Pages, Assets |
| Google Ads | `ads.google.com` | Campaigns & Performance |

### 🤖 Robotics

| Service | URL | Description |
|---------|-----|-------------|
| Tinkr | `tnkr.ai/explore` | Robotics marketplace & community |
| OpenMind | `openmind.com` | AI robotics platform (OM1) |
| dimensionalOS/DIMOS | `github.com/dimensionalOS/dimos` | Humanoid robot OS |
| LeRobot | `huggingface.co/lerobot` | HuggingFace robotics models & datasets |
| VLA-JEPA | `huggingface.co/collections/lerobot/vla-jepa` | Vision-Language-Action + JEPA collection |
| ROS | `ros.org` | Robot Operating System |
| Foxglove | `foxglove.dev` | Robotics observability & ROS bag visualization |
| Gazebo | `gazebosim.org/home` | Open-source robot simulator |

### 🚁 Drones

| Service | URL | Description |
|---------|-----|-------------|
| ArduPilot | `ardupilot.org` | Open-source autopilot for drones, planes & rovers |
| PX4 | `px4.io` | Open-source flight control firmware |
| QGroundControl | `qgroundcontrol.com` | Ground station for ArduPilot & PX4 |
| Pixhawk | `pixhawk.org` | Open hardware autopilot standard |
| CubePilot | `cubepilot.com` | Cube Orange/Black flight controller hardware |
| IRLock | `irlock.com` | Precision landing system for drones |
| MAVROS | `github.com/mavlink/mavros` | MAVLink to ROS bridge |

### 🚀 Fleet

| Service | URL | Port | Description |
|---------|-----|------|-------------|
| Homepage Dashboard | `localhost:3010` | 3010 | This Dashboard — OneMind Ops |
| Legacy UI | `localhost:3000` | 3000 | Next.js Dev Interface |
| DIMOS | `github.com/dimensionalOS/dimos` | — | Humanoid Robot (offline until deployed) |
| OM1 | `github.com/OpenMind/OM1` | — | OpenMind Robot (offline until deployed) |
| Cognitum Seed Fleet | `omos-mbp.tail717ba0.ts.net:8080/devices` | 8080 | IoT Device Registry — Seed Nodes via NATS API |

---

## Known Issues / Gaps

| Item | Status | Notes |
|------|--------|-------|
| `Langfuse` | ✅ Running | `localhost:3004` (omos-cloud) — not yet on dashboard |
| `n8n` | ✅ Running | `localhost:5678` — not yet on dashboard |
| `SwarmClaw` | ✅ Running | `:3456/:3457` — not yet on dashboard |
| `agentmemory viewer` | ✅ Running | `localhost:3113` — not yet on dashboard |
| `ClickHouse` | ✅ Running | HTTP `localhost:8124` — not yet on dashboard |

---

## Architecture Notes

- **Auth boundary:** Tailscale mesh (`tail717ba0.ts.net`) — only enrolled nodes can reach dashboard
- **URL pattern:** `omos-<node>.tail717ba0.ts.net:<port>` for cross-node links
- **omos-cloud Tailscale IP:** `100.111.189.116` (us-east-1 Lightsail)
- **Config location:** `~/homepage/config/` — agent-editable, git-versioned
- **Auto-sync:** `com.onemind.homepage-sync` LaunchAgent watches config dir via `fswatch`
- **Network decision:** Tailscale cloud mesh — see `08 Research/network-mesh-comparison-2026.md`

---

## Adding Services

Edit `~/homepage/config/services.yaml`. Changes auto-sync via fswatch watcher. Format:

```yaml
- Section Name:
    - Service Name:
        href: http://hostname:port
        description: Short description
        icon: mdi-icon-name   # or si-brandname
        ping: http://internal-host:port   # health check URL
```

Icon sources: `mdi-*` = Material Design Icons, `si-*` = Simple Icons (brand logos)