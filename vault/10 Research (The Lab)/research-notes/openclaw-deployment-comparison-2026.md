---
title: "OpenClaw Deployment Methods — Comparison"
type: reference
status: active
created: 2026-04-06
owner: zeus
tags: [type/reference, domain/so, topic/openclaw, topic/deployment, topic/infrastructure]
ai_summary: "Comprehensive comparison of all OpenClaw deployment methods — npm, Docker, Podman, Nix, from source, VPS guides — with recommendation for OneMind full-freedom use case"
---

# OpenClaw Deployment Methods — Full Comparison

> Which deployment gives Zeus full freedom to run the complete OneMind stack?

---

## The Comparison Table

| Method | Setup Time | Full Control | Custom Skills | Memory-Wiki | Dreaming | Multi-Channel | Fork-Friendly | Survives Reboot | Update Path | Best For |
|--------|-----------|-------------|---------------|-------------|----------|---------------|---------------|-----------------|-------------|----------|
| **npm install -g** | 5 min | ✅ Full | ✅ Workspace skills | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ launchd/systemd | `npm update -g openclaw` | Mac daily driver |
| **Installer script** | 5 min | ✅ Full | ✅ Workspace skills | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ launchd/systemd | `openclaw update` | First install on any OS |
| **From source (pnpm)** | 15 min | ✅✅ Maximum | ✅ + source patches | ✅ + can modify | ✅ + can modify | ✅ All 28 | ✅✅ Best | ✅ manual | `git pull && pnpm build` | **Fork development** |
| **Docker (local)** | 20 min | ✅ Full | ✅ Workspace mount | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ `restart: unless-stopped` | `docker compose build` | Isolated testing |
| **Docker (VPS)** | 30 min | ✅ Full | ✅ Workspace mount | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ `restart: unless-stopped` | `git pull && docker compose build && up -d` | **24/7 remote server** |
| **Podman (rootless)** | 25 min | ✅ Full | ✅ Workspace mount | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ Quadlet/systemd | Rebuild image | Security-conscious Linux |
| **Nix (Home Manager)** | 30 min | ⚠️ Declarative | ✅ Workspace skills | ✅ Extension | ✅ Built-in | ✅ All 28 | ⚠️ Config-driven | ✅ launchd/systemd | `home-manager switch` | Reproducible/rollback |
| **Hetzner VPS guide** | 30 min | ✅ Full | ✅ Baked + workspace | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ Docker + systemd | `git pull && rebuild` | **Cheapest 24/7 (~$5/mo)** |
| **OpenClaw.app (macOS)** | 2 min | ⚠️ GUI-driven | ✅ Via Gateway | ✅ Extension | ✅ Built-in | ✅ All 28 | ❌ No fork | ✅ launchd | App auto-update | Mac menu bar convenience |
| **ClawDock** | 15 min | ✅ Full | ✅ Workspace mount | ✅ Extension | ✅ Built-in | ✅ All 28 | ✅ | ✅ Docker | Shell helpers | Docker power users |

---

## Key: What "Full Freedom" Means for OneMind

For the complete OneMind vision (see [[onemind-openclaw-vision-2027]]), the deployment MUST support:

| Requirement | Why | Which Methods Support It |
|-------------|-----|------------------------|
| **Custom fork modifications** | OneMind Gateway needs domain-aware routing, codex bridge, CODEX method automation | From source, Docker (with custom Dockerfile), VPS |
| **Memory-Wiki → Codex bridge** | Agents write back to the vault as Obsidian-native pages | All methods (it's an extension) |
| **Dreaming cycles (cron)** | 3 AM consolidation requires always-on process | VPS Docker, npm+systemd, Nix — NOT Mac (sleeps) |
| **Multi-agent with domain isolation** | 4+ agents (Legacy, Coach, Eden, Empire) each with own workspace | All methods support multi-agent config |
| **All 28 messaging channels** | WhatsApp, Telegram, Slack, Discord, Signal, iMessage, etc. | All methods |
| **Custom skills from ClawHub** | self-improving-agent, ontology, etc. | All methods (workspace skills dir) |
| **Sandbox for untrusted sessions** | Group chats run in Docker sandboxes | Docker, Podman, VPS — NOT bare npm |
| **24/7 uptime** | Phone queries at any hour, Dreaming runs at 3 AM | **VPS only** (Mac sleeps, laptop lids close) |
| **Device nodes (glasses, phone)** | VisionClaw iOS, Meta glasses, Android connect to Gateway | All methods (nodes connect via WS) |
| **Tailscale mesh access** | Access Gateway from anywhere on tailnet | All methods |

---

## Recommendation for OneMind

### The Two-Node Architecture

You don't pick ONE deployment. You run TWO:

```
┌─────────────────────────────────┐     ┌─────────────────────────────────┐
│       NODE 1: VPS (24/7)        │     │      NODE 2: Mac (when open)    │
│                                 │     │                                 │
│  Deployment: Docker (from fork) │     │  Deployment: From source (fork) │
│                                 │     │                                 │
│  ✅ Always-on Gateway           │     │  ✅ Development & testing        │
│  ✅ Dreaming (3 AM cron)        │     │  ✅ Obsidian editing             │
│  ✅ All messaging channels      │     │  ✅ Local Gateway (dev profile)  │
│  ✅ Memory-Wiki compilation     │     │  ✅ Fork modifications           │
│  ✅ Agent sandbox for groups    │     │  ✅ Quick iteration              │
│  ✅ Telegram/WhatsApp gateway   │     │                                 │
│                                 │     │  macOS app optional (menu bar)  │
│  Phone/Glasses connect here     │     │  Connects to VPS via Tailscale  │
│  via Tailscale mesh             │     │  when not running local Gateway │
└─────────────────────────────────┘     └─────────────────────────────────┘
         ↑                                         ↑
    104.248.4.212                              Mac (local)
    (your existing DO droplet)
```

### Why This Is the Right Answer

1. **VPS = production.** Always on. Phone hits it at 3 AM. Dreaming runs. Channels stay connected. Docker gives you isolation + easy updates.

2. **Mac = development.** Fork modifications, testing new skills, writing in Obsidian. Run a `--dev` profile Gateway on a different port (19001) for testing without touching production.

3. **Your existing VPS already runs OpenClaw.** You have NATS, PostgreSQL+pgvector, n8n on 104.248.4.212. Docker Compose just adds another service next to them.

4. **Tailscale connects everything.** Mac, VPS, phone — all on the same mesh. No port forwarding, no public exposure.

### VPS Deployment: Docker From Fork (Recommended)

```bash
# On VPS (104.248.4.212)
cd /opt/onemind
git clone https://github.com/OneMind-OS/onemind-gateway.git
cd onemind-gateway

# Custom Dockerfile with OneMind modifications
docker compose build
docker compose up -d

# Persist state
mkdir -p /root/.openclaw/workspace
chown -R 1000:1000 /root/.openclaw
```

Key config (`openclaw.json`):
```json
{
  "agent": {
    "model": "azure-openai/gpt-5.2"
  },
  "gateway": {
    "port": 18789,
    "bind": "lan",
    "tailscale": {
      "mode": "serve"
    }
  },
  "channels": {
    "telegram": { "botToken": "..." },
    "whatsapp": { "allowFrom": ["..."] }
  }
}
```

### Mac Deployment: From Source (Development)

```bash
# On Mac
cd ~/Documents
git clone https://github.com/OneMind-OS/onemind-gateway.git
cd onemind-gateway

pnpm install && pnpm ui:build && pnpm build

# Dev mode (port 19001, isolated state)
pnpm openclaw --dev gateway --allow-unconfigured

# Or production-like local
pnpm openclaw gateway --port 18789
```

---

## Update Workflow

### VPS (Production)

```bash
ssh root@104.248.4.212
cd /opt/onemind/onemind-gateway

# Cherry-pick from upstream
git remote add upstream https://github.com/openclaw/openclaw.git
git fetch upstream
git cherry-pick <commit>  # specific fixes only

# Rebuild
docker compose build
docker compose up -d
```

### Mac (Development)

```bash
cd ~/Documents/onemind-gateway
git fetch upstream
git cherry-pick <commit>
pnpm install && pnpm build
# Test locally before pushing to VPS
```

---

## What About NanoClaw?

NanoClaw (14K stars) is a **lightweight** OpenClaw runtime — fewer features, smaller footprint. It's designed for edge devices and resource-constrained environments.

**Verdict for OneMind:** Skip it. You want the full Gateway with Memory-Wiki, Dreaming, multi-agent, sandbox, and all 28 channels. NanoClaw drops too many features. It might be relevant later for dedicated hardware nodes (e.g., a robot that only needs basic agent runtime), but not for the main brain.

---

## Cost Comparison

| Deployment | Monthly Cost | Notes |
|-----------|-------------|-------|
| **Your existing DO VPS** | Already paying | Just add OpenClaw container |
| **Hetzner CX22** | ~$5/mo | 2 CPU, 4 GB RAM — enough for Gateway |
| **Hetzner CAX11 (ARM)** | ~$4/mo | ARM, slightly cheaper |
| **Mac (local only)** | $0 | But no 24/7 uptime |
| **Nix/npm local** | $0 | Same limitation |

**Bottom line:** You're already paying for the VPS. Adding OpenClaw Gateway as a Docker service costs $0 additional.

---

*See also: [[hermes-vs-agentzero-vs-openclaw-2026]], [[onemind-openclaw-vision-2027]], [[ai-agent-framework-comparison-2026]]*
