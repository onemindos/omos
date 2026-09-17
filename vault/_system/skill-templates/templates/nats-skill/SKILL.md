---
name: nats-skill
description: >
  [REPLACE THIS] Publishes and subscribes to OneMind OS NATS subjects. Use when
  you need to dispatch tasks to Agno agents, send commands to OM1 robots, query
  system health, or broadcast events across the OneMind runtime network.
compatibility: Requires a running NATS server. NATS_URL env var must be set.
metadata:
  author: zeus
  version: "1.0.0"
---

# NATS Skill

## When to Use

Use this skill to communicate across runtimes via NATS JetStream:

- Dispatch tasks to Agno digital agents
- Send commands to OM1 physical robots
- Query system health and status
- Publish events to the OneMind OS event bus

## Environment Variables Required

| Variable   | Description                                     |
| ---------- | ----------------------------------------------- |
| `NATS_URL` | NATS server URL (e.g., `nats://localhost:4222`) |

## Available Tools

| Tool           | Description                                          |
| -------------- | ---------------------------------------------------- |
| `nats_publish` | Publish a message to a NATS subject                  |
| `nats_request` | Publish and wait for a reply (request/reply pattern) |
| `nats_health`  | Ping the NATS connection and return status           |

## NATS Subject Conventions

```
digital.agent.{agent_id}.task          → Dispatch task to Agno agent
digital.agent.{agent_id}.result        → Agent completed result
physical.robot.{robot_id}.command      → Send command to OM1 robot
physical.robot.{robot_id}.status       → Robot status query
nexus.system.health                    → System-wide health ping
nexus.heartbeat.>                      → All runtime heartbeats
comms.text.agent.legacy                → Legacy's outbound comms
comms.text.human.zeus                  → Zeus's inbound messages
```

## Examples

### Dispatch a task to an Agno agent

```
nats_publish({
  subject: "digital.agent.researcher.task",
  data: {
    task: "Research the latest NATS JetStream features",
    priority: "normal",
    source: "legacy"
  }
})
```

### Send a robot command

```
nats_publish({
  subject: "physical.robot.main-bot.command",
  data: {
    command: "move_to",
    params: { location: "dock", speed: 0.5 }
  }
})
```

### Request-reply — check system health

```
nats_request({
  subject: "nexus.system.health",
  data: { source: "legacy" },
  timeout_ms: 3000
})
```

### Check NATS connection status

```
nats_health()
```

## Priority Levels (for agent dispatches)

| Priority   | When to Use                                           |
| ---------- | ----------------------------------------------------- |
| `critical` | Immediate action required — system failure, emergency |
| `high`     | Important — act within minutes                        |
| `normal`   | Standard task queue                                   |
| `low`      | Background work, no rush                              |

## Notes

- NATS connection is maintained across tool calls (singleton)
- If NATS is unreachable, tools will throw with connection error
- `nats_request` will timeout after `timeout_ms` (default 5000ms)
- Subjects are case-sensitive
