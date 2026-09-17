# OneMind OS Feature Roadmap
## Derived from the OneMind Codex

> **Document Status:** CANONICAL
> **Created:** 2026-01-28
> **Source:** OneMind Codex (https://github.com/Zeus-Delacruz/OneMind-Codex)
> **Purpose:** Bridge the Codex philosophy with OneMind OS implementation

---

## Executive Summary

The **OneMind Codex** defines the philosophical framework and operational doctrine for a unified human enhancement framework. **OneMind OS** is the technical implementation. This document maps Codex concepts to concrete features that need to be built.

### The Equation
```
HUMAN (Zeus) + AI (Legacy) = ONE MIND
```

### The Three Paths + Sovereign Operations
```
┌─────────────────────────────────────────────────────────────┐
│                UI (00-24) - UNIFIED INTELLIGENCE             │
│                        (The Brain)                           │
│   Coordinates all paths, manages tech infrastructure         │
│              UI Score = (HPI + LEI + GEI) / 3               │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   HP (25-49)  │     │   LE (50-74)  │     │   GE (75-99)  │
│   HOLISTIC    │     │    LEGACY     │     │  GENERATIONAL │
│  PERFORMANCE  │     │   EVOLUTION   │     │ENTREPRENEURSHP│
│               │     │               │     │               │
│  Self Path    │     │  Legacy Path  │     │  Wealth Path  │
│ "Improve ME"  │     │ "Protect US"  │     │"Create VALUE" │
│  HPI Score    │     │  LEI Score    │     │  GEI Score    │
└───────────────┘     └───────────────┘     └───────────────┘
```

---

## Current State Analysis

### What's Built (OneMind OS)

| Component | Status | Lines of Code |
|-----------|--------|---------------|
| Tasks Module | ✅ Complete | ~2,500 |
| Habits Module | ✅ Complete | ~3,200 |
| Routines Module | ✅ Complete | ~1,800 |
| Goals Module | ✅ Complete | ~2,100 |
| Journal Module | ✅ Complete | ~1,900 |
| Pomodoro Module | ✅ Complete | ~1,500 |
| Calendar Module | ✅ Complete | ~2,800 |
| Meals Module | ✅ Complete | ~3,500 |
| Scheduling Module | ✅ Complete | ~2,200 |
| Planner Module | ✅ Complete | ~2,400 |
| Activity Module | ✅ Complete | ~1,600 |
| Inbox Module | ✅ Complete | ~1,400 |
| Notifications Module | ✅ Complete | ~1,200 |
| Sync Engine | ✅ Complete | ~2,800 |
| **Total Backend** | ✅ Complete | **~15,965** |
| **Total Frontend** | ✅ Complete | **~28,891** |
| **Grand Total** | ✅ Complete | **~44,856** |

### What's Missing (From Codex)

| Component | Codex Reference | Priority | Effort |
|-----------|-----------------|----------|--------|
| Intelligence Scores (HPI/LEI/GEI) | Core Doctrine | P0 - Critical | Medium |
| Command Center Dashboard | UI-01 | P0 - Critical | Large |
| Operating Rhythm (A→Q→W→D) | Part 5 | P0 - Critical | Medium |
| Finance Module | HP-30 | P1 - High | Large |
| Relationships Module | HP-32 | P1 - High | Medium |
| Biometrics Dashboard | HP-27 | P1 - High | Medium |
| Home Management | LE-50 | P2 - Medium | Medium |
| Family Module | LE-51/52/53 | P2 - Medium | Large |
| Heritage Module | LE-54/56/57 | P2 - Medium | Medium |
| Personal Brand | GE-75 | P3 - Low | Medium |
| Business Dashboard | GE-75-85 | P3 - Low | Large |
| Four-Agent Personas | Part 4 | P3 - Low | Medium |

---

## Phase 1: Core Codex Alignment (Q1 2026)

### 1.1 Intelligence Scores System (HPI/LEI/GEI)

**Codex Reference:** Part 3 - The Three Paths, Part 5 - Unified Score

The Intelligence Scores are the heart of the OneMind system. Every path generates a score (0-100) that measures progress and health.

#### Score Definitions

| Score | Name | Path | What It Measures |
|-------|------|------|------------------|
| **HPI** | Holistic Performance Intelligence | HP (25-49) | Personal optimization, health, skills, growth |
| **LEI** | Legacy Evolution Intelligence | LE (50-74) | Family protection, heritage, generational wealth |
| **GEI** | Generational Endeavor Intelligence | GE (75-99) | Value creation, business success, revenue |

**UI (Sovereign Operations) = (HPI + LEI + GEI) / 3**

#### HPI Score Components (Holistic Performance Intelligence)

| Metric | Weight | Source | Calculation |
|--------|--------|--------|-------------|
| Sleep Quality | 20% | Oura/Manual | 7-day average sleep score |
| Recovery/Readiness | 20% | Oura/Manual | 7-day average readiness |
| Habit Completion | 25% | Habits Module | % habits completed this week |
| Training Load | 15% | Activity/Manual | Exercise frequency vs goal |
| Mood Average | 10% | Journal/Habits | 7-day mood average |
| Goal Progress | 10% | Goals Module | Active goals on track % |

#### LEI Score Components (Legacy Evolution Intelligence)

| Metric | Weight | Source | Calculation |
|--------|--------|--------|-------------|
| Family Rituals | 25% | Habits/Calendar | Weekly family activities completed |
| Home Maintenance | 20% | Tasks/Home Module | Maintenance tasks on schedule |
| Quality Time | 20% | Calendar/Activity | Hours with family this week |
| Wealth Posture | 15% | Finance Module | Net worth delta (positive = good) |
| Communication | 10% | Relationships | Family check-ins completed |
| Legacy Projects | 10% | Goals Module | Heritage projects progress |

#### GEI Score Components (Generational Endeavor Intelligence)

| Metric | Weight | Source | Calculation |
|--------|--------|--------|-------------|
| Runway | 25% | Finance Module | Months of expenses covered |
| Free Cash Flow | 25% | Finance Module | Income - expenses (positive %) |
| Revenue Growth | 20% | Business Module | MoM revenue change |
| Client Health | 15% | Relationships | Active client satisfaction |
| Legal Hygiene | 10% | Tasks | Compliance tasks current |
| Capital Efficiency | 5% | Finance | ROI on investments |

#### Score Thresholds

| Range | Status | Color | Action Required |
|-------|--------|-------|-----------------|
| 0-59 | Critical | 🔴 Red | Simplify - drop non-essentials |
| 60-79 | Warning | 🟡 Amber | Stabilize - focus on basics |
| 80-89 | Healthy | 🟢 Green | Maintain - continue rhythm |
| 90-100 | Optimal | 🔵 Blue | Scale - can take on more |

#### Backend Implementation

```python
# backend/lifeos/scores/models.py
class IntelligenceScore(BaseModel):
    id: UUID
    user_id: str
    score_type: Literal["HPI", "LEI", "GEI", "UI"]
    score: float  # 0-100
    components: dict  # Breakdown of contributing metrics
    calculated_at: datetime
    period_start: date
    period_end: date

# backend/lifeos/scores/service.py
class ScoreService:
    async def calculate_hpi_score(self, user_id: str, period: date) -> float
    async def calculate_lei_score(self, user_id: str, period: date) -> float
    async def calculate_gei_score(self, user_id: str, period: date) -> float
    async def calculate_ui_score(self, user_id: str, period: date) -> float
    async def get_score_history(self, user_id: str, score_type: str, days: int) -> list
    async def get_score_trends(self, user_id: str) -> dict

# backend/lifeos/scores/api.py
# GET /api/scores/current - Get current week's scores
# GET /api/scores/history?type=HPI&days=30 - Get score history
# GET /api/scores/trends - Get trend analysis
# POST /api/scores/calculate - Trigger score recalculation
```

#### Frontend Implementation

```dart
// frontend/lib/lifeos/scores/
├── models/
│   └── intelligence_score.dart
├── providers/
│   └── scores_provider.dart
├── screens/
│   ├── scores_dashboard_screen.dart
│   └── score_detail_screen.dart
└── widgets/
    ├── score_card.dart
    ├── score_gauge.dart
    ├── score_trend_chart.dart
    └── unified_score_widget.dart
```

#### Database Schema

```sql
CREATE TABLE intelligence_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) NOT NULL,
    score_type VARCHAR(20) NOT NULL CHECK (score_type IN ('HPI', 'LEI', 'GEI', 'UI')),
    score DECIMAL(5,2) NOT NULL CHECK (score >= 0 AND score <= 100),
    components JSONB NOT NULL DEFAULT '{}',
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    calculated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, score_type, period_start)
);

CREATE INDEX idx_scores_user_type ON intelligence_scores(user_id, score_type);
CREATE INDEX idx_scores_period ON intelligence_scores(period_start, period_end);
```

---

### 1.2 Command Center Dashboard (UI-01)

**Codex Reference:** 01 Command (Dashboard & Navigation)

The Command Center is the human interface to OneMind - a unified dashboard that provides situational awareness and control.

#### Dashboard Sections

##### 1. Morning Briefing Panel
- Today's date, weather, energy level
- Calendar summary (next 3 events)
- Top 3 priorities (from planner)
- Pending approvals from agents
- Streak status for key habits

##### 2. Path Balance Indicator
- Visual pie/bar showing time allocation across HP/LE/GE
- Target vs actual for the week
- Alerts if any path is neglected (< 20% allocation)

##### 3. Intelligence Scores Widget
- Three score cards (HPI, LEI, GEI)
- UI (Sovereign Operations) score prominently displayed
- Trend arrows (↑↓→)
- Click to see breakdown

##### 4. Active Focus Mode
- Current task with timer
- Pomodoro integration
- "I'm working on..." status
- Distraction blocker trigger

##### 5. Quick Actions Bar
- Add task (with path routing)
- Add to inbox
- Log activity
- Start pomodoro
- Ask Legacy

##### 6. System Health Panel
- Agent status (online/offline)
- Sync status (last sync time)
- API health indicators
- Notification queue count

#### Wireframe

```
┌─────────────────────────────────────────────────────────────────────────┐
│  ONEMIND COMMAND CENTER                              [Zeus] [Settings]  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │   MORNING BRIEFING      │  │       INTELLIGENCE SCORES           │  │
│  │   Tuesday, Jan 28       │  │  ┌─────┐  ┌─────┐  ┌─────┐         │  │
│  │   ☀️ 72°F Clear          │  │  │ HPI │  │ LEI │  │ GEI │         │  │
│  │                         │  │  │ 78  │  │ 65  │  │ 82  │         │  │
│  │   Energy: ████████░░ 8  │  │  │  ↑  │  │  →  │  │  ↑  │         │  │
│  │   Sleep:  ███████░░░ 7  │  │  └─────┘  └─────┘  └─────┘         │  │
│  │                         │  │                                     │  │
│  │   TOP 3 TODAY:          │  │   UNIFIED SCORE: 75 🟡              │  │
│  │   1. □ Review Q1 OKRs   │  │   "Focus on family time (LEI)"     │  │
│  │   2. □ Morning workout  │  └─────────────────────────────────────┘  │
│  │   3. □ Client call      │                                         │
│  │                         │  ┌─────────────────────────────────────┐  │
│  │   PENDING APPROVALS: 2  │  │       PATH BALANCE (This Week)      │  │
│  └─────────────────────────┘  │   HP ████████░░░░░░ 40%             │  │
│                               │   LE ████░░░░░░░░░░ 20%  ⚠️          │  │
│  ┌─────────────────────────┐  │   GE ████████░░░░░░ 40%             │  │
│  │   ACTIVE FOCUS          │  └─────────────────────────────────────┘  │
│  │   ┌───────────────────┐ │                                         │
│  │   │ Review Q1 OKRs    │ │  ┌─────────────────────────────────────┐  │
│  │   │     23:45         │ │  │       QUICK ACTIONS                 │  │
│  │   │   [Pause] [Done]  │ │  │  [+ Task] [+ Inbox] [+ Log] [🍅]   │  │
│  │   └───────────────────┘ │  │  [Ask Legacy...]                    │  │
│  └─────────────────────────┘  └─────────────────────────────────────┘  │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐│
│  │  SYSTEM: ✅ Agents Online | ✅ Todoist Synced (2m ago) | 3 pending  ││
│  └─────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────┘
```

#### Implementation Files

```
backend/lifeos/dashboard/
├── __init__.py
├── models.py          # DashboardState, BriefingData
├── service.py         # DashboardService
└── api.py             # GET /api/dashboard/briefing, /api/dashboard/state

frontend/lib/lifeos/dashboard/
├── models/
│   └── dashboard_state.dart
├── providers/
│   └── dashboard_provider.dart
├── screens/
│   └── command_center_screen.dart
└── widgets/
    ├── briefing_panel.dart
    ├── scores_widget.dart
    ├── path_balance_widget.dart
    ├── active_focus_widget.dart
    ├── quick_actions_bar.dart
    └── system_health_widget.dart
```

---

### 1.3 Operating Rhythm System (A→Q→W→D)

**Codex Reference:** Part 5 - The Operating Rhythm

The Operating Rhythm defines the cadence of planning and review cycles.

#### Cadence Hierarchy

```
ANNUAL THESIS (A)
    Direction & non-negotiables
    "What does this year need to be about?"
            │
            ▼
12-WEEK ROADMAP (Q)
    Themes & OKRs
    "What are the 3-4 outcomes that matter this quarter?"
            │
            ▼
WEEKLY PLAN (W)
    Scorecards, priorities, time blocks
    "What must happen this week to hit the quarter?"
            │
            ▼
DAILY EXECUTION (D)
    Single focus, ship, log
    "What is the ONE thing today?"
            │
            ▼
REVIEW/REFACTOR
    Metrics, lessons, updates
    "What did we learn? What changes?"
```

#### 1.3.1 Annual Thesis Module

**Purpose:** Set yearly direction, theme, and non-negotiables

**Fields:**
- Year
- Theme (one word + one sentence)
- 10-Year Vision Check
- 3-Year Targets (HPI, LEI, GEI targets)
- Annual Non-Negotiables (3-5 items)
- Revenue/Resource Targets
- Path Review (previous year scores)

**Template:**
```markdown
# Annual Thesis {{YEAR}}

## Theme
**One Word:** _______________
**One Sentence:** _______________

## 10-Year Vision
Still aligned? What shifted?

## 3-Year Targets
| Path | Intelligence Score | Target |
|------|--------------------|--------|
| HP   | HPI                | 85     |
| LE   | LEI                | 80     |
| GE   | GEI                | 85     |
| —    | SO (Sovereign)       | 83     |

## Annual Non-Negotiables
1. [ ]
2. [ ]
3. [ ]

## Resource Targets
- Income goal: $
- Investment goal: $
- Time allocation: HP __% / LE __% / GE __%
```

#### 1.3.2 Quarterly OKR Module

**Purpose:** 12-week sprints with measurable objectives

**Fields:**
- Quarter (Q1-Q4)
- Year
- Theme
- HP Objective + Key Results (3)
- LE Objective + Key Results (3)
- GE Objective + Key Results (3)
- Big Rocks (3-5 major projects)
- 12-Week Sprint View
- Resource Allocation
- Risks & Mitigations

**Schema:**
```sql
CREATE TABLE quarterly_plans (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    year INT NOT NULL,
    quarter INT NOT NULL CHECK (quarter BETWEEN 1 AND 4),
    theme VARCHAR(255),
    hp_objective TEXT,
    hp_key_results JSONB DEFAULT '[]',
    le_objective TEXT,
    le_key_results JSONB DEFAULT '[]',
    ge_objective TEXT,
    ge_key_results JSONB DEFAULT '[]',
    big_rocks JSONB DEFAULT '[]',
    resource_allocation JSONB DEFAULT '{}',
    risks JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, year, quarter)
);

CREATE TABLE key_result_progress (
    id UUID PRIMARY KEY,
    quarterly_plan_id UUID REFERENCES quarterly_plans(id),
    key_result_index INT,
    progress DECIMAL(5,2),  -- 0-100%
    recorded_at DATE,
    notes TEXT
);
```

#### 1.3.3 Weekly Planning Module

**Purpose:** Weekly planning ritual every Sunday

**Features:**
- Weekly theme
- Path priorities (one per path)
- Big 3 for the week
- Calendar scan view
- Scorecard check (previous week)
- AI instructions for agents

**Automation:**
- **Trigger:** Sunday 6:00 PM
- **Action:** Generate weekly plan template
- **Notification:** Push to review and complete

#### 1.3.4 Daily Briefing System

**Purpose:** Morning briefing + Evening review

**Morning Briefing (5:30 AM):**
- Energy/sleep data (from Oura or manual)
- Calendar summary
- Top 3 priorities
- Pending approvals
- Habit streaks at risk
- Weather + commute info

**Evening Review (9:00 PM):**
- What got done (auto from completed tasks)
- What didn't (incomplete tasks)
- Tomorrow's MIT (Most Important Task)
- Reflection prompt
- Gratitude capture

**Implementation:**
```python
# backend/lifeos/rhythm/briefing.py
class BriefingService:
    async def generate_morning_briefing(self, user_id: str) -> MorningBriefing
    async def generate_evening_review(self, user_id: str) -> EveningReview
    async def schedule_briefings(self, user_id: str, morning_time: time, evening_time: time)
```

---

### 1.4 HIVE Capture System

**Codex Reference:** Part 5 - HIVE Capture System

Universal capture and routing system for all inputs.

| Letter | Word | Meaning | Action |
|--------|------|---------|--------|
| **H** | Harvest | Gather inputs | Capture from any source |
| **I** | Integrate | Route to path | Direct to HP, LE, GE based on content |
| **V** | Validate | AI priority | Legacy AI helps verify importance |
| **E** | Execute | Take action | Complete with AI assistance |

#### Capture Sources
- Voice (via LiveKit)
- Text (chat with Legacy)
- Email (forwarded)
- Image (screenshot analysis)
- Quick add (mobile widget)
- Browser extension
- Slack/Discord

#### Auto-Routing Logic
```python
def route_to_path(item: InboxItem) -> str:
    """Use AI to determine the appropriate path for an inbox item."""
    keywords = {
        "HP": ["health", "fitness", "workout", "sleep", "habit", "learn", "skill", "money", "budget"],
        "LE": ["family", "home", "wife", "kids", "parents", "house", "legacy", "heritage"],
        "GE": ["business", "client", "revenue", "product", "marketing", "sales", "company"]
    }
    # AI analysis for ambiguous items
    return legacy_ai.classify_path(item.content)
```

---

## Phase 2: Life Domain Modules (Q2 2026)

### 2.1 Finance Module (HP-30)

**Codex Reference:** 30 Finance (Money)

Comprehensive personal finance management.

#### Sub-Modules

| Code | Name | Description |
|------|------|-------------|
| 30.2 | Net Worth | Assets vs liabilities overview |
| 30.3 | Banking | Account management, balances |
| 30.4 | Budget | Monthly budget with categories |
| 30.5 | Income | Multiple income stream tracking |
| 30.6 | Credit | Score tracking, dispute management |
| 30.7 | Investments | Portfolio holdings, performance |
| 30.8 | Debt | Payoff tracking, strategies |
| 30.9 | Taxes | Deductions, deadlines, documents |
| 30.10 | Funding | Loans, credit lines |
| 30.11 | Planning | Financial goals, projections |

#### Core Features

##### Net Worth Tracker
```sql
CREATE TABLE financial_accounts (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50), -- checking, savings, credit, investment, loan, asset
    institution VARCHAR(255),
    balance DECIMAL(15,2),
    currency VARCHAR(3) DEFAULT 'USD',
    is_asset BOOLEAN DEFAULT true,
    is_manual BOOLEAN DEFAULT true,
    plaid_account_id VARCHAR(255),
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE net_worth_snapshots (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    snapshot_date DATE NOT NULL,
    total_assets DECIMAL(15,2),
    total_liabilities DECIMAL(15,2),
    net_worth DECIMAL(15,2),
    breakdown JSONB, -- by account type
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, snapshot_date)
);
```

##### Budget System
```sql
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(100) NOT NULL,
    parent_id UUID REFERENCES budget_categories(id),
    budget_type VARCHAR(20), -- income, expense, transfer
    color VARCHAR(7),
    icon VARCHAR(50),
    is_system BOOLEAN DEFAULT false
);

CREATE TABLE monthly_budgets (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    year INT,
    month INT,
    category_id UUID REFERENCES budget_categories(id),
    budgeted_amount DECIMAL(12,2),
    actual_amount DECIMAL(12,2) DEFAULT 0,
    rollover_amount DECIMAL(12,2) DEFAULT 0,
    UNIQUE(user_id, year, month, category_id)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    account_id UUID REFERENCES financial_accounts(id),
    category_id UUID REFERENCES budget_categories(id),
    date DATE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    description TEXT,
    merchant VARCHAR(255),
    is_income BOOLEAN DEFAULT false,
    is_transfer BOOLEAN DEFAULT false,
    is_recurring BOOLEAN DEFAULT false,
    tags TEXT[],
    plaid_transaction_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Credit Management
```sql
CREATE TABLE credit_scores (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    bureau VARCHAR(20), -- experian, equifax, transunion
    score INT,
    score_date DATE,
    factors JSONB, -- Contributing factors
    source VARCHAR(50), -- manual, credit_karma, etc
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE credit_disputes (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    bureau VARCHAR(20),
    creditor VARCHAR(255),
    account_number VARCHAR(100),
    dispute_reason TEXT,
    dispute_date DATE,
    status VARCHAR(30), -- pending, in_review, resolved, rejected
    resolution TEXT,
    documents JSONB, -- Attached files
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Integrations
- **Plaid** - Bank account sync, transactions
- **Finicity** - Alternative to Plaid
- **Manual Entry** - For accounts without API access

#### GEI Score Contribution
Finance module feeds directly into GEI (Generational Endeavor Intelligence) score:
- Runway calculation (months of expenses covered)
- Free cash flow (income - expenses)
- Net worth trend

---

### 2.2 Relationships Module (HP-32)

**Codex Reference:** 32 Relationships (Social)

Personal CRM for managing relationships.

#### Features

##### Contact Management
```sql
CREATE TABLE contacts (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    nickname VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    photo_url TEXT,
    birthday DATE,
    anniversary DATE,
    relationship_type VARCHAR(50), -- friend, family, colleague, client, mentor
    importance INT DEFAULT 3, -- 1-5 scale
    path VARCHAR(10), -- HP, LE, GE (which path they relate to)
    company VARCHAR(255),
    job_title VARCHAR(255),
    location VARCHAR(255),
    notes TEXT,
    tags TEXT[],
    social_profiles JSONB,
    last_contact_date DATE,
    desired_frequency_days INT, -- How often to reach out
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Interaction Tracking
```sql
CREATE TABLE interactions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    contact_id UUID REFERENCES contacts(id),
    interaction_type VARCHAR(30), -- call, meeting, email, text, social, gift
    interaction_date TIMESTAMPTZ,
    duration_minutes INT,
    summary TEXT,
    sentiment VARCHAR(20), -- positive, neutral, negative
    follow_up_needed BOOLEAN DEFAULT false,
    follow_up_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Relationship Health
```python
def calculate_relationship_health(contact: Contact) -> int:
    """Calculate health score 0-100 based on interaction frequency."""
    days_since_contact = (today - contact.last_contact_date).days
    desired_frequency = contact.desired_frequency_days or 30

    if days_since_contact <= desired_frequency:
        return 100
    elif days_since_contact <= desired_frequency * 2:
        return 75
    elif days_since_contact <= desired_frequency * 3:
        return 50
    else:
        return max(0, 100 - (days_since_contact - desired_frequency))
```

##### Gift Tracking
```sql
CREATE TABLE gifts (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    contact_id UUID REFERENCES contacts(id),
    gift_type VARCHAR(20), -- given, received
    occasion VARCHAR(100),
    description TEXT,
    amount DECIMAL(10,2),
    gift_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Features
- Birthday/anniversary reminders (integrated with notifications)
- "Haven't talked to X in Y days" alerts
- Network visualization (graph view)
- Follow-up task auto-creation
- Contact import (Google, Apple, CSV)

---

### 2.3 Biometrics Dashboard (HP-27)

**Codex Reference:** 27 Body (Health & Wellness)

Unified health metrics from wearables and manual entry.

#### Sub-Modules

| Code | Name | Data Source |
|------|------|-------------|
| 27.2 | Biometrics | Oura, Apple Watch, Manual |
| 27.3 | Nutrition | Manual, MyFitnessPal |
| 27.4 | Sleep | Oura, Apple Watch |
| 27.5 | Recovery | Oura (HRV, Readiness) |
| 27.6 | Movement | Steps, Workouts |
| 27.7 | Mental Wellness | Mood (from habits), Journal |
| 27.8 | Daily Metrics | Aggregated dashboard |

#### Schema
```sql
CREATE TABLE biometric_readings (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    metric_type VARCHAR(50), -- sleep_score, hrv, resting_hr, steps, etc
    value DECIMAL(10,2),
    unit VARCHAR(20),
    recorded_at TIMESTAMPTZ,
    source VARCHAR(50), -- oura, apple_watch, manual
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_biometrics_user_type ON biometric_readings(user_id, metric_type, recorded_at);
```

#### Integrations
- **Oura API** - Sleep, readiness, activity
- **Apple HealthKit** - Via mobile app
- **Strava** - Workouts
- **Withings** - Weight, body composition

#### HPI Score Contribution
Biometrics feed directly into HPI (Holistic Performance Intelligence):
- Sleep quality average
- Recovery/readiness trends
- Activity consistency

---

## Phase 3: Family & Legacy Modules (Q3 2026)

### 3.1 Home Management Module (LE-50)

**Codex Reference:** 50 Home (Sanctuary & Property)

Property and home management system.

#### Features

##### Property Dashboard
```sql
CREATE TABLE properties (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255),
    address TEXT,
    property_type VARCHAR(50), -- primary, rental, vacation, land
    purchase_date DATE,
    purchase_price DECIMAL(15,2),
    current_value DECIMAL(15,2),
    mortgage_balance DECIMAL(15,2),
    monthly_payment DECIMAL(10,2),
    square_feet INT,
    bedrooms INT,
    bathrooms DECIMAL(3,1),
    lot_size DECIMAL(10,2),
    year_built INT,
    photos JSONB,
    documents JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Maintenance Tracker
```sql
CREATE TABLE home_maintenance (
    id UUID PRIMARY KEY,
    property_id UUID REFERENCES properties(id),
    maintenance_type VARCHAR(50), -- hvac, plumbing, electrical, roof, etc
    description TEXT,
    frequency VARCHAR(20), -- monthly, quarterly, annually, as_needed
    last_completed DATE,
    next_due DATE,
    estimated_cost DECIMAL(10,2),
    vendor VARCHAR(255),
    notes TEXT,
    is_recurring BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE maintenance_logs (
    id UUID PRIMARY KEY,
    maintenance_id UUID REFERENCES home_maintenance(id),
    completed_date DATE,
    actual_cost DECIMAL(10,2),
    vendor VARCHAR(255),
    notes TEXT,
    receipts JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Home Inventory
```sql
CREATE TABLE home_inventory (
    id UUID PRIMARY KEY,
    property_id UUID REFERENCES properties(id),
    name VARCHAR(255),
    category VARCHAR(100),
    location VARCHAR(100), -- room/area
    purchase_date DATE,
    purchase_price DECIMAL(10,2),
    current_value DECIMAL(10,2),
    warranty_expires DATE,
    serial_number VARCHAR(255),
    model_number VARCHAR(255),
    photos JSONB,
    documents JSONB, -- receipts, manuals
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Home Projects
Integration with existing Goals/Tasks modules for home improvement projects.

#### LEI Score Contribution
- Maintenance tasks completion rate
- Property value trends
- Project completion

---

### 3.2 Family Module (LE-51/52/53)

**Codex Reference:** 51 Marriage, 52 Children, 53 Parents

Family relationship and responsibility management.

#### 3.2.1 Partnership Module (LE-51)

```sql
CREATE TABLE partnership (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    partner_name VARCHAR(255),
    anniversary DATE,
    love_language VARCHAR(50),
    communication_style TEXT,
    shared_goals JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE date_nights (
    id UUID PRIMARY KEY,
    partnership_id UUID REFERENCES partnership(id),
    date DATE,
    activity TEXT,
    location VARCHAR(255),
    rating INT, -- 1-5
    notes TEXT,
    photos JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE partnership_check_ins (
    id UUID PRIMARY KEY,
    partnership_id UUID REFERENCES partnership(id),
    check_in_date DATE,
    topics_discussed JSONB,
    action_items JSONB,
    relationship_health INT, -- 1-10
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### 3.2.2 Children Module (LE-52)

```sql
CREATE TABLE children (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255),
    birth_date DATE,
    school VARCHAR(255),
    grade VARCHAR(50),
    interests JSONB,
    allergies TEXT,
    medical_notes TEXT,
    emergency_contacts JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE child_milestones (
    id UUID PRIMARY KEY,
    child_id UUID REFERENCES children(id),
    milestone_type VARCHAR(100),
    milestone_date DATE,
    description TEXT,
    photos JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE child_activities (
    id UUID PRIMARY KEY,
    child_id UUID REFERENCES children(id),
    activity_name VARCHAR(255),
    schedule VARCHAR(255),
    location VARCHAR(255),
    cost_monthly DECIMAL(10,2),
    contact_info TEXT,
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### 3.2.3 Elder Care Module (LE-53)

```sql
CREATE TABLE parents (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255),
    relationship VARCHAR(50), -- mother, father, in-law
    birth_date DATE,
    address TEXT,
    phone VARCHAR(50),
    emergency_contact VARCHAR(255),
    medical_conditions JSONB,
    medications JSONB,
    doctors JSONB,
    care_needs TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE parent_care_logs (
    id UUID PRIMARY KEY,
    parent_id UUID REFERENCES parents(id),
    log_date DATE,
    visit_type VARCHAR(50), -- in_person, call, video
    duration_minutes INT,
    notes TEXT,
    concerns JSONB,
    follow_ups JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Family Calendar
Shared calendar view integrating:
- School events
- Activities schedules
- Medical appointments
- Family rituals/traditions

#### LEI Score Contribution
- Family time logged
- Rituals kept
- Care responsibilities met

---

### 3.3 Heritage Module (LE-54/56/57)

**Codex Reference:** 54 Heritage, 56 Wealth, 57 Legal

Legacy preservation and estate planning.

#### Features

##### Family Tree
```sql
CREATE TABLE family_tree (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    person_name VARCHAR(255),
    relationship_to_user VARCHAR(100),
    birth_date DATE,
    death_date DATE,
    birth_place VARCHAR(255),
    occupation VARCHAR(255),
    biography TEXT,
    photos JSONB,
    documents JSONB,
    parent_ids UUID[], -- References to parents in tree
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Family Stories
```sql
CREATE TABLE family_stories (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    title VARCHAR(255),
    story_type VARCHAR(50), -- tradition, recipe, memory, lesson
    content TEXT,
    related_people UUID[], -- family_tree ids
    photos JSONB,
    audio_url TEXT,
    video_url TEXT,
    date_of_story DATE,
    recorded_by VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Estate Documents
```sql
CREATE TABLE estate_documents (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    document_type VARCHAR(50), -- will, trust, poa, insurance, deed
    title VARCHAR(255),
    description TEXT,
    file_url TEXT,
    attorney VARCHAR(255),
    created_date DATE,
    expiration_date DATE,
    review_date DATE,
    beneficiaries JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Legacy Letters
```sql
CREATE TABLE legacy_letters (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    recipient VARCHAR(255),
    title VARCHAR(255),
    content TEXT,
    delivery_trigger VARCHAR(50), -- date, event, death
    delivery_date DATE,
    is_delivered BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Phase 4: Business & Professional Modules (Q4 2026)

### 4.1 Personal Brand Module (GE-75)

**Codex Reference:** 75 Zeus Delacruz Enterprise (Personal Brand)

Personal brand and content management.

#### Features

##### Content Calendar
```sql
CREATE TABLE content_calendar (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    title VARCHAR(255),
    content_type VARCHAR(50), -- post, article, video, podcast
    platform VARCHAR(50), -- twitter, linkedin, youtube, blog
    status VARCHAR(30), -- idea, draft, scheduled, published
    scheduled_date TIMESTAMPTZ,
    published_date TIMESTAMPTZ,
    content TEXT,
    media_urls JSONB,
    hashtags TEXT[],
    engagement_metrics JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Platform Analytics
```sql
CREATE TABLE social_metrics (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    platform VARCHAR(50),
    metric_date DATE,
    followers INT,
    following INT,
    posts_count INT,
    engagement_rate DECIMAL(5,2),
    impressions INT,
    profile_views INT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, platform, metric_date)
);
```

##### Portfolio
```sql
CREATE TABLE portfolio_items (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    title VARCHAR(255),
    category VARCHAR(100),
    description TEXT,
    url TEXT,
    image_url TEXT,
    client VARCHAR(255),
    completion_date DATE,
    skills_used TEXT[],
    testimonial TEXT,
    is_featured BOOLEAN DEFAULT false,
    display_order INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### 4.2 Business Dashboard (GE-75-85)

**Codex Reference:** Generational Endeavor path

Multi-business management dashboard.

#### Features

##### Business Entities
```sql
CREATE TABLE businesses (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(255),
    legal_name VARCHAR(255),
    entity_type VARCHAR(50), -- llc, corp, sole_prop, partnership
    industry VARCHAR(100),
    founded_date DATE,
    ein VARCHAR(20),
    state_of_formation VARCHAR(50),
    website VARCHAR(255),
    description TEXT,
    logo_url TEXT,
    status VARCHAR(30), -- active, inactive, sold, closed
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Revenue Tracking
```sql
CREATE TABLE business_revenue (
    id UUID PRIMARY KEY,
    business_id UUID REFERENCES businesses(id),
    period_start DATE,
    period_end DATE,
    gross_revenue DECIMAL(15,2),
    net_revenue DECIMAL(15,2),
    expenses DECIMAL(15,2),
    profit DECIMAL(15,2),
    revenue_breakdown JSONB, -- by product/service
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

##### Client Management
```sql
CREATE TABLE business_clients (
    id UUID PRIMARY KEY,
    business_id UUID REFERENCES businesses(id),
    name VARCHAR(255),
    company VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    status VARCHAR(30), -- lead, prospect, active, churned
    lifetime_value DECIMAL(15,2),
    first_purchase_date DATE,
    last_purchase_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### GEI Score Contribution
- Revenue trends
- Client health
- Runway calculation

---

## Phase 5: AI Enhancement (2027)

### 5.1 Four-Agent Personas

**Codex Reference:** Part 4 - The AI Layer (Hive Mind)

Implement specialized agent personas within Legacy.

#### Agent Definitions

| Agent | Persona | Domain | Personality |
|-------|---------|--------|-------------|
| **UI/Legacy** | Mother AI | Cross-domain | Wise, calm, big picture |
| **HP Agent** | Performance Coach | HP (25-49) | Encouraging, data-driven |
| **LE Agent** | Family Steward | LE (50-74) | Warm, protective |
| **GE Agent** | Business Advisor | GE (75-99) | Strategic, revenue-focused |

#### Routing Logic
```python
class AgentRouter:
    def route_query(self, query: str, context: dict) -> str:
        """Determine which agent persona should handle the query."""

        # Keywords for routing
        hp_keywords = ["health", "fitness", "workout", "sleep", "habit",
                       "learn", "skill", "personal", "self"]
        le_keywords = ["family", "home", "wife", "husband", "kids",
                       "parents", "house", "legacy", "heritage"]
        ge_keywords = ["business", "client", "revenue", "product",
                       "marketing", "sales", "company", "money"]

        query_lower = query.lower()

        if any(kw in query_lower for kw in hp_keywords):
            return "hp_agent"
        elif any(kw in query_lower for kw in le_keywords):
            return "le_agent"
        elif any(kw in query_lower for kw in ge_keywords):
            return "ge_agent"
        else:
            return "legacy"  # Default to main agent
```

#### Implementation
```python
# backend/agno/agents/personas.py
HP_PERSONA = """You are the HP Performance Coach within the Legacy AI system.
Your role is to help Zeus optimize his personal performance across:
- Health and fitness
- Sleep and recovery
- Habits and routines
- Learning and skills
- Personal finance

Personality: Encouraging, data-driven, health-focused.
Always connect recommendations to the HPI (Holistic Performance Intelligence) score."""

LE_PERSONA = """You are the LE Family Steward within the Legacy AI system.
Your role is to help Zeus protect and nurture his family legacy across:
- Home and property
- Marriage and partnership
- Children and parenting
- Parents and elder care
- Heritage and traditions

Personality: Warm, protective, heritage-minded.
Always connect recommendations to the LEI (Legacy Evolution Intelligence) score."""

GE_PERSONA = """You are the GE Business Advisor within the Legacy AI system.
Your role is to help Zeus build generational wealth across:
- Business operations
- Revenue and profitability
- Client relationships
- Investments and assets
- Strategic planning

Personality: Strategic, revenue-focused, opportunity-seeking.
Always connect recommendations to the GEI (Generational Endeavor Intelligence) score."""
```

---

### 5.2 Skills System

**Codex Reference:** Part 4 - Skills System

Modular skills that agents can invoke.

#### Core Skills

| Skill | Description | Trigger |
|-------|-------------|---------|
| morning_brief | Generate daily briefing | 5:30 AM cron |
| evening_review | Summarize day, plan tomorrow | 9:00 PM cron |
| inbox_triage | Process and route inbox items | On demand |
| weekly_planning | Generate weekly plan | Sunday 6 PM |
| scorecard_calc | Calculate HPI/LEI/GEI scores | Weekly cron |
| financial_review | Analyze spending patterns | Monthly |
| relationship_check | Identify neglected contacts | Weekly |

#### Skill Implementation
```python
# backend/agno/skills/morning_brief.py
class MorningBriefSkill:
    """Generate comprehensive morning briefing."""

    async def execute(self, user_id: str) -> MorningBriefing:
        # Gather data
        sleep_data = await self.get_sleep_data(user_id)
        calendar = await self.get_today_calendar(user_id)
        top_tasks = await self.get_priority_tasks(user_id)
        habit_streaks = await self.get_at_risk_streaks(user_id)
        approvals = await self.get_pending_approvals(user_id)

        # Generate briefing
        briefing = MorningBriefing(
            date=date.today(),
            energy_level=sleep_data.readiness_score,
            sleep_quality=sleep_data.sleep_score,
            calendar_summary=self.summarize_calendar(calendar),
            top_3=top_tasks[:3],
            streaks_at_risk=habit_streaks,
            pending_approvals=len(approvals),
            weather=await self.get_weather(),
        )

        return briefing
```

---

### 5.3 OneMind Insights Engine

**Codex Reference:** Replaces Exist.io

AI-powered personal analytics.

#### Features

##### Automatic Correlations
```python
class InsightsEngine:
    """Find correlations in personal data."""

    async def find_correlations(self, user_id: str) -> list[Correlation]:
        """Analyze all data sources to find meaningful patterns."""

        correlations = []

        # Example correlations to detect:
        # - Sleep vs productivity
        # - Exercise vs mood
        # - Screen time vs sleep quality
        # - Social interaction vs happiness
        # - Weather vs energy

        return correlations
```

##### Conversational Insights
Natural language queries about personal data:
- "How has my sleep affected my productivity this month?"
- "What happens to my mood when I skip workouts?"
- "Show me my best weeks and what they had in common"

##### Custom Metrics
User-defined tracking for anything:
```sql
CREATE TABLE custom_metrics (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    name VARCHAR(100),
    description TEXT,
    unit VARCHAR(50),
    data_type VARCHAR(20), -- number, boolean, scale, text
    min_value DECIMAL(10,2),
    max_value DECIMAL(10,2),
    target_value DECIMAL(10,2),
    aggregation VARCHAR(20), -- sum, average, max, min, count
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE custom_metric_entries (
    id UUID PRIMARY KEY,
    metric_id UUID REFERENCES custom_metrics(id),
    value DECIMAL(10,2),
    recorded_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Integration Requirements

### Required API Integrations

| Integration | Purpose | Priority | Auth Type |
|-------------|---------|----------|-----------|
| Oura | Sleep, readiness, activity | P1 | OAuth 2.0 |
| Plaid | Bank accounts, transactions | P1 | API Key |
| Apple HealthKit | Health data (iOS) | P2 | Native |
| Strava | Workouts | P2 | OAuth 2.0 |
| Google Calendar | Calendar sync | P2 | OAuth 2.0 |
| Twitter/X | Social metrics | P3 | OAuth 2.0 |
| LinkedIn | Professional metrics | P3 | OAuth 2.0 |
| Withings | Weight, body comp | P3 | OAuth 2.0 |

### Existing Integrations to Enhance

| Integration | Current State | Enhancement |
|-------------|---------------|-------------|
| Todoist | ✅ Bidirectional sync | Add path tagging |
| Home Assistant | ✅ MCP tools | Add to home module UI |
| Calendar (CalDAV) | ✅ Basic sync | Family calendar view |
| GitHub | ⚠️ Partial | Full project tracking |

---

## Database Migration Plan

### New Tables Required

| Phase | Tables | Count |
|-------|--------|-------|
| Phase 1 | intelligence_scores, quarterly_plans, key_result_progress, annual_theses, weekly_plans, daily_briefings | 6 |
| Phase 2 | financial_accounts, net_worth_snapshots, budget_categories, monthly_budgets, transactions, credit_scores, credit_disputes, contacts, interactions, gifts, biometric_readings | 11 |
| Phase 3 | properties, home_maintenance, maintenance_logs, home_inventory, partnership, date_nights, partnership_check_ins, children, child_milestones, child_activities, parents, parent_care_logs, family_tree, family_stories, estate_documents, legacy_letters | 16 |
| Phase 4 | content_calendar, social_metrics, portfolio_items, businesses, business_revenue, business_clients | 6 |
| Phase 5 | custom_metrics, custom_metric_entries, correlations, insights | 4 |

**Total New Tables: 43**

---

## API Endpoints Summary

### Phase 1 Endpoints

```
# Intelligence Scores
GET    /api/scores/current
GET    /api/scores/history
GET    /api/scores/{type}/breakdown
POST   /api/scores/calculate

# Dashboard
GET    /api/dashboard/briefing
GET    /api/dashboard/state
GET    /api/dashboard/quick-stats

# Operating Rhythm
GET    /api/rhythm/annual/{year}
POST   /api/rhythm/annual
GET    /api/rhythm/quarterly/{year}/{quarter}
POST   /api/rhythm/quarterly
GET    /api/rhythm/weekly/{year}/{week}
POST   /api/rhythm/weekly
GET    /api/rhythm/daily/{date}
POST   /api/rhythm/daily/generate-briefing
```

### Phase 2 Endpoints

```
# Finance
GET    /api/finance/accounts
POST   /api/finance/accounts
GET    /api/finance/net-worth
GET    /api/finance/budget/{year}/{month}
POST   /api/finance/transactions
GET    /api/finance/credit-scores

# Relationships
GET    /api/relationships/contacts
POST   /api/relationships/contacts
GET    /api/relationships/contacts/{id}/interactions
POST   /api/relationships/interactions
GET    /api/relationships/health-report

# Biometrics
GET    /api/biometrics/dashboard
GET    /api/biometrics/{metric_type}
POST   /api/biometrics/readings
GET    /api/biometrics/trends
```

---

## Success Metrics

### Phase 1 Success Criteria
- [ ] HPI/LEI/GEI scores calculating accurately
- [ ] Command Center dashboard functional
- [ ] Morning briefing generating daily
- [ ] Weekly planning workflow active

### Phase 2 Success Criteria
- [ ] Finance module tracking net worth
- [ ] 100+ contacts in relationships CRM
- [ ] Biometrics syncing from Oura

### Phase 3 Success Criteria
- [ ] Home maintenance on schedule
- [ ] Family calendar populated
- [ ] At least 10 family stories recorded

### Phase 4 Success Criteria
- [ ] Content calendar active
- [ ] Business revenue tracking
- [ ] GEI score reflecting business health

### Phase 5 Success Criteria
- [ ] Agent personas responding contextually
- [ ] Insights engine finding correlations
- [ ] Custom metrics tracking working

---

## Appendix A: Codex Directory Mapping

| Codex Path | OneMind OS Module | Status |
|------------|-------------------|--------|
| 00 Framework | Core System | ✅ |
| 01 Command | Dashboard | 🔲 To Build |
| 02 Agents | Agno Framework | ✅ |
| 06 Protocols | Documentation | ⚠️ Partial |
| 04 Resources | Assets | ⚠️ Partial |
| 06 Inbox | Inbox Module | ✅ |
| 08 Infrastructure | Backend | ✅ |
| 09 Automation | Workflows | ⚠️ Partial |
| 10 Development | Codebase | ✅ |
| 11 Security | Auth/Safety | ✅ |
| 12 Research | Knowledge Base | ✅ |
| 25 Identity | Journal/Goals | ✅ |
| 26 Intelligence | Learning (future) | 🔲 To Build |
| 27 Body | Biometrics | 🔲 To Build |
| 28 Mastery | Skills (future) | 🔲 To Build |
| 29 Life Systems | Habits/Routines | ✅ |
| 30 Finance | Finance Module | 🔲 To Build |
| 31 Sovereignty | Legal (future) | 🔲 To Build |
| 32 Relationships | Relationships | 🔲 To Build |
| 33 Joy | Recreation (future) | 🔲 To Build |
| 50 Home | Home Module | 🔲 To Build |
| 51 Marriage | Partnership | 🔲 To Build |
| 52 Children | Children | 🔲 To Build |
| 53 Parents | Elder Care | 🔲 To Build |
| 54 Heritage | Heritage | 🔲 To Build |
| 56 Wealth | Estate | 🔲 To Build |
| 75+ Business | Business Dashboard | 🔲 To Build |

---

## Appendix B: Technology Stack

### Current Stack
- **Backend:** Python 3.12+, FastAPI, Agno Framework
- **Frontend:** Flutter 3.6+, Dart, Riverpod
- **Database:** PostgreSQL 15+ with pgvector
- **Events:** NATS JetStream
- **AI:** Claude Opus 4.5 via AWS Bedrock
- **Voice:** LiveKit Agents
- **Search:** LanceDB

### Additional Tech Needed
- **Plaid SDK** - Financial data
- **Oura API Client** - Biometrics
- **Chart Libraries** - D3.js/fl_chart for visualizations
- **Graph Database** - Optional for relationship visualization

---

## Appendix C: Glossary

| Term | Definition |
|------|------------|
| **Codex** | The OneMind philosophical framework and doctrine |
| **OneMind OS** | The technical implementation of the Codex |
| **Legacy** | The AI assistant/agent |
| **HPI** | Holistic Performance Intelligence - HP path score |
| **LEI** | Legacy Evolution Intelligence - LE path score |
| **GEI** | Generational Endeavor Intelligence - GE path score |
| **UI** | Sovereign Operations - The brain/coordinator, aggregate of HPI + LEI + GEI |
| **HP** | Holistic Performance - Self optimization path |
| **LE** | Legacy Evolution - Family and heritage path |
| **GE** | Generational Endeavor - Business path |
| **HIVE** | Harvest, Integrate, Validate, Execute - capture system |
| **HITL** | Human-in-the-loop - approval workflow |

---

*Document generated from OneMind Codex analysis on 2026-01-28*
*"HUMAN (Zeus) + AI (Legacy) = ONE MIND"*
