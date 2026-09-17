---
id: so-02-20260528-nats-session-protocol-deep
title: "NATS Session Protocol — Deep Message Structure"
type: reference
status: active
created: 2026-05-28
modified: 2026-05-28
owner: zeus
pillar: SO
domain: "02 Infrastructure (The Body)"
tags: [type/reference, zone/so, domain/so, status/active, infra/nats, agents/sessions, protocol/deep]
depth: 0
---

# NATS Session Protocol — Deep Message Structure

> What's actually inside a session event on the NATS fabric. The subject is the pointer. The stream is the depth.

---

## Subject Anatomy

```
sessions.greenhouse.maintenance.msg.47
│        │           │           │   │
│        │           │           │   └── message sequence in session
│        │           │           └────── message type (msg, meta, close)
│        │           └────────────────── session name (persistence key)
│        └────────────────────────────── agent name
└─────────────────────────────────────── stream namespace
```

---

## Example: User Message (Zeus says "order replacement part")

```
Subject: sessions.greenhouse.maintenance.msg.47

Headers:
  Nats-Msg-Id: sess-greenhouse-maint-msg-47
  X-Session-Id: greenhouse-maintenance-2026
  X-Agent: greenhouse
  X-Owner: zeus
  X-Turn: 47
  X-Role: user
  X-Timestamp: 2026-05-15T14:35:00Z
  X-Source: voice
  X-Correlation-Id: evt-1203847

Payload:
{
  "role": "user",
  "content": "order replacement part",
  "source": {
    "platform": "voice",
    "device": "macbook-pro",
    "location": "office"
  },
  "context_at_time": {
    "triggered_by": "events.alert.greenhouse.vent_failure",
    "greenhouse_state": {
      "zone1_temp": 79.2,
      "east_vent": "jammed",
      "west_vent": "open"
    },
    "preceding_events": [
      {
        "seq": 1203845,
        "subject": "telemetry.greenhouse.vent_east.current",
        "value": "2.3A (spike)"
      },
      {
        "seq": 1203846,
        "subject": "commands.greenhouse.vent_west.open",
        "result": "success"
      },
      {
        "seq": 1203847,
        "subject": "events.alert.greenhouse.vent_failure",
        "severity": "warning"
      }
    ]
  }
}
```

---

## Example: Agent Response (greenhouse agent acts)

```
Subject: sessions.greenhouse.maintenance.msg.48

Headers:
  Nats-Msg-Id: sess-greenhouse-maint-msg-48
  X-Session-Id: greenhouse-maintenance-2026
  X-Agent: greenhouse
  X-Owner: zeus
  X-Turn: 48
  X-Role: assistant
  X-Model: hermes-3b
  X-Tokens-In: 2340
  X-Tokens-Out: 187
  X-Latency-Ms: 890
  X-Tools-Called: saas.execute.todoist.create_task,nats_kv_put
  X-Timestamp: 2026-05-15T14:35:01.890Z

Payload:
{
  "role": "assistant",
  "content": "Created task in Todoist: 'Replace east vent actuator — greenhouse zone 1'. Priority 2, due next Monday. Also updated maintenance log in KV.",
  "tool_calls": [
    {
      "tool": "saas.execute",
      "args": {
        "app": "todoist",
        "action": "create_task",
        "content": "Replace east vent actuator — greenhouse zone 1",
        "priority": 2,
        "due_string": "next Monday"
      },
      "result": {
        "task_id": "8847291",
        "status": "created"
      }
    },
    {
      "tool": "nats_kv_put",
      "args": {
        "bucket": "maintenance",
        "key": "greenhouse.vent_east",
        "value": {
          "status": "failing",
          "symptom": "current_spike_no_position_change",
          "ordered": false,
          "due": "2026-05-19"
        }
      },
      "result": {
        "revision": 3
      }
    }
  ],
  "decision_trace": {
    "reasoning": "User requested part order. East vent actuator showed current spike at 2.3A with no position change = mechanical failure. Creating maintenance task with appropriate urgency.",
    "confidence": 0.98,
    "experience_refs": [
      {
        "stream": "EXPERIENCE",
        "seq": 892341,
        "note": "Similar actuator failure on west vent 3 months ago — same symptom pattern"
      }
    ]
  }
}
```

---

## What This Enables

When Legacy (or any agent) later encounters `sessions.greenhouse.maintenance`, it can:

1. **Pull the full session:** `nats_stream_read("SESSIONS", {filter: "sessions.greenhouse.maintenance.>"})`
2. **See every message** in the conversation (full history)
3. **See what triggered it** — the `X-Correlation-Id` links to the vent failure event
4. **See what the agent did** — tool_calls show exact actions taken
5. **See the reasoning** — decision_trace explains WHY
6. **See experience references** — links to past similar events the agent learned from
7. **See the world state at that moment** — context_at_time captures what was happening

---

## Message Types in a Session Stream

| Subject suffix | Purpose |
|---------------|---------|
| `*.msg.<N>` | Conversation messages (user + assistant turns) |
| `*.meta` | Session metadata (created, status, participants, total tokens) |
| `*.tool.<N>` | Tool execution details (args, result, duration) |
| `*.event.<N>` | Session lifecycle events (created, resumed, paused, closed) |
| `*.summary` | Agent-generated session summary (for fleet dashboard) |

---

## How Legacy Queries Across All Sessions

```python
# Find all sessions where maintenance was discussed
nats_stream_read("SESSIONS", {
    filter: "sessions.*.maintenance.>",
    start_time: "2026-01-01T00:00:00Z"
})

# Find all sessions for a specific agent
nats_stream_read("SESSIONS", {
    filter: "sessions.greenhouse.>",
    headers_match: {"X-Role": "assistant", "X-Tools-Called": "*todoist*"}
})

# Find sessions triggered by a specific event type
nats_stream_read("SESSIONS", {
    filter: "sessions.>",
    headers_match: {"X-Correlation-Id": "evt-*vent_failure*"}
})
```

---

## The Depth Stack

```
SURFACE (what you see in the timeline):
  "Zeus said 'order replacement part' (sessions.greenhouse.maintenance)"

ONE LEVEL DOWN (pull the session):
  Full conversation: 48 messages, 12 tool calls, 3 decisions

TWO LEVELS DOWN (follow correlation IDs):
  The triggering event chain: sensor spike → failure detection → alert → conversation

THREE LEVELS DOWN (follow experience refs):
  Past similar events the agent learned from, months ago

FOUR LEVELS DOWN (query the telemetry stream):
  Raw sensor data at the exact moment — current readings, position data, full state

Every layer is addressable. Every layer is stored in JetStream or Parquet.
The AI has DEPTH OF UNDERSTANDING, not just surface recall.
```
