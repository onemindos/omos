---
id: offline-operation-awareness-hub-and-system-integration
title: "100% Offline Operation, Awareness Hub, and Full System Integration Session"
type: note
source: apple-notes
imported: 2026-03-07
status: staged
tags:
  - apple-notes-import
  - offline-operation
  - mesh-network
  - awareness-hub
  - awareness-modes
  - permissions
  - skills-system
  - redis
  - nats-kv
  - flutter
  - agno
  - novu
  - implementation-plan
  - deployment
  - websocket
  - real-time
  - docker
  - git
---

# 100% Offline Operation, Awareness Hub, and Full System Integration Session

---

## 100% Offline Operation — Fleet Continues Working

- ✅ **100% offline operation** — fleet continues working
- ✅ **5-10km range** — entire property covered
- ✅ **Automatic failover** — transparent to applications
- ✅ **Dual-path alerts** — critical alerts via internet + mesh
- ✅ **Emergency resilience** — works when everything else fails
- ✅ **Zero monthly costs** — no cellular plans
- ✅ **2-year battery life** — with solar: indefinite

---

## Awareness Hub — Naming Options

| Name | Rationale |
|------|-----------|
| **Awareness Hub** | Short, clear, matches existing `/awareness` endpoints |
| **Universal Awareness Hub** | Emphasizes it controls everything |
| **Legacy Control Center** | Focuses on controlling what Legacy can do |
| **Consciousness Control** | Matches existing "consciousness" terminology in code |
| **Mode Control** | Simple, describes the 4-mode system |

---

## Awareness Mode System

The feature includes:

- **4 Modes:** Dormant → Aware → Present → Omnipresent
- **12 Permissions:** What Legacy can/cannot do per mode
- **Context Sources:** What data Legacy has access to
- **Customization:** User-tweakable settings

---

## Implementation Plan — 10 Phases

| # | Phase | Status | Days |
|---|-------|--------|------|
| 1 | Backend Permissions Model | ✅ DONE | 2-3 |
| 2 | NATS KV Integration | ⏳ Pending | 1-2 |
| 3 | WebSocket Real-time Sync | ⏳ Pending | 1-2 |
| 4 | Novu Cloud Integration | ⏳ Pending | 2 |
| 5 | Agent Behavior Integration | ⏳ Pending | 2 |
| 6 | Flutter Provider & Models | ⏳ Pending | 2 |
| 7 | Flutter UI Components | ⏳ Pending | 2-3 |
| 8 | Flutter Inbox & Activity | ⏳ Pending | 2 |
| 9 | Extensibility & Plugins | ⏳ Pending | 2-3 |
| 10 | Testing & Documentation | ⏳ Pending | 1-2 |
| | **Total** | | **17-23 days** |

---

## Skills Loading System

Skills = pre-loaded tool configurations + instructions:

```python
# Agent skill registry
SKILLS = {
    "task_management": {
        "tools": ["todoist_complete", "todoist_create", "linear_create"],
        "instructions": "When completing tasks, check pillar impact...",
        "requires_auth": ["nango:todoist", "nango:linear"]
    },
    "home_control": {
        "tools": ["ha_light", "ha_climate", "ha_scene"],
        "instructions": "Check presence before automation...",
        "requires_auth": []  # No OAuth, direct MQTT
    }
}

# When agent gets a job, inject relevant skill
agent.inject_skill(SKILLS["task_management"])
```

---

## Redis vs NATS KV — Use Case Comparison

| Use Case | Redis | NATS KV |
|----------|-------|---------|
| Agent session state | ✅ Keep | ❌ |
| LiveKit tokens | ✅ Keep | ❌ |
| Awareness state | ❌ | ✅ Better (watch feature) |
| Config | ❌ | ✅ Better |

---

## Flutter Web UI ↔ Agno Backend Architecture

### Summary Diagram

```
 ┌─────────────────────┐     HTTP      ┌─────────────────────┐
 │  Flutter Web UI     │ ──────────────▶│  Agno Backend       │
 │  (port 3005)        │               │  (port 8081)        │
 │                     │               │                     │
 │  AgnoClient.dart    │  /api/agents  │  main.py            │
 │  ───────────────    │ ─────────────▶│  list_agents()      │
 │  getAgents()        │               │  returns 18 agents  │
 │  getTeams()         │  /api/teams   │  list_teams()       │
 │  getModels()        │ ─────────────▶│  returns 4 teams    │
 └─────────────────────┘               └─────────────────────┘
```


---

## Building Agno Features → Flutter UI

The pattern for any new feature:

| Layer | What to Create | Example |
|-------|---------------|---------|
| 1. Backend | Add endpoint in `main.py` | `GET /api/workflows` |
| 2. Model | Add Dart model in `core/models/` | `workflow.dart` + freezed |
| 3. API Client | Add method in `agno_client.dart` | `getWorkflows()` |
| 4. Provider | Add Riverpod provider in `app_providers.dart` | `workflowsProvider` |
| 5. Screen | Create screen that watches provider | `WorkflowsScreen` |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
