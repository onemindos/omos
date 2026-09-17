# Data Flow

```mermaid
graph TD
    subgraph "📥 INPUT"
        V[Voice / Vapi]
        T[Text / Chat]
        O[Obsidian Edit]
        E[External API]
        EM[Email / Message]
    end
    
    subgraph "🔵 C — CAPTURE"
        INB[24 Inbox Queue]
    end
    
    subgraph "🟢 O — ORGANIZE"  
        CLS[Classify Type]
        TPL[Apply Template]
        FM[Add Frontmatter]
        RTE[Route to Domain]
    end
    
    subgraph "🟡 D — DIRECT"
        PRI[Check Priorities<br/>active-goals.md]
        ASN[Assign Agent<br/>SO/HP/LE/GE]
        QUE[Add to Queue]
    end
    
    subgraph "🔴 E — EXECUTE"
        CTX[Load Context<br/>system/ + memory/]
        WRK[Do the Work]
        DOC[Document Output]
        SAV[Save to Domain]
    end
    
    subgraph "⚡ × — MULTIPLY"
        REF[Daily Reflection]
        EXT[Extract Patterns]
        CON[Consolidate Knowledge]
        UPD[Update Memory]
    end
    
    subgraph "📚 KNOWLEDGE BASE"
        SO_K[00-24 SO]
        HP_K[25-49 HP]
        LE_K[50-74 LE]
        GE_K[75-99 GE]
    end
    
    subgraph "🧠 MEMORY"
        MW[Work]
        ML[Learning]
        MR[Relationship]
    end
    
    V & T & O & E & EM --> INB
    INB --> CLS --> TPL --> FM --> RTE
    RTE --> PRI --> ASN --> QUE
    QUE --> CTX --> WRK --> DOC --> SAV
    SAV --> SO_K & HP_K & LE_K & GE_K
    SAV --> REF --> EXT --> CON --> UPD
    UPD --> MW & ML & MR
    
    ML --> |"Insights feed back"| INB
    CTX --> |"loads"| MW & ML & MR
    
    style INB fill:#1a73e8,color:#fff
    style RTE fill:#34a853,color:#fff
    style QUE fill:#fbbc04,color:#333
    style SAV fill:#ea4335,color:#fff
    style UPD fill:#9c27b0,color:#fff
```

## The Loop

Information enters → gets processed through CODEX stages → lands in a domain → gets reviewed → generates new insights → feeds back in. **The system gets smarter with every cycle.**
