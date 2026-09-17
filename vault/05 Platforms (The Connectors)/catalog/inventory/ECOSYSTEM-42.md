# OneMind OS — Ecosystem 42: Complete Integration Catalog

> Every addition to OneMind OS from the OpenClaw ecosystem, with full specs, URLs, install instructions, and architectural connections.

**Last updated**: 2026-03-04
**Source**: OpenClaw ecosystem inventory research (see INVENTORY.md)

### Status Summary

| Status | Count | Items |
|--------|-------|-------|
| ✅ Done | 5 | #5, #6, #10, #13, #14 |
| 🔧 Installing | 3 | #19, #22, #38 |
| ⏳ Partial | 1 | #16 |
| 📋 Next Sprint | 6 | #1, #4, #8, #11, #15, #21 |
| 🔜 Deferred | 27 | Everything else |

---

## TIER 0 — Do First (minutes, instant impact)

---

### #1 — clawsec (Security Suite)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/prompt-security/clawsec |
| **Type** | Skill (8 skills in one package) |
| **Install** | `openclaw skills install clawsec` |
| **Effort** | 10 minutes |
| **Impact** | Critical — we have ZERO security tooling |
| **Status** | 📋 Next sprint — install first before any new skills |
| **Pillar** | UI (Innovative Technology) |

**What it does**: Eight security skills in one package:

1. **NVD Advisory Feed** — pulls real-time CVE alerts. Legacy knows about new vulnerabilities the day they drop.
2. **ClawHub Reputation Checker** — scores any skill before install. Flags the ~900 malicious skills documented by Bitdefender.
3. **Audit Watchdog** — monitors all tool calls for suspicious patterns: rate anomalies, unexpected file access, data exfiltration attempts.
4. **Soul-Guardian** — integrity hashing on SOUL.md, AGENTS.md, USER.md. Alerts if any file is tampered with.
5. **Attack Surface Scanner** — maps exposed endpoints and permissions.
6. **Prompt Injection Detector** — flags attempts to override system prompt.
7. **Permission Auditor** — reviews what each installed skill can access.
8. **Outbound Traffic Monitor** — logs all network requests skills make.

**Why we need it**: Kaspersky found 512 vulnerabilities across the ecosystem (8 critical). Cisco documented ~900 malicious skills. Bitsight scanned 42,900+ exposed instances. Five CVEs assigned and patched. We're running wide open without this.

**OneMind connection**: Protects SOUL.md/AGENTS.md (Legacy's identity), monitors all skill behavior, gates future skill installs.

---

### #2 — ClawRouter (Smart LLM Router)

| Field | Value |
| --- | --- |
| **Repo** | Referenced in awesome lists (community builds) |
| **Type** | Plugin / middleware |
| **Install** | Install as OpenClaw plugin or standalone proxy |
| **Effort** | 30 minutes |
| **Impact** | High — **78% cost savings** documented |
| **Status** | 🔜 Deferred — evaluate after foundry (#8) learns our usage patterns |
| **Pillar** | UI (cost optimization) |

**What it does**: Sits between OpenClaw and LLM providers. Analyzes each prompt's complexity in real-time and routes to the cheapest model that can handle it:

- Simple factual question → GPT-4o-mini ($0.15/M tokens)
- Conversation response → Claude Sonnet ($3/M tokens)
- Complex research/reasoning → Claude Opus ($15/M tokens)

Maintains a performance-per-dollar leaderboard. Learns over time which prompt types each model handles best.

**Why we need it**: We hard-code `anthropic/claude-opus-4-6` as Legacy's default in openclaw.json. Burning $15/M tokens on "what time is it?" and "ok thanks". 78% of Legacy's interactions don't need Opus.

**OneMind connection**: Reads from openclaw.json `models` array (already configured: Opus, Sonnet, GPT-4o, GPT-4o-mini). Routes dynamically instead of defaulting to most expensive.

---

### #3 — Agent Protocol Skill (Pub/Sub Agent Chains)

| Field | Value |
| --- | --- |
| **Repo** | lobehub.com/skills/openclaw-skills-agent-protocol |
| **Type** | Skill (ClawHub) |
| **Install** | `openclaw skills install agent-protocol` |
| **Effort** | 15 minutes |
| **Impact** | Medium-high — declarative multi-agent without custom code |
| **Status** | 🔜 Deferred — evaluate after sessions_spawn (#6) is tested |
| **Pillar** | UI (agent coordination) |

**What it does**: Event-bus communication layer for agent-to-agent messaging within OpenClaw's native session system. Define orchestration chains declaratively:

```text
Research → Summary → Notification
```

Legacy publishes a "research" event → subscriber agent picks it up → produces summary → triggers notification agent. Multi-step workflows as event chains.

**Why we need it**: Different from nexus-bridge (which dispatches to Agno via NATS). This works within OpenClaw's native session system using `sessions_send`. Simpler for purely digital agent chains.

**OneMind connection**: Complements nexus-bridge. Use Agent Protocol for digital-only chains, nexus-bridge for cross-runtime (Agno, OM1) dispatch.

---

### #4 — openclaw-mcp-plugin (Universal MCP Client)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/lunarpulse/openclaw-mcp-plugin |
| **Type** | Plugin |
| **Install** | Install as OpenClaw plugin |
| **Effort** | 30 minutes |
| **Impact** | High — unlocks **13,000+ MCP servers** |
| **Status** | 📋 Next sprint — high priority after clawsec |
| **Pillar** | UI (tool ecosystem) |

**What it does**: Enables OpenClaw to connect to **any MCP server** via Streamable HTTP transport. MCP (Model Context Protocol) has 13,000+ community servers — databases, APIs, file systems, cloud services. Install this plugin once, then Legacy can mount any MCP server as a tool source.

**Examples**:
- Mount a Postgres MCP server → Legacy queries any database
- Mount a Figma MCP server → Legacy inspects designs
- Mount a Stripe MCP server → Legacy checks revenue
- Mount a GitHub MCP server → Legacy manages repos

**Why we need it**: Currently Legacy can only use tools registered as OpenClaw skills. This opens the entire MCP ecosystem.

**OneMind connection**: OpenClaw already supports MCP natively (`@modelcontextprotocol/sdk@1.25.3`). This plugin unlocks remote MCP servers in addition to local ones.

---

### #5 — Workspace Backup (Cron → Private GitHub)

| Field | Value |
| --- | --- |
| **Repo** | Pattern from andreesg/openclaw-docker-config |
| **Type** | Cron job (direct integration) |
| **Install** | Add to crontab |
| **Effort** | 15 minutes |
| **Impact** | Critical — disaster recovery |
| **Status** | ✅ Done — git-sync.lob runs every 15min, local git active |
| **Pillar** | UI (infrastructure) |

**What it does**: Cron job running every 15 minutes that auto-pushes `~/.openclaw/workspace/` to a private GitHub repo:

```bash
*/15 * * * * cd ~/.openclaw/workspace && git add -A && git commit -m "auto-backup $(date +\%Y-\%m-\%d_\%H:\%M)" --allow-empty && git push origin main
```

**Why we need it**: SOUL.md, AGENTS.md, USER.md, memory files, skill configs — Legacy's entire identity lives in this folder. If the Droplet dies, these are gone. We have git sync for Codex but NOT for the OpenClaw workspace.

**OneMind connection**: Recovery time goes from "rebuild from scratch" to `git clone` on a new Droplet.

---

### #6 — sessions_send / sessions_spawn (Built-in Multi-Agent)

| Field | Value |
| --- | --- |
| **Repo** | Built into OpenClaw core (no install needed) |
| **Type** | Built-in feature |
| **Install** | Already available — just use it |
| **Effort** | 1 hour to learn/test |
| **Impact** | High — may simplify our digital dispatch |
| **Status** | ✅ Done — built-in feature, available for use |
| **Pillar** | UI (agent coordination) |

**What it does**: OpenClaw's native multi-agent coordination:

- **`sessions_send(sessionId, message)`** — send a message to any other agent session by ID. Cross-context messaging.
- **`sessions_spawn({ tools, instructions })`** — creates a child agent with restricted tool subset and specific instructions. Child runs, returns results, and dies.
- **Hierarchical sessions** — parent/child with scoped permissions.
- **RFC-0001** (Discussion #10036) proposes formal Agent Teams with SendMessage/broadcast/shutdown.

**Why we need it**: We built nexus-bridge for NATS-based dispatch to Agno. For purely digital agent-to-agent communication (no IoT), the built-in system may be simpler and lower latency.

**OneMind connection**: Evaluate whether digital dispatch (Legacy → researcher) should use sessions_spawn (built-in, zero network hop) vs nexus-bridge (NATS, needed for cross-runtime). Likely answer: sessions_spawn for digital, nexus-bridge for physical.

---

## TIER 1 — Same Day (hours)

---

### #7 — supermemory (Auto-Recall Memory)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/supermemoryai/openclaw-supermemory |
| **Type** | Plugin |
| **Install** | Install as OpenClaw plugin |
| **Effort** | 30 minutes |
| **Impact** | High — makes OmniDB pgvector actually useful |
| **Status** | 🔜 Deferred — evaluate vs OmniDB pgvector semantic_search |
| **Pillar** | UI (memory architecture) |

**What it does**: Two functions:

1. **Auto-recall**: Before every LLM turn, queries a memory store for relevant context and injects it into the prompt. If Zeus says "remember that restaurant?" — supermemory finds the conversation and injects it.
2. **Auto-capture**: After every conversation, extracts key facts and stores them automatically.

**Why we need it**: We built OmniDB with pgvector embeddings and `semantic_search()` but nothing automatically queries it before each response. Legacy has long-term memory written but never reads it back without being asked.

**OneMind connection**: Configure to query OmniDB embeddings on every turn. Closes the read loop on our pgvector investment.

---

### #8 — openclaw-foundry (Self-Upgrading Agent)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/lekt9/openclaw-foundry |
| **Type** | Plugin / meta-extension |
| **Install** | Install as OpenClaw plugin |
| **Effort** | 30 minutes |
| **Impact** | High — agent learns from itself |
| **Status** | 📋 Next sprint — install after clawsec gates skill creation |
| **Pillar** | UI (agent evolution) |

**What it does**: "The forge that forges itself." Observes Legacy's workflow patterns over time. When it detects Legacy doing the same multi-step sequence repeatedly (e.g., "fetch weather → check calendar → compose message → send WhatsApp"), it crystallizes that pattern into a new Lobster workflow or skill automatically. Generates the SKILL.md and optionally index.js, then suggests installation.

**Why we need it**: Instead of us manually writing skills and workflows, Legacy watches what it does repeatedly and auto-generates them.

**OneMind connection**: Works with Lobster (#10) — foundry generates `.lob` files. Feeds into the skill ecosystem over time.

---

### #9 — Composio (500+ App Integrations)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/ComposioHQ/secure-openclaw |
| **Type** | Plugin + Docker container |
| **Install** | Docker container + OpenClaw plugin |
| **Effort** | 1 hour |
| **Impact** | High — instant access to 500+ SaaS APIs |
| **Status** | 🔜 Deferred — MCP plugin (#4) covers most use cases without Docker overhead |
| **Pillar** | All pillars |

**What it does**: Pre-built app integrations — Gmail, Slack, Notion, Linear, Jira, Stripe, GitHub, Google Sheets, Airtable, HubSpot, Salesforce, Intercom, Zendesk, and 480+ more. Each integration handles OAuth, token refresh, rate limiting, and error handling. Wrapped in a hardened shell with Docker isolation, persistent memory, and cron scheduling.

**Why we need it**: Instead of writing custom skills for every SaaS tool, Composio provides the connectors. Legacy says "create a Notion page" → Composio handles OAuth and API.

**OneMind connection**: Direct value for GE pillar (Stripe, HubSpot, Salesforce), HP pillar (health apps, calendars), LE pillar (family shared tools).

---

### #10 — Lobster Workflows (3 Core .lob Files)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/openclaw/lobster |
| **Type** | Workflow files (.lob) |
| **Install** | Write files to `~/.openclaw/workspace/lobster/` |
| **Effort** | 2 hours |
| **Impact** | High — automated daily operations |
| **Status** | ✅ Done — 9 workflows deployed (morning-brief, health-check, git-sync, habit-checker, tech-digest, second-brain, earnings-tracker, project-state, README) |
| **Pillar** | UI (automation) |

**What it does**: Three core workflows:

**morning-brief.lob** — Runs at 7:00 AM via cron:
```text
Calendar API → Todoist tasks → Weather API → News headlines →
Farm sensors (NATS) → Office status → Compose summary → Send WhatsApp
```

**git-sync.lob** — Runs every 15 minutes:
```text
git pull Codex → detect changes → summarize diffs → archive to OmniDB →
optionally alert Zeus of important changes
```

**health-check.lob** — Runs every 5 minutes:
```text
Ping NATS → Ping OmniDB → Ping Agno → Check disk/memory →
Compose status → Alert if anything down
```

**Why we need it**: Lobster is already enabled (`"lobster": {"enabled": true}`) in openclaw.json but we have zero `.lob` files. Workflows are the backbone of a human enhancement framework — things that happen predictably without being asked.

---

### #11 — openclaw-ansible (Hardened Deployment)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/openclaw/openclaw-ansible |
| **Type** | Ansible playbook (infrastructure-as-code) |
| **Install** | Clone + customize + run |
| **Effort** | 2 hours |
| **Impact** | Critical — security hardening |
| **Status** | 📋 Next sprint — run after clawsec for full hardening pass |
| **Pillar** | UI (infrastructure) |

**What it does**: Official Ansible playbook for automated, hardened OpenClaw installation:

1. Installs OpenClaw runtime
2. Configures **Tailscale VPN** (encrypted mesh networking)
3. Sets up **UFW firewall** (blocks all ports except needed)
4. Installs **Fail2ban** (auto-bans IPs after failed login attempts)
5. Configures Docker isolation for containers
6. Sets up systemd service for OpenClaw persistence
7. Configures log rotation
8. Handles SSL/TLS certificates

**Why we need it**: We have a manual `setup-droplet.sh`. This replaces it with repeatable, auditable infrastructure-as-code with security hardening we're currently missing.

**OneMind connection**: Run once to harden existing Droplet. Re-run to provision new locations (office, farm, van nodes).

---

### #12 — openclaw-hardened-ansible (Tier 3+ Security)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/Next-Kick/openclaw-hardened-ansible |
| **Type** | Ansible playbook (advanced security) |
| **Install** | Clone + customize + run |
| **Effort** | 1 hour (evaluate), 4 hours (implement) |
| **Impact** | High — defense in depth |
| **Status** | 🔜 Deferred — evaluate after #11 ansible baseline is running |
| **Pillar** | UI (security) |

**What it does**: Three upgrades beyond the official ansible:

1. **Rootless Podman** instead of Docker — no root daemon, each container runs as unprivileged user.
2. **Squid proxy for egress filtering** — outbound traffic goes through a whitelist proxy. Malicious skill tries to exfiltrate data? Squid blocks it.
3. **LiteLLM credential brokering** — OpenClaw never sees real API keys. LiteLLM holds the keys and proxies LLM requests. OpenClaw compromised? Attacker can't steal API keys.

**Why we need it**: Defense in depth. clawsec (#1) monitors behavior. This hardens the infrastructure. Together they provide layered security.

---

### #13 — awesome-openclaw (195K★ Definitive Collection)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/rohitg00/awesome-openclaw |
| **Type** | Git subtree (reference documentation) |
| **Install** | `git subtree add --prefix=openclaw/docs/ecosystem/awesome-openclaw` |
| **Effort** | 15 minutes |
| **Impact** | Medium — ecosystem discovery layer |
| **Status** | ✅ Done — subtree added |
| **Pillar** | All |

**What it does**: THE definitive curated OpenClaw collection — 80+ ecosystem projects organized by category, hosting guides for every cloud provider (DO, Hetzner, Fly.io, Railway, Coolify), cost comparison calculators, security hardening checklist, deployment tutorials. Powers openclawsearch.com.

**Why we need it**: Offline access to the entire ecosystem discovery layer. When we need to find a skill, tool, or pattern — search here first.

---

### #14 — awesome-openclaw-skills (2,868 Curated Skills)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/VoltAgent/awesome-openclaw-skills |
| **Type** | Git subtree (skill catalog) |
| **Install** | `git subtree add --prefix=openclaw/docs/ecosystem/awesome-openclaw-skills` |
| **Effort** | 15 minutes |
| **Impact** | High — curated skill shopping list with security ratings |
| **Status** | ✅ Done — subtree added |
| **Pillar** | All |

**What it does**: 2,868 skills curated from 5,705+ on ClawHub, organized by category with **security warnings** on each. Categories include: which skills request exec permissions, which make outbound network calls, which access the filesystem.

**Why we need it**: Instead of searching ClawHub blindly, browse a curated catalog with risk ratings. Our skill shopping list.

---

### #15 — Starter Kit Patterns (3 Repos)

| Field | Value |
| --- | --- |
| **Repos** | jeffweisbein, zach-highley, MaximusCarapax starter kits |
| **Type** | Reference (study and steal patterns) |
| **Install** | Read and extract patterns |
| **Effort** | 1 hour |
| **Impact** | Medium — battle-tested patterns |
| **Status** | 📋 Next sprint — extract spend caps + failover patterns |
| **Pillar** | UI (reliability, security) |

**What to extract**:

| Starter Kit | Key Pattern to Steal |
| --- | --- |
| **jeffweisbein** | SOUL.md personality design, daily API spend caps ($50/day), tool usage limits |
| **zach-highley** | Self-healing via launchd/systemd (auto-restart on crash), multi-model failover (Claude → GPT-4o → Gemini), **355-item security audit checklist** |
| **MaximusCarapax** | Aider coding integration, RAG with free Gemini embeddings, Gmail skill patterns |

**OneMind connection**: Apply spend caps to openclaw.json. Implement multi-model failover in agent defaults. Run the 355-item security checklist against our deployment.

---

### #16 — elevenlabs-voices Skill (18 Personas, 32 Languages)

| Field | Value |
| --- | --- |
| **Repo** | ClawHub skill: elevenlabs-voices |
| **Type** | Skill |
| **Install** | `openclaw skills install elevenlabs-voices` |
| **Effort** | 15 minutes |
| **Impact** | Medium — multi-persona, multilingual voice |
| **Status** | ⏳ Partial — ElevenLabs API key + voice ID configured, skill not yet installed |
| **Pillar** | HP (personal interaction), LE (family) |

**What it does**: Production-ready ElevenLabs skill with 18 curated voice personas (not just one voice — different characters for different contexts), 32 language support, streaming audio, SFX sound effects generation, voice design tool (create custom voices from text description), pronunciation dictionary.

**Why we need it**: Legacy could use different voices for: formal business calls, casual family chat, bedtime stories for kids. Multilingual for family members who speak other languages.

---

## TIER 2 — Docker Compose Additions (days)

---

### #17 — n8n (Deterministic Workflow Engine)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/caprihan/openclaw-n8n-stack (reference) |
| **Type** | Docker container |
| **Port** | 5678 |
| **Effort** | 1 hour |
| **Impact** | High — deterministic execution layer |
| **Status** | 🔜 Deferred — already have n8n in Codex infra stack |
| **Pillar** | All pillars |

**What it does**: 400+ built-in integrations (Slack, Gmail, Stripe, databases, APIs). Visual workflow builder at :5678. OpenClaw = reasoning layer, n8n = execution layer. Connected via webhooks:

```text
User asks Legacy something → Legacy reasons about it →
Legacy fires webhook to n8n → n8n workflow runs →
queries DB, formats email, hits APIs → result posts back to OpenClaw
```

Handles retries, error branching, scheduling, data transformation.

**Docker Compose snippet**:
```yaml
n8n:
  image: n8nio/n8n:latest
  container_name: onemind-n8n
  restart: unless-stopped
  profiles: ["full"]
  ports:
    - "5678:5678"
  volumes:
    - n8n_data:/home/node/.n8n
  environment:
    N8N_HOST: 0.0.0.0
    N8N_ENCRYPTION_KEY: ${N8N_ENCRYPTION_KEY}
    WEBHOOK_URL: ${N8N_WEBHOOK_URL:-http://host.docker.internal:5678}
  networks:
    - onemind-network
  extra_hosts:
    - "host.docker.internal:host-gateway"
```

---

### #18 — LiveKit (Real-Time AV / WebRTC)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/livekit/livekit |
| **Type** | Docker container |
| **Ports** | 7880 (HTTP), 7881 (WebRTC TCP), 7882 (WebRTC UDP) |
| **Effort** | 1 hour |
| **Impact** | High — enables smart glasses, voice avatar, video |
| **Status** | 🔜 Deferred — needs VisionClaw (#23) or avatar (#32) use case first |
| **Pillar** | UI (infrastructure), HP (voice interaction) |

**What it does**: Open-source WebRTC server for real-time audio, video, and data streaming. Three use cases:

1. **Smart glasses** (VisionClaw #23) — WebRTC POV streaming from Meta Ray-Bans at 2.5 Mbps
2. **Voice avatar** (#32) — lip-synced visual Legacy with Deepgram + ElevenLabs + LemonSlice
3. **Multi-party rooms** — Zeus, Legacy, and agents in a voice room

**Docker Compose snippet**:
```yaml
livekit:
  image: livekit/livekit-server:latest
  container_name: onemind-livekit
  restart: unless-stopped
  profiles: ["full"]
  ports:
    - "7880:7880"
    - "7881:7881"
    - "7882:7882/udp"
  environment:
    LIVEKIT_KEYS: "${LIVEKIT_API_KEY}:${LIVEKIT_API_SECRET}"
  networks:
    - onemind-network
```

---

### #19 — openclaw-studio (Admin Dashboard)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/grp06/openclaw-studio |
| **Type** | Web application (container or local) |
| **Effort** | 2 hours |
| **Impact** | Medium — admin visibility |
| **Status** | 🔧 Installing — npx openclaw-studio on Mac, connecting to VPS gateway |
| **Pillar** | UI (observability) |

**What it does**: Web-based admin dashboard:

- Visual agent creation/configuration
- Cron job manager (create/edit/delete via UI)
- Exec approval queue (approve/deny shell commands)
- **Tool call visualization** (see every tool call in real-time with input/output)
- **Thinking traces** (see LLM chain-of-thought for debugging)

**Why we need it**: We have no visibility into what Legacy is doing right now. Studio is for Zeus-as-admin, not Zeus-as-user.

---

### #20 — Mem0 (Persistent Memory Layer)

| Field | Value |
| --- | --- |
| **Repo** | https://docs.mem0.ai/integrations |
| **Type** | Library / service (evaluate) |
| **Effort** | 2 hours to evaluate |
| **Impact** | Medium — purpose-built memory vs our general-purpose approach |
| **Status** | 🔜 Deferred — evaluate after supermemory (#7) vs OmniDB comparison |
| **Pillar** | UI (memory architecture) |

**What it does**: Persistent memory layer with first-class OpenClaw support. Automatic deduplication, conflict resolution, and decay (old irrelevant memories fade). Stores: user preferences, conversation history, entity extraction, relationship mapping.

**Why we need it**: Different from our OmniDB approach — Mem0 is purpose-built for agent memory while OmniDB is general-purpose (time-series + vectors).

**Evaluation question**: Does Mem0 on top of OmniDB give better recall than raw pgvector queries? Or is supermemory (#7) sufficient?

---

### #21 — acpx (Universal Agent Bridge)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/openclaw/acpx |
| **Type** | CLI tool (Agent Client Protocol) |
| **Install** | npm install -g acpx |
| **Effort** | 2 hours |
| **Impact** | High — Legacy orchestrates other AI assistants |
| **Status** | 📋 Next sprint — alpha status, Mac only, install when stable |
| **Pillar** | UI (agent orchestration) |

**What it does**: Headless CLI client for Agent Client Protocol (ACP). Bridges external AI agents into OpenClaw sessions. Supported: **Claude Code, GitHub Copilot, Gemini, OpenCode, Pi**, and 14+ others.

**Use cases**:
- Legacy spawns Claude Code via acpx → Claude writes code in sandbox → result flows back → Legacy reviews and commits
- Legacy spawns Gemini session for multimodal analysis (images, video)
- OpenClaw becomes a **universal AI orchestrator**

**OneMind connection**: Issue #28511 proposes expanding acpx to support any ACP agent from a registry of 19+ agents.

---

### #22 — openclaw-mcp (Claude.ai → Legacy Bridge)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/freema/openclaw-mcp |
| **Type** | MCP server |
| **Install** | Deploy as MCP server, configure OAuth |
| **Effort** | 2 hours |
| **Impact** | Medium — control Legacy from Claude Desktop |
| **Status** | 🔧 Installing — Docker on VPS, enabling chatCompletions |
| **Pillar** | UI (access layer) |

**What it does**: MCP server bridging Claude.ai to self-hosted OpenClaw with OAuth 2.1 authentication. "AI assistant orchestrating another AI assistant."

**Use case**: Chatting with Claude on claude.ai → "tell Legacy to check the morning briefing" → Claude calls MCP server → authenticates with OpenClaw gateway → sends message. Two-way: Legacy can also call Claude via MCP.

---

## TIER 3 — Hardware & Subtrees (weeks)

---

### #23 — VisionClaw (Meta Ray-Ban Smart Glasses)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/sseanliu/VisionClaw |
| **Type** | Git subtree + iOS/Android apps |
| **Install** | `git subtree add --prefix=openclaw/hardware/wearables/smart-glasses/visionclaw` |
| **Effort** | 1 week |
| **Impact** | Transformational — Legacy gets eyes |
| **Status** | 🔜 Deferred — requires LiveKit (#18) + Meta Ray-Bans hardware |
| **Pillar** | UI (embodiment), HP (ambient awareness) |
| **Requires** | LiveKit (#18) |

**What it does**: Real-time AI assistant for Meta Ray-Ban smart glasses. Glasses stream POV video via WebRTC at 2.5 Mbps → VisionClaw server runs Gemini Live for real-time visual understanding → OpenClaw executes actions with 56+ tools. iOS (Swift) and Android (Kotlin) apps.

**Example**: Zeus looks at a menu in Italian through glasses → VisionClaw → Gemini sees the text → Legacy: "The daily special is osso buco for €18. It's braised veal shank — you liked the one at that place in Milan."

**OneMind connection**: We already have META_GLASSES_INTEGRATION.md in the Codex planning this. VisionClaw IS the implementation.

---

### #24 — WireClaw ($5 ESP32 NATS Agents)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/M64GitHub/WireClaw |
| **Type** | Git subtree + hardware |
| **Install** | Web flash at wireclaw.io, no toolchain needed |
| **Effort** | 3 days |
| **Impact** | High — spatial awareness in every room |
| **Status** | 🔜 Deferred — hardware purchase needed, evaluate after HA integration |
| **Pillar** | HP (environment), LE (home management) |

**What it does**: AI agent on a $5 ESP32 with 20 built-in tools, Telegram bot, and **native NATS** pub/sub. Rule engine with multi-step chains (if motion AND time > 10pm → lights on AND alert). Persistent local automation — works even if internet is down. Cross-device messaging via NATS at **30ms latency**.

**Deployment**: One per room. Each publishes to `wireclaw.{device_id}.{event}` NATS subjects. Legacy discovers them automatically.

---

### #25 — HeyClawy (ESP32 Voice Companion)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/omeriko9/HeyClawy |
| **Type** | Hardware (ESP32 + speaker) |
| **Effort** | 2 days |
| **Impact** | Medium — voice in every room |
| **Status** | 🔜 Deferred — hardware purchase needed |
| **Pillar** | HP (hands-free interaction) |

**What it does**: ESP32 + speaker + optional display. Wake word via WakeNet ("Hey Legacy"). Captures voice → sends to OpenClaw → response plays on speaker. Local web settings page for WiFi config. Under $10 per unit.

**Use case**: "Hey Legacy, what's on my calendar?" from kitchen, bathroom, office — no phone needed.

---

### #26 — MimiClaw (ESP32-S3 Pure C Agent)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/memovai/mimiclaw |
| **Type** | Hardware firmware |
| **Effort** | 1 day to evaluate |
| **Impact** | Medium — ultra-cheap local reasoning |
| **Status** | 🔜 Deferred — research only, needs ESP32-S3 hardware |
| **Pillar** | UI (edge intelligence) |

**What it does**: OpenClaw-like agent on ESP32-S3 in **pure C**. No Node.js, no Linux, no OS — bare metal firmware. ReAct agent loop (Reason → Act → Observe → Repeat). Persistent memory via SOUL.md/USER.md/MEMORY.md on flash. Telegram integration. Runs on 512KB RAM.

**Why it matters**: Sensors that reason locally about their data before sending to NATS. "Temperature is 95°F — that's unusual for March" happens on the $5 chip.

---

### #27 — OpenClawHomeAssistant (HA Add-on)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/techartdev/OpenClawHomeAssistant |
| **Type** | Home Assistant add-on |
| **Effort** | 2 days |
| **Impact** | High — OpenClaw inside HA |
| **Status** | 🔜 Deferred — evaluate vs REST skill + MCP |
| **Pillar** | HP (environment), LE (home) |

**What it does**: Runs OpenClaw **inside** Home Assistant OS. AI Gateway, Web Terminal, Assist Pipeline integration (OpenClaw becomes conversation agent for "Hey Home"), browser automation via Playwright, persistent storage.

**Different from REST skill**: The REST skill calls HA from outside. This add-on runs OpenClaw inside HA with direct access to all entities, automations, scenes — no network hops.

---

### #28 — openclaw_conversation (HA Custom Component)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/Djelibeybi/openclaw_conversation |
| **Type** | HACS custom component |
| **Effort** | 1 day |
| **Impact** | Medium — lightweight HA integration |
| **Status** | 🔜 Deferred — evaluate alongside #27 |
| **Pillar** | HP (environment) |

**What it does**: Lighter than the full add-on (#27). Installs via HACS. Uses OpenClaw as AI backend for HA's conversation and task platforms. Uses `ha-mcp` for device control via MCP. Can be set as Assist pipeline conversation agent.

**Complements #27**: Add-on runs OpenClaw inside HA; this component connects an external OpenClaw to HA. Different architectures for different setups.

---

### #29 — HomeClaw (Apple HomeKit via MCP)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/omarshahine/HomeClaw |
| **Type** | Mac Catalyst app + MCP |
| **Effort** | 2 days |
| **Impact** | Medium — native HomeKit control |
| **Status** | 🔜 Deferred — Mac Catalyst app, evaluate when HA integration mature |
| **Pillar** | HP (environment) |

**What it does**: Direct Apple HomeKit control — lights, locks, thermostats, scenes, rooms. Mac Catalyst app with native HomeKit framework. CLI tool. TestFlight beta for iOS. Plugin registers HomeKit as OpenClaw skill.

**Use case**: HomeKit devices controlled natively without HA middleman. Use alongside HA — HomeClaw for Apple ecosystem, HA for Zigbee/Z-Wave/MQTT.

---

### #30 — Jupiter Voice (Fully Local STT/TTS)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/SudarshanSuryaprakash/voice-assistant |
| **Type** | Local application (Apple Silicon optimized) |
| **Effort** | 1 day |
| **Impact** | Medium — offline voice, zero latency |
| **Status** | 🔜 Deferred — ElevenLabs working, evaluate as fallback later |
| **Pillar** | UI (resilience) |

**What it does**: Three components:

1. **MLX-accelerated Whisper** — 10x faster STT on Apple Silicon
2. **OpenWakeWord** — custom wake phrase ("Hey Legacy")
3. **Kokoro ONNX TTS** — local text-to-speech, no API calls

**Why we need it**: ElevenLabs goes down = Legacy goes mute. Internet drops = no voice. Jupiter Voice runs 100% local. Configure as failover: ElevenLabs (primary) → Jupiter Voice (fallback).

---

### #31 — openclaw-voice (Browser-Based Voice Chat)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/Purple-Horizons/openclaw-voice |
| **Type** | Web application (self-hosted) |
| **Effort** | 1 day |
| **Impact** | Medium — zero-friction voice access |
| **Status** | 🔜 Deferred — browser voice not priority, Discord voice working |
| **Pillar** | HP (accessibility) |

**What it does**: Browser-based voice chat. Whisper STT runs locally via faster-whisper (audio stays on network). ElevenLabs or Chatterbox for TTS. No app install — open a browser tab on any device and talk to Legacy.

**Use case**: Guests, family members, any device with a browser. No app install, no account, just talk.

---

### #32 — openclaw-voice-avatar (LiveKit Lip-Synced Face)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/openserv-labs/openclaw-voice-avatar |
| **Type** | Application (requires LiveKit #18) |
| **Effort** | 3 days |
| **Impact** | High — Legacy with a visual face |
| **Status** | 🔜 Deferred — requires LiveKit (#18) |
| **Pillar** | HP (personal connection), LE (family) |
| **Requires** | LiveKit (#18) |

**What it does**: Pipeline: Deepgram transcription (Zeus speaks) → OpenClaw Gateway (Legacy thinks) → ElevenLabs TTS (Legacy speaks) → LemonSlice avatar generation (Legacy's face moves in sync with speech). ~1-2 second latency.

**Result**: Legacy has a visual face that talks with lip-sync. Display on tablet, smart display, or Flutter app. Makes Legacy feel like a person, not a text box.

---

## TIER 4 — Strategic / GE Pillar (months)

---

### #33 — Virtual Protocol ACP (Agent Commerce)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/Virtual-Protocol/openclaw-acp |
| **Type** | Plugin + blockchain integration |
| **Effort** | 1 week to evaluate |
| **Impact** | GE pillar — agent-generated revenue |
| **Status** | 🔜 Deferred — blockchain commerce, evaluate when GE pillar active |
| **Pillar** | GE (Generational Empire) |

**What it does**: Agent-to-agent commerce on Base blockchain (Coinbase L2). Auto-provisioned crypto wallets. Browse marketplace, create bounties, list agent services, tokenize capabilities as on-chain assets.

**GE application**: Legacy offers research reports, content creation, data analysis as paid services.

---

### #34 — Network-AI (12-Framework Orchestrator)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/jovanSAPFIONEER/Network-AI |
| **Type** | Framework bridge |
| **Effort** | 1 week |
| **Impact** | High — best tool for each job |
| **Status** | 🔜 Deferred — heavy framework, evaluate after core stable |
| **Pillar** | UI (agent orchestration) |

**What it does**: Single interface to 12 AI frameworks: LangChain, AutoGen, CrewAI, OpenAI Assistants, LlamaIndex, Semantic Kernel, Haystack, DSPy, Agno, MCP, OpenClaw, custom adapters. Shared blackboard, AES-256 encryption, hallucination detection.

**Use case**: Research → LangChain RAG. Structured extraction → DSPy. Code → OpenAI Assistants. Not locked to one framework.

---

### #35 — Langclaw (LangGraph Workflows)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/tisu19021997/langclaw |
| **Type** | Python framework |
| **Effort** | 3 days |
| **Impact** | Medium — more powerful than Lobster for complex workflows |
| **Status** | 🔜 Deferred — Python dependency, Lobster covers current needs |
| **Pillar** | UI (automation) |

**What it does**: Compiles to real LangGraph `CompiledStateGraph` — graph-based state machines with conditional branching, parallel execution, human-in-the-loop checkpoints. RBAC, scheduled tasks, persistent memory, subagent delegation. Multi-channel from day one.

**Trade-off**: Python (vs Lobster's JS-native). Heavier but more powerful for complex branching workflows.

---

### #36 — ExoClaw (WASM-Sandboxed Agent Runtime)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/jbold/exoclaw |
| **Type** | Rust runtime |
| **Effort** | 1 week to evaluate |
| **Impact** | Medium — security for untrusted skills |
| **Status** | 🔜 Deferred — research, clawsec (#1) handles trust gating first |
| **Pillar** | UI (security) |

**What it does**: WASM-sandboxed runtime in Rust. Every agent runs in a WASM sandbox — physically cannot access filesystem, network, or memory outside sandbox. Uses NATS as core message bus (`exoclaw.channel.account.peer`). Graceful fallback if NATS unavailable.

**Use case**: Run untrusted community skills in ExoClaw instead of giving them full Node.js access by default.

---

### #37 — BankrBot Skills (Crypto/DeFi/Polymarket)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/BankrBot/openclaw-skills |
| **Type** | Skill collection |
| **Effort** | 2 days |
| **Impact** | GE pillar — financial automation |
| **Status** | 🔜 Deferred — needs L1 HITL approval flows before any financial skill |
| **Pillar** | GE (Generational Empire) |

**What it does**: Polymarket trading (prediction market bets), crypto trading (buy/sell on exchanges), DeFi operations (yield farming, liquidity provision, swaps). Risk parameters, position sizing, stop-loss logic.

**Critical**: Needs L1 (HITL) approval flows for any transaction. Never give a financial skill L3 autonomy.

---

### #38 — ClawWork Benchmark (Economic Capability Testing)

| Field | Value |
| --- | --- |
| **Repo** | https://github.com/HKUDS/ClawWork |
| **Type** | Benchmark framework |
| **Effort** | 3 days |
| **Impact** | GE pillar — identify monetizable capabilities |
| **Status** | 🔧 Installing — conda env on Mac, React dashboard |
| **Pillar** | GE (Generational Empire) |

**What it does**: Economic survival benchmark: agents start with $10, must complete 220 professional tasks across 44 sectors (writing, coding, design, research, marketing, legal, finance). Top models hit **$1,500+/hour equivalent**.

**Use case**: Run against Legacy + Agno workers to measure economic capability. Identifies which task types are monetizable vs need improvement.

---

### #39 — CrewClaw (Visual Multi-Agent Dashboard)

| Field | Value |
| --- | --- |
| **Repo** | https://www.crewclaw.com/ |
| **Type** | Web platform |
| **Effort** | 2 days to evaluate |
| **Impact** | Medium — visual multi-agent management |
| **Status** | 🔜 Deferred — evaluate after openclaw-studio (#19) is running |
| **Pillar** | UI (observability) |

**What it does**: Visual dashboard for managing 3+ agent teams. Drag-and-drop workflow builder. Real-time monitoring: agent status, current tasks, cost tracking, performance metrics. Alert when agents fail.

**vs openclaw-studio (#19)**: Studio shows individual agent internals. CrewClaw shows multi-agent fleet overview.

---

### #40 — OM1 Sub-Repos (5 Robotics Repos)

| Field | Value |
| --- | --- |
| **Repo** | OpenMind GitHub org (5 repos) |
| **Type** | Research / future bridge |
| **Effort** | 1 week |
| **Impact** | Strategic — physical embodiment |
| **Status** | 🔜 Deferred — research phase, no hardware yet |
| **Pillar** | UI (robotics), LE (physical assistance) |

**Five repos**:

1. **OM1-video-processor** — Docker-based real-time video streaming with face recognition and audio for robots.
2. **OM1-ros2-sdk** — ROS 2 SLAM package for Unitree Go2 quadruped. Navigation, mapping, localization.
3. **om-yolo-ws** — WebSocket-based remote YOLO object detection. Send camera frame → get bounding boxes.
4. **humanoid_teleoperate** — XR device teleoperation for Unitree humanoid robots.
5. **FABRIC** — Decentralized coordination layer for robot identity, trust, and data sharing.

**OneMind connection**: When we build the OpenClaw ↔ OM1 bridge (our unique competitive moat), these repos become directly relevant.

---

### #41 — Clawdstrike (Fleet Security Monitoring)

| Field | Value |
| --- | --- |
| **Repo** | Referenced in awesome lists |
| **Type** | Security tool |
| **Effort** | 2 days to evaluate |
| **Impact** | Medium (future — when we have multiple agents) |
| **Status** | 🔜 Deferred — need multi-agent fleet first |
| **Pillar** | UI (security at scale) |

**What it does**: Runtime security and swarm detection for agent fleets. Monitors behavior patterns across multiple OpenClaw instances. Detects: compromised agents, unusual tool call patterns, coordination attacks, resource abuse, data exfiltration.

**When we need it**: Not now (one Legacy instance). When we scale to community agent + Agno workers + IoT agents.

---

### #42 — Community Platforms (Ecosystem Participation)

| Field | Value |
| --- | --- |
| **Type** | Ongoing participation |
| **Effort** | Continuous |
| **Impact** | Strategic — knowledge, recruiting, positioning |
| **Status** | 🔜 Deferred — join Discord + Reddit when time allows |

**Platforms to join and participate**:

| Platform | Size | URL | Purpose |
| --- | --- | --- | --- |
| Discord | 115,000 | https://discord.gg/clawd | Real-time support, announcements, core dev access |
| Reddit | Active | r/openclaw | Guides, showcases, project sharing |
| X Community | 24,400 | twitter communities | Builder community, configs, snippets |
| GitHub Discussions | Thousands | github.com/openclaw | Official Q&A, RFCs, feature requests |
| Clawstin | In-person | https://clawstin.com/ | Austin meetup |
| Silicon Valley Meetup | In-person | meetup.com | West coast networking |
| AI Tinkerers SF | 40 slots | sf.aitinkerers.org | Curated build sessions |
| Luma Calendar | Global | https://luma.com/claw | Worldwide events |
| Newsletter | Daily | buttondown.com/openclaw-newsletter | Ecosystem updates |

**Value**: Early knowledge of breaking changes, direct access to skill authors, recruitment pipeline for GE pillar, positioning OneMind OS as a reference implementation.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/platform-catalog (Platform Catalog)/inventory (Platform Inventory)/inventory|inventory]]
