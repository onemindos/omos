---
title: "Authentication"
section: "05-SECURITY"
type: note
tags: [nats, security, authentication, token, nkey, jwt, auth-callout]
source: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro
created: 2026-06-14
---

# Authentication

*Who are you?* NATS supports five authentication mechanisms, from a single shared token to fully decentralized JWT. Plus **auth callout** for delegating to an external system.

See also: [[NATS/05-SECURITY/authorization|Authorization]] · [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]] · [[NATS/05-SECURITY/tls-encryption|TLS]]

---

## The Five Mechanisms

| Mechanism | How | Best for |
|-----------|-----|----------|
| **Token** | Single shared bearer token | Dev, simple internal trust |
| **User / Password** | Per-user credentials | Small fixed user sets |
| **NKEY** | ed25519 public-key challenge/response | Servers/services, no shared secrets on the wire |
| **JWT (decentralized)** | Signed user JWT + NKEY seed in a `.creds` file | Production, fleets, multi-tenant, dynamic users |
| **TLS client cert** | Cert subject mapped to a user | Per-device identity (see [[NATS/05-SECURITY/tls-encryption|TLS]]) |

---

## Token

```text
authorization {
  token: "s3cr3t-bus-token"
}
```
Client: `nats.Connect(url, nats.Token("s3cr3t-bus-token"))`. Simplest, weakest — one secret for everyone, no per-user authz.

---

## User / Password

```text
authorization {
  users: [
    { user: "alice", password: "pw1" }
    { user: "bob",   password: "pw2", permissions: { ... } }   # per-user authz
  ]
}
```
Passwords can be **bcrypted** (`nats server passwd` generates the hash). Supports per-user [[NATS/05-SECURITY/authorization|permissions]].

---

## NKEY (ed25519 challenge/response)

NKEYS are NATS's ed25519 key encoding. The **public key** goes in server config; the client holds the **seed** (private key). On connect, the server sends a **nonce**, the client **signs** it with its seed, and the server verifies against the public key — **the private key never crosses the wire**.

```text
authorization {
  users: [
    { nkey: "UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4" }
  ]
}
```
Client: `nats.Connect(url, nats.Nkey(pubkey, sigCB))`. Great for services — rotate by swapping keys, no passwords to leak.

---

## JWT — Decentralized Auth

The production-grade, scalable model. Identity is a **hierarchy of signed JWTs** (Operator → Account → User), and the server only needs the **operator's public key** to verify the whole chain. Users can be minted **offline** without touching server config. Full detail in [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]].

```text
# Server just trusts the operator + resolves account JWTs
operator: "/etc/nats/operator.jwt"
resolver: MEMORY        # or URL/NATS-based account resolver
resolver_preload: { ACC_PUBKEY: "<account jwt>" }
```
Client uses a **`.creds`** file (user JWT + seed):
`nats.Connect(url, nats.UserCredentials("zeus.creds"))`.

---

## Auth Callout — delegate to an external system

**Auth callout** lets the server delegate authentication to **your own service**: on connect, the server packages the client's credentials into a request and publishes it to an auth service (in the system or a dedicated account); your service validates (against LDAP, a DB, an OIDC provider, anything) and returns a **signed user JWT** granting permissions. This bridges NATS auth to any external identity source without baking users into config.

```text
authorization {
  auth_callout {
    issuer: "ABC...ISSUER_NKEY"
    account: "AUTH"
    auth_users: [ "auth_service_user" ]
  }
}
```

Use cases: dynamic users, central IdP integration, per-request fine-grained permissions, short-lived credentials.

---

## `no_auth_user` — anonymous mapping

```text
no_auth_user: "onemind"     # unauthenticated clients become user "onemind"
```
Maps **unauthenticated** connections to a specific configured user (and therefore that user's account/permissions). Convenient for local/dev or a deliberately-open subject space — **but** note the OneMind gotcha: `no_auth_user` lands clients in whatever **account** that user belongs to (GLOBAL), which is why "visibility ≠ persistence" — see [[NATS/05-SECURITY/accounts-multitenancy|Accounts]].

---

## Choosing

```
One trusted internal bus, dev?            → token
Small fixed set of users?                  → user/password (bcrypted)
Services / machine identity, no secrets?   → NKEY
Fleet, multi-tenant, dynamic, production?  → JWT (+ nsc)
Per-device hardware identity?              → TLS client cert (verify_and_map)
Must integrate an external IdP/LDAP/DB?    → auth callout
```

## Sources
- Authentication intro — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro
- Tokens / user-password — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/tokens
- NKEYS — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth
- Auth callout — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_callout
- JWT — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/jwt
