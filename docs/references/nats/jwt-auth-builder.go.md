# jwt-auth-builder.go — NATS JWT Auth Builder

**Repo:** https://github.com/synadia-io/jwt-auth-builder.go
**Publisher:** Synadia (makers of NATS)
**License:** Apache 2.0
**Install:** `go get github.com/synadia-io/jwt-auth-builder.go`

## What it is

Official Synadia Go library for programmatically creating and managing NATS JWTs — operators, accounts, users, signing keys, revocations, imports/exports.

## Why it's relevant to OneMind

Pairs directly with `callout.go` when building the OneMind per-agent NATS auth system. Use it to:

- Generate scoped user JWTs with subject-level permissions per agent
- Manage signing keys for the CORE/FABRIC/INTERFACE account split
- Programmatically issue/revoke credentials without `nsc` CLI
- Store JWTs in JetStream KV (has a `providers/kv` provider built in)

## Key capabilities

| Feature | Use |
|---------|-----|
| Operator/Account/User JWTs | Full credential hierarchy for OneMind accounts |
| Signing keys | Rotate keys without changing operator/account identity |
| Subject permissions | Scope each agent to only its allowed `cmd.*`/`tel.*`/`evt.*` subjects |
| Revocations | Immediately revoke a compromised agent's access |
| KV provider | Store issued JWTs in JetStream KV — no external database needed |
| NSC provider | Bridge to existing nsc-managed credentials if needed |

## When to use

Build this alongside `callout.go` when tackling **"🟡 NATS Auth/ACL — Per-agent subject permissions"** in SO.

## Notes

- The `providers/kv` backend means credentials live entirely in JetStream — no extra infra
- Actively maintained by Synadia, tracks NATS JWT v2 spec
