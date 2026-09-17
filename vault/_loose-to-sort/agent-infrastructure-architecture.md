---
id: so-02-20260527-agent-infra-architecture
title: "Agent Infrastructure Architecture — Node Placement Doctrine"
type: doctrine
status: active
created: 2026-05-27
modified: 2026-05-27
owner: zeus
pillar: SO
domain: "02 Infrastructure (The Body)"
ai_summary: "Architecture decision doctrine for agent placement across the fleet — when to deploy local agents vs SSH-managed nodes, NATS as nervous system, hybrid topology by use case."
tags: [type/doctrine, zone/so, domain/so, status/active, infra/nats, infra/agents, framework/core]
related: ["[[ONE-MIND-NODE-SYSTEM]]", "[[TAILSCALE-SSH]]", "[[synadia-research-summary]]"]
depth: 1
---

# Agent Infrastructure Architecture — Node Placement Doctrine

> **"Don't put agents everywhere just because you can. An agent on a machine is infrastructure to maintain."**
>
> By Zeus Delacruz | Created: May 2026

---

## Context

The OneMind OS fleet spans multiple machine types: local workstations, cloud VPS (Lightsail), IoT/robots, and physical facility nodes. The question: **where do we place AI agents vs. manage centrally via SSH?**

### Current State

- **MBP M2 Max** = brain (Legacy AI on Hermes, NATS hub, central command)
- **Fleet**: omos-core, omos-sentinel, omos-titan (Lightsail), omos-dell (local), omos-ha (offline)
- **Network**: Tailscale mesh + SSH, NATS leaf nodes to NGS/local hub
- **Agent Stack**: Hermes (individual/profiles), OpenClaw (teams), Harmony (orchestration)
- **NATS Agent SDK**: Active — used for SaaS action workers and inter-agent communication

---

## The Rule: Use Case Determines Topology

| Scenario | Pattern | Why |
|----------|---------|-----|
| **Remote facility / another location** | Local agent on-machine, NATS for coordination | Latency, autonomy during network loss, local sensor/actuator loops |
| **Cloud workers / batch compute** | SSH + Legacy manages remotely | No local intelligence needed — they're stateless execution nodes |
| **Smart home / physical automation** | Local agent, NATS events | Sub-second response, local failover if WAN drops |
| **Dev/staging servers** | SSH from central, no local agent | Disposable, no persistent state worth an agent |
| **Second command center (e.g. Office, Barn)** | Local Hermes profile + NATS leaf | Full autonomy + syncs back to central |

---

## When to Put an Agent on the Node

1. It needs to react in <1s to local events (automation, robotics, sensors)
2. It must function autonomously if WAN/Tailscale drops
3. It has its own persistent context (a facility, a team workspace)

## When to Keep SSH + Central Management

1. The node is a compute worker (stateless tasks)
2. You're doing deployments/maintenance (Ansible-style)
3. The workload is request/response with no local state

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│  CENTRAL (MBP M2 Max)                               │
│  Legacy AI + Hermes (default profile) + NATS Server │
│  SSH access to all fleet nodes via Tailscale        │
│  OpenClaw orchestrator (team dispatch)              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─── FACILITY NODE (Barn/Office/Remote) ────┐     │
│  │ Local Hermes (own profile)                 │     │
│  │ NATS leaf → central hub                    │     │
│  │ Local automations (HA, sensors, robots)    │     │
│  │ Can operate fully offline                  │     │
│  └────────────────────────────────────────────┘     │
│                                                     │
│  ┌─── CLOUD WORKER (Lightsail/VPS) ──────────┐     │
│  │ No local agent                             │     │
│  │ SSH via Tailscale for management           │     │
│  │ NATS leaf for task pub/sub                 │     │
│  │ Netdata for monitoring                     │     │
│  │ Stateless — Legacy dispatches work via SSH │     │
│  └────────────────────────────────────────────┘     │
│                                                     │
│  ┌─── ROBOT/IoT NODE ────────────────────────┐     │
│  │ Lightweight agent (nats.swift / embedded)  │     │
│  │ NATS leaf for commands + telemetry         │     │
│  │ Must react locally in real-time            │     │
│  └────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────┘
```

---

## Layer Separation

| Layer | Tool | Purpose |
|-------|------|---------|
| **Nervous System** | NATS | Real-time coordination, events, request/reply between all nodes |
| **Maintenance Plane** | SSH + Tailscale | Provisioning, updates, debugging, Ansible playbooks |
| **Intelligence** | Hermes / OpenClaw | Where AI reasoning happens — central or distributed by use case |
| **Monitoring** | Netdata | Health metrics on every node regardless of agent presence |

**Don't choose NATS or SSH — use both.** NATS is for real-time coordination. SSH is for ops.

---

## NATS Agent SDK Role

The NATS Agent SDK is the communication backbone regardless of agent placement:
- Nodes WITHOUT local agents still run a NATS leaf for task pub/sub and telemetry
- Nodes WITH local agents use NATS for inter-agent coordination and event-driven triggers
- Central (Legacy AI) can dispatch work to any node via NATS request/reply OR SSH — whichever fits

---

## Bottom Line

The question is always: **"Does this node need to *think* locally, or does it just need to *execute* commands from central?"**

- Second facility → local agent (autonomy + latency)
- Cloud VPS → no agent, SSH + NATS (cheaper, simpler)
- Robot → embedded agent (real-time reaction)
- Main command → Legacy AI on MBP, mother of all operations
