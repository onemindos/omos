---
title: "CODEX Method Diagram"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/00-framework]
---

# C.O.D.E.× Method

## The Loop

```mermaid
graph LR
    C["🔵 C\nCAPTURE"]
    O["🟢 O\nORGANIZE"]
    D["🟡 D\nDIRECT"]
    E["🔴 E\nEXECUTE"]
    X["🟣 ×\nMULTIPLY"]

    C -->|"classify"| O
    O -->|"prioritize"| D
    D -->|"assign"| E
    E -->|"reflect"| X
    X -->|"insights\nfeed back"| C

    style C fill:#1a73e8,stroke:#fff,color:#fff
    style O fill:#34a853,stroke:#fff,color:#fff
    style D fill:#fbbc04,stroke:#333,color:#333
    style E fill:#ea4335,stroke:#fff,color:#fff
    style X fill:#9c27b0,stroke:#fff,color:#fff
```

## Each Stage in Detail

```mermaid
graph TD
    subgraph C["🔵 CAPTURE"]
        c1["Voice note / text / API"]
        c2["Quick capture to Inbox"]
        c3["Agent observation"]
    end

    subgraph O["🟢 ORGANIZE"]
        o1["Classify type"]
        o2["Apply template"]
        o3["Route to domain"]
    end

    subgraph D["🟡 DIRECT"]
        d1["Check active goals"]
        d2["Assign to agent"]
        d3["Set priority"]
    end

    subgraph E["🔴 EXECUTE"]
        e1["Load context"]
        e2["Do the work"]
        e3["Save to domain"]
    end

    subgraph X["🟣 × MULTIPLY"]
        x1["Extract patterns"]
        x2["Update memory"]
        x3["Consolidate wisdom"]
    end

    C --> O --> D --> E --> X
    X -.->|"compound\nloop"| C

    style C fill:#1a3a6e,stroke:#4285f4,color:#fff
    style O fill:#1a3a1a,stroke:#34a853,color:#fff
    style D fill:#3a3a1a,stroke:#fbbc04,color:#fff
    style E fill:#3a1a1a,stroke:#ea4335,color:#fff
    style X fill:#2a1a3a,stroke:#9c27b0,color:#fff
```

## Why × Multiply Matters

Most frameworks stop at **Organize** (Second Brain) or **Execute** (GTD).

OneMind's **× Multiply** is where AI extracts patterns, consolidates knowledge, and feeds insights back — making the system self-improving. Every cycle, you get smarter.
