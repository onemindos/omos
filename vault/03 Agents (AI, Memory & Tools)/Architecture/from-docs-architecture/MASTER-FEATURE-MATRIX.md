# OneMind OS Master Feature Matrix
## The Complete Vision: All Features, All Sources, One System

*Document Status: CANONICAL*
*Created: 2026-01-29*
*Purpose: Comprehensive feature mapping for the 24-hour build venture*

---

## Executive Summary

This document maps EVERY feature from:
- **Relevance AI** - Agent workforce, approvals, knowledge, chat
- **Pipedream** - Workflows, triggers, integrations
- **Motion** - Intelligent scheduling, DO DATE, calendar
- **Acorn Land Labs** - Gamification, simulation, agent characters
- **OneMind CODEX** - Paths (HP/LE/GE), intelligence scores, operating rhythm

Into ONE unified system with clear implementation status and priorities.

---

## Part 1: The Complete Feature Matrix

### Table Legend
| Symbol | Meaning |
|--------|---------|
| `[x]` HAVE | Feature exists and is functional |
| `[~]` PARTIAL | Feature exists but needs enhancement |
| `[-]` MISSING | Feature doesn't exist yet |
| `[!]` PRIORITY | Critical for 24hr venture |

---

### 1.1 CHAT INTERFACE

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Context-Aware Chat** | Motion, Relevance | `[~]` Basic chat exists | Full context injection with scores, schedule, state | `[!]` P0 |
| **@Agent Invocation** | Relevance | `[-]` Missing | `@Researcher`, `@Analyst`, `@Writer` routing | `[!]` P0 |
| **Multi-Model Support** | Relevance | `[x]` Have (Agno) | Claude, GPT-4, Gemini selectable per agent | P1 |
| **/Commands** | OneMind | `[-]` Missing | `/brief`, `/schedule`, `/scores`, `/delegate` | `[!]` P0 |
| **Web Browsing in Chat** | Motion, Relevance | `[~]` Have browser tools | Seamless web search with citations | P1 |
| **Document Generation** | Relevance | `[-]` Missing | Generate slides, docs, reports from chat | P2 |
| **Citations by Default** | Motion, Relevance | `[-]` Missing | All claims have sources, click to verify | `[!]` P0 |
| **Multi-Step Workflows** | Relevance, Pipedream | `[x]` Have workflows | Detect complex requests, show plan, execute | `[!]` P0 |
| **Chat History Search** | OneMind | `[~]` Have memory | Semantic search across all conversations | P1 |
| **Voice Chat** | OneMind | `[x]` Have LiveKit | Voice commands, hands-free operation | P2 |
| **Proactive Briefings** | Motion | `[-]` Missing | Morning/evening briefings pushed to chat | P1 |

---

### 1.2 AGENT WORKFORCE

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Agent Personas** | Relevance, Acorn | `[-]` Missing | Avatar, personality traits, communication style | `[!]` P0 |
| **Agent Factory** | Relevance | `[x]` Have | Dynamic agent creation with config | Done |
| **"Invent" Feature** | Relevance | `[-]` Missing | Natural language → working agent | P1 |
| **Agent Cloning** | Relevance | `[-]` Missing | Clone agent with modifications | P2 |
| **Multi-Agent Teams** | Relevance | `[x]` Have (research, dev, finance) | HP Team, LE Team, GE Team by path | P1 |
| **Agent-to-Agent Handoff** | Relevance | `[~]` Have team modes | Explicit handoff with context passing | P1 |
| **Workforce Dashboard** | Relevance, Acorn | `[-]` Missing | Visual grid of all agents with status | `[!]` P1 |
| **Task Queue Visibility** | Relevance | `[-]` Missing | See all pending/running tasks per agent | `[!]` P1 |
| **Agent Activity Logs** | Relevance | `[~]` Have logging | Timeline of agent actions with details | P2 |
| **Specialist Agents** | Agno | `[x]` Have 8 agents | researcher, analyst, writer, coder, finance, house, farm, security | Done |

---

### 1.3 APPROVALS & CONTROL (HITL)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Approval Mode** | Relevance | `[x]` Have | Task pauses for human sign-off | Done |
| **Agent Decides Mode** | Relevance | `[-]` Missing | Runs if confident, escalates if uncertain | `[!]` P0 |
| **Autorun Mode** | Relevance | `[x]` Have (auto_approve) | Full automation, no human needed | Done |
| **Confidence Scoring** | Relevance | `[-]` Missing | Each step outputs confidence (0-1) | P1 |
| **Risk Evaluation** | Relevance | `[-]` Missing | Evaluate risk before auto-approve | P1 |
| **Natural Lang Escalations** | Relevance | `[-]` Missing | "Escalate if involves money" | P1 |
| **Slack/Email Escalations** | Relevance, Pipedream | `[~]` Have notifications | Push approvals to Slack, email, mobile | P1 |
| **Human Feedback Training** | Relevance | `[-]` Missing | Learn from approval patterns | P3 |
| **Approval Timeout** | Agno | `[x]` Have | Auto-escalate after timeout | Done |

---

### 1.4 KNOWLEDGE & MEMORY

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Knowledge Bases** | Relevance | `[x]` Have (LanceDB) | RAG with embeddings | Done |
| **Train Once, Reuse All** | Relevance | `[x]` Have | Shared KB across agents | Done |
| **Auto-Sync Sources** | Relevance | `[-]` Missing | Notion, Google Drive, Supabase sync | P2 |
| **Live Data Fetch** | Motion, Relevance | `[~]` Have tools | Web search during conversations | P1 |
| **Metadata Extraction** | Relevance | `[-]` Missing | Extract entities, dates, amounts | P2 |
| **Unified Context** | OneMind | `[-]` Missing | AI knows EVERYTHING about user | `[!]` P0 |
| **Session Memory** | Agno | `[x]` Have | Per-session conversation history | Done |
| **Long-Term Memory** | OneMind | `[~]` Have DB storage | Cross-session learnings | P1 |
| **Entity Memory** | Agno | `[-]` Missing | Remember people, places, projects | P2 |

---

### 1.5 WORKFLOWS & AUTOMATION

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Workflow Engine** | Pipedream, Agno | `[x]` Have | Multi-step execution with context | Done |
| **Step Types** | Pipedream | `[x]` Have 9 types | step, parallel, condition, router, loop, approval, function, team, agent | Done |
| **HTTP Triggers** | Pipedream | `[-]` Missing | POST /workflows/trigger/{id} | `[!]` P0 |
| **Scheduled Triggers** | Pipedream | `[-]` Missing | Cron-based workflow triggers | P1 |
| **Event Triggers** | Pipedream | `[-]` Missing | On item create/update/delete | P1 |
| **Visual Workflow Builder** | Pipedream | `[-]` Missing | Drag-and-drop React Flow builder | P2 |
| **Workflow Templates** | Pipedream | `[-]` Missing | Pre-built workflow library | P2 |
| **Input/Output Mapping** | Pipedream | `[x]` Have | Data flows between steps | Done |
| **Error Handling** | Pipedream | `[~]` Have basic | Retry, fallback, DLQ | P1 |
| **Workflow Versioning** | Pipedream | `[-]` Missing | Version control for workflows | P3 |

---

### 1.6 SCHEDULING ENGINE (Motion-style)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **DO DATE ≠ DUE DATE** | Motion | `[-]` Missing | AI calculates optimal execution time | `[!]` P1 |
| **Intelligent Scheduling** | Motion | `[-]` Missing | Consider energy, duration, deps, priority | `[!]` P1 |
| **Auto-Reschedule** | Motion | `[-]` Missing | Dynamically adjust when things change | P1 |
| **Calendar Integration** | Motion | `[~]` Have CalDAV | Google Calendar sync, availability check | P1 |
| **Meeting Scheduling** | Motion | `[-]` Missing | AI finds optimal meeting times | P2 |
| **Deep Work Protection** | Motion | `[-]` Missing | Guard focus time blocks | P2 |
| **Capacity Planning** | Motion | `[-]` Missing | Warn when overcommitted | P1 |
| **Work Hours Control** | Motion, Relevance | `[-]` Missing | Define when agents can work | P2 |
| **Bulk Scheduler** | Relevance | `[-]` Missing | Schedule batch of items at once | P2 |
| **Wake Mode** | Relevance | `[-]` Missing | Proactive scheduling of follow-ups | P2 |

---

### 1.7 TASKS & PROJECTS (LifeOS)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Task Management** | Motion, LifeOS | `[x]` Have | CRUD tasks with priorities | Done |
| **Project Management** | Motion | `[~]` Have goals | Projects with milestones | P1 |
| **Universal Items** | CODEX | `[-]` Missing | Single items table for all types | `[!]` P0 |
| **Path-Based Organization** | CODEX | `[-]` Missing | HP, LE, GE categorization | `[!]` P0 |
| **Habit Tracking** | LifeOS | `[x]` Have | Daily habits with streaks | Done |
| **Goals & OKRs** | CODEX | `[x]` Have goals | Quarterly OKRs with key results | P1 |
| **Skills Progression** | CODEX | `[-]` Missing | Track skill development | P2 |
| **Inbox Capture** | Motion, LifeOS | `[x]` Have inbox | Quick capture with auto-routing | Done |
| **HIVE System** | CODEX | `[-]` Missing | Harvest→Integrate→Validate→Execute | P1 |

---

### 1.8 INTEGRATIONS

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Composio (2000+ apps)** | Agno | `[x]` Have | GitHub, Gmail, Slack, Linear, etc. | Done |
| **MCP Servers** | Agno | `[x]` Have | Local tool servers | Done |
| **Pipedream Connect** | Pipedream | `[-]` Optional | 3000+ additional apps | P3 |
| **OAuth Management** | Pipedream | `[~]` Have via Composio | Managed auth for all apps | P2 |
| **Biometrics (Oura)** | CODEX, Acorn | `[-]` Missing | Sleep, readiness, activity data | P1 |
| **Financial (Plaid)** | CODEX | `[-]` Missing | Bank accounts, transactions | P2 |
| **Home Automation** | CODEX, Acorn | `[x]` Have Home Assistant | Smart home control | Done |
| **Calendar Sync** | Motion | `[~]` Have CalDAV | Google, Outlook bidirectional | P1 |
| **Todoist Sync** | LifeOS | `[x]` Have | Bidirectional task sync | Done |

---

### 1.9 INTELLIGENCE SCORES (CODEX)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **HPI Score** | CODEX | `[-]` Missing | Holistic Performance Intelligence (0-100) | `[!]` P0 |
| **LEI Score** | CODEX | `[-]` Missing | Legacy Evolution Intelligence (0-100) | `[!]` P0 |
| **GEI Score** | CODEX | `[-]` Missing | Generational Endeavor Intelligence (0-100) | `[!]` P0 |
| **UI Score** | CODEX | `[-]` Missing | Sovereign Operations = (HPI+LEI+GEI)/3 | `[!]` P0 |
| **Score Components** | CODEX | `[-]` Missing | Breakdown by metric (sleep, habits, revenue, etc.) | P1 |
| **Score Trends** | CODEX | `[-]` Missing | Historical trend analysis (↑↓→) | P1 |
| **Score Thresholds** | CODEX | `[-]` Missing | Critical/Warning/Healthy/Optimal alerts | P1 |
| **Score-Based Insights** | CODEX | `[-]` Missing | AI recommendations based on scores | P1 |

---

### 1.10 OPERATING RHYTHM (CODEX)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Morning Briefing** | Motion, CODEX | `[-]` Missing | 5:30 AM daily briefing workflow | `[!]` P1 |
| **Evening Review** | CODEX | `[-]` Missing | 9:00 PM review and plan tomorrow | P1 |
| **Weekly Planning** | CODEX | `[-]` Missing | Sunday planning ritual | P1 |
| **Quarterly OKRs** | CODEX | `[-]` Missing | 12-week roadmap with objectives | P2 |
| **Annual Thesis** | CODEX | `[-]` Missing | Yearly direction and theme | P2 |
| **Path Balance Check** | CODEX | `[-]` Missing | Alert when paths are imbalanced | P1 |

---

### 1.11 GAMIFICATION (Acorn Labs)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Agent Characters** | Relevance, Acorn | `[-]` Missing | Visual avatars with animations | P2 |
| **XP System** | Acorn | `[-]` Missing | Earn XP for task completion | P2 |
| **Level Progression** | Acorn | `[-]` Missing | Path levels based on scores | P2 |
| **Achievement System** | Acorn | `[-]` Missing | Unlock badges for milestones | P3 |
| **Streak Tracking** | Acorn, LifeOS | `[x]` Have | Streak counters with fire emoji | Done |
| **Path Visualization** | Acorn | `[-]` Missing | Visual progress through HP/LE/GE | P2 |
| **Interactive Simulation** | Acorn | `[-]` Missing | Land Lab style visualization | P3 |

---

### 1.12 REPORTS & ANALYTICS

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Score Dashboard** | CODEX | `[-]` Missing | HPI/LEI/GEI/UI with charts | `[!]` P1 |
| **Activity Reports** | Motion | `[-]` Missing | What got done, time spent | P1 |
| **Path Progress** | CODEX | `[-]` Missing | Progress toward path goals | P1 |
| **Agent Performance** | Relevance | `[-]` Missing | Agent task completion rates | P2 |
| **Financial Overview** | CODEX | `[-]` Missing | Net worth, cash flow, runway | P2 |
| **Biometrics Dashboard** | CODEX | `[-]` Missing | Sleep, recovery, activity trends | P2 |
| **Custom Metrics** | OneMind | `[-]` Missing | Track anything with custom metrics | P2 |

---

### 1.13 CONSCIOUSNESS LAYER (Existing)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Event Bus (NATS)** | Consciousness | `[x]` Have | Real-time event streaming | Done |
| **State Manager** | Consciousness | `[x]` Have | Location, activity, devices, bio | Done |
| **Pattern Engine** | Consciousness | `[x]` Have | Pattern detection, insights | Done |
| **Temporal Awareness** | Consciousness | `[x]` Have | Time context, deadlines | Done |
| **Device Presence** | Consciousness | `[x]` Have | Device registry, activity | Done |
| **Execution Engine** | Consciousness | `[x]` Have | Task dispatch to agents/teams | Done |
| **KV Store** | Consciousness | `[x]` Have | NATS KV for state | Done |
| **Object Store** | Consciousness | `[x]` Have | NATS Object Store for files | Done |

---

### 1.14 CONTEXT ENGINE (NEW - Missing Piece)

| Feature | Source | Current State | Future State | Priority |
|---------|--------|---------------|--------------|----------|
| **Context Aggregation** | OneMind | `[-]` Missing | Aggregate ALL data sources into one | `[!]` P0 |
| **State Aggregator** | OneMind | `[-]` Missing | Wrap consciousness StateManager | `[!]` P0 |
| **Items Aggregator** | OneMind | `[-]` Missing | Query PostgreSQL for active items | `[!]` P0 |
| **Knowledge Aggregator** | OneMind | `[-]` Missing | RAG from LanceDB | `[!]` P0 |
| **Scores Aggregator** | OneMind | `[-]` Missing | Calculate HPI/LEI/GEI/UI | `[!]` P0 |
| **Calendar Aggregator** | OneMind | `[-]` Missing | Google Calendar integration | P1 |
| **History Aggregator** | OneMind | `[-]` Missing | Recent conversations | P1 |
| **Context Cache** | OneMind | `[-]` Missing | Redis-backed with TTL | P1 |
| **Context Builder** | OneMind | `[-]` Missing | Build prompt-ready context text | `[!]` P0 |

---

## Part 2: Feature Sources Summary

### What Comes From Each Platform

```
RELEVANCE AI (Agent Workforce)
═══════════════════════════════════════════════════════════════════════════════
├── Agent Personas (avatars, personality traits)
├── "Invent" Feature (natural language → agent)
├── Agent Cloning
├── Workforce Dashboard (visual agent grid)
├── Task Queue Visibility
├── Smart Approval Modes (approval, agent_decides, autorun)
├── Confidence Scoring
├── Natural Language Escalations
├── Knowledge Bases (train once, reuse all)
├── @Agent Invocation in Chat
├── Bulk Scheduler
├── Wake Mode (proactive scheduling)
└── Human Feedback Training

PIPEDREAM (Workflows & Integrations)
═══════════════════════════════════════════════════════════════════════════════
├── HTTP Triggers (webhook endpoints)
├── Scheduled Triggers (cron)
├── Event Triggers (on data change)
├── Visual Workflow Builder
├── Step Types (code, action, condition)
├── Workflow Templates
├── Input/Output Mapping
├── Error Handling (retry, fallback)
├── 3000+ App Integrations
└── VPC (static IPs)

MOTION (Intelligent Scheduling)
═══════════════════════════════════════════════════════════════════════════════
├── DO DATE ≠ DUE DATE Algorithm
├── Intelligent Scheduling (energy, duration, deps)
├── Auto-Reschedule
├── Calendar Integration
├── Meeting Scheduling
├── Deep Work Protection
├── Capacity Planning
├── Context-Aware Chat
├── Multi-Step Workflow Execution
├── Citations by Default
└── Reports Dashboard

ACORN LAND LABS (Gamification)
═══════════════════════════════════════════════════════════════════════════════
├── Agent Characters (visual avatars)
├── Interactive Simulation View
├── XP/Level Progression
├── Achievement System
├── Streak Tracking
├── Visual Progress Dashboard
├── Sensor Data Integration
└── AI-Powered Suggestions

ONEMIND CODEX (Human Optimization)
═══════════════════════════════════════════════════════════════════════════════
├── Three Paths (HP, LE, GE)
├── Intelligence Scores (HPI, LEI, GEI, UI)
├── Score Components (weighted metrics)
├── Score Thresholds (Critical/Warning/Healthy/Optimal)
├── Universal Items (unified data model)
├── HIVE Capture System
├── Operating Rhythm (A→Q→W→D)
├── Morning/Evening Briefings
├── Weekly/Quarterly/Annual Planning
├── Path Balance Monitoring
├── Finance Module
├── Relationships Module
├── Biometrics Dashboard
├── Home Management
├── Family Module
├── Heritage Module
├── Business Dashboard
└── Four-Agent Personas (HP/LE/GE specialists)
```

---

## Part 3: The Unified System Architecture

```
ONEMIND UNIFIED ARCHITECTURE
═══════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│                           USER INTERFACES                                   │
│                                                                             │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│   │    CHAT      │  │    VOICE     │  │   MOBILE     │  │    WEB       │  │
│   │   (UIO)      │  │  (LiveKit)   │  │  (Flutter)   │  │ (Dashboard)  │  │
│   │              │  │              │  │              │  │              │  │
│   │ @agents      │  │ "Hey Legacy" │  │ Quick add    │  │ Command      │  │
│   │ /commands    │  │ Voice cmds   │  │ Widgets      │  │ Center       │  │
│   │ Workflows    │  │ Hands-free   │  │ Notifications│  │ Reports      │  │
│   └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  │
│          └─────────────────┴─────────────────┴─────────────────┘          │
│                                      │                                     │
│                                      ▼                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         CONTEXT ENGINE (NEW)                                │
│                    "The AI Awareness Layer"                                 │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                                                                     │  │
│   │  Aggregators:                                                       │  │
│   │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ │  │
│   │  │ State  │ │ Items  │ │Knowledge│ │ Scores │ │Calendar│ │History │ │  │
│   │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ │  │
│   │                                                                     │  │
│   │  Context Builder: build_full_context() → All data for AI prompts   │  │
│   │  Context Cache: Redis-backed, event-invalidated                    │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                      │                                     │
│                                      ▼                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         AGENT WORKFORCE (Relevance)                         │
│                                                                             │
│   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐            │
│   │ LEGACY  │ │RESEARCHER│ │ ANALYST │ │ WRITER  │ │ CODER   │            │
│   │   🧠    │ │   🔍    │ │   📊    │ │   ✍️    │ │   💻    │            │
│   │ Primary │ │ Research │ │ Analyze │ │ Draft   │ │ Code    │            │
│   └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘            │
│                                                                             │
│   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐            │
│   │ FINANCE │ │HOUSE MGR│ │ FARM MGR│ │SECURITY │ │HP COACH │            │
│   │   💰    │ │   🏠    │ │   🌾    │ │   🛡️    │ │   🏃    │            │
│   └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘            │
│                                                                             │
│   Approval Modes: [APPROVAL] [AGENT_DECIDES] [AUTORUN]                     │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                        WORKFLOW ENGINE (Pipedream)                          │
│                                                                             │
│   Triggers: [HTTP] [SCHEDULE] [EVENT] [CHAT_COMMAND]                       │
│                                                                             │
│   Steps: [AGENT]→[PARALLEL]→[CONDITION]→[APPROVAL]→[FUNCTION]→[ACTION]    │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                       SCHEDULING ENGINE (Motion)                            │
│                                                                             │
│   DO DATE Calculator: Energy + Duration + Deps + Calendar → Optimal Time   │
│   Auto-Reschedule: Dynamically adjust when changes occur                   │
│   Capacity Planning: Warn when overcommitted                               │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                       INTELLIGENCE ENGINE (CODEX)                           │
│                                                                             │
│   ┌──────────────────────────────────────────────────────────────────────┐ │
│   │                    UNIFIED INTELLIGENCE: 76                           │ │
│   │                                                                      │ │
│   │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐        │ │
│   │  │ HPI: 81 🟢 ↑  │  │ LEI: 65 🟡 →  │  │ GEI: 82 🟢 ↑  │        │ │
│   │  │ Sleep: 85      │  │ Rituals: 70    │  │ Runway: 85     │        │ │
│   │  │ Recovery: 78   │  │ Home: 75       │  │ FCF: 80        │        │ │
│   │  │ Habits: 90     │  │ Time: 55 ⚠️    │  │ Revenue: 85    │        │ │
│   │  │ Training: 70   │  │ Wealth: 70     │  │ Clients: 78    │        │ │
│   │  │ Mood: 80       │  │ Comms: 65      │  │ Legal: 80      │        │ │
│   │  └────────────────┘  └────────────────┘  └────────────────┘        │ │
│   └──────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                        INTEGRATION LAYER                                    │
│                                                                             │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│   │ Composio │  │   MCP    │  │   Oura   │  │  Plaid   │  │Home Asst │  │
│   │ 2000+apps│  │ Servers  │  │Biometrics│  │ Finance  │  │  IoT     │  │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         DATA LAYER                                          │
│                                                                             │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│   │PostgreSQL│  │ LanceDB  │  │  Redis   │  │   NATS   │  │TimescaleDB│  │
│   │ Items    │  │Knowledge │  │  Cache   │  │  Events  │  │Historical │  │
│   │ Scores   │  │Embeddings│  │ Sessions │  │  Pubsub  │  │ Archive  │  │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Part 4: UI/Menu Structure for New System

### 4.1 Navigation Redesign

The current navigation needs to be restructured to accommodate the new features:

```
CURRENT NAVIGATION                    NEW NAVIGATION
════════════════════                  ═══════════════════════════════════════

🏠 Home                        →     🧠 COMMAND CENTER (Dashboard)
📥 Inbox                       →     │   ├── Morning Briefing Panel
✅ Tasks                       →     │   ├── Intelligence Scores Widget
📅 Calendar                    →     │   ├── Path Balance Indicator
🎯 Goals                       →     │   ├── Active Focus Mode
📒 Journal                     →     │   ├── Quick Actions Bar
🔁 Habits                      →     │   └── System Health Panel
🍽️ Meals                       →     │
⏰ Routines                    →     💬 CHAT (Primary Interface)
🍅 Pomodoro                    →     │   ├── UIO Chat (main)
⚙️ Settings                    →     │   ├── @Agent Routing
                                     │   ├── /Commands
                                     │   └── Voice Mode Toggle
                                     │
                                     🤖 WORKFORCE
                                     │   ├── Agent Grid (all agents)
                                     │   ├── Task Queue
                                     │   ├── Pending Approvals
                                     │   └── + Create Agent
                                     │
                                     🛤️ PATHS
                                     │   ├── HP - Holistic Performance
                                     │   │   ├── Tasks (HP)
                                     │   │   ├── Goals (HP)
                                     │   │   ├── Habits (HP)
                                     │   │   └── Biometrics
                                     │   │
                                     │   ├── LE - Legacy Evolution
                                     │   │   ├── Tasks (LE)
                                     │   │   ├── Goals (LE)
                                     │   │   ├── Family
                                     │   │   └── Home
                                     │   │
                                     │   └── GE - Generational Endeavor
                                     │       ├── Tasks (GE)
                                     │       ├── Goals (GE)
                                     │       ├── Business
                                     │       └── Finance
                                     │
                                     📅 SCHEDULE
                                     │   ├── Calendar View
                                     │   ├── Today's Schedule
                                     │   └── Upcoming Deadlines
                                     │
                                     🔄 RHYTHMS
                                     │   ├── Daily (Habits, Routines)
                                     │   ├── Weekly Plan
                                     │   ├── Quarterly OKRs
                                     │   └── Annual Thesis
                                     │
                                     ⚡ WORKFLOWS
                                     │   ├── Active Workflows
                                     │   ├── Workflow Builder
                                     │   └── Templates
                                     │
                                     📊 REPORTS
                                     │   ├── Scores Dashboard
                                     │   ├── Activity Summary
                                     │   ├── Path Progress
                                     │   └── Financial Overview
                                     │
                                     📥 CAPTURE
                                     │   ├── Inbox
                                     │   ├── Quick Add
                                     │   └── Voice Capture
                                     │
                                     🔧 SETTINGS
                                         ├── Profile
                                         ├── Integrations
                                         ├── Notifications
                                         └── Agent Config
```

### 4.2 Primary Screen Layouts

#### Command Center (Home)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ONEMIND COMMAND CENTER                              [Zeus] [🔔 3] [⚙️]     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │   MORNING BRIEFING          │  │       INTELLIGENCE SCORES           │  │
│  │   Wednesday, Jan 29, 2026   │  │                                     │  │
│  │   ☀️ 68°F Sunny              │  │  ┌────────┐ ┌────────┐ ┌────────┐  │  │
│  │                             │  │  │  HPI   │ │  LEI   │ │  GEI   │  │  │
│  │   😴 Sleep: 7h 23m (78)     │  │  │  81 🟢 │ │  65 🟡 │ │  82 🟢 │  │  │
│  │   ⚡ Energy: Moderate       │  │  │   ↑    │ │   →    │ │   ↑    │  │  │
│  │                             │  │  └────────┘ └────────┘ └────────┘  │  │
│  │   📋 TOP 3 TODAY:           │  │                                     │  │
│  │   1. □ Client proposal [GE] │  │   UNIFIED SCORE: 76 🟡              │  │
│  │   2. □ Team standup [GE]    │  │   "Focus on family time (LEI)"     │  │
│  │   3. □ Workout [HP]         │  │                                     │  │
│  │                             │  │  ┌──────────────────────────────┐   │  │
│  │   🔔 PENDING APPROVALS: 2   │  │  │  PATH BALANCE (This Week)    │   │  │
│  └─────────────────────────────┘  │  │  HP █████████░░░  45%        │   │  │
│                                   │  │  LE ████░░░░░░░░░  20% ⚠️    │   │  │
│  ┌─────────────────────────────┐  │  │  GE ███████░░░░░  35%        │   │  │
│  │   ACTIVE FOCUS              │  │  └──────────────────────────────┘   │  │
│  │   ┌───────────────────────┐ │  └─────────────────────────────────────┘  │
│  │   │ Client Proposal       │ │                                           │
│  │   │      45:23            │ │  ┌─────────────────────────────────────┐  │
│  │   │   [⏸️ Pause] [✓ Done] │ │  │       QUICK ACTIONS                 │  │
│  │   └───────────────────────┘ │  │                                     │  │
│  │                             │  │  [+ Task] [+ Inbox] [🎙️ Voice]      │  │
│  │   [Start New Focus]         │  │  [💬 Chat with Legacy]              │  │
│  └─────────────────────────────┘  └─────────────────────────────────────┘  │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  SYSTEM STATUS                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ ✅ Agents: 8 Online | ✅ Todoist: Synced (2m ago) | 🔔 Inbox: 3     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Workforce Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  WORKFORCE                                          [+ New Agent] [⚙️]      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  MY AGENTS                                                                  │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│  │      🧠         │  │      🔍         │  │      📊         │            │
│  │     LEGACY      │  │   RESEARCHER    │  │    ANALYST      │            │
│  │                 │  │                 │  │                 │            │
│  │  ████████░░     │  │  ████░░░░░░     │  │  ██████████     │            │
│  │  Active: 3      │  │  Active: 1      │  │  Idle           │            │
│  │                 │  │                 │  │                 │            │
│  │  "Helping you   │  │  "Researching   │  │  "Ready for     │            │
│  │   optimize"     │  │   marathon..."  │  │   analysis"     │            │
│  │                 │  │                 │  │                 │            │
│  │ [Chat] [Tasks]  │  │ [Chat] [Tasks]  │  │ [Chat] [Tasks]  │            │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘            │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│  │      ✍️         │  │      💻         │  │      💰         │            │
│  │     WRITER      │  │     CODER       │  │    FINANCE      │            │
│  │                 │  │                 │  │                 │            │
│  │  ██████░░░░     │  │  ░░░░░░░░░░     │  │  ████████░░     │            │
│  │  Active: 1      │  │  Idle           │  │  Active: 2      │            │
│  │                 │  │                 │  │                 │            │
│  │  "Drafting      │  │  "Ready for     │  │  "Analyzing     │            │
│  │   blog post"    │  │   code review"  │  │   Q4 cash..."   │            │
│  │                 │  │                 │  │                 │            │
│  │ [Chat] [Tasks]  │  │ [Chat] [Tasks]  │  │ [Chat] [Tasks]  │            │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘            │
│                                                                             │
│  ─────────────────────────────────────────────────────────────────────────  │
│                                                                             │
│  TASK QUEUE                                                 Total: 7 tasks │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │ Agent      │ Task                              │ Status   │ ETA      │  │
│  ├────────────┼───────────────────────────────────┼──────────┼──────────┤  │
│  │ Researcher │ Marathon training plan research   │ Running  │ ~5 min   │  │
│  │ Writer     │ Blog post draft                   │ Running  │ ~15 min  │  │
│  │ Finance    │ Q4 cash flow analysis             │ Running  │ ~10 min  │  │
│  │ Finance    │ Tax deduction research            │ Queued   │ Next     │  │
│  │ Legacy     │ Schedule family dinner            │ Running  │ ~2 min   │  │
│  │ Legacy     │ Book dentist appointments         │ Queued   │ Next     │  │
│  │ Legacy     │ Review investor meeting           │ Pending  │ Approval │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  PENDING APPROVALS                                              2 pending  │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │ 🔔 Review investor meeting request              [Approve] [Reject]   │  │
│  │ 🔔 Send automated email to clients              [Approve] [Reject]   │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Path View (HP Example)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  HP - HOLISTIC PERFORMANCE                                [🏃 HP Coach]     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  SCORE: HPI 81 🟢                           Level: ATHLETE         │   │
│  │  ████████████████████████░░░░░░░░░░  81%    XP: 12,450 / 15,000    │   │
│  │                                                                     │   │
│  │  Components:                                                        │   │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐│   │
│  │  │ SLEEP  │ │RECOVERY│ │ HABITS │ │TRAINING│ │  MOOD  │ │ GOALS  ││   │
│  │  │  85 🟢 │ │  78 🟡 │ │  90 🟢 │ │  70 🟡 │ │  80 🟢 │ │  75 🟡 ││   │
│  │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘│   │
│  │                                                                     │   │
│  │  🔥 STREAKS: Morning Run 23d | Meditation 45d | Sleep 10pm 15d     │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │   HP TASKS                  │  │   HP GOALS                          │  │
│  │                             │  │                                     │  │
│  │   □ Morning run (10 miles)  │  │   🎯 Run a marathon (May)          │  │
│  │   □ Strength training       │  │      ██████░░░░░░░░ 40%            │  │
│  │   □ Meal prep for week      │  │                                     │  │
│  │   ✓ Meditation (done)       │  │   🎯 Sleep 7.5h average           │  │
│  │   ✓ Supplements (done)      │  │      ██████████░░░░ 75%            │  │
│  │                             │  │                                     │  │
│  │   [View All HP Tasks]       │  │   [View All HP Goals]              │  │
│  └─────────────────────────────┘  └─────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │   HP HABITS                 │  │   BIOMETRICS                        │  │
│  │                             │  │                                     │  │
│  │   🔥 Morning Run      23d   │  │   😴 Sleep: 7h 23m (78)            │  │
│  │   🔥 Meditation       45d   │  │   💓 Resting HR: 52 bpm            │  │
│  │   🔥 Sleep by 10pm    15d   │  │   📈 HRV: 45ms                     │  │
│  │   ○ Reading           0d    │  │   🏃 Steps: 8,234                  │  │
│  │                             │  │   🔥 Active Cal: 520               │  │
│  │   [View All HP Habits]      │  │                                     │  │
│  └─────────────────────────────┘  └─────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Part 5: 24-Hour Build Prioritization

### Critical Path (Must Complete)

```
PHASE 1: CONTEXT ENGINE (Hours 1-4)
═══════════════════════════════════════════════════════════════════════════════
Priority: P0 - Critical Foundation

Files to Create:
├── backend/context/__init__.py
├── backend/context/engine.py          # Main context engine
├── backend/context/builder.py         # Context builders
├── backend/context/cache.py           # Redis cache
├── backend/context/models.py          # Pydantic models
├── backend/context/aggregators/
│   ├── __init__.py
│   ├── state.py                       # Wrap consciousness StateManager
│   ├── items.py                       # PostgreSQL items queries
│   ├── knowledge.py                   # LanceDB RAG
│   ├── scores.py                      # Calculate HPI/LEI/GEI/UI
│   ├── calendar.py                    # Google Calendar
│   └── history.py                     # Conversation history

Deliverables:
✓ ContextEngine.get_full_context(user_id) → FullContext
✓ ContextEngine.build_prompt_context(user_id, query) → str
✓ ScoresAggregator.calculate_all(user_id) → IntelligenceScores
✓ Context caching with Redis and TTL


PHASE 2: ENHANCED CHAT (Hours 5-8)
═══════════════════════════════════════════════════════════════════════════════
Priority: P0 - Primary Interface

Files to Create:
├── backend/agno/chat/__init__.py
├── backend/agno/chat/enhanced.py      # Main chat engine
├── backend/agno/chat/commands.py      # /command handlers
├── backend/agno/chat/agents.py        # @agent router
├── backend/agno/chat/citations.py     # Citation tracking

Deliverables:
✓ EnhancedChatEngine.process_message() with context injection
✓ @agent syntax parsing and routing
✓ /command handling (brief, schedule, scores, delegate)
✓ Citation tracking for all claims


PHASE 3: UNIVERSAL ITEMS & SCORES (Hours 9-12)
═══════════════════════════════════════════════════════════════════════════════
Priority: P0 - Data Foundation

Database Migration:
├── migrations/20260129_universal_items.sql
│   ├── CREATE TYPE item_type AS ENUM (...)
│   ├── CREATE TYPE path_type AS ENUM ('HP', 'LE', 'GE')
│   ├── CREATE TABLE intelligence_scores (...)
│   └── ALTER TABLE items ADD COLUMN path path_type

Files to Create:
├── backend/lifeos/scores/__init__.py
├── backend/lifeos/scores/service.py   # Score CRUD
├── backend/lifeos/scores/calculator.py # Calculation logic
├── backend/lifeos/scores/api.py       # REST endpoints

Deliverables:
✓ Intelligence scores table with history
✓ Path categorization on all items
✓ GET /api/scores/current endpoint
✓ Score calculation triggered on data changes


PHASE 4: AGENT ENHANCEMENTS (Hours 13-16)
═══════════════════════════════════════════════════════════════════════════════
Priority: P1 - Workforce Features

Files to Modify:
├── backend/agno/agents/models.py      # Add persona fields

Files to Create:
├── backend/agno/agents/personas.py    # Persona management
├── backend/agno/workflows/triggers/http.py  # HTTP triggers

Database Migration:
├── ALTER TABLE agents ADD COLUMN persona JSONB DEFAULT '{}'

Deliverables:
✓ Agent personas (avatar, emoji, traits, style)
✓ HTTP workflow triggers: POST /workflows/trigger/{id}
✓ Smart approval modes (agent_decides)


PHASE 5: SCHEDULING & BRIEFINGS (Hours 17-20)
═══════════════════════════════════════════════════════════════════════════════
Priority: P1 - Intelligence Features

Files to Create:
├── backend/platform/scheduling/__init__.py
├── backend/platform/scheduling/intelligent.py  # DO DATE algorithm
├── backend/platform/scheduling/briefings.py    # Morning/evening

Deliverables:
✓ DO DATE calculation algorithm
✓ Morning briefing generator
✓ Evening review generator
✓ Briefing push to chat


PHASE 6: FRONTEND INTEGRATION (Hours 21-24)
═══════════════════════════════════════════════════════════════════════════════
Priority: P1 - User Experience

Files to Create:
├── frontend/lib/context/context_service.dart
├── frontend/lib/context/context_provider.dart
├── frontend/lib/shared/widgets/scores_card.dart
├── frontend/lib/shared/widgets/command_center.dart

Deliverables:
✓ Context service connecting to backend
✓ Scores display widget
✓ Updated navigation structure
✓ Command center dashboard
```

---

## Part 6: Linear Project Structure

### Epic 1: Context Engine Foundation
- **Issue 1.1**: Create Context Engine module structure
- **Issue 1.2**: Implement State Aggregator (wrap consciousness)
- **Issue 1.3**: Implement Items Aggregator (PostgreSQL queries)
- **Issue 1.4**: Implement Scores Aggregator (HPI/LEI/GEI/UI)
- **Issue 1.5**: Implement Knowledge Aggregator (LanceDB RAG)
- **Issue 1.6**: Implement Context Builder and Cache
- **Issue 1.7**: Add context API endpoints

### Epic 2: Enhanced Chat System
- **Issue 2.1**: Create Enhanced Chat Engine
- **Issue 2.2**: Implement @agent syntax parsing and routing
- **Issue 2.3**: Implement /command handlers
- **Issue 2.4**: Add citation tracking system
- **Issue 2.5**: Integrate context injection into chat

### Epic 3: Intelligence Scores
- **Issue 3.1**: Create database schema for scores
- **Issue 3.2**: Implement score calculation service
- **Issue 3.3**: Create scores API endpoints
- **Issue 3.4**: Add path categorization to items

### Epic 4: Agent Workforce Enhancement
- **Issue 4.1**: Add persona fields to agent model
- **Issue 4.2**: Implement HTTP workflow triggers
- **Issue 4.3**: Add smart approval modes (agent_decides)
- **Issue 4.4**: Create task queue visibility API

### Epic 5: Intelligent Scheduling
- **Issue 5.1**: Implement DO DATE algorithm
- **Issue 5.2**: Create morning briefing generator
- **Issue 5.3**: Create evening review generator
- **Issue 5.4**: Add calendar integration

### Epic 6: Frontend Integration
- **Issue 6.1**: Create context service and provider
- **Issue 6.2**: Build scores display widget
- **Issue 6.3**: Update navigation structure
- **Issue 6.4**: Build command center dashboard

---

## Part 7: Technology Stack Summary

```
TECHNOLOGY STACK
═══════════════════════════════════════════════════════════════════════════════

LAYER                     TECHNOLOGY          PURPOSE
────────────────────────────────────────────────────────────────────────────────

Frontend                  Flutter 3.6+        Mobile/Web/Desktop apps
                          Dart                Frontend language
                          Riverpod            State management
                          React Flow          Workflow builder (future)

Backend Framework         Python 3.12+        Backend language
                          FastAPI             REST API framework
                          Pydantic            Data validation

AI/Agents                 Agno                Agent framework
                          LangChain/Phidata   LLM orchestration
                          Claude Opus 4.5     Primary LLM (AWS Bedrock)
                          GPT-4o              Alternative LLM
                          LanceDB             Vector embeddings/RAG

Data Storage              PostgreSQL 15+      Primary database
                          Redis               Cache, sessions, hot state
                          LanceDB             Knowledge base embeddings
                          TimescaleDB         Historical time-series

Event Infrastructure      NATS JetStream      Event bus, pub/sub
                          NATS KV             State storage
                          NATS Object Store   File/memory storage

Integrations              Composio            2000+ app integrations
                          MCP Servers         Local tool servers
                          Pipedream           Optional 3000+ apps

Voice                     LiveKit Agents      Voice interface

External APIs             Oura API            Biometrics
                          Plaid API           Financial data
                          Google Calendar     Calendar sync
                          Home Assistant      Smart home

Infrastructure            Docker              Containerization
                          Cloudflare          CDN, DNS, Tunnels
                          AWS Bedrock         LLM API access
```

---

## Conclusion

This document represents the COMPLETE feature vision for OneMind OS, combining:

1. **Relevance AI**: Agent workforce, personas, approvals, knowledge
2. **Pipedream**: Workflows, triggers, integrations
3. **Motion**: Intelligent scheduling, DO DATE, calendar
4. **Acorn Labs**: Gamification, visualization, agent characters
5. **OneMind CODEX**: Paths, scores, operating rhythm, universal items

The 24-hour venture focuses on the **critical foundation**:
- Context Engine (makes AI aware)
- Enhanced Chat (primary interface)
- Intelligence Scores (differentiation)
- Agent Enhancements (workforce features)

Everything else builds on this foundation in subsequent phases.

**Human + AI = ONE MIND** - This is the complete vision.

---

*Document Version: 1.0*
*Created: 2026-01-29*
*Author: Claude Opus 4.5 + Zeus*
