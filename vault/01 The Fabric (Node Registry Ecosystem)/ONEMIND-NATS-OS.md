# One Mind NATS OS
## An Event-Driven AI Operating System for Your Entire Digital Life

**Date:** 2026-06-05  
**Author:** Zeus Delacruz  
**Status:** Planning — Pre-Architecture

---

## The Core Idea

You already have a NATS bus running. You have 10+ Hermes agents, n8n, Node-RED, Home Assistant, SaaS apps, subscribers, and a growing AI stack. Right now these systems talk *around* each other — HTTP, webhooks, MCP tools, direct API calls.

The idea: **make NATS the central nervous system of your entire digital life.** Every service, agent, app, and device becomes a node on the bus. Every action, decision, event, and state change flows through one place. An AI meta-agent sits at the top and governs, routes, analyzes, and learns from everything in real time.

This is not just an infrastructure upgrade. This is a **fundamentally different way to give AI agents experience** — not memory summaries or RAG lookups, but a living, ordered, replayable stream of everything that happened. It surpasses current memory architectures entirely.

---

## Why This Changes Everything

### Current problem with AI memory
Most AI memory systems work like this:
- Event happens → summarize it → store as a fact → retrieve later
- This is lossy, decontextualized, and static

### What NATS JetStream gives you instead
- Event happens → published to stream → stored in order, with timestamp and subject
- Agent can **replay** the stream — not recall a summary but relive the actual sequence
- Causality is preserved — you can see what triggered what
- Patterns emerge over time that no single memory entry could capture

This is how humans build real experience. Not facts — sequences. Not memory — history.

---

## The Architecture

### Layer 1 — The Bus (NATS + JetStream)

NATS is already running. JetStream needs to be enabled for persistence.

```
nats-server --jetstream
```

JetStream turns NATS from a message broker into a **durable event log**. Streams store messages for as long as you want. Consumers can replay from any point in time.

**Core streams to create:**

| Stream | Subjects | Purpose |
|--------|----------|---------|
| `AGENT_EVENTS` | `agent.>` | All agent actions, decisions, outputs |
| `SERVICE_EVENTS` | `service.>` | n8n, Node-RED, HA, SaaS activity |
| `SYSTEM_EVENTS` | `system.>` | Device state, health, metrics |
| `HUMAN_EVENTS` | `human.>` | Your inputs, commands, decisions |
| `SAAS_EVENTS` | `saas.>` | External SaaS webhooks and triggers |
| `GOVERNANCE` | `govern.>` | Meta-agent decisions, overrides, alerts |

---

### Layer 2 — The Nodes (Everything on the Bus)

Every service gets a NATS client and a subject identity. Here is your current stack mapped to the bus:

#### AI Agents (Hermes)
```
agent.ceo.>          ← CEO agent events
agent.cfo.>          ← CFO agent events
agent.ops.>          ← Ops agent events
agent.coder.>        ← Coder agent events
agent.growth.>       ← Growth agent events
agent.sentinel.>     ← Sentinel agent events
agent.eternal.>      ← Eternal agent events
agent.sovereignty.>  ← Sovereignty agent events
agent.unity.>        ← Unity agent events
agent.zenith.>       ← Zenith agent events
```

Each agent publishes to its subject when it acts, and subscribes to subjects it needs to react to.

**How to wire Hermes agents to NATS:**
- Add `nats.py` (Python) or `nats.js` (Node) client to each agent process
- On startup: subscribe to `agent.<name>.in.>` for incoming tasks
- On action: publish to `agent.<name>.out.>` with result payload
- On decision: publish to `agent.<name>.decision` with reasoning + outcome

#### Automation Services
```
service.n8n.>          ← n8n workflow events
service.node-red.>     ← Node-RED flow events
service.ha.>           ← Home Assistant state changes
service.notify.>       ← Notification router events
```

**n8n** — has a native NATS trigger and action node. Configure NATS server address, map workflows to subjects. n8n workflows become NATS-native services.

**Node-RED** — install `node-red-contrib-nats`. Each flow can publish/subscribe. Node-RED becomes a visual NATS node editor.

**Home Assistant bridge** — already a LaunchAgent (`io.onemind.ha-bridge.plist`). Extend it to translate HA state changes into `service.ha.state.<entity>` publish calls.

#### SaaS Apps
```
saas.todoist.>       ← task created, completed, updated
saas.calendar.>      ← events, reminders
saas.email.>         ← incoming, sent, flagged
saas.github.>        ← commits, PRs, issues
saas.slack.>         ← messages, mentions
saas.<any>.>         ← anything with a webhook
```

**How SaaS connects:** A thin HTTP→NATS bridge (you already have the `graphite-bridge` pattern). SaaS app fires a webhook → bridge receives it → publishes to `saas.<app>.<event>` subject. One bridge, all SaaS apps.

#### Devices & System
```
system.mac.>         ← CPU, memory, network, app activity
system.iphone.>      ← location, activity, battery, focus mode
system.ha.>          ← home sensors, switches, automations
```

---

### Layer 3 — The Meta-Agent (Governance + Intelligence)

This is the most important layer. The meta-agent is an AI process that:

1. **Subscribes to `>`** — sees every single message on the entire bus
2. **Routes** — decides which agent or service should handle incoming tasks
3. **Governs** — enforces rules about what agents can and cannot do
4. **Analyzes** — watches streams over time for patterns and anomalies
5. **Intervenes** — can pause, reroute, escalate, or cancel any in-flight action
6. **Learns** — JetStream replay means it can study history to improve routing

```
govern.route.>       ← routing decisions
govern.approve.>     ← approvals required before action executes
govern.alert.>       ← anomalies, escalations
govern.audit.>       ← full decision log
```

**Meta-agent responsibilities:**

| Function | What it does |
|----------|-------------|
| Task Router | Receives tasks, assigns to best agent based on context + history |
| Policy Enforcer | Blocks actions that violate rules (cost limits, scope, permissions) |
| State Aggregator | Maintains current state of all nodes — who is active, busy, errored |
| Pattern Detector | Watches JetStream history, surfaces recurring patterns |
| Experience Engine | Replays past events to give agents contextual experience |
| Escalation Handler | Knows when to wake you up vs handle autonomously |

---

### Layer 4 — Synadia Agent SDK (Native Bus Agents)

The Synadia Agent SDK (NEX — NATS Execution Engine) takes this further. Instead of agents that *use* NATS, you get agents that *live* on NATS:

- Agents deployed as **NATS services** — discoverable via `$SRV` subjects
- Built-in health, stats, and ping per service
- Agents scale automatically via **queue groups** — multiple instances, NATS load balances
- No separate orchestration needed — the bus *is* the orchestrator
- Service mesh without Kubernetes

**Why this is stronger than MCP:**

| MCP | NATS (NEX) |
|-----|-----------|
| Point-to-point connections | One bus |
| Tool per connection | Subject per capability |
| Synchronous request/response | Async + persistent |
| No observability built in | Full audit trail |
| Hard to scale | Queue groups auto-scale |
| Separate server per tool | Everything on the bus |

Teaching agents NATS protocol once gives them access to everything on the bus. No new tool registrations. No new MCP servers. No new connection management. Subscribe to what you need, publish what you do.

---

## The Experience Layer (This is the breakthrough)

With JetStream enabled and everything on the bus, your AI agents have access to something no memory system can replicate: **a complete, ordered, causal record of your digital life.**

**What this enables:**

- **Pattern recognition across systems** — your Todoist tasks correlate with HA activity, iPhone location, and agent decisions. An AI can see your actual behavioral patterns, not just what you told it.
- **Causal chains** — agent can trace: "This SaaS webhook triggered this n8n flow which caused this agent decision which produced this outcome." Full lineage.
- **Temporal context** — not "Zeus prefers X" but "Zeus preferred X on Mondays between 9-11am when he was at home. On travel days he chose Y."
- **Drift detection** — your habits shift over time. The stream captures drift. Agents adapt.
- **Replay for onboarding new agents** — new agent joins, replays last 30 days of streams, understands your system with full context before it ever acts.

This is **experience, not memory.** The difference is the same as the difference between reading a biography and living a life.

---

## Governance Architecture

Two layers of governance:

### Infrastructure-level (NATS server)
- Subject authorization — agent X can only pub/sub to its own subjects
- Rate limiting — agent cannot fire more than N messages per second
- Message size limits — prevents runaway payloads
- Account isolation — separate NATS accounts per trust domain

### Intelligence-level (Meta-agent)
- Watches `govern.approve.>` — certain actions require meta-agent sign-off before execution
- Cost-aware routing — knows which operations are expensive, routes accordingly
- Anomaly detection — agent deviating from normal pattern triggers review
- Human-in-the-loop escalation — meta-agent decides what wakes Zeus vs handles itself

Together these give you **layered governance** — hard rules at infrastructure, soft intelligent rules above.

---

## Implementation Phases

### Phase 1 — Enable JetStream + Define Schema
- Enable JetStream on existing NATS server
- Define and create the 6 core streams
- Establish subject naming convention
- Document the schema

### Phase 2 — Wire Automation Nodes
- Connect n8n via native NATS node
- Connect Node-RED via contrib plugin
- Extend HA bridge to publish state events
- Build/configure HTTP→NATS SaaS bridge

### Phase 3 — Wire Hermes Agents
- Add NATS client to each agent process
- Implement pub/sub on agent subjects
- Agents publish decisions and actions to stream
- Test routing via direct subject calls

### Phase 4 — Build Meta-Agent
- Create meta-agent process (subscribes to `>`)
- Implement routing logic
- Implement governance/approval layer
- Connect to JetStream for history access

### Phase 5 — Synadia NEX Migration
- Evaluate migrating agents to NEX native services
- Implement service discovery
- Enable queue group auto-scaling
- Replace remaining MCP tools with NATS subjects

### Phase 6 — Experience Engine
- Build JetStream replay tooling
- Implement pattern detection on streams
- Enable new agent onboarding via stream replay
- Build dashboard for stream visualization

---

## What You Have Today vs What You're Building

| Today | NATS OS |
|-------|---------|
| Agents run in silos | All agents on one bus |
| Memory = summaries | Memory = ordered event streams |
| MCP tools per connection | One bus, subject per capability |
| Webhooks per SaaS app | One bridge, all SaaS on bus |
| No causal history | Full causal audit trail |
| Governance = code logic | Governance = infrastructure + meta-agent |
| Agents learn from prompts | Agents experience from streams |

---

## The Vision

Your Mac becomes a **sovereign AI node** — a self-contained intelligence platform where every piece of your digital life flows through a governed, persistent, observable bus. The meta-agent is the intelligence layer above it all: routing, learning, governing, and experiencing everything in real time.

This is not a productivity tool. This is infrastructure for a new kind of human-AI collaboration — where the AI doesn't just respond to your requests but has genuine context, history, and experience of your life, enforced by hard governance rules you control.

You own the bus. You own the streams. You own the agents. Nothing leaves unless you say so.

---

*One Mind Framework — NATS OS Initiative*  
*Zeus Delacruz — 2026*
