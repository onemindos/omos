---
id: so-01-20260527-fleet-harness-decision-framework
title: "Fleet Harness Decision Framework — When to Use Which Agent Runtime"
type: doctrine
status: active
created: 2026-05-27
modified: 2026-05-27
owner: zeus
pillar: SO
domain: "01 Fleet"
tags: [type/doctrine, zone/so, domain/so, status/active, infra/hermes, infra/openclaw, agents/fleet, architecture/canonical]
depth: 0
---

# Fleet Harness Decision Framework

> When to use which agent runtime for each domain.

## The Rule

Every harness in the fleet MUST:
1. Register as a NATS micro-service ($SRV)
2. Report heartbeats via `agents.hb.<agent>`
3. Accept work via `agents.prompt.<agent>` (Synadia Agent SDK)
4. Stream telemetry to Langfuse
5. Have persistent sessions (own DB)
6. Self-manage until escalation needed

---

## Decision Matrix

| Use Case | Harness | Why |
|----------|---------|-----|
| **Web-heavy business** (SaaS, marketing, sales, customer support) | OpenClaw | TS-native, great at web tools, canvas, messaging platforms |
| **Physical/IoT** (greenhouse, barn, drones, robots) | Hermes | Python sensors, CV, ML, hardware SDKs |
| **Coding agent** | Claude Code / Codex (summoned via ACP) | Purpose-built, not a gateway |
| **Personal domain** (home, family, health) | Legacy (Hermes main) | Needs full access, deep memory, trust |
| **New business / client project** | OpenClaw gateway (spin up, tune, hand off) | Isolated, portable, TS web ecosystem |
| **Edge node** (remote location, limited hardware) | Hermes lightweight (Pi/Jetson) | Python, small footprint, NATS leaf |
| **Voice interface** | Qwen2.5-Omni (specialist) | Audio-native, summoned via ACP |
| **Creative / media gen** | ComfyUI / specialist | GPU-bound, session-based |

---

## Fleet Tiers

### Tier 1 — BRAIN (always on, MBP)
- **Hermes (Legacy)** — orchestrator, personal domain, fleet manager
- Always-on, NATS-native, full tool access, self-improving
- Sees all harnesses, reviews all work, manages escalations

### Tier 2 — BUSINESS HARNESSES (per-company)
- **OpenClaw instances** — one per business entity
- Each registers on NATS via Synadia Agent SDK
- Each has own sessions, memory, skills tuned to that business
- Legacy sees all, injects work, reviews output
- Portable: can hand off to employees or buyers

### Tier 3 — PHYSICAL NODES (local hardware)
- **Hermes lightweight** — greenhouse, office, barn, drones
- Python for sensors/ML/CV, NATS leaf for connectivity
- Autonomous until escalation needed
- Works offline (local JetStream store-and-forward)

### Tier 4 — SPECIALIST AGENTS (summoned on demand)
- Claude Code / Codex / OpenCode — coding
- Qwen2.5-Omni — voice interface
- ComfyUI — image/video generation
- Summoned via ACP/NATS, session-based, terminate when done

---

## Scaling Pattern

```
Year 1 (now):
  Legacy (Hermes) + 1 OpenClaw (OneMind LLC)
  Langfuse observing both
  NATS connecting them

Year 2:
  Legacy + 3 OpenClaw businesses + greenhouse (Hermes) + office (Hermes)
  All on NATS, all in Langfuse
  Flutter app shows full fleet

Year 3:
  Legacy + 5-10 business harnesses + 3 physical nodes + drone fleet
  Each business spun up in hours (clone template, tune, deploy)
  Each physical node autonomous but reports to Legacy
  Full voice control across all of it
```

---

## Key Principles

1. **New business = new gateway.** Never mix business data in one harness.
2. **Physical location = own node.** Greenhouse agent runs ON greenhouse hardware.
3. **Legacy stays light.** Orchestrate, don't execute grunt work.
4. **NATS is the only bus.** Every harness speaks NATS or doesn't exist.
5. **Langfuse is the only observer.** Every harness reports or is blind to fleet.
6. **$SRV is the only registry.** Every harness is discoverable or is dead.
7. **Sessions persist until explicitly closed.** Work continues across days/weeks.
8. **Language is per-domain; the wire is universal.** Each harness's language is chosen by ecosystem fit (Python for physical/sensors/CV/ML, TS for web/messaging/SaaS), never globally. The runtime language is an implementation detail — the only thing that must be language-agnostic is the fabric (NATS subjects + JetStream + a common event schema). Any runtime that speaks the wire protocol is interchangeable to the brain. The tax for this freedom: the event schema, auth, and Langfuse tracing must be identical across every language.
9. **No harness is a monolith. Legacy is a constellation.** No single process is memory + watchers + realtime voice + reasoning + tools at once. Legacy is composed of single-purpose services (Presence, Cortex, watchers, memory, security) unified by NATS. The runtime reasons; sidecars sense and speak. See `DUAL-PROCESS-PRESENCE-ARCHITECTURE.md`.
10. **The conversation lives in the fabric, not in any device.** A realtime session binds to a routable *node* (speaker/mic/camera) chosen by your current zone. The session persists; only the I/O endpoint moves. This decoupling is the omnipresence primitive — Legacy follows you room to room, home to office, without losing the thread.
