# OneMind Execution Engine

> Orchestrates inbox tasks to agents, teams, and workflows with HITL support.

## Overview

The Execution Engine is the "glue layer" that transforms the inbox from a passive notification system into an **active work orchestration hub**. When a task arrives in the inbox, the engine:

1. **Classifies** the task (category, urgency, complexity)
2. **Selects** the appropriate execution target (agent, team, or workflow)
3. **Injects** context from the inbox item
4. **Executes** and tracks progress
5. **Updates** the inbox on completion

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      INBOX ITEM ARRIVES                          │
│  Type: task | Title: "Review PR #456" | Priority: high          │
│  Context: {pr_url, files_changed, author}                       │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      EXECUTION ENGINE                            │
├─────────────────────────────────────────────────────────────────┤
│  1. Classify task → "code_review"                               │
│  2. Match execution config (from DB or defaults)                │
│  3. Resolve target: agent/team/workflow                         │
│  4. Build context from inbox item                               │
│  5. Execute and track via NATS events                           │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    EXECUTION TARGETS                             │
├─────────────┬─────────────┬─────────────────────────────────────┤
│   AGENT     │    TEAM     │           WORKFLOW                  │
│             │             │                                     │
│ Single AI   │ Multi-agent │ Multi-step with                    │
│ execution   │ coordination│ approval gates                     │
│             │             │                                     │
│ Legacy      │ Code Review │ Deploy Workflow                     │
│ Research    │ Dev Team    │ Morning Briefing                    │
│ Briefing    │             │ Task Triage                         │
└─────────────┴─────────────┴─────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    PROGRESS TRACKING                             │
├─────────────────────────────────────────────────────────────────┤
│  NATS Events: engine.progress.*                                 │
│  WebSocket: Real-time UI updates                                │
│  Database: inbox_item status updates                            │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    HITL (IF NEEDED)                              │
├─────────────────────────────────────────────────────────────────┤
│  1. Workflow pauses at approval step                            │
│  2. Approval request created                                    │
│  3. You approve/deny from Flutter                               │
│  4. Workflow continues                                          │
└─────────────────────────────────────────────────────────────────┘
```

## API Endpoints

### Inbox Dispatch Endpoints

```bash
# Dispatch an inbox item (auto-classifies and routes)
POST /api/inbox/{item_id}/dispatch
{
  "execution_type": "agent",     # Optional: agent, team, workflow
  "execution_target": "legacy-mother-ai",  # Optional: specific target
  "team_mode": "coordinate",     # Optional: for teams
  "require_approval": false      # Optional: override HITL
}

# Classify without dispatching
GET /api/inbox/{item_id}/classify
# Returns: category, suggested_target, requires_approval
```

### Dispatcher API

```bash
# Full dispatch with all options
POST /api/engine/dispatch
{
  "inbox_item_id": "uuid-123",
  "execution_type": "workflow",
  "execution_target": "deploy-workflow",
  "additional_context": {"version": "2.1.0"},
  "require_approval": true,
  "actor": "zeus"
}

# Get execution status
GET /api/engine/status/{run_id}

# List active executions
GET /api/engine/active

# List execution configs
GET /api/engine/configs

# Create custom config
POST /api/engine/configs
{
  "name": "Code Review Team",
  "execution_type": "team",
  "execution_target": "code-review-team",
  "task_categories": ["code_review"],
  "keywords": ["review", "PR", "pull request"],
  "team_mode": "coordinate"
}

# Health check
GET /api/engine/health
```

## Task Classification

The engine uses **keyword matching** and **heuristics** to classify tasks:

| Category | Keywords | Default Target |
|----------|----------|----------------|
| `code_review` | review, PR, pull request | code-review-team |
| `research` | research, investigate, find | research-agent |
| `deployment` | deploy, release, ship | deploy-workflow (HITL) |
| `writing` | write, draft, compose | Legacy agent |
| `planning` | morning, triage, evening | Workflows |
| `general` | (default) | Legacy agent |

### Classification Response

```json
{
  "inbox_item_id": "uuid-123",
  "category": "code_review",
  "confidence": 0.8,
  "keywords_matched": ["review", "PR"],
  "suggested_execution_type": "team",
  "suggested_target": "code-review-team",
  "requires_approval": false,
  "reason": "Matched 2 keywords"
}
```

## Execution Configs

Execution configs define how task categories map to execution targets:

```sql
-- Database table
CREATE TABLE execution_configs (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    execution_type VARCHAR(20) NOT NULL,  -- agent, team, workflow
    execution_target VARCHAR(255) NOT NULL,
    task_categories TEXT[],
    keywords TEXT[],
    requires_approval BOOLEAN,
    approval_steps TEXT[],
    team_mode VARCHAR(20),
    timeout_seconds INT DEFAULT 600,
    is_active BOOLEAN DEFAULT TRUE
);
```

### Default Configs

| Name | Type | Target | Categories | HITL |
|------|------|--------|------------|------|
| Legacy Agent | agent | legacy-mother-ai | general, communication | No |
| Morning Briefing | workflow | morning-briefing | planning | No |
| Task Triage | workflow | task-triage | planning | No |
| Focus Mode | workflow | focus-mode | planning | No |
| Evening Review | workflow | evening-review | planning | No |

## Context Injection

The engine injects context from inbox items into agent/team/workflow execution:

```python
context = {
    "inbox_item_id": "uuid-123",
    "title": "Review PR #456",
    "priority": "high",
    "source": "github",
    "type": "task",
    "body": "Security fixes in auth module",
    "pillar": "IT",  # Four Pillars alignment
    "due_at": "2026-01-26T00:00:00Z",
    "action_url": "https://github.com/...",
    "metadata": {
        "pr_url": "...",
        "files_changed": ["auth.py", "api.py"]
    }
}
```

## Progress Tracking

### NATS Events

```
engine.progress.{inbox_item_id}
{
  "inbox_item_id": "uuid-123",
  "run_id": "run-456",
  "status": "running",  # pending, running, awaiting_approval, completed, failed
  "step": "deploy",
  "step_index": 3,
  "total_steps": 4,
  "message": "Deploying to production...",
  "percent_complete": 75.0,
  "timestamp": "2026-01-25T12:00:00Z"
}
```

### WebSocket Updates

The Progress Tracker broadcasts to WebSocket clients:

```json
{
  "type": "dispatcher.progress",
  "data": {
    "inbox_item_id": "uuid-123",
    "run_id": "run-456",
    "status": "completed",
    "message": "Task completed"
  }
}
```

## HITL Integration

For workflows with `requires_approval=True`:

1. **Workflow pauses** at approval step
2. **Approval created** in `approvals` table
3. **Inbox item updated** with `approval_id` link
4. **Notification sent** via notification router
5. **You approve/deny** from Flutter or API
6. **Workflow continues** with your decision

### Approval Flow

```python
# Workflow definition
deploy_workflow = Workflow(
    steps=[
        WorkflowStep(name="tests", agent=test_agent),
        WorkflowStep(name="build", agent=build_agent),
        WorkflowStep(
            name="deploy",
            agent=deploy_agent,
            requires_approval=True,  # HITL gate
        ),
        WorkflowStep(name="verify", agent=verify_agent),
    ]
)
```

## File Structure

```
backend/consciousness/engine/
├── __init__.py           # Exports
├── models.py             # Pydantic models (~200 lines)
├── engine.py             # Main execution engine logic (~600 lines)
├── progress.py           # Progress tracker (~250 lines)
└── api.py                # FastAPI endpoints (~200 lines)
```

## Usage Examples

### Dispatch from UI

```typescript
// Flutter/React
const dispatch = async (itemId: string) => {
  const response = await fetch(`/api/inbox/${itemId}/dispatch`, {
    method: 'POST',
  });
  const result = await response.json();
  console.log(`Dispatched: ${result.run_id}`);
};
```

### Dispatch with Override

```bash
# Force workflow execution with HITL
curl -X POST "http://localhost:8081/api/inbox/uuid-123/dispatch" \
  -H "Content-Type: application/json" \
  -d '{
    "execution_type": "workflow",
    "execution_target": "deploy-workflow",
    "require_approval": true
  }'
```

### Assign to Team

```bash
# Assign to code review team
curl -X POST "http://localhost:8081/api/inbox/uuid-123/dispatch" \
  -H "Content-Type: application/json" \
  -d '{
    "execution_type": "team",
    "execution_target": "code-review-team",
    "team_mode": "coordinate"
  }'
```

## Integration with Existing Systems

### Inbox Service

The engine integrates with the inbox service:
- Updates `inbox_items.status` on dispatch
- Records history in `inbox_item_history`
- Broadcasts updates via WebSocket

### Approvals Handler

For HITL workflows:
- Uses `core/approvals/handler.py` for paused run management
- Creates approval records in `approvals` table
- Links approval to inbox item via `approval_id`

### Notification Router

Completion/failure triggers notifications:
- Uses `core/notifications/router.py`
- Channels based on priority and awareness mode

## Database Tables

### execution_configs

Stores routing rules for task categories.

### inbox_items (updated)

- `status`: Updated to "reading" on dispatch, "handled" on complete
- `task_result`: Stores execution output
- `approval_id`: Links to pending approval (if HITL)

### inbox_item_history (updated)

Records dispatch events:
```json
{
  "action": "dispatched",
  "actor": "zeus",
  "new_value": {
    "status": "running",
    "run_id": "run-456",
    "execution_type": "workflow"
  }
}
```

## Troubleshooting

### Dispatcher Not Available

```bash
# Check health
curl http://localhost:8081/api/engine/health

# Returns:
{
  "status": "healthy",
  "dispatcher_available": true,
  "db_available": true,
  "nats_available": true,
  "agent_os_available": true
}
```

### Execution Not Starting

1. Check if agent/team/workflow exists
2. Check execution configs
3. Check NATS connection
4. Check agent_os initialization

### Progress Not Updating

1. Check NATS subscription
2. Check WebSocket connection
3. Check progress tracker is running

## Flutter UI

The Flutter app provides a **Dispatch Dialog** that fetches live data from the registry API.

### Dispatch Button Locations

1. **Task cards** - Red "Dispatch" button on task-type inbox items
2. **Long-press menu** - "Dispatch" option in action sheet for any item

### Dispatch Dialog

When triggered, the dialog:

1. **Loads registry** - Fetches `/api/registry` for available targets
2. **Shows type selector** - Agent/Team/Workflow chips with counts
3. **Shows target dropdown** - Lists registered items with descriptions
4. **Team mode** - Coordinate/Collaborate/Route (only for teams)
5. **HITL toggle** - "Require Approval" switch

```
┌─────────────────────────────────────────┐
│ 🚀 Dispatch Task                        │
├─────────────────────────────────────────┤
│ ┌─────────────────────────────────────┐ │
│ │ Review PR #456                      │ │
│ │ Security fixes in auth module       │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ Dispatch To                             │
│ ┌─────────┐ ┌─────────┐ ┌───────────┐  │
│ │ Agent   │ │ Team    │ │ Workflow  │  │
│ │ 🤖 (2)  │ │ 👥 (9)  │ │ 🌳 (5)    │  │
│ └─────────┘ └─────────┘ └───────────┘  │
│                                         │
│ Target                                  │
│ ┌─────────────────────────────────────┐ │
│ │ Dev Team                    4 members│ │
│ │ Development team for software tasks │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ Team Mode                               │
│ ┌─────────────────────────────────────┐ │
│ │ Coordinate                          │ │
│ │ Leader delegates to best member     │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ ☐ Require Approval                      │
│   Pause for human approval before...    │
│                                         │
│              [Cancel]  [🚀 Dispatch]    │
└─────────────────────────────────────────┘
```

### Files

| File | Purpose |
|------|---------|
| `ui/lib/features/inbox/screens/inbox_screen.dart` | `_DispatchDialog` widget (~400 lines) |
| `ui/lib/features/inbox/providers/unified_inbox_provider.dart` | `getDispatchTargets()`, `DispatchTargets`, `DispatchTarget` |
| `ui/lib/features/inbox/widgets/inbox_item_card.dart` | Dispatch button in `_ActionButtons` |

### Data Flow

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Dispatch       │     │  Provider       │     │  Backend        │
│  Dialog         │     │                 │     │                 │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ initState()     │────▶│ getDispatch     │────▶│ GET /api/       │
│                 │     │ Targets()       │     │ registry        │
│                 │◀────│                 │◀────│                 │
│ _targets =      │     │ DispatchTargets │     │ {agents, teams, │
│   loaded        │     │   .fromJson()   │     │  workflows}     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ User clicks     │────▶│ dispatch()      │────▶│ POST /api/inbox │
│ "Dispatch"      │     │                 │     │ /{id}/dispatch  │
│                 │◀────│                 │◀────│                 │
│ Close + snackbar│     │ result          │     │ {run_id, status}│
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

### Key Classes

```dart
/// Registry response containing all dispatch targets
class DispatchTargets {
  final List<DispatchTarget> agents;
  final List<DispatchTarget> teams;
  final List<DispatchTarget> workflows;

  List<DispatchTarget> getTargets(String type); // 'agent'|'team'|'workflow'
}

/// A single dispatch target (agent, team, or workflow)
class DispatchTarget {
  final String id;          // 'dev-team', 'legacy-mother-ai'
  final String name;        // 'Dev Team', 'Legacy'
  final String description; // 'Development team for...'
  final List<String> members; // ['Coder', 'Reviewer'] (teams only)
  final List<String> tags;    // ['development', 'it']
}
```

### States

| State | SO |
|-------|-----|
| Loading | Spinner + "Loading dispatch targets..." |
| Error | Error icon + message + "Retry" button |
| Empty | "No agents registered" info box |
| Ready | Type chips + target dropdown + dispatch button |
| Dispatching | Spinner in button + disabled |

### Usage

```dart
// Trigger from inbox screen
void _showDispatchDialog(InboxItem item) {
  showDialog(
    context: context,
    builder: (context) => _DispatchDialog(
      item: item,
      onDispatch: (type, target, teamMode, requireApproval) async {
        final result = await notifier.dispatch(
          item.id,
          executionType: type,
          executionTarget: target,
          teamMode: teamMode,
          requireApproval: requireApproval,
        );
        // Show success/error snackbar
      },
    ),
  );
}
```

## Related Documentation

- [Notification System](NOTIFICATION-SYSTEM.md) - Multi-channel notifications
- [HITL Approvals](../operations/HITL-APPROVALS.md) - Human-in-the-loop system
- [Unified Inbox](../UNIFIED-INBOX-SYSTEM.md) - Inbox architecture
- [Agno Workflows](../legacy/AGNO_FEATURES_ANALYSIS.md) - Workflow capabilities
