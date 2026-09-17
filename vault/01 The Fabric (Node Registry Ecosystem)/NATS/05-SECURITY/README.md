---
title: "Security — Section Index"
section: "05-SECURITY"
type: index
tags: [nats, security, auth, accounts, tls, jwt, nkeys]
created: 2026-06-14
---

# 05 — Security

Everything about securing a NATS deployment: encryption (TLS), authentication (who are you), authorization (what may you do), multi-tenancy (accounts), and decentralized auth (JWT/NKEYS).

## Notes
- [[NATS/05-SECURITY/tls-encryption|TLS Encryption]] — server TLS, mTLS, cert rotation, WebSocket/leaf TLS
- [[NATS/05-SECURITY/authentication|Authentication]] — token, user/password, NKEY, JWT, TLS-cert auth, auth callout
- [[NATS/05-SECURITY/authorization|Authorization]] — subject-level allow/deny permissions, publish/subscribe rules
- [[NATS/05-SECURITY/accounts-multitenancy|Accounts & Multi-Tenancy]] — isolated subject spaces, exports/imports, system account
- [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS (Decentralized Auth)]] — operator/account/user JWT hierarchy, `nsc`, creds files

## The three layers
```
TLS         — is the channel encrypted? (transport)
AuthN       — who are you?               (token / user-pass / NKEY / JWT / cert)
AuthZ       — what subjects may you pub/sub? (per-user allow/deny)
Accounts    — whose subject space?       (tenant isolation + exports/imports)
```

## See also
- OneMind security model: `~/onemindos-core/docs/NATS-DOCTRINE.md`, the live accounts (CORE/FABRIC/GLOBAL/INTERFACE/SAAS/VOICE), and the SYS-user pitfall in the `nats-server` ops skill.
