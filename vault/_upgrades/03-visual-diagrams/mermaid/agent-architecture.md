# Agent Architecture

```mermaid
graph TB
    subgraph "🧠 Shared Intelligence (_codex/)"
        SYS[System Context<br/>Identity, Mission, Goals, Beliefs]
        SKILLS[Skills Library<br/>Capture, Organize, Direct, Execute, Multiply]
        MEM[Memory<br/>Work / Learning / Relationship]
        TRAITS[Trait System<br/>Expertise × Personality × Approach]
    end
    
    subgraph "⚡ SO Agent — Sentinel"
        SO_T[Traits: systems-thinking,<br/>analytical, methodical]
        SO_D[Domains: 00-24<br/>AI, Tech, Infra, Security]
        SO_R[Role: Central orchestrator<br/>Routes cross-domain tasks]
    end
    
    subgraph "💪 HP Agent — Coach"  
        HP_T[Traits: wellness,<br/>encouraging, consistent]
        HP_D[Domains: 25-49<br/>Health, Skills, Finance]
        HP_R[Role: Personal optimizer<br/>Habit & health tracking]
    end
    
    subgraph "🏡 LE Agent — Haven"
        LE_T[Traits: family-focused,<br/>patient, protective]
        LE_D[Domains: 50-74<br/>Family, Home, Heritage]
        LE_R[Role: Family steward<br/>Legacy preservation]
    end
    
    subgraph "💰 GE Agent — Empire"
        GE_T[Traits: business,<br/>bold, strategic]
        GE_D[Domains: 75-99<br/>Business, Revenue, Products]
        GE_R[Role: Wealth creator<br/>Revenue & growth]
    end
    
    SYS --> SO_T & HP_T & LE_T & GE_T
    SKILLS --> SO_R & HP_R & LE_R & GE_R
    MEM --> SO_D & HP_D & LE_D & GE_D
    TRAITS --> SO_T & HP_T & LE_T & GE_T
    
    SO_R <--> |"Cross-domain<br/>coordination"| HP_R
    SO_R <--> |"Cross-domain<br/>coordination"| LE_R
    SO_R <--> |"Cross-domain<br/>coordination"| GE_R
```

```mermaid
sequenceDiagram
    participant Z as Zeus
    participant SO as SO (Sentinel)
    participant HP as HP (Coach)
    participant GE as GE (Empire)
    
    Z->>SO: "Set up the new product landing page"
    Note over SO: Cross-domain task detected
    SO->>SO: Technical infrastructure (domain 01, 10)
    SO->>GE: Marketing copy needed (domain 78)
    SO->>HP: Check if this aligns with weekly goals (domain 29)
    
    GE-->>SO: Copy ready, brand-aligned
    HP-->>SO: Confirmed: aligns with Q1 priority #1
    SO-->>Z: Landing page deployed ✅
    
    Note over SO,GE: All agents share memory tiers
```

## The Pattern

- **Each agent owns 25 domains** — clear boundaries, no overlap
- **SO is the router** — cross-domain tasks go through Sentinel
- **All agents share** — identity context, skills, memory, and traits
- **Traits differentiate** — same skills, different personality and approach
