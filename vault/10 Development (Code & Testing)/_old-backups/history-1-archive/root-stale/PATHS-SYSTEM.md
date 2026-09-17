# OneMind OS PATHS System

> **Document Status:** CANONICAL
> **Created:** 2026-01-28
> **Section:** PATHS - The Three Paths
> **Philosophy:** Life management across HP/LE/GE domains

---

## Overview

PATHS is the life management hub of OneMind OS, representing the **human side** of the OneMind equation. While LEGACY manages the AI workforce, PATHS manages Zeus's life across three doctrinal domains derived from the OneMind Codex.

### Core Equation

```
HUMAN (Zeus) + AI (Legacy) = ONE MIND
  ↓              ↓
PATHS         LEGACY
```

PATHS embodies the structured approach to life optimization, ensuring that personal growth, family legacy, and entrepreneurial ventures are all managed with intention and AI-assisted intelligence.

---

## The Three Paths

The Three Paths framework originates from the OneMind Codex and represents the three domains of a well-lived life:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           THE THREE PATHS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│         ┌───────────────┐                                                   │
│         │      HP       │                                                   │
│         │   Holistic    │                                                   │
│         │  Performance  │                                                   │
│         └───────┬───────┘                                                   │
│                 │                                                           │
│     ┌───────────┴───────────┐                                               │
│     │                       │                                               │
│     ▼                       ▼                                               │
│ ┌───────────────┐   ┌───────────────┐                                       │
│ │      LE       │   │      GE       │                                       │
│ │    Legacy     │   │  Generational │                                       │
│ │   Evolution   │   │Entrepreneurship                                       │
│ └───────────────┘   └───────────────┘                                       │
│                                                                             │
│  "Optimize self (HP) to serve family (LE) and create value (GE)"           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### HP - Holistic Performance

**Focus:** Self-optimization, health, habits, personal effectiveness

**Philosophy:** You cannot pour from an empty cup. HP represents the foundation - optimizing your own performance so you have the capacity to serve your family and build your ventures.

**Domains:**
| Domain | Description | Features |
|--------|-------------|----------|
| Physical | Body optimization | Health tracking, meals, exercise |
| Mental | Cognitive performance | Focus sessions, journaling, meditation |
| Temporal | Time mastery | Scheduling, routines, pomodoro |
| Habitual | Behavior engineering | Habit tracking, streaks, triggers |

**Current Features:**
- Habits - Daily habit tracking with streaks
- Routines - Sequential action chains
- Journal - Reflection and logging
- Pomodoro - Focus session timer
- Meals - Nutrition planning (recipes, shopping, pantry)
- Scheduling - Time blocking
- Planner - Daily/weekly planning

**Metrics:** Feeds into the **HPI (Holistic Performance Intelligence)** score via SENSES.

---

### LE - Legacy Evolution

**Focus:** Family, home, heritage, relationships, generational continuity

**Philosophy:** What you build means nothing if it doesn't outlast you. LE represents building and nurturing the human connections and home environment that create lasting legacy.

**Domains:**
| Domain | Description | Future Features |
|--------|-------------|-----------------|
| Family | Household members | Family profiles, shared calendars |
| Home | Living environment | Home management, maintenance schedules |
| Relationships | Social connections | Contact management, relationship health |
| Heritage | Cultural continuity | Family traditions, memory preservation |

**Planned Features:**
- Family Module - Family member profiles and management
- Shared Tasks - Tasks assigned to/from family members
- Shared Calendar - Family-wide event coordination
- Relationship Tracker - Connection health monitoring
- Home Management - Property maintenance, inventory

**Metrics:** Feeds into the **LEI (Legacy Evolution Intelligence)** score via SENSES.

---

### GE - Generational Empire

**Focus:** Business, wealth, value creation, financial independence

**Philosophy:** True wealth is built across generations. GE represents the entrepreneurial ventures and financial strategies that create lasting prosperity.

**Domains:**
| Domain | Description | Future Features |
|--------|-------------|-----------------|
| Finance | Wealth management | Budgets, investments, net worth |
| Business | Venture management | Projects, clients, revenue |
| Assets | Property & holdings | Real estate, vehicles, valuables |
| Strategy | Long-term planning | Goals, milestones, projections |

**Planned Features:**
- Finance Module - Budget tracking, investment monitoring
- Business Dashboard - Revenue, clients, projects
- Asset Tracker - Property, vehicles, valuables inventory
- Net Worth Calculator - Comprehensive wealth view
- Financial Goals - Savings targets, retirement planning

**Metrics:** Feeds into the **GEI (Generational Empire Intelligence)** score via SENSES.

---

## Path Intersection

Features in PATHS often span multiple paths. This is intentional - life doesn't fit into neat boxes.

### Intersection Matrix

| Feature | HP | LE | GE | Example |
|---------|----|----|-----|---------|
| Tasks | ✓ | ✓ | ✓ | Personal task (HP), family chore (LE), business deliverable (GE) |
| Projects | ✓ | ✓ | ✓ | Fitness program (HP), home renovation (LE), product launch (GE) |
| Calendar | ✓ | ✓ | ✓ | Workout (HP), family dinner (LE), client meeting (GE) |
| Goals | ✓ | ✓ | ✓ | Run marathon (HP), college fund (LE), revenue target (GE) |
| Meals | ✓ | ✓ | - | Personal nutrition (HP), family dinners (LE) |
| Journal | ✓ | ✓ | ✓ | Daily reflection (HP), family memories (LE), business insights (GE) |
| Habits | ✓ | - | - | Primarily personal optimization |
| Routines | ✓ | ✓ | - | Morning routine (HP), family evening routine (LE) |
| Finance | - | ✓ | ✓ | Family budget (LE), business finances (GE) |

### Tagging System (Future)

Each item in PATHS will support path tagging:

```dart
class Task {
  final String id;
  final String title;
  final Set<PathType> paths; // {PathType.hp, PathType.le, PathType.ge}
  final String? familyMemberId; // For LE shared tasks
  final String? businessId; // For GE business tasks
}
```

---

## Current Feature Map

### Navigation Structure

```
PATHS/
├── Tasks ────────────── /tasks
├── Projects ─────────── /projects
├── Planner ──────────── /planner
├── Calendar ─────────── /calendar
├── Habits ───────────── /habits
├── Routines ─────────── /routines
├── Scheduling ───────── /scheduling
├── Pomodoro ─────────── /pomodoro
├── Journal ──────────── /journal
├── Goals ────────────── /goals
└── Meals/
    ├── Meal Plan ────── /meals
    ├── Recipes ──────── /meals/recipes
    ├── Shopping ─────── /meals/shopping
    └── Pantry ───────── /meals/pantry
```

### Feature Details

| Feature | Route | Path | Description |
|---------|-------|------|-------------|
| Tasks | `/tasks` | HP/LE/GE | Task management with priorities, due dates |
| Projects | `/projects` | HP/LE/GE | Multi-task initiatives with milestones |
| Planner | `/planner` | HP | Daily/weekly planning view |
| Calendar | `/calendar` | HP/LE/GE | Time-based event scheduling |
| Habits | `/habits` | HP | Habit tracking with streaks |
| Routines | `/routines` | HP/LE | Sequential action chains |
| Scheduling | `/scheduling` | HP | Time blocking interface |
| Pomodoro | `/pomodoro` | HP | Focus session timer (25/5 pattern) |
| Journal | `/journal` | HP/LE/GE | Reflection and logging |
| Goals | `/goals` | HP/LE/GE | Objective tracking with key results |
| Meal Plan | `/meals` | HP/LE | Weekly meal planning |
| Recipes | `/meals/recipes` | HP/LE | Recipe library and discovery |
| Shopping | `/meals/shopping` | HP/LE | Grocery list management |
| Pantry | `/meals/pantry` | HP/LE | Inventory tracking |

---

## Multi-User Architecture (Future)

PATHS will evolve to support multiple family members, enabling true LE functionality.

### User Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                     FAMILY STRUCTURE                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐                                                │
│  │    Zeus     │  ◄── Primary user (full access)               │
│  │   (Owner)   │                                                │
│  └──────┬──────┘                                                │
│         │                                                       │
│    ┌────┴────┬────────────┐                                     │
│    ▼         ▼            ▼                                     │
│ ┌──────┐ ┌──────┐    ┌──────┐                                   │
│ │Spouse│ │Child │    │Child │  ◄── Family members              │
│ │(Admin)│ │(User)│    │(User)│     (scoped access)              │
│ └──────┘ └──────┘    └──────┘                                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Shared Features

| Feature | Sharing Model |
|---------|---------------|
| Tasks | Assignable to/from family members |
| Calendar | Shared family calendar + personal calendars |
| Goals | Family goals visible to all, personal goals private |
| Meals | Shared meal plan, individual dietary preferences |
| Projects | Family projects with member assignments |

### Permission Levels

| Role | Permissions |
|------|-------------|
| Owner | Full access, system configuration, billing |
| Admin | Manage family members, shared resources |
| User | View shared items, manage own items |
| Child | Limited view, assigned tasks only |

---

## Integration with OneMind OS

### SENSES Integration

PATHS data feeds into SENSES for intelligence scoring:

```
PATHS (Input)              SENSES (Analysis)           Output
─────────────              ─────────────────           ──────
Habits completed    ───►   HPI Score Algorithm  ───►   Holistic Performance Intelligence
Exercise logged     ───►   (HP Path)
Sleep data          ───►

Family events       ───►   LEI Score Algorithm  ───►   Legacy Evolution Intelligence
Relationship logs   ───►   (LE Path)
Shared tasks        ───►

Revenue tracked     ───►   GEI Score Algorithm  ───►   Generational Empire Intel
Investments         ───►   (GE Path)
Business metrics    ───►
```

### LEGACY Integration

AI agents in LEGACY can assist with PATHS features:

| Agent Use Case | Example |
|----------------|---------|
| Task Assistant | "Schedule my high-priority tasks for tomorrow" |
| Meal Planner | "Plan healthy meals for the week under $200" |
| Goal Coach | "Break down my marathon goal into weekly milestones" |
| Family Coordinator | "Find a time that works for family dinner this week" |
| Financial Advisor | "Analyze my spending patterns this month" |

### HIVE Integration

Physical devices in HIVE enhance PATHS tracking:

| Device | PATHS Integration |
|--------|-------------------|
| Watch | Habit tracking, exercise logging, reminders |
| Glasses | AR task overlays, recipe display while cooking |
| Home Assistant | Routine triggers, family presence detection |
| Vehicles | Commute tracking, errand optimization |

---

## Data Models

### Core Entities

```dart
/// The three paths
enum PathType {
  hp,  // Holistic Performance
  le,  // Legacy Evolution
  ge,  // Generational Empire
}

/// Family member (future)
class FamilyMember {
  final String id;
  final String name;
  final FamilyRole role;
  final String? photoUrl;
  final DateTime? birthDate;
}

/// Task with path tagging (future enhancement)
class Task {
  final String id;
  final String title;
  final String? description;
  final Set<PathType> paths;
  final String? assignedTo;      // FamilyMember ID
  final String? assignedBy;      // FamilyMember ID
  final DateTime? dueDate;
  final TaskPriority priority;
  final bool isCompleted;
}

/// Goal with path alignment
class Goal {
  final String id;
  final String title;
  final PathType primaryPath;
  final Set<PathType> secondaryPaths;
  final List<KeyResult> keyResults;
  final DateTime targetDate;
  final double progress;
}
```

---

## Roadmap

### Phase 1: Foundation (Current)
- [x] Basic task management
- [x] Habit tracking
- [x] Meal planning system
- [x] Calendar integration
- [x] Goal tracking
- [x] Pomodoro timer
- [x] Journal

### Phase 2: Path Tagging
- [ ] Add path tags to all entities
- [ ] Path-filtered views
- [ ] Path-based analytics
- [ ] Cross-path correlations in SENSES

### Phase 3: Family Module (LE)
- [ ] Family member profiles
- [ ] Shared calendar
- [ ] Task assignment
- [ ] Family goals
- [ ] Relationship tracking

### Phase 4: Finance Module (GE)
- [ ] Budget tracking
- [ ] Investment monitoring
- [ ] Net worth calculation
- [ ] Financial goals
- [ ] Business dashboard

### Phase 5: Multi-User
- [ ] Family member accounts
- [ ] Permission system
- [ ] Shared/private item distinction
- [ ] Family activity feed

---

## Design System

### Colors by Path

| Path | Primary Color | Hex Code | Usage |
|------|---------------|----------|-------|
| HP | Blue | `#3B82F6` | Health, performance, personal growth |
| LE | Purple | `#8B5CF6` | Family, legacy, relationships |
| GE | Green | `#22C55E` | Business, wealth, entrepreneurship |
| Mixed | Red | `#FF0000` | OneMind primary accent |

### Path Indicators

```
┌─────────────────────────────────────┐
│  [HP] Morning Workout               │  ← Single path tag
│  ────────────────────────────────── │
│  [HP][LE] Family Hike               │  ← Multi-path tag
│  ────────────────────────────────── │
│  [HP][LE][GE] Life Planning Review  │  ← All paths
└─────────────────────────────────────┘
```

---

## File Locations

| File | Purpose |
|------|---------|
| `frontend/lib/lifeos/tasks/` | Task management module |
| `frontend/lib/lifeos/habits/` | Habit tracking module |
| `frontend/lib/lifeos/meals/` | Meal planning module |
| `frontend/lib/lifeos/calendar/` | Calendar module |
| `frontend/lib/lifeos/goals/` | Goal tracking module |
| `frontend/lib/lifeos/journal/` | Journal module |
| `frontend/lib/shared/widgets/app_drawer.dart` | PATHS navigation |

---

## Related Documentation

- [Navigation System](./NAVIGATION-SYSTEM.md) - Overall navigation architecture
- [OneMind Codex](../../CODEX-FEATURE-ROADMAP.md) - Feature roadmap and philosophy
- [Skills System](./SKILLS-SYSTEM.md) - AI skill packages

---

## Summary

PATHS is not just a collection of productivity tools - it's a philosophical framework for intentional living. By organizing life into three paths (HP/LE/GE), OneMind OS helps Zeus:

1. **Optimize** personal performance (HP)
2. **Nurture** family and legacy (LE)
3. **Build** generational wealth (GE)

All while LEGACY (AI) assists and SENSES monitors progress toward the ultimate goal: a well-lived life that creates lasting impact.

---

*"Optimize self (HP) to serve family (LE) and create value (GE)"*

*Part of the PATHS section - Life Management layer of OneMind OS*
