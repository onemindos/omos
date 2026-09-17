---
title: "System Architecture Diagram"
type: diagram
created: 2026-03-27
modified: 2026-04-30
tags: [type/diagram, domain/so, topic/architecture, topic/openclaw]
---

# OneMind System Architecture

```mermaid
graph TB
    subgraph HUMAN["👤 HUMAN LAYER"]
        direction LR
        Z["Zeus (Owner)"]
        OBS["Obsidian Vault"]
        VSC["VS Code"]
        Z --- OBS
        Z --- VSC
    end

    subgraph INTEL["🧠 INTELLIGENCE LAYER (OpenClaw / LEGACY)"]
        direction TB
        subgraph SYS["System Context"]
            direction LR
            s1["Identity"]
            s2["Mission"]
            s3["Goals"]
            s4["Beliefs"]
            s5["Preferences"]
        end
        subgraph SKL["Skills — C·O·D·E·×"]
            direction LR
            sk1["Capture"]
            sk2["Organize"]
            sk3["Direct"]
            sk4["Execute"]
            sk5["× Multiply"]
        end
        subgraph MEM["Memory (PostgreSQL + pgvector)"]
            direction LR
            m1["🔨 Work"]
            m2["📚 Learning"]
            m3["🤝 Relationship"]
        end
    end

    subgraph AGENTS["🤖 AGENT LAYER (OpenClaw)"]
        direction LR
        MASTER["👑 LEGACY\nMaster Orchestrator"]
        SO["⚡ Grid · Guardian\nSO 00-24"]
        HP["💪 Spartan · Coach\nHP 25-49"]
        LE["🏡 Eden · Haven\nLE 50-74"]
        GE["💰 Oracle · Empire\nGE 75-99"]
        MASTER --> SO & HP & LE & GE
    end

    subgraph KNOWLEDGE["📚 100 DOMAINS"]
        direction LR
        K1["00-24\nSovereign Ops"]
        K2["25-49\nHolistic Perf"]
        K3["50-74\nLegacy Evo"]
        K4["75-99\nGen Endeavor"]
    end

    subgraph OPS["⚙️ OPERATIONS"]
        direction LR
        o1["Infrastructure"]
        o2["Security"]
        o3["N8N Workflows"]
        o4["NATS JetStream"]
    end

    subgraph PLATFORM["☁️ PLATFORMS"]
        direction LR
        p1["GitHub"]
        p2["VPS (omos-core · omos-sentinel)"]
        p3["OpenClaw"]
        p4["LiveKit"]
        p5["AWS Bedrock"]
    end

    HUMAN --> INTEL
    INTEL --> AGENTS
    AGENTS --> KNOWLEDGE
    AGENTS --> OPS
    OPS --> PLATFORM

    style HUMAN fill:#2d2d3f,stroke:#888,color:#fff
    style INTEL fill:#1a1a4e,stroke:#4285f4,color:#fff
    style AGENTS fill:#1a3a1a,stroke:#34a853,color:#fff
    style KNOWLEDGE fill:#3a2a1a,stroke:#fbbc04,color:#fff
    style OPS fill:#3a1a1a,stroke:#ea4335,color:#fff
    style PLATFORM fill:#2a1a3a,stroke:#9c27b0,color:#fff
```

## Reading Guide

**Top to Bottom = Human to Machine**

| Layer | What | Where |
|-------|------|-------|
| Human | Zeus interacts via Obsidian + VS Code | Physical |
| Intelligence | Identity, skills, memory — the AI brain | OpenClaw / LEGACY |
| Agents | LEGACY master + 8 domain agents (SO/HP/LE/GE) | omos-core + omos-sentinel |
| Knowledge | 100 domains of life documentation | Codex root folders |
| Operations | NATS events, N8N automation, config | omos-core |
| Platforms | AWS Bedrock, LiveKit, GitHub, VPS | Cloud + VPS |


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
