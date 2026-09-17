---
domain: 06
name: Agent-Native Infrastructure
tags: [openclaw, hermes, mem0, agentmail, agentphone, integrations, agent-tools]
status: research
updated: 2026-05-14
---

# Agent-Native Infrastructure

> Tools built *for* AI agents — not adapted from human tools. Each gives agents real-world presence: memory, identity, communication, compute, payments.

Researched: 2026-05-14

---

## ✅ Confirmed OpenClaw / Hermes Integrations

These tools have **dedicated, first-party documentation pages** for OpenClaw and/or Hermes. They are not "probably compatible" — they explicitly support and name OpenClaw/Hermes.

---

### 1. Mem0 — Long-Term Memory

**What it does:** Persistent cross-session memory for agents. Automatically extracts facts from conversations, stores them, and injects relevant memories before each turn.

| Field | Value |
|-------|-------|
| Site | https://mem0.ai |
| Pricing | Free (10K add req/mo) → $19 → $79 → $250/mo → Enterprise custom |
| Funding | $24M raised |
| GitHub | https://github.com/mem0ai/mem0 |

**OpenClaw Integration — CONFIRMED**
- Dedicated docs page: https://docs.mem0.ai/integrations/openclaw
- Plugin name: `@mem0/openclaw-mem0`
- Install: `openclaw plugins install @mem0/openclaw-mem0`
- Or chat-based setup (paste into any channel): `Setup Mem0 from mem0.ai/claw-setup`
- Requires: OpenClaw >= 2026.4.25

```json
// openclaw.json — minimal config
{
  "plugins": {
    "slots": { "memory": "openclaw-mem0" },
    "entries": {
      "openclaw-mem0": {
        "enabled": true,
        "config": {
          "apiKey": "${MEM0_API_KEY}",
          "userId": "zeus",
          "skills": {
            "triage": { "enabled": true },
            "recall": { "enabled": true, "tokenBudget": 1500, "rerank": true },
            "dream": { "enabled": true },
            "domain": "companion"
          }
        }
      }
    }
  }
}
```

**How it works (3-stage pipeline):**
1. **Triage** — Extracts durable facts after each turn (background thread, non-blocking)
2. **Recall** — Injects relevant memories before each turn (reranked, scoped: session or long-term)
3. **Dream** — Periodic consolidation: merges duplicates, resolves conflicts, prunes stale entries

**8 agent tools added:**
`memory_search`, `memory_add`, `memory_get`, `memory_list`, `memory_update`, `memory_delete`, `memory_event_list`, `memory_event_status`

**CLI commands:**
```bash
openclaw mem0 search "what languages does the user know"
openclaw mem0 search "what does Zeus prefer" --scope long-term
openclaw mem0 list --top-k 20
openclaw mem0 status
openclaw mem0 list --json  # machine-readable
```

**Self-hosted (no Mem0 key) option:**
```bash
openclaw mem0 init --mode open-source  # interactive wizard
# Supports: Ollama (local), OpenAI, Anthropic for LLM + embeddings
# Vector stores: Qdrant, PGVector, in-memory
```

**Privacy notes:**
- Credential patterns (`sk-`, `m0-`, `ghp_`, `AKIA`, `Bearer`, `password=`, etc.) are never stored — 5 independent guard layers
- `openclaw.json` stored at `0o600` (owner-read-only)
- Set `MEM0_TELEMETRY=false` to opt out of anonymous usage stats

---

**Hermes Agent Integration — CONFIRMED**
- Dedicated docs page: https://docs.mem0.ai/integrations/hermes
- Hermes is the open-source AI agent CLI by Nous Research (github.com/NousResearch/hermes-agent)
- Install Hermes: `curl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash`
- Setup: `hermes memory setup` → select `mem0` → enter API key
- Or manual: `hermes config set memory.provider mem0`

**3 agent tools added to Hermes:**
`mem0_profile`, `mem0_search`, `mem0_conclude`

**Memory timing in Hermes:**
1. Before response: cached Mem0 results from previous turn injected into system prompt (zero-latency)
2. After response: (user msg, assistant reply) pair sent to Mem0 in background thread
3. Background prefetch: Mem0 search pre-loads memories for the *next* turn while current turn completes

**Hermes config** (`~/.hermes/mem0.json`):
```json
{
  "api_key": "m0-...",
  "user_id": "zeus",
  "agent_id": "legacy",
  "rerank": true
}
```

---

### 2. AgentMail — Email Inboxes for Agents

**What it does:** Gives agents real email inboxes, accessible via API. Two-way conversations, threads, attachments, webhooks. Not a human email client — built entirely for programmatic agents.

| Field | Value |
|-------|-------|
| Site | https://agentmail.to |
| Pricing | Free (3 inboxes, 3K emails/mo) → $20/mo (10 inboxes) → $200/mo (150 inboxes) → Enterprise |
| Backed by | Y Combinator (S25) |
| SOC 2 | Compliant |
| MCP server | Yes |

**OpenClaw Integration — CONFIRMED**
- Dedicated docs page: https://docs.agentmail.to/integrations/openclaw
- Install via clawhub: `npx clawhub@latest install agentmail`
- Or via skills CLI: `npx skills add agentmail-to/agentmail-skills`
- Docs note: "Works with Claude Code, Cursor, **OpenClaw**, Codex, and other compatible agents"

```json
// openclaw.json
{
  "skills": {
    "entries": {
      "agentmail": {
        "enabled": true,
        "env": { "AGENTMAIL_API_KEY": "your-api-key-here" }
      }
    }
  }
}
```

**Get API key:** https://console.agentmail.to

**Verify installation:**
```bash
openclaw skills list --eligible  # should show "agentmail"
```

**Features included in the OpenClaw skill:**
- Inbox management (create scalable inboxes on-demand)
- Send emails (text + HTML)
- Thread management
- Attachments (Base64)
- Drafts (human-in-the-loop approval before sending)
- Pods (multi-tenant isolation)
- Real-time events (WebSocket + webhook)

**Real-time email → OpenClaw notifications:**
Register AgentMail webhook pointing to your OpenClaw webhook URL so agents are notified immediately when emails arrive:
```bash
curl -X POST -H "Authorization: Bearer $AGENTMAIL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-openclaw-webhook-url", "event_types": ["message.received"]}' \
  https://api.agentmail.to/v0/webhooks
```

**One Mind application:** Each agent (Legacy, Zenith, Unity, Eternal, Sovereignty, Sentinel) gets its own `@agentmail.to` inbox. Legacy can route emails to domain agents. Eternal gets a business email for client comms. Zenith gets a health tracking inbox.

---

### 3. AgentPhone — Phone Numbers for Agents

**What it does:** Gives agents real US/CA phone numbers that can send SMS, receive calls, and hold voice conversations — all via a single unified webhook. Voice calls are transcribed live and delivered as text; agent replies as text and AgentPhone handles speech synthesis.

| Field | Value |
|-------|-------|
| Site | https://agentphone.ai |
| Pricing | $3/number/month (usage-based) · $150/mo Pro (iMessage-enabled, RCS fallback, no per-msg cost) · Enterprise custom |
| Backed by | Y Combinator |
| Status | Free during beta (no credit card) |
| MCP server | Yes (`npx agentphone-mcp`) |

**OpenClaw Integration — CONFIRMED (explicitly named)**

From the AgentPhone website:
> *"Works with Claude Code, Cursor, **OpenClaw**, Windsurf, and any MCP client."*

There is also a dedicated use case section on the homepage called:
> **"Give your OpenClaw/Hermes agent a phone number"**

This is the only tool in this list that explicitly names both OpenClaw AND Hermes by name on its public homepage.

**MCP config (works natively with OpenClaw MCP support):**
```json
{
  "mcpServers": {
    "agentphone": {
      "command": "npx",
      "args": ["agentphone-mcp"]
    }
  }
}
```

**API quickstart:**
```python
from agentphone import AgentPhone

client = AgentPhone(api_key="ap_...")

# Create agent persona
agent = client.agents.create(
    name="Legacy",
    voice_mode="hosted",
    system_prompt="You are Legacy, Zeus's AI orchestrator."
)

# Buy a number
number = client.numbers.buy(agent_id=agent.id)
print(f"Legacy's number: {number.phone_number}")

# Make an outbound call
call = client.calls.create(
    agent_id=agent.id,
    to_number="+14155557890"
)
```

**REST API (one call to get a number):**
```bash
# Create agent
curl -X POST "https://api.agentphone.ai/v1/agents" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{"name": "Legacy"}'

# Provision number
curl -X POST "https://api.agentphone.ai/v1/numbers" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{}'

# Attach number to agent
curl -X POST "https://api.agentphone.ai/v1/agents/{agent_id}/numbers" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{"numberId": "NUMBER_ID"}'
```

**FAQ from their docs:**
- *"How is this different from Twilio?"* → "Twilio gives you low-level telecom building blocks. AgentPhone is designed specifically for AI agents, with a unified webhook for messages and calls, automatic conversation threading, real-time transcription, and MCP support already built in."
- Port existing numbers from Twilio, Telnyx, etc. for free on any plan.

**One Mind application:** Legacy gets a primary number for direct contact with Zeus. Each domain agent gets a number for their area (Zenith for health check-ins, Eternal for business calls, Sovereignty for infra alerts). 2FA inbox for agent autonomy workflows.

---

## 🔧 MCP-Compatible (Work With OpenClaw via MCP Protocol)

These tools expose an MCP server. Since OpenClaw supports MCP natively, they work — but they don't name OpenClaw explicitly in their docs.

---

### 4. Kapso — WhatsApp for Agents

**What it does:** Official Meta Business Partner API for WhatsApp. Agents get WhatsApp access with full two-way messaging, webhooks, flows (mini-apps), and a workflow builder.

| Field | Value |
|-------|-------|
| Site | https://kapso.ai |
| Pricing | Free (2,000 msgs/mo, $2 AI credits) → paid tiers |
| MCP/Skills | `npx skills add gokapso/agent-skills` (5 skills) |
| Docs | https://docs.kapso.ai |

**Agent skill install:**
```bash
npx skills add gokapso/agent-skills
```

This installs 5 WhatsApp skills compatible with Claude Code, Cursor, Codex, and any skills-compatible agent (including OpenClaw via skills support).

**OpenClaw compatibility:** Via skills protocol (same as AgentMail's secondary install method). Not explicitly named in Kapso docs but uses the same `npx skills add` mechanism that OpenClaw supports.

**Features:** Sandbox mode (test without verification), webhook delivery tracking, API logs, message status tracking, serverless JS functions, broadcast campaigns.

---

### 5. x402 — HTTP-Native Agent Payments

**What it does:** Open protocol (Linux Foundation) for internet-native micropayments. When an agent makes an HTTP request that requires payment, the server responds with HTTP 402, agent pays instantly with stablecoins, access granted — no accounts, no API keys, no KYC.

| Field | Value |
|-------|-------|
| Site | https://x402.org |
| Backed by | Coinbase |
| Protocol fees | Zero |
| Volume (last 30d) | $24.24M across 75.41M transactions |
| Networks | USDC on Base (primary), multi-chain expanding |

**One-line integration (seller side):**
```js
app.use(paymentMiddleware({
  "GET /weather": {
    accepts: [...],
    description: "Weather data",
  },
}));
```

**How it works for agents:**
1. Agent sends HTTP request → receives `402 Payment Required`
2. Agent pays instantly with stablecoins (no signup, no approval, no API key)
3. API access granted

**OpenClaw compatibility:** Protocol-level. Any agent that can make HTTP requests can use x402. No specific OpenClaw integration page, but Cal.com already accepts x402 for calendar bookings, and AgentMail supports x402 for inbox creation (email inboxes purchasable over HTTP).

**One Mind application:** Agents autonomously pay for tools, APIs, data feeds, and services. Grid pays for cloud resources. Eternal pays for market data. Legacy pays for research tools.

---

## 🌐 Agent Infrastructure (API-Based, No MCP — Integrate via Skills/Webhooks)

---

### 6. E2B — Secure Compute Sandboxes

**What it does:** Isolated cloud sandboxes (microVMs) where agents can safely execute code. Each sandbox is a real Linux environment. Used by Manus, Perplexity, HuggingFace, Gumloop.

| Field | Value |
|-------|-------|
| Site | https://e2b.dev |
| Pricing | Free (100 credits) → $150/mo Pro + usage ($0.000028/s per 2 vCPU) |
| Funding | $21M Series A |
| GitHub | https://github.com/e2b-dev/e2b |

**Use cases:** Grid/Sovereignty agents running untrusted code. Code interpreter for Eternal (financial modeling). Any agent task that needs a real execution environment without risking the VPS.

---

### 7. Browserbase — Agent Browser Sessions

**What it does:** Real production-grade browsers for agents to use the web. Includes an identity layer (authenticate as a human without getting blocked), session replay, and observability. Stagehand is their open-source browser automation framework.

| Field | Value |
|-------|-------|
| Site | https://browserbase.com |
| Free tier | Yes |
| MCP | Via Stagehand |

**Products:**
- **Browsers** — Real browser sessions, anti-bot identity layer
- **Stagehand** — Open-source AI browser automation framework (most popular)
- **Director** — Visual UI for building browser agents
- **Web Data APIs** — Search/fetch APIs to retrieve structured web data

**One Mind application:** Legacy/Trinity for deep web research. Sovereignty for monitoring services. Eternal for market research.

---

### 8. Mem0 Open-Source (Self-Hosted) — Already Covered Above

See Section 1. The self-hosted mode uses Ollama + Qdrant (fully local, no external API) or OpenAI + PGVector.

---

## 🎙️ Voice

### 9. Vapi — Voice AI Agent Platform

**What it does:** Build, test, and deploy voice AI agents. Agents get a voice persona, handle inbound/outbound calls, with real-time transcription and STT/TTS.

| Field | Value |
|-------|-------|
| Site | https://vapi.ai |
| Pricing | Usage-based ($0.05/min average, varies by provider) |
| Integration | REST API + SDKs |

**Works with AgentPhone:** Vapi handles the voice AI layer; AgentPhone handles the phone number routing. Stack them: AgentPhone number → Vapi agent for voice handling.

**One Mind application:** Each agent gets a distinct voice persona via ElevenLabs voices, deployed through Vapi. Legacy gets one voice, Zenith another, etc.

---

### 10. ElevenLabs Conversational AI

**What it does:** Voice agent platform with distinct voice identities. Each agent can have a unique voice — different timbre, accent, personality.

| Field | Value |
|-------|-------|
| Site | https://elevenlabs.io/conversational-ai |
| Pricing | Usage-based |

**One Mind application:** Assign each agent a distinct voice. Legacy sounds authoritative. Zenith sounds energetic. Unity sounds warm.

---

## 📊 Observability

### 11. AgentOps — Agent Tracing & Cost Monitoring

**What it does:** Trace, debug, and deploy reliable AI agents. Tracks LLM calls, tool use, multi-agent interactions. Session replay. Cost monitoring across 400+ LLMs.

| Field | Value |
|-------|-------|
| Site | https://agentops.ai |
| Free tier | Yes |
| Integrations | OpenAI, CrewAI, Autogen, 400+ LLMs |

**One Mind application:** Track token spend per agent. Debug tool call chains. Replay failed sessions. Monitor when agents hit errors.

---

## 🔐 Auth

### 12. Scalekit — Auth Stack for Agentic Apps

**What it does:** OAuth 2.1 for MCP servers + delegated user consent for agents acting on behalf of users. Purpose-built for the agentic era.

| Field | Value |
|-------|-------|
| Site | https://scalekit.com |
| Funding | $5.5M Seed |
| Use case | If you expose agents as MCP servers to other users/tools |

**One Mind application:** If One Mind agents are ever exposed as MCP servers to third parties, Scalekit handles the auth layer so external users can authorize agents to act on their behalf securely.

---

## 🗺️ OpenClaw/Hermes Compatibility Summary

| Tool | Explicitly names OpenClaw | Explicitly names Hermes | MCP Compatible | Install Method |
|------|:---:|:---:|:---:|---|
| **Mem0** | ✅ Dedicated docs page | ✅ Dedicated docs page | ✅ | `openclaw plugins install @mem0/openclaw-mem0` |
| **AgentMail** | ✅ Dedicated docs page | ❌ | ✅ | `npx clawhub@latest install agentmail` |
| **AgentPhone** | ✅ Named on homepage + FAQ | ✅ Named on homepage | ✅ | `npx agentphone-mcp` (MCP server) |
| **Kapso** | ❌ | ❌ | Via skills | `npx skills add gokapso/agent-skills` |
| **x402** | ❌ | ❌ | Protocol-level | HTTP middleware |
| **E2B** | ❌ | ❌ | Via API | REST API / Python SDK |
| **Browserbase** | ❌ | ❌ | Via Stagehand MCP | REST API / Stagehand |
| **Vapi** | ❌ | ❌ | Via API | REST API |
| **AgentOps** | ❌ | ❌ | Via SDK | Python/JS SDK |
| **Scalekit** | ❌ | ❌ | MCP OAuth | REST API |

---

## 🚀 Recommended Priority for One Mind

| Priority | Tool | Reason |
|----------|------|--------|
| **1** | Mem0 | Biggest immediate upgrade — fixes stateless memory. One command: `openclaw plugins install @mem0/openclaw-mem0`. Works with all 6 agents individually scoped. |
| **2** | AgentPhone | Free beta · YC-backed · Explicitly supports OpenClaw + Hermes · MCP native · $3/mo when it launches. Legacy gets a number Zeus can call. |
| **3** | AgentMail | Per-agent email identities. Install: `npx clawhub@latest install agentmail`. Free tier has 3 inboxes (start with Legacy, Eternal, Sovereignty). |
| **4** | x402 | No account needed. When agents need to autonomously pay for APIs or data, this is the cleanest path. |
| **5** | E2B | For Grid/Sovereignty code execution tasks. Keeps agent workloads off the VPS. |
| **6** | Browserbase | For Legacy/Trinity deep research. Replaces basic curl with real browser sessions. |

---

## 🔗 Source URLs

- AgentMail: https://agentmail.to · Docs: https://docs.agentmail.to · OpenClaw: https://docs.agentmail.to/integrations/openclaw
- AgentPhone: https://agentphone.ai · Docs: https://docs.agentphone.ai · Pricing: https://agentphone.ai/pricing
- Mem0: https://mem0.ai · Pricing: https://mem0.ai/pricing · OpenClaw: https://docs.mem0.ai/integrations/openclaw · Hermes: https://docs.mem0.ai/integrations/hermes
- Kapso: https://kapso.ai · Docs: https://docs.kapso.ai
- x402: https://x402.org · Docs: https://docs.x402.org
- E2B: https://e2b.dev · Pricing: https://e2b.dev/pricing
- Browserbase: https://browserbase.com
- Vapi: https://vapi.ai
- AgentOps: https://agentops.ai
- Scalekit: https://scalekit.com
