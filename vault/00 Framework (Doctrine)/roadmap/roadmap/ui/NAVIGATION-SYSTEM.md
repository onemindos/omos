# OneMind OS Navigation System

> **Document Status:** CANONICAL
> **Created:** 2026-01-28
> **Purpose:** Define the navigation architecture and naming conventions aligned with OneMind philosophy

---

## Philosophy

The OneMind navigation structure reflects the core equation:

```
HUMAN (Zeus) + AI (Legacy) = ONE MIND
```

Each navigation section represents a distinct aspect of the unified system, creating a cohesive experience that bridges human intent with AI capability.

---

## Navigation Sections

### Visual Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        ONEMIND OPERATING SYSTEM                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│  │   UIO   │  │  PATHS  │  │ SENSES  │  │ LEGACY  │  │  HIVE   │  │  CORE   │
│  │         │  │         │  │         │  │         │  │         │  │         │
│  │ Command │  │  Life   │  │  Intel  │  │   AI    │  │Physical │  │ System  │
│  │ Center  │  │ Mgmt    │  │  Layer  │  │Workforce│  │Interface│  │ Config  │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
│       │            │            │            │            │            │
│       ▼            ▼            ▼            ▼            ▼            ▼
│   Dashboard    Tasks       Awareness     Agents      Home Asst   Settings
│   Command      Habits      Consciousness Teams       Vehicles    Databases
│   Eagle Eye    Goals       Presence      Tools       Robotics    Profile
│   Voice        Calendar    Health Score  Skills      Drones      Integrations
│   Inbox        Journal     Trends        Memory      Glasses
│   Activity     Meals       Correlations  Sessions    Watch
│   Overview     Planner     Alerts        Workflows
│   Metrics      Projects    Log           Approvals
│                Pomodoro                  Knowledge
│                Routines                  MCP
│                Scheduling                Culture
│                                          Traces
│                                          Evals
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section Definitions

### 1. UIO - Sovereign Operations Operations

**Purpose:** The command center and primary interface for OneMind operations.

**Philosophy:** UIO represents the "brain" of the system - where Zeus observes, commands, and orchestrates all operations. This is the unified view that brings everything together.

**Contains:**
| Screen | Route | Description |
|--------|-------|-------------|
| Dashboard | `/dashboard` | System overview and quick stats |
| Command | `/` | Primary AI interaction terminal |
| Eagle Eye | `/eagle-eye` | High-level situational awareness |
| Voice | `/voice` | Voice-based AI interaction |
| Inbox | `/inbox` | Universal capture and triage |
| Activity | `/activity` | Real-time notification feed |
| Overview | `/sight` | System-wide visibility |
| Edge | `/sight/edge` | Edge device connections |
| Pipeline | `/sight/pipeline` | Task and workflow pipeline |
| Metrics | `/metrics` | Performance and analytics |

**Visual Design:**
- Primary accent: Red (#FF0000)
- Active indicator: Red glow effect
- Icon style: Tactical/operational

---

### 2. PATHS - The Three Paths (HP/LE/GE)

**Purpose:** Life management across the three doctrinal paths.

**Philosophy:** PATHS embodies the Three Paths from the OneMind Codex:
- **HP (Holistic Performance):** Self-optimization, health, habits
- **LE (Legacy Evolution):** Family, home, heritage
- **GE (Generational Endeavor):** Business, wealth, value creation

**Contains:**
| Screen | Route | Path | Description |
|--------|-------|------|-------------|
| Tasks | `/tasks` | HP/LE/GE | Task management |
| Projects | `/projects` | HP/LE/GE | Multi-task initiatives |
| Planner | `/planner` | HP | Daily/weekly planning |
| Calendar | `/calendar` | HP/LE/GE | Time-based scheduling |
| Habits | `/habits` | HP | Habit tracking |
| Routines | `/routines` | HP | Routine sequences |
| Scheduling | `/scheduling` | HP | Time blocking |
| Pomodoro | `/pomodoro` | HP | Focus sessions |
| Journal | `/journal` | HP | Reflection and logging |
| Goals | `/goals` | HP/LE/GE | Objective tracking |
| Meal Plan | `/meals` | HP | Nutrition planning |
| Recipes | `/meals/recipes` | HP | Recipe library |
| Shopping | `/meals/shopping` | HP | Grocery lists |
| Pantry | `/meals/pantry` | HP | Inventory tracking |

**Future Expansion:**
- Finance Module (GE path)
- Family Module (LE path)
- Business Dashboard (GE path)
- Relationships (HP/LE path)

---

### 3. SENSES - Awareness & Intelligence Layer

**Purpose:** How the system perceives, monitors, and understands.

**Philosophy:** SENSES represents the nervous system of OneMind - the awareness layer that monitors everything. It provides the intelligence data that feeds into the HPI/LEI/GEI scores.

**Contains:**
| Screen | Route | Description |
|--------|-------|-------------|
| Awareness | `/awareness` | Current awareness state |
| Consciousness | `/consciousness` | System consciousness view |
| Presence | `/presence` | Location and context awareness |
| Health Score | `/intel` | Personal analytics dashboard |
| Trends | `/intel/trends` | Long-term pattern analysis |
| Correlations | `/intel/correlations` | Cross-metric correlations |
| Alerts | `/intel/alerts` | Important notifications |
| Log | `/intel/log` | Manual data entry |

**Integration:**
- Feeds data to HPI (Holistic Performance Intelligence) score
- Connects to biometrics (Oura, Apple Watch)
- Correlates across all life domains

---

### 4. LEGACY - AI Workforce

**Purpose:** The AI agents, teams, and capabilities that work alongside Zeus.

**Philosophy:** LEGACY is named after the AI partner in the OneMind equation. This section manages the AI workforce - the agents, teams, tools, and knowledge that execute tasks autonomously or with human-in-the-loop approval.

**Contains:**

#### Core Agents
| Screen | Route | Description |
|--------|-------|-------------|
| Agents | `/agents` | Individual AI agent management |
| Teams | `/teams` | Multi-agent team configurations |
| Models | `/models` | LLM model selection and config |

#### Automation
| Screen | Route | Description |
|--------|-------|-------------|
| Workflows | `/workflows` | Automated workflow definitions |
| Approvals | `/approvals` | HITL approval queue |
| Guardrails | `/guardrails` | Safety constraints |

#### Capabilities
| Screen | Route | Description |
|--------|-------|-------------|
| Tools | `/tools` | Function/tool library |
| Skills | `/skills` | Skill packages (progressive disclosure) |
| Knowledge | `/knowledge` | Knowledge base and RAG |
| MCP | `/mcp` | Model Context Protocol servers |

#### Memory & Learning
| Screen | Route | Description |
|--------|-------|-------------|
| Memory | `/memory` | Agent memory management |
| Sessions | `/sessions` | Conversation history |
| Entities | `/entities` | Named entity recognition |
| Culture | `/culture` | Agent culture and values |

#### Monitoring
| Screen | Route | Description |
|--------|-------|-------------|
| Traces | `/traces` | Execution traces and logs |
| Evals | `/evals` | Agent evaluation results |

---

### 5. HIVE - Physical Interface

**Purpose:** The physical world connection - IoT, robotics, and wearables.

**Philosophy:** HIVE represents the networked mesh of physical devices that extend OneMind into the real world. The name evokes a connected network of sensors and actuators working together.

**Contains:**
| Screen | Route | Description |
|--------|-------|-------------|
| Home Assistant | `/home-assistant` | Smart home control |
| Vehicles | `/vehicles` | Vehicle fleet management (cars, motorcycles, jets, etc.) |
| Robotics | `/robotics` | Robot management |
| Drones | `/drones` | Drone fleet control |
| Glasses | `/glasses` | AR glasses interface |
| Watch | `/watch` | Smartwatch connection |

**Vehicle Integration:** See [GitHub Issue #97](https://github.com/Zeus-Delacruz/OneMind-OS/issues/97) for the full vehicle integration roadmap including:
- Multi-vehicle type support (cars, motorcycles, jets, boats, etc.)
- AI-powered software module installation
- Edge deployment for real-time telemetry
- OBD-II and manufacturer API integrations

**Future Expansion:**
- Security cameras
- Environmental sensors
- Automated equipment

---

### 6. CORE - System Foundation

**Purpose:** Configuration, settings, and system administration.

**Philosophy:** CORE is the foundation layer - the settings and infrastructure that everything else builds upon.

**Contains:**
| Screen | Route | Description |
|--------|-------|-------------|
| Profile | `/settings/profile` | User profile settings |
| Databases | `/databases` | Database management |
| Integrations | `/settings/integrations` | Third-party connections |
| Settings | `/settings` | System configuration |

---

## Navigation Methods

OneMind provides three navigation methods:

### 1. App Drawer (Side Menu)
- **Trigger:** Hamburger menu or swipe from left
- **Style:** Accordion sections with subgroups
- **Features:**
  - Favorites (star icons)
  - Badge counts for pending items
  - Section item counts
  - Subgroup headers (in LEGACY section)

### 2. Bottom Navigation
- **Trigger:** Always visible on mobile
- **Style:** Customizable pinned items
- **Features:**
  - Max 5 pinned items
  - Reorderable via drag
  - Customization modal

### 3. Command Palette
- **Trigger:** `⌘K` (Mac) or `Ctrl+K` (Windows)
- **Style:** Spotlight-like search
- **Features:**
  - Fuzzy search across all screens
  - Recent screens
  - Quick actions

---

## File Locations

| File | Purpose |
|------|---------|
| [app_drawer.dart](frontend/lib/shared/widgets/app_drawer.dart) | Main navigation drawer |
| [bottom_nav_customization_modal.dart](frontend/lib/shared/widgets/bottom_nav_customization_modal.dart) | Bottom nav customization |
| [command_palette.dart](frontend/lib/shared/widgets/command_palette.dart) | Command palette search |
| [navigation_provider.dart](frontend/lib/platform/providers/navigation_provider.dart) | Navigation state management |
| [favorites_provider.dart](frontend/lib/platform/providers/favorites_provider.dart) | Favorites management |
| [recent_screens_provider.dart](frontend/lib/platform/providers/recent_screens_provider.dart) | Recent screens tracking |

---

## Design System

### Colors
| Element | Color | Usage |
|---------|-------|-------|
| Primary Accent | `#FF0000` | Active states, highlights |
| Background | `#0A0A0A` | Main background |
| Surface | `#0D0D0D` | Headers, cards |
| Text Primary | `#FFFFFF` | Main text |
| Text Secondary | `#FFFFFF70` | Muted text |
| Success | `#00FF00` | Connected, success states |
| Warning | `#FFD700` | Warnings, cautions |

### Typography
- **Section Headers:** 12px, weight 700, monospace, letter-spacing 2px
- **Nav Items:** 14px, weight 400-600, letter-spacing 0.5px
- **Subheaders:** 10px, weight 600, monospace, letter-spacing 1.5px
- **Badges:** 10px, weight 600-700

### Components

#### Accordion Section
```dart
_AccordionSection(
  title: 'SECTION_NAME',
  currentRoute: currentRoute,
  expandedSection: _expandedSection,
  onExpand: (section) => setState(() { ... }),
  items: [...],
)
```

#### Subheader (within LEGACY section)
```dart
const _DrawerNavItem.subheader('SUBGROUP NAME'),
```

#### Navigation Item
```dart
_DrawerNavItem('Label', Icons.icon_name, '/route'),
_DrawerNavItem('Label', Icons.icon_name, '/route', badgeCount: count),
```

---

## Related Documentation

- **OneMind Codex:** [CODEX-FEATURE-ROADMAP.md](../../CODEX-FEATURE-ROADMAP.md)
- **PATHS System:** [PATHS-SYSTEM.md](./PATHS-SYSTEM.md) - The Three Paths (HP/LE/GE)
- **Skills System:** [SKILLS-SYSTEM.md](./SKILLS-SYSTEM.md)
- **Vehicles System:** [VEHICLES-SYSTEM.md](./VEHICLES-SYSTEM.md)

---

*"HUMAN (Zeus) + AI (Legacy) = ONE MIND"*
