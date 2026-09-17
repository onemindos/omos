# OneMind Mobile-First Sprint

> **Sprint Start:** 2026-01-29
> **Total Issues:** 56 active
> **Total Points:** 181
> **Agents:** 8 parallel agents
> **Goal:** Transform OneMind into mobile-first, agent-powered OS

---

## Sprint Overview

### Transformation Summary
```
BEFORE                              AFTER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
77 screens                    →     12 primary screens
6 drawer sections             →     4 bottom nav + capture FAB
Separate task/habit/goal      →     Universal Assets
Agents hidden in menus        →     Agents always visible
Manual everything             →     AI-powered automation
Desktop layouts               →     Mobile-optimized
Navigate to create            →     Capture anywhere
Module-based thinking         →     Conversation-first
```

### Architecture Stack
```
┌─────────────────────────────────────────────────────────────┐
│                       MOBILE UI (Flutter)                    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  Home   │ │ Command │ │ Execute │ │ Profile │ + Capture │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
├─────────────────────────────────────────────────────────────┤
│                    AGNO AGENT LAYER                          │
│  Legacy (Mother) → Routes to Specialists via @mentions       │
│  /Commands → Quick Actions | Workflows → Automation          │
├─────────────────────────────────────────────────────────────┤
│                    CONTEXT ENGINE                            │
│  State + Assets + Scores + Calendar + History + Biometrics   │
├─────────────────────────────────────────────────────────────┤
│                    DATA LAYER (Life OS)                      │
│  Universal Assets | Memory | Knowledge                       │
├─────────────────────────────────────────────────────────────┤
│                    INTEGRATIONS                              │
│  [Todoist] [Calendar] [Oura] [Home Assistant] [Finance]     │
└─────────────────────────────────────────────────────────────┘
```

---

## Existing Codebase Status

### Backend - READY TO BUILD ON

| Module | Status | Location |
|--------|--------|----------|
| Tasks | ✅ Production | `backend/lifeos/tasks/` |
| Habits | ✅ Production | `backend/lifeos/habits/` |
| Goals | ✅ Production | `backend/lifeos/goals/` |
| Calendar | ✅ Complete | `backend/lifeos/calendar/` |
| Inbox | ✅ Complete | `backend/lifeos/inbox/` |
| Journal | ✅ Complete | `backend/lifeos/journal/` |
| Meals | ✅ Complete | `backend/lifeos/meals/` |
| Notifications | ✅ Complete | `backend/lifeos/notifications/` |
| Agents | ✅ Mature | `backend/agno/agents/` |
| Teams | ✅ Ready | `backend/agno/teams/` |
| Workflows | ✅ Ready | `backend/agno/workflows/` |
| Memory | ✅ Ready | `backend/agno/memory/` |
| Consciousness | ✅ Ready | `backend/consciousness/` |

### Frontend - NEEDS REFACTORING

| Component | Status | Location |
|-----------|--------|----------|
| Navigation | 🔄 Needs refactor | `frontend/lib/platform/router/` |
| Chat Screen | ✅ Works | `frontend/lib/agno/chat/` |
| Shared Widgets | ✅ Ready | `frontend/lib/shared/widgets/` |
| Providers | ✅ Ready | `frontend/lib/platform/providers/` |
| 77 Screens | 🔄 Consolidate → 12 | Various modules |

---

## 8-Agent Assignments

### Agent Distribution

| Agent | Focus | Model | Issues | Points |
|-------|-------|-------|--------|--------|
| **Agent 1** | Universal Assets | Opus 4.5 | 3 | 13 |
| **Agent 2** | Context Engine | Opus 4.5 | 9 | 22 |
| **Agent 3** | Intelligence & Scheduling | Opus 4.5 | 9 | 24 |
| **Agent 4** | Agent Routing & Coordination | Opus 4.5 | 8 | 30 |
| **Agent 5** | Workflows & Automation | Opus 4.5 | 12 | 40 |
| **Agent 6** | Mobile Navigation | Opus 4.5 | 1 | 8 |
| **Agent 7** | Home & Execute UI | Opus 4.5 | 3 | 13 |
| **Agent 8** | Command, Capture & Support | Opus 4.5 | 11 | 51 |

### Dependency Graph

```
Week 1 (Foundation)
═══════════════════════════════════════════════════════════
Agent 1 (Assets) ──────────────────────────────────┐
                                                   │
Agent 6 (Navigation) ─────────────────────────────┐│
                                                  ││
                                                  ↓↓
Week 2 (Core Systems)
═══════════════════════════════════════════════════════════
Agent 2 (Context) ←── Agent 1                      │
         │                                         │
         ↓                                         │
Agent 3 (Scores/Scheduling) ←── Agent 2            │
         │                                         │
         ↓                                         ↓
Agent 4 (Routing) ←── Agent 2             Agent 7 (Home/Execute) ←── Agent 6
         │                                         │
         ↓                                         ↓
Agent 5 (Workflows) ←── Agent 4           Agent 8 (Command/Capture) ←── Agents 4,6,7

Week 3 (Integration)
═══════════════════════════════════════════════════════════
All Agents → Integration & Polish
```

---

## Agent 1: Universal Assets (Opus 4.5)

**Focus:** Create the unified Asset system that replaces separate task/habit/goal modules

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-467 | Universal Asset Schema | 5 | P0 |
| OMOS-468 | Asset Unified API | 5 | P0 |
| OMOS-460 | Planner Drag-Drop Backend | 3 | P0 |

**Total: 3 issues, 13 points**

### Key Files to Create/Modify

```
backend/lifeos/assets/
├── __init__.py
├── models.py          # Universal Asset schema
├── service.py         # AssetService with CRUD
├── api.py            # REST endpoints
├── compatibility.py  # Maps existing modules → assets
└── migrations/
    └── 001_create_assets.py

backend/lifeos/planner/service.py  # Add drag-drop methods
```

### Dependencies
- **None** - This is foundational, can start immediately

### Success Criteria
- [ ] Asset model with all types (task, habit, goal, event, meal, journal, project)
- [ ] CRUD API working: `GET/POST/PATCH/DELETE /api/v1/assets`
- [ ] Compatibility layer for existing task/habit/goal endpoints
- [ ] Planner drag-drop operations functional
- [ ] Migration applied without data loss

### Detailed Requirements

**OMOS-467: Universal Asset Schema**
```python
class AssetType(str, Enum):
    TASK = "task"
    HABIT = "habit"
    GOAL = "goal"
    EVENT = "event"
    MEAL = "meal"
    JOURNAL = "journal"
    PROJECT = "project"
    ROUTINE = "routine"
    NOTE = "note"
    REMINDER = "reminder"

class Asset(BaseModel):
    id: UUID
    asset_type: AssetType
    title: str
    description: Optional[str] = None

    # Classification
    pillar: Optional[Pillar] = None  # HP, LE, GE
    priority: Priority = Priority.NORMAL
    energy_level: Optional[EnergyLevel] = None

    # Timing
    do_date: Optional[datetime] = None
    due_date: Optional[datetime] = None
    duration_minutes: Optional[int] = None
    recurring: Optional[RecurrenceRule] = None

    # Status
    status: AssetStatus = AssetStatus.PENDING
    progress: float = 0.0
    streak: Optional[int] = None  # For habits

    # Relations
    parent_id: Optional[UUID] = None
    project_id: Optional[UUID] = None
    linked_assets: List[UUID] = []

    # Metadata
    tags: List[str] = []
    context: Dict[str, Any] = {}
    source: str = "manual"

    # Timestamps
    created_at: datetime
    updated_at: datetime
    completed_at: Optional[datetime] = None
```

**OMOS-468: API Endpoints**
```
GET    /api/v1/assets              # List with filters (type, pillar, status, date range)
POST   /api/v1/assets              # Create any asset type
GET    /api/v1/assets/{id}         # Get single asset
PATCH  /api/v1/assets/{id}         # Update asset
DELETE /api/v1/assets/{id}         # Delete asset
POST   /api/v1/assets/{id}/complete # Complete/check-off
GET    /api/v1/assets/today        # Today's assets
GET    /api/v1/assets/search       # Full-text search
```

---

## Agent 2: Context Engine (Opus 4.5)

**Focus:** Build the Context Engine that aggregates ALL user data for AI operations

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-416 | Context Engine Framework | 3 | P0 |
| OMOS-417 | State Aggregator | 2 | - |
| OMOS-418 | Items Aggregator | 2 | - |
| OMOS-419 | Scores Aggregator | 2 | - |
| OMOS-420 | Calendar Aggregator | 2 | - |
| OMOS-421 | History Aggregator | 2 | - |
| OMOS-451 | Assets Aggregator | 3 | P1 |
| OMOS-456 | Biometrics Aggregator (Oura) | 3 | P1 |
| OMOS-479 | Photo Capture Processing | 3 | P1 |

**Total: 9 issues, 22 points**

### Key Files to Create/Modify

```
backend/consciousness/context/
├── __init__.py
├── engine.py         # Main ContextEngine class
├── models.py         # UnifiedContext, MinimalContext
├── api.py           # Context API endpoints
└── aggregators/
    ├── __init__.py
    ├── state.py      # Location, energy, focus
    ├── items.py      # Tasks, habits, goals (from assets)
    ├── scores.py     # PI/LI/WI (pulls from score service)
    ├── calendar.py   # Calendar events
    ├── history.py    # Recent activity
    ├── assets.py     # Physical/digital assets (Home Assistant, finances)
    └── biometrics.py # Oura integration

backend/agno/capture/
└── photo_processor.py  # GPT-4 Vision integration
```

### Dependencies
- **Agent 1** - Needs Universal Asset schema for Items Aggregator

### Success Criteria
- [ ] `ContextEngine.build()` returns complete context in <500ms
- [ ] All 7 aggregators implemented and tested
- [ ] Oura data flowing into biometrics
- [ ] Photo capture processing with GPT-4 Vision
- [ ] Caching layer reduces repeated calls
- [ ] API endpoint: `GET /api/v1/context`

### Detailed Requirements

**OMOS-416: Context Engine Framework**
```python
class ContextEngine:
    def __init__(self, db: Database, oura: Optional[OuraClient] = None):
        self.db = db
        self.oura = oura
        self.cache = ContextCache()
        self.aggregators = {
            'state': StateAggregator(),
            'items': ItemsAggregator(),
            'scores': ScoresAggregator(),
            'calendar': CalendarAggregator(),
            'history': HistoryAggregator(),
            'assets': AssetsAggregator(),
            'biometrics': BiometricsAggregator(oura),
        }

    async def build(self, user_id: str) -> UnifiedContext:
        """Build full context for AI operations - target <500ms."""
        results = await asyncio.gather(*[
            agg.aggregate(user_id)
            for agg in self.aggregators.values()
        ])

        return UnifiedContext(
            user_id=user_id,
            timestamp=datetime.utcnow(),
            state=results[0],
            items=results[1],
            scores=results[2],
            calendar=results[3],
            history=results[4],
            assets=results[5],
            biometrics=results[6],
        )

    async def build_minimal(self, user_id: str) -> MinimalContext:
        """Build minimal context for quick operations."""
        # Only state + today's items
        ...
```

---

## Agent 3: Intelligence & Scheduling (Opus 4.5)

**Focus:** Build Intelligence Scores (PI/LI/WI) and intelligent scheduling (DO DATE)

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-422 | HPI Calculator | 3 | P0 |
| OMOS-423 | LEI Calculator | 2 | - |
| OMOS-424 | GEI Calculator | 2 | - |
| OMOS-425 | UI Calculator (Unified Index) | 2 | - |
| OMOS-426 | Score History & Trending | 2 | - |
| OMOS-431 | DO DATE Algorithm | 5 | P1 |
| OMOS-432 | Auto-Reschedule | 3 | - |
| OMOS-433 | Energy Matching | 3 | - |
| OMOS-434 | Path Balance Scheduling | 2 | - |

**Total: 9 issues, 24 points**

### Key Files to Create/Modify

```
backend/lifeos/scores/
├── __init__.py
├── hpi.py           # Holistic Performance Index
├── lei.py           # Legacy Evolution Index
├── gei.py           # Generational Endeavor Index
├── unified.py       # Unified Index (average)
├── history.py       # Score persistence & trending
└── api.py           # Score endpoints

backend/lifeos/scheduling/
├── __init__.py
├── do_date.py       # DO DATE algorithm
├── auto_reschedule.py
├── energy_matching.py
├── path_balance.py
└── api.py
```

### Dependencies
- **Agent 2** - Needs Context Engine for data (biometrics, items)

### Success Criteria
- [ ] All 4 intelligence scores calculating correctly
- [ ] Score history stored with 7/30/90 day trends
- [ ] DO DATE algorithm scheduling tasks optimally
- [ ] Energy matching using Oura readiness score
- [ ] Path balance tracking HP/LE/GE distribution
- [ ] API endpoints for all scores and scheduling

### Detailed Requirements

**HPI Components & Weights:**
- Sleep Quality: 20% (from Oura)
- Recovery Score: 20% (from Oura)
- Habit Completion: 25% (from habits)
- Task Completion: 20% (from tasks)
- Energy Levels: 15% (from Oura)

**DO DATE Algorithm Factors:**
- Task priority and deadline
- Estimated duration
- Available time slots
- Energy prediction for slot
- Path balance (HP/LE/GE)
- User preferences (morning/afternoon person)

---

## Agent 4: Agent Routing & Coordination (Opus 4.5)

**Focus:** Build @mention routing, /commands, Mother Agent coordination, and memory system

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-474 | @Mention Agent Routing Mobile | 3 | P1 |
| OMOS-475 | /Command Quick Actions Mobile | 3 | P1 |
| OMOS-453 | Mother Agent Coordination System | 5 | P1 |
| OMOS-454 | Specialist Agent Definitions | 3 | P1 |
| OMOS-430 | Smart Chat Engine Integration | 3 | P0 |
| OMOS-429 | Multi-Step Workflow from Chat | 5 | - |
| OMOS-452 | Risk Detection for Agent Actions | 3 | P0 |
| OMOS-457 | Memory Blocks with Semantic Search | 5 | P1 |

**Total: 8 issues, 30 points**

### Key Files to Create/Modify

```
backend/agno/routing/
├── __init__.py
├── parser.py        # Message parsing for @mentions, /commands
├── router.py        # Route to correct agent
├── commands.py      # Command registry and handlers
└── orchestrator.py  # Mother Agent coordination

backend/agno/coordination/
├── __init__.py
├── mother.py        # Mother Agent (Legacy) coordination
├── handoff.py       # Agent-to-agent handoffs
└── review.py        # Output review system

backend/agno/safety/
├── __init__.py
└── risk_detector.py # Risk level assessment

backend/agno/memory/
├── blocks.py        # Core memory blocks
└── archival.py      # Semantic search over memories
```

### Dependencies
- **Agent 2** - Needs Context Engine for context injection

### Success Criteria
- [ ] @mentions route to correct specialist agents
- [ ] /commands execute immediately with formatted response
- [ ] Mother Agent coordinates multi-agent tasks
- [ ] Risk detection blocks dangerous actions
- [ ] Memory blocks persist user knowledge
- [ ] Semantic search over archival memory

### Detailed Requirements

**Supported @Mentions:**
- `@legacy` - General AI (default)
- `@research` - Deep research
- `@code` - Code assistance
- `@plan` - Planning/scheduling
- `@health` - Health/HP queries
- `@finance` - Finance/GE queries
- `@write` - Content creation

**Supported /Commands:**
- `/brief` - Morning briefing
- `/tasks` - Today's tasks
- `/habits` - Today's habits
- `/scores` - Current scores
- `/schedule` - Today's schedule
- `/focus 30` - Start focus session
- `/add [text]` - Quick add task
- `/help` - List commands

**Risk Levels:**
- LOW - Read operations, auto-execute
- MEDIUM - Create/update, agent decides
- HIGH - Delete, send external, require approval
- CRITICAL - Credentials, finances, always block

---

## Agent 5: Workflows & Automation (Sonnet)

**Focus:** Build workflow triggers, automation, and operating rhythm

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-435 | HTTP Triggers | 3 | P1 |
| OMOS-436 | Schedule Triggers (Cron) | 3 | - |
| OMOS-437 | Smart Approvals (agent_decides) | 3 | - |
| OMOS-438 | Event Triggers | 3 | - |
| OMOS-439 | AI Triage Engine | 3 | P1 |
| OMOS-441 | Task Extraction | 3 | - |
| OMOS-443 | Morning Briefing Workflow | 3 | P2 |
| OMOS-444 | Evening Review Workflow | 3 | - |
| OMOS-445 | Weekly Planning Automation | 3 | - |
| OMOS-446 | Quarterly OKR Tracking | 3 | - |
| OMOS-458 | Cross-Domain Insights | 5 | P1 |
| OMOS-459 | Workflow Templates Library | 5 | P2 |

**Total: 12 issues, 40 points**

### Key Files to Create/Modify

```
backend/agno/workflows/triggers/
├── __init__.py
├── http.py          # Webhook trigger
├── schedule.py      # Cron-based trigger
└── events.py        # Event-based trigger

backend/agno/workflows/workflow/steps/
└── approval.py      # Enhanced with agent_decides mode

backend/lifeos/inbox/
└── intelligence.py  # AI triage engine

backend/lifeos/insights/
├── __init__.py
├── engine.py        # Insights engine
├── correlations.py  # Find correlations
├── anomalies.py     # Detect anomalies
└── predictions.py   # Generate predictions

backend/agno/workflows/templates/
├── __init__.py
├── library.py       # Template library
├── morning_startup.yaml
├── evening_shutdown.yaml
├── weekly_review.yaml
└── ... (10 templates)
```

### Dependencies
- **Agent 4** - Needs routing system for workflow triggers

### Success Criteria
- [ ] HTTP triggers working with webhook URLs
- [ ] Schedule triggers running on cron
- [ ] Smart approvals with agent_decides mode
- [ ] Event triggers firing on system events
- [ ] AI triage classifying inbox items
- [ ] Morning/evening/weekly workflows defined
- [ ] 10+ workflow templates ready
- [ ] Insights engine finding correlations

---

## Agent 6: Mobile Navigation (Sonnet)

**Focus:** Complete redesign of navigation for mobile-first experience

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-466 | Mobile Navigation Redesign | 8 | P0 |

**Total: 1 issue, 8 points** (Large foundational issue)

### Key Files to Create/Modify

```
frontend/lib/platform/router/
└── app_router.dart  # Consolidated routes (77 → 12)

frontend/lib/platform/providers/
└── navigation_provider.dart  # Updated nav state

frontend/lib/shared/widgets/
├── app_shell.dart   # New bottom nav with FAB
└── app_drawer.dart  # Simplified for secondary nav
```

### Dependencies
- **None** - Can start immediately

### Success Criteria
- [ ] New bottom nav: [Home] [Command] [Execute] [Profile]
- [ ] Center FAB for universal capture
- [ ] Routes consolidated from 77 to 12 primary
- [ ] Drawer simplified for secondary navigation
- [ ] Smooth transitions between screens
- [ ] Gesture navigation support

### Detailed Requirements

**New Bottom Nav Structure:**
```dart
class AppShell extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      body: child,
      drawer: const SimplifiedDrawer(),
      floatingActionButton: const CaptureButton(),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      bottomNavigationBar: BottomAppBar(
        shape: const CircularNotchedRectangle(),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            NavItem(icon: Icons.home, label: 'Home', route: '/home'),
            NavItem(icon: Icons.chat, label: 'Command', route: '/command'),
            const SizedBox(width: 48), // Space for FAB
            NavItem(icon: Icons.flash_on, label: 'Execute', route: '/execute'),
            NavItem(icon: Icons.person, label: 'Profile', route: '/profile'),
          ],
        ),
      ),
    );
  }
}
```

**Route Consolidation:**
```
PRIMARY (Bottom Nav)
├── /home         ← Dashboard, scores, agents
├── /command      ← Chat with Legacy
├── /execute      ← Today's tasks/habits
├── /profile      ← Settings, config
└── /capture      ← Universal input (modal)

SECONDARY (From Drawer)
├── /tasks        ← Full asset list
├── /calendar     ← Time-based view
├── /insights     ← Analytics
├── /goals        ← Three Paths progress
├── /agents       ← Agent management
├── /knowledge    ← Memory and docs
└── /settings     ← Deep configuration
```

---

## Agent 7: Home & Execute UI (Sonnet)

**Focus:** Build the Home dashboard and Execute view screens

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-470 | Home Dashboard Mobile | 5 | P0 |
| OMOS-471 | Execute View (Today) | 5 | P0 |
| OMOS-480 | Smart Suggestions Widget | 3 | P1 |

**Total: 3 issues, 13 points**

### Key Files to Create/Modify

```
frontend/lib/home/
├── screens/
│   └── home_screen.dart
└── widgets/
    ├── score_card.dart
    ├── unified_score.dart
    ├── active_agents_card.dart
    ├── focus_card.dart
    └── quick_stats.dart

frontend/lib/execute/
├── screens/
│   └── execute_screen.dart
└── widgets/
    ├── focus_block.dart
    ├── task_list.dart
    ├── habit_row.dart
    ├── schedule_timeline.dart
    └── suggestion_card.dart
```

### Dependencies
- **Agent 6** - Needs navigation structure

### Success Criteria
- [ ] Home shows unified score + PI/LI/WI
- [ ] Home shows today's focus and active agents
- [ ] Execute shows today's tasks, habits, schedule
- [ ] Smart suggestions appearing contextually
- [ ] Pull-to-refresh working
- [ ] Skeleton loading states

### Detailed Requirements

**Home Dashboard Layout:**
```
┌─────────────────────────────────────┐
│ OneMind                    🔔 ⚡    │
├─────────────────────────────────────┤
│         UNIFIED SCORE: 84           │
│    PI:87  •  LI:82  •  WI:83       │
├─────────────────────────────────────┤
│ Today's Focus                       │
│ 🎯 Complete MVP launch prep    78%  │
├─────────────────────────────────────┤
│ Active Agents (3)                   │
│ • Planning: Optimizing schedule     │
│ • Research: Gathering data          │
├─────────────────────────────────────┤
│ Quick Stats                         │
│ Tasks: 5 • Habits: 7/10 • Streak: 23│
└─────────────────────────────────────┘
```

**Execute View Layout:**
```
┌─────────────────────────────────────┐
│ ⚡ Execute            Wed Jan 29    │
├─────────────────────────────────────┤
│ Focus Block: Deep Work 🧠  9-11am   │
│ ☐ Code review PR #123      45min   │
│ ☐ Write proposal           30min   │
├─────────────────────────────────────┤
│ Habits Today                 5/7    │
│ [✓Med] [✓Exer] [☐Journ] [✓Read]   │
├─────────────────────────────────────┤
│ Up Next                             │
│ • 11:00 Team standup (30min)       │
│ • 12:00 Lunch break                │
├─────────────────────────────────────┤
│ Suggested by Legacy:                │
│ "Move budget review to Friday"      │
│ [Accept] [Dismiss]                  │
└─────────────────────────────────────┘
```

---

## Agent 8: Command, Capture & Support (Sonnet)

**Focus:** Build Command chat, Universal Capture, and all support screens

### Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| OMOS-469 | Universal Capture Component | 5 | P0 |
| OMOS-472 | Command Chat Mobile | 5 | P0 |
| OMOS-478 | Voice Capture Integration | 3 | P1 |
| OMOS-473 | Agent Visibility Layer | 3 | P1 |
| OMOS-476 | Insights Mobile View | 3 | P1 |
| OMOS-477 | Profile & Settings Mobile | 3 | P1 |
| OMOS-481 | Agent Status Dashboard | 3 | P2 |
| OMOS-482 | Memory UI Component | 3 | P2 |
| OMOS-483 | Knowledge Mobile UI | 3 | P2 |
| OMOS-484 | Mobile Approval UI | 3 | P2 |
| OMOS-485 | Integration & Polish Pass | 5 | P1 |

**Total: 11 issues, 39 points**

### Key Files to Create/Modify

```
frontend/lib/command/
├── screens/
│   └── command_screen.dart
└── widgets/
    ├── chat_messages.dart
    ├── chat_input.dart
    ├── mention_autocomplete.dart
    ├── command_autocomplete.dart
    ├── agent_indicator.dart
    └── thinking_loader.dart

frontend/lib/capture/
├── screens/
│   └── capture_modal.dart
└── widgets/
    ├── capture_input.dart
    ├── voice_capture.dart
    ├── photo_capture.dart
    ├── quick_types.dart
    └── classification_preview.dart

frontend/lib/insights/
├── screens/
│   └── insights_screen.dart
└── widgets/
    ├── score_trends.dart
    ├── correlation_card.dart
    └── path_balance.dart

frontend/lib/profile/
├── screens/
│   └── profile_screen.dart
└── widgets/
    ├── profile_header.dart
    ├── settings_section.dart
    └── integration_status.dart

frontend/lib/agents/
├── screens/
│   └── agents_screen.dart
└── widgets/
    ├── agent_list.dart
    ├── agent_status.dart
    └── approval_card.dart

frontend/lib/knowledge/
├── screens/
│   └── knowledge_screen.dart
└── widgets/
    ├── document_list.dart
    ├── memory_list.dart
    └── search_bar.dart
```

### Dependencies
- **Agent 4** - Needs routing for @mentions and /commands
- **Agent 6** - Needs navigation structure
- **Agent 7** - Needs home/execute for integration

### Success Criteria
- [ ] Command chat with @mention and /command autocomplete
- [ ] Universal capture modal (text, voice, photo)
- [ ] Voice transcription via Whisper
- [ ] Agent visibility throughout app
- [ ] Insights screen with trends
- [ ] Profile with settings
- [ ] Agent status dashboard
- [ ] Memory and knowledge browsing
- [ ] Approval UI for HITL
- [ ] Final polish with animations and accessibility

---

## Communication Protocol

### Issue Status Updates

When agents complete work, they must update Linear issues to **In Review** state:

```python
# Update issue to In Review
mutation {
  issueUpdate(id: "<issue_uuid>", input: { stateId: "<in_review_state_id>" }) {
    success
  }
}
```

### Handoff Format

When agent completes task and another depends on it:

```markdown
## Handoff: Agent 1 → Agent 2

### Completed
- Universal Asset schema in `backend/lifeos/assets/models.py`
- CRUD API in `backend/lifeos/assets/api.py`
- Migration applied

### For Agent 2
- Import Asset model from `backend.lifeos.assets.models`
- Use AssetService for queries
- API base: `/api/v1/assets`

### Notes
- Backward compatibility layer ready
- Existing task/habit endpoints still work
```

---

## Key Patterns to Follow

### Backend
- **Service Layer Pattern** - Business logic in `service.py`
- **Pydantic Models** - Use for all schemas
- **Async/Await** - All database and external calls
- **FastAPI Routers** - REST endpoints in `api.py`

### Frontend
- **State Management** - Riverpod (StateNotifierProvider, FutureProvider)
- **API Calls** - Use existing ApiClient in `platform/services/api_client.dart`
- **Styling** - Use TacticalColors and existing theme
- **Navigation** - Use GoRouter `context.go()` and `context.push()`
- **Models** - Use freezed for immutable models

### Code Style
- Python: Black formatter, type hints everywhere
- Dart: Flutter linter, null safety

---

## Success Metrics

### Week 1 (Foundation)
- [ ] Agent 1: Asset API returning data
- [ ] Agent 6: New nav working on device
- [ ] Agent 2: Context Engine building

### Week 2 (Core Systems)
- [ ] Agent 3: Scores calculating correctly
- [ ] Agent 4: @mentions routing correctly
- [ ] Agent 7: Home/Execute screens functional

### Week 3 (Integration)
- [ ] Agent 5: Workflows running
- [ ] Agent 8: All screens connected
- [ ] All: No critical bugs, ready for production

---

## Linear Project Info

- **Project ID:** `5d760a6b-9bbb-4eb0-bb32-d15884f9997b`
- **Team ID:** `5dcfd649-4f48-47be-8419-5ef5550e5b76`
- **API Key:** `lin_api_skPs6T60yTr5eW24diYGNtWZDhGQGfCOWLAxhOEA`

---

**Sprint Start: 2026-01-29 | 56 Issues | 8 Agents | Mobile-First Evolution**
