---
id: so-01-openclaw-system-specs-20260501
domain: SO
area: "01"
type: reference
status: active
agent: Legacy
tags: [domain/so, area/01, type/reference, status/active, openclaw/system]
ai_summary: "OpenClaw system architecture — Asset schema, Context Engine, routing, memory, API endpoints, mobile design system. Extracted from implementation sprint docs before deletion."
updated: 2026-05-01
---

# OpenClaw System Specs — Reference

> *Technical architecture reference. Extracted from the Jan 2026 build sprint docs before `implementation/` folder deletion.*

---

## Universal Asset Schema

### Asset Types (10)
`task` | `habit` | `goal` | `event` | `meal` | `journal` | `project` | `routine` | `note` | `reminder`

### Priority Scale
| Value | Label | Default |
|-------|-------|---------|
| `p0` | Critical | — |
| `p1` | High | — |
| `p2` | Medium | ✅ default |
| `p3` | Low | — |

### Source Values
`manual` | `capture` | `agent` | `sync`

### Key Fields
| Field | Type | Notes |
|-------|------|-------|
| `do_date` | datetime | When to DO it (scheduled work time) |
| `due_date` | datetime | Hard deadline |
| `duration_minutes` | int | Estimated time |
| `energy_level` | enum | high / medium / low |
| `pillar` | enum | hp / le / ge |
| `progress` | float | 0.0 → 1.0 |
| `streak` | int | Habits only |
| `parent_id` | UUID | Subtasks / sub-goals |
| `linked_assets` | UUID[] | Cross-asset links |
| `context` | JSONB | Freeform metadata |
| `source` | string | manual / capture / agent / sync |

### DB Indexes (7)
`asset_type`, `pillar`, `status`, `do_date`, `due_date`, `project_id`, + full-text gin index on `title + description`

### Backward Compatibility
Old endpoints `/api/tasks`, `/api/habits`, `/api/goals` still work — they are shims into the Asset system via `compatibility.py`. No breaking changes.

---

## Context Engine

### Performance Targets
- `build()` full context: **< 500ms**
- `build_minimal()` state + today items only: **< 100ms**
- Cache TTL: **60 seconds**

### 7 Aggregators (run in parallel)
| Aggregator | Data Source | Produces |
|------------|-------------|---------|
| `state` | Home Assistant + DB | Location, energy, focus mode, active pillar |
| `items` | Asset service | Tasks, habits, goals — today |
| `scores` | Score service | PI/LI/WI/Unified + trends |
| `calendar` | CalDAV | Today + tomorrow events, free slots |
| `history` | DB | Recent completions, conversations, patterns |
| `assets` | HA + DB | Home devices, wearables, financial summary, vehicles |
| `biometrics` | Oura API | Sleep, readiness, HRV, resting HR |

**Failure handling:** `return_exceptions=True` — if one aggregator fails, the rest still return. No cascading failure.

### Oura → Energy Mapping
| Readiness Score | Energy Level |
|-----------------|-------------|
| ≥ 85 | `high` |
| 65–84 | `medium` |
| < 65 | `low` |

### Oura Fields Captured
`sleep_score`, `sleep_duration` (hrs), `readiness_score`, `hrv_balance`, `resting_heart_rate`

**Recommendation rules:**
- sleep_score < 70 → "Sleep score low - consider earlier bedtime"
- readiness_score < 70 → "Recovery low - recommend lighter day"

### Context Prompt Format
Exact format injected into every AI system prompt:

```
## Current State
- Location: {home|office|mobile|unknown}
- Energy Level: {high|medium|low}
- Active Focus: {focus_mode or null}
- Current Path: {hp|le|ge or null}

## Intelligence Scores
- PI (Performance): {score}/100 ({up|down|stable})
- LI (Legacy): {score}/100 ({trend})
- WI (Wealth): {score}/100 ({trend})
- Unified: {score}/100

## Today's Items
Tasks ({count}):
  - {title} ({priority})
  ... (top 5)

Habits ({completed}/{total}):
  ✓/○ {title} | ...

## Calendar
  - {HH:MM}: {title}
  ... (today's events, top 5)

## Biometrics (Oura)
- Sleep Score: {score}/100
- Readiness: {score}/100
- HRV: {hrv} ms
```

### Photo Capture
Model: `gpt-4-vision-preview`
Returns JSON: `{extracted_text, context, suggested_type, suggested_title, suggested_pillar, confidence}`

---

## Agent Routing & Coordination

### @Mention Alias Map
| Alias(es) | Routes To |
|-----------|-----------|
| `legacy`, `l` | Legacy agent |
| `research`, `researcher`, `r` | Research agent |
| `code`, `coder`, `c` | Code agent |
| `plan`, `planner`, `planning`, `p` | Plan agent |
| `health`, `hp` | Health agent |
| `finance`, `money`, `ge` | Finance agent |
| `write`, `writer`, `w` | Write agent |
| `analyst`, `analyze`, `analysis` | Analyst agent |

Unknown alias → falls back to Legacy.

### /Command List (8 commands)
| Command | Action |
|---------|--------|
| `/brief` | Generate morning briefing |
| `/tasks` | Show today's tasks |
| `/habits` | Show today's habits |
| `/scores` | Show current PI/LI/WI scores |
| `/schedule` | Show today's schedule |
| `/focus` | Start focus session |
| `/add {text}` | Quick add a task |
| `/help` | List all commands |

### Smart Chat Processing Priority
1. `/command` → execute command handler
2. `@mention` → route to specialist(s)
3. Workflow detection → multi-step plan
4. Normal → chat with Legacy

### Mother Agent (Legacy) Execution Modes
- **DIRECT** — handle request without specialists
- **RESEARCH** — delegate to research agent
- **CODE** — delegate to code agent
- **MULTI** — parallel/sequential specialist execution with review + synthesis

### Risk Detection System
| Risk Level | Trigger Conditions | Mode |
|------------|-------------------|------|
| LOW | Read operations | `autorun` |
| MEDIUM | Create / update; affected_count > 10 | `agent_decides` |
| HIGH | Delete, external comms (email/SMS/webhook), irreversible | `approval` |
| CRITICAL | Financial (transfer/payment/purchase), credentials/secrets | `approval` (always block) |

### Memory Architecture

**Core Memory Blocks** (always in context):
| Block | Contents |
|-------|----------|
| `zeus_profile` | Name, role, communication style, goals |
| `pillar_hp` | Health focus, key habits, HPI target |
| `pillar_le` | Family priorities, weekly rituals |
| `pillar_ge` | Business name, focus, key metrics |
| `onemind_sys` | System name, version, features |

**Archival Memory:**
- Semantic search via pgvector (`embedding <=>` cosine distance)
- Default top_k: 5 results
- Stores: past decisions, preferences, insights
- Table: `archival_memory (user_id, content, embedding, metadata)`

**Memory Block table:** `memory_blocks (user_id, block_name, content, updated_at)`

---

## API Endpoint Map

### Assets
| Endpoint | Method | Action |
|----------|--------|--------|
| `/api/v1/assets` | GET | List with filters (type/pillar/status/dates) |
| `/api/v1/assets` | POST | Create any asset type |
| `/api/v1/assets/today` | GET | All assets for today |
| `/api/v1/assets/search?q=` | GET | Full-text search |
| `/api/v1/assets/{id}` | GET | Single asset |
| `/api/v1/assets/{id}` | PATCH | Update |
| `/api/v1/assets/{id}` | DELETE | Delete |
| `/api/v1/assets/{id}/complete` | POST | Mark complete |

### Scores
| Endpoint | Method | Action |
|----------|--------|--------|
| `/api/v1/scores/all` | GET | All scores at once |
| `/api/v1/scores/hpi` | GET | PI score |
| `/api/v1/scores/lei` | GET | LI score |
| `/api/v1/scores/gei` | GET | WI score |
| `/api/v1/scores/unified` | GET | Unified index |
| `/api/v1/scores/history?days=30` | GET | Score history for trending |

### Scheduling
| Endpoint | Method | Action |
|----------|--------|--------|
| `/api/v1/scheduling/tasks/{id}/schedule` | POST | Get optimal DO DATE slots (returns top 3) |
| `/api/v1/scheduling/balance` | GET | Current HP/LE/GE path balance |
| `/api/v1/scheduling/auto-reschedule` | POST | Trigger auto-reschedule check |

### Context
| Endpoint | Method | Action |
|----------|--------|--------|
| `/api/v1/context` | GET | Full unified context (<500ms) |
| `/api/v1/context/minimal` | GET | Quick context (<100ms) |
| `/api/v1/context/prompt` | GET | Context as formatted prompt text |
| `/api/v1/context/invalidate` | POST | Invalidate cache |

### Other
| Endpoint | Method | Action |
|----------|--------|--------|
| `/api/v1/focus/today` | GET | Active focus session |
| `/api/v1/agents/active` | GET | Real-time agent status list |
| `/api/v1/stats/today` | GET | Quick stats |
| `/api/v1/suggestions/next` | GET | Next AI suggestion |
| `/api/v1/workflows/trigger/{token}` | POST | Webhook trigger (HMAC-SHA256 via `X-Signature`) |

### Workflow Webhook Security
HMAC-SHA256: `hmac.new(secret, str(payload), sha256).hexdigest()` — verified against `X-Signature` header.

---

## Cross-Domain Insights Engine

**4 insight types:**
1. **Correlations** — "When X, Y happens" (Pearson r > 0.5, p < 0.05)
2. **Anomalies** — "X is unusual" (z-score > ±2 std from 7-day mean)
3. **Trends** — "X is declining" (first vs second half of 7-day window, >±5pt change)
4. **Spillovers** — "X affects Y" (cross-domain causal patterns)

Returns top 10 insights sorted by `impact` score.

---

## Mobile Design System

### Color Tokens
| Token | Hex | Usage |
|-------|-----|-------|
| Background | `#0A0A0A` | Screen backgrounds |
| Cards | `#1A1A1A` | Cards, bottom nav, drawer |
| Surface | `#2A2A2A` | Input fields |
| Borders | `#333333` | Dividers, card borders |
| Accent | `#FF0000` | FAB, active nav, CTAs (tactical red) |
| Habit complete | `#00FF00` | Habit chip completed state |

### Navigation Structure
**Bottom Nav (4 items + FAB):**
```
[Home 🏠]  [Command 💬]  [+FAB]  [Execute ⚡]  [Profile 👤]
```
- FAB: centered, circular notch, `notchMargin: 8`
- FAB opens Universal Capture modal (DraggableScrollableSheet, init 0.6, max 0.9)

**Drawer (3 sections):**
- **Life OS:** Tasks / Calendar / Goals / Habits
- **Intelligence:** Insights / Agents / Knowledge
- **System:** Home Assistant / Settings

**Route map:**
| Old | New |
|-----|-----|
| `/` | → `/home` |
| `/dashboard` | → `/home` |
| `/chat` | → `/command` |
| `/today` | → `/execute` |

**Secondary routes:** `/tasks`, `/calendar`, `/goals`, `/habits`, `/insights`, `/agents`, `/knowledge`, `/home-assistant`, `/settings`

### Universal Capture Flow
1. User opens modal (FAB or programmatic)
2. Types text / records voice / takes photo
3. Taps "Save" → **AI classifies first** → shows classification preview
4. User accepts or modifies classification
5. Asset created and saved

**Capture input types:** text, voice (mic), photo (camera), file attachment
**Quick type buttons:** Task / Note / Event

### Quick Action Chips (empty Command screen)
`📊 Analyze → /scores` | `📋 Tasks → /tasks` | `🔄 Schedule → /schedule` | `📝 Brief → /brief`

---

## Backend Module Map

```
backend/
├── lifeos/
│   ├── assets/         # Universal Asset (models, service, api, compatibility, migrations)
│   ├── scores/         # PI/LI/WI calculators + history
│   ├── scheduling/     # DO DATE, auto-reschedule, energy matching, path balance
│   ├── insights/       # Cross-domain insights engine (correlations, anomalies, trends)
│   ├── tasks/          # Legacy CRUD (shims to assets)
│   ├── habits/         # Legacy CRUD (shims to assets)
│   ├── goals/          # Legacy CRUD (shims to assets)
│   ├── calendar/       # CalDAV integration
│   ├── inbox/          # GTD-style processing + AI triage
│   ├── journal/        # Daily reflection
│   ├── meals/          # Nutritional tracking
│   ├── planner/        # Time block management (move/resize/swap/assign)
│   └── scheduling/     # Event scheduling
├── consciousness/
│   ├── bus.py          # NATS JetStream event bus
│   ├── kv_store.py     # NATS KV Store for state
│   ├── awareness/      # StateManager, device presence
│   └── context/        # ContextEngine + 7 aggregators
├── agno/
│   ├── routing/        # Parser, router, commands, orchestrator (SmartChatEngine)
│   ├── coordination/   # Mother agent, handoffs, specialist definitions
│   ├── safety/         # RiskDetector
│   ├── memory/         # Core blocks + archival (pgvector)
│   ├── workflows/      # HTTP/schedule/event triggers, approval engine
│   └── capture/        # PhotoProcessor (GPT-4V)
└── integrations/
    ├── homeassistant/  # HA client
    └── oura/           # Oura ring client
```

### Sprint Stats (completed Jan 2026)
- 56 Linear issues (OMOS-301 → OMOS-485)
- 181 story points
- 8 parallel Claude agents (Opus 4.5 for backend, Sonnet for frontend)
- 3-week sprint: Foundation → UI Build → Integration


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
