---
title: "Subject Grammar — The Locked Spine"
domain: SO
type: law
status: LOCKED
version: 1.0.0
locked: 2026-06-05
owner: zeus
authority: constitutional
tags: [type/law, domain/so, status/locked, nats, grammar, spine]
ai_summary: "The single invariant of One Mind OS — the addressing grammar for every node on the NATS bus. LOCKED 2026-06-05. Subject law is plane.class.id.path with 5 planes (evt/tel/cmd/qry/gov). Changing this rewrites every subject in the system; it is deliberately frozen. Everything else in the architecture conforms to this; this conforms to nothing."
---

# Subject Grammar — The Locked Spine

> 🔒 **STATUS: LOCKED — v1.0.0 — 2026-06-05.**
> This is the one invariant. Every node, agent, rule, and stream conforms to it. It conforms to nothing.
> Changing this grammar rewrites every subject in the system. It is frozen on purpose. Amendments require an explicit version bump and a migration plan — never a casual edit.

---

## The Law

```
<plane>.<class>.<id>.<path>
```

Four segments. Read left to right: **what kind of message · what kind of thing · which one · the detail.**

Example: `cmd.drone.skyhawk-1.land` = a command, to a drone, the one called skyhawk-1, telling it to land.

---

## Segment 1 — `plane` (the verb: what kind of message this is)

Five planes. Exactly five. This is the most load-bearing decision in the system — it's how the autonomic layer routes and governs without an LLM ever reading the message.

| Plane | Meaning | Governed | Persistence | Example |
|-------|---------|----------|-------------|---------|
| **`evt`** | Something **happened** (past tense, factual, immutable) | audited | JetStream + ClickHouse, 5yr | `evt.service.todoist.task_completed` |
| **`tel`** | A **measurement** (a reading, a sample, ephemeral) | soft | ephemeral / short TTL | `tel.human.zeus.heartrate` |
| **`cmd`** | **Mutate** something (an order, has effects) | **HARD — ACL enforced** | logged as `evt` on success | `cmd.device.plug-a886.turn_on` |
| **`qry`** | **Read** something (request/reply, no effects) | soft | not persisted | `qry.account.mercury.balance` |
| **`gov`** | **Meta** — the system governing itself (escalation, audit, policy) | self-audited | JetStream, permanent | `gov.escalate.security` |

**Why this split:** the bus carries thousands of messages/sec and 99.9% of them are `tel` and `evt` that need a *rule*, not a *thought*. The plane is the first routing decision — a deterministic daemon (`subscriber.py`, the autonomic layer) reads the plane and acts in sub-milliseconds at zero token cost. Only `gov.escalate.>` wakes the cortex (Legacy). **The plane is what makes governance deterministic instead of theatrical.**

**The hard rule:** `cmd` is the only plane that mutates state, so `cmd` is the only plane governed by hard ACLs at the bus level. You cannot fake a command into an `evt`. Mutation has exactly one door, and that door has a lock.

---

## Segment 2 — `class` (what kind of thing)

The class is the node's *kind* — drawn from The Fabric's class list (the 18 node classes + control classes). The class token IS the vault class IS the bus class. One word, three meanings, zero translation.

Examples: `human`, `pet`, `device`, `sensor`, `robot`, `drone`, `vehicle`, `machine`, `agent`, `service`, `org`, `account`, `holding`, `location`, `land`, `crop`, `mode`, `fleet`.

The class list itself is **versioned vocabulary** (additive, cheap to extend) — it is NOT part of the locked grammar. You add a class by adding a folder + a token. You never change the *shape* `plane.class.id.path`.

---

## Segment 3 — `id` (which specific one)

The unique identifier of the individual node.

### id naming law

- **`lowercase-hyphen`** only. `main-house`, `skyhawk-1`, `plug-a886`.
- **NEVER a dot.** A dot is a subject delimiter — an id with a dot in it shatters the address. This is the one character that breaks everything.
- **Singletons = the name:** `legacy`, `mbp-hub`, `mercury`, `todoist`.
- **Instances = name + suffix:** `skyhawk-1`, `car-1`, `meter-7616`, `plug-a886`.

The id in the subject IS the id in the registry frontmatter IS the folder/file identity. Same token everywhere.

---

## Segment 4 — `path` (the detail)

Everything after the id is free hierarchy — the specific signal, command verb, sub-resource, or reading. Use as many dots as the structure needs.

- `tel.human.zeus.heartrate`
- `tel.human.zeus.sleep.deep`
- `cmd.service.github.issues.create`
- `qry.location.main-house.presence`
- `*.account.mercury.>` (wildcard subscription to everything about Mercury)

---

## The 3 Modeling Rules (how to decide a subject without guessing)

These resolve 95% of "what subject does this get?" questions. They are doctrine, frozen alongside the grammar.

### Rule 1 — Telemetry publishes under what it is ABOUT, not what produced it.

The Oura ring is a `device`. But Zeus's heart rate is about *Zeus*, so it publishes under `tel.human.zeus.heartrate` — not `tel.device.oura.heartrate`. The Bouncie tracker is a `device`; the car's location is about the *car*: `tel.vehicle.car-1.location`. **The producer is an implementation detail. The subject names the subject.**

### Rule 2 — Identity goes in the subject. Properties go in the registry.

The subject carries only `class` + `id`. Everything else — location, owner, institution, IP, specs, access level, pillar — lives as frontmatter *properties* in the node's vault file. `location` is a property of a sensor, never a routing token. `pillar` is a ClickHouse column, never a subject segment. **If you're tempted to put it in the subject, ask: is this WHICH node, or a FACT about the node? Only "which" goes in the subject.**

### Rule 3 — Work is EVENTS, not classes.

A task, a goal, a payment, a message, a calendar entry — these are not nodes. They are things that *happen to* nodes. A completed task is `evt.service.todoist.task_completed`. A payment is `evt.account.mercury.payment_sent`. **Nouns that persist and have telemetry are nodes. Verbs that occur are events.** This is what dissolved the old Entity/Node split: the "entity" work-types were always events.

---

## What is LOCKED vs what is OPEN

| Locked (frozen — version bump required) | Open (additive — extend freely) |
|---|---|
| The 4-segment shape `plane.class.id.path` | The class vocabulary (add `crop`, `livestock`, etc.) |
| The 5 planes and their semantics | New ids (every new node) |
| `cmd` = the only mutating, hard-governed plane | New paths under any node |
| The id naming law (lowercase-hyphen, no dots) | New modeling examples |
| The 3 modeling rules | — |

**The test of a good design from here:** every new capability must express itself in this grammar without bending it. If something can't be said in `plane.class.id.path`, the design is wrong — not the grammar.

---

## Worked Examples (the grammar in anger)

```
cmd.device.plug-a886.turn_on            turn on a smart plug
tel.device.plug-a886.power              its live power draw
qry.location.main-house.presence        who/what is home right now
cmd.drone.skyhawk-1.land                land a drone
tel.drone.skyhawk-1.battery             its battery level
tel.human.zeus.heartrate                Zeus's heart rate (from Oura, a device)
tel.vehicle.car-1.location              car GPS (from Bouncie, a device)
cmd.service.todoist.task.create         create a task
evt.service.todoist.task_completed      a task got completed (fact)
qry.account.mercury.balance             read a bank balance
evt.account.mercury.payment_sent        a payment happened (fact)
*.org.onemind.>                         everything about the OneMind org
cmd.mode.incognito                      flip the whole system to incognito
cmd.fleet.shutdown                      command an entire fleet at once
gov.escalate.security                   wake the cortex — security judgment needed
```

---

> **One grammar. Every noun. One mind.**
> The vault registry and the NATS namespace are the same namespace, and this is the law that makes them one.
