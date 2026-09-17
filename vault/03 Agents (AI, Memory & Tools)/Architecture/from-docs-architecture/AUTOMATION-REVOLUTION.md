# OneMind Automation Revolution
## Combining Relevance AI + Pipedream + Motion + Agno for Human Optimization at Scale

*Research synthesis for building a revolutionary automation platform*

---

## Executive Summary

OneMind already has a powerful foundation with **Agno agents**, **workflows**, and **approvals**. By incorporating patterns from **Relevance AI** (workforce management), **Pipedream** (integration layer), **Motion** (intelligent scheduling), and **Acorn Land Labs** (gamification), we can build a revolutionary human optimization platform.

**Key Insight**: Relevance AI uses Pipedream as their integration backbone. We can do the same - or build our own integration layer using MCP + Composio.

---

## Research Findings

### 1. Relevance AI - Agent Workforce Platform

#### What They Do Well:
```
┌─────────────────────────────────────────────────────────────────┐
│                    RELEVANCE AI PATTERNS                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   AGENTS    │    │  WORKFORCE  │    │   TOOLS     │         │
│  │ ─────────── │    │ ─────────── │    │ ─────────── │         │
│  │ • Personas  │    │ • Teams     │    │ • No-code   │         │
│  │ • LLM swap  │    │ • Clone     │    │ • 2000+ int │         │
│  │ • Skills    │    │ • Observe   │    │ • Bulk ops  │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  APPROVALS  │    │  KNOWLEDGE  │    │  SCHEDULE   │         │
│  │ ─────────── │    │ ─────────── │    │ ─────────── │         │
│  │ • 3 modes   │    │ • Train 1x  │    │ • Recurring │         │
│  │ • Escalate  │    │ • Reuse all │    │ • Bulk      │         │
│  │ • Auto/Man  │    │ • Sync      │    │ • Wake mode │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│                                                                 │
│  KEY FEATURE: "Invent" - Describe agent → Get working agent    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Approval Modes (Critical Pattern):
1. **Approval Mode**: Tasks pause for human sign-off
2. **Agent Decision Mode**: Runs if confident, escalates if uncertain
3. **Autorun Mode**: Full automation, no human needed

#### What We Already Have in Agno:
- ✅ Agent factory with dynamic creation
- ✅ Teams (research, dev, finance)
- ✅ Workflow approval steps
- ✅ Knowledge bases with RAG
- ✅ Composio integration (actions)
- ⚠️ Missing: Agent personas with avatars
- ⚠️ Missing: "Invent" natural language agent creation
- ⚠️ Missing: Queue manager UI

---

### 2. Pipedream - Integration Layer

#### Architecture:
```
┌─────────────────────────────────────────────────────────────────┐
│                    PIPEDREAM ARCHITECTURE                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  TRIGGER ──► STEP 1 ──► STEP 2 ──► STEP N ──► OUTPUT           │
│                                                                 │
│  Step Types:                                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   TRIGGER   │  │   ACTION    │  │    CODE     │             │
│  │ ─────────── │  │ ─────────── │  │ ─────────── │             │
│  │ • HTTP      │  │ • Pre-built │  │ • Node.js   │             │
│  │ • Webhook   │  │ • 3000 apps │  │ • Python    │             │
│  │ • Schedule  │  │ • No-code   │  │ • Go/Bash   │             │
│  │ • Event     │  │             │  │             │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
│  Data Flow: steps.trigger.event → steps.step1.output → ...     │
│                                                                 │
│  VPC: Unique static IPs per workspace for secure access        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Key Features:
- **Event-driven**: Workflows execute on trigger events
- **Sequential steps**: Data flows from step to step
- **Props**: Reusable step configurations
- **CLI**: Local development and testing

#### What We Already Have in Agno:
- ✅ Workflow engine with steps
- ✅ Parallel, condition, router, loop steps
- ✅ Context passing between steps
- ✅ Function steps (custom Python)
- ⚠️ Missing: HTTP trigger endpoints
- ⚠️ Missing: Visual workflow builder
- ⚠️ Missing: Step marketplace/templates

---

### 3. Motion - Intelligent Scheduling

#### Core Concept:
```
┌─────────────────────────────────────────────────────────────────┐
│                    MOTION SCHEDULING                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│           DO DATE ≠ DUE DATE                                    │
│                                                                 │
│  ┌───────────┐     ┌───────────────────────────────────────┐   │
│  │  DUE DATE │     │              DO DATE                   │   │
│  │           │ →   │  AI calculates optimal execution time  │   │
│  │ 2024-02-15│     │  based on: energy, duration, deps,     │   │
│  │           │     │  calendar, deadlines, priority         │   │
│  └───────────┘     └───────────────────────────────────────┘   │
│                                                                 │
│  Features:                                                      │
│  • Context-aware AI chat with full data access                  │
│  • Multi-step workflow execution                                │
│  • Citations by default (show reasoning)                        │
│  • Human approval before destructive actions                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### What We Already Have:
- ✅ ONEMIND-CODEX.md defines the scheduling algorithm
- ⚠️ Missing: Implementation of DO DATE calculation
- ⚠️ Missing: Calendar integration for auto-scheduling

---

### 4. Acorn Land Labs - Gamification/Simulation

#### Interactive Elements:
```
┌─────────────────────────────────────────────────────────────────┐
│                    LAND LAB SIMULATOR                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                2D BIRDS-EYE VIEW                         │   │
│  │  ┌────┐  ┌────┐  ┌────┐                                 │   │
│  │  │🏠  │  │🌱  │  │💧  │  Interactive systems            │   │
│  │  │House│  │Garden│ │Water│                              │   │
│  │  └────┘  └────┘  └────┘                                 │   │
│  │                                                          │   │
│  │  Input/Output Tracking:                                  │   │
│  │  • Raw materials  • Energy  • Food  • Water              │   │
│  │                                                          │   │
│  │  AI Features:                                            │   │
│  │  • Sensor data integration (realtime)                    │   │
│  │  • Cost estimation                                       │   │
│  │  • Labor planning                                        │   │
│  │  • AI-powered suggestions/alerts                         │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### Application to OneMind:
- **Agent avatars/personas** like characters
- **Visual representation** of agent workforce
- **Real-time dashboards** showing agent activity
- **Gamified progress** toward PATH goals

---

## Architecture Mapping: OneMind Integration

```
┌─────────────────────────────────────────────────────────────────┐
│                 ONEMIND AUTOMATION STACK                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                       UIO CHAT                           │   │
│  │  (Motion-style: Context-aware, multi-step, citations)    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                            │                                    │
│                            ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   LEGACY ORCHESTRATOR                    │   │
│  │  (Relevance-style: Agent workforce management)           │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │                                                          │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │   │
│  │  │Researcher│ │Analyst  │ │Writer   │ │Coder    │       │   │
│  │  │   🔍    │ │   📊    │ │   ✍️    │ │   💻    │       │   │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │   │
│  │  │Finance  │ │House Mgr│ │Farm Mgr │ │Security │       │   │
│  │  │   💰    │ │   🏠    │ │   🌾    │ │   🛡️    │       │   │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │   │
│  │                                                          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                            │                                    │
│                            ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                   WORKFLOW ENGINE                        │   │
│  │  (Pipedream-style: Triggers → Steps → Actions)           │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │  Step Types:                                             │   │
│  │  [TRIGGER]→[AGENT]→[APPROVAL]→[FUNCTION]→[ACTION]       │   │
│  │  [PARALLEL] [CONDITION] [ROUTER] [LOOP] [TEAM]          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                            │                                    │
│                            ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 INTEGRATION LAYER                        │   │
│  │  (Pipedream/MCP: Connect to external apps)               │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │   │
│  │  │Composio │ │  MCP    │ │Pipedream│ │ Custom  │       │   │
│  │  │ Tools   │ │ Servers │ │ Connect │ │  APIs   │       │   │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Fast-Build Feature Roadmap

### Tier 1: Quick Wins (Days)
*Leverage existing Agno infrastructure*

| Feature | Effort | Impact | Implementation |
|---------|--------|--------|----------------|
| **Agent Personas** | 1 day | High | Add avatar, personality fields to agent config |
| **Approval Modes** | 1 day | High | Add "confident_autorun" mode to approval step |
| **Workflow Triggers** | 2 days | High | Add HTTP webhook trigger to workflow engine |
| **Queue Dashboard** | 2 days | Medium | Frontend: List pending tasks/approvals |
| **Agent Cloning** | 1 day | Medium | API: Clone agent with modifications |

### Tier 2: Medium Effort (1-2 Weeks)

| Feature | Effort | Impact | Implementation |
|---------|--------|--------|----------------|
| **Natural Language Agent Creation** | 1 week | Very High | "Invent" feature using LLM to generate agent config |
| **Bulk Scheduler** | 3 days | High | Schedule batch of items across agents |
| **Visual Workflow Builder** | 1 week | High | React Flow + Agno workflow models |
| **Wake Mode** | 3 days | Medium | Agent proactively schedules follow-ups |
| **Knowledge Sync** | 3 days | Medium | Auto-sync from Notion, Supabase |

### Tier 3: Revolutionary Features (2-4 Weeks)

| Feature | Effort | Impact | Implementation |
|---------|--------|--------|----------------|
| **Pipedream Integration** | 2 weeks | Very High | Connect to 3000+ apps via Pipedream |
| **Intelligent Scheduling** | 2 weeks | Very High | DO DATE algorithm from CODEX |
| **Agent Simulation View** | 2 weeks | High | Interactive visualization of workforce |
| **Multi-Agent Chains** | 1 week | High | Agent-to-agent handoff and feedback |
| **Escalation System** | 1 week | High | Smart routing to human when uncertain |

---

## Detailed Implementation Plans

### 1. Agent Personas (Quick Win)

```python
# backend/agno/agents/models.py - Add to AgentCreateRequest

class AgentPersona(BaseModel):
    """Visual identity and personality for agents."""
    avatar_url: Optional[str] = None
    emoji: Optional[str] = None  # Fallback icon
    personality_traits: Optional[List[str]] = None  # ["analytical", "creative"]
    communication_style: Optional[str] = None  # "formal", "casual", "concise"
    voice_description: Optional[str] = None  # For future voice agents

# Usage:
agent = create_agent_from_config({
    "name": "Research Agent",
    "persona": {
        "avatar_url": "/avatars/researcher.png",
        "emoji": "🔍",
        "personality_traits": ["curious", "thorough"],
        "communication_style": "academic"
    }
})
```

### 2. Smart Approval Modes (Quick Win)

```python
# backend/agno/workflows/workflow/steps/approval.py - Enhance

class SmartApproval(Approval):
    """Enhanced approval with confidence-based autorun."""

    def __init__(
        self,
        name: str,
        description: str,
        mode: Literal["always_approve", "agent_decides", "always_ask"] = "agent_decides",
        confidence_threshold: float = 0.8,
        risk_evaluator: Optional[str] = None,
        **kwargs
    ):
        super().__init__(name=name, description=description, **kwargs)
        self.mode = mode
        self.confidence_threshold = confidence_threshold
        self.risk_evaluator = risk_evaluator

    async def execute(self, context: WorkflowContext) -> StepResult:
        if self.mode == "always_approve":
            return StepResult(success=True, approved=True, content="Auto-approved")

        if self.mode == "agent_decides":
            # Evaluate confidence from previous step
            confidence = context.get("step_confidence", 0.5)
            risk_level = await self._evaluate_risk(context)

            if confidence >= self.confidence_threshold and risk_level == "low":
                return StepResult(success=True, approved=True,
                    content=f"Auto-approved (confidence: {confidence:.0%})")
            # Fall through to human approval

        return await super().execute(context)
```

### 3. HTTP Workflow Triggers (Quick Win)

```python
# backend/agno/workflows/triggers/http.py - New file

from fastapi import APIRouter, BackgroundTasks
from backend.agno.workflows.workflow.executor import WorkflowExecutor

router = APIRouter(prefix="/workflows/trigger", tags=["triggers"])

@router.post("/{workflow_id}")
async def trigger_workflow(
    workflow_id: str,
    background_tasks: BackgroundTasks,
    payload: dict = {},
):
    """
    HTTP trigger endpoint for workflows.
    POST /workflows/trigger/{workflow_id} with JSON body
    """
    executor = WorkflowExecutor()
    run_id = await executor.start_async(
        workflow_id=workflow_id,
        input=payload,
        trigger_type="http"
    )
    return {"run_id": run_id, "status": "started"}

@router.get("/{workflow_id}/webhook-url")
async def get_webhook_url(workflow_id: str):
    """Get the webhook URL for a workflow."""
    base_url = os.getenv("API_BASE_URL", "http://localhost:8081")
    return {"url": f"{base_url}/workflows/trigger/{workflow_id}"}
```

### 4. Natural Language Agent Creation ("Invent")

```python
# backend/agno/agents/invent.py - New feature

class AgentInventor:
    """Create agents from natural language descriptions."""

    INVENT_PROMPT = """
    Based on this description, generate an agent configuration:

    Description: {description}

    Generate a JSON config with:
    - name: Short agent name
    - role: What role this agent plays
    - instructions: Detailed system prompt
    - persona: avatar emoji, personality traits, communication style
    - tools: List of tool names from available tools
    - model: Recommended model (gpt-4o, claude-sonnet, etc)

    Available tools: {available_tools}
    """

    async def invent(self, description: str) -> AgentConfig:
        """Create an agent from natural language."""
        # Get available tools
        tools = load_all_tools()
        tool_names = [t.name for t in tools]

        # Ask LLM to generate config
        response = await get_model("gpt-4o").ainvoke(
            self.INVENT_PROMPT.format(
                description=description,
                available_tools=tool_names
            )
        )

        # Parse and validate
        config = json.loads(response.content)
        return AgentConfig(**config)

# API endpoint
@router.post("/agents/invent")
async def invent_agent(request: InventRequest):
    """Create agent from description: 'I need an agent that researches companies'"""
    inventor = AgentInventor()
    config = await inventor.invent(request.description)
    agent = create_and_register_agent(config)
    return AgentResponse.from_agent(agent)
```

### 5. Pipedream Integration Layer

```python
# backend/agno/integrations/pipedream.py - New integration

class PipedreamIntegration:
    """
    Use Pipedream as integration backbone.
    Similar to how Relevance AI uses it.
    """

    def __init__(self):
        self.api_key = os.getenv("PIPEDREAM_API_KEY")
        self.base_url = "https://api.pipedream.com/v1"

    async def trigger_workflow(
        self,
        workflow_url: str,
        payload: dict
    ) -> dict:
        """Trigger a Pipedream workflow."""
        async with httpx.AsyncClient() as client:
            response = await client.post(workflow_url, json=payload)
            return response.json()

    async def list_integrations(self) -> List[dict]:
        """List available Pipedream integrations."""
        # Could use Pipedream Connect API
        pass

    async def execute_action(
        self,
        app: str,
        action: str,
        params: dict
    ) -> dict:
        """Execute a Pipedream action."""
        # Use Pipedream's action execution API
        pass

# Tool wrapper for agents
class PipedreamTool(BaseTool):
    """Tool for agents to trigger Pipedream workflows."""
    name = "pipedream"
    description = "Trigger external workflows via Pipedream"

    async def _arun(self, workflow_url: str, payload: dict) -> str:
        pd = PipedreamIntegration()
        result = await pd.trigger_workflow(workflow_url, payload)
        return json.dumps(result)
```

### 6. Intelligent Scheduling (Motion-style)

```python
# backend/platform/scheduling/intelligent.py - New module

class IntelligentScheduler:
    """
    Motion-style scheduling: DO DATE ≠ DUE DATE
    AI calculates optimal execution time.
    """

    async def calculate_do_date(
        self,
        item: Item,
        user_id: str
    ) -> datetime:
        """Calculate optimal DO DATE for an item."""

        # Get user's calendar
        calendar = await self.get_user_calendar(user_id)

        # Get energy patterns
        energy_patterns = await self.get_energy_patterns(user_id)

        # Build scheduling prompt
        prompt = f"""
        Calculate optimal DO DATE for this item:

        Item: {item.title}
        Due Date: {item.due_date}
        Duration: {item.duration} minutes
        Energy Required: {item.energy_required}
        Priority: {item.priority}
        Dependencies: {item.dependencies}

        User's calendar availability: {calendar.get_free_slots()}
        User's peak energy times: {energy_patterns.peak_times}

        Rules:
        1. Schedule before due date with buffer
        2. Match energy requirements to user's energy patterns
        3. Respect blocked time
        4. Group similar tasks
        5. Consider dependencies

        Return optimal datetime in ISO format.
        """

        response = await get_model("gpt-4o-mini").ainvoke(prompt)
        return datetime.fromisoformat(response.content.strip())

    async def auto_schedule_inbox(self, user_id: str) -> int:
        """Auto-schedule all inbox items."""
        items = await self.get_inbox_items(user_id)
        scheduled = 0

        for item in items:
            if not item.do_date:
                item.do_date = await self.calculate_do_date(item, user_id)
                await self.update_item(item)
                scheduled += 1

        return scheduled
```

---

## Integration Strategy: MCP vs Pipedream

### Option A: Pipedream as Integration Backbone
```
OneMind Agents → Pipedream Workflows → 3000+ Apps
```
- **Pros**: Massive app library, managed auth, visual builder
- **Cons**: External dependency, cost at scale

### Option B: MCP + Composio Native
```
OneMind Agents → MCP Servers + Composio → Apps
```
- **Pros**: Full control, no external dependency, free
- **Cons**: Fewer integrations, more work

### Option C: Hybrid (Recommended)
```
OneMind Agents → MCP/Composio (core) + Pipedream (specialized)
```
- Use MCP for: Local tools, custom integrations, file ops
- Use Composio for: Common SaaS (GitHub, Slack, etc.)
- Use Pipedream for: Complex multi-app workflows, 1000+ apps

---

## Frontend Requirements

### Agent Workforce Dashboard
```
┌─────────────────────────────────────────────────────────────────┐
│  LEGACY WORKFORCE                               [+ New Agent]   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │      🔍         │  │      📊         │  │      ✍️         │ │
│  │   Researcher    │  │    Analyst      │  │     Writer      │ │
│  │                 │  │                 │  │                 │ │
│  │ Status: Idle    │  │ Status: Working │  │ Status: Waiting │ │
│  │ Tasks: 0        │  │ Tasks: 3        │  │ Tasks: 1        │ │
│  │ [Assign] [Clone]│  │ [View] [Pause]  │  │ [Approve]       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                 │
│  QUEUE                                           [View All →]   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ ⏳ Research competitor analysis      Analyst    Running   │  │
│  │ ⏳ Write blog post draft            Writer     Waiting   │  │
│  │ ✓ Summarize Q4 metrics             Analyst    Complete  │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  PENDING APPROVALS                                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ 🔔 Review research before publishing   [Approve] [Reject] │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Visual Workflow Builder
```
┌─────────────────────────────────────────────────────────────────┐
│  WORKFLOW: Morning Briefing                        [Save] [Run] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│    [HTTP Trigger] ──► [Researcher] ──► [Analyst] ──┐           │
│                                                     │           │
│                       ┌─────────────────────────────┘           │
│                       ▼                                         │
│                   [Approval] ──► [Writer] ──► [Email Action]   │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  STEPS:                                                         │
│  1. HTTP Trigger    - Webhook: /api/workflows/trigger/xyz       │
│  2. Researcher      - Agent: researcher, Prompt: "Research..."  │
│  3. Analyst         - Agent: analyst, Prompt: "Analyze..."      │
│  4. Approval        - Mode: agent_decides, Timeout: 60min       │
│  5. Writer          - Agent: writer, Prompt: "Write briefing"   │
│  6. Email Action    - Composio: gmail.send_email                │
│                                                                 │
│  [+ Add Step]  [+ Add Parallel]  [+ Add Condition]              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Priority Implementation Order

### Week 1: Foundation
1. ✅ Agent Personas (avatars, personality)
2. ✅ Smart Approval Modes
3. ✅ HTTP Workflow Triggers
4. ✅ Queue Dashboard API

### Week 2: Intelligence
5. Natural Language Agent Creation ("Invent")
6. Bulk Scheduler
7. Wake Mode (proactive scheduling)

### Week 3: Integration
8. Pipedream Integration (optional)
9. Knowledge Auto-Sync
10. Multi-Agent Chains

### Week 4: Experience
11. Visual Workflow Builder (React Flow)
12. Agent Workforce Dashboard
13. Intelligent Scheduling (DO DATE)

---

## Database Schema Additions

```sql
-- Agent personas
ALTER TABLE agents ADD COLUMN persona JSONB DEFAULT '{}';
-- Example: {"avatar_url": "...", "emoji": "🔍", "traits": ["analytical"]}

-- Workflow triggers
CREATE TABLE workflow_triggers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID REFERENCES workflows(id),
    trigger_type VARCHAR(50) NOT NULL,  -- http, schedule, event
    config JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Task queue
CREATE TABLE task_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id VARCHAR(100) NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    input JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    result JSONB,
    error TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Intelligent scheduling
ALTER TABLE items ADD COLUMN do_date TIMESTAMPTZ;
ALTER TABLE items ADD COLUMN scheduling_metadata JSONB DEFAULT '{}';
-- Example: {"energy_match": 0.9, "calendar_slot": "...", "reason": "..."}
```

---

## Conclusion

OneMind already has a powerful foundation. By adding:

1. **Agent Personas** → Interactive, engaging workforce
2. **Smart Approvals** → Balanced automation with control
3. **HTTP Triggers** → External system integration
4. **Invent Feature** → Zero-code agent creation
5. **Intelligent Scheduling** → Motion-style DO DATE
6. **Visual Builder** → Pipedream-style workflow creation

We can create a revolutionary human optimization platform that combines the best of:
- **Relevance AI**: Workforce management, approvals
- **Pipedream**: Integration layer, workflows
- **Motion**: Intelligent scheduling
- **Acorn Labs**: Gamification, simulation

The revolution starts with what you already have. Build incrementally, ship fast.
