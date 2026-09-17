---
title: "System Architecture Diagram"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/00-framework]
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

    subgraph INTEL["🧠 INTELLIGENCE LAYER &lpar;_codex/&rpar;"]
        direction TB
        subgraph SYS["System Context"]
            direction LR
            s1["Identity"]
            s2["Mission"]
            s3["Goals"]
            s4["Beliefs"]
            s5["Preferences"]
        end
        subgraph SKL["Skills"]
            direction LR
            sk1["Capture"]
            sk2["Organize"]
            sk3["Direct"]
            sk4["Execute"]
            sk5["× Multiply"]
        end
        subgraph MEM["Memory"]
            direction LR
            m1["🔨 Work"]
            m2["📚 Learning"]
            m3["🤝 Relationship"]
        end
    end

    subgraph AGENTS["🤖 AGENT LAYER"]
        direction LR
        SO["⚡ Sentinel\nSO 00-24"]
        HP["💪 Coach\nHP 25-49"]
        LE["🏡 Haven\nLE 50-74"]
        GE["💰 Empire\nGE 75-99"]
    end

    subgraph KNOWLEDGE["📚 100 DOMAINS"]
        direction LR
        K1["00-24\nSovereign Ops"]
        K2["25-49\nHolistic Perf"]
        K3["50-74\nLegacy Evo"]
        K4["75-99\nGen Endeavor"]
    end

    subgraph OPS["⚙️ OPERATIONS &lpar;_system/&rpar;"]
        direction LR
        o1["Infrastructure"]
        o2["Security"]
        o3["Hooks"]
        o4["Deploy"]
    end

    subgraph PLATFORM["☁️ PLATFORMS"]
        direction LR
        p1["GitHub"]
        p2["VPS"]
        p3["Letta"]
        p4["Vapi"]
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
| Intelligence | Identity, skills, memory — the AI brain | `_codex/` |
| Agents | 4 domain agents, each covering 25 domains | Letta + OpenClaw |
| Knowledge | 100 domains of life documentation | Root folders |
| Operations | Scripts, configs, automation | `_system/` |
| Platforms | External services | Cloud/VPS |
