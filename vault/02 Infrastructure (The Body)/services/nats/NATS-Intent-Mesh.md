---
title: NATS Intent Mesh — Personal Command Language
category: NATS / Messaging
domain: 00-24 SO / 02 Infrastructure
status: active
created: 2026-05-24
tags: [nats, intent, commands, automation, onemind-os, agents, ha, parallel]
---

# NATS Intent Mesh — Personal Command Language

> **Core concept:** Zeus's words + context = deterministic parallel action across every system in sub-second time.

---

## The Vision

Not just automation — a **personal command language** where short inputs + live context = multi-domain parallel execution.

```
You say something short
        │
        ▼
  Intent Parser (understands YOU)
        │
        ▼
  Context Layer ── location, focus_mode, time, status
        │
        ▼
  Intent → N subjects (parallel NATS publish)
        │
   ┌────┼────┬────┬────┐
   ▼    ▼    ▼    ▼    ▼
  HA  Todo  Lock Email Notify
        │
        ▼
   Results aggregated → sub-second response
```

---

## NATS Subject Schema

### Top-level Namespace

```
cmd.ha.*          → Home Assistant commands
cmd.ha.lights.*   → Light controls
cmd.ha.lock.*     → Lock controls
cmd.ha.climate.*  → Thermostat / HVAC
cmd.ha.scene.*    → Scene activations
cmd.ha.switch.*   → Switch toggles
cmd.ha.media.*    → Media player

cmd.todoist.*     → Todoist task management
cmd.email.*       → Fastmail / email actions
cmd.notify.*      → Notifications (Telegram, push, TTS)
cmd.calendar.*    → Calendar actions

ctx.set           → Update Zeus Context KV store
ctx.get           → Read Zeus Context KV store

agent.prompt.*    → Dispatch to named agent
```

### Example: `cmd.ha.lights.on`

```
cmd  .  ha  .  lights  .  on
 │       │       │         │
 │       │       │         └── Action (on / off / set / toggle)
 │       │       └──────────── Domain (lights / lock / climate / scene)
 │       └──────────────────── Service namespace (Home Assistant)
 └──────────────────────────── Category (command)
```

---

## Zeus Context Layer

Context is always injected. Zeus never repeats where he is.

```bash
nats kv get zeus-context location      # home | office | driving | travel
nats kv get zeus-context focus_mode    # deep-work | active | rest | unavailable
nats kv get zeus-context status        # online | offline | dnd
nats kv get zeus-context time_of_day   # morning | day | evening | night
```

Same input → different output based on context:

```
"I'm leaving"

  IF location = home:
    cmd.ha.lights.off        (whole house)
    cmd.ha.lock.lock         (front door)
    cmd.ha.climate.away      (thermostat)
    ctx.set location=driving

  IF location = office:
    cmd.ha.lights.off        (office zone)
    cmd.todoist.create       ("EOD wrap")
    ctx.set location=driving
```

---

## Personal Command Vocabulary (Zeus's Shorthand)

| Command | Context | Expands To |
|---|---|---|
| `"wrap"` | any | EOD: HA off, lock, tomorrow tasks surfaced |
| `"lock in"` | home | Focus lighting + DND + surface today's tasks |
| `"I'm up"` | morning | Lights on + briefing + weather + calendar |
| `"wheels up"` | home | Lock + lights off + thermostat away + location=travel |
| `"clean slate"` | any | Close tasks, clear notifications, reset context |
| `"fire it"` | pending action | Execute last queued intent |
| `"hold"` | any | Pause automations, DND, status=unavailable |
| `"game time"` | evening/home | Cinema scene + TV on + lights dim/colored |

---

## Parallel Execution Model

All subjects fire simultaneously — no sequential waiting:

```
nats pub cmd.ha.lights.set    {...}  ─┐
nats pub cmd.ha.lock.lock     {...}  ─┤── all fire at once (~100µs per publish)
nats pub cmd.email.send       {...}  ─┤   total wall-clock: sub-second
nats pub cmd.todoist.create   {...}  ─┘
```

Reply aggregation waits max **800ms** — collects all results, reports status.

---

## Domain Workers (Always-On Subscribers)

| Worker | Subscribes To | Action |
|---|---|---|
| `ha-worker` | `cmd.ha.>` | HA REST API calls |
| `todoist-worker` | `cmd.todoist.>` | Todoist REST API |
| `email-worker` | `cmd.email.>` | Fastmail SMTP/API |
| `lock-worker` | `cmd.ha.lock.>` | HA lock entities |
| `notify-worker` | `cmd.notify.>` | Telegram / push / TTS |
| `context-worker` | `ctx.set` | Updates zeus-context KV |

---

## Architecture: 5 Layers

### Layer 1 — Intent Parser
Fast classifier that maps Zeus's shorthand → subject array + payloads.

```
Input:  "lights down, lock up, send alex the notes"
Output: [
  { subject: "cmd.ha.lights.set",  payload: { brightness: 30 } },
  { subject: "cmd.ha.lock.lock",   payload: { entity: "front_door" } },
  { subject: "cmd.email.send",     payload: { to: "alex", attach: "last_notes" } }
]
```

### Layer 2 — Context Injector
Auto-appends location / focus_mode / time to every intent before dispatch.

### Layer 3 — NATS Subject Router
Publishes all subjects simultaneously. ~100µs per message.

### Layer 4 — Domain Workers
Persistent NATS subscribers, one per domain. Stateless. Always on.

### Layer 5 — Reply Aggregator
Collects parallel replies (max 800ms window), surfaces status to Zeus.

---

## OneMind OS vs. Other Automation Systems

| System | Capability |
|---|---|
| Alexa / Siri | Voice → single action, no context, no parallelism |
| HA Automations | Trigger → fixed action, no language, no cross-app |
| Zapier / Make | Event → sequential steps, 1–30 second delay |
| **OneMind OS** | Language + context → parallel multi-domain sub-second execution |

**The difference is NATS as the universal nervous system.** Every worker is already listening. No polling, no webhook latency, no sequential chains. Publish → everything fires at once.

---

## Build Status

| Component | Status |
|---|---|
| NATS server | ✅ Running |
| Zeus Context KV | ✅ Running |
| HA integration (MCP) | ✅ Active |
| Todoist integration (MCP) | ✅ Active |
| Intent Parser | 🔲 To build |
| Domain Workers (persistent) | 🔲 To build |
| Personal Vocabulary Registry | 🔲 To build |
| Reply Aggregator | 🔲 To build |
| Context Injector | 🔲 To build |

---

## Related Notes

- [[NATS]] — Core NATS setup and config
- [[NATS-SaaS-Worker]] — Composio/SaaS event worker
- [[Hermes-Overview]] — Legacy agent architecture
