---
title: "Subject Design Patterns"
section: "14-REFERENCE"
type: note
tags: [nats, subjects, design, naming, patterns, conventions, wildcards, anti-patterns]
source: https://docs.nats.io/nats-concepts/subjects
created: 2026-06-14
---

# Subject Design Patterns

Subjects are the API of your system on NATS. Good subject design makes routing, security, and scaling natural; bad design fights you forever. This is the conventions + anti-patterns reference.

See also: [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects (concepts)]] · [[NATS/05-SECURITY/accounts-multitenancy|Accounts]] · [[NATS/14-REFERENCE/glossary|Glossary]]

---

## The rules of subjects

- **Tokens** are dot-separated: `a.b.c`. Case-sensitive.
- **`*`** matches exactly one token. **`>`** matches one-or-more trailing tokens (must be last).
- **Allowed chars:** alphanumerics; avoid spaces, dots-in-a-token, and (by convention) `*`/`>`/`.`/null/whitespace in literal tokens.
- **Reserved:** subjects beginning with `$` are system/internal (`$JS.API.>`, `$SYS.>`, `$KV.>`, `$O.>`). Don't publish there.
- **Length:** keep it sane; very deep subjects cost a little matching time but mostly hurt readability.

---

## The golden pattern: general → specific, left → right

Order tokens from **broadest** to **narrowest** so wildcards are useful:

```
<domain>.<entity>.<id>.<event>
orders.us.12345.created
sensors.barn.temp-01.reading
```

This lets subscribers express exactly the scope they want:

```
orders.>                    # everything orders
orders.us.>                 # all US orders
orders.*.12345.>            # one order, any region
orders.*.*.created          # all "created" events, any order
```

If you put the variable/high-cardinality token first (`12345.orders.created`), wildcards become useless and security rules can't generalize. **Hierarchy first, identity in the middle, verb last.**

---

## Hierarchy for routing & security

Subject structure *is* your authorization model — permissions are subject patterns:

```text
permissions: {
  publish:   ["orders.us.>"]      # this app may only publish US orders
  subscribe: ["orders.us.*.created", "_INBOX.>"]
}
```

So design the hierarchy with **security boundaries** in mind: group by what should be allowed/denied together. Account exports/imports also work on subject prefixes. See [[NATS/05-SECURITY/authorization|Authorization]].

---

## Verbs / event types as the last token

Put the action at the end so you can subscribe per-action across all entities:

```
account.*.created
account.*.updated
account.*.deleted
account.*.>            # all account lifecycle events
```

Pair with **queue groups** for work distribution: N workers all `sub account.*.created --queue creators` → load-balanced.

---

## Request-reply subjects

- Service subjects: `svc.<service>.<endpoint>` (e.g. `svc.billing.charge`).
- Replies go to auto-generated `_INBOX.<nuid>` — always allow `subscribe: ["_INBOX.>"]` for requesters.
- With the Services API, group endpoints under a service prefix for discovery/stats. See [[NATS/04-SERVICES-MICRO/services-micro|Services API]].

---

## JetStream subject capture

A stream binds a set of subjects; design so one stream captures a coherent subject tree:

```
stream ORDERS  ← subjects: ["orders.>"]
stream SENSORS ← subjects: ["sensors.>"]
```

Avoid overlapping streams capturing the same subject (ambiguous storage). Use **filtered consumers** to slice a broad stream (`--filter orders.us.>`) rather than many narrow streams.

---

## Anti-patterns (don't)

| Anti-pattern | Why it hurts | Do instead |
|--------------|-------------|------------|
| High-cardinality token **first** (`<uuid>.event`) | wildcards & perms can't generalize | `event.<uuid>` (id later) |
| **Encoding data** in subjects (`order.total.4999.usd`) | explodes subject space, leaks data into routing | put data in the **payload/headers** |
| One mega-subject for everything (`events`) | no filtering, no per-scope security | hierarchical tokens |
| Subjects that **change meaning** by position | impossible to reason about | fixed positional schema |
| Publishing to `$...` reserved subjects | collides with system | use your own namespace |
| Unbounded dynamic tokens (per-request literal) | sublist churn | bounded hierarchy + payload ids |
| Mixing separators (`orders-us.new`) | inconsistent matching | dots only between tokens |

---

## A reference schema

```
<org>.<domain>.<entity>.<id>.<event>

acme.orders.order.12345.created
acme.iot.sensor.barn-temp-01.reading
acme.svc.billing.charge            # request-reply endpoint
```

- `org` — tenant/account boundary (often the account itself, so omit if accounts already isolate).
- `domain` — bounded context (orders, iot, billing).
- `entity` + `id` — the thing and which one.
- `event` — the verb, last, for cross-entity action subscriptions.

---

## OneMind Relevance (this is the fleet's subject law)

Recommended fleet schema:

```
omos.<node>.<domain>.<event>[.<id>]

omos.barn.telemetry.reading
omos.barn.telemetry.reading.temp-01
omos.vehicle-1.command.move
omos.house.agent.heartbeat
omos.>                     # whole fleet
omos.barn.>                # one site
omos.*.telemetry.>         # all telemetry, any node
omos.*.command.>           # all commands (lock down publish!)
```

- **`node` second** → per-node security + per-site subscriptions (a leaf only publishes its own `omos.<node>.>`).
- **`domain` third** (`telemetry`/`command`/`agent`/`ledger`) → maps to streams and to Core-vs-JetStream tier (telemetry=Core/drop-ok, command/ledger=JetStream/durable).
- **`event` + optional `id` last** → cross-node action subscriptions and queue-group workers.
- Lock **`publish: omos.<node>.>`** per node via JWT so a compromised edge node can't impersonate the fleet. Commands flow hub→node on tightly-scoped subjects.

This schema *is* a design decision — capture it in the SO spine before services proliferate; retrofitting subject law across hundreds of nodes is painful.

## Sources
- Subjects — https://docs.nats.io/nats-concepts/subjects
- Subject-based messaging — https://docs.nats.io/nats-concepts/what-is-nats/walkthrough_setup
- ADR: subject mapping — https://github.com/nats-io/nats-architecture-and-design
