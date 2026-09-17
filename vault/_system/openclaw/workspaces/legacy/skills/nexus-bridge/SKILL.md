# SKILL.md — nexus-bridge

## What This Skill Does

Connects Legacy to the OneMind OS NATS nervous system. Use this to:
- Dispatch tasks to Agno digital agents (researcher, coder, analyst, etc.)
- Send commands to OpenMind OM1 physical robots
- Query system health across all runtimes
- Publish system events

## When To Use It

Use `nexus_dispatch` when you need a specialized agent to handle a task:
- Research tasks → dispatch to "researcher"
- Code tasks → dispatch to "coder"
- Analysis → dispatch to "analyst"
- Browser tasks → dispatch to "web_automator"
- Physical actions → dispatch to OM1 via robot commands

Use `nexus_query` to check system state:
- Agent health
- Robot status
- Location state
- System metrics

## NATS Subject Conventions

```
digital.agent.{agent_id}.task      → dispatch task to Agno agent
digital.agent.{agent_id}.result    → agent completed a task
physical.robot.{robot_id}.command  → send command to robot
physical.robot.{robot_id}.status   → robot status
physical.location.{id}.state       → location state query
nexus.system.health                → system health ping
nexus.heartbeat.>                  → all runtime heartbeats
comms.text.agent.legacy            → Legacy's outbound comms
comms.text.human.zeus              → Zeus's inbound messages
```

## Priority Levels

- `critical`: Immediate action required (system failure, emergency)
- `high`: Important, act within minutes
- `normal`: Standard task queue
- `low`: Background work, no rush

## Examples

Dispatch a research task:
```
nexus_dispatch({
  agent: "researcher",
  task: "Research latest OM1 firmware updates and summarize key changes",
  priority: "normal",
  context: { domain: "UI", project: "OneMind OS" }
})
```

Send robot command:
```
nexus_dispatch({
  target: "robot",
  robot_id: "om1-home-01",
  command: "navigate",
  params: { destination: "front_door" }
})
```

Check system health:
```
nexus_query({ subject: "nexus.system.health" })
```
