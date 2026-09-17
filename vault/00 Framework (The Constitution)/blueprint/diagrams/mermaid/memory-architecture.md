---
title: "Memory Architecture Diagram"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/00-framework]
---

# Memory Architecture

## 3-Tier System

```mermaid
graph TD
    INPUT["📥 Inputs"]
    ROUTER{"Memory\nRouter"}

    subgraph WORK["🔨 WORK TIER"]
        w1["Active tasks"]
        w2["Session state"]
        w3["Current projects"]
    end

    subgraph LEARN["📚 LEARNING TIER"]
        l1["signals.jsonl\n&lpar;append-only&rpar;"]
        l2["patterns.md\n&lpar;consolidated&rpar;"]
        l3["errors.md"]
    end

    subgraph REL["🤝 RELATIONSHIP TIER"]
        r1["How Zeus works"]
        r2["Trust levels"]
        r3["Preferences learned"]
    end

    INPUT --> ROUTER
    ROUTER -->|"task state"| WORK
    ROUTER -->|"feedback"| LEARN
    ROUTER -->|"about Zeus"| REL

    WORK -->|"on completion"| LEARN
    LEARN -->|"weekly review"| l2
    REL -->|"pattern detection"| r2

    style WORK fill:#1a2a4e,stroke:#4285f4,color:#fff
    style LEARN fill:#1a3a1a,stroke:#34a853,color:#fff
    style REL fill:#3a3a1a,stroke:#fbbc04,color:#fff
```

## Session Lifecycle

```mermaid
sequenceDiagram
    participant A as 🤖 Agent
    participant W as 🔨 Work
    participant L as 📚 Learning
    participant R as 🤝 Relationship

    Note over A: ── Session Start ──
    A->>W: Load current-focus.json
    A->>L: Scan recent signals
    A->>R: Load interaction notes

    Note over A: ── During Work ──
    A->>W: Update active task

    Note over A: ── Zeus gives feedback ──
    A->>L: Log signal to signals.jsonl

    Note over A: ── Session End ──
    A->>W: Archive completed work
    A->>L: Extract lessons learned
    A->>R: Update preferences
```

## What Each Tier Stores

| Tier | Purpose | Lifecycle | Format |
|------|---------|-----------|--------|
| **Work** | What's happening now | Created → archived | JSON |
| **Learning** | What we've learned | Append-only → consolidated weekly | JSONL + MD |
| **Relationship** | How to work with Zeus | Updated after significant interactions | Markdown |

## Location

```
_codex/memory/
├── work/
│   └── current-focus.json
├── learning/
│   ├── signals.jsonl
│   └── patterns.md
└── relationship/
    └── interaction-notes.md
```


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
