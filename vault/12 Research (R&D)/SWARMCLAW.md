# SwarmClaw — Architecture & Operations Reference

> SwarmClaw v0.7.7 — Multi-provider AI orchestration cockpit for OneMind OS
> Deployed at: `https://onemind-main.tail717ba0.ts.net:8443` (Tailscale) / `http://104.248.4.212:3456` (direct)
> Last updated: March 2026

---

## What SwarmClaw Is

SwarmClaw is the **meta layer / cockpit** for OneMind OS. It is NOT an OpenClaw replacement — it's a multi-provider orchestration dashboard that can connect to 15+ LLM providers simultaneously, manage multiple OpenClaw gateways, run LangGraph workflows, and handle background task scheduling.

```
SwarmClaw (Meta Layer — cockpit)
    ↕ WebSocket RPC (Ed25519)
OpenClaw Gateway (Engine Layer — Legacy + 8 domain agents)
    ↕ NATS JetStream
Agno AgentOS / OM1 / LeRobot
```

**Key distinction**: OpenClaw agents (Legacy, Spartan, Oracle, etc.) get SwarmClaw's *chat passthrough* — they bypass the LangGraph pipeline. Direct API agents (Claude, GPT, Gemini) get the full SwarmClaw feature stack (tools, MCP, memory, plugins, RAG).

---

## Deployment (VPS — `root@104.248.4.212`)

| Item | Value |
|------|-------|
| Port | `3456` |
| Service | `swarmclaw.service` (systemd, user `onemind`) |
| Project path | `/home/onemind/onemind-os/command-center/swarmclaw/` |
| Data path | `/home/onemind/onemind-os/command-center/swarmclaw/data/` |
| Backups | `/home/onemind/backups/swarmclaw/` (every 6h via cron) |
| Tailscale HTTPS | `https://onemind-main.tail717ba0.ts.net:8443` |
| Auth key | `49359e4f6cd67002f2e2b9b2491a6dbb` |
| CREDENTIAL_SECRET | `0642d37ecac9b93a86767760f609fb78cb5f24af943d787bcfc09669c6dce606` |
| Next.js build | Standalone (`npm run build` → `.next/standalone/`) |

### Service Management

```bash
ssh root@104.248.4.212 "systemctl status swarmclaw"
ssh root@104.248.4.212 "systemctl restart swarmclaw"
ssh root@104.248.4.212 "journalctl -u swarmclaw -f"
```

### Post-Build Steps (REQUIRED after every `npm run build`)

```bash
# 1. Copy env
cp .env.local .next/standalone/.env.local
chown onemind:onemind .next/standalone/.env.local

# 2. Static files symlink (or Next.js serves 404s)
ln -sfn ../../static .next/standalone/.next/static
chown -h onemind:onemind .next/standalone/.next/static

# 3. Restart
systemctl restart swarmclaw
```

---

## Key Bugs Fixed (Do Not Revert)

| Bug | Fix | File |
|-----|-----|------|
| `agentId: 'main'` hardcoded | Use `normalizeOpenClawAgentId(session._openclawAgentName)` | `src/lib/server/chat-execution.ts` |
| OpenClaw hitting LangGraph loop | Add `'openclaw'` to `NON_LANGGRAPH_PROVIDER_IDS` | `src/lib/provider-sets.ts` |
| `.env.local` ACCESS_KEY mismatch | Copy correct `.env.local` to standalone dir post-build | Post-build script |
| Credential API field `key` vs `apiKey` | Use `apiKey` field | Credentials setup |
| Static files 404 | Symlink `.next/static` → `.next/standalone/.next/static` | Post-build script |

**Critical line** — `src/lib/provider-sets.ts`:
```typescript
export const NON_LANGGRAPH_PROVIDER_IDS = new Set(['claude-cli', 'codex-cli', 'opencode-cli', 'openclaw'])
```
OpenClaw is in bypass list — this is intentional. OpenClaw agents have their own tools/memory natively.

---

## All Supported Providers

### Direct API (full SwarmClaw feature stack: tools, MCP, memory, plugins)

| Provider | Type | Notable Models |
|----------|------|----------------|
| OpenAI | Direct API | gpt-4o, gpt-4o-mini, gpt-4.1, o3, o3-mini, o4-mini |
| Anthropic | Direct API | claude-sonnet-4-6, claude-opus-4-6, claude-haiku-4-5 |
| Google Gemini | OpenAI-compat | gemini-2.5-pro, gemini-2.5-flash, gemini-2.5-flash-lite |
| DeepSeek | OpenAI-compat | deepseek-chat, deepseek-reasoner |
| Groq | OpenAI-compat | llama-3.3-70b, deepseek-r1-distill, qwen-qwq-32b |
| Together AI | OpenAI-compat | Llama-4-Maverick, DeepSeek-R1, Qwen2.5-72B |
| Mistral | OpenAI-compat | mistral-large, mistral-small, magistral-medium |
| xAI (Grok) | OpenAI-compat | grok-3, grok-3-fast, grok-3-mini |
| Fireworks AI | OpenAI-compat | deepseek-r1, llama, qwen3-235b |
| Ollama | Local HTTP | 50+ models (qwen3.5, llama3.3, etc.) |
| Custom | Any OpenAI-compat | Any model via custom baseUrl + API key |

**Missing natively** (workarounds):
- Azure OpenAI → use Custom provider with Azure endpoint
- AWS Bedrock → use Agno (Python runtime) via NATS
- Alibaba Qwen direct → use Ollama (local) or Fireworks/Together (cloud)

### Bypass (passthrough — no LangGraph/tools/memory from SwarmClaw side)

| Provider | Type | Notes |
|----------|------|-------|
| OpenClaw | WebSocket + Ed25519 | Your 9 OpenClaw agents; tools/memory handled by gateway |
| Claude Code CLI | Local CLI | Requires Claude auth on VPS |
| OpenAI Codex CLI | Local CLI | Requires OpenAI CLI access |
| OpenCode CLI | Local CLI | Multi-provider CLI |

### AI Frameworks Used Internally

| Framework | Purpose |
|-----------|---------|
| LangChain | Universal LLM abstraction for all 15 providers |
| LangGraph | Multi-agent workflows, checkpointed orchestration |

No AutoGen, CrewAI, Semantic Kernel, DSPy, LlamaIndex — just LangChain + LangGraph.

---

## Multi-Gateway Support

SwarmClaw fully supports **multiple OpenClaw gateways simultaneously**.

- Each gateway is a **Gateway Profile** (stored in SwarmClaw's SQLite)
- Profiles have: endpoint URL, wsUrl, credentialId, tags, use-case label, default model
- Agents routed to gateways via: explicit `gatewayProfileId` → tag matching → use-case matching → default
- Maintains **multiple WebSocket connections** at once (keyed by `profile:ID` or `agent:ID`)
- **Default fallback**: `ws://127.0.0.1:18789` (localhost — correct for VPS deployment)

**Gateway URL normalization** (automatic):
```
localhost:18789                      → ws://localhost:18789
https://onemind-main.tail717ba0.ts.net → wss://onemind-main.tail717ba0.ts.net
http://192.168.1.100:18789           → ws://192.168.1.100:18789
```

Use `ws://127.0.0.1:18789` for the main VPS gateway (same machine, no Tailscale overhead).

---

## OpenClaw Config Manageable from SwarmClaw

SwarmClaw has 20 dedicated `/api/openclaw/*` routes:

| Feature | Route | Notes |
|---------|-------|-------|
| Agent personality files | `/api/openclaw/agent-files` | SOUL.md, IDENTITY.md, USER.md, TOOLS.md, HEARTBEAT.md, AGENTS.md, MEMORY.md |
| Cron jobs | `/api/openclaw/cron` | List, create, update, delete |
| Execution approvals | `/api/openclaw/approvals` | Approve/deny pending tool executions |
| Gateway reload mode | `/api/openclaw/gateway` | hot / hybrid / full |
| Skills | `/api/openclaw/skills` | List, enable/disable, install, remove |
| Sandbox env allowlist | `/api/openclaw/sandbox-env` | Docker sandbox env vars |
| Agent permissions | `/api/openclaw/permissions` | TOOL / ORCHESTRATOR / OBSERVER roles |
| Config health check | `/api/openclaw/doctor` | Guided issue detection + repair |
| Deploy | `/api/openclaw/deploy` | Local, VPS via SSH, or hosted (Render/Fly/Railway) |
| Agent sync | `/api/openclaw/sync` | Bidirectional SwarmClaw ↔ gateway agent sync |
| Chat history | `/api/openclaw/history` | Fetch, paginate, clear |
| Network discovery | `/api/openclaw/discover` | Auto-discover gateways via mDNS/UDP |

**NOT manageable from SwarmClaw** (gateway-side only):
- Channel config (Discord, iMessage, Telegram) — set in `openclaw.json`
- Hooks — set in `openclaw.json`
- Raw workspace JS/TS files — only markdown agent files

---

## What SwarmClaw Manages Natively (Its Own Features)

These features are NOT OpenClaw-native — they're SwarmClaw's own system, for direct-API agents:

| Feature | Details |
|---------|---------|
| **Task Board** | Kanban tasks with status tracking |
| **Schedules + Daemon** | Background cron scheduler |
| **LangGraph Workflows** | Multi-step orchestration with checkpoints |
| **Chatrooms** | Multi-agent rooms with routing |
| **Connectors** | Discord, Slack, Telegram, WhatsApp, Teams, Matrix, email bridges |
| **Memory Graph** | Per-agent hybrid FTS5 + vector search (sqlite-vec) |
| **Plugins** | 30+ built-in session tools (shell, web, git, wallet, http, email, calendar, etc.) |
| **MCP Servers** | Model Context Protocol server attachments per agent |
| **Secrets Vault** | Encrypted credential storage |
| **Smart Deploy** | OpenClaw local/VPS/hosted deployment from UI |

---

## Voice Mode

Voice is **partially functional**:

| Direction | Implementation | Status |
|-----------|---------------|--------|
| **Input (mic → text)** | Browser Web Speech API | Works in Chrome/Edge/Safari |
| **Input (server-side)** | ElevenLabs STT or OpenAI Whisper | Via connector pipeline |
| **Output (text → speech)** | ElevenLabs API (`eleven_multilingual_v2`) | Works with API key |

**To enable**: Settings → add `elevenLabsApiKey` + set `elevenLabsEnabled: true`

**Not available**: Real-time voice conversation (no WebRTC/LiveKit integration in SwarmClaw UI).

---

## OpenClaw WebSocket + Ed25519 Explained

**Why WebSocket** (not HTTP like other providers):
- OpenClaw is a stateful runtime (agent memory, tool sessions, device connections, approval queues)
- WebSocket is full-duplex — gateway can push events TO SwarmClaw (approvals, heartbeats)
- HTTP is stateless request-response only

**Ed25519 device authentication flow**:
1. SwarmClaw has an Ed25519 keypair. `deviceId = SHA256(raw_public_key)`
2. Gateway sends a random **nonce** challenge
3. SwarmClaw signs: `"v3|deviceId|gateway-client|backend|operator|scopes|timestamp|token|nonce|platform|"` with its private key
4. Gateway verifies signature with public key → if valid → connected
5. Each nonce is unique → replay attacks impossible

SwarmClaw shares OpenClaw CLI's keypair if present (`~/.openclaw/identity/device.json`), so it appears as the same "device" to your gateway.

---

## Data Architecture

| Store | Technology | Location | Purpose |
|-------|-----------|----------|---------|
| SwarmClaw primary | SQLite (`swarmclaw.db`) | `data/swarmclaw.db` | Agents, sessions, tasks, schedules, credentials |
| SwarmClaw memory | SQLite (`memory.db`) | `data/memory.db` | Per-agent memory (FTS5 + sqlite-vec) |
| OpenClaw | SQLite (own `memory.db`) + markdown | `~/.openclaw/` | Agent memory, conversation history |
| Codex | Markdown files | `codex/` | Obsidian vault, all agents read this |

No live data sharing between SwarmClaw and OpenClaw during chat — OpenClaw agents get pure passthrough (message + agentId only).

**Backup** (VPS cron, every 6h, 7-day retention):
```bash
/home/onemind/backups/backup-swarmclaw.sh
```

---

## UI Decision: SwarmClaw vs OpenClaw Studio

**Use SwarmClaw as primary UI.** It's a strict superset of Studio's OpenClaw management capabilities.

| Studio-only advantage | Workaround in SwarmClaw |
|-----------------------|------------------------|
| Slightly cleaner single-gateway admin UX | SwarmClaw has same routes, just more options |

SwarmClaw adds: multi-gateway, 15 providers, LangGraph, task board, connectors, memory graph, plugins, Smart Deploy.

---

## OpenViking — Recommended POC

**OpenViking** (`volcengine/OpenViking`) is a context database for AI agents that has a native OpenClaw plugin.

**Benchmark** (LoCoMo10 long-range dialogue, 1,540 cases):

| Setup | Accuracy | Input Tokens |
|-------|---------|-------------|
| OpenClaw native memory | 35.65% | 24.6M |
| OpenClaw + LanceDB | 44.55% | 51.6M |
| **OpenClaw + OpenViking** | **52.08%** | **4.3M** |

**+46% accuracy, -83% token cost** over native memory.

**What it offers**:
- Hierarchical virtual filesystem (`viking://`) for memories, resources, skills
- L0/L1/L2 tiered context loading (abstract → overview → full detail)
- Directory recursive retrieval (better than flat RAG)
- Session auto-memory extraction
- Codex vault maps naturally to `viking://resources/`

**Integration path**: Run as a service on VPS (port 1933), wire up the OpenClaw memory plugin, point it at the codex directory.

```bash
pip install openviking
openviking-server  # runs on :1933
```

---

## Local Development

```bash
# SwarmClaw local
cd command-center/swarmclaw && npm install && npm run dev  # http://localhost:3456

# Setup after local start (run once to configure gateway + credentials)
python3 /tmp/swarmclaw_setup.py

# OpenClaw Studio (lightweight single-gateway admin)
cd command-center/openclaw-studio && npm install && npm run dev  # http://localhost:3000
```
