# Agent 6: Mobile Navigation (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Flutter Mobile Navigation Redesign
> **Issues:** 1 (large) | **Points:** 8
> **Dependencies:** None - Start immediately

---

## Your Mission

You are Agent 6. Your job is the CRITICAL task of redesigning OneMind's navigation from 77 desktop screens to 12 mobile-first screens. This is the foundation that ALL other frontend agents (7, 8) depend on.

**Start immediately - you have no dependencies.**

---

## Your Assigned Issue

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-466** | Mobile Navigation Redesign | 8 | P0-Critical |

---

## Current State (What Exists)

### Current Navigation Structure
```
frontend/lib/
├── platform/router/app_router.dart      # GoRouter with 77+ routes
├── platform/providers/navigation_provider.dart  # Nav state
├── shared/widgets/
│   ├── app_shell.dart                   # Current bottom nav (5 items)
│   └── app_drawer.dart                  # Full drawer with 6 sections
```

### Current Bottom Nav (5 items)
```dart
// Current structure
[Command] [Life OS] [Inbox] [Awareness] [More →]
```

### Current Drawer (6 sections)
```
UIO (11 screens)
PATHS (14 screens)
SENSES (8 screens)
LEGACY (19 screens)
HIVE (6 screens)
CORE (4 screens)
```

---

## Target State (Build This)

### New Bottom Nav (4 items + FAB)
```
┌─────────────────────────────────────────────────────────┐
│                    SCREEN CONTENT                        │
├─────────────────────────────────────────────────────────┤
│                      [+ FAB]                             │  ← Floating center
├─────────────────────────────────────────────────────────┤
│  [Home]    [Command]         [Execute]    [Profile]     │
│   🏠          💬       +        ⚡           👤          │
└─────────────────────────────────────────────────────────┘
```

### New Route Structure
```
PRIMARY (Bottom Nav)
├── /home         ← Dashboard, scores, agents
├── /command      ← Chat with Legacy
├── /execute      ← Today's tasks/habits
├── /profile      ← Settings, config
└── /capture      ← Universal input (modal from FAB)

SECONDARY (From Drawer or Navigation)
├── /tasks        ← Full asset list (filtered)
├── /calendar     ← Time-based view
├── /insights     ← Analytics
├── /goals        ← Three Paths progress
├── /agents       ← Agent management
├── /knowledge    ← Memory and docs
└── /settings     ← Deep configuration (from Profile)
```

---

## Implementation

### Step 1: Update app_shell.dart

**Modify:** `frontend/lib/shared/widgets/app_shell.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class AppShell extends ConsumerWidget {
  final Widget child;

  const AppShell({required this.child, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final currentPath = GoRouterState.of(context).uri.path;

    return Scaffold(
      body: child,
      drawer: const SimplifiedDrawer(),
      floatingActionButton: _buildCaptureFAB(context),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      bottomNavigationBar: _buildBottomNav(context, currentPath),
    );
  }

  Widget _buildCaptureFAB(BuildContext context) {
    return FloatingActionButton(
      onPressed: () => _showCaptureModal(context),
      backgroundColor: const Color(0xFFFF0000), // Tactical red
      elevation: 4,
      child: const Icon(Icons.add, size: 28, color: Colors.white),
    );
  }

  void _showCaptureModal(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => const CaptureModal(),
    );
  }

  Widget _buildBottomNav(BuildContext context, String currentPath) {
    return BottomAppBar(
      shape: const CircularNotchedRectangle(),
      notchMargin: 8,
      color: const Color(0xFF1A1A1A),
      child: SizedBox(
        height: 60,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            _NavItem(
              icon: Icons.home_outlined,
              activeIcon: Icons.home,
              label: 'Home',
              route: '/home',
              isActive: currentPath == '/home' || currentPath == '/',
              onTap: () => context.go('/home'),
            ),
            _NavItem(
              icon: Icons.chat_bubble_outline,
              activeIcon: Icons.chat_bubble,
              label: 'Command',
              route: '/command',
              isActive: currentPath == '/command',
              onTap: () => context.go('/command'),
            ),
            const SizedBox(width: 56), // Space for FAB
            _NavItem(
              icon: Icons.flash_on_outlined,
              activeIcon: Icons.flash_on,
              label: 'Execute',
              route: '/execute',
              isActive: currentPath == '/execute',
              onTap: () => context.go('/execute'),
            ),
            _NavItem(
              icon: Icons.person_outline,
              activeIcon: Icons.person,
              label: 'Profile',
              route: '/profile',
              isActive: currentPath == '/profile',
              onTap: () => context.go('/profile'),
            ),
          ],
        ),
      ),
    );
  }
}

class _NavItem extends StatelessWidget {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final String route;
  final bool isActive;
  final VoidCallback onTap;

  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.route,
    required this.isActive,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final color = isActive ? const Color(0xFFFF0000) : Colors.grey;

    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              isActive ? activeIcon : icon,
              color: color,
              size: 24,
            ),
            const SizedBox(height: 4),
            Text(
              label,
              style: TextStyle(
                color: color,
                fontSize: 10,
                fontWeight: isActive ? FontWeight.w600 : FontWeight.normal,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// Placeholder for Agent 8 to implement
class CaptureModal extends StatelessWidget {
  const CaptureModal({super.key});

  @override
  Widget build(BuildContext context) {
    return DraggableScrollableSheet(
      initialChildSize: 0.5,
      minChildSize: 0.3,
      maxChildSize: 0.9,
      builder: (context, scrollController) {
        return Container(
          decoration: const BoxDecoration(
            color: Color(0xFF1A1A1A),
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            children: [
              const SizedBox(height: 12),
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[600],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              const Text(
                'Capture',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const Expanded(
                child: Center(
                  child: Text(
                    'Agent 8 will implement capture UI',
                    style: TextStyle(color: Colors.grey),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
```

### Step 2: Simplified Drawer

**Modify:** `frontend/lib/shared/widgets/app_drawer.dart`

```dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SimplifiedDrawer extends StatelessWidget {
  const SimplifiedDrawer({super.key});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      backgroundColor: const Color(0xFF0A0A0A),
      child: SafeArea(
        child: Column(
          children: [
            // Header
            Container(
              padding: const EdgeInsets.all(20),
              child: Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: const Color(0xFFFF0000),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Icon(Icons.psychology, color: Colors.white),
                  ),
                  const SizedBox(width: 12),
                  const Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'OneMind',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        'Mobile-First',
                        style: TextStyle(
                          color: Colors.grey,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),

            const Divider(color: Color(0xFF333333)),

            // Secondary Navigation
            Expanded(
              child: ListView(
                padding: const EdgeInsets.symmetric(vertical: 8),
                children: [
                  _DrawerSection(
                    title: 'Life OS',
                    items: [
                      _DrawerItem(icon: Icons.task_alt, label: 'Tasks', route: '/tasks'),
                      _DrawerItem(icon: Icons.calendar_today, label: 'Calendar', route: '/calendar'),
                      _DrawerItem(icon: Icons.flag, label: 'Goals', route: '/goals'),
                      _DrawerItem(icon: Icons.trending_up, label: 'Habits', route: '/habits'),
                    ],
                  ),
                  _DrawerSection(
                    title: 'Intelligence',
                    items: [
                      _DrawerItem(icon: Icons.insights, label: 'Insights', route: '/insights'),
                      _DrawerItem(icon: Icons.smart_toy, label: 'Agents', route: '/agents'),
                      _DrawerItem(icon: Icons.memory, label: 'Knowledge', route: '/knowledge'),
                    ],
                  ),
                  _DrawerSection(
                    title: 'System',
                    items: [
                      _DrawerItem(icon: Icons.home, label: 'Home Assistant', route: '/home-assistant'),
                      _DrawerItem(icon: Icons.settings, label: 'Settings', route: '/settings'),
                    ],
                  ),
                ],
              ),
            ),

            // Footer
            const Divider(color: Color(0xFF333333)),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.circle, size: 8, color: Colors.green[400]),
                  const SizedBox(width: 8),
                  const Text(
                    'Connected',
                    style: TextStyle(color: Colors.grey, fontSize: 12),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _DrawerSection extends StatelessWidget {
  final String title;
  final List<_DrawerItem> items;

  const _DrawerSection({required this.title, required this.items});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Text(
            title.toUpperCase(),
            style: const TextStyle(
              color: Colors.grey,
              fontSize: 11,
              fontWeight: FontWeight.w600,
              letterSpacing: 1,
            ),
          ),
        ),
        ...items,
      ],
    );
  }
}

class _DrawerItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final String route;

  const _DrawerItem({
    required this.icon,
    required this.label,
    required this.route,
  });

  @override
  Widget build(BuildContext context) {
    final currentPath = GoRouterState.of(context).uri.path;
    final isActive = currentPath == route;

    return ListTile(
      dense: true,
      leading: Icon(
        icon,
        color: isActive ? const Color(0xFFFF0000) : Colors.grey,
        size: 20,
      ),
      title: Text(
        label,
        style: TextStyle(
          color: isActive ? Colors.white : Colors.grey[400],
          fontSize: 14,
        ),
      ),
      selected: isActive,
      selectedTileColor: const Color(0xFF1A1A1A),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      onTap: () {
        Navigator.pop(context); // Close drawer
        context.go(route);
      },
    );
  }
}
```

### Step 3: Update Router

**Modify:** `frontend/lib/platform/router/app_router.dart`

```dart
import 'package:go_router/go_router.dart';
import 'package:flutter/material.dart';

// Import your screens
import '../../shared/widgets/app_shell.dart';

final appRouter = GoRouter(
  initialLocation: '/home',
  routes: [
    // Shell route wraps all main screens
    ShellRoute(
      builder: (context, state, child) => AppShell(child: child),
      routes: [
        // PRIMARY (Bottom Nav)
        GoRoute(
          path: '/home',
          name: 'home',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: HomeScreen(), // Agent 7 will build
          ),
        ),
        GoRoute(
          path: '/command',
          name: 'command',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: CommandScreen(), // Agent 8 will build
          ),
        ),
        GoRoute(
          path: '/execute',
          name: 'execute',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: ExecuteScreen(), // Agent 7 will build
          ),
        ),
        GoRoute(
          path: '/profile',
          name: 'profile',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: ProfileScreen(), // Agent 8 will build
          ),
        ),

        // SECONDARY (From Drawer)
        GoRoute(path: '/tasks', name: 'tasks', builder: (c, s) => const TasksScreen()),
        GoRoute(path: '/calendar', name: 'calendar', builder: (c, s) => const CalendarScreen()),
        GoRoute(path: '/goals', name: 'goals', builder: (c, s) => const GoalsScreen()),
        GoRoute(path: '/habits', name: 'habits', builder: (c, s) => const HabitsScreen()),
        GoRoute(path: '/insights', name: 'insights', builder: (c, s) => const InsightsScreen()),
        GoRoute(path: '/agents', name: 'agents', builder: (c, s) => const AgentsScreen()),
        GoRoute(path: '/knowledge', name: 'knowledge', builder: (c, s) => const KnowledgeScreen()),
        GoRoute(path: '/home-assistant', name: 'homeAssistant', builder: (c, s) => const HomeAssistantScreen()),
        GoRoute(path: '/settings', name: 'settings', builder: (c, s) => const SettingsScreen()),
      ],
    ),

    // Redirect old routes to new ones
    GoRoute(path: '/', redirect: (_, __) => '/home'),
    GoRoute(path: '/dashboard', redirect: (_, __) => '/home'),
    GoRoute(path: '/chat', redirect: (_, __) => '/command'),
    GoRoute(path: '/today', redirect: (_, __) => '/execute'),
  ],
);

// Placeholder screens (Agents 7, 8 will replace)
class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color(0xFF0A0A0A),
      body: Center(
        child: Text('Home - Agent 7', style: TextStyle(color: Colors.white)),
      ),
    );
  }
}

class CommandScreen extends StatelessWidget {
  const CommandScreen({super.key});
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color(0xFF0A0A0A),
      body: Center(
        child: Text('Command - Agent 8', style: TextStyle(color: Colors.white)),
      ),
    );
  }
}

class ExecuteScreen extends StatelessWidget {
  const ExecuteScreen({super.key});
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color(0xFF0A0A0A),
      body: Center(
        child: Text('Execute - Agent 7', style: TextStyle(color: Colors.white)),
      ),
    );
  }
}

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});
  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color(0xFF0A0A0A),
      body: Center(
        child: Text('Profile - Agent 8', style: TextStyle(color: Colors.white)),
      ),
    );
  }
}

// Add placeholder classes for other screens...
class TasksScreen extends StatelessWidget {
  const TasksScreen({super.key});
  @override
  Widget build(BuildContext context) => const Center(child: Text('Tasks'));
}
// ... etc for all secondary screens
```

---

## Success Criteria

- [ ] New bottom nav with 4 items + center FAB
- [ ] FAB opens capture modal
- [ ] Drawer simplified with 3 sections
- [ ] All routes working without 404s
- [ ] Old routes redirect correctly
- [ ] Smooth transitions between screens
- [ ] Works on both iOS and Android

---

## When Done

### 1. Update Linear Issue to "In Review"

Update OMOS-466

### 2. Post Handoff

```markdown
## Handoff: Agent 6 → Agents 7, 8

### Completed
- ✅ OMOS-466: Mobile Navigation Redesign
- New bottom nav in `app_shell.dart`
- Simplified drawer in `app_drawer.dart`
- Consolidated routes in `app_router.dart`
- CaptureModal placeholder ready

### For Agent 7 (Home & Execute)
- Implement `HomeScreen` at `/home`
- Implement `ExecuteScreen` at `/execute`
- Screens will appear in AppShell automatically

### For Agent 8 (Command, Capture, Profile)
- Implement `CommandScreen` at `/command`
- Implement `ProfileScreen` at `/profile`
- Replace `CaptureModal` placeholder with full implementation
- Implement secondary screens (Insights, Agents, Knowledge)

### Navigation Usage
```dart
// Navigate to primary screens
context.go('/home');
context.go('/command');
context.go('/execute');
context.go('/profile');

// Open capture modal
_showCaptureModal(context); // Already in AppShell
```

### Color Scheme
- Background: `#0A0A0A`
- Cards: `#1A1A1A`
- Borders: `#333333`
- Accent: `#FF0000` (tactical red)
```

---

**You are the first frontend agent to complete. Agents 7 and 8 are waiting on you. Ship fast, test on device, document clearly.**
