---
type: system-status
category: operations
updated: 2026-05-26
tags:
  - status
  - health
  - monitoring
  - dashboard
---

# OneMind OS — System Status

**Last Updated:** 2026-05-26 (Monday)
**Overall Status:** 🟢 Operational

---

## Core Systems

| **System** | **Status** | **Location** | **Notes** |
|------------|-----------|--------------|-----------|
| **Mac M2 Max** | 🟢 Online | Zeuss-MBP.local | 24/7 clamshell, 32GB RAM |
| **Hermes Agent** | 🟢 Running | ~/.hermes/ | Default profile, Legacy AI |
| **NATS Server** | 🟢 Running | localhost:4222 | OneMind fabric |
| **Obsidian Vault** | 🟢 Synced | ~/Documents/onemind/ | Git backed |
| **Homepage Dashboard** | 🟢 Running | http://localhost:9119 | Hermes WebUI |

---

## Agent Fleet

| **Agent** | **Status** | **Platform** | **Role** |
|-----------|-----------|--------------|---------|
| **Legacy AI** | 🟢 Active | Hermes | Command agent — primary interface |
| **Claude Code** | 🟡 On-demand | Claude Code CLI | Deep coding, features, PRs |
| **Gemini CLI** | 🟡 On-demand | Gemini CLI | Coding, research, analysis |
| **Codex / OpenCode** | 🟡 On-demand | CLI | Fast coding, quick tasks |

---

## MCP Servers

| **Server** | **Status** | **Tools** | **Transport** |
|------------|-----------|-----------|---------------|
| **nats** | 🟢 Enabled | NATS request/reply | stdio |
| **cua-driver** | 🟢 Enabled | macOS desktop control | stdio |
| **todoist** | 🟢 Enabled | Tasks, projects, labels | stdio |
| **filesystem** | 🟢 Enabled | File CRUD | stdio |
| **context7** | 🟢 Enabled | Live library docs | stdio |
| **excalidraw** | 🟢 Enabled | Diagram canvas | stdio |
| **telegram** | 🟢 Enabled | Messaging | stdio |
| **fastmail** | 🟢 Enabled | Email (JMAP) | HTTP |

---

## Infrastructure

### Networking
| **Service** | **Status** | **Details** |
|-------------|-----------|-------------|
| **Netbird** | 🟢 Running | net.onemindos.dev, self-hosted VPN mesh |
| **Tailscale** | 🟡 Pending decom | Legacy, migrating to Netbird |

### Storage
| **Location** | **Usage** | **Purpose** |
|--------------|-----------|-------------|
| **~/Documents/onemind/** | ~500MB | Obsidian vault (OneMind OS knowledge base) |
| **~/.hermes/** | ~2GB | Hermes config, skills, sessions, cache |
| **~/.nats/** | ~100MB | NATS config and data |
| **~/homepage/** | ~50MB | Homepage dashboard Docker stack |

### Services (Docker)
| **Container** | **Status** | **Port** | **Purpose** |
|---------------|-----------|----------|-------------|
| **homepage** | 🟢 Running | 3010 | Dashboard UI |
| **netbird-management** | 🟢 Running | 8081 | Netbird admin |
| **netbird-signal** | 🟢 Running | 10000 | WebRTC signaling |

---

## AI Providers

| **Provider** | **Status** | **Models** | **Usage** |
|--------------|-----------|------------|-----------|
| **AWS Bedrock** | 🟢 Active | Claude Sonnet 4.5, 4.6 | Primary |
| **Azure OpenAI** | 🟢 Active | GPT-5.x, GPT-4.1, o4-mini | Fallback |
| **GitHub Copilot** | 🟢 Active | Claude Sonnet 4 | MCP client |
| **Anthropic Direct** | 🟡 Available | Claude 3.5+ | Emergency fallback |

---

## Integrations

| **Platform** | **Status** | **Access Method** | **Notes** |
|--------------|-----------|-------------------|-----------|
| **Todoist** | 🟢 Connected | MCP server | Tasks, projects, labels |
| **Fastmail** | 🟢 Connected | MCP server (HTTP) | Email via JMAP |
| **Home Assistant** | 🟢 Connected | NATS bridge | 185 entities cached |
| **Telegram** | 🟢 Connected | MCP server | Messaging |
| **1Password** | 🟢 Integrated | Safari autofill | Credential vault |
| **Composio** | 🟢 Connected | NATS worker | 2000+ SaaS actions |
| **GitHub** | 🟢 Connected | Composio / gh CLI | Repos, PRs, issues |

---

## Pending

### Hardware
- [ ] **Pi Zero 2 W** — IoT sensor nodes (temp, motion, doors)
- [ ] **Pi 5 + AI HAT 2** — Edge vision

### Software
- [ ] **OBS Studio + NATS** — Video automation (deferred)
- [ ] **YOLOv8** — Object detection / LPR (deferred until Pi 5)
- [ ] **Resource monitoring** — CPU, RAM, disk (Netdata or Prometheus)
- [ ] **Auto-backup** — Hermes config + Homepage config

---

## Recent Changes

### 2026-05-26
- ✅ Vault root cleanup — 17 skills pruned, stale files deleted/moved
- ✅ onemind.md + codex.md renamed (lowercase, developer convention)
- ✅ DOMAIN-REGISTRY.md moved to SO/02 Infrastructure/network/dns/
- ✅ SCHEMA.md moved to SO/00 Framework/reference/
- ✅ MASTER-INDEX.md deleted (stale, will rebuild on cron)
- ✅ CHANGELOG.md deleted (documented dead VPS system)
- ✅ ONE-MIND-FRAMEWORK.md + onemind.md merged into onemind.md (full line-for-line)
- ✅ ROADS-Community-Vision.md moved to GE/80 Innovation/ideas/
- ✅ agents.md, claude.md, gemini.md, README.md all updated to current stack
- ✅ Peekaboo MCP removed

### 2026-05-24
- ✅ Installed Excalidraw MCP Server
- ✅ Composio NATS worker deployed (v3.0.0)
- ✅ Home Assistant migrated to NATS bridge

---

## Known Issues

### Active
- ⚠️ **Tailscale pending decom** — Netbird is the target, migration not complete
- ⚠️ **MASTER-INDEX stale/deleted** — Needs cron rebuild once vault settles

### Resolved
- ✅ Obsidian Git SSH error — Fixed (known_hosts refresh)
- ✅ Hermes compaction threshold — Raised to 0.8
- ✅ Peekaboo MCP conflicts — Removed

---

## Monitoring & Alerts

### Automated Checks
- **Hermes Gateway:** Auto-continues after 900s timeout warning
- **NATS Heartbeats:** agents.hb.* every 30s
- **Docker Health:** Homepage checks container status

### Manual Checks (Daily)
- [ ] Hermes dashboard (http://localhost:9119)
- [ ] Obsidian vault Git sync status
- [ ] NATS server logs (`nats server check`)

---

## Backup Status

| **System** | **Frequency** | **Destination** | **Status** |
|------------|---------------|-----------------|------------|
| **Obsidian Vault** | Real-time (Git) | github.com/OneMind-OS/One-Mind | 🟢 Active |
| **Hermes Config** | Manual | Git tracked | 🟡 Pending auto |
| **Homepage Config** | Manual | Git tracked | 🟡 Pending auto |

### Recovery
- **Vault:** `git clone git@github.com:OneMind-OS/One-Mind.git`
- **Hermes:** Restore from `~/.hermes/`
- **NATS:** Config in `~/.nats/`

---

## Emergency

### Critical Systems Down
1. **Zeus** — Direct hardware access
2. **Legacy AI** — Restart via `hermes gateway restart` or WebUI reload

### Service Escalation
- **NATS:** `~/.nats/nats-server.log`
- **Hermes:** `hermes logs`
- **Docker:** `docker logs <container>`
- **Git Sync:** Check `~/.ssh/known_hosts` + GitHub SSH keys

---

## Related Files
- **NATS Infrastructure:** `00-24 SO/02 Infrastructure (The Body)/NATS-INFRASTRUCTURE.md`
- **OBS Integration Plan:** `00-24 SO/02 Infrastructure (The Body)/obs-nats-integration-plan.md`
- **Hermes Config:** `~/.hermes/config.yaml`

---

**Hermes Dashboard:** http://localhost:9119
**Overall Health:** 🟢 All critical systems operational ★
