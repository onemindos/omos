---
id: diagram-agent-hive-map
title: "Agent Hive Map"
type: diagram
created: 2026-02-27
modified: 2026-04-30
tags: [type/diagram, domain/so, topic/architecture, topic/agents, topic/openclaw]
---

# Agent Hive Map

> The OpenClaw runtime: how Zeus reaches LEGACY, how LEGACY commands domain agents, and how everything connects through NATS JetStream.

```mermaid
graph TB
    subgraph ZEUS["👤 ZEUS (Human)"]
        direction LR
        FL[Flutter App]
        VC[Voice / LiveKit]
        TG[Telegram]
        DC[Discord]
    end

    subgraph OC["👑 OPENCLAW — LEGACY's Runtime (omos-core)"]
        direction TB
        SOUL[SOUL.md\nPersonality & Tone]
        AGENTS_MD[AGENTS.md\nOperating Rules]
        HB[HEARTBEAT.md\nProactive Check / 30m]
        LEGACY_CORE[LEGACY\nClaude Sonnet 4.6]
        TOOLS[MCP Tools\nTodoist · GitHub · Home Assistant]
    end

    subgraph NATS["🔀 NATS JetStream — Nervous System"]
        direction LR
        N1[omos.comms.text.human.zeus]
        N2[omos.agent.*.task]
        N3[omos.agent.*.result]
        N4[omos.physical.robot.*.command]
        N5[omos.system.health]
    end

    subgraph DOMAIN_AGENTS["🤖 Domain Agents (OpenClaw)"]
        direction TB
        GR[Grid\nSO — Infra]
        GD[Guardian\nSO — Security]
        SP[Spartan\nHP — Performance]
        CO[Coach\nHP — Habits]
        OR[Oracle\nGE — Finance]
        EM[Empire\nGE — Business]
        ED[Eden\nLE — Homestead]
        HA[Haven\nLE — Family]
    end

    subgraph PHYSICAL["🦾 Physical Layer"]
        ROB[FarmBot / Robots]
        CAM[Frigate / Cameras]
        HA_SYS[Home Assistant]
    end

    FL & VC & TG & DC --> LEGACY_CORE
    SOUL & AGENTS_MD & HB --> LEGACY_CORE
    LEGACY_CORE --> TOOLS
    TOOLS <--> N1
    LEGACY_CORE <--> N2
    N2 --> GR & GD & SP & CO & OR & EM & ED & HA
    GR & GD & SP & CO & OR & EM & ED & HA --> N3
    N3 --> LEGACY_CORE
    N4 --> ROB & CAM & HA_SYS
    N5 -.->|health ping| LEGACY_CORE

    style ZEUS fill:#1a1a2e,stroke:#7c6af7,color:#fff
    style OC fill:#16213e,stroke:#e94560,color:#fff
    style NATS fill:#0f3460,stroke:#00b4d8,color:#fff
    style DOMAIN_AGENTS fill:#1a1a2e,stroke:#4ade80,color:#fff
    style PHYSICAL fill:#1a1a2e,stroke:#fb923c,color:#fff
```

## Key Relationships

| From | To | Via | Purpose |
|------|----|-----|---------|
| Zeus | LEGACY | Any channel | Human → AI command |
| LEGACY | Domain agents | NATS `omos.agent.*.task` | Delegate domain-specific tasks |
| Domain agents | LEGACY | NATS `omos.agent.*.result` | Return results |
| LEGACY | Physical | NATS `omos.physical.robot.*.command` | Control robots/home |
| LEGACY | Zeus | Same channel | Response / proactive nudge |

## Infrastructure

| Component | Host | Port |
|-----------|------|------|
| OpenClaw (LEGACY) | omos-core 100.123.126.34 | Gateway 11434 |
| OpenClaw (agents) | omos-sentinel 100.66.235.38 | Gateway 11434 |
| NATS JetStream | omos-core | 4222 |
| PostgreSQL memory | omos-core | 5432 |
| LiveKit voice | omos-core | 7880 |


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
