---
type: fleet-roster
agent: all
status: active
updated: 2026-05-11
---
# AGENTS.md — OneMind Fleet Roster

## Fleet Architecture (6 agents)

```
                    ┌─────────────────┐
                    │  ⚡ Legacy       │
                    │  (Orchestrator)  │
                    └────────┬────────┘
         ┌──────────┬────────┼────────┬──────────┐
         │          │        │        │          │
    ┌────┴───┐ ┌────┴──┐ ┌──┴───┐ ┌──┴──────┐ ┌┴────────┐
    │Zenith  │ │Unity  │ │Eternal│ │Sovereign│ │SENTINEL │
    │  HP    │ │  LE   │ │  GE  │ │   SO    │ │(indep)  │
    └────────┘ └───────┘ └──────┘ └─────────┘ └─────────┘
```

## Active Agents

| Agent | Emoji | Pillar | Model | Node |
|---|---|---|---|---|
| **Legacy** | ⚡ | ALL | Sonnet 4.6 | omos-core |
| **Zenith** | 🏔 | HP 25-49 | Sonnet 4.6 | omos-core |
| **Unity** | 🌿 | LE 50-74 | Sonnet 4.6 | omos-core |
| **Eternal** | 👑 | GE 75-99 | Sonnet 4.6 | omos-core |
| **Sovereignty** | 🌐 | SO 00-24 | Sonnet 4.6 | omos-core |
| **Sentinel** | 👁️ | Watchdog | Haiku 4.5 | omos-sentinel |

## Routing Rules

1. Legacy routes ALL domain requests to the appropriate super-agent
2. Super-agents load domain skills on demand — no sub-agent processes
3. Cross-domain requests go through Legacy only
4. Sentinel is independent — does NOT depend on Core agents
5. Sentinel does NOT handle domain requests — observation only
