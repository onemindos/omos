---
title: "OneMind vs Second Brain"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/00-framework]
---

# OneMind vs Second Brain

```mermaid
graph TD
    subgraph SB["📓 SECOND BRAIN"]
        sb1["✅ Capture"]
        sb2["✅ Organize"]
        sb3["❌ Direct"]
        sb4["❌ Execute"]
        sb5["❌ Multiply"]
        sb1 --> sb2
        sb2 -.->|"stops here"| sb3
    end

    subgraph OM["🧠 ONE MIND"]
        om1["✅ Capture"]
        om2["✅ Organize"]
        om3["✅ Direct"]
        om4["✅ Execute"]
        om5["✅ × Multiply"]
        om1 --> om2 --> om3 --> om4 --> om5
        om5 -.->|"loops back"| om1
    end

    style SB fill:#2a2a2a,stroke:#666,color:#aaa
    style OM fill:#1a1a3e,stroke:#9c27b0,color:#fff
```

## Feature Comparison

```mermaid
graph LR
    subgraph SECOND["📓 Second Brain"]
        direction TB
        f1["Note-taking"]
        f2["Folder organization"]
        f3["Tagging"]
        f4["Search"]
        f5["Manual review"]
    end

    subgraph ONEMIND["🧠 One Mind"]
        direction TB
        g1["Note-taking + Voice + API capture"]
        g2["100-domain auto-routing"]
        g3["Structured frontmatter + tags"]
        g4["AI-powered tiered search"]
        g5["AI-driven × Multiply reflection"]
        g6["4 AI domain agents"]
        g7["3-tier persistent memory"]
        g8["Lifecycle hooks"]
        g9["Security validation"]
        g10["Generational knowledge transfer"]
    end

    style SECOND fill:#2a2a2a,stroke:#666,color:#aaa
    style ONEMIND fill:#1a1a3e,stroke:#9c27b0,color:#fff
```

## The Bottom Line

| | Second Brain | One Mind |
|---|---|---|
| **What it is** | Note-taking system | Human enhancement framework |
| **AI role** | Optional add-on | Native partner |
| **Scope** | Your notes | Your entire life (100 domains) |
| **Action** | Capture + organize | Capture → organize → direct → execute → multiply |
| **Memory** | Your notebook | 3-tier AI memory that learns |
| **Agents** | None | 4 specialized domain agents |
| **Time horizon** | Your lifetime | Multi-generational |
| **Self-improving** | No | Yes — × Multiply feedback loop |

> *"They gave you a Second Brain. I'm giving you One Mind."*
