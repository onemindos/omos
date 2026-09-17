# Agent 7: Home & Execute UI (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Flutter Home Dashboard & Execute Screen
> **Issues:** 3 | **Points:** 13
> **Dependencies:** Agent 6 (Navigation)

---

## Your Mission

You are Agent 7. Your job is building the two most important mobile screens:

1. **Home Screen** - Dashboard with scores, today's focus, active agents
2. **Execute Screen** - Today's actionable items (tasks, habits, schedule)

These are the screens users see most. Make them beautiful and functional.

---

## Your Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-470** | Home Dashboard Mobile | 5 | P0-Critical |
| **OMOS-471** | Execute View (Today) | 5 | P0-Critical |
| **OMOS-480** | Smart Suggestions Widget | 3 | P1-High |

---

## Dependencies

**Wait for Agent 6** to complete:
- Navigation structure with routes
- AppShell with new bottom nav

Once Agent 6 posts handoff, you can begin.

---

## Files to Create

```
frontend/lib/home/
├── screens/
│   └── home_screen.dart
└── widgets/
    ├── unified_score_card.dart
    ├── score_row.dart
    ├── focus_card.dart
    ├── active_agents_card.dart
    └── quick_stats.dart

frontend/lib/execute/
├── screens/
│   └── execute_screen.dart
└── widgets/
    ├── focus_block_card.dart
    ├── task_item.dart
    ├── habit_toggle_row.dart
    ├── schedule_timeline.dart
    └── suggestion_card.dart
```

---

## OMOS-470: Home Dashboard

### Design

```
┌─────────────────────────────────────┐
│ OneMind                    🔔 ⚡    │
├─────────────────────────────────────┤
│         UNIFIED SCORE: 84           │
│         ████████████░░░░            │
│    PI:87  •  LI:82  •  WI:83       │
├─────────────────────────────────────┤
│ Today's Focus                       │
│ 🎯 Complete MVP launch prep    78%  │
│ ━━━━━━━━━━━━━━━━░░░░               │
├─────────────────────────────────────┤
│ Active Agents (3)                   │
│ 🤖 Planning: Optimizing schedule    │
│ 🔬 Research: Gathering data         │
│ ❤️ Health: Analyzing patterns       │
│                     [View All →]    │
├─────────────────────────────────────┤
│ Quick Stats                         │
│ Tasks: 5  │  Habits: 7/10  │  🔥 23 │
└─────────────────────────────────────┘
```

### Implementation

```dart
// frontend/lib/home/screens/home_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../widgets/unified_score_card.dart';
import '../widgets/score_row.dart';
import '../widgets/focus_card.dart';
import '../widgets/active_agents_card.dart';
import '../widgets/quick_stats.dart';

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scores = ref.watch(scoresProvider);
    final focus = ref.watch(todayFocusProvider);
    final agents = ref.watch(activeAgentsProvider);
    final stats = ref.watch(quickStatsProvider);

    return Scaffold(
      backgroundColor: const Color(0xFF0A0A0A),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: const Text(
          'OneMind',
          style: TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          IconButton(
            icon: const Badge(
              label: Text('3'),
              child: Icon(Icons.notifications_outlined, color: Colors.white),
            ),
            onPressed: () => context.push('/notifications'),
          ),
          IconButton(
            icon: const Icon(Icons.bolt, color: Color(0xFFFF0000)),
            onPressed: () => _showAwarenessSheet(context),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(scoresProvider);
          ref.invalidate(activeAgentsProvider);
        },
        color: const Color(0xFFFF0000),
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Unified Score
            scores.when(
              data: (data) => UnifiedScoreCard(
                unified: data.unified,
                pi: data.pi,
                li: data.li,
                wi: data.wi,
              ),
              loading: () => const ScoreCardSkeleton(),
              error: (e, _) => ErrorCard(message: e.toString()),
            ),

            const SizedBox(height: 16),

            // Today's Focus
            focus.when(
              data: (data) => FocusCard(focus: data),
              loading: () => const FocusCardSkeleton(),
              error: (e, _) => const SizedBox.shrink(),
            ),

            const SizedBox(height: 16),

            // Active Agents
            agents.when(
              data: (data) => ActiveAgentsCard(agents: data),
              loading: () => const AgentsCardSkeleton(),
              error: (e, _) => const SizedBox.shrink(),
            ),

            const SizedBox(height: 16),

            // Quick Stats
            stats.when(
              data: (data) => QuickStats(stats: data),
              loading: () => const QuickStatsSkeleton(),
              error: (e, _) => const SizedBox.shrink(),
            ),
          ],
        ),
      ),
    );
  }
}
```

### Score Card Widget

```dart
// frontend/lib/home/widgets/unified_score_card.dart
class UnifiedScoreCard extends StatelessWidget {
  final double unified;
  final double pi;
  final double li;
  final double wi;

  const UnifiedScoreCard({
    required this.unified,
    required this.pi,
    required this.li,
    required this.wi,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A1A),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFF333333)),
      ),
      child: Column(
        children: [
          // Unified Score
          Text(
            'UNIFIED SCORE',
            style: TextStyle(
              color: Colors.grey[400],
              fontSize: 12,
              letterSpacing: 1.5,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            unified.toStringAsFixed(0),
            style: const TextStyle(
              color: Colors.white,
              fontSize: 48,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          // Progress bar
          ClipRRect(
            borderRadius: BorderRadius.circular(4),
            child: LinearProgressIndicator(
              value: unified / 100,
              backgroundColor: const Color(0xFF333333),
              valueColor: AlwaysStoppedAnimation(_getScoreColor(unified)),
              minHeight: 8,
            ),
          ),
          const SizedBox(height: 16),
          // Individual scores
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _ScoreChip(label: 'PI', value: pi, color: Colors.green),
              _ScoreChip(label: 'LI', value: li, color: Colors.blue),
              _ScoreChip(label: 'WI', value: wi, color: Colors.orange),
            ],
          ),
        ],
      ),
    );
  }

  Color _getScoreColor(double score) {
    if (score >= 80) return Colors.green;
    if (score >= 60) return Colors.orange;
    return Colors.red;
  }
}

class _ScoreChip extends StatelessWidget {
  final String label;
  final double value;
  final Color color;

  const _ScoreChip({
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(
          label,
          style: TextStyle(color: Colors.grey[400], fontSize: 12),
        ),
        const SizedBox(height: 4),
        Row(
          children: [
            Container(
              width: 8,
              height: 8,
              decoration: BoxDecoration(
                color: color,
                shape: BoxShape.circle,
              ),
            ),
            const SizedBox(width: 4),
            Text(
              value.toStringAsFixed(0),
              style: const TextStyle(
                color: Colors.white,
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
```

### Active Agents Card

```dart
// frontend/lib/home/widgets/active_agents_card.dart
class ActiveAgentsCard extends StatelessWidget {
  final List<AgentStatus> agents;

  const ActiveAgentsCard({required this.agents, super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A1A),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFF333333)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Active Agents (${agents.length})',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
              TextButton(
                onPressed: () => context.push('/agents'),
                child: const Text(
                  'View All',
                  style: TextStyle(color: Color(0xFFFF0000)),
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          ...agents.take(3).map((agent) => _AgentRow(agent: agent)),
          if (agents.isEmpty)
            Center(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Text(
                  'No agents currently active',
                  style: TextStyle(color: Colors.grey[500]),
                ),
              ),
            ),
        ],
      ),
    );
  }
}

class _AgentRow extends StatelessWidget {
  final AgentStatus agent;

  const _AgentRow({required this.agent});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(
            width: 32,
            height: 32,
            decoration: BoxDecoration(
              color: agent.color.withOpacity(0.2),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(agent.icon, color: agent.color, size: 18),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  agent.name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Text(
                  agent.currentTask,
                  style: TextStyle(
                    color: Colors.grey[400],
                    fontSize: 12,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          _PulsingIndicator(),
        ],
      ),
    );
  }
}
```

---

## OMOS-471: Execute View

### Design

```
┌─────────────────────────────────────┐
│ ⚡ Execute            Wed Jan 29    │
├─────────────────────────────────────┤
│ Focus Block: Deep Work 🧠  9-11am   │
│ ┌─────────────────────────────────┐ │
│ │ ○ Code review PR #123    45min │ │
│ │ ○ Write proposal         30min │ │
│ └─────────────────────────────────┘ │
├─────────────────────────────────────┤
│ Habits Today                  5/7   │
│ [✓Med] [✓Ex] [○Jrn] [✓Rd] [○Slp]  │
├─────────────────────────────────────┤
│ Up Next                             │
│ ┃ 11:00  Team standup    30min     │
│ ┃ 12:00  Lunch break               │
│ ┃ 13:00  Client call     1h        │
├─────────────────────────────────────┤
│ 💡 Suggested by Legacy:             │
│ "Move budget review to Friday -     │
│  your energy is low today"          │
│           [Accept] [Dismiss]        │
└─────────────────────────────────────┘
```

### Implementation

```dart
// frontend/lib/execute/screens/execute_screen.dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

class ExecuteScreen extends ConsumerWidget {
  const ExecuteScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final todayAssets = ref.watch(todayAssetsProvider);
    final suggestions = ref.watch(suggestionsProvider);

    return Scaffold(
      backgroundColor: const Color(0xFF0A0A0A),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Row(
          children: [
            const Icon(Icons.flash_on, color: Color(0xFFFF0000)),
            const SizedBox(width: 8),
            const Text(
              'Execute',
              style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: Center(
              child: Text(
                DateFormat('EEE MMM d').format(DateTime.now()),
                style: TextStyle(color: Colors.grey[400]),
              ),
            ),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(todayAssetsProvider),
        color: const Color(0xFFFF0000),
        child: todayAssets.when(
          data: (data) => ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Focus Block
              if (data.focusBlock != null)
                FocusBlockCard(block: data.focusBlock!),

              const SizedBox(height: 16),

              // Habits Quick Toggle
              HabitToggleRow(habits: data.habits),

              const SizedBox(height: 16),

              // Today's Tasks
              _SectionHeader(title: 'Tasks', count: data.tasks.length),
              ...data.tasks.map((task) => TaskItem(task: task)),

              const SizedBox(height: 16),

              // Schedule Timeline
              _SectionHeader(title: 'Up Next'),
              ScheduleTimeline(events: data.events),

              const SizedBox(height: 16),

              // AI Suggestion
              suggestions.when(
                data: (sug) => sug != null
                    ? SuggestionCard(suggestion: sug)
                    : const SizedBox.shrink(),
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),
            ],
          ),
          loading: () => const ExecuteScreenSkeleton(),
          error: (e, _) => ErrorWidget(message: e.toString()),
        ),
      ),
      floatingActionButton: FloatingActionButton.small(
        onPressed: () => _showQuickAdd(context),
        backgroundColor: const Color(0xFF1A1A1A),
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }
}
```

### Habit Toggle Row

```dart
// frontend/lib/execute/widgets/habit_toggle_row.dart
class HabitToggleRow extends ConsumerWidget {
  final List<Habit> habits;

  const HabitToggleRow({required this.habits, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final completed = habits.where((h) => h.completedToday).length;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A1A),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFF333333)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Habits Today',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
              Text(
                '$completed/${habits.length}',
                style: TextStyle(
                  color: Colors.grey[400],
                  fontSize: 14,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: habits.map((habit) => _HabitChip(
              habit: habit,
              onToggle: () => ref.read(habitsProvider.notifier).toggle(habit.id),
            )).toList(),
          ),
        ],
      ),
    );
  }
}

class _HabitChip extends StatelessWidget {
  final Habit habit;
  final VoidCallback onToggle;

  const _HabitChip({required this.habit, required this.onToggle});

  @override
  Widget build(BuildContext context) {
    final completed = habit.completedToday;

    return GestureDetector(
      onTap: onToggle,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: completed ? const Color(0xFF00FF00).withOpacity(0.2) : const Color(0xFF2A2A2A),
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: completed ? const Color(0xFF00FF00) : const Color(0xFF444444),
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              completed ? Icons.check_circle : Icons.circle_outlined,
              size: 16,
              color: completed ? const Color(0xFF00FF00) : Colors.grey,
            ),
            const SizedBox(width: 6),
            Text(
              habit.shortName,
              style: TextStyle(
                color: completed ? Colors.white : Colors.grey[400],
                fontSize: 13,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

---

## OMOS-480: Suggestion Card

```dart
// frontend/lib/execute/widgets/suggestion_card.dart
class SuggestionCard extends StatelessWidget {
  final Suggestion suggestion;

  const SuggestionCard({required this.suggestion, super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A1A),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: const Color(0xFFFF0000).withOpacity(0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.lightbulb_outline, color: Color(0xFFFF0000), size: 20),
              const SizedBox(width: 8),
              Text(
                'Suggested by Legacy',
                style: TextStyle(
                  color: Colors.grey[400],
                  fontSize: 12,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            suggestion.message,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
              height: 1.4,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              TextButton(
                onPressed: () => _dismiss(context),
                child: Text(
                  'Dismiss',
                  style: TextStyle(color: Colors.grey[400]),
                ),
              ),
              const SizedBox(width: 8),
              ElevatedButton(
                onPressed: () => _accept(context),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFFFF0000),
                  foregroundColor: Colors.white,
                ),
                child: const Text('Accept'),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
```

---

## Providers to Create

```dart
// frontend/lib/home/providers/home_providers.dart

final scoresProvider = FutureProvider<Scores>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/scores/all');
});

final todayFocusProvider = FutureProvider<Focus?>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/focus/today');
});

final activeAgentsProvider = FutureProvider<List<AgentStatus>>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/agents/active');
});

final quickStatsProvider = FutureProvider<QuickStats>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/stats/today');
});

// frontend/lib/execute/providers/execute_providers.dart

final todayAssetsProvider = FutureProvider<TodayAssets>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/assets/today');
});

final suggestionsProvider = FutureProvider<Suggestion?>((ref) async {
  final client = ref.watch(apiClientProvider);
  return await client.get('/api/v1/suggestions/next');
});
```

---

## Success Criteria

- [ ] Home screen shows unified score + individual scores
- [ ] Home screen shows today's focus with progress
- [ ] Home screen shows active agents (real-time)
- [ ] Home screen shows quick stats
- [ ] Execute screen shows today's tasks
- [ ] Execute screen shows habit toggles that work
- [ ] Execute screen shows schedule timeline
- [ ] Suggestion card with accept/dismiss
- [ ] Pull-to-refresh on both screens
- [ ] Skeleton loading states
- [ ] No crashes, smooth animations

---

## When Done

Update issues OMOS-470, 471, 480 to "In Review" and post handoff to Agent 8.
