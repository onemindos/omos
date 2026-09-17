---
title: "Agent Architecture Diagram"
type: diagram
created: 2026-03-27
modified: 2026-04-30
tags: [type/diagram, domain/so, topic/agents, topic/openclaw]
---

# Agent Architecture

## The Nine Agents (OpenClaw)

```mermaid
graph TD
    subgraph SHARED["🧠 Shared Intelligence (NATS JetStream)"]
        ctx["zeus_profile block\nIdentity · Mission · Goals"]
        mem["PostgreSQL + pgvector\nMemory · Entities · Culture"]
        skills["OpenClaw Tools\nTodoist · GitHub · MCP · Home Assistant"]
    end

    subgraph LEGACY_NODE["👑 LEGACY (Master Orchestrator — SO)"]
        lg["Claude Sonnet 4.6 · Full L0-L3 access\nRoutes to all 8 domain agents via NATS"]
    end

    subgraph SO_AGENTS["⚡ Sovereign Operations (SO 00-24)"]
        GRID["Grid\nInfrastructure & DevOps\n00-24 systems"]
        GUARDIAN["Guardian\nSecurity & Compliance\n00-24 security"]
        TRINITY["Trinity\n9th Agent — Wildcard"]
    end

    subgraph HP_AGENTS["💪 Holistic Performance (HP 25-49)"]
        SPARTAN["Spartan\nPhysical performance · fitness\n25-49"]
        COACH["Coach\nMental performance · habits\n25-49"]
    end

    subgraph LE_AGENTS["🏡 Legacy Evolution (LE 50-74)"]
        EDEN["Eden\nHomestead · family environment\n50-74"]
        HAVEN["Haven\nRelationships · heritage\n50-74"]
    end

    subgraph GE_AGENTS["💰 Generational Empire (GE 75-99)"]
        ORACLE["Oracle\nFinancial intelligence · investments\n75-99"]
        EMPIRE["Empire\nBusiness operations · revenue\n75-99"]
    end

    SHARED --> LEGACY_NODE
    LEGACY_NODE -->|"NATS digital.agent.grid.task"| SO_AGENTS
    LEGACY_NODE -->|"NATS digital.agent.spartan.task"| HP_AGENTS
    LEGACY_NODE -->|"NATS digital.agent.eden.task"| LE_AGENTS
    LEGACY_NODE -->|"NATS digital.agent.oracle.task"| GE_AGENTS

    style SHARED fill:#1a1a3e,stroke:#888,color:#fff
    style LEGACY_NODE fill:#2a1a4e,stroke:#7c6af7,color:#fff
    style SO_AGENTS fill:#1a2a4e,stroke:#4285f4,color:#fff
    style HP_AGENTS fill:#1a3a1a,stroke:#34a853,color:#fff
    style LE_AGENTS fill:#3a3a1a,stroke:#fbbc04,color:#fff
    style GE_AGENTS fill:#3a1a1a,stroke:#ea4335,color:#fff
```

## Cross-Domain Task Flow

```mermaid
sequenceDiagram
    participant Z as 👤 Zeus
    participant LG as 👑 LEGACY
    participant GR as ⚡ Grid
    participant EM as 💰 Empire
    participant SP as 💪 Spartan

    Z->>LG: "Launch the product landing page"
    Note over LG: Detects cross-domain task

    LG->>GR: Infra setup (SO 01, 10) via NATS
    LG->>EM: Need marketing copy (GE 78) via NATS
    LG->>SP: Aligns with weekly goals? (HP 29) via NATS

    GR-->>LG: Infra ready ✅
    EM-->>LG: Copy ready ✅
    SP-->>LG: Confirmed — aligns with Q1 #1 ✅

    LG-->>Z: Landing page deployed 🚀

    Note over LG,EM: All agents share zeus_profile via NATS
```

## Autonomy Levels

| Level | Name | Who Acts | Examples |
|-------|------|----------|---------|
| L3 | Full autonomy | Agent alone | Read files, route, search, respond |
| L2 | Delegated | Agent + LEGACY orchestration | Research, code, web automation |
| L1 | HITL | Zeus approves | Financial moves, external comms |
| L0 | Human only | Zeus only | Legal, health emergencies, deletions |

## Agent Roster

| Agent | Domain | Primary Tools | Autonomy |
|-------|--------|--------------|---------|
| **LEGACY** | SO (master) | All tools | L0-L3 full |
| **Grid** | SO 00-24 | SSH, Docker, GitHub, infra APIs | L2-L3 |
| **Guardian** | SO 00-24 | Security scans, audit logs | L1-L2 |
| **Trinity** | SO 00-24 | Flexible | L2-L3 |
| **Spartan** | HP 25-49 | Health APIs, calendar, fitness | L2-L3 |
| **Coach** | HP 25-49 | Habit tracker, goals | L2-L3 |
| **Oracle** | GE 75-99 | Finance APIs, market data | L1-L2 |
| **Empire** | GE 75-99 | Business tools, CRM | L1-L2 |
| **Eden** | LE 50-74 | Home Assistant, FarmBot | L2-L3 |
| **Haven** | LE 50-74 | Family calendar, heritage | L2-L3 |


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
