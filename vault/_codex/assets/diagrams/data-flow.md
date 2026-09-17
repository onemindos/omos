---
title: "Data Flow Diagram"
type: diagram
created: 2026-03-27
tags: [type/diagram, domain/00-framework]
---

# Data Flow — How Information Moves

```mermaid
graph TD
    subgraph IN["📥 INPUTS"]
        i1["🎤 Voice"]
        i2["💬 Chat"]
        i3["📝 Obsidian"]
        i4["🔌 API"]
        i5["📧 Email"]
    end

    subgraph CAPTURE["🔵 CAPTURE"]
        inbox["24 Inbox Queue"]
    end

    subgraph ORGANIZE["🟢 ORGANIZE"]
        classify["Classify type"]
        template["Apply template"]
        route["Route to domain"]
    end

    subgraph DIRECT["🟡 DIRECT"]
        goals["Check goals"]
        assign["Assign agent"]
        queue["Priority queue"]
    end

    subgraph EXECUTE["🔴 EXECUTE"]
        context["Load context"]
        work["Do the work"]
        save["Save output"]
    end

    subgraph MULTIPLY["🟣 × MULTIPLY"]
        reflect["Reflect"]
        extract["Extract patterns"]
        consolidate["Update memory"]
    end

    subgraph VAULT["📚 KNOWLEDGE"]
        v1["00-24 SO"]
        v2["25-49 HP"]
        v3["50-74 LE"]
        v4["75-99 GE"]
    end

    subgraph MEMORY["🧠 MEMORY"]
        mem1["Work"]
        mem2["Learning"]
        mem3["Relationship"]
    end

    IN --> inbox
    inbox --> classify --> template --> route
    route --> goals --> assign --> queue
    queue --> context --> work --> save
    save --> v1 & v2 & v3 & v4
    save --> reflect --> extract --> consolidate
    consolidate --> mem1 & mem2 & mem3
    mem2 -.->|"insights\nloop back"| inbox

    style CAPTURE fill:#1a3a6e,stroke:#4285f4,color:#fff
    style ORGANIZE fill:#1a3a1a,stroke:#34a853,color:#fff
    style DIRECT fill:#3a3a1a,stroke:#fbbc04,color:#fff
    style EXECUTE fill:#3a1a1a,stroke:#ea4335,color:#fff
    style MULTIPLY fill:#2a1a3a,stroke:#9c27b0,color:#fff
```

## The Key Insight

Information flows **one direction** through CODEX stages, but the × Multiply step creates a **feedback loop** — insights from reflection feed back into Capture. This is what makes OneMind self-improving.
