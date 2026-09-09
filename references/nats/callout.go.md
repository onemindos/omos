# callout.go — NATS Auth Callout Service Library

**Repo:** https://github.com/synadia-io/callout.go
**Publisher:** Synadia (makers of NATS)
**License:** Apache 2.0
**Install:** `go get github.com/synadia-io/callout.go`

## What it is

Official Synadia Go library for implementing a NATS **auth callout service**. When a client connects to NATS, the server calls out to your auth service to make the authentication decision dynamically — instead of static credential files.

## Why it's relevant to OneMind

OneMind has an open task: **"🟡 NATS Auth/ACL — Per-agent subject permissions"** (SO). Dynamic auth callout is the right pattern for this:

- Each agent/service gets its own subject permissions at connect time
- The auth service decides what `cmd.*`, `evt.*`, `tel.*` subjects each actor can publish/subscribe to
- No static creds files to manage per agent — permissions are generated dynamically based on actor identity

## Typical flow

```
Agent connects to NATS
    → NATS server calls out to auth service (your callout.go implementation)
    → Auth service returns signed JWT with subject permissions
    → Agent is admitted with scoped access
```

## When to use

Build this when tackling the per-agent ACL task. Pairs with `jwt-auth-builder.go` to generate the scoped JWTs.

## Notes

- Works with NATS account model (CORE/FABRIC/INTERFACE split in OneMind)
- Supports delegated signing keys for zero-trust patterns
- Examples in repo: delegated auth, dynamic accounts, default sentinel
