---
title: "05 Stack — Arsenal Ecosystem (REDIRECT)"
type: redirect
status: migrated
created: 2026-01-13
modified: 2026-06-05
owner: zeus
tags: [type/redirect, domain/so, status/migrated]
ai_summary: "REDIRECT STUB. The 173-file arsenal (every tool/platform/service) was collapsed into The Fabric on 2026-06-05. Services now live as nodes at 01 The Fabric/nodes/digital/services-registry/. The arsenal is no longer a folder — it is a VIEW: query category:digital AND class:service."
---

# 05 Stack — Arsenal Ecosystem → **MIGRATED**

> 🔀 **This folder is a redirect stub.** On **2026-06-05** the entire 173-file arsenal was collapsed into **The Fabric** as part of the OneMind NATS-OS unification.

## Where it went

```
OLD:  00-24 SO/01 The Fabric (Node Registry)/<domain>/...
NEW:  00-24 SO/01 The Fabric (Node Registry)/nodes/digital/services-registry/<domain>/...
```

All 10 domain groupings were preserved (`ai/ browser/ business/ communication/ creative/ development/ hardware/ infrastructure/ iot/ personal/`). Domain is now a navigable **property**, not a bus token.

## The arsenal is no longer a folder — it is a VIEW

Every platform became a **node** in The Fabric with class `service`. "The arsenal" is now a query:

```
category: digital  AND  class: service
```

Bus address pattern: `<plane>.service.<id>.<path>` — e.g.
- `cmd.service.todoist.create_task`
- `qry.service.github.list_issues`
- `tel.service.nats.connections`

## Why

A service IS a node you can write to (`cmd`) or read from (`qry`). Keeping a separate "Stack" folder duplicated the registry. One namespace: the vault registry and the NATS namespace are now the same namespace.

## What stayed here

- This redirect stub (`05 Stack.md`)
- `archive/` — dead/replaced platforms (pending disposition — flagged for discussion, not yet moved)

---

> Looking for a specific platform? It's at:
> `01 The Fabric (Node Registry)/nodes/digital/services-registry/<domain>/<platform>.md`
