---
title: "NUID & NKEYS"
section: "10-PROTOCOL-INTERNALS"
type: note
tags: [nats, nuid, nkeys, ed25519, crypto, identity, inbox, internals]
source: https://github.com/nats-io/nuid
created: 2026-06-14
---

# NUID & NKEYS

Two small but foundational primitives. **NUID** generates the fast unique identifiers NATS uses everywhere internally (inboxes, client IDs, dedup tokens). **NKEYS** is the Ed25519-based identity system that underpins decentralized JWT auth.

See also: [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS (Security)]] · [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]] · [[NATS/05-SECURITY/jwt-nkeys-decentralized|nsc & nk]]

---

## NUID — NATS Unique Identifier

A **NUID** is a 22-character, URL-safe, collision-resistant identifier optimized for **speed**, not cryptographic randomness. Repo: `nats-io/nuid`.

How it works:
- A NUID has a **fixed random prefix** (12 base-62 chars, ~71 bits of entropy) generated once per instance, plus a **sequential suffix** (10 base-62 chars) that increments by a random step.
- When the sequence space is exhausted (or periodically), the prefix is re-randomized.
- Generating the next NUID is essentially a counter increment + encode — **no locking on a global RNG per call**, so it's extremely fast (tens of millions/sec).

Where NATS uses NUIDs:
- **Reply inboxes** — `_INBOX.<NUID>` for every request-reply round-trip. The uniqueness guarantees one requester's reply doesn't land in another's subscription. See [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]].
- **Client identifiers**, subscription IDs, and other internal handles.
- A natural choice for **`Nats-Msg-Id`** dedup tokens when you don't have a domain key. See [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]].

> NUID is for *uniqueness at speed*, not unpredictability. Don't use a raw NUID as a security token — use NKEYS/JWT for anything that must be unforgeable.

```go
import "github.com/nats-io/nuid"
id := nuid.Next()            // e.g. "1k9aX…q3" (22 chars)
inbox := "_INBOX." + nuid.Next()
```

---

## NKEYS — Ed25519 identity

**NKEYS** are NATS's public-key identity system, built on **Ed25519** (fast, modern elliptic-curve signatures). Repo: `nats-io/nkeys`. Each NKEY is an Ed25519 keypair encoded with a **role prefix** and a CRC, so the *type* of an identity is visible in its first character.

| Prefix | Role | Used for |
|--------|------|----------|
| `O` | **Operator** | Root of trust for a deployment (signs accounts) |
| `A` | **Account** | A tenant identity (signs users); owns subject space |
| `U` | **User** | An end identity (a client/agent) that connects |
| `S` | **Seed** | The *private* key (the secret you guard) — `SO…`, `SA…`, `SU…` |
| `N` | **Server** | Server identity |
| `C` | **Cluster** | Cluster identity |

A **seed** (`S…`) is the private key; from it you derive the **public** key (`O…`/`A…`/`U…`). You never share the seed.

```bash
# Generate a user NKEY with nk
nk -gen user -pubout            # prints the public U… key
nk -gen user                    # prints the seed SU… (private — guard it)
```

### Challenge-response auth (no shared secret on the wire)
NKEYS enable authentication where the **private key never leaves the client**:
1. Server's `INFO` includes a random **nonce**.
2. Client **signs the nonce** with its seed and sends the signature + its public NKEY in `CONNECT`.
3. Server verifies the signature against the public key it trusts. No password is transmitted.

### NKEYS → JWT → decentralized auth
NKEYS are the signing keys behind **JWT-based decentralized auth**: an **Operator** NKEY signs **Account** JWTs; an Account NKEY signs **User** JWTs. The server is configured to trust the Operator's public key and can then validate the whole chain **offline** — no central auth database. This is the model `nsc` manages and what the fleet uses for agent identity. Full treatment in [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS (Security)]].

```
Operator (O / SO)  ──signs──▶  Account JWT (A / SA)  ──signs──▶  User JWT (U / SU)
        │                                                              │
   trusted by server                                          presented at CONNECT
```

---

## NUID vs NKEYS — don't confuse them

| | NUID | NKEYS |
|--|------|-------|
| Purpose | Fast unique **IDs** | Cryptographic **identity** |
| Security | Not a secret | Ed25519 keypair (seed is secret) |
| Speed | ~tens of M/sec | Sign/verify per auth |
| Used for | Inboxes, CIDs, dedup tokens | Auth, JWT signing, server/account/user identity |

---

## OneMind Relevance
- Every agent/service in the fleet gets a **User NKEY/JWT** identity (managed via `nsc`), trusted up a chain to our **Operator** key — that's how "hundreds of nodes" authenticate without a central password store and work offline at the edge.
- `_INBOX.<NUID>` reply subjects are behind every request-reply call our workers make.
- Seeds (`S…`) are **secrets** → they live in 1Password (`op://…`), never in markdown or repo. See `~/onemindos-core/docs/NATS-DOCTRINE.md`.

## Sources
- NUID — https://github.com/nats-io/nuid
- NKEYS — https://github.com/nats-io/nkeys
- JWT — https://github.com/nats-io/jwt
- Decentralized auth — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/jwt
