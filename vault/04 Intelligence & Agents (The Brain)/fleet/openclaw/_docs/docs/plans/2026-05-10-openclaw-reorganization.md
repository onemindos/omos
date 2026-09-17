# OpenClaw Vault Reorganization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restructure `openclaw (OpenClaw Platform)/` to mirror the official docs.openclaw.ai hierarchy so every piece of documentation, config, and reference has a predictable home.

**Architecture:** The official docs has 17 top-level sections (start, gateway, channels, nodes, platforms, concepts, automation, tools, plugins, clawhub, providers, install, cli, reference, security, help, web). We map our existing content into those exact sections. Nothing is deleted — everything is moved or kept in place.

**Tech Stack:** Bash (mv/mkdir), Markdown files, Obsidian vault conventions

---

## Section Map (Official Docs → Our Vault)

| New Folder | Official Docs | Source Content |
|---|---|---|
| `start/` | start/ | guides/OPENCLAW-COURSE.md, guides/OPENCLAW-PRACTICAL-GUIDE.md, guides/OpenClaw-*.md |
| `gateway/` | gateway/ | deploy/docs/ (DEPLOYMENT, NODES, VPS-GATEWAY, VPS-SERVICES, DISASTER-RECOVERY, etc.), openclaw-reference/OPENCLAW-ARCHITECTURE.md, OPENCLAW-CONFIG-REFERENCE.md |
| `channels/` | channels/ | Channels/CHANNELS.md, interfaces/discord.md, interfaces/telegram.md |
| `nodes/` | nodes/ | interfaces/ios-app.md, interfaces/conversal-voice-architecture.md, interfaces/TALK-MODE-CONFIG.md |
| `platforms/` | platforms/ | interfaces/macos-app.md, interfaces/onemind-app.md, openclaw-reference/openclaw-peekaboo-mac-automation.md |
| `concepts/` | concepts/ | openclaw-reference/OPENCLAW_DEEP_DIVE.md, OPENCLAW-MASTER-REFERENCE.md, ACTIVE-MEMORY-PLUGIN.md, NATIVE-MEMORY-ARCHITECTURE.md, OPENCLAW-CONCURRENCY.md, MULTI-USER.md, DEVICE-ARCHITECTURE.md, openclaw-reference/OPENCLAW_CAPABILITIES.md |
| `automation/` | automation/ | _system/openclaw/hooks/ (all), deploy/scripts/setup-cron-reports.sh |
| `tools/` | tools/ | OpenClaw-Skills-Catalog/, OpenClaw-Usecases/, openclaw-reference/The complete map*.md |
| `plugins/` | plugins/ | openclaw-reference/OMCOS_PLUGIN_CATALOG.md, interfaces/swarmclaw-web.md |
| `clawhub/` | clawhub/ | (new — placeholder for skill publishing docs) |
| `providers/` | providers/ | openclaw-reference/PROVIDERS.md, OpenClaw-Hosting-Costs/AI-Model-API-Costs.md |
| `install/` | install/ | deploy/ (renamed + restructured) |
| `cli/` | cli/ | openclaw-reference/cli.md |
| `reference/` | reference/ | openclaw-reference/OPENCLAW_INTEGRATION_ANALYSIS.md, openclaw-reference/The complete OpenClaw ecosystem inventory.md, openclaw-reference/QUICK-REFERENCE.md, openclaw-reference/openclaw-compass-artifact.md, planning/ |
| `security/` | security/ | OpenClaw-Hosting-Costs/Security-CVEs-Hardening.md, openclaw-reference/Discord × OpenClaw*.md |
| `help/` | help/ | deploy/docs/PRODUCTION_BEST_PRACTICES.md, deploy/docs/OPENCLAW-GATEWAY-DEPLOYMENT.md, deploy/docs/OPENCLAW-GATEWAY-HANDOFF.md, deploy/docs/OPENCLAW-LIVE-STATE.md, guides/OPENCLAW-DEEP-DIVE.md |
| `web/` | web/ | interfaces/gateway-api.md, guides/GATEWAY-DASHBOARD-ACCESS.md, deploy/docs/OPENCLAW-OTEL-DATADOG.md |
| `costs/` | (custom) | OpenClaw-Hosting-Costs/Hosting-Providers-Comparison.md |

---

### Task 1: Create all new section folders

**Files:**
- Create: `start/`, `gateway/`, `channels/`, `nodes/`, `platforms/`, `concepts/`, `automation/`, `tools/`, `plugins/`, `clawhub/`, `providers/`, `install/`, `cli/`, `reference/`, `security/`, `help/`, `web/`, `costs/`

**Step 1: Run mkdir for all sections**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"
mkdir -p "$BASE/start"
mkdir -p "$BASE/gateway"
mkdir -p "$BASE/channels"
mkdir -p "$BASE/nodes"
mkdir -p "$BASE/platforms"
mkdir -p "$BASE/concepts"
mkdir -p "$BASE/automation"
mkdir -p "$BASE/tools/skills-catalog"
mkdir -p "$BASE/tools/usecases"
mkdir -p "$BASE/plugins"
mkdir -p "$BASE/clawhub"
mkdir -p "$BASE/providers"
mkdir -p "$BASE/install"
mkdir -p "$BASE/cli"
mkdir -p "$BASE/reference"
mkdir -p "$BASE/security"
mkdir -p "$BASE/help"
mkdir -p "$BASE/web"
mkdir -p "$BASE/costs"
```

**Step 2: Verify folders created**

```bash
ls "$BASE/" | sort
```

---

### Task 2: Move openclaw-reference/ files to correct sections

**Files:**
- Modify: `openclaw-reference (OpenClaw Reference)/` → distribute to sections

**Step 1: Move architecture/concepts files → concepts/**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"
SRC="$BASE/openclaw-reference (OpenClaw Reference)"
mv "$SRC/OPENCLAW_DEEP_DIVE.md" "$BASE/concepts/"
mv "$SRC/OPENCLAW-MASTER-REFERENCE.md" "$BASE/concepts/"
mv "$SRC/ACTIVE-MEMORY-PLUGIN.md" "$BASE/concepts/"
mv "$SRC/NATIVE-MEMORY-ARCHITECTURE.md" "$BASE/concepts/"
mv "$SRC/OPENCLAW-CONCURRENCY.md" "$BASE/concepts/"
mv "$SRC/MULTI-USER.md" "$BASE/concepts/"
mv "$SRC/DEVICE-ARCHITECTURE.md" "$BASE/concepts/"
mv "$SRC/OPENCLAW_CAPABILITIES.md" "$BASE/concepts/"
mv "$SRC/openclaw-deep-dive.md" "$BASE/concepts/"
mv "$SRC/OPENCLAW-INTERFACES.md" "$BASE/concepts/"
```

**Step 2: Move gateway/config files → gateway/**

```bash
mv "$SRC/OPENCLAW-ARCHITECTURE.md" "$BASE/gateway/"
mv "$SRC/OPENCLAW_ARCHITECTURE_REFERENCE.md" "$BASE/gateway/"
mv "$SRC/OPENCLAW-CONFIG-REFERENCE.md" "$BASE/gateway/"
```

**Step 3: Move providers → providers/**

```bash
mv "$SRC/PROVIDERS.md" "$BASE/providers/"
```

**Step 4: Move CLI → cli/**

```bash
mv "$SRC/cli.md" "$BASE/cli/"
```

**Step 5: Move plugins/tools → plugins/**

```bash
mv "$SRC/OMCOS_PLUGIN_CATALOG.md" "$BASE/plugins/"
mv "$SRC/OPENCLAW_INTEGRATION_ANALYSIS.md" "$BASE/plugins/"
```

**Step 6: Move macOS/platform → platforms/**

```bash
mv "$SRC/openclaw-peekaboo-mac-automation.md" "$BASE/platforms/"
```

**Step 7: Move reference docs → reference/**

```bash
mv "$SRC/QUICK-REFERENCE.md" "$BASE/reference/"
mv "$SRC/openclaw-compass-artifact.md" "$BASE/reference/"
mv "$SRC/The complete OpenClaw ecosystem inventory.md" "$BASE/reference/"
mv "$SRC/The complete map of every tool in the OpenClaw ecosystem.md" "$BASE/reference/"
mv "$SRC/what-openclaw-is-weaker-on-and-the-fix.md" "$BASE/reference/"
mv "$SRC/legacy-openclaw-diagrams.html" "$BASE/reference/"
```

**Step 8: Move remaining to security/**

```bash
mv "$SRC/Discord × OpenClaw - Complete Technical Reference for One Mind OS.md" "$BASE/security/"
```

**Step 9: Verify source is now empty**

```bash
ls "$SRC/"
```

---

### Task 3: Move guides/ and interfaces/ files

**Step 1: Move guides → start/ and help/**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"
SRC="$BASE/guides (Guides)"
mv "$SRC/OPENCLAW-COURSE.md" "$BASE/start/"
mv "$SRC/OPENCLAW-PRACTICAL-GUIDE.md" "$BASE/start/"
mv "$SRC/OpenClaw - the 220k-star AI agent redefining personal automation.md" "$BASE/start/"
mv "$SRC/OPENCLAW-DEEP-DIVE.md" "$BASE/help/"
mv "$SRC/GATEWAY-DASHBOARD-ACCESS.md" "$BASE/web/"
```

**Step 2: Move interfaces → channels / nodes / platforms / web**

```bash
SRC="$BASE/interfaces (Interfaces)"
# Channels
mv "$SRC/discord.md" "$BASE/channels/"
mv "$SRC/telegram.md" "$BASE/channels/"
# Nodes / voice
mv "$SRC/ios-app.md" "$BASE/nodes/"
mv "$SRC/conversal-voice-architecture.md" "$BASE/nodes/"
mv "$SRC/TALK-MODE-CONFIG.md" "$BASE/nodes/"
mv "$SRC/twilio.md" "$BASE/nodes/"
# Platforms
mv "$SRC/macos-app.md" "$BASE/platforms/"
mv "$SRC/onemind-app.md" "$BASE/platforms/"
# Web
mv "$SRC/gateway-api.md" "$BASE/web/"
mv "$SRC/swarmclaw-web.md" "$BASE/web/"
```

**Step 3: Verify both source dirs are empty**

```bash
ls "$BASE/guides (Guides)/" 2>/dev/null || echo "EMPTY"
ls "$BASE/interfaces (Interfaces)/" 2>/dev/null || echo "EMPTY"
```

---

### Task 4: Reorganize deploy/ into install/ structure

**Step 1: Copy deploy/docs/ files to gateway/ and help/**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"
SRC="$BASE/deploy (Deployment)/docs (Deployment Docs)"

# Gateway-focused docs
mv "$SRC/DEPLOYMENT.md" "$BASE/gateway/"
mv "$SRC/NODES.md" "$BASE/gateway/"
mv "$SRC/VPS-GATEWAY.md" "$BASE/gateway/"
mv "$SRC/VPS-SERVICES.md" "$BASE/gateway/"
mv "$SRC/DISASTER-RECOVERY.md" "$BASE/gateway/"
mv "$SRC/OPENCLAW-OTEL-DATADOG.md" "$BASE/web/"

# Help/runbook docs
mv "$SRC/OPENCLAW-GATEWAY-DEPLOYMENT.md" "$BASE/help/"
mv "$SRC/OPENCLAW-GATEWAY-HANDOFF.md" "$BASE/help/"
mv "$SRC/OPENCLAW-LIVE-STATE.md" "$BASE/help/"
mv "$SRC/PRODUCTION_BEST_PRACTICES.md" "$BASE/help/"
mv "$SRC/OC VPS STATUS.md" "$BASE/help/"
```

**Step 2: Rename General (Overview) → config-templates inside deploy/**

```bash
mv "$BASE/deploy (Deployment)/General (Overview)" "$BASE/deploy (Deployment)/config-templates"
```

**Step 3: Move patches into deploy/scripts/**

```bash
mv "$BASE/patches (Patches)/patch-universal-voice.sh" "$BASE/deploy (Deployment)/scripts (Deploy Scripts)/"
```

**Step 4: Move deploy/docs (now empty) — remove it**

```bash
ls "$BASE/deploy (Deployment)/docs (Deployment Docs)/" 2>/dev/null || echo "EMPTY"
```

**Step 5: Rename deploy folder to install**

```bash
mv "$BASE/deploy (Deployment)" "$BASE/install (Install & Deploy)"
```

---

### Task 5: Consolidate remaining top-level folders

**Step 1: Move OpenClaw-Hosting-Costs → costs/ and providers/**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"
mv "$BASE/OpenClaw-Hosting-Costs (Hosting Costs)/Hosting-Providers-Comparison.md" "$BASE/costs/"
mv "$BASE/OpenClaw-Hosting-Costs (Hosting Costs)/AI-Model-API-Costs.md" "$BASE/providers/"
mv "$BASE/OpenClaw-Hosting-Costs (Hosting Costs)/Security-CVEs-Hardening.md" "$BASE/security/"
```

**Step 2: Move Skills Catalog → tools/skills-catalog/**

```bash
mv "$BASE/OpenClaw-Skills-Catalog (Skills Catalog)/"* "$BASE/tools/skills-catalog/"
```

**Step 3: Move Usecases → tools/usecases/**

```bash
mv "$BASE/OpenClaw-Usecases (Use Cases)/"* "$BASE/tools/usecases/"
```

**Step 4: Move planning/ → reference/planning/**

```bash
mkdir -p "$BASE/reference/planning"
mv "$BASE/planning (Planning)/"* "$BASE/reference/planning/"
```

**Step 5: Move Channels/CHANNELS.md → channels/**

```bash
mv "$BASE/Channels/CHANNELS.md" "$BASE/channels/"
```

---

### Task 6: Write README.md for every section + master README

One README per section — short, with: Purpose, What's Here, Link to official docs section.

**Master README.md** (`openclaw (OpenClaw Platform)/README.md`):

```markdown
---
domain: 02
name: OpenClaw Platform
agent: Grid
status: active
updated: 2026-05-10
---

# OpenClaw Platform

OpenClaw is the AI agent gateway running all 19 One Mind agents across omos-core and omos-sentinel.

This folder mirrors the structure of [docs.openclaw.ai](https://docs.openclaw.ai) so every topic maps to a known official docs section.

## Sections

| Folder | What It Covers | Official Docs |
|---|---|---|
| [start/](start/) | Setup, onboarding, lore | [docs.openclaw.ai/start](https://docs.openclaw.ai/start/getting-started.md) |
| [gateway/](gateway/) | Configuration, health, logging, security | [docs.openclaw.ai/gateway](https://docs.openclaw.ai/gateway/index.md) |
| [channels/](channels/) | Telegram, Discord, WhatsApp, etc. | [docs.openclaw.ai/channels](https://docs.openclaw.ai/channels/index.md) |
| [nodes/](nodes/) | iOS/Android, voice, camera, media | [docs.openclaw.ai/nodes](https://docs.openclaw.ai/nodes/index.md) |
| [platforms/](platforms/) | macOS, iOS, Android, Linux apps | [docs.openclaw.ai/platforms](https://docs.openclaw.ai/platforms/index.md) |
| [concepts/](concepts/) | Memory, soul, sessions, architecture | [docs.openclaw.ai/concepts](https://docs.openclaw.ai/concepts/architecture.md) |
| [automation/](automation/) | Cron, hooks, taskflow, standing orders | [docs.openclaw.ai/automation](https://docs.openclaw.ai/automation/index.md) |
| [tools/](tools/) | Skills, browser, search, exec, use cases | [docs.openclaw.ai/tools](https://docs.openclaw.ai/tools/index.md) |
| [plugins/](plugins/) | Plugin SDK, manifests, community plugins | [docs.openclaw.ai/plugins](https://docs.openclaw.ai/plugins/architecture.md) |
| [clawhub/](clawhub/) | Publishing and installing skills | [docs.openclaw.ai/clawhub](https://docs.openclaw.ai/clawhub/index.md) |
| [providers/](providers/) | AI model providers and costs | [docs.openclaw.ai/providers](https://docs.openclaw.ai/providers/index.md) |
| [install/](install/) | Deployment: VPS, Docker, scripts, systemd | [docs.openclaw.ai/install](https://docs.openclaw.ai/install/index.md) |
| [cli/](cli/) | CLI reference | [docs.openclaw.ai/cli](https://docs.openclaw.ai/cli/index.md) |
| [reference/](reference/) | Master reference, templates, ecosystem inventory | [docs.openclaw.ai/reference](https://docs.openclaw.ai/reference/credits.md) |
| [security/](security/) | Threat model, CVEs, hardening | [docs.openclaw.ai/security](https://docs.openclaw.ai/gateway/security/index.md) |
| [help/](help/) | Troubleshooting, runbooks, live state | [docs.openclaw.ai/help](https://docs.openclaw.ai/help/index.md) |
| [web/](web/) | Dashboard, WebChat, TUI, API | [docs.openclaw.ai/web](https://docs.openclaw.ai/web/index.md) |
| [costs/](costs/) | Hosting provider cost comparison | — |
| [workspaces/](workspaces/) | Per-agent workspace files | — |

## Our Setup

- **Platform version:** OpenClaw 2026.5.x
- **Nodes:** omos-core (primary), omos-sentinel (failover)
- **Provider:** GitHub Copilot → claude-sonnet-4.6
- **Agents:** 19 (see [workspaces/](workspaces/))
```

**Per-section READMEs** — each gets a short `README.md` with Purpose + What's Here + Official Docs link.

---

### Task 7: Remove now-empty old folders

**Only remove a folder after confirming it is empty.**

```bash
BASE="/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)"

# Check each before removing
for dir in \
  "openclaw-reference (OpenClaw Reference)" \
  "guides (Guides)" \
  "interfaces (Interfaces)" \
  "OpenClaw-Hosting-Costs (Hosting Costs)" \
  "OpenClaw-Skills-Catalog (Skills Catalog)" \
  "OpenClaw-Usecases (Use Cases)" \
  "planning (Planning)" \
  "patches (Patches)" \
  "Channels"
do
  count=$(ls "$BASE/$dir/" 2>/dev/null | grep -v '.DS_Store' | wc -l)
  echo "$dir: $count files remaining"
done
```

Only rmdir the ones showing 0 files.

---

### Task 8: Commit

```bash
cd /Users/zeuslegacy/Documents/onemind
git add -A "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/"
git commit -m "refactor(openclaw): restructure platform folder to mirror docs.openclaw.ai hierarchy

- 17 sections now match official docs structure (start, gateway, channels, nodes,
  platforms, concepts, automation, tools, plugins, clawhub, providers, install,
  cli, reference, security, help, web)
- openclaw-reference/ distributed to proper concept/gateway/cli/reference sections
- guides/ split into start/ and help/
- interfaces/ split into channels/ nodes/ platforms/ web/
- deploy/ renamed to install/, internal docs moved to gateway/ and help/
- OpenClaw-Hosting-Costs/ merged into providers/ costs/ security/
- Skills-Catalog + Usecases moved into tools/
- planning/ moved into reference/planning/
- Added master README.md with full section index"
```
