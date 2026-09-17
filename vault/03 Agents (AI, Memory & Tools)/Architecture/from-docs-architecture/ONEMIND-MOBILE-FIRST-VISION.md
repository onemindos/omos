# OneMind OS: Mobile-First Vision & Evolution

> **The Intelligence of You - Everywhere**
> **Last Updated:** 2026-01-29
> **Goal:** Mobile-optimized, agent-powered personal operating system

---

## Executive Summary

This document outlines the evolution of OneMind OS from its current 77-screen desktop-first architecture to a **mobile-first, agent-powered platform** called **1INTEL**. The core insight: most interactions should happen through **natural conversation with intelligent agents**, not navigating through screens.

### Vision Statement

> **"Your entire life managed through a single conversation interface, backed by autonomous agents operating 24/7."**

---

## Part 1: Current State Analysis

### Current Architecture (Problems)

```
CURRENT: 77 Screens across 6 drawer sections
├── UIO (11 screens) - Command Center
├── PATHS (14 screens) - Life Management
├── SENSES (8 screens) - Awareness
├── LEGACY (19 screens) - AI Workforce
├── HIVE (6 screens) - Hardware
└── CORE (4 screens) - System
```

**Issues for Mobile:**
1. **Too Many Screens** - 77 screens is overwhelming on mobile
2. **Drawer Navigation** - Requires multiple taps to reach anything
3. **Desktop-First Layout** - Wide layouts don't work on phone
4. **Scattered Data Entry** - Tasks, habits, goals all separate forms
5. **No Universal Capture** - Different input for every module
6. **Agents Hidden** - AI workforce buried in menus

### Current Bottom Nav (5 items)
```
[Command] [Life OS] [Inbox] [Awareness] [More →]
```

---

## Part 2: Mobile-First Redesign

### Design Philosophy

1. **Conversation-First** - Most actions via natural language
2. **Universal Capture** - One input for everything
3. **Contextual Surfaces** - Right info at right time
4. **Agent-Powered** - Agents do the work, you decide
5. **Glanceable Insights** - Scores and status at a glance
6. **Minimal Navigation** - 3 taps max to anything

### New Bottom Navigation (4 items + Capture)

```
┌─────────────────────────────────────────────────────────┐
│                    SCREEN CONTENT                        │
├─────────────────────────────────────────────────────────┤
│                     [+ CAPTURE]                          │  ← Floating center
├─────────────────────────────────────────────────────────┤
│  [Home]    [Command]    [Execute]    [Profile]          │
│   🏠          💬           ⚡           👤                │
└─────────────────────────────────────────────────────────┘
```

### Core Mobile Screens (12 Primary)

| Screen | Purpose | Features |
|--------|---------|----------|
| **Home** | Dashboard + Scores | PI/LI/WI scores, today's focus, active agents |
| **Command** | Chat Interface | Natural language, @agents, /commands |
| **Execute** | Action Hub | Today's tasks, habits, active workflows |
| **Profile** | Settings & Status | Preferences, agent config, integrations |
| **Capture** | Universal Input | Voice/text/photo → auto-categorized |
| **Tasks** | Task Management | Universal task view (all asset types) |
| **Calendar** | Time View | Schedule, events, time blocks |
| **Insights** | Analytics | Trends, correlations, predictions |
| **Agents** | Agent Control | Active agents, approvals, workflows |
| **Knowledge** | Memory & Docs | Search everything, add context |
| **Goals** | Three Paths | HP/LE/GE progress and OKRs |
| **Settings** | Configuration | Deep settings (accessible from Profile) |

---

## Part 3: Universal Asset System

### The Problem

Currently, different "things" have different screens:
- Tasks → `/tasks`
- Habits → `/habits`
- Goals → `/goals`
- Calendar Events → `/calendar`
- Meals → `/meals`
- Journal Entries → `/journal`
- Projects → `/projects`

**Users don't think in modules. They think in "things to do/track."**

### The Solution: Universal Asset

Every trackable item becomes an **Asset** with a common schema:

```python
class Asset(BaseModel):
    """Universal asset that can represent any trackable item"""

    # Identity
    id: UUID
    asset_type: AssetType  # task, habit, goal, event, meal, journal, project, routine
    title: str
    description: Optional[str]

    # Classification
    pillar: Pillar  # HP, LE, GE
    path: Path  # Specific path within pillar
    priority: Priority  # P0-P3
    energy_level: EnergyLevel  # high, medium, low

    # Timing
    do_date: Optional[datetime]  # When to do it (Motion-style)
    due_date: Optional[datetime]  # Hard deadline
    duration_minutes: Optional[int]
    recurring: Optional[RecurrenceRule]

    # Status
    status: AssetStatus  # pending, in_progress, completed, deferred, archived
    progress: float  # 0.0 to 1.0
    streak: Optional[int]  # For habits

    # Relations
    parent_id: Optional[UUID]  # For subtasks/sub-goals
    project_id: Optional[UUID]
    linked_assets: List[UUID]

    # Metadata
    tags: List[str]
    context: Dict[str, Any]  # Flexible context
    source: str  # Where it came from (capture, agent, sync, etc.)

    # Timestamps
    created_at: datetime
    updated_at: datetime
    completed_at: Optional[datetime]
```

### Asset Types & Views

```
ASSET TYPES
├── task          → Completable action
├── habit         → Recurring tracked behavior
├── goal          → Long-term objective
├── event         → Calendar occurrence
├── meal          → Food/nutrition tracking
├── journal       → Reflection/log entry
├── project       → Container for tasks
├── routine       → Sequence of habits/tasks
├── note          → Knowledge capture
└── reminder      → Time-triggered notification
```

### Single Unified View

Instead of 10 different screens, one **Smart View** with filters:

```
┌─────────────────────────────────────────────────────────┐
│ ⚡ Execute                              🔍 ⚙️ Filter    │
├─────────────────────────────────────────────────────────┤
│ [All] [Tasks] [Habits] [Goals] [Events] [Today] [HP]   │
├─────────────────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ☐ Complete quarterly review           HP • P1 • 2h │ │
│ │ 🔄 Morning meditation                 HP • Streak 23│ │
│ │ 🎯 Launch MVP                         GE • 45%     │ │
│ │ 📅 Team standup 10am                  GE • 30min   │ │
│ │ 🍽️ Log breakfast                      HP • Pending │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## Part 4: Agent Architecture at Scale

### How Agno Connects to UI

```
┌─────────────────────────────────────────────────────────┐
│                      USER INTERFACE                      │
│   [Command Chat] [Capture] [Execute] [Insights]         │
├─────────────────────────────────────────────────────────┤
│                    AGENT GATEWAY                         │
│   Natural Language → Intent → Route → Agent → Action    │
├─────────────────────────────────────────────────────────┤
│                    AGENT LAYER                           │
│   ┌─────────┐   ┌─────────┐   ┌─────────┐              │
│   │ Mother  │ → │ Router  │ → │Specialist│              │
│   │ (Legacy)│   │         │   │ Agents  │              │
│   └─────────┘   └─────────┘   └─────────┘              │
│        ↓             ↓             ↓                    │
│   ┌─────────┐   ┌─────────┐   ┌─────────┐              │
│   │ Teams   │   │Workflows│   │  Tools  │              │
│   └─────────┘   └─────────┘   └─────────┘              │
├─────────────────────────────────────────────────────────┤
│                    DATA LAYER                            │
│   [Assets] [Memory] [Knowledge] [Context] [Scores]      │
├─────────────────────────────────────────────────────────┤
│                    EXTERNAL                              │
│   [Todoist] [Calendar] [Home Assistant] [Finance]       │
└─────────────────────────────────────────────────────────┘
```

### Agent Roles & Visibility

| Agent | Purpose | User Sees |
|-------|---------|-----------|
| **Legacy (Mother)** | Coordinator, router, memory | Main chat interface |
| **Planning Agent** | Task breakdown, scheduling | Suggestions in Execute |
| **Research Agent** | Information gathering | Research results |
| **Code Agent** | Development tasks | Code suggestions |
| **Analysis Agent** | Data analysis, insights | Insights screen |
| **Writer Agent** | Content creation | Draft content |
| **Health Agent** | HP pillar monitoring | Health insights |
| **Finance Agent** | GE pillar, investments | Finance insights |

### Autonomous Agent Operations

Agents run **autonomously** in background:

```
AUTONOMOUS OPERATIONS (24/7)
├── Morning Briefing (6am)
│   └── Summarize day, priorities, calendar
├── Inbox Triage (continuous)
│   └── Categorize, suggest actions
├── Task Optimization (hourly)
│   └── Reschedule, optimize energy matching
├── Health Monitoring (continuous)
│   └── Oura data, activity tracking
├── Evening Review (9pm)
│   └── Daily summary, tomorrow prep
└── Weekly Planning (Sunday 8am)
    └── Week review, goal progress
```

### Context Engine (Full Context)

```python
class ContextEngine:
    """Aggregates ALL context for agent operations"""

    def build_context(self) -> AgentContext:
        return AgentContext(
            # Identity
            user_profile=self.get_user_profile(),

            # Current State
            location=self.get_location(),       # Home Assistant
            time_context=self.get_time_context(),
            device=self.get_device_context(),

            # Life State
            today_assets=self.get_today_assets(),
            active_projects=self.get_active_projects(),
            current_goals=self.get_current_goals(),

            # Historical
            recent_completions=self.get_recent_completions(),
            patterns=self.get_patterns(),
            preferences=self.get_preferences(),

            # Scores
            pi_score=self.get_pi_score(),
            li_score=self.get_li_score(),
            wi_score=self.get_wi_score(),

            # Health
            biometrics=self.get_oura_data(),
            sleep_quality=self.get_sleep_quality(),
            energy_level=self.get_energy_level(),

            # External
            calendar_events=self.get_calendar(),
            weather=self.get_weather(),
            commute=self.get_commute_status(),
        )
```

---

## Part 5: New UX Flows

### Flow 1: Universal Capture

```
User taps [+ Capture]
    ↓
┌─────────────────────────────────────────────────────────┐
│                    CAPTURE                               │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐ │
│  │ What's on your mind?                               │ │
│  │ __________________________________________________ │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                          │
│  [🎤 Voice] [📷 Photo] [📎 File] [⏰ Remind]            │
│                                                          │
│  Quick: [Task] [Habit] [Note] [Event]                   │
├─────────────────────────────────────────────────────────┤
│  Recent Captures:                                        │
│  • "Call mom about birthday" → Task (LE) ✓             │
│  • Photo of whiteboard → Note (GE) ✓                   │
└─────────────────────────────────────────────────────────┘
```

**Agent Processing:**
```
User: "Remind me to review Q1 financials next week"
    ↓
Legacy Agent parses:
  - Intent: Create reminder/task
  - Content: Review Q1 financials
  - When: Next week (DO DATE algorithm)
  - Pillar: GE (finance-related)
  - Priority: P1 (inferred from content)
    ↓
Creates Asset (type: task)
    ↓
Confirms: "Created 'Review Q1 financials' for Monday 9am (GE pillar)"
```

### Flow 2: Smart Command

```
User opens Command screen
    ↓
┌─────────────────────────────────────────────────────────┐
│ 💬 Command                            Context: Home 🏠  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Legacy: Good morning Zeus. You have 5 tasks for       │
│  today (3 HP, 2 GE). Your PI score dropped 3 points    │
│  yesterday - want me to analyze why?                    │
│                                                          │
│  Quick Actions:                                          │
│  [📊 Analyze] [📋 Show Tasks] [🔄 Reschedule]          │
│                                                          │
├─────────────────────────────────────────────────────────┤
│  @ Research what happened to PI score                   │
│  ________________________________________________[Send]│
│                                                          │
│  [@Agent] [/Command] [+ Attach]                         │
└─────────────────────────────────────────────────────────┘
```

**@Mentions route to specialists:**
- `@research` → Research Agent
- `@code` → Code Agent
- `@plan` → Planning Agent
- `@health` → Health Agent
- `@finance` → Finance Agent

**/Commands for quick actions:**
- `/task Create new task` → Direct task creation
- `/schedule tomorrow` → Show tomorrow's schedule
- `/focus 30` → Start 30-min focus session
- `/brief` → Get current briefing
- `/goals` → Show goal progress

### Flow 3: Execute View (Today)

```
┌─────────────────────────────────────────────────────────┐
│ ⚡ Execute                    Wed Jan 29 • PI:87 🟢     │
├─────────────────────────────────────────────────────────┤
│ Focus Block: Deep Work 🧠                 9:00-11:00   │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ☐ Code review for PR #123            GE • 45min    │ │
│ │ ☐ Write project proposal             GE • 30min    │ │
│ └─────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────┤
│ Habits Today                              5/7 done     │
│ [✓ Meditate] [✓ Exercise] [☐ Journal] [✓ Read]       │
├─────────────────────────────────────────────────────────┤
│ Up Next                                                 │
│ • 11:00 Team standup (30min)                           │
│ • 12:00 Lunch break                                    │
│ • 13:00 Client call (1h)                               │
├─────────────────────────────────────────────────────────┤
│ Suggested by Legacy:                                    │
│ "Move 'Budget review' to Friday - energy low today"    │
│ [Accept] [Dismiss] [Modify]                            │
└─────────────────────────────────────────────────────────┘
```

### Flow 4: Home Dashboard

```
┌─────────────────────────────────────────────────────────┐
│ OneMind                           🔔 3  ⚡ Awareness    │
├─────────────────────────────────────────────────────────┤
│                                                          │
│         ┌───────────────────────────────┐               │
│         │     UNIFIED SCORE: 84         │               │
│         │         ████████░░            │               │
│         └───────────────────────────────┘               │
│                                                          │
│    PI: 87 🟢       LI: 82 🟡       WI: 83 🟢           │
│    Health         Legacy         Wealth                 │
│                                                          │
├─────────────────────────────────────────────────────────┤
│ Today's Focus                                           │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ 🎯 Complete MVP launch prep                          │ │
│ │ Progress: ████████░░ 78%                            │ │
│ └─────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────┤
│ Active Agents                          3 working       │
│ • Planning Agent: Optimizing tomorrow's schedule       │
│ • Research Agent: Gathering market data                │
│ • Health Agent: Analyzing sleep patterns               │
│                                                          │
│ [View All Agents →]                                     │
├─────────────────────────────────────────────────────────┤
│ Quick Stats                                             │
│ Tasks: 5 today • Habits: 7/10 • Streak: 23 days       │
└─────────────────────────────────────────────────────────┘
```

---

## Part 6: Linear Issue Mapping

### New Issues Needed for Mobile-First

| Issue | Title | Epic | Priority | Points |
|-------|-------|------|----------|--------|
| OMOS-470 | Mobile Navigation Redesign | Frontend | P0 | 8 |
| OMOS-471 | Universal Asset Schema | Backend | P0 | 5 |
| OMOS-472 | Universal Capture Component | Frontend | P0 | 5 |
| OMOS-473 | Asset Unified API | Backend | P0 | 5 |
| OMOS-474 | Home Dashboard Mobile | Frontend | P0 | 5 |
| OMOS-475 | Execute View (Today) | Frontend | P0 | 5 |
| OMOS-476 | Command Chat Mobile | Frontend | P0 | 5 |
| OMOS-477 | Agent Visibility Layer | Frontend | P1 | 3 |
| OMOS-478 | @Mention Agent Routing | Backend | P1 | 3 |
| OMOS-479 | /Command Quick Actions | Backend | P1 | 3 |
| OMOS-480 | Insights Mobile View | Frontend | P1 | 3 |
| OMOS-481 | Profile & Settings Mobile | Frontend | P1 | 3 |
| OMOS-482 | Voice Capture Integration | Frontend | P1 | 3 |
| OMOS-483 | Photo Capture Processing | Backend | P1 | 3 |
| OMOS-484 | Smart Suggestions Widget | Frontend | P1 | 3 |
| OMOS-485 | Agent Status Dashboard | Frontend | P2 | 3 |

### Existing Issues to Combine/Update

| Existing | Combine With | New Focus |
|----------|--------------|-----------|
| OMOS-460-465 | Life OS | Part of Universal Asset |
| OMOS-320-323 | Smart Chat | Mobile Command UI |
| OMOS-370-373 | Command Center | Mobile Home/Execute |
| OMOS-350-353 | Inbox | Universal Capture |

---

## Part 7: Agent Build Assignments

### 8 Parallel Agents for Build

```
┌─────────────────────────────────────────────────────────┐
│                    BUILD SWARM (8 Agents)               │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  AGENT 1: Backend Architecture                          │
│  └─ Universal Asset schema, migrations, API             │
│                                                          │
│  AGENT 2: Context Engine                                │
│  └─ Full context aggregation, Oura, Home Assistant      │
│                                                          │
│  AGENT 3: Agent Routing                                 │
│  └─ @mentions, /commands, Mother Agent orchestration    │
│                                                          │
│  AGENT 4: Mobile Navigation                             │
│  └─ Bottom nav, new routing, screen reorganization      │
│                                                          │
│  AGENT 5: Home & Execute Screens                        │
│  └─ Dashboard, scores, today view, focus blocks         │
│                                                          │
│  AGENT 6: Command & Capture                             │
│  └─ Chat UI, universal capture, voice integration       │
│                                                          │
│  AGENT 7: Insights & Profile                            │
│  └─ Analytics views, settings, agent config             │
│                                                          │
│  AGENT 8: Integration & Polish                          │
│  └─ Wire everything, transitions, animations            │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Agent 1: Backend Architecture (Opus 4.5)

**Scope:**
- Universal Asset model and migrations
- Asset CRUD API endpoints
- Backward compatibility layer (existing modules → assets)
- Asset search and filtering

**Files:**
```
backend/lifeos/assets/
├── models.py          # Asset schema
├── service.py         # Business logic
├── api.py            # REST endpoints
├── migrations/       # Database migrations
└── compatibility.py  # Legacy module mapping
```

**Dependencies:** None (foundational)

### Agent 2: Context Engine (Opus 4.5)

**Scope:**
- Full context aggregator
- Oura biometrics integration
- Home Assistant context
- Score calculations (PI/LI/WI)

**Files:**
```
backend/agno/context/
├── engine.py         # Main context engine
├── aggregators/      # Data aggregators
│   ├── assets.py
│   ├── biometrics.py
│   ├── location.py
│   └── scores.py
└── api.py           # Context API
```

**Dependencies:** Agent 1 (Assets)

### Agent 3: Agent Routing (Opus 4.5)

**Scope:**
- @mention parsing and routing
- /command registration and execution
- Mother Agent orchestration
- Workflow triggers

**Files:**
```
backend/agno/routing/
├── parser.py        # Message parsing
├── router.py        # Agent routing
├── commands.py      # /command registry
└── orchestrator.py  # Mother Agent logic
```

**Dependencies:** Agent 2 (Context)

### Agent 4: Mobile Navigation (Sonnet)

**Scope:**
- New bottom nav (4 items + capture)
- Route consolidation (77 → 12 primary)
- Navigation provider updates
- Drawer simplification

**Files:**
```
frontend/lib/
├── platform/router/app_router.dart    # Updated routes
├── platform/providers/navigation_provider.dart
├── shared/widgets/app_shell.dart      # New bottom nav
└── shared/widgets/app_drawer.dart     # Simplified
```

**Dependencies:** None (can start immediately)

### Agent 5: Home & Execute Screens (Sonnet)

**Scope:**
- Home dashboard with scores
- Execute view (today's focus)
- Asset list component (universal)
- Focus block widget

**Files:**
```
frontend/lib/
├── home/
│   ├── screens/home_screen.dart
│   └── widgets/
│       ├── score_card.dart
│       ├── active_agents.dart
│       └── quick_stats.dart
├── execute/
│   ├── screens/execute_screen.dart
│   └── widgets/
│       ├── focus_block.dart
│       ├── habit_row.dart
│       └── up_next.dart
└── shared/widgets/asset_list.dart
```

**Dependencies:** Agent 4 (Navigation)

### Agent 6: Command & Capture (Sonnet)

**Scope:**
- Mobile-optimized chat UI
- Universal capture modal
- Voice input integration
- Quick action buttons

**Files:**
```
frontend/lib/
├── command/
│   ├── screens/command_screen.dart
│   └── widgets/
│       ├── chat_input.dart
│       ├── suggestion_chips.dart
│       └── agent_indicator.dart
├── capture/
│   ├── screens/capture_modal.dart
│   └── widgets/
│       ├── voice_capture.dart
│       ├── photo_capture.dart
│       └── quick_types.dart
```

**Dependencies:** Agent 3 (Routing), Agent 4 (Navigation)

### Agent 7: Insights & Profile (Sonnet)

**Scope:**
- Insights/analytics mobile view
- Profile screen with settings
- Agent configuration UI
- Integration status

**Files:**
```
frontend/lib/
├── insights/
│   ├── screens/insights_screen.dart
│   └── widgets/
│       ├── trend_chart.dart
│       ├── correlation_card.dart
│       └── prediction_widget.dart
├── profile/
│   ├── screens/profile_screen.dart
│   └── widgets/
│       ├── settings_section.dart
│       ├── agent_config.dart
│       └── integration_status.dart
```

**Dependencies:** Agent 4 (Navigation)

### Agent 8: Integration & Polish (Sonnet)

**Scope:**
- Wire all components together
- Screen transitions and animations
- Error handling and loading states
- Final polish and testing

**Files:**
```
frontend/lib/
├── shared/
│   ├── animations/
│   ├── transitions/
│   └── error_handling/
└── platform/
    └── services/
        └── asset_service.dart
```

**Dependencies:** Agents 1-7 (final integration)

---

## Part 8: Build Timeline

### Week 1: Foundation

| Day | Agent 1 | Agent 2 | Agent 3 | Agent 4 |
|-----|---------|---------|---------|---------|
| Mon | Asset schema | Context design | Routing design | Nav restructure |
| Tue | Asset API | Aggregators | Parser impl | Route consolidation |
| Wed | Migrations | Oura integration | Commands | Bottom nav |
| Thu | Compatibility | HA context | Router impl | Drawer simplify |
| Fri | Testing | Score calc | Orchestrator | Navigation provider |

### Week 2: UI Build

| Day | Agent 5 | Agent 6 | Agent 7 | Agent 8 |
|-----|---------|---------|---------|---------|
| Mon | Home screen | Chat UI | Insights screen | Integration setup |
| Tue | Score cards | Capture modal | Trends | Wire home |
| Wed | Execute view | Voice capture | Profile | Wire execute |
| Thu | Focus blocks | Photo capture | Agent config | Wire command |
| Fri | Asset list | Quick actions | Settings | Wire insights |

### Week 3: Integration

| Day | All Agents |
|-----|------------|
| Mon | Full integration testing |
| Tue | Bug fixes and refinements |
| Wed | Performance optimization |
| Thu | Animation polish |
| Fri | Final testing and review |

---

## Part 9: After State Vision

### What OneMind Looks Like After

```
BEFORE (Desktop-First)                 AFTER (Mobile-First)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
77 screens                      →      12 primary screens
6 drawer sections              →      4 bottom nav + capture
Navigate to create             →      Capture anywhere
Agents hidden in menu          →      Agents visible & active
Separate task/habit/goal      →      Universal assets
Desktop layouts               →      Mobile-optimized
Manual scheduling             →      AI-powered DO DATE
Passive data display          →      Active suggestions
Module-based thinking         →      Conversation-first

User Experience:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"Open drawer → expand PATHS  →      "Hey Legacy, what should
 → tap Tasks → find task →           I focus on today?"
 create → fill form → save"
                                     "Got it. Starting focus
                                      block for MVP prep."
```

### Mobile Screen Hierarchy

```
PRIMARY (Bottom Nav)
├── Home                    ← Dashboard, scores, agents
├── Command                 ← Chat with Legacy
├── Execute                 ← Today's tasks/habits
├── Profile                 ← Settings, config
└── [+ Capture]            ← Universal input (center FAB)

SECONDARY (From Primary)
├── Tasks (from Execute)    ← Full asset list
├── Calendar (from Home)    ← Time-based view
├── Insights (from Home)    ← Analytics deep dive
├── Goals (from Home)       ← Three Paths progress
├── Agents (from Home)      ← Agent management
├── Knowledge (from Command)← Memory and docs
└── Settings (from Profile) ← Deep configuration

MODALS (Overlay)
├── Capture                 ← Quick input
├── Asset Detail           ← View/edit any asset
├── Agent Approval         ← HITL decisions
└── Quick Actions          ← Context menu
```

### Key Metrics After Build

| Metric | Before | After |
|--------|--------|-------|
| Screens to navigate | 77 | 12 |
| Taps to create task | 5+ | 1 (capture) |
| Agent visibility | 0% | 100% |
| Voice input | No | Yes |
| Photo capture | No | Yes |
| Autonomous agents | Manual | 24/7 |
| Mobile optimization | 30% | 100% |
| Context awareness | Partial | Full |

---

## Part 10: Linear Issues Summary

### Total Issues After Evolution

```
EXISTING ISSUES (50)
├── Context Engine (6)      → Keep, enhance
├── Intelligence Scores (5) → Keep, integrate
├── Smart Chat (4)          → Evolve to Command
├── Scheduling (4)          → Keep
├── Workflow Automation (4) → Keep
├── Inbox Intelligence (4)  → Evolve to Capture
├── Operating Rhythm (4)    → Keep
├── Frontend (4)            → Replace with new
├── Life OS Completion (6)  → Merge into Assets
├── Advanced Features (9)   → Keep

NEW ISSUES (16)
├── Mobile Navigation (1)
├── Universal Assets (3)
├── Mobile Screens (6)
├── Voice/Photo Capture (2)
├── Agent Visibility (2)
└── Integration (2)

TOTAL: 66 issues
```

### Priority Matrix

```
P0 (Must Have for Mobile Launch)
├── Universal Asset Schema
├── Mobile Navigation Redesign
├── Universal Capture
├── Home Dashboard Mobile
├── Execute View
├── Command Chat Mobile
└── Context Engine (existing)

P1 (Core Experience)
├── Agent Visibility Layer
├── @Mention Routing
├── /Command System
├── Insights Mobile
├── Voice Capture
├── Photo Processing
└── Scheduling (existing)

P2 (Enhancement)
├── Agent Status Dashboard
├── Advanced Workflows
├── Operating Rhythm
└── Memory Blocks
```

---

## Conclusion

This evolution transforms OneMind OS from a **feature-rich but complex desktop application** into a **mobile-first, agent-powered personal operating system**.

The key shifts:
1. **From screens to conversations** - Command is primary interface
2. **From modules to assets** - Everything is a universal trackable item
3. **From manual to autonomous** - Agents work 24/7
4. **From hidden to visible** - See what agents are doing
5. **From navigate to capture** - One-tap input for everything

With 8 agents working in parallel, we can complete this evolution in **3 weeks**, resulting in a truly mobile-optimized platform where the primary interaction is natural conversation with intelligent agents who have full context of your life.

---

**Ready to build?** Start by creating the Linear issues, then deploy the 8-agent swarm.
