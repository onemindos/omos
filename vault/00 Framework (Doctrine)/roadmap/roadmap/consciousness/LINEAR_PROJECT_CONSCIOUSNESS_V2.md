# Linear Project: OneMind Consciousness v2.0

> **Project Type:** Sprint
> **Duration:** 4 Weeks
> **Goal:** Production-ready personal productivity + human enhancement system

---

## Project Setup in Linear

1. Go to Linear → Settings → Projects → New Project
2. Name: `OneMind Consciousness v2.0`
3. Lead: Zeus
4. Target Date: +4 weeks from start
5. Create the following Labels:
   - `backend` (blue)
   - `frontend` (purple)
   - `ai` (green)
   - `infrastructure` (gray)
   - `p0-critical` (red)
   - `p1-high` (orange)
   - `p2-medium` (yellow)

---

## Epic 1: Context Engine (Foundation)

> **Priority:** P0 - Critical (Week 1)
> **Label:** backend, ai, p0-critical

### OMOS-301: Context Engine Framework

**Title:** Build Context Engine framework for universal AI awareness

**Description:**
Create the core Context Engine that aggregates ALL user data into a unified context object for AI consumption.

**Acceptance Criteria:**
- [ ] `ContextEngine` class in `backend/consciousness/context/engine.py`
- [ ] Async data aggregation from all LifeOS services
- [ ] `UnifiedContext` Pydantic model with all fields
- [ ] `build()` method returns complete context in <500ms
- [ ] Caching layer for frequently accessed data

**Technical Details:**
```python
# Location: backend/consciousness/context/engine.py
class ContextEngine:
    async def build(self, user_id: str) -> UnifiedContext
    async def build_minimal(self, user_id: str) -> MinimalContext
    async def invalidate_cache(self, user_id: str)
```

**Estimate:** 3 points
**Labels:** backend, ai, p0-critical

---

### OMOS-302: State Aggregator

**Title:** Implement State Aggregator for real-time user state

**Description:**
Build the state aggregator component that tracks current user state including energy, focus, location, and active context.

**Acceptance Criteria:**
- [ ] Pull from StateManager (`backend/consciousness/awareness/state.py`)
- [ ] Include current energy level from biometrics
- [ ] Include current focus mode
- [ ] Include location context (home/office/mobile)
- [ ] Include active path (HP/LE/GE)

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-303: Items Aggregator

**Title:** Implement Items Aggregator for tasks, habits, goals

**Description:**
Build the items aggregator that pulls active tasks, today's habits, and current goals into context.

**Acceptance Criteria:**
- [ ] Pull active tasks with priorities and due dates
- [ ] Pull today's habits with completion status
- [ ] Pull active goals with progress percentages
- [ ] Pull active routines for today
- [ ] Format for AI consumption (concise, structured)

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-304: Scores Aggregator

**Title:** Implement Scores Aggregator for HPI/LEI/GEI/UI

**Description:**
Build the scores aggregator that calculates and includes intelligence scores in context.

**Acceptance Criteria:**
- [ ] Include current HPI (Holistic Performance Index)
- [ ] Include current LEI (Legacy Evolution Index)
- [ ] Include current GEI (Generational Endeavor Index)
- [ ] Include UI (Unified Index = average)
- [ ] Include trend data (up/down/stable)
- [ ] Include component breakdowns

**Depends on:** OMOS-310 (Score Calculator)
**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-305: Calendar Aggregator

**Title:** Implement Calendar Aggregator for schedule awareness

**Description:**
Build the calendar aggregator that provides schedule context including today's events, upcoming deadlines, and available time slots.

**Acceptance Criteria:**
- [ ] Pull today's calendar events
- [ ] Pull next 3 days of events
- [ ] Calculate available time slots
- [ ] Include meeting context (with who, about what)
- [ ] Include deadline proximity warnings

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-306: History Aggregator

**Title:** Implement History Aggregator for pattern awareness

**Description:**
Build the history aggregator that provides recent activity patterns and conversation history.

**Acceptance Criteria:**
- [ ] Include last 10 completed tasks
- [ ] Include recent conversation summary
- [ ] Include habit completion patterns (7-day)
- [ ] Include productivity patterns by time of day
- [ ] Include path distribution (HP/LE/GE balance)

**Estimate:** 2 points
**Labels:** backend, ai

---

## Epic 2: Intelligence Scores

> **Priority:** P0 - Critical (Week 1)
> **Label:** backend, ai, p0-critical

### OMOS-310: HPI Calculator (Holistic Performance Index)

**Title:** Implement HPI score calculation

**Description:**
Build the Holistic Performance Index calculator that measures personal health and performance.

**Components & Weights:**
- Sleep Quality: 20%
- Recovery Score: 20%
- Habit Completion: 25%
- Task Completion: 20%
- Energy Levels: 15%

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/scores/hpi.py`
- [ ] Pull data from Oura/health integrations
- [ ] Calculate weighted score 0-100
- [ ] Include component breakdown
- [ ] Store historical scores for trending
- [ ] API endpoint: `GET /api/v1/scores/hpi`

**Estimate:** 3 points
**Labels:** backend, ai, p0-critical

---

### OMOS-311: LEI Calculator (Legacy Evolution Index)

**Title:** Implement LEI score calculation

**Description:**
Build the Legacy Evolution Index calculator that measures family/legacy building progress.

**Components & Weights:**
- Family Time: 25%
- Legacy Tasks: 25%
- Relationship Health: 20%
- Teaching/Mentoring: 15%
- Estate Planning: 15%

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/scores/lei.py`
- [ ] Pull data from LifeOS legacy path
- [ ] Calculate weighted score 0-100
- [ ] Include component breakdown
- [ ] API endpoint: `GET /api/v1/scores/lei`

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-312: GEI Calculator (Generational Endeavor Index)

**Title:** Implement GEI score calculation

**Description:**
Build the Generational Endeavor Index calculator that measures business/wealth building progress.

**Components & Weights:**
- Revenue Progress: 25%
- Business Tasks: 25%
- Investment Growth: 20%
- Network Building: 15%
- Skill Development: 15%

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/scores/gei.py`
- [ ] Pull data from LifeOS business path
- [ ] Calculate weighted score 0-100
- [ ] Include component breakdown
- [ ] API endpoint: `GET /api/v1/scores/gei`

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-313: UI Calculator (Unified Index)

**Title:** Implement Unified Index score calculation

**Description:**
Build the master Unified Index that combines HPI, LEI, and GEI into a single score.

**Calculation:**
```
UI = (HPI + LEI + GEI) / 3
```

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/scores/unified.py`
- [ ] Calculate as average of three indices
- [ ] Include balance indicator (which path needs attention)
- [ ] Include trend (improving/declining/stable)
- [ ] API endpoint: `GET /api/v1/scores/unified`
- [ ] API endpoint: `GET /api/v1/scores/all`

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-314: Score History & Trending

**Title:** Implement score history tracking and trend analysis

**Description:**
Build the persistence and trending system for all intelligence scores.

**Acceptance Criteria:**
- [ ] Database model for score history
- [ ] Daily score snapshots
- [ ] 7-day, 30-day, 90-day trends
- [ ] Trend direction detection (up/down/stable)
- [ ] API endpoint: `GET /api/v1/scores/history?days=30`

**Estimate:** 2 points
**Labels:** backend, infrastructure

---

## Epic 3: Smart Chat

> **Priority:** P0 - Critical (Week 1-2)
> **Label:** backend, ai, p0-critical

### OMOS-320: @Agent Routing

**Title:** Implement @agent mention routing in chat

**Description:**
Build the ability to route chat messages to specific specialist agents using @mentions.

**Supported Agents:**
- `@legacy` - General AI assistant (default)
- `@researcher` - Deep research tasks
- `@analyst` - Data analysis
- `@writer` - Content creation
- `@coder` - Code assistance
- `@finance` - Financial analysis

**Acceptance Criteria:**
- [ ] Create `backend/agno/chat/routing.py`
- [ ] Parse @mentions from message text
- [ ] Route to appropriate agent
- [ ] Fall back to Legacy for unknown mentions
- [ ] Return agent identity in response

**Example:**
```
User: @researcher find the latest studies on sleep optimization
Response: [Researcher Agent] Here are the latest studies...
```

**Estimate:** 3 points
**Labels:** backend, ai, p0-critical

---

### OMOS-321: /Commands System

**Title:** Implement slash command system in chat

**Description:**
Build the slash command system for quick actions from chat.

**Commands:**
| Command | Action |
|---------|--------|
| `/brief` | Generate morning briefing |
| `/schedule` | Show today's schedule |
| `/inbox` | Process inbox items |
| `/tasks` | List active tasks |
| `/scores` | Show current scores |
| `/workflow [name]` | Run a workflow |
| `/delegate [task]` | Delegate to agent |
| `/add [task]` | Quick add task |
| `/habit [name]` | Log habit completion |

**Acceptance Criteria:**
- [ ] Create `backend/agno/chat/commands.py`
- [ ] Parse /commands from message text
- [ ] Execute command handler
- [ ] Return formatted response
- [ ] Help command: `/help` lists all commands

**Estimate:** 3 points
**Labels:** backend, ai, p0-critical

---

### OMOS-322: Multi-Step Workflow from Chat

**Title:** Enable workflow creation from natural language chat

**Description:**
Build the ability to create and execute multi-step workflows from natural language requests.

**Example:**
```
User: Research competitors, analyze findings, and draft a summary report
System: Creating 3-step workflow:
  1. Research competitors (researcher agent)
  2. Analyze findings (analyst agent)
  3. Draft summary (writer agent)
Execute? [Yes/No]
```

**Acceptance Criteria:**
- [ ] Detect multi-step intent in messages
- [ ] Generate workflow definition from intent
- [ ] Show workflow preview for approval
- [ ] Execute workflow and stream progress
- [ ] Store generated workflows for reuse

**Estimate:** 5 points
**Labels:** backend, ai

---

### OMOS-323: Smart Chat Engine Integration

**Title:** Integrate all chat features into unified SmartChatEngine

**Description:**
Build the unified SmartChatEngine that combines @mentions, /commands, workflows, and context-aware chat.

**Acceptance Criteria:**
- [ ] Create `backend/agno/chat/engine.py`
- [ ] Process order: @mention → /command → workflow → normal chat
- [ ] Inject full context into every chat
- [ ] Stream responses via SSE
- [ ] Update existing `/api/v1/chat` endpoint

**Estimate:** 3 points
**Labels:** backend, ai, p0-critical

---

## Epic 4: Intelligent Scheduling

> **Priority:** P1 - High (Week 2)
> **Label:** backend, ai, p1-high

### OMOS-330: DO DATE Algorithm

**Title:** Implement Motion-style DO DATE scheduling algorithm

**Description:**
Build the intelligent scheduling algorithm that calculates the optimal time to work on tasks (DO DATE) rather than just tracking due dates.

**Algorithm Factors:**
- Task priority and deadline
- Estimated duration
- Available time slots
- Energy prediction for time slot
- Path balance (HP/LE/GE)
- Time preferences (morning/afternoon)

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/scheduling/do_date.py`
- [ ] Calculate optimal DO DATE for any task
- [ ] Score time slots based on factors
- [ ] Return top 3 recommended slots
- [ ] API endpoint: `POST /api/v1/tasks/{id}/schedule`

**Estimate:** 5 points
**Labels:** backend, ai, p1-high

---

### OMOS-331: Auto-Reschedule

**Title:** Implement automatic rescheduling on conflicts

**Description:**
Build the auto-reschedule system that adjusts DO DATEs when conflicts arise or tasks aren't completed.

**Triggers:**
- Calendar event added during DO DATE slot
- Task not started by DO DATE
- Task taking longer than estimated
- Higher priority task added

**Acceptance Criteria:**
- [ ] Detect scheduling conflicts
- [ ] Calculate new optimal slot
- [ ] Notify user of reschedule
- [ ] Update calendar integration
- [ ] Log reschedule reasons

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-332: Energy Matching

**Title:** Implement energy-aware task scheduling

**Description:**
Build the energy matching system that schedules hard tasks during high-energy periods and easy tasks during low-energy periods.

**Logic:**
- Hard tasks → High energy slots (energy > 80)
- Medium tasks → Medium energy slots (50-80)
- Easy tasks → Low energy slots (< 50)

**Acceptance Criteria:**
- [ ] Predict energy levels by time of day
- [ ] Use historical energy patterns
- [ ] Match task difficulty to energy slots
- [ ] Include in DO DATE scoring
- [ ] Pull from Oura/health data when available

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-333: Path Balance Scheduling

**Title:** Implement path-balanced scheduling

**Description:**
Build the path balance system that ensures HP, LE, and GE tasks are distributed throughout the week.

**Goal Distribution:**
- HP (Self): ~40% of productive time
- LE (Legacy): ~30% of productive time
- GE (Business): ~30% of productive time

**Acceptance Criteria:**
- [ ] Track path distribution per day/week
- [ ] Warn when path is neglected
- [ ] Factor into DO DATE scoring
- [ ] Suggest path-specific tasks when imbalanced
- [ ] API endpoint: `GET /api/v1/scheduling/balance`

**Estimate:** 2 points
**Labels:** backend, ai

---

## Epic 5: Workflow Automation

> **Priority:** P1 - High (Week 2)
> **Label:** backend, infrastructure, p1-high

### OMOS-340: HTTP Triggers

**Title:** Implement HTTP webhook triggers for workflows

**Description:**
Build the HTTP trigger system that allows workflows to be started via webhook URLs.

**Acceptance Criteria:**
- [ ] Create `backend/agno/workflows/triggers/http.py`
- [ ] Generate unique webhook URLs per workflow
- [ ] Accept POST with JSON payload
- [ ] Inject payload into workflow context
- [ ] Secure with API key or signature
- [ ] API endpoint: `POST /api/v1/workflows/{id}/trigger`

**Example:**
```
POST /api/v1/workflows/morning-brief/trigger
Authorization: Bearer xxx
{"extra_context": "Focus on project deadlines"}
```

**Estimate:** 3 points
**Labels:** backend, infrastructure, p1-high

---

### OMOS-341: Schedule Triggers (Cron)

**Title:** Implement cron-based schedule triggers for workflows

**Description:**
Build the schedule trigger system that runs workflows on a time-based schedule.

**Acceptance Criteria:**
- [ ] Create `backend/agno/workflows/triggers/schedule.py`
- [ ] Support cron syntax for scheduling
- [ ] Integrate with APScheduler or similar
- [ ] Store schedule definitions in database
- [ ] API endpoints for CRUD on schedules

**Examples:**
```python
# Morning briefing at 7am daily
schedule = "0 7 * * *"

# Weekly review on Sunday at 6pm
schedule = "0 18 * * 0"
```

**Estimate:** 3 points
**Labels:** backend, infrastructure

---

### OMOS-342: Smart Approvals (agent_decides)

**Title:** Implement smart approval modes for workflows

**Description:**
Enhance the existing Approval step with three approval modes based on Relevance AI pattern.

**Modes:**
- `approval` - Always require human approval
- `agent_decides` - Agent auto-approves if confident (>85%), asks if unsure
- `autorun` - Only stop for critical/destructive actions

**Acceptance Criteria:**
- [ ] Extend `backend/agno/workflows/workflow/steps/approval.py`
- [ ] Add confidence calculation to agent responses
- [ ] Implement mode-based approval logic
- [ ] Log all auto-approval decisions
- [ ] UI notification for pending approvals

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-343: Event Triggers

**Title:** Implement event-based triggers for workflows

**Description:**
Build the event trigger system that starts workflows based on system events.

**Events:**
- `inbox.item_added` - New inbox item
- `task.due_soon` - Task due within 24 hours
- `habit.missed` - Habit not completed by deadline
- `score.dropped` - Score dropped significantly
- `calendar.event_starting` - Event starting in 15 min

**Acceptance Criteria:**
- [ ] Create `backend/agno/workflows/triggers/events.py`
- [ ] Event subscription model
- [ ] Trigger matching on event type
- [ ] Inject event data into workflow context

**Estimate:** 3 points
**Labels:** backend, infrastructure

---

## Epic 6: Inbox Intelligence

> **Priority:** P1 - High (Week 2-3)
> **Label:** backend, ai, p1-high

### OMOS-350: AI Triage Engine

**Title:** Implement AI-powered inbox triage

**Description:**
Build the AI triage engine that automatically classifies and routes inbox items.

**Classification Output:**
- Path (HP, LE, GE)
- Priority (high, medium, low)
- Action type (task, note, event, reference, archive)
- Suggested title
- Due date (if applicable)

**Acceptance Criteria:**
- [ ] Create `backend/lifeos/inbox/intelligence.py`
- [ ] Analyze inbox item content with AI
- [ ] Return structured triage result
- [ ] Confidence score for classification
- [ ] API endpoint: `POST /api/v1/inbox/{id}/triage`

**Estimate:** 3 points
**Labels:** backend, ai, p1-high

---

### OMOS-351: Path Routing

**Title:** Implement automatic path routing for inbox items

**Description:**
Build the path routing system that assigns inbox items to HP, LE, or GE paths.

**Path Detection:**
- HP: health, fitness, personal growth, self-care
- LE: family, relationships, estate, teaching
- GE: business, money, investing, networking

**Acceptance Criteria:**
- [ ] Detect path from content keywords
- [ ] Use AI for ambiguous items
- [ ] Allow manual override
- [ ] Track routing accuracy

**Estimate:** 2 points
**Labels:** backend, ai

---

### OMOS-352: Task Extraction

**Title:** Implement task extraction from natural language

**Description:**
Build the task extraction system that creates tasks from natural language inbox items.

**Example:**
```
Input: "Need to call mom about Sunday dinner and also pick up groceries"
Output:
  - Task 1: "Call mom about Sunday dinner" (LE path)
  - Task 2: "Pick up groceries" (HP path)
```

**Acceptance Criteria:**
- [ ] Extract multiple tasks from single input
- [ ] Assign path and priority
- [ ] Detect due dates from text
- [ ] Create tasks in LifeOS
- [ ] API endpoint: `POST /api/v1/inbox/{id}/extract`

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-353: Quick Capture Processing

**Title:** Implement quick capture inbox processing

**Description:**
Build the quick capture system that processes voice memos, screenshots, and quick notes.

**Acceptance Criteria:**
- [ ] Process text notes
- [ ] Transcribe voice memos (Whisper)
- [ ] OCR screenshots/images
- [ ] Auto-triage after processing
- [ ] API endpoint: `POST /api/v1/inbox/capture`

**Estimate:** 3 points
**Labels:** backend, ai

---

## Epic 7: Operating Rhythm

> **Priority:** P2 - Medium (Week 3)
> **Label:** backend, ai, p2-medium

### OMOS-360: Morning Briefing Workflow

**Title:** Implement automated morning briefing workflow

**Description:**
Build the morning briefing workflow that generates a personalized daily briefing.

**Briefing Contents:**
- Today's schedule overview
- Priority tasks for the day
- Habits to complete
- Current scores and trends
- Weather/commute info
- Motivational insight

**Acceptance Criteria:**
- [ ] Create workflow definition
- [ ] Schedule trigger at 7am
- [ ] Generate briefing with AI
- [ ] Send to preferred channel (app notification, email)
- [ ] API endpoint: `GET /api/v1/briefings/morning`

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-361: Evening Review Workflow

**Title:** Implement automated evening review workflow

**Description:**
Build the evening review workflow that summarizes the day and plans tomorrow.

**Review Contents:**
- Tasks completed vs planned
- Habits completed
- Score changes
- Key accomplishments
- Tomorrow's priorities

**Acceptance Criteria:**
- [ ] Create workflow definition
- [ ] Schedule trigger at 8pm
- [ ] Generate review with AI
- [ ] Prompt for reflection input
- [ ] API endpoint: `GET /api/v1/briefings/evening`

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-362: Weekly Planning Automation

**Title:** Implement weekly planning automation

**Description:**
Build the weekly planning system that helps plan the upcoming week.

**Planning Steps:**
1. Review last week's accomplishments
2. Check quarterly OKR progress
3. Identify key priorities for week
4. Schedule DO DATEs for priorities
5. Balance path distribution

**Acceptance Criteria:**
- [ ] Create workflow definition
- [ ] Schedule trigger Sunday 6pm
- [ ] Interactive planning session
- [ ] Generate weekly plan document
- [ ] API endpoint: `POST /api/v1/planning/weekly`

**Estimate:** 3 points
**Labels:** backend, ai

---

### OMOS-363: Quarterly OKR Tracking

**Title:** Implement quarterly OKR tracking system

**Description:**
Build the OKR tracking system that monitors quarterly goals and key results.

**Features:**
- OKR definition and storage
- Progress tracking (manual + automatic)
- Weekly check-in prompts
- Quarterly review generation

**Acceptance Criteria:**
- [ ] OKR data model
- [ ] CRUD API endpoints
- [ ] Progress calculation from linked tasks
- [ ] Quarterly review workflow
- [ ] API endpoints: `GET/POST /api/v1/okrs`

**Estimate:** 3 points
**Labels:** backend

---

## Epic 8: Frontend (Flutter)

> **Priority:** P1 - High (Week 3-4)
> **Label:** frontend, p1-high

### OMOS-370: Command Center Dashboard

**Title:** Build Command Center dashboard in Flutter

**Description:**
Build the new Command Center dashboard that serves as the home screen.

**Components:**
- Score widgets (HPI, LEI, GEI, UI)
- Today's schedule timeline
- Priority tasks list
- Quick capture button
- Morning/Evening briefing cards

**Acceptance Criteria:**
- [ ] Create `frontend/lib/lifeos/dashboard/command_center.dart`
- [ ] Score widget components
- [ ] Schedule timeline widget
- [ ] Task list widget
- [ ] Quick capture FAB
- [ ] Pull-to-refresh

**Estimate:** 5 points
**Labels:** frontend, p1-high

---

### OMOS-371: Score Widgets

**Title:** Build score visualization widgets

**Description:**
Build the reusable score widgets that display HPI, LEI, GEI, and UI.

**Widget Features:**
- Circular progress indicator
- Score value (0-100)
- Trend indicator (up/down/stable)
- Tap for detail breakdown
- Animated transitions

**Acceptance Criteria:**
- [ ] Create `frontend/lib/lifeos/widgets/score_widget.dart`
- [ ] Circular progress design
- [ ] Color coding by score level
- [ ] Trend arrow indicator
- [ ] Detail bottom sheet on tap

**Estimate:** 3 points
**Labels:** frontend

---

### OMOS-372: Chat UI Enhancements

**Title:** Enhance chat UI for @mentions and /commands

**Description:**
Enhance the existing chat UI to support @mentions, /commands, and workflow progress.

**Enhancements:**
- @mention autocomplete
- /command autocomplete
- Workflow progress indicator
- Agent identity badges
- Rich message formatting

**Acceptance Criteria:**
- [ ] Autocomplete popup for @ and /
- [ ] Agent avatar in responses
- [ ] Workflow step progress UI
- [ ] Code block formatting
- [ ] Markdown support

**Estimate:** 3 points
**Labels:** frontend

---

### OMOS-373: Navigation Restructure

**Title:** Implement new navigation structure

**Description:**
Restructure the app navigation based on the three paths (HP, LE, GE).

**New Structure:**
```
Bottom Nav:
- Home (Command Center)
- Chat (UIO)
- Paths (HP/LE/GE tabs)
- Inbox
- More (Settings, etc.)
```

**Acceptance Criteria:**
- [ ] New bottom navigation
- [ ] Path-based tab view
- [ ] Smooth transitions
- [ ] Deep linking support
- [ ] Badge indicators

**Estimate:** 3 points
**Labels:** frontend

---

## Summary

| Epic | Issues | Total Points | Priority |
|------|--------|--------------|----------|
| Context Engine | 6 | 13 | P0 |
| Intelligence Scores | 5 | 11 | P0 |
| Smart Chat | 4 | 14 | P0 |
| Intelligent Scheduling | 4 | 13 | P1 |
| Workflow Automation | 4 | 12 | P1 |
| Inbox Intelligence | 4 | 11 | P1 |
| Operating Rhythm | 4 | 12 | P2 |
| Frontend | 4 | 14 | P1 |
| **Total** | **35** | **100** | - |

---

## Week-by-Week Plan

### Week 1 (P0 - Foundation)
- OMOS-301 to 306: Context Engine
- OMOS-310 to 314: Intelligence Scores
- OMOS-320, 321, 323: Smart Chat core

### Week 2 (P1 - Productivity)
- OMOS-322: Multi-step workflows
- OMOS-330 to 333: Intelligent Scheduling
- OMOS-340 to 343: Workflow Automation

### Week 3 (P1/P2 - Intelligence)
- OMOS-350 to 353: Inbox Intelligence
- OMOS-360 to 363: Operating Rhythm
- OMOS-370, 371: Dashboard start

### Week 4 (Frontend + Polish)
- OMOS-372, 373: Chat + Navigation
- Testing and bug fixes
- Production deployment

---

## Quick Import (CSV Format)

```csv
Title,Description,Priority,Labels,Estimate,Epic
"Context Engine Framework","Build Context Engine for universal AI awareness",1,"backend,ai,p0-critical",3,"Context Engine"
"State Aggregator","Implement real-time user state aggregation",2,"backend,ai",2,"Context Engine"
"Items Aggregator","Aggregate tasks, habits, goals into context",2,"backend,ai",2,"Context Engine"
"Scores Aggregator","Include HPI/LEI/GEI/UI in context",2,"backend,ai",2,"Context Engine"
"Calendar Aggregator","Schedule awareness in context",2,"backend,ai",2,"Context Engine"
"History Aggregator","Pattern awareness in context",2,"backend,ai",2,"Context Engine"
"HPI Calculator","Holistic Performance Index calculation",1,"backend,ai,p0-critical",3,"Intelligence Scores"
"LEI Calculator","Legacy Evolution Index calculation",2,"backend,ai",2,"Intelligence Scores"
"GEI Calculator","Generational Endeavor Index calculation",2,"backend,ai",2,"Intelligence Scores"
"UI Calculator","Unified Index calculation",2,"backend,ai",2,"Intelligence Scores"
"Score History & Trending","Score persistence and trend analysis",2,"backend,infrastructure",2,"Intelligence Scores"
"@Agent Routing","Route chat to specialist agents via @mentions",1,"backend,ai,p0-critical",3,"Smart Chat"
"/Commands System","Slash command system for quick actions",1,"backend,ai,p0-critical",3,"Smart Chat"
"Multi-Step Workflow from Chat","Create workflows from natural language",2,"backend,ai",5,"Smart Chat"
"Smart Chat Engine Integration","Unified chat engine",1,"backend,ai,p0-critical",3,"Smart Chat"
"DO DATE Algorithm","Motion-style intelligent scheduling",1,"backend,ai,p1-high",5,"Intelligent Scheduling"
"Auto-Reschedule","Automatic rescheduling on conflicts",2,"backend,ai",3,"Intelligent Scheduling"
"Energy Matching","Energy-aware task scheduling",2,"backend,ai",3,"Intelligent Scheduling"
"Path Balance Scheduling","HP/LE/GE distribution",2,"backend,ai",2,"Intelligent Scheduling"
"HTTP Triggers","Webhook triggers for workflows",1,"backend,infrastructure,p1-high",3,"Workflow Automation"
"Schedule Triggers (Cron)","Time-based workflow triggers",2,"backend,infrastructure",3,"Workflow Automation"
"Smart Approvals","agent_decides approval mode",2,"backend,ai",3,"Workflow Automation"
"Event Triggers","Event-based workflow triggers",2,"backend,infrastructure",3,"Workflow Automation"
"AI Triage Engine","AI-powered inbox classification",1,"backend,ai,p1-high",3,"Inbox Intelligence"
"Path Routing","Auto-assign HP/LE/GE paths",2,"backend,ai",2,"Inbox Intelligence"
"Task Extraction","Extract tasks from natural language",2,"backend,ai",3,"Inbox Intelligence"
"Quick Capture Processing","Process voice/image/text captures",2,"backend,ai",3,"Inbox Intelligence"
"Morning Briefing Workflow","Automated daily briefing",2,"backend,ai,p2-medium",3,"Operating Rhythm"
"Evening Review Workflow","Automated daily review",2,"backend,ai",3,"Operating Rhythm"
"Weekly Planning Automation","Weekly planning system",2,"backend,ai",3,"Operating Rhythm"
"Quarterly OKR Tracking","OKR tracking system",2,"backend",3,"Operating Rhythm"
"Command Center Dashboard","New home dashboard in Flutter",1,"frontend,p1-high",5,"Frontend"
"Score Widgets","Score visualization components",2,"frontend",3,"Frontend"
"Chat UI Enhancements","@mentions and /commands UI",2,"frontend",3,"Frontend"
"Navigation Restructure","New path-based navigation",2,"frontend",3,"Frontend"
```

---

*Last Updated: 2026-01-29*
*Project: OneMind Consciousness v2.0*
