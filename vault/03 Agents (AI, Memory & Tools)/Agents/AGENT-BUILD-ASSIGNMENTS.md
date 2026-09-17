# OneMind Mobile-First: 8-Agent Build Assignments

> **Parallel Agent Swarm for Rapid Development**
> **Last Updated:** 2026-01-29
> **Goal:** Complete mobile-first evolution in 3 weeks

---

## Overview

This document defines the work assignments for 8 parallel agents using Opus 4.5 (for complex tasks) and Sonnet (for UI/frontend tasks). Each agent has specific responsibilities, Linear issues, and dependencies.

---

## Agent Summary

| Agent | Focus | Model | Linear Issues | Dependencies |
|-------|-------|-------|---------------|--------------|
| **Agent 1** | Backend Architecture | Opus 4.5 | OMOS-467, 468 | None (foundational) |
| **Agent 2** | Context Engine | Opus 4.5 | OMOS-301-306, 479 | Agent 1 |
| **Agent 3** | Agent Routing | Opus 4.5 | OMOS-320, 321, 474, 475 | Agent 2 |
| **Agent 4** | Mobile Navigation | Sonnet | OMOS-466 | None (can start immediately) |
| **Agent 5** | Home & Execute | Sonnet | OMOS-470, 471, 480 | Agent 4 |
| **Agent 6** | Command & Capture | Sonnet | OMOS-469, 472, 478 | Agents 3, 4 |
| **Agent 7** | Insights & Profile | Sonnet | OMOS-473, 476, 477, 481-484 | Agent 4 |
| **Agent 8** | Integration & Polish | Sonnet | OMOS-485 | Agents 1-7 |

---

## Week 1: Foundation Phase

### Day 1-2

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 1** | Universal Asset schema | `backend/lifeos/assets/models.py`, migrations |
| **Agent 2** | Context Engine framework | `backend/agno/context/engine.py` |
| **Agent 3** | Message parser design | `backend/agno/routing/parser.py` |
| **Agent 4** | New bottom nav | `frontend/lib/shared/widgets/app_shell.dart` |

### Day 3-4

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 1** | Asset CRUD API | `backend/lifeos/assets/api.py`, `service.py` |
| **Agent 2** | Aggregators implementation | `backend/agno/context/aggregators/*.py` |
| **Agent 3** | @mention routing | `backend/agno/routing/router.py` |
| **Agent 4** | Route consolidation | `frontend/lib/platform/router/app_router.dart` |

### Day 5

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 1** | Backward compatibility | `backend/lifeos/assets/compatibility.py` |
| **Agent 2** | Score calculations | `backend/agno/context/aggregators/scores.py` |
| **Agent 3** | /command system | `backend/agno/routing/commands.py` |
| **Agent 4** | Navigation provider | `frontend/lib/platform/providers/navigation_provider.dart` |

---

## Week 2: UI Build Phase

### Day 1-2

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 5** | Home dashboard | `frontend/lib/home/screens/home_screen.dart` |
| **Agent 6** | Chat UI | `frontend/lib/command/screens/command_screen.dart` |
| **Agent 7** | Insights screen | `frontend/lib/insights/screens/insights_screen.dart` |
| **Agent 8** | Integration setup | Wire providers, services |

### Day 3-4

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 5** | Execute view | `frontend/lib/execute/screens/execute_screen.dart` |
| **Agent 6** | Capture modal | `frontend/lib/capture/screens/capture_modal.dart` |
| **Agent 7** | Profile screen | `frontend/lib/profile/screens/profile_screen.dart` |
| **Agent 8** | Wire home/execute | Integration testing |

### Day 5

| Agent | Task | Files to Create/Modify |
|-------|------|------------------------|
| **Agent 5** | Suggestions widget | `frontend/lib/shared/widgets/suggestion_card.dart` |
| **Agent 6** | Voice capture | `frontend/lib/capture/widgets/voice_capture.dart` |
| **Agent 7** | Agent dashboard | `frontend/lib/agents/screens/agents_screen.dart` |
| **Agent 8** | Wire command/capture | Integration testing |

---

## Week 3: Integration Phase

### Day 1-2

| All Agents | Task |
|------------|------|
| Full end-to-end testing |
| API integration verification |
| Data flow testing |

### Day 3-4

| All Agents | Task |
|------------|------|
| Bug fixes and refinements |
| Performance optimization |
| Memory leak checks |

### Day 5

| All Agents | Task |
|------------|------|
| Animation polish |
| Final testing |
| Documentation |

---

## Agent 1: Backend Architecture (Opus 4.5)

### Responsibility
Create the foundational Universal Asset system that unifies all trackable items.

### Linear Issues
- **OMOS-467**: Universal Asset Schema (5pts)
- **OMOS-468**: Asset Unified API (5pts)

### Files to Create

```
backend/lifeos/assets/
├── __init__.py
├── models.py          # Asset, AssetType, AssetStatus models
├── service.py         # AssetService with CRUD operations
├── api.py            # REST endpoints
├── compatibility.py  # Mapping from existing modules
└── migrations/
    └── 001_create_assets.py
```

### Key Implementation

```python
# models.py
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
    pillar: Optional[Pillar] = None
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
    streak: Optional[int] = None

    # Relations
    parent_id: Optional[UUID] = None
    project_id: Optional[UUID] = None
    linked_assets: List[UUID] = []

    # Metadata
    tags: List[str] = []
    context: Dict[str, Any] = {}
    source: str = "manual"
```

### Success Criteria
- [ ] All existing task/habit/goal data migrated to assets
- [ ] API returns <100ms for list operations
- [ ] Backward compatible with existing endpoints
- [ ] OpenAPI spec generated

### Dependencies
None - this is foundational

---

## Agent 2: Context Engine (Opus 4.5)

### Responsibility
Build the Context Engine that aggregates ALL user data for AI operations.

### Linear Issues
- **OMOS-301**: Context Engine Framework (3pts)
- **OMOS-302**: State Aggregator (2pts)
- **OMOS-303**: Items Aggregator (2pts)
- **OMOS-304**: Scores Aggregator (2pts)
- **OMOS-305**: Calendar Aggregator (2pts)
- **OMOS-306**: History Aggregator (2pts)
- **OMOS-479**: Photo Capture Processing (3pts)

### Files to Create

```
backend/agno/context/
├── __init__.py
├── engine.py         # Main ContextEngine class
├── models.py         # UnifiedContext, MinimalContext models
├── api.py           # Context API endpoints
└── aggregators/
    ├── __init__.py
    ├── state.py      # Current state (energy, focus, location)
    ├── assets.py     # Tasks, habits, goals
    ├── scores.py     # PI/LI/WI calculations
    ├── calendar.py   # Calendar events
    ├── biometrics.py # Oura integration
    └── history.py    # Recent activity
```

### Key Implementation

```python
# engine.py
class ContextEngine:
    def __init__(self, db: Database, oura: Optional[OuraClient] = None):
        self.db = db
        self.oura = oura
        self.cache = ContextCache()

    async def build(self, user_id: str) -> UnifiedContext:
        """Build full context for AI operations."""
        # Parallel aggregation for speed
        state, assets, scores, calendar, history = await asyncio.gather(
            self.aggregate_state(user_id),
            self.aggregate_assets(user_id),
            self.aggregate_scores(user_id),
            self.aggregate_calendar(user_id),
            self.aggregate_history(user_id),
        )

        return UnifiedContext(
            user_id=user_id,
            timestamp=datetime.utcnow(),
            state=state,
            assets=assets,
            scores=scores,
            calendar=calendar,
            history=history,
        )

    async def build_minimal(self, user_id: str) -> MinimalContext:
        """Build minimal context for quick operations."""
        # Only essential data
        ...
```

### Success Criteria
- [ ] Full context builds in <500ms
- [ ] Oura data integrated when available
- [ ] Caching reduces repeated calls
- [ ] Context API documented

### Dependencies
- Agent 1: Universal Asset schema

---

## Agent 3: Agent Routing (Opus 4.5)

### Responsibility
Build the @mention and /command routing system for chat.

### Linear Issues
- **OMOS-320**: @Agent Routing (3pts)
- **OMOS-321**: /Commands System (3pts)
- **OMOS-474**: @Mention Agent Routing Mobile (3pts)
- **OMOS-475**: /Command Quick Actions Mobile (3pts)

### Files to Create

```
backend/agno/routing/
├── __init__.py
├── parser.py        # Message parsing for @mentions, /commands
├── router.py        # Route to correct agent
├── commands.py      # Command registry and handlers
└── orchestrator.py  # Mother Agent coordination
```

### Key Implementation

```python
# parser.py
class MessageParser:
    def parse(self, message: str) -> ParsedMessage:
        """Parse message for @mentions and /commands."""
        mentions = self._extract_mentions(message)
        command = self._extract_command(message)
        clean_text = self._clean_text(message)

        return ParsedMessage(
            raw=message,
            mentions=mentions,
            command=command,
            text=clean_text,
        )

    def _extract_mentions(self, message: str) -> List[AgentMention]:
        pattern = r'@(\w+)'
        matches = re.findall(pattern, message)
        return [self._resolve_agent(m) for m in matches]

# commands.py
class CommandRegistry:
    commands = {
        "brief": BriefCommand,
        "tasks": TasksCommand,
        "habits": HabitsCommand,
        "scores": ScoresCommand,
        "schedule": ScheduleCommand,
        "focus": FocusCommand,
        "add": AddCommand,
        "help": HelpCommand,
    }

    async def execute(self, command: str, args: str, context: Context) -> str:
        handler = self.commands.get(command)
        if handler:
            return await handler().execute(args, context)
        return f"Unknown command: {command}. Type /help for available commands."
```

### Success Criteria
- [ ] All @mentions route correctly
- [ ] All /commands execute
- [ ] Unknown mentions fall back to Legacy
- [ ] Help command lists all options

### Dependencies
- Agent 2: Context Engine (for command execution context)

---

## Agent 4: Mobile Navigation (Sonnet)

### Responsibility
Redesign navigation for mobile-first experience.

### Linear Issues
- **OMOS-466**: Mobile Navigation Redesign (8pts)

### Files to Modify

```
frontend/lib/
├── platform/router/app_router.dart         # Consolidated routes
├── platform/providers/navigation_provider.dart
├── shared/widgets/app_shell.dart           # New bottom nav
└── shared/widgets/app_drawer.dart          # Simplified drawer
```

### Key Implementation

```dart
// app_shell.dart - New bottom nav
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

class CaptureButton extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      onPressed: () => CaptureModal.show(context),
      backgroundColor: TacticalColors.primary,
      child: const Icon(Icons.add, size: 28),
    );
  }
}
```

### Success Criteria
- [ ] 4 bottom nav items + center FAB
- [ ] Routes consolidated from 77 to 12 primary
- [ ] Drawer simplified for secondary nav
- [ ] Smooth transitions

### Dependencies
None - can start immediately

---

## Agent 5: Home & Execute Screens (Sonnet)

### Responsibility
Build the Home dashboard and Execute view.

### Linear Issues
- **OMOS-470**: Home Dashboard Mobile (5pts)
- **OMOS-471**: Execute View (Today) (5pts)
- **OMOS-480**: Smart Suggestions Widget (3pts)

### Files to Create

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

### Key Implementation

```dart
// home_screen.dart
class HomeScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scores = ref.watch(scoresProvider);
    final agents = ref.watch(activeAgentsProvider);

    return BaseScreen(
      title: 'OneMind',
      actions: [NotificationBell(), AwarenessIndicator()],
      body: RefreshIndicator(
        onRefresh: () => ref.refresh(scoresProvider.future),
        child: ListView(
          children: [
            UnifiedScoreCard(score: scores.unified),
            ScoreRow(pi: scores.pi, li: scores.li, wi: scores.wi),
            FocusCard(focus: ref.watch(todayFocusProvider)),
            ActiveAgentsCard(agents: agents),
            QuickStats(stats: ref.watch(quickStatsProvider)),
          ],
        ),
      ),
    );
  }
}

// execute_screen.dart
class ExecuteScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final today = ref.watch(todayAssetsProvider);

    return BaseScreen(
      title: 'Execute',
      subtitle: DateFormat('EEE MMM d').format(DateTime.now()),
      body: ListView(
        children: [
          if (today.focusBlock != null)
            FocusBlockCard(block: today.focusBlock!),
          HabitQuickToggle(habits: today.habits),
          TaskList(tasks: today.tasks),
          ScheduleTimeline(events: today.events),
          if (today.suggestions.isNotEmpty)
            SuggestionCard(suggestion: today.suggestions.first),
        ],
      ),
    );
  }
}
```

### Success Criteria
- [ ] Home shows all key widgets
- [ ] Execute shows today's actionable items
- [ ] Pull-to-refresh working
- [ ] Skeleton loading states

### Dependencies
- Agent 4: Mobile Navigation (for routing)

---

## Agent 6: Command & Capture (Sonnet)

### Responsibility
Build the Command chat screen and Universal Capture modal.

### Linear Issues
- **OMOS-469**: Universal Capture Component (5pts)
- **OMOS-472**: Command Chat Mobile (5pts)
- **OMOS-478**: Voice Capture Integration (3pts)

### Files to Create

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
```

### Key Implementation

```dart
// command_screen.dart
class CommandScreen extends ConsumerStatefulWidget {
  @override
  ConsumerState<CommandScreen> createState() => _CommandScreenState();
}

class _CommandScreenState extends ConsumerState<CommandScreen> {
  final _inputController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    final messages = ref.watch(chatMessagesProvider);
    final isThinking = ref.watch(isAgentThinkingProvider);

    return BaseScreen(
      title: 'Command',
      actions: [ContextIndicator()],
      body: Column(
        children: [
          Expanded(
            child: ChatMessages(messages: messages),
          ),
          if (isThinking) ThinkingLoader(),
          QuickActionChips(),
          ChatInput(
            controller: _inputController,
            onSend: _sendMessage,
            onVoice: _startVoice,
          ),
        ],
      ),
    );
  }
}

// capture_modal.dart
class CaptureModal extends ConsumerStatefulWidget {
  static Future<void> show(BuildContext context) {
    return showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (_) => const CaptureModal(),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return DraggableScrollableSheet(
      initialChildSize: 0.6,
      maxChildSize: 0.9,
      builder: (_, scrollController) {
        return Container(
          decoration: TacticalDecoration.modal,
          child: Column(
            children: [
              CaptureInput(onSubmit: _processCapture),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  VoiceButton(onCapture: _processVoice),
                  PhotoButton(onCapture: _processPhoto),
                  AttachButton(onCapture: _processFile),
                ],
              ),
              QuickTypeButtons(onSelect: _setType),
              if (_classification != null)
                ClassificationPreview(classification: _classification!),
              RecentCaptures(),
            ],
          ),
        );
      },
    );
  }
}
```

### Success Criteria
- [ ] Chat fully functional with streaming
- [ ] @mention and /command autocomplete
- [ ] Voice capture with transcription
- [ ] Photo capture with analysis
- [ ] Classification preview before save

### Dependencies
- Agent 3: Agent Routing (for @mention and /command processing)
- Agent 4: Mobile Navigation (for modal integration)

---

## Agent 7: Insights & Profile (Sonnet)

### Responsibility
Build analytics views and profile/settings screens.

### Linear Issues
- **OMOS-473**: Agent Visibility Layer (3pts)
- **OMOS-476**: Insights Mobile View (3pts)
- **OMOS-477**: Profile & Settings Mobile (3pts)
- **OMOS-481**: Agent Status Dashboard (3pts)
- **OMOS-482**: Memory UI Component (3pts)
- **OMOS-483**: Knowledge Mobile UI (3pts)
- **OMOS-484**: Mobile Approval UI (3pts)

### Files to Create

```
frontend/lib/insights/
├── screens/
│   └── insights_screen.dart
└── widgets/
    ├── score_trends.dart
    ├── correlation_card.dart
    ├── prediction_card.dart
    └── path_balance.dart

frontend/lib/profile/
├── screens/
│   └── profile_screen.dart
└── widgets/
    ├── profile_header.dart
    ├── settings_section.dart
    ├── agent_config.dart
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

### Success Criteria
- [ ] Insights shows trends and correlations
- [ ] Profile shows user info and settings
- [ ] Agents screen shows all agent status
- [ ] Knowledge screen allows browsing/search

### Dependencies
- Agent 4: Mobile Navigation (for routing)

---

## Agent 8: Integration & Polish (Sonnet)

### Responsibility
Wire all components together with animations and polish.

### Linear Issues
- **OMOS-485**: Integration & Polish Pass (5pts)

### Files to Modify
All screens and components for:
- Consistent animations
- Loading states
- Error handling
- Accessibility
- Performance

### Key Tasks

```dart
// Animations
class PageTransitions {
  static Widget fade(Widget child) => FadeTransition(
    opacity: _animation,
    child: child,
  );

  static Widget slideUp(Widget child) => SlideTransition(
    position: _slideAnimation,
    child: child,
  );
}

// Loading States
class SkeletonLoading {
  static Widget scoreCard() => ShimmerContainer(
    child: Container(
      height: 100,
      decoration: BoxDecoration(
        color: TacticalColors.surface,
        borderRadius: BorderRadius.circular(12),
      ),
    ),
  );
}

// Error Boundaries
class ErrorBoundary extends StatelessWidget {
  final Widget child;
  final Widget Function(Object error)? onError;

  @override
  Widget build(BuildContext context) {
    return ErrorWidget.builder = (details) {
      return onError?.call(details.exception) ??
        ErrorCard(error: details.exception.toString());
    };
  }
}
```

### Success Criteria
- [ ] 60fps animations throughout
- [ ] Consistent skeleton loading
- [ ] All errors handled gracefully
- [ ] VoiceOver/TalkBack accessible
- [ ] No memory leaks
- [ ] Offline mode functional

### Dependencies
- Agents 1-7 must complete first

---

## Dependency Graph

```
Week 1 (Foundation)
==================
Agent 1 ──────────────────────────────────┐
(Backend)                                  │
                                          ↓
Agent 2 ─────────────────────────────→ Agent 3
(Context)                               (Routing)
                                          │
Agent 4 ──────────────────────────────────┼───────────────────┐
(Navigation)                              │                   │
                                          │                   │
Week 2 (UI Build)                         │                   │
=================                         ↓                   ↓
Agent 5 ←─────────────────────────── Agent 4            Agent 6
(Home/Execute)                                         (Command/Capture)
                                                            │
Agent 7 ←─────────────────────────── Agent 4               │
(Insights/Profile)                                          │
                                                            │
Week 3 (Integration)                                        │
====================                                        │
Agent 8 ←─────────────────────────── All Agents ←──────────┘
(Polish)
```

---

## Communication Protocol

### Daily Sync Points
- **Morning (9am)**: Status check, blockers
- **Afternoon (3pm)**: Progress review, coordination
- **Evening (6pm)**: End-of-day summary

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

## Success Metrics

### Week 1
- [ ] Agent 1: Asset API returning data
- [ ] Agent 2: Context Engine building in <500ms
- [ ] Agent 3: @mentions routing correctly
- [ ] Agent 4: New nav working on device

### Week 2
- [ ] Agent 5: Home/Execute screens functional
- [ ] Agent 6: Chat/Capture working end-to-end
- [ ] Agent 7: All secondary screens built
- [ ] Agent 8: First integration pass complete

### Week 3
- [ ] All screens connected
- [ ] No critical bugs
- [ ] Performance targets met
- [ ] Ready for production

---

## Getting Started

### For Each Agent

1. **Read this document** fully
2. **Review assigned Linear issues** for details
3. **Check dependencies** - wait if blocked
4. **Create files** in specified locations
5. **Follow patterns** from existing code
6. **Test locally** before marking complete
7. **Post handoff** when done

### Key Patterns to Follow

- **State Management**: Riverpod (StateNotifierProvider, FutureProvider)
- **API Calls**: Use existing ApiClient in `platform/services/api_client.dart`
- **Styling**: Use TacticalColors and existing theme
- **Navigation**: Use GoRouter context.go() and context.push()
- **Models**: Use freezed for immutable models

---

**Ready to deploy the swarm? Each agent should start with their Day 1 tasks.**
