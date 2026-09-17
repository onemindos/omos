---
id: diagram-agent-hive-map
title: "Agent Hive Map"
type: diagram
created: 2026-02-27
modified: 2026-02-27
tags: [type/diagram, domain/so, topic/architecture, topic/agents]
---

# Agent Hive Map

> The 4-runtime architecture: how Zeus reaches Legacy, how Legacy commands workers, and how everything connects through NATS JetStream.

```mermaid
graph TB
    subgraph ZEUS["👤 ZEUS (Human)"]
        direction LR
        WA[WhatsApp]
        TG[Telegram]
        DC[Discord]
        VC[Voice]
        FL[Flutter App]
    end

    subgraph OC["⚡ OPENCLAW — Legacy's Runtime"]
        direction TB
        SOUL[SOUL.md\nPersonality & Tone]
        AGENTS_MD[AGENTS.md\nOperating Rules]
        HB[HEARTBEAT.md\nProactive Check / 30m]
        LEGACY_CORE[LEGACY AI\nClaude Sonnet]
        SKILLS[Skills\nnexus-bridge · agno-dispatch]
    end

    subgraph NATS["🔀 NATS JetStream — Nervous System"]
        direction LR
        N1[comms.text.human.zeus]
        N2[digital.agent.*.task]
        N3[digital.agent.*.result.legacy]
        N4[physical.robot.*.command]
        N5[nexus.system.health]
    end

    subgraph AGNO["🤖 AGNO AgentOS — Digital Workers"]
        direction TB
        R[Researcher]
        AN[Analyst]
        CO[Coder]
        WEB[Web Automator]
        DB[DB Manager]
        MM[Memory Manager]
    end

    subgraph OM1["🦾 OM1 — Physical Runtime"]
        ROB[Robots / Sensors]
        CAM[Cameras / Vision]
        MOT[Motors / Actuators]
        LR[LeRobot Skills]
    end

    WA & TG & DC & VC & FL --> LEGACY_CORE
    SOUL & AGENTS_MD & HB --> LEGACY_CORE
    LEGACY_CORE --> SKILLS
    SKILLS <--> N1
    SKILLS <--> N2
    N2 --> R & AN & CO & WEB & DB & MM
    R & AN & CO & WEB & DB & MM --> N3
    N3 --> LEGACY_CORE
    N4 --> ROB & CAM & MOT
    LR --> MOT
    N5 -.->|health ping| LEGACY_CORE

    style ZEUS fill:#1a1a2e,stroke:#7c6af7,color:#fff
    style OC fill:#16213e,stroke:#e94560,color:#fff
    style NATS fill:#0f3460,stroke:#00b4d8,color:#fff
    style AGNO fill:#1a1a2e,stroke:#4ade80,color:#fff
    style OM1 fill:#1a1a2e,stroke:#fb923c,color:#fff
```

## Key Relationships

| From | To | Via | Purpose |
|------|----|-----|---------|
| Zeus | Legacy | Any channel | Human → AI command |
| Legacy | Agno workers | NATS `digital.agent.*.task` | Delegate complex tasks |
| Workers | Legacy | NATS `digital.agent.*.result.legacy` | Return results |
| Legacy | OM1 | NATS `physical.robot.*.command` | Control physical devices |
| Legacy | Zeus | Same channel | Response / proactive nudge |
