# Resgate — Real-Time NATS API Gateway

**Repo:** https://github.com/resgateio/resgate
**Module:** `github.com/resgateio/resgate`
**Language:** Go
**License:** MIT
**Docs:** https://resgate.io

## What it is

Resgate is a real-time API gateway that sits between NATS and browser clients. It syncs NATS resources to connected clients via WebSocket using the **RES protocol** — clients subscribe to NATS resources and receive live updates automatically. No polling, no custom WebSocket server to build.

```
Browser/App                  Resgate                    NATS Bus
    │                           │                           │
    ├── subscribe(resource) ───►│                           │
    │                           ├── get.{resource} ────────►│
    │                           │◄── response ──────────────┤
    │◄── initial state ─────────┤                           │
    │                           │   event.{resource}.>  ◄──►│
    │◄── live change events ────┤  (subscribed)             │
```

## Why it's relevant to OneMind

| Use case | How |
|----------|-----|
| `omos-ui` real-time updates | Browser subscribes to NATS resources (agent roster, stream stats, entity state) — updates push automatically, no polling |
| Entity state live feed | `ent.{type}.{id}.state` events on NATS → pushed to all subscribed browser clients instantly |
| Bus topology live view | `$SYS.*` subjects expose cluster health — resgate bridges this to the browser |
| Flutter cockpit sync | RES client libraries available for Dart — same protocol, same NATS backend |

## RES Protocol subjects

Resgate maps NATS subjects to resources:

```
get.{service}.{resource}          → fetch resource state
call.{service}.{resource}.{method} → call a resource method
auth.{service}.{resource}.{method} → authenticate
event.{service}.{resource}.>       → live change events
```

## Client libraries

| Platform | Library |
|----------|---------|
| JavaScript/TypeScript | `resclient` — https://github.com/resgateio/resclient |
| Dart/Flutter | `resclient-dart` — https://github.com/resgateio/resclient-dart |
| Go | Native service SDK — `github.com/resgateio/go-res` |

## Run with Docker

```bash
docker run -p 8080:8080 resgateio/resgate --nats nats://nats:4222
```

Or point at the OneMind NATS cluster:

```bash
docker run -p 8080:8080 resgateio/resgate \
  --nats nats://nats.nats.svc.cluster.local:4222 \
  --apikey your-api-key
```

## Key docs in the local clone

| File | Content |
|------|---------|
| `docs/gateway-flow.md` | Full data flow: cache, subscribe, change events |
| `docs/res-client-protocol.md` | Browser ↔ resgate wire protocol |
| `docs/res-protocol-v1.2-update.md` | Latest protocol changes |

## Notes

- Already cloned locally at `integrations/nats/resgate/` in omos-core local repo
- Official Docker image: `resgateio/resgate`
- Actively maintained, used in production by multiple companies
- Alternative to building a custom WebSocket bridge for `omos-ui`
