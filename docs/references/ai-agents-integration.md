# `integrations/ai-agents/`

AI agents, agentic frameworks, and the agent communication layer for One Mind OS.
All agents target the **Synadia Agent Protocol v0.3** — every agent is a NATS Micro
service on `agents.prompt.<name>.zeus.<session>`.

---

## Repos

### `synadia/` — The Agentic Layer ✅ ACTIVE
The official **Synadia Agent Protocol for NATS** SDK and tooling. This is the
standard wire format for all One Mind Hive agents. Contains:

| Subfolder | What |
|-----------|------|
| `synadia-agents` | Caller + Host SDKs (TypeScript + Python). Pre-built channels for Claude Code, OpenClaw, Hermes, OpenCode, PI, DeerFlow, Flue. |
| `hermes-nats-gateway` | Hermes (Legacy) plugin — puts this agent on NATS as `agents.prompt.hermes.zeus.*`. Install via `hermes plugins install`. |
| `synadia-agents-labs` | Courseware — echo agents, tool-calling, LLM agents, edge-sensor capstone. Reference only. |
| `synadia-agent-sdk-docs` | Wire protocol spec v0.3. Source of truth for subject format and message framing. |
| `remarkable-nats-agents` | AI chat over NATS on reMarkable e-paper tablet. Staged until hardware available. |

**Discovery:**
```bash
nats req '$SRV.INFO.agents' '' --replies=0 --timeout=2s
```

---

### `orp` — Protocol Bridge ✅ ACTIVE
**Language:** Rust | **Size:** 45MB binary | **Tests:** 1,500+

Open Reality Protocol — fuses 40+ external protocols into a single
cryptographically-signed real-time data stream. Includes TAK Protocol v1,
MAVLink, MQTT, and more. mTLS, OIDC, AES-256-GCM, Argon2id. This is the
**protocol layer** beneath the agent bus — not an agent itself.

Use when: a new protocol needs bridging into NATS that omos-bridge doesn't
cover yet.

---

### `NyxID` — API Credential Gateway ✅ ACTIVE
**Language:** Rust | **Role:** MCP server / API proxy for agents

API credential management and routing proxy for AI agents. Agents request
credentials through NyxID rather than holding API keys directly. Connects
agents to any external API with centralized credential governance.

Not an identity provider — it's an API access proxy that sits between agents
and external services.

---

### `cordum` — Agent Governance ✅ ACTIVE
**Language:** Node.js | **Role:** Compliance firewall for AI agents

Control plane for AI agent safety and compliance. Governs what agents are
allowed to do before they do it — policy enforcement, audit trail, safety
constraints. Helm chart available for containerized deployment.

Sits between the Synadia Agent bus and any agent that needs governed access
to sensitive systems (home automation, financial APIs, etc.).

---

### `infermesh` — LLM Inference Mesh 🟡 STAGED (Phase 3)
**Language:** Go | **Role:** Geographic LLM load balancing over NATS

Routes LLM inference requests across GPU nodes, cloud API adapters, and edge
devices via NATS. OpenAI-compatible HTTP proxy — swap `baseURL`, zero code
changes. Automatic failover across regions.

**Activate when:** second inference node available (VPS GPU, second Mac, etc.).
Currently a single LiteLLM node — load balancing across 1 node is meaningless.

---

### `iron-concept` — Open-Source SCADA 🟡 STAGED (Phase 3+)
**Language:** Rust (edge) + Elixir/Phoenix (UI) | **Status:** Architecture phase, pre-v0.1

Industrial automation platform: Rust edge agent → NATS JetStream → Phoenix
LiveView dashboard → TimescaleDB historian. Modbus, OPC-UA, S7 protocol support.
Apache 2.0.

**Activate when:** industrial/environmental hardware arrives (PLCs, Modbus
devices, greenhouses, water systems). No runnable code yet — all specs and ADRs.

---

## Subject Namespace

All agents use Synadia Agent Protocol v0.3:

```
agents.prompt.<token>.<owner>.<session>   # invoke agent
agents.status.<token>.<owner>.<session>   # health check
agents.hb.<token>.<owner>.<session>       # heartbeat / liveness
```

## One Mind Hive — Agent Registry

| Agent | Token | Status |
|-------|-------|--------|
| Legacy (Hermes) | `hermes` | Install hermes-nats-gateway |
| OpenClaw | `oc` | Install @synadia-ai/nats-channel |
| Claude Code | `cc` | Install claude-channel-nats |
| Spartan | `spartan` | Build — Python AgentService |
| Oracle | `oracle` | Build — Python AgentService |
| Empire | `empire` | Build — Python AgentService |
| Eden | `eden` | Build — Python AgentService |
| Guardian | `guardian` | Build — Python AgentService |
