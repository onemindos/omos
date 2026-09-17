# Proactive Engine Implementation — Full Changelog

**Date**: February 20, 2026  
**Scope**: Transform OneMind OS from reactive-only to proactive + reactive  
**Architecture Decision**: Hybrid Agno-First — Legacy becomes a full Agno agent, Nexus slims to infrastructure only

---

## The Problem

OneMind OS had a fundamental gap: **Agno agents are reactive**. They only respond when a user sends a message. But the vision for Legacy (the user's AI command center) requires **proactive behavior** — monitoring systems 24/7, catching problems before the user notices, running periodic checks, and acting autonomously.

The existing Nexus layer tried to solve this by wrapping an Agno agent inside a custom `Legacy` class with its own dispatch system. This created duplication:
- Nexus had its own dispatcher (~400 lines) that replicated what Agno already does
- Nexus had its own Legacy class (~580 lines) that was "an agent with extra steps"
- The system had two command paths (Agno's `/agents/{id}/runs` AND Nexus's `/nexus/command`)

**Research inputs**:
- **OpenClaw** (open-source AI gateway, 214k stars) — heartbeat loop pattern, HEARTBEAT.md checklist, HEARTBEAT_OK suppression, cron vs heartbeat separation
- **ada_local** (local AI assistant) — voice pipeline architecture, function routing, proactive monitoring

---

## What Was Done (In Order)

### 1. Architecture Design Document

**File**: `docs/PROACTIVE_ARCHITECTURE.md` (NEW — ~420 lines)

**Why**: Before writing code, the entire design needed to be documented. This captures the research findings from OpenClaw and ada_local, the architectural decisions, the implementation phases, and the sidebar consolidation spec. It serves as the north star for all subsequent changes.

**What it covers**:
- The reactive gap problem analysis
- Research summary from OpenClaw (heartbeat, cron, nodes, skills) and ada_local (voice, function routing)
- What changes vs what stays
- New architecture diagram (Legacy as Agno agent → Heartbeat → NATS → Nexus slim)
- 7-phase implementation plan
- Sidebar consolidation spec (COMMAND CENTER)
- Deployment topology (local, cloud, mobile)
- Key principles

---

### 2. Heartbeat Loop

**File**: `backend/heartbeat.py` (NEW — ~260 lines)

**Why**: This is the core solution to the proactive gap. Instead of waiting for user input, the heartbeat loop wakes Legacy up on a configurable interval (default: 30 minutes) and asks it to check on things. This is directly inspired by OpenClaw's heartbeat pattern — a periodic "are things OK?" loop that only escalates when something needs attention.

**What it does**:
- `HeartbeatConfig` — reads configuration from environment variables:
  - `HEARTBEAT_ENABLED` (default: true)
  - `HEARTBEAT_INTERVAL_MINUTES` (default: 30)
  - `HEARTBEAT_ACTIVE_HOURS_START` / `END` (default: 6–23, so it doesn't wake you at 3am)
- `HeartbeatLoop` — the main engine:
  - Takes `agents_dict` (to find Legacy) and optional `on_alert` callback
  - `start()` / `stop()` — lifecycle management
  - `_loop()` — runs forever, sleeping between ticks
  - `_tick()` — the core cycle:
    1. Check if within active hours (skip if not)
    2. Read `backend/workspace/HEARTBEAT.md` (the checklist)
    3. Build a prompt with the checklist + current time
    4. Run an Agno agent turn (Legacy thinks through the checklist)
    5. Process the response
  - `_process_response()` — implements the **HEARTBEAT_OK protocol**:
    - If Legacy's response contains "HEARTBEAT_OK", everything is fine → suppress (no noise)
    - If NOT, something needs attention → publish alert to NATS topic `system.heartbeat.alert`
- Singleton pattern via `create_heartbeat()` / `get_heartbeat()`

**Design choice — why HEARTBEAT_OK suppression?**: Inspired by OpenClaw. If the system checks 48 times a day and everything is fine 47 times, you don't want 47 notifications. The agent only speaks up when something is wrong. This is the difference between useful proactivity and spam.

---

### 3. Heartbeat Checklist

**File**: `backend/workspace/HEARTBEAT.md` (NEW)

**Why**: The heartbeat loop needs to know *what* to check. Rather than hardcoding a list, we use a Markdown file that Legacy reads each cycle. This is editable — the user (or Legacy itself) can modify the checklist without touching code.

**Current checklist items**:
1. Check NATS for interesting events since last check
2. Scan for failed agent/team runs
3. Review any pending tasks or deadlines
4. Light check-in — is the user likely active?
5. Summarize any running workflows
6. Quick health check on connected systems

---

### 4. NATS Tools for Legacy

**File**: `backend/tools/nats_tools.py` (NEW — ~230 lines)

**Why**: For Legacy to be truly aware of the system, it needs to observe the NATS event bus (where ALL system events flow) and publish events when it takes action. These tools give Legacy "eyes and hands" on the bus, turning it from a chat agent into a system-aware orchestrator.

**Three tools created**:

1. **`nats_observe(topic, limit)`** — Subscribe to a NATS topic for 2 seconds and collect messages
   - *Why*: Legacy needs to see what's happening on the bus (agent activity, system events, heartbeats)
   - Returns: list of messages with subject, data, timestamp

2. **`nats_publish(topic, message, event_type)`** — Publish an event to the NATS bus
   - *Why*: When Legacy decides to take action (assign a task, alert the user, coordinate agents), it needs to announce it on the bus so other systems can react
   - Publishes: structured event with source="legacy", timestamp, event_type

3. **`system_overview()`** — Get a complete snapshot of the system
   - *Why*: Legacy needs situational awareness — how many agents, which teams, NATS status, heartbeat status, available tools
   - Returns: agents list, teams list, NATS connection status, heartbeat status, tool availability

---

### 5. Tool Registry Updates

**File**: `backend/tools/registry.py` (MODIFIED)

**Why**: Tools in OneMind OS are registered centrally so agent presets can reference them by name. The 3 new NATS tools needed to be registered here so Legacy (and any future agent) can use them.

**Changes**:
- Added 3 tool specs to `BUILTIN_TOOL_SPECS`:
  - `nats_observe` — category: BUILTIN, description, parameter schema
  - `nats_publish` — category: BUILTIN, description, parameter schema
  - `system_overview` — category: BUILTIN, description, no parameters
- Added `get_tool_availability()` function — returns a dict of all registered tool categories and their tool names. Used by `system_overview()` tool so Legacy can see what tools are available.

---

### 6. Legacy Agent Preset

**File**: `backend/agents/presets.py` (MODIFIED)

**Why**: This is the architectural pivot. Previously, Legacy was a custom Python class in `nexus/legacy.py` that wrapped an Agno agent and added its own dispatch/observation layer. Now Legacy is a **first-class Agno agent** defined as a preset — just like every other agent in the system. This means:
- Legacy gets ALL of Agno's features for free (memory, sessions, streaming, tool use, reasoning)
- Legacy is accessible via the standard `/agents/legacy/runs` endpoint
- No custom code path to maintain
- The heartbeat loop can use Legacy just like any other Agno agent

**Preset configuration**:
- **Position**: First entry in `AGENT_PRESETS` dict (it's the primary agent)
- **Model**: `gpt-4o` (best reasoning for a super-orchestrator)
- **Tools**: `duckduckgo`, `clickup`, `nats_observe`, `nats_publish`, `system_overview`
- **Memory**: All enabled — `enable_user_memories`, `enable_session_summaries`, `enable_agentic_context`
- **Reasoning**: Enabled with max 15 steps
- **History**: `num_history_runs: 20` (deep context window)
- **Instructions**: Extensive system prompt covering:
  - Identity (Legacy, Zeus's AI command center)
  - Proactive behavior (heartbeat protocol, HEARTBEAT_OK)
  - System awareness (NATS observation, agent coordination)
  - Capabilities (search, task management, system monitoring)
  - Personality (direct, mission-focused, no fluff)

---

### 7. Main.py Wiring

**File**: `backend/main.py` (MODIFIED — 3 changes)

**Why**: The heartbeat loop needs to start when the server boots and stop when it shuts down. This wires it into the existing startup/shutdown lifecycle.

**Change 1 — Startup (Step 6a)**:
After the workflow scheduler starts, we:
1. Import `create_heartbeat` from `backend.heartbeat`
2. Create the heartbeat singleton with `agents_dict`
3. Call `heartbeat.start()` to begin the loop
4. Wrapped in try/except so heartbeat failure doesn't crash the server

**Change 2 — Shutdown**:
In the `shutdown_event()` handler:
1. Import `get_heartbeat`
2. If heartbeat exists, call `hb.stop()` to cleanly terminate the loop

**Change 3 — Boot Banner**:
Updated the ASCII art boot banner to show:
- "Legacy: Proactive + Reactive (Jarvis Mode)"
- "Heartbeat: Always Watching"
This is cosmetic but signals that the proactive engine is active.

---

### 8. Sidebar Consolidation

**File**: `frontend/lib/widgets/sidebar.dart` (MODIFIED)

**Why**: The user explicitly asked to combine COMMAND + SYSTEM + NEXUS + GAME into a single section. The old sidebar had 7 sections — too many clicks, too much scrolling, important items buried under collapsed sections. The new layout puts everything command-relevant in one always-visible block.

**Before** (7 sections):
```
COMMAND (always expanded) — Nexus, Inbox, Chat, Tasks, Calendar, Activity
WORKFORCE — Agents, Teams, Machines, Assets
WORKSPACE — Projects, Documents, Knowledge, Memories
AUTOMATION — Workflows, Tools & MCP, Sessions
SYSTEM (collapsed) — Mission Control, Events, Analytics
GAME (collapsed) — Commander, Quests, Skill Tree, HQ, Tactical Base, Roster, Game Hub
SETTINGS (collapsed) — Settings, Integrations
```

**After** (5 sections):
```
COMMAND CENTER (always expanded)
  ── Core Ops ──
  Nexus, Inbox, Chat, Tasks, Calendar, Activity
  ── System ──
  Mission Control, Events, Analytics, Heartbeat [LIVE]
  ── Commander ──
  Commander, Game Hub

WORKFORCE — Agents, Teams, Machines, Assets
WORKSPACE — Projects, Documents, Knowledge, Memories
AUTOMATION — Workflows, Tools & MCP, Sessions
SETTINGS — Settings, Integrations
```

**Key decisions**:
- System monitoring (Mission Control, Events, Analytics) moved up to COMMAND CENTER because monitoring is a command activity, not a buried afterthought
- **Heartbeat** added as a new nav item with a green LIVE badge — gives visibility into the proactive engine
- Game items consolidated to just Commander (profile) and Game Hub (the rest are accessible from Game Hub)
- Collapsed sections state updated: only `workspace`, `automation`, `settings` start collapsed (removed `system` and `game` since they're gone)

---

### 9. Nexus Core Slimming

**File**: `backend/nexus/core.py` (REWRITTEN — 276 → ~130 lines)

**Why**: With Legacy as an Agno agent and dispatch handled by Agno natively, Nexus no longer needs to instantiate a `Legacy` class or a `Dispatcher`. The core's job is reduced to:
1. **Framework Registry** — keep the unified entity view across Agno + future frameworks
2. **Agent Awareness** — keep NATS awareness (subscribe agents to topics)

**Removed**:
- `from backend.nexus.dispatcher import Dispatcher` — Agno handles dispatch
- `from backend.nexus.legacy import Legacy` — Legacy is an Agno preset
- `self.dispatcher = Dispatcher(self.registry)` — gone
- `self.legacy = Legacy(self)` — gone
- `command()`, `dispatch()`, `fan_out()`, `broadcast()` — all convenience methods removed
- `assign_task()` — removed (was coupling Nexus dispatch with task engine)
- `overview()` simplified — no more legacy/dispatcher stats

**Kept**:
- `register_agno()`, `register_openmind()`, `register_custom()` — framework registration
- `start()` / `shutdown()` — lifecycle (awareness activation, NATS events)
- `overview()` — simplified to frameworks + awareness + agent/team counts
- Singleton pattern (`get_nexus()` / `reset_nexus()`)

---

### 10. Archive Legacy & Dispatcher Modules

**Files**: `backend/nexus/legacy.py` → `legacy.py.archived`, `backend/nexus/dispatcher.py` → `dispatcher.py.archived`

**Why**: These files are no longer imported by anything. Rather than deleting them (losing the code), we archive them so they can be referenced if any specific logic needs to be recovered later. Combined they were ~970 lines of code that duplicated Agno functionality.

**What they contained**:
- `legacy.py` (577 lines): SystemEvent dataclass, EventLog (in-memory deque), Legacy class with NATS subscriptions, command routing, event observation, awareness state
- `dispatcher.py` (392 lines): DispatchPriority enum, DispatchResult, DispatchHistory, Dispatcher class with single/fan-out/broadcast/smart dispatch

---

### 11. Nexus Endpoints Rewrite

**File**: `backend/nexus/endpoints.py` (REWRITTEN — 448 → ~280 lines)

**Why**: The old endpoints relied heavily on `nexus.legacy` and `nexus.dispatcher` which no longer exist. The rewrite:
- Removes dead endpoints that called into the old Legacy/Dispatcher
- Keeps useful endpoints (command intake, entities, modalities, awareness)
- Adds new heartbeat endpoints
- Rewires the WebSocket stream to connect directly to NATS

**Removed endpoints** (6):
- `POST /nexus/fan-out` — used old Dispatcher.fan_out()
- `POST /nexus/assign` — used old Nexus.assign_task()
- `GET /nexus/events` — used old Legacy.event_log
- `GET /nexus/events/search` — used old Legacy.search_events()
- `GET /nexus/history` — used old Dispatcher.history
- `GET /nexus/legacy/awareness` — used old Legacy.get_awareness()

**Kept endpoints** (6):
- `POST /nexus/command` — simplified, now publishes NexusMessage to NATS and hints to use `/agents/legacy/runs` for direct interaction
- `GET /nexus/overview` — enriched with heartbeat status
- `GET /nexus/entities` — unchanged, lists all entities with filters
- `GET /nexus/modalities` — unchanged, modality inventory
- `GET /nexus/awareness` — unchanged, NATS awareness status
- `POST /nexus/awareness/subscribe` — unchanged, dynamic topic subscription

**New endpoints** (2):
- `GET /nexus/heartbeat/status` — returns heartbeat config, running state, tick count, last tick time, active hours
- `POST /nexus/heartbeat/trigger` — manually triggers a heartbeat tick (for testing/debugging)

**WebSocket rewrite** (`WS /nexus/stream`):
- Old: Monkey-patched `nexus.legacy.event_log.record` and registered on `nexus.dispatcher._listeners`
- New: Subscribes directly to NATS subjects (`nexus.>`, `system.>`, `comms.>`) and forwards events. Cleaner, no monkey-patching, works with NATS directly.

---

### 12. Nexus __init__.py Update

**File**: `backend/nexus/__init__.py` (MODIFIED)

**Why**: The public API of the nexus package changed — `Legacy` and `Dispatcher` classes no longer exist. The __init__.py needed to stop exporting them.

**Before**: Exported `Nexus`, `Legacy`, `FrameworkRegistry`, `Dispatcher`, `NexusMessage`, etc.  
**After**: Exports `Nexus`, `get_nexus`, `reset_nexus`, `FrameworkRegistry`, `NexusMessage`, etc.

Updated the architecture ASCII diagram in the docstring to reflect the new Agno-first hybrid model.

---

## Verification

After all changes:
- **21/21 backend modules** import successfully (tested via `importlib.import_module`)
- **0 frontend errors** in sidebar.dart (verified via diagnostics)
- **0 errors** in all modified nexus files
- No circular imports, no broken references

---

## What's Next

| Priority | Task | Status |
|----------|------|--------|
| 1 | Slim `nexus/registry.py` to FrameworkAdapter ABC only | Not started |
| 2 | Enhanced cron scheduler (DB-backed, replace naive scheduler.py) | Not started |
| 3 | Frontend heartbeat screen (`/heartbeat` route) | Not started |
| 4 | Voice pipeline (wake word → STT → router → LLM → TTS) | Future |
| 5 | Node protocol for hardware (FarmBot, Meshtastic) | Future |

---

## File Summary

| File | Action | Lines Before → After |
|------|--------|---------------------|
| `docs/PROACTIVE_ARCHITECTURE.md` | Created | — → ~420 |
| `backend/heartbeat.py` | Created | — → ~260 |
| `backend/workspace/HEARTBEAT.md` | Created | — → ~20 |
| `backend/tools/nats_tools.py` | Created | — → ~230 |
| `backend/tools/registry.py` | Modified | +3 tool specs, +1 function |
| `backend/agents/presets.py` | Modified | +1 Legacy preset (~60 lines) |
| `backend/main.py` | Modified | +heartbeat startup/shutdown/banner |
| `frontend/lib/widgets/sidebar.dart` | Modified | 7 sections → 5 sections |
| `backend/nexus/core.py` | Rewritten | 276 → ~130 |
| `backend/nexus/endpoints.py` | Rewritten | 448 → ~280 |
| `backend/nexus/__init__.py` | Modified | Removed Legacy/Dispatcher exports |
| `backend/nexus/legacy.py` | Archived | 577 → .archived |
| `backend/nexus/dispatcher.py` | Archived | 392 → .archived |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
