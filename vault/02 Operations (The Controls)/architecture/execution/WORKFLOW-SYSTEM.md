# OneMind OS Workflow System

> Agent-orchestrated workflows with triggers, HITL, and templates

## Overview

The workflow system enables multi-step AI agent orchestration with:
- **8 Step Types**: Sequential, parallel, conditional, routing, loops, approvals, functions, teams
- **3 Trigger Types**: Schedule (cron), Event (NATS), HTTP (webhooks)
- **HITL Support**: Human-in-the-loop approval gates
- **Templates**: Pre-built workflow patterns

## Architecture

```
backend/agno/workflows/workflow/
├── models.py          # Pydantic models for all workflow entities
├── executor.py        # Workflow execution engine
├── builder.py         # JSON → Step instance converter
├── registry.py        # Workflow storage and lookup
├── storage.py         # PostgreSQL persistence
├── templates.py       # Pre-built workflow templates
├── api.py             # Custom FastAPI endpoints
├── steps/             # Step type implementations
│   ├── base.py        # Abstract base step
│   ├── step.py        # Single agent step
│   ├── parallel.py    # Concurrent execution
│   ├── condition.py   # Conditional branching
│   ├── router.py      # Dynamic step selection
│   ├── loop.py        # Iterative execution
│   ├── approval.py    # HITL gates
│   ├── function.py    # Custom Python functions
│   └── team.py        # Team coordination
└── triggers/          # Automatic workflow execution
    ├── models.py      # Trigger models
    ├── schedule.py    # Cron-based triggers
    ├── event.py       # NATS event triggers
    ├── http.py        # Webhook triggers
    └── api.py         # Trigger CRUD endpoints
```

## Step Types

### 1. Step (Basic Agent Execution)

Execute a single agent with a prompt.

```json
{
  "type": "step",
  "name": "Research Topic",
  "agent_id": "researcher",
  "prompt_template": "Research: {topic}",
  "output_key": "research_results",
  "timeout_seconds": 120
}
```

**Features**:
- Input caching (OMOS-292) - skip re-execution on cache hit
- Step reference resolution: `{previous_step.field}`
- Timeout and retry configuration

### 2. Parallel (Concurrent Execution)

Run multiple steps simultaneously.

```json
{
  "type": "parallel",
  "name": "Gather Data",
  "steps": [
    {"type": "step", "name": "Calendar", "agent_id": "calendar-agent"},
    {"type": "step", "name": "Tasks", "agent_id": "task-agent"},
    {"type": "step", "name": "Weather", "agent_id": "weather-agent"}
  ],
  "fail_fast": true,
  "merge_outputs": true
}
```

**Features**:
- `fail_fast`: Stop all steps if one fails
- `merge_outputs`: Combine all step outputs into single dict

### 3. Condition (Branching)

Route to different steps based on evaluation.

```json
{
  "type": "condition",
  "name": "Priority Router",
  "evaluator": "priority",
  "branches": {
    "high": {"type": "step", "agent_id": "urgent-handler"},
    "medium": {"type": "step", "agent_id": "standard-handler"},
    "low": {"type": "step", "agent_id": "batch-handler"}
  },
  "default_branch": "medium"
}
```

**Built-in Evaluators**:
| Evaluator | Returns | Description |
|-----------|---------|-------------|
| `priority` | high/medium/low | Keyword-based priority detection |
| `task_type` | code/writing/research/review/general | Task classification |
| `pillar` | hp/le/ge/general | Three Paths framework |
| `sentiment` | positive/negative/neutral | Sentiment analysis |

**Custom Evaluator**:
```python
from backend.agno.workflows.workflow.steps.condition import register_evaluator

@register_evaluator("custom_eval")
def my_evaluator(context: dict) -> str:
    if context.get("urgent"):
        return "high"
    return "normal"
```

### 4. Router (Dynamic Step Selection)

Select one or more steps to execute based on input.

```json
{
  "type": "router",
  "name": "Expert Selector",
  "selector": "multi_expert",
  "choices": [
    {"type": "step", "name": "Code Expert", "agent_id": "coder"},
    {"type": "step", "name": "Writing Expert", "agent_id": "writer"},
    {"type": "step", "name": "Research Expert", "agent_id": "researcher"}
  ],
  "fallback": {"type": "step", "agent_id": "generalist"}
}
```

**Built-in Selectors**:
| Selector | Behavior |
|----------|----------|
| `topic` | Match keywords in step names |
| `expertise` | Route to single best expert |
| `multi_expert` | Select multiple relevant experts |

### 5. Loop (Iterative Execution)

Repeat steps until condition is met.

```json
{
  "type": "loop",
  "name": "Iterative Refinement",
  "steps": [
    {"type": "step", "name": "Draft", "agent_id": "writer"},
    {"type": "step", "name": "Review", "agent_id": "reviewer"}
  ],
  "end_condition": "quality_achieved",
  "max_iterations": 5,
  "min_iterations": 2
}
```

**Built-in End Conditions**:
| Condition | Stops When |
|-----------|------------|
| `word_count_500` | Content has 500+ words |
| `word_count_1000` | Content has 1000+ words |
| `has_conclusion` | Content contains conclusion markers |
| `quality_achieved` | `quality_achieved=True` in context |

### 6. Approval (HITL Gate)

Pause workflow for human approval.

```json
{
  "type": "approval",
  "name": "Review Before Publish",
  "description": "Please review the content before publishing",
  "timeout_minutes": 60,
  "auto_approve_mode": "low_risk",
  "require_notes": false
}
```

**Auto-Approve Modes** (OMOS-437):
| Mode | Behavior |
|------|----------|
| `null` | Always require human approval |
| `"low_risk"` | Auto-approve if LLM assesses as low risk |
| `"always"` | Skip approval entirely |

**Approval Flow**:
1. Workflow pauses at approval step
2. Notification sent to inbox
3. User approves/rejects via UI or API
4. Workflow resumes or fails

### 7. Function (Custom Python)

Execute custom Python logic.

```json
{
  "type": "function",
  "name": "Process Data",
  "executor": "transform_data",
  "output_key": "processed_data"
}
```

**Built-in Executors**:
| Executor | Description |
|----------|-------------|
| `uppercase` | Transform content to uppercase |
| `word_count` | Count words in content |
| `extract_summary` | Extract first 500 chars |
| `set_quality_achieved` | Set flag for loop termination |

**Custom Executor**:
```python
from backend.agno.workflows.workflow.steps.function import register_executor

@register_executor("my_transform")
def my_transform(context: dict) -> dict:
    data = context.get("raw_data", "")
    return {"transformed": data.upper(), "length": len(data)}
```

### 8. Team (Multi-Agent Coordination)

Execute an Agno team with coordination mode.

```json
{
  "type": "team",
  "name": "Dev Team Review",
  "team_id": "dev-team",
  "mode": "coordinate",
  "prompt_template": "Review this code: {code}",
  "output_key": "review_results",
  "include_member_responses": true
}
```

**Team Modes**:
| Mode | Description |
|------|-------------|
| `coordinate` | Leader assigns tasks to specific members |
| `collaborate` | All members work on same task, share context |
| `delegate` | Leader picks one member, returns directly |
| `route` | Leader routes to best member(s) |
| `sequential` | Members execute in order, chain outputs |

---

## Triggers

### Schedule Triggers (Cron)

Execute workflows on a schedule using APScheduler.

```python
# API: POST /triggers/schedule
{
    "workflow_id": "morning-briefing",
    "name": "Daily Morning Briefing",
    "cron_expression": "0 6 * * *",
    "timezone": "America/Chicago",
    "input_template": {"user_id": "zeus"}
}
```

**Cron Presets**:
| Preset | Cron | Description |
|--------|------|-------------|
| `morning_briefing` | `0 6 * * *` | 6 AM daily |
| `evening_review` | `0 21 * * *` | 9 PM daily |
| `weekly_planning` | `0 20 * * 0` | Sunday 8 PM |
| `quarterly_okr` | `0 9 1 1,4,7,10 *` | First of quarter |
| `hourly_check` | `0 * * * *` | Every hour |
| `every_15_minutes` | `*/15 * * * *` | Every 15 min |

**API Endpoints**:
```
POST   /triggers/schedule              # Create trigger
POST   /triggers/schedule/from-preset  # Create from preset
GET    /triggers/schedule              # List triggers
GET    /triggers/schedule/{id}         # Get trigger
PATCH  /triggers/schedule/{id}         # Update trigger
DELETE /triggers/schedule/{id}         # Delete trigger
POST   /triggers/schedule/{id}/pause   # Pause trigger
POST   /triggers/schedule/{id}/resume  # Resume trigger
GET    /triggers/schedule/upcoming     # Get upcoming runs
```

### Event Triggers (NATS)

Execute workflows when NATS events occur.

```python
# API: POST /triggers/event
{
    "workflow_id": "task-triage",
    "name": "Triage New Tasks",
    "subject_pattern": "inbox.item.created",
    "filters": {"type": "task"},
    "debounce_seconds": 5,
    "input_mapping": {
        "item_id": "$.id",
        "title": "$.title"
    }
}
```

**Built-in Event Patterns**:
| Pattern | Subject | Filters |
|---------|---------|---------|
| `inbox_task_created` | `inbox.item.created` | `{type: task}` |
| `inbox_approval_created` | `inbox.item.created` | `{type: approval}` |
| `habit_completed` | `lifeos.habit.completed` | - |
| `habit_missed` | `lifeos.habit.missed` | - |
| `goal_progress_updated` | `lifeos.goal.updated` | `{progress_delta_gte: 10}` |
| `task_completed` | `lifeos.task.completed` | - |
| `score_changed` | `lifeos.score.updated` | - |

**Filter Operators**:
- `key: value` - Exact match
- `key_gte: value` - Greater than or equal
- `key_lte: value` - Less than or equal
- `key_in: [values]` - In list

**API Endpoints**:
```
POST   /triggers/event               # Create trigger
POST   /triggers/event/from-builtin  # Create from pattern
GET    /triggers/event               # List triggers
GET    /triggers/event/{id}          # Get trigger
PATCH  /triggers/event/{id}          # Update trigger
DELETE /triggers/event/{id}          # Delete trigger
POST   /triggers/event/{id}/pause    # Pause (unsubscribe)
POST   /triggers/event/{id}/resume   # Resume (resubscribe)
```

### HTTP Triggers (Webhooks)

Execute workflows via external HTTP requests.

```python
# API: POST /triggers/http
{
    "workflow_id": "github-integration",
    "name": "GitHub Webhook",
    "allowed_methods": ["POST"],
    "require_signature": true,
    "transform_template": "{\"pr_number\": {{ pull_request.number }}}"
}

# Response:
{
    "trigger_id": "abc123",
    "url": "/hooks/xyz789",
    "secret": "your-hmac-secret",
    "instructions": "POST with X-Webhook-Signature header"
}
```

**Signature Verification**:
```bash
# Generate signature
echo -n '{"event": "push"}' | openssl dgst -sha256 -hmac "your-secret"

# Send request
curl -X POST https://api.example.com/hooks/xyz789 \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: sha256=abc123..." \
  -d '{"event": "push"}'
```

**API Endpoints**:
```
POST   /triggers/http                       # Create webhook
GET    /triggers/http                       # List webhooks
GET    /triggers/http/{id}                  # Get webhook
PATCH  /triggers/http/{id}                  # Update webhook
DELETE /triggers/http/{id}                  # Delete webhook
POST   /triggers/http/{id}/regenerate-secret # New secret

# Webhook receiver (external calls hit this):
POST   /hooks/{token}                       # Trigger workflow
```

---

## Templates

Pre-built workflow patterns available via API.

```python
# List templates
GET /workflows/templates

# Create from template
POST /workflows/from-template
{
    "template_id": "research",
    "name": "My Research Workflow",
    "agent_overrides": {
        "researcher": "my-custom-researcher"
    }
}
```

**Built-in Templates**:

| Template | Steps | Use Case |
|----------|-------|----------|
| `research` | Research → Write → Review | Information gathering |
| `code_review` | Parallel review (security, style, logic) | Code analysis |
| `task_triage` | Condition-based routing by priority | Task management |
| `morning_briefing` | Parallel data → Synthesize | Daily briefing |
| `iterative_research` | Loop until quality achieved | Deep research |
| `content_approval` | Create → Approval → Publish | Content workflow |
| `expert_routing` | Router to domain experts | Complex queries |

---

## API Reference

### Native AgentOS Endpoints

```
POST   /workflows                    # Create workflow
GET    /workflows                    # List workflows
GET    /workflows/{id}               # Get workflow
DELETE /workflows/{id}               # Delete workflow
POST   /workflows/{id}/runs          # Run workflow
DELETE /workflows/{id}/runs/{run_id} # Cancel run
```

### Custom OneMind Endpoints

```
PUT    /workflows/{id}                      # Update workflow
PUT    /workflows/{id}/activate             # Enable/disable
POST   /workflows/{id}/runs/{run_id}/resume # Resume paused (HITL)
GET    /workflows/templates                 # List templates
GET    /workflows/templates/{id}            # Get template
POST   /workflows/from-template             # Create from template
GET    /workflows/stats                     # Statistics
GET    /workflows/active-runs               # Active runs
GET    /workflows/step-types                # UI builder schema
```

### Unified Trigger Endpoints

```
GET    /triggers                     # List all triggers
GET    /triggers/presets             # Get cron/event presets
GET    /triggers/stats               # Trigger statistics
```

---

## Database Schema

```sql
-- Workflow definitions
CREATE TABLE workflows (
    id UUID PRIMARY KEY,
    workflow_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    definition JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    ...
);

-- Workflow runs
CREATE TABLE workflow_runs (
    id UUID PRIMARY KEY,
    run_id VARCHAR(100) UNIQUE NOT NULL,
    workflow_id VARCHAR(100) REFERENCES workflows,
    status VARCHAR(50) DEFAULT 'pending',
    input JSONB,
    output JSONB,
    context JSONB DEFAULT '{}',
    ...
);

-- Step runs
CREATE TABLE workflow_step_runs (
    id UUID PRIMARY KEY,
    run_id VARCHAR(100) REFERENCES workflow_runs,
    step_name VARCHAR(255) NOT NULL,
    step_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    ...
);

-- Triggers
CREATE TABLE workflow_triggers (
    id UUID PRIMARY KEY,
    workflow_id VARCHAR(100) REFERENCES workflows,
    trigger_type VARCHAR(20) NOT NULL,  -- schedule, event, http
    config JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    ...
);
```

---

## Usage Examples

### Example 1: Morning Briefing Workflow

```json
{
  "name": "Morning Briefing",
  "steps": [
    {
      "type": "parallel",
      "name": "Gather Data",
      "steps": [
        {
          "type": "step",
          "name": "Calendar",
          "agent_id": "calendar-agent",
          "prompt_template": "Get today's events for {user_id}",
          "output_key": "calendar"
        },
        {
          "type": "step",
          "name": "Tasks",
          "agent_id": "task-agent",
          "prompt_template": "Get priority tasks for {user_id}",
          "output_key": "tasks"
        },
        {
          "type": "step",
          "name": "Weather",
          "agent_id": "weather-agent",
          "prompt_template": "Get weather for {location}",
          "output_key": "weather"
        }
      ]
    },
    {
      "type": "step",
      "name": "Synthesize",
      "agent_id": "briefing-agent",
      "prompt_template": "Create morning briefing from:\nCalendar: {calendar}\nTasks: {tasks}\nWeather: {weather}"
    }
  ]
}
```

### Example 2: Content Pipeline with Approval

```json
{
  "name": "Content Pipeline",
  "steps": [
    {
      "type": "step",
      "name": "Research",
      "agent_id": "researcher",
      "prompt_template": "Research: {topic}",
      "output_key": "research"
    },
    {
      "type": "step",
      "name": "Write",
      "agent_id": "writer",
      "prompt_template": "Write article based on: {research}",
      "output_key": "draft"
    },
    {
      "type": "approval",
      "name": "Editorial Review",
      "description": "Review article before publishing",
      "timeout_minutes": 120
    },
    {
      "type": "step",
      "name": "Publish",
      "agent_id": "publisher",
      "prompt_template": "Publish: {draft}"
    }
  ]
}
```

### Example 3: Automated with Schedule Trigger

```python
# Create workflow
workflow = await create_workflow({
    "name": "Daily Digest",
    "steps": [...]
})

# Add schedule trigger
await create_schedule_trigger({
    "workflow_id": workflow.id,
    "name": "Evening Digest",
    "cron_expression": "0 18 * * *",
    "timezone": "America/Chicago"
})
```

---

## Composio Integration (Future)

Composio provides **800+ pre-built integrations** with managed OAuth for AI agents. When ready to add external service integrations, Composio eliminates the need for custom API wrappers.

### What Composio Provides

| Feature | Description |
|---------|-------------|
| **800+ Integrations** | Slack, Gmail, GitHub, Notion, Linear, Salesforce, etc. |
| **Managed OAuth** | No token refresh handling - Composio manages auth |
| **Agno Native** | First-class support via `composio_agno` package |
| **Function Calling** | Tools designed for LLM function calling |
| **Triggers** | Built-in triggers for webhooks from external services |

### Installation

```bash
pip install composio-agno
composio login  # Authenticate with Composio
composio add slack  # Connect user's Slack account
composio add github  # Connect user's GitHub account
```

### Using Composio Tools in Workflow Steps

```python
from composio_agno import ComposioToolSet, App, Action
from agno.agent import Agent

# Initialize Composio
composio = ComposioToolSet()

# Option 1: Get all tools for an app
slack_tools = composio.get_tools(apps=[App.SLACK])

# Option 2: Get specific actions only
tools = composio.get_tools(actions=[
    Action.SLACK_SEND_MESSAGE,
    Action.GITHUB_CREATE_ISSUE,
    Action.NOTION_CREATE_PAGE,
])

# Create agent with Composio tools
notification_agent = Agent(
    id="notifier",
    name="Notification Agent",
    model=get_model("sonnet"),
    tools=tools,
    instructions="Send notifications to appropriate channels based on urgency."
)
```

### Workflow Step with Composio Agent

```json
{
  "type": "step",
  "name": "Notify Team",
  "agent_id": "notifier",
  "prompt_template": "Send a Slack message to #team-updates: {summary}",
  "output_key": "notification_result"
}
```

### Available App Categories

| Category | Apps |
|----------|------|
| **Communication** | Slack, Discord, Teams, Gmail, Twilio |
| **Project Management** | Linear, Jira, Asana, Notion, Monday |
| **Development** | GitHub, GitLab, Bitbucket, Vercel |
| **CRM** | Salesforce, HubSpot, Pipedrive |
| **Productivity** | Google Calendar, Todoist, Airtable |
| **Storage** | Google Drive, Dropbox, Box |
| **Analytics** | Mixpanel, Amplitude, Segment |

### Why Composio vs Custom Integrations

| Aspect | Custom | Composio |
|--------|--------|----------|
| OAuth handling | Manual token refresh | Managed |
| API wrappers | Build per service | Pre-built |
| Rate limiting | Manual | Built-in |
| Schema updates | Manual maintenance | Auto-updated |
| Time to integrate | Days per service | Minutes |

### Composio Triggers (External → Workflow)

Composio can also trigger workflows from external events:

```python
from composio_agno import ComposioToolSet, App

composio = ComposioToolSet()

# Get triggers from services
triggers = composio.get_triggers(apps=[App.GITHUB, App.SLACK])

# Example: GitHub PR opened → Start code review workflow
# Configure in Composio dashboard or programmatically
```

### Implementation Plan

1. **Install**: `pip install composio-agno`
2. **Configure**: Add Composio API key to environment
3. **Create integration agents**: Build agents with Composio tools
4. **Use in workflows**: Reference agents in workflow steps
5. **Add triggers**: Connect external webhooks to workflows

This replaces the need for:
- Custom OAuth token management
- Individual API wrappers per service
- Webhook configuration per service
- Rate limit handling
- Schema maintenance

---

## Initialization

Triggers are automatically initialized on app startup:

```python
# In agentos_setup.py lifespan
async def lifespan(app: FastAPI):
    # ... agent initialization ...

    # Initialize workflow triggers
    await initialize_workflow_triggers()

    yield

    # Cleanup
    await shutdown_workflow_triggers()
```

This loads:
1. Schedule triggers into APScheduler
2. Event triggers subscribe to NATS
3. HTTP triggers ready on-demand


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
