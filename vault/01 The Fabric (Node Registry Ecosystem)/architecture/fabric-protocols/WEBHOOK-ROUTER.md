---
id: webhook-router-20260430
title: "Webhook Router"
type: protocol
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/protocol, domain/so, system/todoist, system/webhooks, system/agents]
ai_summary: "Real-time event router that receives webhooks from Todoist and external systems, classifies them, and dispatches to the correct agent immediately."
---

# Webhook Router

> Instant response to the world. When something happens, the right agent knows immediately.

---

## Why Real-Time Matters

The heartbeat guarantees 15-minute coverage. But some events need **immediate** response:
- Zeus adds a P1 task → Agent should start within seconds, not wait 15 min
- A deploy fails → Alert agent immediately
- A client emails → Capture and classify now
- A cron completes → Trigger dependent tasks instantly
- Zeus sends a voice command → Route to agent NOW

The Webhook Router is the **real-time nervous system** that fills the gaps between heartbeats.

---

## Architecture

```text
EXTERNAL EVENTS                    WEBHOOK ROUTER                      AGENTS
──────────────                    ──────────────                      ──────

Todoist webhook ──────────┐
                          │
N8N automation ───────────┤       ┌─────────────────────┐
                          ├──────►│  INGRESS ENDPOINT   │
GitHub webhook ───────────┤       │  (VPS: port 9090)   │
                          │       └──────────┬──────────┘
Email (via Zapier) ───────┤                  │
                          │                  ▼
Voice command ────────────┤       ┌─────────────────────┐
                          │       │     CLASSIFIER      │
Home Assistant ───────────┘       │                     │
                                  │  • Parse event type │
                                  │  • Identify domain  │
                                  │  • Check priority   │
                                  │  • Load context     │
                                  └──────────┬──────────┘
                                             │
                                             ▼
                                  ┌─────────────────────┐
                                  │    DISPATCHER        │──────► Legacy
                                  │                     │──────► Spartan
                                  │  • Route to agent   │──────► Oracle
                                  │  • Attach context   │──────► Eden
                                  │  • Log event        │──────► Empire
                                  └─────────────────────┘
```

---

## Supported Event Sources

### Todoist Webhooks

Todoist sends webhooks for all task lifecycle events:

| Event Type | Trigger | Action |
|------------|---------|--------|
| `item:added` | New task created | Classify → route to domain agent |
| `item:updated` | Task modified (priority, due, labels) | Re-evaluate scheduling |
| `item:completed` | Task marked done | Log completion, trigger dependents |
| `item:uncompleted` | Task re-opened | Re-add to agent queue |
| `item:deleted` | Task removed | Clean up any scheduled crons |
| `note:added` | Comment added to task | Check if it's an instruction from Zeus |

**Todoist Webhook Setup:**
```text
URL: https://your-vps.com:9090/webhooks/todoist
Events: All item events + note events
Auth: X-Todoist-Hmac-SHA256 header validation
```

### System Events (Internal)

| Source | Events | Example |
|--------|--------|---------|
| **N8N** | Workflow completion, errors | "Weekly report generated" |
| **GitHub** | Push, PR, issue created | "New commit to onemind-codex" |
| **Cron system** | Cron completed/failed | "Deploy cron finished" |
| **Home Assistant** | Device events, routines | "Zeus left home" / "Bedtime routine" |
| **Email** | New email matching rules | "Client invoice received" |
| **Voice** | Transcribed commands | "Legacy, schedule a meeting" |

---

## Event Classification Logic

```text
EVENT ARRIVES
│
├─► Is it from Todoist?
│   ├─► YES → Check task labels/project for domain routing
│   │         • @health/@fitness → Spartan
│   │         • @mind/@money → Oracle
│   │         • @family/@home → Eden
│   │         • @business/@career → Empire
│   │         • No label / @system → Legacy (meta-agent)
│   │
│   └─► Check priority
│       • P1 → IMMEDIATE dispatch (interrupt current work)
│       • P2 → Dispatch within 60 seconds
│       • P3 → Queue for next heartbeat
│       • P4 → Queue (no urgency)
│
├─► Is it a system event?
│   ├─► Error/failure → Escalate to Legacy (meta-agent)
│   ├─► Completion → Check for dependent tasks, trigger next
│   └─► Informational → Log only, no action needed
│
├─► Is it from Zeus directly (voice/message)?
│   └─► IMMEDIATE dispatch to addressed agent
│       • "Legacy, do X" → Legacy
│       • "Spartan, log workout" → Spartan
│       • Unaddressed → Legacy routes it
│
└─► Unknown source?
    └─► Log as anomaly, alert Legacy for review
```

---

## Dispatch Protocol

When the router decides to dispatch to an agent:

```yaml
dispatch_payload:
  event_id: "evt-20260430-001"
  timestamp: "2026-04-30T14:32:05Z"
  source: "todoist"
  event_type: "item:added"
  priority: "p1"
  target_agent: "empire"
  context:
    task_id: "todoist-8234567890"
    task_title: "Send proposal to Acme Corp"
    task_due: "2026-04-30T17:00:00Z"
    project: "Business"
    labels: ["@business", "@urgent"]
    description: "Final proposal with pricing"
  codex_context:
    domain: "GE (75-99)"
    relevant_files:
      - "75-99 GE/80 Ventures/acme-proposal.md"
    active_goals:
      - "Close 3 new clients by end of April"
  instructions: "Execute immediately — P1 with today deadline"
```

---

## Zeus Direct Commands

Zeus can send real-time commands to agents through multiple interfaces:

| Interface | How It Works | Example |
|-----------|--------------|---------|
| **Todoist task** | Add P1 task with agent label | "Deploy hotfix @legacy @p1" |
| **Todoist comment** | Comment on existing task | "@legacy please prioritize this" |
| **Voice** | Speech-to-text → webhook | "Hey Legacy, check the server" |
| **Chat** | Direct message interface | VS Code Copilot / OpenClaw |
| **Signal** | Quick text command | "/legacy deploy staging" |

### Command Syntax

```text
@{agent} {action} {details}

Examples:
  @legacy deploy staging
  @spartan log workout: 30 min run
  @oracle research Bitcoin ETF options
  @eden order groceries from last week's list
  @empire send invoice to client-xyz
```

---

## Security & Validation

| Check | Implementation |
|-------|---------------|
| **Webhook auth** | HMAC-SHA256 signature validation for Todoist |
| **Source whitelist** | Only accept from known IPs/services |
| **Rate limiting** | Max 100 events/minute before throttling |
| **Payload validation** | Schema check before classification |
| **Replay prevention** | Event ID deduplication (24h window) |
| **Scope enforcement** | Agents can only act within their domain |

---

## Failure Handling

| Scenario | Response |
|----------|----------|
| Webhook endpoint down | Events queue in source (Todoist retries 3x) |
| Classification fails | Route to Legacy as fallback (meta-agent handles all) |
| Target agent offline | Queue event, deliver at next heartbeat |
| Duplicate event | Deduplicate by event_id, ignore repeats |
| Zeus unreachable for escalation | Hold escalation queue, batch-deliver when available |

---

## Implementation

| Component | Technology | Location |
|-----------|-----------|----------|
| Ingress endpoint | Express.js / Fastify on VPS | Port 9090, HTTPS |
| Classifier | Rule engine (YAML rules) + LLM fallback | `_system/webhook-rules.yaml` |
| Dispatcher | NATS pub/sub OR direct HTTP to agent | Agent-specific endpoints |
| Event log | Append-only YAML in Codex | `_system/events/` |
| Monitoring | N8N dashboard + Todoist alert task | P1 task on 3 consecutive failures |


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
