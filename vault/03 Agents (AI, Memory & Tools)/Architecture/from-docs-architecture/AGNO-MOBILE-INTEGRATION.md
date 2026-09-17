# Agno Framework Integration for Mobile-First OneMind

> **How Every Agno Feature Powers the Mobile Experience**
> **Last Updated:** 2026-01-29

---

## Overview

This document maps every Agno framework feature to its role in the mobile-first OneMind OS experience. Agno is the backbone powering intelligent agents, and understanding this connection is critical for the build.

---

## Part 1: Agno Core Features → Mobile UI

### 1. Agent Reasoning (`reasoning=True`)

**Agno Feature:**
```python
Agent(
    name="Legacy",
    reasoning=True,  # Enable multi-step reasoning
    ...
)
```

**Mobile Integration:**
- Powers the **Command Chat** with complex query handling
- Shows **"Thinking..."** indicator when reasoning
- Enables multi-step task breakdown in chat

**UI Component:**
```dart
// command_screen.dart
if (agentState.isReasoning) {
  ThinkingLoader(text: "Legacy is reasoning...")
}
```

---

### 2. Agent Memory (`enable_user_memories`, `enable_agentic_memory`)

**Agno Feature:**
```python
Agent(
    enable_user_memories=True,    # Remember user facts
    enable_agentic_memory=True,   # Agent self-learning
    enable_session_summaries=True, # Session context
    ...
)
```

**Mobile Integration:**
- **Command Chat**: Context-aware responses based on memory
- **Capture Processing**: Auto-categorization from learned preferences
- **Suggestions**: Personalized based on patterns

**UI Components:**
```dart
// Memory indicator in chat
MemoryIndicator(
  memories: context.recentMemories,
  onTap: () => showMemoryDetails(),
)

// Knowledge screen
KnowledgeScreen(
  userMemories: memoryService.getUserMemories(),
  agentMemories: memoryService.getAgentMemories(),
)
```

---

### 3. Knowledge Base (RAG)

**Agno Feature:**
```python
Agent(
    knowledge_base=LanceDBKnowledgeBase(
        table_name="shared-knowledge",
        search_type=SearchType.vector,
    ),
    search_knowledge=True,
    ...
)
```

**Mobile Integration:**
- **Knowledge Screen**: Browse and search all documents
- **Command Chat**: Semantic search through `/search` command
- **Capture**: Auto-link to relevant knowledge

**UI Component:**
```dart
// Knowledge search in Command
SearchResults(
  query: userQuery,
  results: knowledgeService.semanticSearch(userQuery),
  onSelect: (doc) => showDocument(doc),
)
```

---

### 4. Tools Registry

**Agno Feature:**
```python
from backend.agno.tools.registry import get_tool, list_tools

# Available tools
tools = [
    "todoist",        # Task management
    "github",         # Code operations
    "home_assistant", # Smart home
    "finance",        # Financial data
    "browser",        # Web browsing
    "comms",          # Communications
    ...
]
```

**Mobile Integration:**
- **Execute Screen**: Actions powered by tools
- **Capture**: Routes to appropriate tool
- **Quick Actions**: Tool-specific shortcuts

**Tool → UI Mapping:**

| Tool | Mobile Feature |
|------|----------------|
| `todoist` | Task creation, completion in Execute |
| `life_os` | Habits, routines, meals tracking |
| `home_assistant` | Home status card on Home screen |
| `finance` | WI score calculation, finance insights |
| `comms` | Notification badges, quick reply |
| `browser` | Research results in chat |

---

### 5. Skills System (Progressive Disclosure)

**Agno Feature:**
```python
from backend.agno.skills import get_skill_instructions

# Skills loaded on-demand
skills = [
    "shared/three-paths",      # HP/LE/GE framework
    "shared/communication",    # How to communicate
    "specialist/task-planning", # Task breakdown
    "domain/smart-home",       # Home automation
    ...
]
```

**Mobile Integration:**
- **Agent Personas**: Different skill sets per agent
- **Suggestions**: Skill-aware recommendations
- **Command Parsing**: Skill-specific /commands

**Skills → UI:**

| Skill | UI Impact |
|-------|-----------|
| `three-paths` | Pillar badges (HP/LE/GE) everywhere |
| `task-planning` | Smart task breakdown suggestions |
| `smart-home` | Home Assistant widget integration |
| `code-review` | Code-related suggestions |

---

### 6. Teams & Coordination

**Agno Feature:**
```python
from backend.agno.teams import create_team, TeamMode

team = Team(
    name="Research Team",
    mode=TeamMode.COLLABORATE,  # Or COORDINATE, DELEGATE, ROUTE, SEQUENTIAL
    members=[researcher, analyst, writer],
)
```

**Mobile Integration:**
- **Active Agents Widget**: Shows which team is working
- **Approvals Screen**: Team decisions needing input
- **Insights Screen**: Results from team analysis

**Team Visibility:**
```dart
// Home screen widget
ActiveAgentsCard(
  activeTeams: agentService.getActiveTeams(),
  onTeamTap: (team) => showTeamDetails(team),
  children: [
    TeamRow(
      name: "Research Team",
      status: "Analyzing market data",
      members: 3,
    ),
  ],
)
```

---

### 7. Workflows

**Agno Feature:**
```python
from backend.agno.workflows import WorkflowExecutor

workflow = Workflow(
    steps=[
        Step(type="AGENT", agent="researcher"),
        Step(type="PARALLEL", agents=["analyst", "writer"]),
        Step(type="APPROVAL", message="Review findings?"),
    ]
)
```

**Mobile Integration:**
- **Operating Rhythm**: Automated workflows (morning briefing)
- **Workflow Status**: Progress indicators
- **Approvals**: In-line approval requests

**Workflow Types:**

| Workflow | Mobile Trigger | SO |
|----------|---------------|-----|
| Morning Briefing | 6am auto | Push notification + Home card |
| Task Triage | New capture | Background, results in Execute |
| Research | @research mention | Progress in Command chat |
| Evening Review | 9pm auto | Summary card on Home |

---

### 8. Approvals (Human-in-the-Loop)

**Agno Feature:**
```python
from backend.agno.approvals import ApprovalService

# Create approval request
await approval_service.create_approval(
    action="delete_all_tasks",
    risk_level="high",
    context={"count": 50},
    requires_confirmation=True,
)
```

**Mobile Integration:**
- **Push Notifications**: Approval requests
- **Approval Badge**: Count on Home screen
- **Quick Actions**: Approve/Reject/Modify inline

**Approval UI:**
```dart
// Approval modal
ApprovalModal(
  action: "Delete 50 tasks",
  risk: RiskLevel.high,
  context: approvalContext,
  onApprove: () => approvalService.approve(id),
  onReject: () => approvalService.reject(id),
  onModify: () => showModifySheet(id),
)
```

---

### 9. Context Provider

**Agno Feature:**
```python
from backend.agno.agents.context import get_dynamic_context

context = await get_dynamic_context()
# Returns: location, time, device, sun_state, etc.
```

**Mobile Integration:**
- **Home Screen**: Location-aware greeting
- **Command Chat**: Context in every response
- **Suggestions**: Context-aware recommendations

**Context → UI:**

| Context | UI Display |
|---------|-----------|
| Location | "Good morning Zeus, you're at Home" |
| Time | Different suggestions AM vs PM |
| Sun state | Dark/light theme suggestion |
| Device presence | Show relevant device controls |

---

### 10. Structured Outputs (`response_model`)

**Agno Feature:**
```python
from pydantic import BaseModel

class TaskBreakdown(BaseModel):
    title: str
    subtasks: List[str]
    estimated_hours: float
    pillar: str

agent = Agent(
    response_model=TaskBreakdown,
    ...
)
```

**Mobile Integration:**
- **Capture Processing**: Structured asset creation
- **Smart Suggestions**: Typed suggestion objects
- **Analytics**: Structured insights data

**Example Flow:**
```
User captures: "Plan the product launch"
    ↓
Agent returns TaskBreakdown:
{
  "title": "Product Launch Plan",
  "subtasks": [
    "Create launch timeline",
    "Prepare marketing materials",
    "Set up analytics",
    "Draft announcement"
  ],
  "estimated_hours": 8.0,
  "pillar": "GE"
}
    ↓
UI shows structured task creation dialog
```

---

## Part 2: Agno Agent Types → Mobile Screens

### Mother Agent (Legacy)

**Configuration:**
```python
legacy = Agent(
    agent_id="legacy-mother-ai",
    name="Legacy",
    model=get_model("claude-sonnet"),
    reasoning=True,
    enable_user_memories=True,
    enable_agentic_memory=True,
    tools=LEGACY_TOOLS,
    knowledge_base=shared_knowledge,
)
```

**Primary Mobile Screen:** Command

**Features:**
- Main chat interface
- Routes to specialists via @mentions
- Executes /commands
- Provides suggestions
- Manages approvals

---

### Specialist Agents

| Agent | Model | Primary Screen | Mobile Features |
|-------|-------|---------------|-----------------|
| Planning Agent | Sonnet | Execute | Task breakdown, scheduling suggestions |
| Research Agent | Sonnet | Command | Research results, citations |
| Code Agent | Sonnet | Command | Code snippets, PR summaries |
| Analysis Agent | Sonnet | Insights | Trends, correlations |
| Health Agent | Sonnet | Home | PI score, health suggestions |
| Finance Agent | Opus | Insights | WI score, investment alerts |
| Writer Agent | Sonnet | Command | Draft content |

---

### Teams

| Team | Members | Mobile Trigger | Output Screen |
|------|---------|---------------|---------------|
| Research Team | Researcher + Analyst + Writer | @research | Command |
| C-Suite | CEO + CFO + COO + CTO | Strategic queries | Insights |
| Dev Team | Coder + Analyst | @code + analysis | Command |
| Content Team | Writer + Researcher + Editor | @write + context | Command |

---

## Part 3: Agno Features → Linear Issues

### Feature: Agent Reasoning

**Linear Issues:**
- OMOS-320: @Agent Routing (3pts) - Route messages to specialist agents
- OMOS-323: Smart Chat Engine Integration (3pts)

**Mobile Implementation:**
- Reasoning indicator in chat
- Multi-step progress display
- Result streaming

---

### Feature: Memory System

**Linear Issues:**
- OMOS-455: Memory Blocks with Semantic Search (3pts)
- New: OMOS-486: Memory UI Component (3pts)

**Mobile Implementation:**
- Memory indicator badge
- Memory browser screen
- Auto-memory creation from capture

---

### Feature: Knowledge Base

**Linear Issues:**
- Existing RAG issues
- New: OMOS-487: Knowledge Mobile UI (3pts)

**Mobile Implementation:**
- Document list view
- Search interface
- Upload capability

---

### Feature: Tools

**Linear Issues:**
- All existing tool issues remain
- New: OMOS-488: Tool Status Dashboard (2pts)

**Mobile Implementation:**
- Tool integration status in Profile
- Tool-specific quick actions

---

### Feature: Skills

**Linear Issues:**
- Existing skill system
- New: OMOS-489: Skill Visibility (2pts)

**Mobile Implementation:**
- Agent skill badges
- Skill-aware suggestions

---

### Feature: Teams

**Linear Issues:**
- Existing team issues
- OMOS-477: Agent Visibility Layer (3pts)

**Mobile Implementation:**
- Active teams widget
- Team progress indicators

---

### Feature: Workflows

**Linear Issues:**
- OMOS-340-343: Workflow Automation
- OMOS-360-363: Operating Rhythm

**Mobile Implementation:**
- Workflow status cards
- Schedule triggers
- Push notifications

---

### Feature: Approvals

**Linear Issues:**
- OMOS-342: Smart Approvals (3pts)
- New: OMOS-490: Mobile Approval UI (3pts)

**Mobile Implementation:**
- Approval push notifications
- Inline approval actions
- Risk level badges

---

### Feature: Context

**Linear Issues:**
- OMOS-301-306: Context Engine
- OMOS-302: State Aggregator

**Mobile Implementation:**
- Context indicator in header
- Location-aware features
- Time-based suggestions

---

### Feature: Structured Outputs

**Linear Issues:**
- Part of capture processing
- New: OMOS-491: Structured Asset Creation (3pts)

**Mobile Implementation:**
- Type-safe capture parsing
- Structured suggestion display

---

## Part 4: Complete Agno → Mobile Feature Map

```
┌─────────────────────────────────────────────────────────┐
│                    AGNO FRAMEWORK                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  AGENTS                    → COMMAND SCREEN             │
│  ├── Legacy (Mother)       │  Main chat interface      │
│  ├── Specialists           │  @mention routing         │
│  └── Teams                 │  Team collaboration       │
│                                                          │
│  TOOLS                     → EXECUTE SCREEN            │
│  ├── todoist               │  Task management          │
│  ├── life_os               │  Habits, routines         │
│  ├── home_assistant        │  Home status              │
│  └── finance               │  WI tracking              │
│                                                          │
│  MEMORY                    → KNOWLEDGE SCREEN          │
│  ├── User memories         │  Preference storage       │
│  ├── Agent memories        │  Learned patterns         │
│  └── Session summaries     │  Conversation context     │
│                                                          │
│  KNOWLEDGE BASE            → KNOWLEDGE SCREEN          │
│  ├── Documents             │  Document browser         │
│  ├── Semantic search       │  Search interface         │
│  └── RAG retrieval         │  Context injection        │
│                                                          │
│  WORKFLOWS                 → HOME SCREEN               │
│  ├── Morning briefing      │  Daily summary card       │
│  ├── Task triage           │  Background processing    │
│  └── Evening review        │  Day summary              │
│                                                          │
│  APPROVALS                 → APPROVALS MODAL           │
│  ├── HITL requests         │  Push notifications       │
│  ├── Risk assessment       │  Risk badges              │
│  └── Quick actions         │  Inline approve/reject    │
│                                                          │
│  CONTEXT                   → EVERYWHERE                │
│  ├── Location              │  Location indicator       │
│  ├── Time                  │  Time-aware suggestions   │
│  ├── Device state          │  Relevant controls        │
│  └── User state            │  Energy-based scheduling  │
│                                                          │
│  SKILLS                    → AGENT CARDS               │
│  ├── Three Paths           │  Pillar badges           │
│  ├── Domain skills         │  Specialist routing       │
│  └── Progressive load      │  On-demand expertise      │
│                                                          │
│  SCORES                    → HOME DASHBOARD            │
│  ├── PI (Health)           │  HP pillar score          │
│  ├── LI (Legacy)           │  LE pillar score          │
│  ├── WI (Wealth)           │  GE pillar score          │
│  └── SO (Sovereign)          │  Master score             │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Part 5: Implementation Checklist

### Backend (Agno Enhancements)

- [ ] Universal Asset model integrates with all tools
- [ ] Context Engine aggregates all data sources
- [ ] @mention parser routes to correct agent
- [ ] /command registry with all quick actions
- [ ] Mother Agent orchestration for complex queries
- [ ] Memory persistence for user preferences
- [ ] Approval queue with push notification triggers
- [ ] Workflow scheduler for operating rhythm
- [ ] Score calculation from aggregated data
- [ ] Structured output models for capture

### Frontend (Mobile UI)

- [ ] Command screen with full chat features
- [ ] Execute screen with tool-powered actions
- [ ] Home dashboard with all widgets
- [ ] Capture modal with voice/photo
- [ ] Approval modal with risk display
- [ ] Active agents widget
- [ ] Memory/knowledge browser
- [ ] Score display components
- [ ] Context indicator
- [ ] Notification integration

---

## Summary

Every Agno feature has a direct mobile UI manifestation:

| Agno Feature | Mobile Screen | Primary Component |
|--------------|--------------|-------------------|
| Agent Reasoning | Command | ThinkingLoader |
| Memory | Knowledge | MemoryList |
| Knowledge Base | Knowledge | DocumentBrowser |
| Tools | Execute | ActionButtons |
| Skills | Command | AgentBadge |
| Teams | Home | ActiveTeamsCard |
| Workflows | Home | WorkflowStatusCard |
| Approvals | Modal | ApprovalSheet |
| Context | Header | ContextIndicator |
| Scores | Home | ScoreCard |

This tight integration ensures Agno's power is fully accessible on mobile.
