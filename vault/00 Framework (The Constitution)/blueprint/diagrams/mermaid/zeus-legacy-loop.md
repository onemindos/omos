---
id: diagram-zeus-legacy-loop
title: "Zeus + Legacy Loop"
type: diagram
created: 2026-02-27
modified: 2026-02-27
tags: [type/diagram, domain/so, topic/architecture, topic/agents, topic/legacy]
---

# Zeus + Legacy Loop

> The human↔AI feedback cycle: inbound requests, proactive outreach, and continuous learning across every interaction.

```mermaid
sequenceDiagram
    participant Z as 👤 Zeus
    participant CH as 📡 Channel<br/>(Voice/Flutter/Telegram)
    participant L as 👑 LEGACY<br/>(OpenClaw)
    participant MEM as 🧠 Memory<br/>(SOUL · USER · HEARTBEAT)
    participant N as 🔀 NATS
    participant W as 🤖 Domain Agents<br/>(OpenClaw)

    Note over Z,W: ── INBOUND LOOP ──

    Z->>CH: Message / Voice / Request
    CH->>L: Webhook / WebSocket
    L->>MEM: Load Zeus profile + context
    MEM-->>L: Identity · Preferences · State

    alt Simple request (L3 autonomy)
        L-->>Z: Direct response via channel
    else Needs a worker (L2)
        L->>N: digital.agent.{id}.task
        N->>W: Route to domain agent (Grid/Spartan/Oracle/etc.)
        W-->>N: digital.agent.{id}.result
        N-->>L: Result
        L-->>Z: Synthesized response
    else HITL required (L1)
        L->>Z: Approval request via Discord/Telegram
        Z-->>L: Approve / Deny
        L->>N: Execute if approved
    end

    Note over Z,W: ── PROACTIVE LOOP (every 30m) ──

    L->>MEM: Run HEARTBEAT checklist
    MEM-->>L: Pending tasks · Events · Alerts
    L->>Z: Proactive nudge if needed
    Z-->>L: Acknowledgement / instruction

    Note over Z,W: ── LEARNING LOOP ──

    L->>MEM: Write daily memory log
    L->>MEM: Update Zeus profile signals
    MEM-->>L: Refined context next session
```

## Autonomy Levels

| Level | Name | Who Acts | Examples |
|-------|------|----------|---------|
| L3 | Full autonomy | Legacy alone | Read files, route, search, respond |
| L2 | Delegated | Legacy + Workers | Research, code, web automation |
| L1 | HITL | Zeus approves | Financial moves, external comms, purchases |
| L0 | Human only | Zeus only | Legal, health emergencies, deletions |

## Loop Frequencies

- **Inbound Loop** — triggered on every message/request (real-time)
- **Proactive Loop** — every 30 minutes via HEARTBEAT checklist
- **Learning Loop** — end of each day, daily memory log written


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
