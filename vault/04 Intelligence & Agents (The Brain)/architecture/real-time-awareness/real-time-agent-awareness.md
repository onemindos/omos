# Real-Time Agent Awareness System

**Status**: Design Complete, Implementation Pending
**Created**: 2025-02-05
**Owner**: Zeus Legacy

## Overview

This document outlines the architecture for making OneMind OS agents **proactive** through real-time event subscriptions. Instead of agents only responding when explicitly called, they will continuously monitor the NATS event stream and take autonomous actions based on their configuration and permissions.

## Current State (Reactive)

```
User → API Request → Agent Run → Response
                ↓
         NATS Events (side effect)
```

Agents are **stateless workers** that:
- Only run when explicitly called via API
- Have no awareness of ongoing events
- Cannot take proactive actions
- Miss opportunities for automation

## Target State (Proactive)

```
External Event → NATS Stream → Agent Subscriptions → Autonomous Actions
                                       ↓
                                  Activity Feed
                                       ↓
                                  User Approval (if needed)
```

Agents become **continuous processes** that:
- Subscribe to relevant event patterns
- Monitor the event stream in real-time
- Autonomously decide when to act
- Create activity items for user visibility
- Request approval via HITL when needed

## Architecture

### 1. Event Subscriptions

Each agent can subscribe to NATS event patterns based on their role:

```python
# backend/agno/agents/subscriptions.py

AGENT_SUBSCRIPTIONS = {
    "legacy-mother-ai": [
        "events.>",           # All events for context
        "urgent.>",           # Urgent events for immediate action
    ],
    "task-manager": [
        "events.todoist.>",   # All Todoist events
        "events.calendar.>",  # Calendar for deadline awareness
    ],
    "home-controller": [
        "events.home.>",      # All home automation events
        "events.presence.>",  # Presence detection
    ],
    "security-agent": [
        "events.home.motion.*",
        "events.home.door.*",
        "biometric.alert.>",  # Critical biometric alerts
    ],
    "researcher": [
        "events.github.pr.opened",
        "events.github.issue.created",
    ],
}
```

### 2. Agent Runtime Manager

A new service manages long-running agent processes:

```python
# backend/agno/runtime/agent_runtime.py

class AgentRuntime:
    """
    Manages long-running agent processes that subscribe to events.

    Each agent runs in its own async task, monitoring its subscriptions
    and taking autonomous actions based on its configuration.
    """

    def __init__(self, agent_id: str, subscriptions: list[str]):
        self.agent_id = agent_id
        self.subscriptions = subscriptions
        self.event_bus = None  # NATS EventBus
        self.agent = None  # Agno Agent instance
        self._running = False

    async def start(self):
        """Start the agent runtime and subscribe to events."""
        # Load agent configuration
        self.agent = await load_agent(self.agent_id)

        # Connect to NATS
        self.event_bus = await create_event_bus()

        # Subscribe to agent's event patterns
        for pattern in self.subscriptions:
            await self.event_bus.subscribe(
                pattern=pattern,
                handler=self._handle_event
            )

        # Start listening
        self._running = True
        await self.event_bus.listen(queue_group=f"agent-{self.agent_id}")

    async def _handle_event(self, event: Event):
        """
        Handle incoming event and decide if action is needed.

        Flow:
        1. Evaluate if event is relevant (agent's decision)
        2. If relevant, create activity item or take action
        3. If action requires approval, pause for HITL
        4. Track decision in agent's learning system
        """
        try:
            # Run agent with event context
            decision = await self._evaluate_event(event)

            if decision.should_act:
                if decision.requires_approval:
                    # Create activity item for user approval
                    await self._create_activity_item(event, decision)
                else:
                    # Take autonomous action
                    await self._execute_action(event, decision)

            # Log decision for learning
            await self._log_decision(event, decision)

        except Exception as e:
            logger.error(f"Error handling event {event}: {e}")

    async def _evaluate_event(self, event: Event) -> Decision:
        """
        Use agent to evaluate if event requires action.

        Prompt: "Event: {event}. Should you take action? Why?"
        """
        # Run agent in evaluation mode (cheap, fast model)
        response = await self.agent.run(
            message=f"Event: {event.to_dict()}. Is this relevant to your role? Should you take action?",
            model="haiku",  # Fast, cheap for evaluation
            stream=False,
        )

        # Parse decision from response
        return self._parse_decision(response)

    async def _execute_action(self, event: Event, decision: Decision):
        """Execute autonomous action based on decision."""
        # Run agent with full context
        result = await self.agent.run(
            message=decision.action_prompt,
            model=self.agent.model,
            session_id=f"event-{event.sequence}",
        )

        # Broadcast result as new event
        await self.event_bus.publish(Event(
            type=f"action.{self.agent_id}.completed",
            source=self.agent_id,
            data={
                "trigger_event": event.to_dict(),
                "result": result,
            },
        ))
```

### 3. Agent Configuration

Agents declare their event subscriptions and permissions in presets:

```python
# backend/agno/agents/presets.py

TASK_MANAGER = {
    "agent_id": "task-manager",
    "name": "Task Manager",

    # Event subscriptions (NEW)
    "subscriptions": [
        "events.todoist.task.created",
        "events.todoist.task.completed",
        "events.todoist.task.due_soon",
        "events.calendar.meeting.soon",
    ],

    # Autonomous actions (NEW)
    "autonomous_actions": {
        "enabled": True,
        "requires_approval": ["task.create", "task.update"],  # These need HITL
        "auto_approve": ["task.reminder"],  # These can run autonomously
    },

    # Existing config...
}
```

### 4. Activity Feed Integration

When agents decide to act, they create activity items:

```python
async def _create_activity_item(self, event: Event, decision: Decision):
    """Create activity item for user review."""
    from backend.consciousness.activity.service import ActivityService

    activity = await ActivityService.create_item(
        type="approval" if decision.requires_approval else "notification",
        title=decision.title,
        description=decision.description,
        priority=event.priority,
        metadata={
            "trigger_event": event.to_dict(),
            "agent_id": self.agent_id,
            "proposed_action": decision.action_prompt,
        },
        assignee_agent_id=self.agent_id if not decision.requires_approval else None,
    )

    # If approval needed, pause agent action
    if decision.requires_approval:
        await activity.wait_for_approval()

    # Execute once approved
    await self._execute_action(event, decision)
```

### 5. Startup Integration

In `backend/agno/onemindos.py` lifespan:

```python
# Start Agent Runtimes (proactive agents)
agent_runtimes = []
try:
    from backend.agno.runtime.agent_runtime import start_agent_runtimes
    from backend.agno.agents.subscriptions import AGENT_SUBSCRIPTIONS

    agent_runtimes = await start_agent_runtimes(AGENT_SUBSCRIPTIONS)
    logger.info(f"Started {len(agent_runtimes)} proactive agents")
except Exception as e:
    logger.debug(f"Agent runtimes not available: {e}")
```

## Implementation Phases

### Phase 5.1: Core Infrastructure (Week 1)
- [ ] Create `backend/agno/runtime/` directory
- [ ] Implement `AgentRuntime` class
- [ ] Implement `start_agent_runtimes()` function
- [ ] Add subscriptions config to agent presets
- [ ] Test with single agent (task-manager)

### Phase 5.2: Decision Logic (Week 2)
- [ ] Implement event evaluation system
- [ ] Add decision logging to Learning Machines
- [ ] Create decision templates for common patterns
- [ ] Test autonomous decisions with approval flow

### Phase 5.3: Activity Feed Integration (Week 3)
- [ ] Wire agent actions to activity feed
- [ ] Implement approval waiting mechanism
- [ ] Add activity notifications to frontend
- [ ] Test full proactive workflow

### Phase 5.4: Multi-Agent Coordination (Week 4)
- [ ] Enable multiple agents per event pattern
- [ ] Implement agent-to-agent communication
- [ ] Add conflict resolution for competing actions
- [ ] Load testing with all agents active

## Benefits

### For Users
- **Proactive Automation**: Agents handle routine tasks without prompting
- **Intelligent Filtering**: Only surface what needs human attention
- **Context Awareness**: Agents see full event history, not just current request
- **Reduced Cognitive Load**: System handles monitoring and triage

### For Agents
- **Continuous Learning**: Every event is a learning opportunity
- **Better Context**: Access to full event stream vs. single request
- **Autonomous Growth**: Agents improve through observation and decision logging
- **Coordination**: Agents can collaborate on complex tasks

## Example Workflows

### 1. Smart Task Management

```
Event: Calendar meeting created "Project Review @ 2pm"
  ↓
Task Manager Agent:
  - Sees calendar event
  - Checks if related Todoist project exists
  - Creates pre-meeting checklist task
  - Schedules reminder 1 hour before
  ↓
Activity Item: "Created pre-meeting tasks for Project Review"
  ↓
User: [Approve] or [Modify] or [Reject]
```

### 2. Home Automation

```
Event: Motion detected in bedroom @ 6:45am
  ↓
Home Controller Agent:
  - Sees motion event
  - Checks time (morning routine)
  - Checks biometric data (awake, not sleeping)
  - Decides to start morning routine
  ↓
Autonomous Action: (no approval needed)
  - Turn on bedroom lights (25% brightness)
  - Start coffee maker
  - Open blinds
```

### 3. Security Alert

```
Event: Door opened @ 2:30am
  ↓
Security Agent:
  - Sees door event
  - Checks presence data (no one home)
  - Evaluates threat level (HIGH)
  ↓
Urgent Activity Item: "Door opened while away - camera snapshot attached"
  ↓
User: [Call Police] or [It's me] or [False alarm]
```

## Performance Considerations

### Event Volume
- Expected: 1000-5000 events/day
- Peak: 10,000 events/day
- Per agent: 100-500 relevant events/day

### Latency
- Event → Agent: <50ms (NATS push)
- Evaluation: 200-500ms (Haiku model)
- Decision logging: <10ms (async)
- Total: <1 second from event to decision

### Resource Usage
- 9 agents × 1 runtime = 9 background tasks
- Memory: ~50MB per agent runtime
- CPU: Idle (event-driven), spikes during evaluation

### Scaling
- Agents use NATS queue groups (horizontal scaling)
- Multiple instances share event load automatically
- No coordination needed between instances

## Security & Safety

### Permission System
- Agents declare required permissions in config
- User can enable/disable autonomous actions per agent
- All actions logged to audit trail
- HITL approval for sensitive operations

### Rate Limiting
- Max 10 autonomous actions per agent per hour
- Backoff after failures
- Alert user if agent behavior is anomalous

### Monitoring
- Track agent decision accuracy
- Alert on excessive activity creation
- Dashboard for agent health and activity

## Future Enhancements

### Multi-Agent Teams (Phase 6)
- Agents can spawn sub-agents for delegation
- Team coordination via NATS messaging
- Shared context and memory

### Predictive Actions (Phase 7)
- Agents analyze patterns to predict needs
- Proactive suggestions before events occur
- "You usually X when Y happens, should I?"

### Learning from Feedback (Phase 8)
- Track approval/rejection rates
- Adjust decision thresholds automatically
- Personalize agent behavior per user

## Conclusion

Real-time agent awareness transforms OneMind OS from a reactive chatbot system into a **proactive AI operating system**. Agents become living processes that continuously monitor your life, make intelligent decisions, and act autonomously within defined permissions.

This is the foundation for true AI assistance - where the system anticipates needs, handles routine tasks, and only surfaces what truly requires human attention.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
