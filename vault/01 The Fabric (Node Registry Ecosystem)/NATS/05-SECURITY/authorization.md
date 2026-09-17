---
title: "Authorization"
section: "05-SECURITY"
type: note
tags: [nats, security, authorization, permissions, allow, deny]
source: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization
created: 2026-06-14
---

# Authorization

*What may you do?* Once authenticated, a user's **permissions** decide which subjects they can **publish** to and **subscribe** from — subject-level allow/deny.

See also: [[NATS/05-SECURITY/authentication|Authentication]] · [[NATS/05-SECURITY/accounts-multitenancy|Accounts]] · [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects & Wildcards]]

---

## Permission Model

Each user has `permissions` with `publish` and `subscribe` rules, each holding `allow` and/or `deny` subject lists (wildcards supported):

```text
authorization {
  users: [
    {
      user: "telemetry-writer", password: "..."
      permissions: {
        publish:   { allow: ["TEL.>"] }                # may publish only under TEL.*
        subscribe: { deny:  ["secrets.>"] }            # may subscribe to anything except secrets.*
      }
    },
    {
      user: "readonly"
      permissions: {
        publish:   { deny: [">"] }                      # cannot publish anything
        subscribe: { allow: ["events.>", "_INBOX.>"] }  # read events + receive replies
      }
    }
  ]
}
```

- **allow** lists what's permitted; **deny** overrides allow (deny wins on conflict).
- Use `>` (everything) and `*` (one token) [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|wildcards]].
- Omitting `publish`/`subscribe` = allow all for that direction (so lock down explicitly).

---

## Don't Forget `_INBOX` for Request-Reply

A user that makes **requests** must be allowed to **subscribe** to its reply inbox (`_INBOX.>`), and a user that **responds** must be allowed to **publish** to `_INBOX.>` (or the specific reply subject). Forgetting this silently breaks request-reply:

```text
permissions: {
  publish:   { allow: ["math.>", "_INBOX.>"] }   # responder publishes replies
  subscribe: { allow: ["math.>", "_INBOX.>"] }   # requester receives replies
}
```

---

## Allow-Responses (least privilege for services)

For a pure responder, instead of opening `_INBOX.>` broadly, use **`allow_responses`** — the server lets the user publish a reply **only** to the reply-subject of a message it just received, optionally bounded by count and time:

```text
permissions: {
  subscribe: { allow: ["service.requests.>"] }
  allow_responses: { max: 1, expires: "2s" }    # may reply once, within 2s
}
```
This is the tightest req-reply permission — a service can answer requests but cannot publish anywhere else.

---

## Queue Permissions

You can restrict which **queue groups** a user may join (so a worker can only join its intended group), expressed in the subscribe permission with a queue qualifier where supported by the config/version.

---

## Where Permissions Live

| Auth model | Permissions defined in |
|------------|------------------------|
| token | n/a (single identity) |
| user/password, nkey (config) | inline `permissions` per user |
| **JWT** | encoded **inside the user JWT** (signed by the account) — see [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]] |
| accounts | account boundaries are a *separate, stronger* isolation than per-user permissions |

> **Permissions vs Accounts:** per-user permissions restrict a user *within* its account's subject space. **Accounts** (next note) are a hard wall — subjects in account A are invisible to account B entirely, regardless of permissions. Use permissions for intra-account least-privilege; use accounts for tenant isolation.

---

## Reload

Authorization changes in the config file apply on **`nats-server --signal reload`** (SIGHUP) — *for config-based users/permissions*. (Account/system-account topology changes need a full restart; see the ops pitfalls.)

## Sources
- Authorization — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization
- allow_responses — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/authorization#allow-responses-permissions
- Subjects & wildcards — https://docs.nats.io/nats-concepts/subjects
