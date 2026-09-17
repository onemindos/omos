---
title: "JWT & NKEYS — Decentralized Authentication"
section: "05-SECURITY"
type: note
tags: [nats, security, jwt, nkeys, nsc, operator, decentralized, creds]
source: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/jwt
created: 2026-06-14
---

# JWT & NKEYS — Decentralized Authentication

The production-grade, scalable identity model. Identity is a **signed three-tier hierarchy** — Operator → Account → User — encoded as JWTs and signed with NKEYS. The server trusts only the **operator public key** and verifies the whole chain cryptographically. **Users can be created offline, with no server config change.**

See also: [[NATS/05-SECURITY/authentication|Authentication]] · [[NATS/05-SECURITY/accounts-multitenancy|Accounts]]

---

## NKEYS — the key primitive

**NKEYS** are NATS's encoding of **ed25519** keypairs. Each entity (operator, account, user) has a keypair. Keys are prefixed by role so you can tell them apart at a glance:

| Prefix | Role |
|--------|------|
| **O** | Operator public key |
| **A** | Account public key |
| **U** | User public key |
| **S** | **Seed** (the private key — `SO…`, `SA…`, `SU…`) |
| **N** | Server |
| **C** | Cluster |

Generate with **`nk`**: `nk -gen user -pubout`. The **seed** (`S…`) is the secret; the public key (`O/A/U…`) is shareable. Signing a server nonce with the seed proves identity without sending the secret.

---

## The Trust Hierarchy

```
        OPERATOR  (root of trust; the server is configured with its public key)
           │ signs
           ▼
        ACCOUNT   (a tenant: owns subjects, JetStream, limits, exports/imports)
           │ signs
           ▼
        USER      (an identity within the account; carries permissions)
```

- The **operator** signs **account** JWTs.
- Each **account** signs its own **user** JWTs (and embeds the user's permissions + limits).
- A NATS server configured with the **operator JWT/public key** can verify any account and user in the chain — so **new users are minted offline** by the account holder, never touching server config. This is the decentralization: issuance is delegated down the tree.

---

## `nsc` — the management tool

**`nsc`** ("NATS Security Configuration") creates and manages the whole hierarchy:

```shell
nsc add operator OneMind                       # create operator
nsc add account CORE                            # create account under operator
nsc add user --account CORE zeus                # create user (mints a JWT + seed)
nsc generate creds --account CORE --name zeus > zeus.creds   # export creds file
nsc edit user zeus --allow-pub "core.>" --allow-sub "core.>,_INBOX.>"  # set perms
nsc push -A                                     # push account JWTs to a resolver
nsc list keys                                   # inspect
```

`nsc` stores JWTs and (optionally) seeds locally and can **push** account JWTs to an **account resolver** the server reads.

---

## The `.creds` File

A **credentials file** bundles a **user JWT** + the user's **NKEY seed**:

```
-----BEGIN NATS USER JWT-----
eyJ0eXAiOiJKV1Q...
------END NATS USER JWT------
-----BEGIN USER NKEY SEED-----
SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ...
------END USER NKEY SEED------
```

The client points at it and the library handles the challenge/response:

```go
nc, _ := nats.Connect(url, nats.UserCredentials("zeus.creds"))
```
```python
nc = await nats.connect(url, user_credentials="zeus.creds")
```

> Protect `.creds` like a private key — it *contains* the seed. `nsc generate creds` produces it; store in 1Password, never in git.

---

## Account Resolvers — how the server gets account JWTs

The operator JWT alone lets the server verify signatures, but it needs each **account JWT** to know that account's config/limits. Three resolver modes:

| Resolver | How accounts are loaded |
|----------|-------------------------|
| **MEMORY** | Accounts preloaded in config (`resolver_preload`) — static, small deployments |
| **URL** | Server fetches account JWTs from an HTTP endpoint |
| **NATS-based (full)** | Accounts stored *in NATS itself* (the built-in account server); `nsc push` updates them live — the scalable production choice |

```text
operator: "/etc/nats/operator.jwt"
resolver: {
  type: full
  dir: "/etc/nats/jwt"
  allow_delete: true
  interval: "2m"
}
system_account: "ADXU...SYS"
```

---

## Why JWT for a Fleet

| Benefit | Detail |
|---------|--------|
| **Offline user issuance** | Account holders mint users without server access or restarts |
| **Embedded permissions + limits** | Each user JWT carries its own authz + connection/data limits |
| **Revocation** | Accounts can revoke users; resolvers propagate |
| **Multi-tenant by construction** | Account = tenant, with its own signing authority |
| **No shared secrets on the wire** | NKEY challenge/response; seed never sent |
| **Rotation** | Re-issue JWTs/keys without touching the server's trust anchor (operator) |

---

## OneMind Relevance

OneMind currently runs simpler auth (token / user-password per account over Tailscale). The **JWT + nsc** model is the documented growth path for "hundreds of nodes": an **operator** as root of trust, an **account per domain/tenant**, and **per-device user creds** minted offline as the fleet grows — without ever editing and restarting the hub server for each new node. This directly serves the sovereignty/scale doctrine. When the fleet outgrows static creds, `nsc` + a full NATS-based resolver is the move.

## Sources
- JWT auth — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/jwt
- NKEYS — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth
- `nsc` tooling — https://docs.nats.io/using-nats/nats-tools/nsc
- Account resolvers — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/jwt/resolver
- `nk` / `nsc` repos — https://github.com/nats-io/nkeys · https://github.com/nats-io/nsc
