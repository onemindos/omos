---
title: "Agent Architecture Diagram"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/03-agents]
---

# Agent Architecture

## The Four Agents

```mermaid
graph TD
    subgraph SHARED["🧠 Shared Intelligence"]
        ctx["System Context\nIdentity · Mission · Goals"]
        skills["Skills Library\nC · O · D · E · ×"]
        mem["Memory Tiers\nWork · Learning · Relationship"]
        traits["Trait System\nExpertise × Personality × Approach"]
    end

    subgraph SENTINEL["⚡ SENTINEL &lpar;SO Agent&rpar;"]
        st["systems-thinking · analytical · methodical"]
        sd["Domains 00-24: AI, Tech, Infra, Security"]
    end

    subgraph COACH["💪 COACH &lpar;HP Agent&rpar;"]
        ct["wellness · encouraging · consistent"]
        cd["Domains 25-49: Health, Skills, Finance"]
    end

    subgraph HAVEN["🏡 HAVEN &lpar;LE Agent&rpar;"]
        ht["family-focused · patient · protective"]
        hd["Domains 50-74: Family, Home, Heritage"]
    end

    subgraph EMPIRE["💰 EMPIRE &lpar;GE Agent&rpar;"]
        et["business · bold · strategic"]
        ed["Domains 75-99: Business, Revenue, Products"]
    end

    SHARED --> SENTINEL
    SHARED --> COACH
    SHARED --> HAVEN
    SHARED --> EMPIRE

    SENTINEL <-->|"cross-domain\ncoordination"| COACH
    SENTINEL <-->|"cross-domain\ncoordination"| HAVEN
    SENTINEL <-->|"cross-domain\ncoordination"| EMPIRE

    style SHARED fill:#1a1a3e,stroke:#888,color:#fff
    style SENTINEL fill:#1a2a4e,stroke:#4285f4,color:#fff
    style COACH fill:#1a3a1a,stroke:#34a853,color:#fff
    style HAVEN fill:#3a3a1a,stroke:#fbbc04,color:#fff
    style EMPIRE fill:#3a1a1a,stroke:#ea4335,color:#fff
```

## Cross-Domain Task Flow

```mermaid
sequenceDiagram
    participant Z as 👤 Zeus
    participant SO as ⚡ Sentinel
    participant GE as 💰 Empire
    participant HP as 💪 Coach

    Z->>SO: "Launch the product landing page"
    Note over SO: Detects cross-domain task

    SO->>SO: Infra setup (domain 01, 10)
    SO->>GE: Need marketing copy (domain 78)
    SO->>HP: Does this align with weekly goals? (domain 29)

    GE-->>SO: Copy ready ✅
    HP-->>SO: Confirmed — aligns with Q1 #1 ✅

    SO-->>Z: Landing page deployed 🚀

    Note over SO,GE: All agents share memory tiers
```

## Trait Composition

Each agent is composed from 3 trait dimensions:

| Agent | Expertise | Personality | Approach |
|-------|-----------|-------------|----------|
| Sentinel | systems-thinking, security | analytical, measured | methodical |
| Coach | wellness, research | encouraging, enthusiastic | consistent |
| Haven | family-focused, creative | patient, measured | protective |
| Empire | business, research | bold, analytical | strategic |
